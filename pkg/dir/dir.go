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

func Root() (root string, err error) {
	root, err = os.UserHomeDir()
	if err != nil {
		return
	}

	root += "/.krypt"
	err = os.MkdirAll(root, 0755)
	return
}

func Pass() (string, error) {
	return rootDir("passwords")
}

func KeyFile() (string, error) {
	return rootFile("key")
}

func WriteKey(key []byte) (err error) {
	path, err := KeyFile()
	if err != nil {
		return
	}

	os.WriteFile(path, key, 0644)
	return
}

func Key() (key []byte, err error) {
	path, err := KeyFile()
	if err != nil {
		return
	}

	key, err = os.ReadFile(path)
	if len(key) != 32 {
		err = ErrInvalidKey
	}

	return
}

func WriteSalt(salt []byte) (err error) {
	path, err := rootFile("salt")
	if err != nil {
		return
	}

	err = os.WriteFile(path, salt, 0644)
	return
}

func Salt() (salt []byte, err error) {
	path, err := rootFile("salt")
	if err != nil {
		return
	}

	salt, err = os.ReadFile(path)
	if len(salt) != 8 {
		err = ErrInvalidSalt
	}

	return
}

func WriteChecksum(hash []byte) (err error) {
	path, err := rootFile("checksum")
	if err != nil {
		return
	}

	err = os.WriteFile(path, hash, 0644)
	return
}

func Checksum() (checksum []byte, err error) {
	path, err := rootFile("checksum")
	if err != nil {
		return
	}

	checksum, err = os.ReadFile(path)
	if len(checksum) != 32 {
		err = ErrInvalidChecksum
	}

	return
}

func rootFile(path string) (s string, err error) {
	s, err = Root()
	if err != nil {
		return
	}

	s += "/" + path
	return
}

func rootDir(path string) (s string, err error) {
	s, err = rootFile(path)
	if err == nil {
		err = os.MkdirAll(s, 0755)
	}

	return
}
