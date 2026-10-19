package cli

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/YorkHost-fr/yorkhost-cli/internal/output"
)

// machine : ce qui distingue VPS et serveurs de jeu (mêmes verbes).
type machine struct {
	kind     string // type de service : vps, game
	resource string // segment d'API : vps, game-servers
	label    string // « VPS », « serveur de jeu »
	plural   string
	fields   []output.Field
}

func (m machine) getCmd(a *App) *cobra.Command {
	return &cobra.Command{
		Use:               "get [id|nom]",
		Short:             "État et ressources en direct",
		Args:              noArgsOrOne,
		ValidArgsFunction: a.serviceCompletion(m.kind),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := a.resolveService(ctx(cmd), args, m.kind)
			if err != nil {
				return err
			}
			var obj map[string]any
			if err := a.Client().Get(ctx(cmd), path(m.resource, id), nil, &obj); err != nil {
				return err
			}
			return a.Printer.Object(obj, obj, m.fields)
		},
	}
}

func (m machine) listCmd(a *App) *cobra.Command {
	var status string
	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "Lister vos " + m.plural,
		Args:    cobra.NoArgs,
		RunE:    func(cmd *cobra.Command, _ []string) error { return a.listServices(cmd, m.kind, status, 0) },
	}
	cmd.Flags().StringVar(&status, "status", "", "filtrer par statut : "+strings.Join(serviceStatuses, ", "))
	_ = cmd.RegisterFlagCompletionFunc("status", fixedCompletion(serviceStatuses...))
	return cmd
}

// powerCmd : start / stop / restart / kill.
func (m machine) powerCmd(a *App, action, short string, destructive bool) *cobra.Command {
	var wait bool
	var timeout time.Duration
	cmd := &cobra.Command{
		Use:               action + " [id|nom]",
		Short:             short,
		Args:              noArgsOrOne,
		ValidArgsFunction: a.serviceCompletion(m.kind),
		RunE: func(cmd *cobra.Command, args []string) error {
			c := ctx(cmd)
			id, err := a.resolveService(c, args, m.kind)
			if err != nil {
				return err
			}
			if destructive {
				if err := a.Prompt.Confirm(fmt.Sprintf("%s le %s %s ? Les connexions en cours seront coupées.", verb(action), m.label, id), a.Yes); err != nil {
					return err
				}
			}
			if _, err := a.Client().Post(c, path(m.resource, id, "power"), map[string]string{"action": action}, nil); err != nil {
				return err
			}
			if !wait {
				a.Printer.Success("%s %s : %s demandé.", m.label, id, action)
				if a.Printer.Machine() {
					return a.Printer.Raw(map[string]any{"id": id, "action": action, "accepted": true})
				}
				return nil
			}
			want := "running"
			if action == "stop" || action == "kill" {
				want = "offline"
			}
			if action == "restart" {
				// Laisse le temps à l'arrêt d'être visible avant d'attendre le retour.
				time.Sleep(5 * time.Second)
			}
			a.Printer.Info("Attente de l'état %q…", want)
			var last map[string]any
			err = a.waitFor(c, m.label+" "+id, timeout, func() (bool, string, error) {
				last = nil
				if err := a.Client().Get(c, path(m.resource, id), nil, &last); err != nil {
					return false, "", err
				}
				st := str(last, "state")
				return st == want, st, nil
			})
			if err != nil {
				return err
			}
			a.Printer.Success("%s %s : %s.", m.label, id, want)
			if a.Printer.Machine() {
				return a.Printer.Raw(last)
			}
			return nil
		},
	}
	cmd.Flags().BoolVarP(&wait, "wait", "w", false, "attendre que l'action soit effective")
	cmd.Flags().DurationVar(&timeout, "timeout", 5*time.Minute, "durée maximale d'attente avec --wait")
	return cmd
}

func verb(action string) string {
	switch action {
	case "stop":
		return "Arrêter"
	case "restart":
		return "Redémarrer"
	case "kill":
		return "Tuer (arrêt forcé, sans sauvegarde)"
	}
	return "Démarrer"
}

func (m machine) metricsCmd(a *App) *cobra.Command {
	return &cobra.Command{
		Use:               "metrics [id|nom]",
		Short:             "CPU et RAM de la dernière heure",
		Long:              "CPU et RAM de la dernière heure. En -o wide, chaque point est listé ; en -o json, la série brute.",
		Args:              noArgsOrOne,
		ValidArgsFunction: a.serviceCompletion(m.kind),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := a.resolveService(ctx(cmd), args, m.kind)
			if err != nil {
				return err
			}
			var obj map[string]any
			if err := a.Client().Get(ctx(cmd), path(m.resource, id, "metrics"), nil, &obj); err != nil {
				return err
			}
			points := list(obj, "points")
			if a.Printer.Machine() {
				return a.Printer.Raw(obj)
			}
			if a.Printer.Format == output.Wide {
				for _, p := range points {
					p["time"] = time.UnixMilli(int64(num(p, "t"))).Local().Format("15:04:05")
				}
				return a.Printer.List(obj, points, []output.Column{
					{Header: "Heure", Path: "time"},
					{Header: "CPU", Path: "cpu", Format: output.Percent},
					{Header: "RAM", Path: "ram", Format: output.MB},
				})
			}
			if len(points) == 0 {
				a.Printer.Info("Aucune mesure sur la dernière heure.")
				return nil
			}
			cpu := series(points, "cpu")
			ram := series(points, "ram")
			ramLine := fmt.Sprintf("moy %s · max %s", output.MB(avg(ram)), output.MB(maxOf(ram)))
			if lim := num(obj, "ramLimitMb"); lim > 0 {
				ramLine += " / " + output.MB(lim)
			}
			summary := map[string]any{
				"cpu":    fmt.Sprintf("moy %s · max %s", output.Percent(avg(cpu)), output.Percent(maxOf(cpu))),
				"cpuS":   Sparkline(cpu),
				"ram":    ramLine,
				"ramS":   Sparkline(ram),
				"points": float64(len(points)),
			}
			return a.Printer.Object(obj, summary, []output.Field{
				{Header: "CPU", Path: "cpu"},
				{Header: "", Path: "cpuS", Color: output.Cyan},
				{Header: "RAM", Path: "ram"},
				{Header: "", Path: "ramS", Color: output.Cyan},
				{Header: "Mesures", Path: "points"},
			})
		},
	}
}

func series(points []map[string]any, k string) []float64 {
	out := make([]float64, 0, len(points))
	for _, p := range points {
		if v, ok := p[k].(float64); ok {
			out = append(out, v)
		}
	}
	return out
}

func avg(v []float64) float64 {
	if len(v) == 0 {
		return 0
	}
	s := 0.0
	for _, x := range v {
		s += x
	}
	return s / float64(len(v))
}

func maxOf(v []float64) float64 {
	m := 0.0
	for _, x := range v {
		if x > m {
			m = x
		}
	}
	return m
}

// Sparkline : ▁▂▃▄▅▆▇█ sur 60 colonnes au plus.
func Sparkline(v []float64) string {
	if len(v) == 0 {
		return "-"
	}
	const width = 60
	if len(v) > width {
		step := float64(len(v)) / width
		s := make([]float64, width)
		for i := range s {
			s[i] = v[int(float64(i)*step)]
		}
		v = s
	}
	bars := []rune("▁▂▃▄▅▆▇█")
	lo, hi := v[0], v[0]
	for _, x := range v {
		if x < lo {
			lo = x
		}
		if x > hi {
			hi = x
		}
	}
	var b strings.Builder
	for _, x := range v {
		i := 0
		if hi > lo {
			i = int((x - lo) / (hi - lo) * float64(len(bars)-1))
		}
		b.WriteRune(bars[i])
	}
	return b.String()
}

// ------------------------------------------------------------ VPS

func vpsFields() []output.Field {
	return []output.Field{
		{Header: "ID", Path: "id"},
		{Header: "Hostname", Path: "hostname"},
		{Header: "Produit", Path: "productName"},
		{Header: "État", Path: "state", Color: output.Status},
		{Header: "Node", Path: "node"},
		{Header: "IP", Path: "ips"},
		{Header: "Utilisateur", Path: "username"},
		{Header: "CPU", Path: "cpuPct", Format: output.Percent},
		{Header: "RAM", Value: usage("ramMb", "ramMaxMb", "ramPct")},
		{Header: "Disque", Value: usage("diskMb", "diskMaxMb", "")},
		{Header: "Uptime", Path: "uptimeSeconds", Format: output.Seconds},
		{Header: "Cycle", Path: "billingCycle"},
		{Header: "Créé le", Path: "registeredAt"},
		{Header: "Expire le", Path: "expiresAt"},
		{Header: "Infos", Value: func(m map[string]any) any {
			var lines []string
			for _, f := range list(m, "fields") {
				lines = append(lines, str(f, "label")+" : "+str(f, "value"))
			}
			return strings.Join(lines, "\n")
		}, Wide: true},
	}
}

// usage : « 1.2 Go / 4 Go (30 %) ».
func usage(used, max, pct string) func(map[string]any) any {
	return func(m map[string]any) any {
		u, ok := m[used].(float64)
		if !ok {
			return nil
		}
		s := output.MB(u)
		if mx, ok := m[max].(float64); ok && mx > 0 {
			s += " / " + output.MB(mx)
			p, ok := m[pct].(float64)
			if !ok {
				p, ok = u/mx*100, true
			}
			if ok {
				s += " (" + output.Percent(p) + ")"
			}
		}
		return s
	}
}

var vps = machine{kind: "vps", resource: "vps", label: "VPS", plural: "VPS", fields: vpsFields()}

func newVPSCmd(a *App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vps",
		Short: "VPS : état, démarrage, arrêt, métriques, SSH",
		Long: `Gère vos VPS. L'ID est celui du service (yorkhost vps list) ; un nom de
service est aussi accepté, et l'ID peut être omis avec « yorkhost config set default_vps <id> ».`,
		Example: `  yorkhost vps list
  yorkhost vps get 1234
  yorkhost vps restart 1234 --wait
  yorkhost vps ssh 1234`,
	}
	cmd.AddCommand(
		vps.listCmd(a),
		vps.getCmd(a),
		vps.powerCmd(a, "start", "Démarrer un VPS", false),
		vps.powerCmd(a, "stop", "Arrêter un VPS", true),
		vps.powerCmd(a, "restart", "Redémarrer un VPS", true),
		vps.metricsCmd(a),
		newVPSSSHCmd(a),
	)
	return cmd
}

func newVPSSSHCmd(a *App) *cobra.Command {
	var user string
	var port int
	cmd := &cobra.Command{
		Use:   "ssh [id|nom] [-- commande…]",
		Short: "Se connecter en SSH à l'IP principale du VPS",
		Example: `  yorkhost vps ssh 1234
  yorkhost vps ssh 1234 --user debian -- uptime`,
		ValidArgsFunction: a.serviceCompletion("vps"),
		RunE: func(cmd *cobra.Command, args []string) error {
			var ref, rest []string
			if dash := cmd.ArgsLenAtDash(); dash >= 0 {
				ref, rest = args[:dash], args[dash:]
			} else {
				ref = args
			}
			if len(ref) > 1 {
				return usageErr(errors.New("un seul VPS ; placez la commande après --"))
			}
			id, err := a.resolveService(ctx(cmd), ref, "vps")
			if err != nil {
				return err
			}
			var obj map[string]any
			if err := a.Client().Get(ctx(cmd), path("vps", id), nil, &obj); err != nil {
				return err
			}
			ips, _ := obj["ips"].([]any)
			if len(ips) == 0 {
				return &exitError{ExitNotFound, fmt.Errorf("aucune IP connue pour le VPS %s", id)}
			}
			host := output.Text(ips[0])
			if user == "" {
				user = str(obj, "username")
			}
			if user == "" {
				user = "root"
			}
			sshArgs := []string{}
			if port != 0 && port != 22 {
				sshArgs = append(sshArgs, "-p", fmt.Sprint(port))
			}
			sshArgs = append(sshArgs, user+"@"+host)
			sshArgs = append(sshArgs, rest...)
			return a.runSSH(sshArgs)
		},
	}
	cmd.Flags().StringVarP(&user, "user", "l", "", "utilisateur (défaut : celui du VPS, sinon root)")
	cmd.Flags().IntVar(&port, "port", 22, "port SSH")
	return cmd
}

// runSSH : lance ssh en relayant le terminal ; son code de sortie est le nôtre.
func (a *App) runSSH(args []string) error {
	bin, err := exec.LookPath("ssh")
	if err != nil {
		return errors.New("client ssh introuvable dans le PATH (installez OpenSSH)")
	}
	a.Printer.Info("$ ssh %s", strings.Join(args, " "))
	c := exec.Command(bin, args...)
	c.Stdin, c.Stdout, c.Stderr = os.Stdin, os.Stdout, os.Stderr
	if err := c.Run(); err != nil {
		var ee *exec.ExitError
		if errors.As(err, &ee) {
			return silentExit(ee.ExitCode())
		}
		return err
	}
	return nil
}

// ------------------------------------------------------------ jeux

func gameFields() []output.Field {
	return []output.Field{
		{Header: "ID", Path: "id"},
		{Header: "Nom", Path: "name"},
		{Header: "Jeu", Path: "game"},
		{Header: "État", Path: "state", Color: output.Status},
		{Header: "Adresse", Value: func(m map[string]any) any {
			ad, _ := m["address"].(map[string]any)
			if ad == nil || str(ad, "ip") == "" {
				return nil
			}
			return str(ad, "ip") + ":" + str(ad, "port")
		}},
		{Header: "Ports", Value: func(m map[string]any) any {
			var ps []string
			for _, al := range list(m, "allocations") {
				s := str(al, "ip") + ":" + str(al, "port")
				if al["primary"] == true {
					s += " (principal)"
				}
				ps = append(ps, s)
			}
			return strings.Join(ps, "\n")
		}, Wide: true},
		{Header: "Node", Path: "node"},
		{Header: "CPU", Path: "cpuPct", Format: output.Percent},
		{Header: "RAM", Value: usage("ramMb", "ramMaxMb", "ramPct")},
		{Header: "Disque", Value: usage("diskMb", "diskMaxMb", "")},
		{Header: "Uptime", Path: "uptimeSeconds", Format: output.Seconds},
		{Header: "Panel", Path: "panelUrl"},
	}
}

var game = machine{kind: "game", resource: "game-servers", label: "serveur de jeu", plural: "serveurs de jeu", fields: gameFields()}

func newGameCmd(a *App) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "game",
		Aliases: []string{"game-server", "gs"},
		Short:   "Serveurs de jeu : état, démarrage, arrêt, métriques",
		Long: `Gère vos serveurs de jeu (FiveM, Minecraft, Rust, bots…). L'ID est celui du
service ; un nom est accepté, et l'ID peut être omis avec « yorkhost config set default_game <id> ».`,
		Example: `  yorkhost game list
  yorkhost game get mon-serveur
  yorkhost game restart 5678 --wait
  yorkhost game open 5678          # ouvre le panel`,
	}
	open := &cobra.Command{
		Use:               "open [id|nom]",
		Short:             "Ouvrir le panel du serveur dans le navigateur",
		Args:              noArgsOrOne,
		ValidArgsFunction: a.serviceCompletion("game"),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := a.resolveService(ctx(cmd), args, "game")
			if err != nil {
				return err
			}
			var obj map[string]any
			if err := a.Client().Get(ctx(cmd), path("game-servers", id), nil, &obj); err != nil {
				return err
			}
			u := str(obj, "panelUrl")
			if u == "" {
				return &exitError{ExitNotFound, errors.New("aucune URL de panel pour ce serveur")}
			}
			a.openURL(u)
			return nil
		},
	}
	cmd.AddCommand(
		game.listCmd(a),
		game.getCmd(a),
		game.powerCmd(a, "start", "Démarrer un serveur de jeu", false),
		game.powerCmd(a, "stop", "Arrêter un serveur de jeu", true),
		game.powerCmd(a, "restart", "Redémarrer un serveur de jeu", true),
		game.powerCmd(a, "kill", "Tuer un serveur de jeu bloqué (arrêt forcé)", true),
		game.metricsCmd(a),
		open,
	)
	return cmd
}

// openURL : ouvre le navigateur, ou affiche le lien (toujours sur stdout,
// pour qu'il reste copiable / scriptable).
func (a *App) openURL(u string) {
	if a.Printer.Machine() {
		_ = a.Printer.Raw(map[string]any{"url": u})
		return
	}
	fmt.Fprintln(a.Stdout, u)
	if a.OpenBrowser != nil && a.OpenBrowser(u) == nil {
		a.Printer.Info("Ouvert dans le navigateur.")
	}
}
