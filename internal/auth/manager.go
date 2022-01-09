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

// Manager represents an authentication manager
type Manager interface {
	// getter and setter for the encryption key
	Key() ([]byte, error)
	SetKey([]byte) error

	// getter and setter for the master password checksum
	Checksum() ([]byte, error)
	SetChecksum([]byte) error

	// getter and setter from the key derivation salt
	Salt() ([]byte, error)
	SetSalt([]byte) error
}
