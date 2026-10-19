// Package prompt : questions interactives (confirmation, saisie, secret).
// Hors terminal, aucune question n'est posée : une action à confirmer exige
// --yes, pour qu'un script ne reste jamais bloqué sur stdin.
package prompt

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"golang.org/x/term"
)

// ErrNoTTY : confirmation impossible sans terminal.
var ErrNoTTY = errors.New("confirmation requise : relancez avec --yes (pas de terminal interactif)")

// ErrAborted : l'utilisateur a refusé.
var ErrAborted = errors.New("annulé")

// Prompter : entrée/sortie des questions (remplaçable dans les tests).
type Prompter struct {
	In          io.Reader
	Out         io.Writer
	Interactive bool
	reader      *bufio.Reader
}

// Default : stdin/stderr, interactif si les deux sont des terminaux.
func Default() *Prompter {
	return &Prompter{
		In: os.Stdin, Out: os.Stderr,
		Interactive: term.IsTerminal(int(os.Stdin.Fd())) && term.IsTerminal(int(os.Stderr.Fd())),
	}
}

func (p *Prompter) line() (string, error) {
	if p.reader == nil {
		p.reader = bufio.NewReader(p.In)
	}
	s, err := p.reader.ReadString('\n')
	// EOF après une saisie sans retour à la ligne : la saisie compte.
	if err != nil && (!errors.Is(err, io.EOF) || s == "") {
		return "", err
	}
	return strings.TrimSpace(s), nil
}

// Confirm : question oui/non (défaut non). yes = --yes déjà donné.
func (p *Prompter) Confirm(question string, yes bool) error {
	if yes {
		return nil
	}
	if !p.Interactive {
		return ErrNoTTY
	}
	fmt.Fprintf(p.Out, "%s [o/N] ", question)
	s, err := p.line()
	if err != nil {
		return ErrAborted
	}
	switch strings.ToLower(s) {
	case "o", "oui", "y", "yes":
		return nil
	}
	return ErrAborted
}

// ConfirmTyped : l'utilisateur doit retaper expected (actions irréversibles).
func (p *Prompter) ConfirmTyped(question, expected string, yes bool) error {
	if yes {
		return nil
	}
	if !p.Interactive {
		return ErrNoTTY
	}
	fmt.Fprintf(p.Out, "%s\nTapez %q pour confirmer : ", question, expected)
	s, err := p.line()
	if err != nil || s != expected {
		return ErrAborted
	}
	return nil
}

// Input : saisie libre (def proposé si vide).
func (p *Prompter) Input(label, def string) (string, error) {
	if !p.Interactive {
		if def != "" {
			return def, nil
		}
		return "", fmt.Errorf("%s : valeur requise (pas de terminal interactif)", label)
	}
	if def != "" {
		fmt.Fprintf(p.Out, "%s [%s] : ", label, def)
	} else {
		fmt.Fprintf(p.Out, "%s : ", label)
	}
	s, err := p.line()
	if err != nil {
		return "", err
	}
	if s == "" {
		return def, nil
	}
	return s, nil
}

// Secret : saisie masquée (clé API, valeur de secret).
func (p *Prompter) Secret(label string) (string, error) {
	if !p.Interactive {
		return "", fmt.Errorf("%s : saisie interactive impossible, passez la valeur par stdin", label)
	}
	fmt.Fprintf(p.Out, "%s : ", label)
	if f, ok := p.In.(*os.File); ok && term.IsTerminal(int(f.Fd())) {
		b, err := term.ReadPassword(int(f.Fd()))
		fmt.Fprintln(p.Out)
		if err != nil {
			return "", err
		}
		return strings.TrimSpace(string(b)), nil
	}
	return p.line()
}
