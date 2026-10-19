// Package version : version du binaire, injectée au build par goreleaser
// (-ldflags "-X …/version.Version=… -X …/version.Commit=… -X …/version.Date=…").
package version

import (
	"fmt"
	"runtime"
	"runtime/debug"
)

var (
	Version = "dev"
	Commit  = ""
	Date    = ""
)

func init() {
	// `go install …@v1.2.3` : pas de ldflags, mais la version du module est
	// connue du runtime.
	if Version != "dev" {
		return
	}
	if info, ok := debug.ReadBuildInfo(); ok {
		if v := info.Main.Version; v != "" && v != "(devel)" {
			Version = v
		}
		for _, s := range info.Settings {
			switch s.Key {
			case "vcs.revision":
				if Commit == "" {
					Commit = s.Value
				}
			case "vcs.time":
				if Date == "" {
					Date = s.Value
				}
			}
		}
	}
}

// UserAgent : en-tête envoyé à l'API (journal client, statistiques).
func UserAgent() string {
	return fmt.Sprintf("yorkhost-cli/%s (%s/%s)", Version, runtime.GOOS, runtime.GOARCH)
}

// String : ligne affichée par `yorkhost version`.
func String() string {
	s := "yorkhost " + Version
	if Commit != "" {
		c := Commit
		if len(c) > 12 {
			c = c[:12]
		}
		s += " (" + c
		if Date != "" {
			s += ", " + Date
		}
		s += ")"
	}
	return s + fmt.Sprintf(" %s/%s %s", runtime.GOOS, runtime.GOARCH, runtime.Version())
}
