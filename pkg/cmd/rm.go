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
	rootCmd.AddCommand(rmCmd)
}

var rmCmd = &cobra.Command{
	Use:   "rm [name]",
	Short: "remove a password from krypt",
	Args:  cobra.ExactArgs(1),
	Long: heredoc.Doc(`
		Logout clears the file which stores your database key,
		so that accessing the passwords requires logging in with
		the master password.
	`),
	Run: rm,
}

func rm(cmd *cobra.Command, args []string) {
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

	var p string
	for hash := range passwords {
		p = hash
		break
	}

	err = pass.Remove(p)
	if err != nil {
		term.Errorln(err)
		return
	}

	fmt.Println("Deleted password.")
}
