package auth

import (
	"testing"
)

func TestValidate(t *testing.T) {
	// TODO: write test cases
	tc := []struct {
		name  string
		creds *Creds
		data  []byte
		equal bool
	}{}

	for _, c := range tc {
		t.Run(c.name, func(t *testing.T) {
			equal := c.creds.Validate(c.data)
			if equal != c.equal {
				t.Errorf("(*Creds).Validate(): want %v, got %v", c.equal, equal)
			}
		})
	}
}

var (
	emptyBytes    = []byte{}
	nonEmptyBytes = []byte{1}
)

var testCases = []struct {
	name       string
	creds      *Creds
	registered bool
	loggedIn   bool
}{
	{
		name: "test for empty key and hash",
		creds: &Creds{
			Key:  emptyBytes,
			Hash: emptyBytes,
		},
		registered: false,
		loggedIn:   false,
	},
	{
		name: "test for non-empty key and hash",
		creds: &Creds{
			Key:  nonEmptyBytes,
			Hash: nonEmptyBytes,
		},
		registered: true,
		loggedIn:   true,
	},
}

func TestRegistered(t *testing.T) {
	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			registered := c.creds.Registered()
			if registered != c.registered {
				t.Errorf("(*Creds).Registered(): want %v, got %v", c.registered, registered)
			}
		})
	}
}

func TestLoggedIn(t *testing.T) {
	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			loggedIn := c.creds.LoggedIn()
			if loggedIn != c.loggedIn {
				t.Errorf("(*Creds).LoggedIn(): want %v, got %v", c.loggedIn, loggedIn)
			}
		})
	}
}
