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
		Use:   "list [name]",
		Short: "un-encrypt and fetch a password from krypt using the filters",
		Args:  cobra.ExactArgs(1),
		Long: heredoc.Doc(`
			List all the passwords which match the provided filters. If no filters
			are provided, all the passwords are listed.
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			return list(c.PassManager, c.Creds, args[0])
		},
	}

	return cmd
}

func list(passMan pass.Manager, creds *auth.Creds, ident string) error {
	if !creds.LoggedIn() {
		return cmdutil.ErrNoLogin
	}

	pass, err := pass.GetS(passMan, ident, creds.Key)
	if err != nil {
		return err
	}

	fmt.Println(pass.String())
	return nil
}
