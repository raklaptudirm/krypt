package auth

import "github.com/raklaptudirm/krypt/pkg/dir"

type Auth struct {
	Key []byte
}

func Get() (*Auth, error) {
	key, err := dir.Key()
	if err != nil {
		return &Auth{}, err
	}

	return &Auth{
		Key: key,
	}, nil
}
