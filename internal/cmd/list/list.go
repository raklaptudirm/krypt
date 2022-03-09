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

package list

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
		Use:   "list [regexp]",
		Short: "list the passwords which match the provided regexp",
		Args:  cobra.RangeArgs(0, 1),
		Long: heredoc.Doc(`
			List all the passwords which match the provided regular expression. If no
			regular expression is provided, all the passwords are listed.

			The printed passwords are censored by default.
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			regex := "" // empty string matches all
			if len(args) > 0 {
				regex = args[0]
			}
			return list(c.PassManager, c.Creds, regex)
		},
	}

	return cmd
}

func list(passMan pass.Manager, creds *auth.Creds, ident string) error {
	if !creds.LoggedIn() {
		return cmdutil.ErrNoLogin
	}

	passwords, err := pass.Filter(passMan, creds.Key, ident)
	if err != nil {
		return err
	}

	length := len(passwords) - 1
	for i, pass := range passwords {
		fmt.Println(pass.String())

		// print a newline for all but last
		if i != length {
			fmt.Println()
		}
	}

	return nil
}
