// Package output : rendu des résultats (table, wide, json, yaml) et couleurs.
//
// Règle : stdout porte le résultat (pipeable : `-o json | jq`), stderr les
// messages de progression et les erreurs. En json/yaml, la réponse de l'API
// est imprimée telle quelle — le contrat est celui de l'API, pas du CLI.
package output

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"golang.org/x/term"
	"gopkg.in/yaml.v3"
)

// Formats.
const (
	Table = "table"
	Wide  = "wide"
	JSON  = "json"
	YAML  = "yaml"
)

// Printer : destination et format.
type Printer struct {
	Format string
	Out    io.Writer
	Err    io.Writer
	Color  bool
	Quiet  bool
}

// New : couleurs si stdout est un terminal et ni NO_COLOR ni --no-color.
func New(format string, noColor bool) *Printer {
	color := !noColor && os.Getenv("NO_COLOR") == "" && os.Getenv("TERM") != "dumb" && term.IsTerminal(int(os.Stdout.Fd()))
	return &Printer{Format: format, Out: os.Stdout, Err: os.Stderr, Color: color}
}

// Machine : sortie destinée à un programme (json/yaml).
func (p *Printer) Machine() bool { return p.Format == JSON || p.Format == YAML }

// Column : colonne d'une table.
type Column struct {
	Header string
	// Path : chemin pointé dans l'objet (ex. "address.ip").
	Path string
	// Value : calcul libre (prioritaire sur Path).
	Value func(map[string]any) any
	// Format : mise en forme de la valeur (défaut : Text).
	Format func(any) string
	// Wide : colonne affichée seulement en -o wide.
	Wide bool
	// Max : troncature en mode table (0 = aucune).
	Max int
	// Color : couleur selon la valeur affichée.
	Color func(string) string
}

// Raw : json/yaml de v.
func (p *Printer) Raw(v any) error {
	switch p.Format {
	case YAML:
		// Passage par JSON : clés et types identiques à la sortie json.
		b, err := json.Marshal(v)
		if err != nil {
			return err
		}
		var generic any
		if err := json.Unmarshal(b, &generic); err != nil {
			return err
		}
		enc := yaml.NewEncoder(p.Out)
		enc.SetIndent(2)
		defer enc.Close()
		return enc.Encode(generic)
	default:
		enc := json.NewEncoder(p.Out)
		enc.SetIndent("", "  ")
		enc.SetEscapeHTML(false)
		return enc.Encode(v)
	}
}

// List : table des items, ou raw (json/yaml).
func (p *Printer) List(raw any, items []map[string]any, cols []Column) error {
	if p.Machine() {
		return p.Raw(raw)
	}
	if len(items) == 0 {
		p.Info("Aucun résultat.")
		return nil
	}
	var shown []Column
	for _, c := range cols {
		if !c.Wide || p.Format == Wide {
			shown = append(shown, c)
		}
	}
	rows := make([][]string, len(items))
	plain := make([][]string, len(items))
	widths := make([]int, len(shown))
	for i, c := range shown {
		widths[i] = utf8.RuneCountInString(c.Header)
	}
	for r, it := range items {
		rows[r] = make([]string, len(shown))
		plain[r] = make([]string, len(shown))
		for i, c := range shown {
			s := cell(it, c)
			if c.Max > 0 && p.Format != Wide {
				s = truncate(s, c.Max)
			}
			plain[r][i] = s
			if p.Color && c.Color != nil {
				rows[r][i] = c.Color(s)
			} else {
				rows[r][i] = s
			}
			if w := utf8.RuneCountInString(s); w > widths[i] {
				widths[i] = w
			}
		}
	}
	var b strings.Builder
	for i, c := range shown {
		h := strings.ToUpper(c.Header)
		if p.Color {
			h = Bold(h)
		}
		writeCell(&b, h, utf8.RuneCountInString(c.Header), widths[i], i == len(shown)-1)
	}
	b.WriteByte('\n')
	for r := range rows {
		for i := range shown {
			writeCell(&b, rows[r][i], utf8.RuneCountInString(plain[r][i]), widths[i], i == len(shown)-1)
		}
		b.WriteByte('\n')
	}
	_, err := io.WriteString(p.Out, b.String())
	return err
}

func writeCell(b *strings.Builder, s string, visible, width int, last bool) {
	b.WriteString(s)
	if !last {
		b.WriteString(strings.Repeat(" ", width-visible+3))
	}
}

func cell(it map[string]any, c Column) string {
	var v any
	if c.Value != nil {
		v = c.Value(it)
	} else {
		v = Lookup(it, c.Path)
	}
	if c.Format != nil {
		return c.Format(v)
	}
	return Text(v)
}

// Field : ligne d'une fiche détaillée.
type Field = Column

// Object : fiche « Clé : valeur », ou raw (json/yaml). Les champs dont la
// valeur est vide sont omis, sauf en -o wide.
func (p *Printer) Object(raw any, obj map[string]any, fields []Field) error {
	if p.Machine() {
		return p.Raw(raw)
	}
	width := 0
	for _, f := range fields {
		if w := utf8.RuneCountInString(f.Header); w > width {
			width = w
		}
	}
	var b strings.Builder
	for _, f := range fields {
		if f.Wide && p.Format != Wide {
			continue
		}
		s := cell(obj, f)
		if s == "-" && p.Format != Wide {
			continue
		}
		label := f.Header
		pad := strings.Repeat(" ", width-utf8.RuneCountInString(label))
		if p.Color {
			label = Dim(label)
			if f.Color != nil {
				s = f.Color(s)
			}
		}
		// Valeurs multilignes : indentées sous le libellé.
		s = strings.ReplaceAll(s, "\n", "\n"+strings.Repeat(" ", width+3))
		fmt.Fprintf(&b, "%s%s   %s\n", label, pad, s)
	}
	_, err := io.WriteString(p.Out, b.String())
	return err
}

// Info : message de progression (stderr, muet avec --quiet).
func (p *Printer) Info(format string, args ...any) {
	if p.Quiet {
		return
	}
	fmt.Fprintf(p.Err, format+"\n", args...)
}

// Success : confirmation d'une action (stdout en mode humain, rien en
// json/yaml où la commande imprime son propre résultat).
func (p *Printer) Success(format string, args ...any) {
	if p.Machine() || p.Quiet {
		return
	}
	msg := fmt.Sprintf(format, args...)
	if p.Color {
		msg = Green("✓") + " " + msg
	} else {
		msg = "✓ " + msg
	}
	fmt.Fprintln(p.Out, msg)
}

// Warn : avertissement (stderr).
func (p *Printer) Warn(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	if p.Color {
		msg = Yellow("!") + " " + msg
	} else {
		msg = "! " + msg
	}
	fmt.Fprintln(p.Err, msg)
}

// Lookup : valeur au chemin "a.b.c" (nil si absent).
func Lookup(m map[string]any, path string) any {
	if path == "" {
		return nil
	}
	var cur any = m
	for _, part := range strings.Split(path, ".") {
		mm, ok := cur.(map[string]any)
		if !ok {
			return nil
		}
		cur = mm[part]
	}
	return cur
}

// ------------------------------------------------------------ formats

// Text : rendu générique d'une valeur JSON.
func Text(v any) string {
	switch t := v.(type) {
	case nil:
		return "-"
	case string:
		if t == "" {
			return "-"
		}
		return t
	case bool:
		if t {
			return "oui"
		}
		return "non"
	case float64:
		return Number(t)
	case int:
		return strconv.Itoa(t)
	case json.Number:
		return t.String()
	case []any:
		if len(t) == 0 {
			return "-"
		}
		parts := make([]string, len(t))
		for i, x := range t {
			parts[i] = Text(x)
		}
		return strings.Join(parts, ", ")
	case map[string]any:
		keys := make([]string, 0, len(t))
		for k := range t {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		parts := make([]string, len(keys))
		for i, k := range keys {
			parts[i] = k + "=" + Text(t[k])
		}
		return strings.Join(parts, " ")
	}
	return fmt.Sprint(v)
}

// Number : sans décimales inutiles.
func Number(f float64) string {
	if f == math.Trunc(f) && math.Abs(f) < 1e15 {
		return strconv.FormatInt(int64(f), 10)
	}
	return strconv.FormatFloat(f, 'f', 2, 64)
}

// Money : montant en euros (ou devise donnée).
func Money(currency string) func(any) string {
	return func(v any) string {
		f, ok := v.(float64)
		if !ok {
			return Text(v)
		}
		cur := currency
		if cur == "" || cur == "EUR" {
			cur = "€"
		}
		return strings.Replace(strconv.FormatFloat(f, 'f', 2, 64), ".", ",", 1) + " " + cur
	}
}

// Percent : 12.3 → "12.3 %".
func Percent(v any) string {
	f, ok := v.(float64)
	if !ok {
		return Text(v)
	}
	return strconv.FormatFloat(f, 'f', 1, 64) + " %"
}

// MB : mégaoctets lisibles.
func MB(v any) string {
	f, ok := v.(float64)
	if !ok {
		return Text(v)
	}
	if f >= 1024 {
		return strconv.FormatFloat(f/1024, 'f', 1, 64) + " Go"
	}
	return Number(f) + " Mo"
}

// Bytes : octets lisibles.
func Bytes(v any) string {
	f, ok := v.(float64)
	if !ok {
		return Text(v)
	}
	units := []string{"o", "Kio", "Mio", "Gio", "Tio"}
	i := 0
	for f >= 1024 && i < len(units)-1 {
		f /= 1024
		i++
	}
	if i == 0 {
		return Number(f) + " " + units[i]
	}
	return strconv.FormatFloat(f, 'f', 1, 64) + " " + units[i]
}

// Seconds : durée lisible (uptime).
func Seconds(v any) string {
	f, ok := v.(float64)
	if !ok {
		return Text(v)
	}
	return Duration(time.Duration(f) * time.Second)
}

// Duration : 3j 4h 12min.
func Duration(d time.Duration) string {
	if d < time.Minute {
		return fmt.Sprintf("%ds", int(d.Seconds()))
	}
	days := int(d.Hours()) / 24
	h := int(d.Hours()) % 24
	m := int(d.Minutes()) % 60
	switch {
	case days > 0:
		return fmt.Sprintf("%dj %dh", days, h)
	case h > 0:
		return fmt.Sprintf("%dh %02dmin", h, m)
	}
	return fmt.Sprintf("%dmin", m)
}

// Date : RFC 3339 → "2006-01-02 15:04" (heure locale) ; autre texte inchangé.
func Date(v any) string {
	s, ok := v.(string)
	if !ok || s == "" {
		return Text(v)
	}
	for _, layout := range []string{time.RFC3339Nano, time.RFC3339, "2006-01-02 15:04:05"} {
		if t, err := time.Parse(layout, s); err == nil {
			return t.Local().Format("2006-01-02 15:04")
		}
	}
	return s
}

// Relative : échéance relative (« dans 1h 20min », « il y a 3min »).
func Relative(v any) string {
	s, ok := v.(string)
	if !ok || s == "" {
		return Text(v)
	}
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return s
	}
	d := time.Until(t)
	if d >= 0 {
		return "dans " + Duration(d)
	}
	return "il y a " + Duration(-d)
}

func truncate(s string, max int) string {
	s = strings.ReplaceAll(s, "\n", " ")
	if utf8.RuneCountInString(s) <= max {
		return s
	}
	r := []rune(s)
	return string(r[:max-1]) + "…"
}
