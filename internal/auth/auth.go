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

package auth

import (
	"reflect"

	"github.com/raklaptudirm/krypt/pkg/crypto"
)

// Creds stored authentication information that can be used to validate
// the master password and decrypt stored passwords.
type Creds struct {
	Key  []byte
	Hash []byte
}

// Validate checks if the provided byte array matches the master password
// checksum.
func (a *Creds) Validate(b []byte) bool {
	hash := crypto.Checksum(b)
	return reflect.DeepEqual(hash, a.Hash)
}

// Registered checks if the user is registered from the credentials.
func (a *Creds) Registered() bool {
	// password hash is non zero
	return len(a.Hash) > 0
}

// LoggedIn checks if the user is logged in from the credentials.
func (a *Creds) LoggedIn() bool {
	// key value is non zero
	return len(a.Key) > 0
}

// Get returns a pointer to a Creds instance with the credential information
// from the provided auth manager.
func Get(man Manager) *Creds {
	auth := &Creds{}

	sum, err := man.Checksum()
	if err == nil {
		auth.Hash = sum
	}

	key, err := man.Key()
	if err == nil {
		auth.Key = key
	}

	return auth
}
