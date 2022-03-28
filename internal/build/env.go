package build

import (
	"laptudirm.com/x/krypt/internal/auth"
	"laptudirm.com/x/krypt/internal/manager"
	"laptudirm.com/x/krypt/pkg/pass"
)

// values of the managers can be manually changed
// to make krypt work in any environment.
var PassManager pass.Manager = manager.Pass
var AuthManager auth.Manager = manager.Auth
