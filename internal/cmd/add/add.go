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

package add

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
	"laptudirm.com/x/krypt/internal/auth"
	"laptudirm.com/x/krypt/internal/cmdutil"
	"laptudirm.com/x/krypt/pkg/pass"
)

func NewCmd(c *cmdutil.Context) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "add",
		Short: "add a new password to krypt encrypted with your key",
		Args:  cobra.NoArgs,
		Long: heredoc.Doc(`
			Add a new password to krypt, with the provided Name, UserID, and Password.
			The password will be encrypted and stored, and can only be accessed while
			logged in.

			This password can be fetched, edited, or removed using the other krypt
			commands.
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			return add(c.PassManager, c.Creds)
		},
	}

	return cmd
}

func add(passMan pass.Manager, creds *auth.Creds) error {
	if !creds.LoggedIn() {
		return cmdutil.ErrNoLogin
	}

	var password pass.Password
	questions := []*survey.Question{
		{
			Name:   "Name",
			Prompt: &survey.Input{Message: "Name"},
		},
		{
			Name:   "UserID",
			Prompt: &survey.Input{Message: "Username"},
		},
		{
			Name:   "Password",
			Prompt: &survey.Password{Message: "Password"},
		},
	}

	err := survey.Ask(questions, &password)
	if err != nil {
		return err
	}

	err = password.Write(passMan, creds.Key)
	if err != nil {
		return err
	}

	fmt.Println("Added new password.")
	return nil
}
