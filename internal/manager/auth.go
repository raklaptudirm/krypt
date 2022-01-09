package manager

import (
	"os"
	"path/filepath"
)

type auth struct {
	Dir string // source directory
}

func (a *auth) Key() ([]byte, error) {
	path := filepath.Join(a.Dir, "key")
	return os.ReadFile(path)
}

func (a *auth) SetKey(data []byte) error {
	path := filepath.Join(a.Dir, "key")
	return os.WriteFile(path, data, 0644)
}

func (a *auth) Checksum() ([]byte, error) {
	path := filepath.Join(a.Dir, "checksum")
	return os.ReadFile(path)
}

func (a *auth) SetChecksum(data []byte) error {
	path := filepath.Join(a.Dir, "checksum")
	return os.WriteFile(path, data, 0644)
}

func (a *auth) Salt() ([]byte, error) {
	path := filepath.Join(a.Dir, "salt")
	return os.ReadFile(path)
}

func (a *auth) SetSalt(data []byte) error {
	path := filepath.Join(a.Dir, "salt")
	return os.WriteFile(path, data, 0644)
}
