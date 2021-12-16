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
	rootCmd.AddCommand(addCmd)
}

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a password to krypt.",
	Args:  cobra.NoArgs,
	Long: heredoc.Doc(`
		Add a new password with a name, username and password to
		krypt. You will later be able to edit and manipulate this
		password.
	`),
	Run: func(cmd *cobra.Command, args []string) {
		name, err := term.Input("name: ")
		if err != nil {
			term.Errorln(err)
			return
		}

		user, err := term.Input("username: ")
		if err != nil {
			term.Errorln(err)
			return
		}

		p, err := term.Pass("password: ")
		if err != nil {
			term.Errorln(err)
			return
		}

		password := pass.Password{
			Name:     name,
			UserID:   user,
			Password: string(p),
		}

		err = password.Write()
		if err != nil {
			term.Errorln(err)
			return
		}

		fmt.Println("Added new password.")
	},
}
