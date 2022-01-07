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
	"strings"

	"github.com/raklaptudirm/krypt/pkg/crypto"
)

type ErrDecode struct {
	ct []byte
}

func (e ErrDecode) Error() string {
	return fmt.Sprintf("Error decoding %v", e.ct)
}

type Password struct {
	Name     string
	UserID   string
	Password string
	Checksum []byte
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

func (p *Password) String() string {
	hidden := strings.Repeat("*", len(p.Password))
	return fmt.Sprintf("Name: %v\nUsername: %v\nPassword: %v\n", p.Name, p.UserID, hidden)
}

func (p *Password) Write(man Manager, key []byte) error {
	data, err := p.encode(key)
	if err != nil {
		return err
	}

	err = man.Write(data)
	return err
}

func GetS(man Manager, ident string, key []byte) (pass *Password, err error) {
	pbs, err := man.Passwords()
	if err != nil {
		return
	}

	for _, pb := range pbs {
		hash := hex.EncodeToString((crypto.Checksum(pb)))

		if strings.HasPrefix(hash, ident) {
			pass, err = decode(pb, key)
			return
		}

		pass, err = decode(pb, key)
		if err == nil && strings.Contains(pass.Name, ident) {
			return
		}
	}

	err = fmt.Errorf("no password matched %v", ident)
	return
}
