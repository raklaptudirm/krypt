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

package cmdutil

import (
	"os"

	"laptudirm.com/x/krypt/internal/auth"
	"laptudirm.com/x/krypt/internal/build"
	"laptudirm.com/x/krypt/pkg/pass"
)

// Context represents the context in which the krypt commands will
// get executed.
type Context struct {
	ExeFile     string       // the krypt executable
	Creds       *auth.Creds  // the user credentials
	Version     *Version     // the krypt version
	PassManager pass.Manager // the password manager
	AuthManager auth.Manager // the authentication manager
}

// NewContext fetches all the context information and returns a
// pointer to a Context instance.
func NewContext() *Context {
	exec, err := os.Executable()
	if err != nil {
		exec = ""
	}

	return &Context{
		ExeFile:     exec,
		Creds:       auth.Get(build.AuthManager),
		Version:     NewVersion(build.Version, build.Date),
		PassManager: build.PassManager,
		AuthManager: build.AuthManager,
	}
}
