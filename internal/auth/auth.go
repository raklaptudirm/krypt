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
	"github.com/raklaptudirm/krypt/pkg/dir"
)

type Creds struct {
	Key  []byte
	Hash []byte
}

func (a *Creds) Validate(b []byte) bool {
	hash := crypto.Checksum(b)
	return reflect.DeepEqual(hash, a.Hash)
}

func (a *Creds) Registered() bool {
	return len(a.Hash) > 0
}

func (a *Creds) LoggedIn() bool {
	return len(a.Key) > 0
}

func Get() *Creds {
	auth := &Creds{}

	sum, err := dir.Checksum()
	if err == nil {
		auth.Hash = sum
	}

	key, err := dir.Key()
	if err == nil {
		auth.Key = key
	}

	return auth
}
