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

package login

import (
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/raklaptudirm/krypt/internal/auth"
	"github.com/raklaptudirm/krypt/pkg/cmdutil"
	"github.com/raklaptudirm/krypt/pkg/crypto"
	"github.com/raklaptudirm/krypt/pkg/dir"
	"github.com/raklaptudirm/krypt/pkg/term"
)

type LoginOptions struct {
	Creds *auth.Creds
}

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &LoginOptions{
		Creds: f.Creds,
	}

	var cmd = &cobra.Command{
		Use:   "login",
		Short: "login to krypt with your registered master password",
		Args:  cobra.NoArgs,
		Long: heredoc.Doc(`
			Login stores your provided password in a file so that
			you do not need to enter your	password multiple times.
	
			Remember to logout once you are finished,	otherwise the
			file with your key will	remain and other people may get
			access to it.
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			return login(opts)
		},
	}

	return cmd
}

func login(opts *LoginOptions) error {
	loggedIn := len(opts.Creds.Key) != 0
	if loggedIn {
		return fmt.Errorf("already logged in")
	}

	pw, err := term.Pass("Enter password: ")
	if err != nil {
		return err
	}

	if !opts.Creds.Validate(pw) {
		return fmt.Errorf("wrong password")
	}

	// use previously generated random salt for key generation
	salt, err := dir.Salt()
	if err != nil {
		return err
	}

	key := crypto.Pbkdf2(pw, salt)

	err = dir.WriteKey(key)
	if err == nil {
		fmt.Println("Logged in.")
	}

	return err
}
