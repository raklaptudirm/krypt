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
	rootCmd.AddCommand(getCmd)
}

var getCmd = &cobra.Command{
	Use:   "get [name]",
	Short: "Get a stored krypt password.",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("no of args need to be 1")
		}
		return nil
	},
	Long: heredoc.Doc(`
		Add a new password with a name, username and password to
		krypt. You will later be able to edit and manipulate this
		password.
	`),
	Run: func(cmd *cobra.Command, args []string) {
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
		} else {
			fmt.Println(passwords[0].String())
		}
	},
}
