package auth

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/YorkHost-fr/yorkhost-cli/internal/api"
	"github.com/YorkHost-fr/yorkhost-cli/internal/config"
)

// fakeAS : serveur d'autorisation minimal, mêmes règles que yorkhost-api
// (PKCE S256, redirect_uri exacte, refresh à usage unique).
type fakeAS struct {
	mu        sync.Mutex
	redirects map[string]string // client_id → redirect_uri
	codes     map[string]string // code → code_challenge
	refresh   map[string]bool   // refresh → déjà utilisé
	issued    int
	revoked   []string
}

func newFakeAS(t *testing.T) (*fakeAS, *httptest.Server) {
	f := &fakeAS{redirects: map[string]string{}, codes: map[string]string{}, refresh: map[string]bool{}}
	mux := http.NewServeMux()
	mux.HandleFunc("POST /client/oauth/register", func(w http.ResponseWriter, r *http.Request) {
		var in struct {
			RedirectURIs []string `json:"redirect_uris"`
			Method       string   `json:"token_endpoint_auth_method"`
		}
		json.NewDecoder(r.Body).Decode(&in)
		if len(in.RedirectURIs) != 1 || in.Method != "none" {
			w.WriteHeader(400)
			return
		}
		f.mu.Lock()
		id := fmt.Sprintf("yhc_%d", len(f.redirects)+1)
		f.redirects[id] = in.RedirectURIs[0]
		f.mu.Unlock()
		w.WriteHeader(201)
		json.NewEncoder(w).Encode(map[string]any{"client_id": id})
	})
	mux.HandleFunc("GET /client/oauth/authorize", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		f.mu.Lock()
		defer f.mu.Unlock()
		if f.redirects[q.Get("client_id")] != q.Get("redirect_uri") || q.Get("code_challenge_method") != "S256" {
			w.WriteHeader(400)
			return
		}
		// Consentement simulé : redirection immédiate vers l'application.
		code := fmt.Sprintf("code%d", len(f.codes))
		f.codes[code] = q.Get("code_challenge")
		http.Redirect(w, r, q.Get("redirect_uri")+"?code="+code+"&state="+url.QueryEscape(q.Get("state")), http.StatusFound)
	})
	mux.HandleFunc("POST /client/oauth/token", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		f.mu.Lock()
		defer f.mu.Unlock()
		bad := func() {
			w.WriteHeader(400)
			fmt.Fprint(w, `{"error":"invalid_grant","error_description":"nope"}`)
		}
		switch r.PostForm.Get("grant_type") {
		case "authorization_code":
			ch, ok := f.codes[r.PostForm.Get("code")]
			sum := sha256.Sum256([]byte(r.PostForm.Get("code_verifier")))
			if !ok || base64.RawURLEncoding.EncodeToString(sum[:]) != ch {
				bad()
				return
			}
			delete(f.codes, r.PostForm.Get("code"))
		case "refresh_token":
			rt := r.PostForm.Get("refresh_token")
			used, ok := f.refresh[rt]
			if !ok || used {
				bad()
				return
			}
			f.refresh[rt] = true
		default:
			bad()
			return
		}
		f.issued++
		rt := fmt.Sprintf("yhr_%d", f.issued)
		f.refresh[rt] = false
		json.NewEncoder(w).Encode(map[string]any{
			"access_token": fmt.Sprintf("yhk_access_%d", f.issued), "token_type": "Bearer",
			"expires_in": 3600, "refresh_token": rt, "scope": "account:read vps:read",
		})
	})
	mux.HandleFunc("POST /client/oauth/revoke", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		f.mu.Lock()
		f.revoked = append(f.revoked, r.PostForm.Get("token"))
		f.mu.Unlock()
		fmt.Fprint(w, `{}`)
	})
	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)
	return f, srv
}

func TestLoginFlow(t *testing.T) {
	_, srv := newFakeAS(t)
	c := api.New(srv.URL, nil, "test")
	l := &Login{
		API: c,
		// Le « navigateur » suit la redirection jusqu'au callback loopback.
		OpenBrowser: func(u string) error {
			go func() {
				res, err := http.Get(u)
				if err == nil {
					res.Body.Close()
				}
			}()
			return nil
		},
		Timeout: 10 * time.Second,
	}
	res, err := l.Run(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if res.Token.AccessToken != "yhk_access_1" || res.Token.RefreshToken != "yhr_1" || res.ClientID != "yhc_1" {
		t.Fatalf("res = %+v", res)
	}
}

func TestLoginStateMismatch(t *testing.T) {
	_, srv := newFakeAS(t)
	l := &Login{
		API: api.New(srv.URL, nil, "test"),
		OpenBrowser: func(u string) error {
			q, _ := url.Parse(u)
			cb := q.Query().Get("redirect_uri") + "?code=x&state=forged"
			go http.Get(cb)
			return nil
		},
		Timeout: 5 * time.Second,
	}
	if _, err := l.Run(context.Background()); err == nil {
		t.Fatal("state forgé accepté")
	}
}

func TestProfileTokensRefreshAndRotation(t *testing.T) {
	as, srv := newFakeAS(t)
	cfgPath := filepath.Join(t.TempDir(), "config.yaml")
	f, _ := config.LoadFrom(cfgPath)
	as.refresh["yhr_seed"] = false
	f.Ensure("default").Auth = &config.Auth{Type: config.AuthOAuth, AccessToken: "expired", RefreshToken: "yhr_seed", ClientID: "yhc_x", ExpiresAt: time.Now().Add(-time.Minute)}
	if err := f.Save(); err != nil {
		t.Fatal(err)
	}

	// Deux process (deux lectures du fichier) rafraîchissent en même temps :
	// un seul doit consommer le refresh token, l'autre relire le résultat.
	var wg sync.WaitGroup
	toks := make([]string, 2)
	errs := make([]error, 2)
	for i := range 2 {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			fi, _ := config.LoadFrom(cfgPath)
			ts := &ProfileTokens{File: fi, Profile: "default", OAuth: api.New(srv.URL, nil, "test")}
			toks[i], errs[i] = ts.Token(context.Background())
		}(i)
	}
	wg.Wait()
	for i := range 2 {
		if errs[i] != nil {
			t.Fatalf("process %d : %v", i, errs[i])
		}
	}
	if toks[0] != toks[1] || as.issued != 1 {
		t.Fatalf("toks = %v, émis = %d (attendu 1 seul rafraîchissement)", toks, as.issued)
	}
	g, _ := config.LoadFrom(cfgPath)
	a := g.Profile("default").Auth
	if a.RefreshToken != "yhr_1" || a.AccessToken != "yhk_access_1" || !a.ExpiresAt.After(time.Now()) || len(a.Scopes) != 2 {
		t.Fatalf("auth enregistrée = %+v", a)
	}

	// Sur un 401, Refresh force une rotation.
	ts := &ProfileTokens{File: g, Profile: "default", OAuth: api.New(srv.URL, nil, "test")}
	if tok, _ := ts.Token(context.Background()); tok != "yhk_access_1" {
		t.Fatalf("tok = %s", tok)
	}
	ok, err := ts.Refresh(context.Background())
	if err != nil || !ok {
		t.Fatalf("Refresh = %v %v", ok, err)
	}
	if tok, _ := ts.Token(context.Background()); tok != "yhk_access_2" {
		t.Fatalf("après rotation tok = %s", tok)
	}
}

func TestDeadRefreshClearsSession(t *testing.T) {
	_, srv := newFakeAS(t)
	cfgPath := filepath.Join(t.TempDir(), "config.yaml")
	f, _ := config.LoadFrom(cfgPath)
	f.Ensure("default").Auth = &config.Auth{Type: config.AuthOAuth, AccessToken: "old", RefreshToken: "yhr_revoked", ClientID: "yhc_x"}
	f.Save()
	ts := &ProfileTokens{File: f, Profile: "default", OAuth: api.New(srv.URL, nil, "test")}
	if _, err := ts.Token(context.Background()); !errors.Is(err, ErrInvalidGrant) {
		t.Fatalf("err = %v", err)
	}
	g, _ := config.LoadFrom(cfgPath)
	if g.Profile("default").Auth.RefreshToken != "" {
		t.Fatal("refresh mort conservé")
	}
}

func TestAPIKeyPassthrough(t *testing.T) {
	f, _ := config.LoadFrom(filepath.Join(t.TempDir(), "c.yaml"))
	f.Ensure("default").Auth = &config.Auth{Type: config.AuthAPIKey, APIKey: "yhk_client_v1_k"}
	ts := &ProfileTokens{File: f, Profile: "default"}
	tok, err := ts.Token(context.Background())
	if err != nil || tok != "yhk_client_v1_k" {
		t.Fatalf("tok = %q err = %v", tok, err)
	}
	if ok, _ := ts.Refresh(context.Background()); ok {
		t.Fatal("une clé API ne se rafraîchit pas")
	}
}

func TestPKCE(t *testing.T) {
	v, c := pkce()
	if len(v) != 43 {
		t.Fatalf("verifier de %d caractères", len(v))
	}
	sum := sha256.Sum256([]byte(v))
	if base64.RawURLEncoding.EncodeToString(sum[:]) != c {
		t.Fatal("challenge incorrect")
	}
}
