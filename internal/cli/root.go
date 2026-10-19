// Package cli : arbre de commandes `yorkhost`.
//
// Conventions (style scw / ovhcloud / aws) :
//   - `yorkhost <ressource> <verbe> [ID] [options]` ;
//   - options globales : --profile, --output, --api-url, --yes, --quiet,
//     --debug, --no-color ; variables YORKHOST_* équivalentes ;
//   - priorité : option > variable d'environnement > profil > défaut.
package cli

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/spf13/cobra"

	"github.com/YorkHost-fr/yorkhost-cli/internal/api"
	"github.com/YorkHost-fr/yorkhost-cli/internal/auth"
	"github.com/YorkHost-fr/yorkhost-cli/internal/config"
	"github.com/YorkHost-fr/yorkhost-cli/internal/output"
	"github.com/YorkHost-fr/yorkhost-cli/internal/prompt"
	"github.com/YorkHost-fr/yorkhost-cli/internal/version"
)

// App : état partagé par les commandes d'une exécution.
type App struct {
	// Options globales.
	ProfileFlag string
	OutputFlag  string
	APIURLFlag  string
	Yes         bool
	Quiet       bool
	Debug       bool
	NoColor     bool

	File    *config.File
	Profile string
	Printer *output.Printer
	Prompt  *prompt.Prompter

	// Remplaçables dans les tests.
	Stdin       io.Reader
	Stdout      io.Writer
	Stderr      io.Writer
	OpenBrowser func(string) error
	Interactive *bool

	client *api.Client
}

// NewApp : application avec les entrées/sorties du process.
func NewApp() *App {
	return &App{Stdin: os.Stdin, Stdout: os.Stdout, Stderr: os.Stderr, OpenBrowser: openBrowser}
}

// APIURL : URL de l'API résolue.
func (a *App) APIURL() string {
	if a.APIURLFlag != "" {
		return strings.TrimRight(a.APIURLFlag, "/")
	}
	if v := strings.TrimSpace(os.Getenv("YORKHOST_API_URL")); v != "" {
		return strings.TrimRight(v, "/")
	}
	if p := a.File.Profile(a.Profile); p != nil && p.APIURL != "" {
		return p.APIURL
	}
	return config.DefaultAPIURL
}

// envToken : clé API passée par l'environnement (CI, scripts).
func envToken() string {
	for _, k := range []string{"YORKHOST_API_KEY", "YORKHOST_TOKEN"} {
		if v := strings.TrimSpace(os.Getenv(k)); v != "" {
			return v
		}
	}
	return ""
}

// bare : client sans identifiants (OAuth, routes publiques).
func (a *App) bare() *api.Client {
	c := api.New(a.APIURL(), nil, version.UserAgent())
	if a.Debug {
		c.Debug = a.Stderr
	}
	return c
}

// Client : client authentifié (clé d'environnement, sinon profil).
func (a *App) Client() *api.Client {
	if a.client != nil {
		return a.client
	}
	var ts api.TokenSource
	if t := envToken(); t != "" {
		ts = api.StaticToken(t)
	} else {
		ts = &auth.ProfileTokens{File: a.File, Profile: a.Profile, OAuth: a.bare()}
	}
	c := api.New(a.APIURL(), ts, version.UserAgent())
	if a.Debug {
		c.Debug = a.Stderr
	}
	a.client = c
	return c
}

// init : configuration, profil, sortie (PersistentPreRunE de la racine).
func (a *App) init() error {
	f, err := config.Load()
	if err != nil {
		return err
	}
	a.File = f
	a.Profile = a.ProfileFlag
	if a.Profile == "" {
		a.Profile = strings.TrimSpace(os.Getenv("YORKHOST_PROFILE"))
	}
	if a.Profile == "" {
		a.Profile = f.Active()
	}
	if !config.ValidProfileName(a.Profile) {
		return usageErr(fmt.Errorf("nom de profil invalide %q", a.Profile))
	}
	format := a.OutputFlag
	if format == "" {
		format = strings.TrimSpace(os.Getenv("YORKHOST_OUTPUT"))
	}
	if format == "" {
		if p := f.Profile(a.Profile); p != nil && p.Output != "" {
			format = p.Output
		}
	}
	if format == "" {
		format = config.DefaultOutput
	}
	if !config.ValidOutput(format) {
		return usageErr(fmt.Errorf("format de sortie inconnu %q (table, wide, json, yaml)", format))
	}
	a.Printer = output.New(format, a.NoColor)
	a.Printer.Out, a.Printer.Err, a.Printer.Quiet = a.Stdout, a.Stderr, a.Quiet
	if a.Stdout != os.Stdout {
		a.Printer.Color = false
	}
	a.Prompt = prompt.Default()
	a.Prompt.In, a.Prompt.Out = a.Stdin, a.Stderr
	if a.Interactive != nil {
		a.Prompt.Interactive = *a.Interactive
	}
	if os.Getenv("YORKHOST_YES") == "1" || strings.EqualFold(os.Getenv("YORKHOST_YES"), "true") {
		a.Yes = true
	}
	return nil
}

// ctx : contexte de la commande.
func ctx(cmd *cobra.Command) context.Context {
	if c := cmd.Context(); c != nil {
		return c
	}
	return context.Background()
}

// NewRootCmd : arbre complet.
func NewRootCmd(a *App) *cobra.Command {
	cobra.EnableCommandSorting = false
	root := &cobra.Command{
		Use:   "yorkhost",
		Short: "CLI officiel YorkHost : VPS, serveurs de jeu, facturation, support, sandboxes",
		Long: `yorkhost pilote votre compte YorkHost depuis le terminal : VPS, serveurs de jeu,
services, factures, tickets, anti-DDoS, commandes et sandboxes agents.

Commencez par vous connecter (navigateur, A2F requise) :

  yorkhost login

Puis par exemple :

  yorkhost service list
  yorkhost vps restart 1234 --wait
  yorkhost game get mon-serveur -o json | jq .address
  yorkhost invoice list --status unpaid

Documentation : https://github.com/YorkHost-fr/yorkhost-cli`,
		SilenceUsage:      true,
		SilenceErrors:     true,
		Version:           version.Version,
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error { return a.init() },
	}
	root.SetVersionTemplate(version.String() + "\n")
	root.SetIn(a.Stdin)
	root.SetOut(a.Stdout)
	root.SetErr(a.Stderr)
	root.SetFlagErrorFunc(func(c *cobra.Command, err error) error { return usageErr(err) })

	pf := root.PersistentFlags()
	pf.StringVarP(&a.ProfileFlag, "profile", "p", "", "profil de configuration (env YORKHOST_PROFILE)")
	pf.StringVarP(&a.OutputFlag, "output", "o", "", "format de sortie : table, wide, json, yaml (env YORKHOST_OUTPUT)")
	pf.StringVar(&a.APIURLFlag, "api-url", "", "URL de l'API (env YORKHOST_API_URL)")
	pf.BoolVarP(&a.Yes, "yes", "y", false, "répondre oui aux confirmations (env YORKHOST_YES=1)")
	pf.BoolVarP(&a.Quiet, "quiet", "q", false, "n'afficher que le résultat")
	pf.BoolVar(&a.Debug, "debug", false, "journal des requêtes HTTP sur stderr (sans secrets)")
	pf.BoolVar(&a.NoColor, "no-color", false, "désactiver les couleurs (env NO_COLOR)")
	_ = root.RegisterFlagCompletionFunc("output", fixedCompletion("table", "wide", "json", "yaml"))
	_ = root.RegisterFlagCompletionFunc("profile", func(*cobra.Command, []string, string) ([]string, cobra.ShellCompDirective) {
		f, err := config.Load()
		if err != nil {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		return f.ProfileNames(), cobra.ShellCompDirectiveNoFileComp
	})

	root.AddGroup(
		&cobra.Group{ID: "start", Title: "Démarrer :"},
		&cobra.Group{ID: "infra", Title: "Infrastructure :"},
		&cobra.Group{ID: "account", Title: "Compte et facturation :"},
		&cobra.Group{ID: "tools", Title: "Outils :"},
	)
	add := func(group string, cmds ...*cobra.Command) {
		for _, c := range cmds {
			c.GroupID = group
			root.AddCommand(c)
		}
	}
	add("start", newLoginCmd(a), newLogoutCmd(a), newWhoamiCmd(a), newAuthCmd(a))
	add("infra", newServiceCmd(a), newVPSCmd(a), newGameCmd(a), newSandboxCmd(a), newAntiDDoSCmd(a))
	add("account", newAccountCmd(a), newInvoiceCmd(a), newTicketCmd(a), newCatalogCmd(a), newOrderCmd(a))
	add("tools", newConfigCmd(a), newAPICmd(a), newMCPCmd(a), newVersionCmd(a))
	root.AddCommand(newGenDocsCmd(root))
	root.SetHelpCommandGroupID("tools")
	root.SetCompletionCommandGroupID("tools")
	frenchHelp(root)
	return root
}

// frenchHelp : libellés cobra traduits.
func frenchHelp(root *cobra.Command) {
	root.InitDefaultHelpCmd()
	root.InitDefaultCompletionCmd()
	for _, c := range root.Commands() {
		switch c.Name() {
		case "help":
			c.Short = "Aide sur une commande"
		case "completion":
			c.Short = "Script d'autocomplétion pour bash, zsh, fish ou PowerShell"
			c.Long = `Génère le script d'autocomplétion du shell.

  bash :       source <(yorkhost completion bash)
  zsh :        yorkhost completion zsh > "${fpath[1]}/_yorkhost"
  fish :       yorkhost completion fish > ~/.config/fish/completions/yorkhost.fish
  PowerShell : yorkhost completion powershell | Out-String | Invoke-Expression`
			for _, sh := range c.Commands() {
				sh.Short = "Script d'autocomplétion pour " + sh.Name()
				if f := sh.Flags().Lookup("no-descriptions"); f != nil {
					f.Usage = "sans les descriptions des propositions"
				}
			}
		}
	}
	root.Flags().BoolP("version", "v", false, "version de yorkhost")
	cobra.AddTemplateFunc("frenchUse", func(s string) string { return strings.Replace(s, "[flags]", "[options]", 1) })
	root.SetUsageTemplate(usageTemplate)
	var walk func(*cobra.Command)
	walk = func(c *cobra.Command) {
		if c.Flags().Lookup("help") == nil {
			c.Flags().BoolP("help", "h", false, "aide de "+c.CommandPath())
		}
		for _, sub := range c.Commands() {
			walk(sub)
		}
	}
	walk(root)
}

// translateCobra : messages d'erreur de cobra en français.
func translateCobra(msg string) string {
	r := strings.NewReplacer(
		"unknown command", "commande inconnue",
		"unknown flag", "option inconnue",
		"unknown shorthand flag", "option courte inconnue",
		"\nDid you mean this?\n", "\nVouliez-vous dire :\n",
		"Did you mean this?", "Vouliez-vous dire :",
		"\" for \"", "\" pour \"",
		" in -", " dans -",
		"flag needs an argument", "l'option attend une valeur",
		"invalid argument", "argument invalide",
		"required flag(s)", "option(s) requise(s)",
		" not set", " absente(s)",
	)
	msg = r.Replace(msg)
	for _, re := range cobraArgs {
		msg = re.re.ReplaceAllString(msg, re.repl)
	}
	return msg
}

var cobraArgs = []struct {
	re   *regexp.Regexp
	repl string
}{
	{regexp.MustCompile(`accepts (\d+) arg\(s\), received (\d+)`), "attend $1 argument(s), $2 reçu(s)"},
	{regexp.MustCompile(`accepts at most (\d+) arg\(s\), received (\d+)`), "attend au plus $1 argument(s), $2 reçu(s)"},
	{regexp.MustCompile(`accepts between (\d+) and (\d+) arg\(s\), received (\d+)`), "attend entre $1 et $2 arguments, $3 reçu(s)"},
	{regexp.MustCompile(`requires at least (\d+) arg\(s\), only received (\d+)`), "attend au moins $1 argument(s), $2 reçu(s)"},
}

const usageTemplate = `Utilisation :{{if .Runnable}}
  {{frenchUse .UseLine}}{{end}}{{if .HasAvailableSubCommands}}
  {{.CommandPath}} [commande]{{end}}{{if gt (len .Aliases) 0}}

Alias :
  {{.NameAndAliases}}{{end}}{{if .HasExample}}

Exemples :
{{.Example}}{{end}}{{if .HasAvailableSubCommands}}{{$cmds := .Commands}}{{if eq (len .Groups) 0}}

Commandes :{{range $cmds}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{else}}{{range $group := .Groups}}

{{.Title}}{{range $cmds}}{{if (and (eq .GroupID $group.ID) (or .IsAvailableCommand (eq .Name "help")))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{if not .AllChildCommandsHaveGroup}}

Autres commandes :{{range $cmds}}{{if (and (eq .GroupID "") (or .IsAvailableCommand (eq .Name "help")))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

Options :
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}

Options globales :
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableSubCommands}}

Lancez "{{.CommandPath}} [commande] --help" pour l'aide d'une commande.{{end}}
`

// Execute : lance le CLI et renvoie le code de sortie.
func Execute(ctx context.Context, a *App, args []string) int {
	root := NewRootCmd(a)
	root.SetArgs(args)
	err := root.ExecuteContext(ctx)
	if err == nil {
		return 0
	}
	code := ExitCode(err)
	if errors.Is(err, context.Canceled) {
		return code
	}
	msg := err.Error()
	if code == ExitUsage {
		msg = translateCobra(msg)
	}
	var ue *usageError
	if (errors.As(err, &ue) || code == ExitUsage) && !strings.Contains(msg, "--help") {
		msg += "\nVoir « yorkhost --help » ou « yorkhost <commande> --help »."
	}
	prefix := "Erreur : "
	if a.Printer != nil && a.Printer.Color {
		prefix = output.Red("Erreur") + " : "
	}
	fmt.Fprintln(a.Stderr, prefix+msg)
	return code
}

// fixedCompletion : complétion par liste fixe.
func fixedCompletion(values ...string) func(*cobra.Command, []string, string) ([]string, cobra.ShellCompDirective) {
	return func(*cobra.Command, []string, string) ([]string, cobra.ShellCompDirective) {
		return values, cobra.ShellCompDirectiveNoFileComp
	}
}
