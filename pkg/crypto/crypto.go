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

package crypto

import (
	"fmt"

	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/pbkdf2"
)

// Checksum wraps the bcrypt.GenerateFromPassword method and returns the
// data's checksum.
func PassChecksum(data []byte) []byte {
	cost := 12
	sum, _ := bcrypt.GenerateFromPassword(data, cost)
	return sum
}

func CompareChecksum(hash, pass []byte) bool {
	err := bcrypt.CompareHashAndPassword(hash, pass)
	return err == nil
}

// Checksum wraps the sha256.Sum256 method to return a []byte instead of
// a [32]byte array.
func Checksum(data []byte) []byte {
	checksum := sha256.Sum256(data)
	return checksum[:]
}

// Encrypt encrypts src with key using the AES encryption algorithm. It
// automatically generates a random salt and appends to the front of the
// ciphertext.
func Encrypt(src []byte, key []byte) (enc []byte, err error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return enc, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return
	}

	nonce := RandBytes(aesgcm.NonceSize())    // random iv
	enc = aesgcm.Seal(nonce, nonce, src, nil) // append to iv
	return
}

var ErrNoNonce = fmt.Errorf("ciphertext smaller than nonce")

// Decrypt decrypts ct with key with the AES encryption algorithm. It
// automatically extracts the salt from the ct.
func Decrypt(ct []byte, key []byte) (clt []byte, err error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return clt, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return
	}

	nonceSize := aesgcm.NonceSize()
	if len(ct) < nonceSize {
		// ciphertext can't be smaller than iv
		err = ErrNoNonce
		return
	}

	nonce, ct := ct[:nonceSize], ct[nonceSize:] // extract iv
	clt, err = aesgcm.Open(nil, nonce, ct, nil)
	return
}

// DeriveKey generates an AES algorithm key from pw, using the SHA-256 hash
// algorithm and the provided salt.
func DeriveKey(pw []byte, salt []byte) (key []byte) {
	iter := 4096 // no of pbkdf2 iterations
	klen := 32   // length of key in bytes

	algo := sha256.New // hash algorithm for pbkdf

	key = pbkdf2.Key(pw, salt, iter, klen, algo)
	return
}

// RandBytes generates an array of cryptographically secure random bytes
// with the provided length.
func RandBytes(len int) []byte {
	// generate len random bytes
	b := make([]byte, len)
	rand.Read(b)
	return b
}
