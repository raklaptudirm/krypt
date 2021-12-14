package crypto

import (
	"fmt"
	"sync"
	"time"

	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"math/rand"

	"github.com/raklaptudirm/krypt/pkg/dir"
	"golang.org/x/crypto/pbkdf2"
)

func EncryptWithKey(src []byte) (enc []byte, err error) {
	key, err := dir.Key()
	if err != nil {
		return
	}

	enc, err = Encrypt(src, key)
	return
}

func DecryptWithKey(ct []byte) (clt []byte, err error) {
	key, err := dir.Key()
	if err != nil {
		return
	}

	clt, err = Decrypt(ct, key)
	return
}

func Sha256(data []byte) []byte {
	checksum := sha256.Sum256(data)
	return checksum[:]
}

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
		err = fmt.Errorf("ciphertext smaller than nonce")
		return
	}

	nonce, ct := ct[:nonceSize], ct[nonceSize:] // extract iv
	clt, err = aesgcm.Open(nil, nonce, ct, nil)
	return
}

func Pbkdf2(pw []byte, salt []byte) (key []byte) {
	iter := 4096 // no of pbkdf2 iterations
	klen := 32   // length of key in bytes

	algo := sha256.New // hash algorithm for pbkdf

	key = pbkdf2.Key(pw, salt, iter, klen, algo)
	return
}

func HashString(data []byte) string {
	return fmt.Sprintf("%x", data)
}

var setSeed sync.Once

func RandBytes(len int) []byte {
	// set rand seed once
	setSeed.Do(func() {
		rand.Seed(time.Now().UnixNano())
	})

	// generate len random bytes
	b := make([]byte, len)
	rand.Read(b)
	return b
}
