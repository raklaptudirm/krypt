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
	"os"
	"path/filepath"
)

type auth struct {
	Dir string // source directory
}

// Key gets the current key from the keyfile.
func (a *auth) Key() ([]byte, error) {
	path := filepath.Join(a.Dir, "key")
	return os.ReadFile(path)
}

// SetKey puts the provided key into the keyfile.
func (a *auth) SetKey(data []byte) error {
	path := filepath.Join(a.Dir, "key")
	return os.WriteFile(path, data, 0600)
}

// Checksum gets the current checksum from the file.
func (a *auth) Checksum() ([]byte, error) {
	path := filepath.Join(a.Dir, "checksum")
	return os.ReadFile(path)
}

// SetChecksum puts the provided checksum into the file.
func (a *auth) SetChecksum(data []byte) error {
	path := filepath.Join(a.Dir, "checksum")
	return os.WriteFile(path, data, 0600)
}

// Salt gets the current salt from the file.
func (a *auth) Salt() ([]byte, error) {
	path := filepath.Join(a.Dir, "salt")
	return os.ReadFile(path)
}

// SetSalt puts the provided salt into the file.
func (a *auth) SetSalt(data []byte) error {
	path := filepath.Join(a.Dir, "salt")
	return os.WriteFile(path, data, 0600)
}
