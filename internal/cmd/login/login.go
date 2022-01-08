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
	"github.com/raklaptudirm/krypt/internal/cmdutil"
	"github.com/raklaptudirm/krypt/pkg/crypto"
	"github.com/raklaptudirm/krypt/pkg/term"
)

func NewCmd(c *cmdutil.Context) *cobra.Command {
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
			return login(c.AuthManager, c.Creds)
		},
	}

	return cmd
}

func login(authMan auth.Manager, creds *auth.Creds) error {
	loggedIn := len(creds.Key) != 0
	if loggedIn {
		return fmt.Errorf("already logged in")
	}

	pw, err := term.Pass("Enter password: ")
	if err != nil {
		return err
	}

	if !creds.Validate(pw) {
		return fmt.Errorf("wrong password")
	}

	// use previously generated random salt for key generation
	salt, err := authMan.Salt()
	if err != nil {
		return err
	}

	key := crypto.DeriveKey(pw, salt)

	err = authMan.SetKey(key)
	if err == nil {
		fmt.Println("Logged in.")
	}

	return err
}
