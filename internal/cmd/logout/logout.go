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

package logout

import (
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/raklaptudirm/krypt/internal/auth"
	"github.com/raklaptudirm/krypt/internal/cmdutil"
	"github.com/spf13/cobra"
)

type LogoutOptions struct {
	Creds *auth.Creds
	Auth  auth.Manager
}

func NewCmd(c *cmdutil.Context) *cobra.Command {
	opts := &LogoutOptions{
		Creds: c.Creds,
		Auth:  c.AuthManager,
	}

	var cmd = &cobra.Command{
		Use:   "logout",
		Short: "log off krypt by removing the encryption key file",
		Args:  cobra.NoArgs,
		Long: heredoc.Doc(`
			Logout clears the file which stores your database key,
			so that accessing the passwords requires logging in with
			the master password.
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			return logout(opts)
		},
	}

	return cmd
}

func logout(opts *LogoutOptions) error {
	loggedIn := len(opts.Creds.Key) != 0
	if loggedIn {
		opts.Auth.SetKey([]byte{})
		fmt.Println("Logged out.")
		return nil
	}

	// not logged in
	return fmt.Errorf("you are not logged in")
}
