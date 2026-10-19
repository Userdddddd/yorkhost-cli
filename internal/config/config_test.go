package config

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"
)

func TestSaveLoadRoundTrip(t *testing.T) {
	p := filepath.Join(t.TempDir(), "sub", "config.yaml")
	f, err := LoadFrom(p)
	if err != nil {
		t.Fatal(err)
	}
	f.ActiveProfile = "pro"
	pr := f.Ensure("pro")
	if err := pr.Set("output", "json"); err != nil {
		t.Fatal(err)
	}
	pr.Auth = &Auth{Type: AuthOAuth, AccessToken: "a", RefreshToken: "r", ExpiresAt: time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC), ClientID: "yhc_x"}
	if err := f.Save(); err != nil {
		t.Fatal(err)
	}
	if runtime.GOOS != "windows" {
		st, err := os.Stat(p)
		if err != nil {
			t.Fatal(err)
		}
		if st.Mode().Perm() != 0o600 {
			t.Fatalf("droits = %v, attendu 0600", st.Mode().Perm())
		}
	}
	g, err := LoadFrom(p)
	if err != nil {
		t.Fatal(err)
	}
	if g.Active() != "pro" || g.Profile("pro").Output != "json" || g.Profile("pro").Auth.RefreshToken != "r" || !g.Profile("pro").Auth.ExpiresAt.Equal(pr.Auth.ExpiresAt) {
		t.Fatalf("relu : %+v", g.Profile("pro"))
	}
}

func TestSetValidation(t *testing.T) {
	p := &Profile{}
	if err := p.Set("output", "xml"); err == nil {
		t.Fatal("format invalide accepté")
	}
	if err := p.Set("api_url", "ftp://x"); err == nil {
		t.Fatal("URL invalide acceptée")
	}
	if err := p.Set("api_url", "https://api.example.test/"); err != nil || p.APIURL != "https://api.example.test" {
		t.Fatalf("api_url = %q err = %v", p.APIURL, err)
	}
	if err := p.Set("nope", "x"); err == nil {
		t.Fatal("réglage inconnu accepté")
	}
}

func TestProfileNames(t *testing.T) {
	for name, ok := range map[string]bool{"default": true, "pro-2": true, "a.b_c": true, "": false, "-x": false, "a b": false} {
		if ValidProfileName(name) != ok {
			t.Errorf("ValidProfileName(%q) != %v", name, ok)
		}
	}
}

func TestLock(t *testing.T) {
	f, _ := LoadFrom(filepath.Join(t.TempDir(), "config.yaml"))
	unlock, err := f.Lock(time.Second)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.Lock(200 * time.Millisecond); err == nil {
		t.Fatal("double verrou accordé")
	}
	unlock()
	unlock2, err := f.Lock(time.Second)
	if err != nil {
		t.Fatal(err)
	}
	unlock2()
}
