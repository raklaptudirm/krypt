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

package edit

import (
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/raklaptudirm/krypt/internal/auth"
	"github.com/raklaptudirm/krypt/internal/cmdutil"
	"github.com/raklaptudirm/krypt/pkg/pass"
	"github.com/raklaptudirm/krypt/pkg/term"
	"github.com/spf13/cobra"
)

type EditOptions struct {
	Creds    *auth.Creds
	Password *pass.Password
	PassHash string
}

func NewCmd(c *cmdutil.Context) *cobra.Command {
	opts := &EditOptions{
		Creds: c.Creds,
	}

	var cmd = &cobra.Command{
		Use:   "edit [name]",
		Short: "edit a stored password in krypt",
		Args:  cobra.ExactArgs(1),
		Long: heredoc.Doc(`
			Edit is used to edit a password in krypt, delete the
			previous one, and store the new one.
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			name, pass, err := pass.GetS(args[0], c.Creds.Key)
			if err != nil {
				return err
			}

			opts.Password = pass
			opts.PassHash = name
			return edit(opts)
		},
	}

	return cmd
}

func edit(opts *EditOptions) error {
	if !opts.Creds.LoggedIn() {
		return cmdutil.ErrNoLogin
	}

	name, err := term.Input("name: ")
	if err != nil {
		return err
	}
	if name == "" {
		name = opts.Password.Name
	}

	user, err := term.Input("username: ")
	if err != nil {
		return err
	}
	if user == "" {
		user = opts.Password.UserID
	}

	p, err := term.Pass("password: ")
	if err != nil {
		return err
	}
	if len(p) == 0 {
		p = []byte(opts.Password.Password)
	}

	password := pass.Password{
		Name:     name,
		UserID:   user,
		Password: string(p),
	}

	err = pass.Remove(opts.PassHash)
	if err != nil {
		return err
	}

	err = password.Write(opts.Creds.Key)
	if err != nil {
		return err
	}

	fmt.Println("Edited password.")
	return nil
}
