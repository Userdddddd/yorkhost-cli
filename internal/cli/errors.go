package cli

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/YorkHost-fr/yorkhost-cli/internal/api"
	"github.com/YorkHost-fr/yorkhost-cli/internal/auth"
	"github.com/YorkHost-fr/yorkhost-cli/internal/prompt"
)

// Codes de sortie (documentés dans docs/scripting.md).
const (
	ExitOK          = 0
	ExitError       = 1 // erreur générique
	ExitUsage       = 2 // commande, option ou argument invalide
	ExitAuth        = 3 // non connecté, clé invalide, permission manquante
	ExitNotFound    = 4 // ressource introuvable
	ExitConflict    = 5 // état incompatible (409), confirmation expirée
	ExitRateLimited = 6 // 429
	ExitUnavailable = 7 // API injoignable ou en erreur (5xx)
	ExitAborted     = 8 // confirmation refusée ou délai d'attente dépassé (--wait)
	ExitInterrupted = 130
)

type usageError struct{ err error }

func (e *usageError) Error() string { return e.err.Error() }
func (e *usageError) Unwrap() error { return e.err }

func usageErr(err error) error { return &usageError{err} }

// exitError : code de sortie imposé (ex. `sandbox exec` renvoie le code de la commande).
type exitError struct {
	code int
	err  error
}

func (e *exitError) Error() string {
	if e.err == nil {
		return ""
	}
	return e.err.Error()
}

// silentExit : code de sortie sans message.
type silentExit int

func (e silentExit) Error() string { return "" }

// ExitCode : code de sortie d'une erreur.
func ExitCode(err error) int {
	var se silentExit
	if errors.As(err, &se) {
		return int(se)
	}
	var ee *exitError
	if errors.As(err, &ee) {
		return ee.code
	}
	var ue *usageError
	if errors.As(err, &ue) {
		return ExitUsage
	}
	if errors.Is(err, context.Canceled) {
		return ExitInterrupted
	}
	if errors.Is(err, prompt.ErrAborted) || errors.Is(err, prompt.ErrNoTTY) || errors.Is(err, errWaitTimeout) {
		return ExitAborted
	}
	if errors.Is(err, api.ErrNotLoggedIn) || errors.Is(err, auth.ErrInvalidGrant) {
		return ExitAuth
	}
	var ne *api.NetworkError
	if errors.As(err, &ne) {
		return ExitUnavailable
	}
	var ae *api.Error
	if errors.As(err, &ae) {
		switch {
		case ae.Status == http.StatusUnauthorized, ae.Status == http.StatusForbidden && ae.Code != "spend_cap":
			return ExitAuth
		case ae.Status == http.StatusNotFound:
			return ExitNotFound
		case ae.Status == http.StatusConflict:
			return ExitConflict
		case ae.Status == http.StatusTooManyRequests:
			return ExitRateLimited
		case ae.Status >= 500:
			return ExitUnavailable
		case ae.Status == http.StatusBadRequest || ae.Status == http.StatusUnprocessableEntity:
			return ExitUsage
		}
		return ExitError
	}
	// Erreurs cobra (commande inconnue, argument manquant…).
	msg := err.Error()
	for _, p := range []string{"unknown command", "unknown flag", "unknown shorthand", "accepts ", "requires at least", "requires at most", "invalid argument", "required flag"} {
		if strings.Contains(msg, p) {
			return ExitUsage
		}
	}
	return ExitError
}
