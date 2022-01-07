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

type ListOptions struct {
	Creds *auth.Creds
	Ident string
	Pass  pass.Manager
}

func NewCmd(c *cmdutil.Context) *cobra.Command {
	opts := &ListOptions{
		Creds: c.Creds,
		Pass:  c.PassManager,
	}

	var cmd = &cobra.Command{
		Use:   "list [name]",
		Short: "un-encrypt and fetch a password from krypt using the filters",
		Args:  cobra.ExactArgs(1),
		Long: heredoc.Doc(`
			List all the passwords which match the provided filters. If no filters
			are provided, all the passwords are listed.
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Ident = args[0]
			return list(opts)
		},
	}

	return cmd
}

func list(opts *ListOptions) error {
	if !opts.Creds.LoggedIn() {
		return cmdutil.ErrNoLogin
	}

	pass, err := pass.GetS(opts.Pass, opts.Ident, opts.Creds.Key)
	if err != nil {
		return err
	}

	fmt.Println(pass.String())
	return nil
}
