package output

import "strings"

func wrap(code, s string) string { return "\x1b[" + code + "m" + s + "\x1b[0m" }

func Bold(s string) string   { return wrap("1", s) }
func Dim(s string) string    { return wrap("2", s) }
func Red(s string) string    { return wrap("31", s) }
func Green(s string) string  { return wrap("32", s) }
func Yellow(s string) string { return wrap("33", s) }
func Cyan(s string) string   { return wrap("36", s) }

// Status : couleur d'un statut (service, VPS, serveur, facture, sandbox).
func Status(s string) string {
	switch strings.ToLower(s) {
	case "active", "running", "paid", "oui", "open", "confirmed", "ok":
		return Green(s)
	case "suspended", "offline", "stopped", "unpaid", "expired", "failed", "destroyed", "cancelled", "non":
		return Red(s)
	case "pending", "starting", "stopping", "provisioning", "claim_pending", "transferring", "refunded", "possible", "unknown":
		return Yellow(s)
	}
	return s
}
