package cli

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/YorkHost-fr/yorkhost-cli/internal/api"
	"github.com/YorkHost-fr/yorkhost-cli/internal/auth"
	"github.com/YorkHost-fr/yorkhost-cli/internal/config"
	"github.com/YorkHost-fr/yorkhost-cli/internal/output"
)

// Format des clés client (hostbill-apiclient, README « Format des clés »).
var apiKeyRe = regexp.MustCompile(`^yhk_client_v1_[0-9A-Za-z]{16}_[0-9A-Za-z]{40}$`)

func newLoginCmd(a *App) *cobra.Command {
	var useKey, withToken, noBrowser bool
	var scopes []string
	cmd := &cobra.Command{
		Use:   "login",
		Short: "Se connecter (navigateur + A2F, ou clé API)",
		Long: `Connecte le profil courant à votre compte YorkHost.

Par défaut, le navigateur s'ouvre sur votre espace client : connectez-vous
(double authentification requise), choisissez les permissions et les services
accessibles au CLI, puis validez. Le CLI reçoit un jeton d'accès (1 h) et un
jeton de rafraîchissement (30 jours glissants), renouvelés automatiquement.
L'accès se retire dans l'espace client → Clés API → Applications connectées.

Avec --api-key, le CLI utilise une clé créée dans l'espace client
(Clés API → Nouvelle clé) : pratique pour un serveur ou une CI. La clé est
demandée en saisie masquée, ou lue sur stdin avec --with-token.`,
		Example: `  yorkhost login
  yorkhost login --no-browser            # machine distante : ouvrez le lien vous-même
  yorkhost login --api-key
  echo "$YORKHOST_KEY" | yorkhost login --with-token --profile ci
  yorkhost login --scopes services:read,vps:read,vps:power`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			c := ctx(cmd)
			if envToken() != "" {
				a.Printer.Warn("YORKHOST_API_KEY est définie : elle reste prioritaire sur le profil tant qu'elle existe.")
			}
			if useKey || withToken {
				return a.loginWithKey(c, withToken)
			}
			return a.loginOAuth(c, noBrowser, scopes)
		},
	}
	f := cmd.Flags()
	f.BoolVar(&useKey, "api-key", false, "utiliser une clé API (saisie masquée)")
	f.BoolVar(&withToken, "with-token", false, "lire la clé API sur stdin")
	f.BoolVar(&noBrowser, "no-browser", false, "ne pas ouvrir le navigateur, afficher le lien")
	f.StringSliceVar(&scopes, "scopes", nil, "permissions demandées (défaut : toutes, à affiner sur la page de consentement)")
	_ = cmd.RegisterFlagCompletionFunc("scopes", fixedCompletion(auth.Scopes...))
	return cmd
}

func (a *App) loginWithKey(c context.Context, fromStdin bool) error {
	var key string
	var err error
	if fromStdin {
		key, err = a.readText("-", "clé API")
	} else {
		a.Printer.Info("Créez une clé dans l'espace client : https://clients.yorkhost.fr/index.php?cmd=apiclient")
		key, err = a.Prompt.Secret("Clé API (yhk_client_v1_…)")
	}
	if err != nil {
		return err
	}
	key = strings.TrimSpace(key)
	if !apiKeyRe.MatchString(key) {
		return usageErr(errors.New("format de clé invalide : attendu yhk_client_v1_<16 car.>_<40 car.>"))
	}
	client := api.New(a.APIURL(), api.StaticToken(key), a.bare().UserAgent)
	client.Debug = a.bare().Debug
	account, err := checkAccount(c, client)
	if err != nil {
		return err
	}
	p := a.File.Ensure(a.Profile)
	a.revokeOAuth(c, p.Auth)
	p.Auth = &config.Auth{Type: config.AuthAPIKey, APIKey: key}
	if a.File.ActiveProfile == "" {
		a.File.ActiveProfile = a.Profile
	}
	if err := a.File.Save(); err != nil {
		return err
	}
	a.loggedIn(account)
	return nil
}

func (a *App) loginOAuth(c context.Context, noBrowser bool, scopes []string) error {
	l := &auth.Login{
		API:    a.bare(),
		Scopes: scopes,
		Notify: func(format string, args ...any) { a.Printer.Info(format, args...) },
	}
	if !noBrowser {
		l.OpenBrowser = a.OpenBrowser
	}
	res, err := l.Run(c)
	if err != nil {
		return err
	}
	p := a.File.Ensure(a.Profile)
	a.revokeOAuth(c, p.Auth)
	p.Auth = &config.Auth{ClientID: res.ClientID, RedirectURI: res.RedirectURI}
	auth.Apply(p.Auth, &res.Token, time.Now())
	if a.File.ActiveProfile == "" {
		a.File.ActiveProfile = a.Profile
	}
	if err := a.File.Save(); err != nil {
		return err
	}
	client := api.New(a.APIURL(), api.StaticToken(res.Token.AccessToken), a.bare().UserAgent)
	account, err := checkAccount(c, client)
	if err != nil {
		return err
	}
	a.loggedIn(account)
	return nil
}

// checkAccount : valide les identifiants. Une clé sans account:read reste
// valable (403 missing_scope) : on ne lit simplement pas le profil.
func checkAccount(c context.Context, client *api.Client) (map[string]any, error) {
	var account map[string]any
	err := client.Get(c, "/client/v1/account", nil, &account)
	var ae *api.Error
	if errors.As(err, &ae) && ae.Code == "missing_scope" {
		return nil, nil
	}
	return account, err
}

func (a *App) loggedIn(account map[string]any) {
	who := ""
	if account != nil {
		who = strings.TrimSpace(str(account, "firstName") + " " + str(account, "lastName"))
		if e := str(account, "email"); e != "" {
			who += " <" + e + ">"
		}
	}
	if who != "" {
		a.Printer.Success("Connecté en tant que %s (profil %q).", strings.TrimSpace(who), a.Profile)
	} else {
		a.Printer.Success("Connecté (profil %q).", a.Profile)
	}
	if a.Printer.Machine() {
		_ = a.Printer.Raw(map[string]any{"profile": a.Profile, "account": account})
	}
}

// revokeOAuth : révoque une ancienne session OAuth (meilleur effort).
func (a *App) revokeOAuth(c context.Context, old *config.Auth) {
	if old == nil || old.Type != config.AuthOAuth || old.ClientID == "" {
		return
	}
	tok := old.RefreshToken
	if tok == "" {
		tok = old.AccessToken
	}
	if tok == "" {
		return
	}
	rc, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()
	if err := auth.Revoke(rc, a.bare(), old.ClientID, tok); err != nil && a.Debug {
		a.Printer.Warn("révocation de l'ancienne session : %v", err)
	}
}

func newLogoutCmd(a *App) *cobra.Command {
	var all bool
	cmd := &cobra.Command{
		Use:   "logout",
		Short: "Se déconnecter (révoque la session OAuth, oublie les identifiants)",
		Long: `Supprime les identifiants du profil. Une session OAuth est aussi révoquée
côté YorkHost. Une clé API n'est pas révoquée : faites-le dans l'espace client
(Clés API) si elle ne doit plus servir.`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			names := []string{a.Profile}
			if all {
				names = a.File.ProfileNames()
			}
			done := 0
			for _, n := range names {
				p := a.File.Profile(n)
				if p == nil || p.Auth.Empty() {
					continue
				}
				a.revokeOAuth(ctx(cmd), p.Auth)
				if p.Auth.Type == config.AuthAPIKey {
					a.Printer.Info("Profil %q : clé API oubliée localement (révocable dans l'espace client → Clés API).", n)
				}
				p.Auth = nil
				done++
			}
			if done == 0 {
				a.Printer.Info("Aucune session à fermer.")
				return nil
			}
			if err := a.File.Save(); err != nil {
				return err
			}
			a.Printer.Success("Déconnecté (%d profil(s)).", done)
			return nil
		},
	}
	cmd.Flags().BoolVar(&all, "all", false, "tous les profils")
	return cmd
}

func newWhoamiCmd(a *App) *cobra.Command {
	return &cobra.Command{
		Use:   "whoami",
		Short: "Compte et identifiants utilisés",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			var account map[string]any
			if err := a.Client().Get(ctx(cmd), "/client/v1/account", nil, &account); err != nil {
				return err
			}
			info := a.authInfo()
			for k, v := range info {
				account["_cli_"+k] = v
			}
			if a.Printer.Machine() {
				return a.Printer.Raw(map[string]any{"account": stripCLI(account), "cli": info})
			}
			return a.Printer.Object(account, account, append(accountFields(),
				output.Field{Header: "Profil", Path: "_cli_profile"},
				output.Field{Header: "API", Path: "_cli_api_url"},
				output.Field{Header: "Authentification", Path: "_cli_auth"},
				output.Field{Header: "Session expire", Path: "_cli_expires_at", Format: output.Date},
				output.Field{Header: "Permissions", Path: "_cli_scopes", Wide: true},
			))
		},
	}
}

func stripCLI(m map[string]any) map[string]any {
	out := map[string]any{}
	for k, v := range m {
		if !strings.HasPrefix(k, "_cli_") {
			out[k] = v
		}
	}
	return out
}

// authInfo : description des identifiants effectifs (sans secret).
func (a *App) authInfo() map[string]any {
	info := map[string]any{"profile": a.Profile, "api_url": a.APIURL()}
	if t := envToken(); t != "" {
		info["auth"] = "clé API (variable d'environnement) " + maskToken(t)
		return info
	}
	p := a.File.Profile(a.Profile)
	if p == nil || p.Auth.Empty() {
		info["auth"] = "aucune"
		return info
	}
	switch p.Auth.Type {
	case config.AuthOAuth:
		info["auth"] = "session navigateur (OAuth)"
		if !p.Auth.ExpiresAt.IsZero() {
			info["expires_at"] = p.Auth.ExpiresAt.Format(time.RFC3339)
		}
		if len(p.Auth.Scopes) > 0 {
			scopes := make([]any, len(p.Auth.Scopes))
			for i, s := range p.Auth.Scopes {
				scopes[i] = s
			}
			info["scopes"] = scopes
		}
	default:
		info["auth"] = "clé API " + maskToken(p.Auth.APIKey)
	}
	return info
}

// maskToken : identifiant public de la clé, secret masqué.
func maskToken(t string) string {
	if m := regexp.MustCompile(`^(yhk_client_v1_[0-9A-Za-z]{16})_`).FindStringSubmatch(t); m != nil {
		return m[1] + "_…"
	}
	if len(t) > 8 {
		return t[:8] + "…"
	}
	return "…"
}

func newAuthCmd(a *App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "auth",
		Short: "Sessions et identifiants (status, token, login, logout)",
	}
	status := &cobra.Command{
		Use:   "status",
		Short: "Identifiants de chaque profil (sans les secrets)",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			var rows []map[string]any
			for _, n := range a.File.ProfileNames() {
				p := a.File.Profile(n)
				row := map[string]any{"profile": n, "active": n == a.File.Active(), "api_url": p.APIURL}
				if p.APIURL == "" {
					row["api_url"] = config.DefaultAPIURL
				}
				switch {
				case p.Auth.Empty():
					row["auth"] = "aucune"
				case p.Auth.Type == config.AuthOAuth:
					row["auth"] = "oauth"
					row["detail"] = "rafraîchie automatiquement"
					if p.Auth.RefreshToken == "" {
						row["detail"] = "expirée : yorkhost login"
					}
				default:
					row["auth"] = "api_key"
					row["detail"] = maskToken(p.Auth.APIKey)
				}
				rows = append(rows, row)
			}
			if envToken() != "" {
				a.Printer.Info("YORKHOST_API_KEY est définie et prioritaire sur les profils.")
			}
			return a.Printer.List(rows, rows, []output.Column{
				{Header: "Profil", Path: "profile"},
				{Header: "Actif", Path: "active", Format: func(v any) string {
					if v == true {
						return "*"
					}
					return ""
				}},
				{Header: "Auth", Path: "auth"},
				{Header: "Détail", Path: "detail"},
				{Header: "API", Path: "api_url"},
			})
		},
	}
	token := &cobra.Command{
		Use:     "token",
		Short:   "Afficher le jeton d'accès courant (pour curl, scripts)",
		Example: `  curl -H "Authorization: Bearer $(yorkhost auth token)" https://api.yorkhost.fr/client/v1/services`,
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			t, err := a.Client().Tokens.Token(ctx(cmd))
			if err != nil {
				return err
			}
			fmt.Fprintln(a.Stdout, t)
			return nil
		},
	}
	cmd.AddCommand(status, token, newLoginCmd(a), newLogoutCmd(a))
	return cmd
}
