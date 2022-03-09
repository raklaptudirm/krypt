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

package pass

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"

	"github.com/raklaptudirm/krypt/pkg/crypto"
)

// ErrDecode represents some data which was not able to be decoded
// by decode.
type ErrDecode struct {
	ct []byte // ciphertext
}

// Error represents ErrDecode in a string.
func (e ErrDecode) Error() string {
	return fmt.Sprintf("Error decoding %x", e.ct)
}

// Password represents a krypt password.
type Password struct {
	Name     string // password name
	UserID   string // password username
	Password string // password
	Checksum []byte // checksum of encrypted password
}

func (p *Password) encode(key []byte) (b []byte, err error) {
	b = append([]byte(p.Name), '\n')
	b = append(b, []byte(p.UserID)...)
	b = append(b, '\n')
	b = append(b, []byte(p.Password)...)
	b, err = crypto.Encrypt(b, key)
	return
}

func decode(b []byte, key []byte) (pass *Password, err error) {
	hash := crypto.Checksum(b)
	b, err = crypto.Decrypt(b, key)
	if err != nil {
		return
	}

	lines := bytes.Split(b, []byte{'\n'})
	if len(lines) != 3 {
		err = ErrDecode{b}
		return
	}

	pass = &Password{
		Name:     string(lines[0]),
		UserID:   string(lines[1]),
		Password: string(lines[2]),
		Checksum: hash,
	}
	return
}

// String represents Password as a string.
func (p *Password) String() string {
	hidden := strings.Repeat("*", len(p.Password))
	return fmt.Sprintf("Name: %v\nUsername: %v\nPassword: %v", p.Name, p.UserID, hidden)
}

// Write encrypts the password with the provided key and writes it to the provided manager.
func (p *Password) Write(man Manager, key []byte) error {
	data, err := p.encode(key)
	if err != nil {
		return err
	}

	return man.Write(data)
}

func Get(man Manager, key []byte) ([]Password, error) {
	dataset, err := man.Passwords()
	if err != nil {
		return nil, err
	}

	var passwords []Password

	for _, data := range dataset {
		password, err := decode(data, key)
		if err == nil {
			passwords = append(passwords, *password)
		}
	}

	if len(passwords) == 0 {
		return nil, fmt.Errorf("no passwords in database")
	}

	return passwords, nil
}

// Get fetches a list of password from the provided password manager whose names match
// the provided regular expression.
func Filter(man Manager, key []byte, ident string) ([]Password, error) {
	regex, err := regexp.Compile(ident)
	if err != nil {
		return nil, err
	}

	var filtered []Password

	passwords, err := Get(man, key)
	if err != nil {
		return nil, err
	}

	for _, password := range passwords {
		if regex.Match([]byte(password.Name)) {
			filtered = append(passwords, password)
		}
	}

	if len(filtered) == 0 {
		return nil, fmt.Errorf("no passwords matched %#v", ident)
	}

	return filtered, nil
}
