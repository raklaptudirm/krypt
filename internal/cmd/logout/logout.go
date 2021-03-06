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
	"github.com/spf13/cobra"
	"laptudirm.com/x/krypt/internal/auth"
	"laptudirm.com/x/krypt/internal/cmdutil"
)

func NewCmd(c *cmdutil.Context) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "logout",
		Short: "logout logs the user out of krypt and blocks access",
		Args:  cobra.NoArgs,
		Long: heredoc.Doc(`
			Logout logs you out from krypt, preventing access to the database without
			the master password. Use this once you are done using krypt.

			After logging out, you need to log back in to access the database.
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			return logout(c.AuthManager, c.Creds)
		},
	}

	return cmd
}

func logout(authMan auth.Manager, creds *auth.Creds) error {
	loggedIn := len(creds.Key) != 0
	if loggedIn {
		authMan.SetKey([]byte{})
		fmt.Println("Logged out.")
		return nil
	}

	// not logged in
	return fmt.Errorf("you are not logged in")
}
