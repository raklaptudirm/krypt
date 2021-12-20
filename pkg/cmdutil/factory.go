package cmdutil

import (
	"github.com/raklaptudirm/krypt/internal/auth"
)

type Factory struct {
	Executable string
	Auth       *auth.Auth
}
