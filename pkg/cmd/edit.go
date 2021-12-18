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

package cmd

import (
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/raklaptudirm/krypt/pkg/pass"
	"github.com/raklaptudirm/krypt/pkg/term"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(editCmd)
}

var editCmd = &cobra.Command{
	Use:   "edit [name]",
	Short: "edit a stored password in krypt",
	Args:  cobra.ExactArgs(1),
	Long: heredoc.Doc(`
		Edit is used to edit a password in krypt, delete the
		previous one, and store the new one.
	`),
	Run: edit,
}

func edit(cmd *cobra.Command, args []string) {
	passwords, err := pass.Get(pass.Filter{
		Type: pass.FilterName,
		Data: args[0],
	})
	if err != nil {
		term.Errorln(err)
		return
	}

	if len(passwords) == 0 {
		term.Errorln("No passwords matched provided filters.")
		return
	}

	var currentPass pass.Password
	var currentPassName string
	for pName, password := range passwords {
		currentPass = password
		currentPassName = pName
		break
	}

	name, err := term.Input("name: ")
	if err != nil {
		term.Errorln(err)
		return
	}
	if name == "" {
		name = currentPass.Name
	}

	user, err := term.Input("username: ")
	if err != nil {
		term.Errorln(err)
		return
	}
	if user == "" {
		user = currentPass.UserID
	}

	p, err := term.Pass("password: ")
	if err != nil {
		term.Errorln(err)
		return
	}
	if len(p) == 0 {
		p = []byte(currentPass.Password)
	}

	password := pass.Password{
		Name:     name,
		UserID:   user,
		Password: string(p),
	}

	err = pass.Remove(currentPassName)
	if err != nil {
		term.Errorln(err)
		return
	}

	err = password.Write()
	if err != nil {
		term.Errorln(err)
		return
	}

	fmt.Println("Edited password.")
}
