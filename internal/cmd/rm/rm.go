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

package rm

import (
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/raklaptudirm/krypt/internal/auth"
	"github.com/raklaptudirm/krypt/internal/cmdutil"
	"github.com/raklaptudirm/krypt/pkg/pass"
	"github.com/spf13/cobra"
)

func NewCmd(c *cmdutil.Context) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "rm regexp",
		Short: "remove a password which matches the provided regexp",
		Args:  cobra.ExactArgs(1),
		Long: heredoc.Doc(`
			Remove deletes the password matching the regular expression from krypt.

			The password to delete is the first password from the list of password
			names that match the provided regular expression.
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			ps, err := pass.Filter(c.PassManager, c.Creds.Key, args[0])
			if err != nil {
				return err
			}

			return rm(c.PassManager, c.Creds, ps[0].Checksum)
		},
	}

	return cmd
}

func rm(passMan pass.Manager, creds *auth.Creds, checksum []byte) error {
	if !creds.LoggedIn() {
		return cmdutil.ErrNoLogin
	}

	err := passMan.Delete(checksum)
	if err != nil {
		return err
	}

	fmt.Println("Deleted password.")
	return nil
}
