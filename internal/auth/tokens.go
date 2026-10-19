package auth

import (
	"context"
	"errors"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/YorkHost-fr/yorkhost-cli/internal/api"
	"github.com/YorkHost-fr/yorkhost-cli/internal/config"
)

// ProfileTokens : jetons d'un profil du fichier de configuration. Une clé API
// est renvoyée telle quelle ; un jeton OAuth est rafraîchi avant son échéance
// (ou sur un 401), sous verrou, et le nouveau couple est enregistré aussitôt
// (le refresh token précédent n'est plus valable).
type ProfileTokens struct {
	File    *config.File
	Profile string
	// OAuth : client utilisé pour /client/oauth/token (même URL d'API).
	OAuth *api.Client
	Now   func() time.Time

	mu      sync.Mutex
	current string
}

func (p *ProfileTokens) now() time.Time {
	if p.Now != nil {
		return p.Now()
	}
	return time.Now()
}

func (p *ProfileTokens) auth() *config.Auth {
	if pr := p.File.Profile(p.Profile); pr != nil {
		return pr.Auth
	}
	return nil
}

// Token : jeton Bearer à utiliser.
func (p *ProfileTokens) Token(ctx context.Context) (string, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	a := p.auth()
	if a.Empty() {
		return "", api.ErrNotLoggedIn
	}
	if a.Type != config.AuthOAuth {
		p.current = a.APIKey
		return a.APIKey, nil
	}
	if a.AccessToken != "" && p.now().Before(a.ExpiresAt) {
		p.current = a.AccessToken
		return a.AccessToken, nil
	}
	if err := p.refreshLocked(ctx, ""); err != nil {
		return "", err
	}
	p.current = p.auth().AccessToken
	return p.current, nil
}

// Refresh : appelé sur un 401. Faux pour une clé API (rien à rafraîchir).
func (p *ProfileTokens) Refresh(ctx context.Context) (bool, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	a := p.auth()
	if a.Empty() || a.Type != config.AuthOAuth || a.RefreshToken == "" {
		return false, nil
	}
	if err := p.refreshLocked(ctx, p.current); err != nil {
		return false, err
	}
	p.current = p.auth().AccessToken
	return true, nil
}

// refreshLocked : rafraîchit sauf si un autre process l'a déjà fait (jeton
// relu différent de rejected et encore valide).
func (p *ProfileTokens) refreshLocked(ctx context.Context, rejected string) error {
	unlock, err := p.File.Lock(15 * time.Second)
	if err != nil {
		return err
	}
	defer unlock()
	if fresh, err := config.LoadFrom(p.File.FilePath()); err == nil {
		*p.File = *fresh
	}
	pr := p.File.Profile(p.Profile)
	if pr == nil || pr.Auth.Empty() {
		return api.ErrNotLoggedIn
	}
	a := pr.Auth
	if a.AccessToken != "" && a.AccessToken != rejected && p.now().Before(a.ExpiresAt) {
		return nil
	}
	if a.RefreshToken == "" || a.ClientID == "" {
		return ErrInvalidGrant
	}
	tok, err := Exchange(ctx, p.OAuth, a.ClientID, url.Values{
		"grant_type":    {"refresh_token"},
		"refresh_token": {a.RefreshToken},
	})
	if err != nil {
		if errors.Is(err, ErrInvalidGrant) {
			// Session morte : on oublie les jetons pour ne pas réessayer en boucle.
			a.AccessToken, a.RefreshToken = "", ""
			_ = p.File.Save()
		}
		return err
	}
	Apply(a, tok, p.now())
	return p.File.Save()
}

// Apply : enregistre un jeton OAuth dans les identifiants d'un profil.
func Apply(a *config.Auth, tok *Token, now time.Time) {
	a.Type = config.AuthOAuth
	a.APIKey = ""
	a.AccessToken = tok.AccessToken
	if tok.RefreshToken != "" {
		a.RefreshToken = tok.RefreshToken
	}
	a.ExpiresAt = tok.ExpiresAt(now).UTC().Truncate(time.Second)
	if tok.Scope != "" {
		a.Scopes = strings.Fields(tok.Scope)
	}
}
