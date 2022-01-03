package crypto

import (
	"fmt"
	"reflect"
	"testing"
)

func TestChecksum(t *testing.T) {
	// TODO: add test cases
	tc := []struct {
		data []byte
		hash []byte
	}{}

	for _, c := range tc {
		name := fmt.Sprintf("check hash for %x", c.data)
		t.Run(name, func(t *testing.T) {
			hash := Checksum(c.data)
			if !reflect.DeepEqual(hash, c.hash) {
				t.Errorf("Checksum(): want %x, got %x", c.hash, hash)
			}
		})
	}
}

func TestEncryptDecrypt(t *testing.T) {
	for i := 0; i < 5; i++ { // 5 is arbitrary
		key := RandBytes(32)
		src := RandBytes(32)

		name := fmt.Sprintf("test encrypt-decrypt for %x", src)

		t.Run(name, func(t *testing.T) {
			enc, err := Encrypt(src, key)
			if err != nil {
				t.Errorf("Encrypt(): %v", err)
			}

			clt, err := Decrypt(enc, key)
			if err != nil {
				t.Errorf("Encrypt(): %v", err)
			}

			if !reflect.DeepEqual(clt, src) {
				t.Errorf("Encrypt() or Decrypt(): want %x, got %x", src, clt)
			}
		})
	}
}

func TestDeriveKey(t *testing.T) {
	// TODO: add test cases
	tc := []struct {
		data []byte
		salt []byte
		key  []byte
	}{}

	for _, c := range tc {
		name := fmt.Sprintf("check key derivation for %x[%x]", c.data, c.salt)
		t.Run(name, func(t *testing.T) {
			key := DeriveKey(c.data, c.salt)
			if !reflect.DeepEqual(key, c.key) {
				t.Errorf("Checksum(): want %x, got %x", c.key, key)
			}
		})
	}
}
