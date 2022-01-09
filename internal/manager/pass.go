package manager

import (
	"encoding/hex"
	"os"
	"path/filepath"

	"github.com/raklaptudirm/krypt/pkg/crypto"
)

type pass struct {
	Dir string // source directory
}

func (p *pass) Password(hash []byte) ([]byte, error) {
	name := hex.EncodeToString(hash)
	return os.ReadFile(filepath.Join(p.Dir, name))
}

func (p *pass) Passwords() ([][]byte, error) {
	var passwords [][]byte

	elements, err := os.ReadDir(p.Dir)
	if err != nil {
		return [][]byte{}, err
	}

	for _, element := range elements {
		if element.IsDir() {
			continue
		}

		path := filepath.Join(p.Dir, element.Name())
		data, err := os.ReadFile(path)
		if err != nil {
			return [][]byte{}, err
		}

		passwords = append(passwords, data)
	}

	return passwords, nil
}

func (p *pass) Write(data []byte) error {
	path := hex.EncodeToString(crypto.Checksum(data))
	path = filepath.Join(p.Dir, path)

	return os.WriteFile(path, data, 0644)
}

func (p *pass) Delete(hash []byte) error {
	path := hex.EncodeToString(hash)
	path = filepath.Join(p.Dir, path)

	return os.Remove(path)
}
