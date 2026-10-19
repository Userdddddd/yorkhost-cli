// Package auth : connexion OAuth 2.1 (navigateur + PKCE + redirection
// loopback, RFC 8252) et source de jetons avec rafraîchissement automatique.
//
// Serveur d'autorisation : yorkhost-api (/client/oauth/*). La connexion et le
// consentement (permissions, services) se font dans l'espace client HostBill,
// A2F obligatoire. Le jeton d'accès (1 h) a le même format qu'une clé API
// (yhk_client_v1_…) ; le refresh token (30 j glissants) est à usage unique.
package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/YorkHost-fr/yorkhost-cli/internal/api"
)

// ClientName : nom affiché sur la page de consentement et dans
// « Applications connectées ».
const ClientName = "YorkHost CLI"

// Scopes demandés par défaut : tout le catalogue (le client choisit
// ensuite sur la page de consentement).
var Scopes = []string{
	"account:read", "services:read", "billing:read",
	"tickets:read", "tickets:write",
	"vps:read", "vps:power", "vps:manage", "vps:reinstall",
	"game:read", "game:power", "game:console", "game:files:read", "game:manage", "game:reinstall",
	"dns:read", "dns:write",
	"antiddos:read", "antiddos:write",
	"orders:write",
	"sandbox:read", "sandbox:write",
}

// Token : réponse de /client/oauth/token.
type Token struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
}

// ExpiresAt : échéance absolue (avec une marge de sécurité).
func (t *Token) ExpiresAt(now time.Time) time.Time {
	d := time.Duration(t.ExpiresIn) * time.Second
	if d <= 0 {
		d = time.Hour
	}
	return now.Add(d - 60*time.Second)
}

// Login : paramètres de la connexion par navigateur.
type Login struct {
	API    *api.Client
	Scopes []string
	// OpenBrowser : ouvre l'URL (nil = ne rien ouvrir, l'URL est affichée).
	OpenBrowser func(string) error
	// Notify : messages à l'utilisateur (URL à ouvrir, attente).
	Notify func(format string, args ...any)
	// Timeout : durée maximale d'attente du retour navigateur.
	Timeout time.Duration
}

// Result : jeton obtenu et application OAuth enregistrée.
type Result struct {
	Token       Token
	ClientID    string
	RedirectURI string
}

// Run : enregistre une application publique pour cette redirection
// loopback, ouvre la page d'autorisation et attend le code.
func (l *Login) Run(ctx context.Context) (*Result, error) {
	if l.Timeout == 0 {
		l.Timeout = 10 * time.Minute
	}
	notify := l.Notify
	if notify == nil {
		notify = func(string, ...any) {}
	}
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return nil, fmt.Errorf("impossible d'ouvrir un port local pour la redirection : %w", err)
	}
	defer ln.Close()
	redirectURI := fmt.Sprintf("http://127.0.0.1:%d/callback", ln.Addr().(*net.TCPAddr).Port)

	clientID, err := Register(ctx, l.API, redirectURI)
	if err != nil {
		return nil, err
	}
	verifier, challenge := pkce()
	state := randomString(24)
	scopes := l.Scopes
	if len(scopes) == 0 {
		scopes = Scopes
	}
	q := url.Values{
		"response_type":         {"code"},
		"client_id":             {clientID},
		"redirect_uri":          {redirectURI},
		"scope":                 {strings.Join(scopes, " ")},
		"state":                 {state},
		"code_challenge":        {challenge},
		"code_challenge_method": {"S256"},
	}
	authURL := l.API.BaseURL + "/client/oauth/authorize?" + q.Encode()

	type cb struct {
		code string
		err  error
	}
	done := make(chan cb, 1)
	mux := http.NewServeMux()
	mux.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		v := r.URL.Query()
		var res cb
		switch {
		case v.Get("state") != state:
			res.err = errors.New("réponse d'autorisation invalide (state différent) : relancez yorkhost login")
		case v.Get("error") != "":
			res.err = fmt.Errorf("autorisation refusée : %s", strings.TrimSpace(v.Get("error")+" "+v.Get("error_description")))
		case v.Get("code") == "":
			res.err = errors.New("réponse d'autorisation sans code")
		default:
			res.code = v.Get("code")
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Cache-Control", "no-store")
		if res.err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, page("Connexion refusée", html.EscapeString(res.err.Error())))
		} else {
			fmt.Fprint(w, page("Connexion réussie", "Vous pouvez fermer cet onglet et revenir au terminal."))
		}
		select {
		case done <- res:
		default:
		}
	})
	srv := &http.Server{Handler: mux, ReadHeaderTimeout: 10 * time.Second}
	go func() { _ = srv.Serve(ln) }()
	defer func() { _ = srv.Close() }()

	opened := false
	if l.OpenBrowser != nil {
		opened = l.OpenBrowser(authURL) == nil
	}
	if opened {
		notify("Ouverture du navigateur pour vous connecter à votre espace client YorkHost…")
		notify("Si rien ne s'ouvre, copiez ce lien :\n\n  %s\n", authURL)
	} else {
		notify("Ouvrez ce lien dans un navigateur (sur cette machine) pour vous connecter :\n\n  %s\n", authURL)
	}
	notify("En attente de l'autorisation (A2F requise)…")

	ctx, cancel := context.WithTimeout(ctx, l.Timeout)
	defer cancel()
	var res cb
	select {
	case res = <-done:
	case <-ctx.Done():
		if errors.Is(ctx.Err(), context.Canceled) {
			return nil, ctx.Err()
		}
		return nil, errors.New("délai dépassé en attendant l'autorisation dans le navigateur")
	}
	if res.err != nil {
		return nil, res.err
	}
	tok, err := Exchange(ctx, l.API, clientID, url.Values{
		"grant_type":    {"authorization_code"},
		"code":          {res.code},
		"redirect_uri":  {redirectURI},
		"code_verifier": {verifier},
	})
	if err != nil {
		return nil, err
	}
	return &Result{Token: *tok, ClientID: clientID, RedirectURI: redirectURI}, nil
}

// Register : enregistrement dynamique (RFC 7591) d'un client public.
func Register(ctx context.Context, c *api.Client, redirectURI string) (string, error) {
	resp, err := c.Do(ctx, api.Request{
		Method: http.MethodPost, Path: "/client/oauth/register", NoAuth: true,
		Body: map[string]any{
			"client_name":                ClientName,
			"client_uri":                 "https://github.com/YorkHost-fr/yorkhost-cli",
			"redirect_uris":              []string{redirectURI},
			"token_endpoint_auth_method": "none",
			"grant_types":                []string{"authorization_code", "refresh_token"},
			"response_types":             []string{"code"},
		},
	})
	if err != nil {
		return "", fmt.Errorf("enregistrement OAuth : %w", err)
	}
	var out struct {
		ClientID string `json:"client_id"`
	}
	if err := resp.Decode(&out); err != nil {
		return "", err
	}
	if out.ClientID == "" {
		return "", errors.New("enregistrement OAuth : client_id absent de la réponse")
	}
	return out.ClientID, nil
}

// Exchange : POST /client/oauth/token (formulaire).
func Exchange(ctx context.Context, c *api.Client, clientID string, form url.Values) (*Token, error) {
	form.Set("client_id", clientID)
	body, status, err := postForm(ctx, c, "/client/oauth/token", form)
	if err != nil {
		return nil, err
	}
	if status != http.StatusOK {
		return nil, oauthErr(status, body)
	}
	var t Token
	if err := json.Unmarshal(body, &t); err != nil || t.AccessToken == "" {
		return nil, fmt.Errorf("réponse de jeton illisible (HTTP %d)", status)
	}
	return &t, nil
}

// Revoke : POST /client/oauth/revoke (RFC 7009) — révoque tout le grant.
func Revoke(ctx context.Context, c *api.Client, clientID, token string) error {
	body, status, err := postForm(ctx, c, "/client/oauth/revoke", url.Values{"client_id": {clientID}, "token": {token}})
	if err != nil {
		return err
	}
	if status != http.StatusOK {
		return oauthErr(status, body)
	}
	return nil
}

// ErrInvalidGrant : refresh token refusé (révoqué, expiré, déjà utilisé) —
// il faut se reconnecter.
var ErrInvalidGrant = errors.New("session expirée ou révoquée : relancez `yorkhost login`")

func oauthErr(status int, body []byte) error {
	var e struct {
		Error       string `json:"error"`
		Description string `json:"error_description"`
	}
	_ = json.Unmarshal(body, &e)
	if e.Error == "invalid_grant" {
		if e.Description != "" {
			return fmt.Errorf("%w (%s)", ErrInvalidGrant, e.Description)
		}
		return ErrInvalidGrant
	}
	if e.Error != "" {
		return fmt.Errorf("OAuth : %s (%s, HTTP %d)", e.Description, e.Error, status)
	}
	return fmt.Errorf("OAuth : HTTP %d", status)
}

func postForm(ctx context.Context, c *api.Client, path string, form url.Values) ([]byte, int, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.BaseURL+path, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, 0, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")
	if c.UserAgent != "" {
		req.Header.Set("User-Agent", c.UserAgent)
	}
	res, err := c.HTTP.Do(req)
	if err != nil {
		return nil, 0, &api.NetworkError{URL: c.BaseURL, Err: err}
	}
	defer res.Body.Close()
	b, err := io.ReadAll(io.LimitReader(res.Body, 1<<20))
	if err != nil {
		return nil, 0, &api.NetworkError{URL: c.BaseURL, Err: err}
	}
	return b, res.StatusCode, nil
}

// pkce : code_verifier (43 car. base64url) et code_challenge S256.
func pkce() (verifier, challenge string) {
	verifier = randomString(32)
	sum := sha256.Sum256([]byte(verifier))
	return verifier, base64.RawURLEncoding.EncodeToString(sum[:])
}

func randomString(n int) string {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	return base64.RawURLEncoding.EncodeToString(b)
}

func page(title, msg string) string {
	return `<!doctype html><html lang="fr"><head><meta charset="utf-8"><meta name="viewport" content="width=device-width,initial-scale=1"><title>YorkHost CLI</title>
<style>body{font-family:system-ui,sans-serif;background:#0b1020;color:#e7eaf3;display:grid;place-items:center;min-height:100vh;margin:0}
main{max-width:28rem;padding:2rem;text-align:center}h1{font-size:1.4rem}p{color:#aab1c5}</style></head>
<body><main><h1>` + title + `</h1><p>` + msg + `</p></main></body></html>`
}
