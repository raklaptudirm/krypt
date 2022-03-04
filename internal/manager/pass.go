// Copyright © 2021 Rak Laptudirm <raklaptudirm@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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

// Password gets the password data with the provided checksum.
func (p *pass) Password(hash []byte) ([]byte, error) {
	name := hex.EncodeToString(hash)
	return os.ReadFile(filepath.Join(p.Dir, name))
}

// Passwords gets all the password data in the manager.
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

// Write writes the password data into the manager.
func (p *pass) Write(data ...[]byte) error {
	for _, pass := range data {
		path := hex.EncodeToString(crypto.Checksum(pass))
		path = filepath.Join(p.Dir, path)

		err := os.WriteFile(path, pass, 0600)
		if err != nil {
			return err
		}
	}

	return nil
}

// Delete deletes the password data with the given checksum.
func (p *pass) Delete(hashes ...[]byte) error {
	for _, hash := range hashes {
		path := hex.EncodeToString(hash)
		path = filepath.Join(p.Dir, path)

		err := os.Remove(path)
		if err != nil {
			return err
		}
	}

	return nil
}
