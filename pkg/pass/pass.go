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
	"encoding/hex"
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
	return fmt.Sprintf("Name: %v\nUsername: %v\nPassword: %v\n", p.Name, p.UserID, hidden)
}

// Write encrypts the password with the provided key and writes it to the provided manager.
func (p *Password) Write(man Manager, key []byte) error {
	data, err := p.encode(key)
	if err != nil {
		return err
	}

	return man.Write(data)
}

// GetS fetches a single password from the provided password manager according to
// the provided identifier.
func GetS(man Manager, ident string, key []byte) (pass *Password, err error) {
	ident = strings.ToLower(ident)

	pbs, err := man.Passwords()
	if err != nil {
		return
	}

	for _, pb := range pbs {
		hash := hex.EncodeToString((crypto.Checksum(pb)))

		// check if string matches checksum
		if strings.HasPrefix(hash, ident) {
			return decode(pb, key)
		}

		pass, err = decode(pb, key)
		// check if string matches password name
		if err == nil && strings.Contains(strings.ToLower(pass.Name), ident) {
			return
		}
	}

	err = fmt.Errorf("no password matched %v", ident)
	return
}

func Get(man Manager, ident string, key []byte) (pass []Password, err error) {
	regex, err := regexp.Compile(ident)
	if err != nil {
		return
	}

	pbs, err := man.Passwords()
	if err != nil {
		return
	}

	for _, pb := range pbs {
		p, err := decode(pb, key)
		// check if string matches password name
		if err == nil && regex.Match([]byte(p.Name)) {
			pass = append(pass, *p)
		}
	}

	err = fmt.Errorf("no password matched %v", ident)
	return
}
