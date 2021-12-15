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

package dir

import (
	"fmt"
	"os"
)

var (
	ErrInvalidChecksum = fmt.Errorf("the checksum is invalid")
	ErrInvalidSalt     = fmt.Errorf("the salt is invalid")
	ErrInvalidKey      = fmt.Errorf("the key is invalid")
)

// file permission integers
const (
	drwxr_xr_x = 0755
	_rw_r__r__ = 0644
)

func Root() (root string, err error) {
	root, err = os.UserHomeDir()
	if err != nil {
		return
	}

	root += "/.krypt"
	err = os.MkdirAll(root, drwxr_xr_x)
	return
}

func Pass() (string, error) {
	return rootDir("passwords")
}

// key writing and fetching functions

func WriteKey(key []byte) error {
	return writeFile("key", key)
}

func Key() ([]byte, error) {
	return readFile("key", 32, ErrInvalidKey)
}

// salt writing and fetching functions

func WriteSalt(salt []byte) error {
	return writeFile("salt", salt)
}

func Salt() ([]byte, error) {
	return readFile("salt", 8, ErrInvalidSalt)
}

// checksum writing and fetching functions

func WriteChecksum(hash []byte) error {
	return writeFile("checksum", hash)
}

func Checksum() ([]byte, error) {
	return readFile("checksum", 32, ErrInvalidChecksum)
}

// rootFile gets the path of the file provided from the the krypt root
func rootFile(path string) (s string, err error) {
	s, err = Root()
	if err != nil {
		return
	}

	s += "/" + path
	return
}

// rootDir gets the path of the directory provided from the krypt root,
// and creates it if it does not exist.
func rootDir(path string) (s string, err error) {
	s, err = rootFile(path)
	if err == nil {
		err = os.MkdirAll(s, drwxr_xr_x)
	}

	return
}

// writeFile writes data to root file path
func writeFile(path string, data []byte) (err error) {
	path, err = rootFile(path)
	if err != nil {
		return
	}

	err = os.WriteFile(path, data, _rw_r__r__)
	return
}

// readFile reads data from root file path and checks if it's length is equal to
// expLen, otherwise returns lenErr
func readFile(path string, expLen int, lenErr error) (data []byte, err error) {
	path, err = rootFile(path)
	if err != nil {
		return
	}

	data, err = os.ReadFile(path)
	if len(data) != expLen {
		err = lenErr
	}

	return
}
