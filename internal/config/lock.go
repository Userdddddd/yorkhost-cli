package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// staleLock : un verrou plus vieux est abandonné (process tué en plein
// rafraîchissement).
const staleLock = 30 * time.Second

// Lock : verrou exclusif sur le fichier de configuration, pour qu'un seul
// process rafraîchisse le jeton OAuth (les refresh tokens sont à usage
// unique : deux rafraîchissements concurrents feraient échouer le second).
// Fichier créé en O_EXCL, portable sans flock.
func (f *File) Lock(timeout time.Duration) (unlock func(), err error) {
	if f.path == "" {
		p, err := Path()
		if err != nil {
			return nil, err
		}
		f.path = p
	}
	if err := os.MkdirAll(filepath.Dir(f.path), 0o700); err != nil {
		return nil, err
	}
	lp := f.path + ".lock"
	deadline := time.Now().Add(timeout)
	for {
		fh, err := os.OpenFile(lp, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o600)
		if err == nil {
			fmt.Fprintf(fh, "%d\n", os.Getpid())
			fh.Close()
			return func() { os.Remove(lp) }, nil
		}
		if !errors.Is(err, os.ErrExist) {
			return nil, err
		}
		if st, serr := os.Stat(lp); serr == nil && time.Since(st.ModTime()) > staleLock {
			os.Remove(lp)
			continue
		}
		if time.Now().After(deadline) {
			return nil, fmt.Errorf("configuration verrouillée par un autre process (%s)", lp)
		}
		time.Sleep(100 * time.Millisecond)
	}
}
