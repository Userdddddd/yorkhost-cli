// Package config : fichier de configuration et profils du CLI.
//
// Un seul fichier YAML (droits 0600, il contient des secrets) :
//
//	active_profile: default
//	profiles:
//	  default:
//	    api_url: https://api.yorkhost.fr
//	    output: table
//	    auth:
//	      type: oauth            # ou api_key
//	      access_token: yhk_client_v1_…
//	      refresh_token: yhr_client_v1_…
//	      expires_at: 2026-10-19T20:00:00Z
//	      client_id: yhc_…
//
// Emplacement : $YORKHOST_CONFIG_FILE, sinon $XDG_CONFIG_HOME/yorkhost/config.yaml
// (~/.config/yorkhost/config.yaml, y compris sur macOS), %AppData%\yorkhost\config.yaml
// sous Windows.
package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

const (
	DefaultAPIURL  = "https://api.yorkhost.fr"
	DefaultProfile = "default"
	DefaultOutput  = "table"
)

// Types d'authentification d'un profil.
const (
	AuthAPIKey = "api_key"
	AuthOAuth  = "oauth"
)

// Auth : identifiants d'un profil.
type Auth struct {
	Type         string    `yaml:"type,omitempty"`
	APIKey       string    `yaml:"api_key,omitempty"`
	AccessToken  string    `yaml:"access_token,omitempty"`
	RefreshToken string    `yaml:"refresh_token,omitempty"`
	ExpiresAt    time.Time `yaml:"expires_at,omitempty"`
	ClientID     string    `yaml:"client_id,omitempty"`
	RedirectURI  string    `yaml:"redirect_uri,omitempty"`
	Scopes       []string  `yaml:"scopes,omitempty"`
}

// Empty : aucun identifiant enregistré.
func (a *Auth) Empty() bool {
	return a == nil || (a.APIKey == "" && a.AccessToken == "" && a.RefreshToken == "")
}

// Profile : réglages d'un profil (un compte, une URL d'API).
type Profile struct {
	APIURL       string `yaml:"api_url,omitempty"`
	Output       string `yaml:"output,omitempty"`
	DefaultVPS   string `yaml:"default_vps,omitempty"`
	DefaultGame  string `yaml:"default_game,omitempty"`
	SSHPublicKey string `yaml:"ssh_public_key,omitempty"`
	Auth         *Auth  `yaml:"auth,omitempty"`
}

// File : contenu du fichier de configuration.
type File struct {
	ActiveProfile string              `yaml:"active_profile,omitempty"`
	Profiles      map[string]*Profile `yaml:"profiles,omitempty"`

	path string
}

// Keys : réglages modifiables par `yorkhost config set`.
var Keys = []struct{ Name, Help string }{
	{"api_url", "URL de l'API YorkHost (défaut " + DefaultAPIURL + ")"},
	{"output", "format de sortie par défaut : table, wide, json, yaml"},
	{"default_vps", "VPS utilisé quand l'ID est omis (yorkhost vps …)"},
	{"default_game", "serveur de jeu utilisé quand l'ID est omis (yorkhost game …)"},
	{"ssh_public_key", "fichier de clé publique SSH pour les sandboxes (défaut ~/.ssh/id_ed25519.pub)"},
}

var profileNameRe = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9_.-]{0,63}$`)

// ValidProfileName : nom de profil acceptable (sert aussi de clé YAML).
func ValidProfileName(name string) bool { return profileNameRe.MatchString(name) }

// Path : chemin du fichier de configuration.
func Path() (string, error) {
	if p := strings.TrimSpace(os.Getenv("YORKHOST_CONFIG_FILE")); p != "" {
		return p, nil
	}
	if runtime.GOOS == "windows" {
		dir, err := os.UserConfigDir()
		if err != nil {
			return "", err
		}
		return filepath.Join(dir, "yorkhost", "config.yaml"), nil
	}
	if x := os.Getenv("XDG_CONFIG_HOME"); x != "" {
		return filepath.Join(x, "yorkhost", "config.yaml"), nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config", "yorkhost", "config.yaml"), nil
}

// Load : lit le fichier (absent = configuration vide).
func Load() (*File, error) {
	p, err := Path()
	if err != nil {
		return nil, err
	}
	return LoadFrom(p)
}

// LoadFrom : lit un fichier donné.
func LoadFrom(p string) (*File, error) {
	f := &File{path: p, Profiles: map[string]*Profile{}}
	b, err := os.ReadFile(p)
	if errors.Is(err, os.ErrNotExist) {
		return f, nil
	}
	if err != nil {
		return nil, fmt.Errorf("lecture de %s : %w", p, err)
	}
	if err := yaml.Unmarshal(b, f); err != nil {
		return nil, fmt.Errorf("configuration illisible (%s) : %w", p, err)
	}
	if f.Profiles == nil {
		f.Profiles = map[string]*Profile{}
	}
	for name, pr := range f.Profiles {
		if pr == nil {
			f.Profiles[name] = &Profile{}
		}
	}
	f.path = p
	return f, nil
}

// FilePath : chemin d'où le fichier a été lu (et où il sera écrit).
func (f *File) FilePath() string { return f.path }

// Save : écriture atomique (fichier temporaire puis renommage), droits 0600.
func (f *File) Save() error {
	if f.path == "" {
		p, err := Path()
		if err != nil {
			return err
		}
		f.path = p
	}
	if err := os.MkdirAll(filepath.Dir(f.path), 0o700); err != nil {
		return err
	}
	b, err := yaml.Marshal(f)
	if err != nil {
		return err
	}
	header := []byte("# Configuration du CLI YorkHost — contient des secrets, ne pas partager.\n")
	tmp, err := os.CreateTemp(filepath.Dir(f.path), ".config-*.yaml")
	if err != nil {
		return err
	}
	defer os.Remove(tmp.Name())
	if err := tmp.Chmod(0o600); err != nil && runtime.GOOS != "windows" {
		tmp.Close()
		return err
	}
	if _, err := tmp.Write(append(header, b...)); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	return os.Rename(tmp.Name(), f.path)
}

// ProfileNames : noms triés.
func (f *File) ProfileNames() []string {
	names := make([]string, 0, len(f.Profiles))
	for n := range f.Profiles {
		names = append(names, n)
	}
	sort.Strings(names)
	return names
}

// Active : profil actif enregistré (DefaultProfile s'il n'y en a pas).
func (f *File) Active() string {
	if f.ActiveProfile != "" {
		return f.ActiveProfile
	}
	return DefaultProfile
}

// Profile : profil existant, nil sinon.
func (f *File) Profile(name string) *Profile { return f.Profiles[name] }

// Ensure : profil existant ou créé vide.
func (f *File) Ensure(name string) *Profile {
	p := f.Profiles[name]
	if p == nil {
		p = &Profile{}
		f.Profiles[name] = p
	}
	return p
}

// Get : valeur d'un réglage (`yorkhost config get`).
func (p *Profile) Get(key string) (string, error) {
	switch key {
	case "api_url":
		return p.APIURL, nil
	case "output":
		return p.Output, nil
	case "default_vps":
		return p.DefaultVPS, nil
	case "default_game":
		return p.DefaultGame, nil
	case "ssh_public_key":
		return p.SSHPublicKey, nil
	}
	return "", unknownKey(key)
}

// Set : modifie un réglage (valeur vide = retour au défaut).
func (p *Profile) Set(key, value string) error {
	value = strings.TrimSpace(value)
	switch key {
	case "api_url":
		if value != "" && !strings.HasPrefix(value, "https://") && !strings.HasPrefix(value, "http://") {
			return fmt.Errorf("api_url doit commencer par https://")
		}
		p.APIURL = strings.TrimRight(value, "/")
	case "output":
		if value != "" && !ValidOutput(value) {
			return fmt.Errorf("format inconnu %q (table, wide, json, yaml)", value)
		}
		p.Output = value
	case "default_vps":
		p.DefaultVPS = value
	case "default_game":
		p.DefaultGame = value
	case "ssh_public_key":
		p.SSHPublicKey = value
	default:
		return unknownKey(key)
	}
	return nil
}

// ValidOutput : formats de sortie connus.
func ValidOutput(o string) bool {
	switch o {
	case "table", "wide", "json", "yaml":
		return true
	}
	return false
}

func unknownKey(key string) error {
	names := make([]string, len(Keys))
	for i, k := range Keys {
		names[i] = k.Name
	}
	return fmt.Errorf("réglage inconnu %q (réglages : %s)", key, strings.Join(names, ", "))
}
