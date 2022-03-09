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

	"github.com/AlecAivazis/survey/v2"
	"github.com/MakeNowJust/heredoc"
	"github.com/raklaptudirm/krypt/internal/auth"
	"github.com/raklaptudirm/krypt/internal/cmdutil"
	"github.com/raklaptudirm/krypt/pkg/pass"
	"github.com/spf13/cobra"
)

func NewCmd(c *cmdutil.Context) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "edit regexp",
		Short: "edit the password which matches the provided regexp",
		Args:  cobra.ExactArgs(1),
		Long: heredoc.Doc(`
			Edit is used to edit a password stored in krypt. All of the details can
			be edited, or kept the same by providing an empty input.

			The password to edit is the first password from the list of password
			names that match the provided regular expression.
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			p, err := pass.Filter(c.PassManager, c.Creds.Key, args[0])
			if err != nil {
				return err
			}

			return edit(c.PassManager, c.Creds, &p[0])
		},
	}

	return cmd
}

func edit(passMan pass.Manager, creds *auth.Creds, p *pass.Password) error {
	if !creds.LoggedIn() {
		return cmdutil.ErrNoLogin
	}

	var password pass.Password
	questions := []*survey.Question{
		{
			Name: "Name",
			Prompt: &survey.Input{
				Message: "Name",
				Default: p.Name,
			},
		},
		{
			Name: "UserID",
			Prompt: &survey.Input{
				Message: "Username",
				Default: p.UserID,
			},
		},
		{
			Name: "Password",
			Prompt: &survey.Password{
				Message: "Password",
			},
		},
	}

	err := survey.Ask(questions, &password)
	if err != nil {
		return err
	}

	err = passMan.Delete(p.Checksum)
	if err != nil {
		return err
	}

	err = password.Write(passMan, creds.Key)
	if err != nil {
		return err
	}

	fmt.Println("Edited password.")
	return nil
}
