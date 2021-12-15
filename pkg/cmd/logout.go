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
	"github.com/raklaptudirm/krypt/pkg/dir"
	"github.com/raklaptudirm/krypt/pkg/term"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(logoutCmd)
}

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Log off your krypt database.",
	Args:  cobra.NoArgs,
	Long: heredoc.Doc(`
		Logout clears the file which stores your database key,
		so that accessing the passwords requires logging in with
		the master password.
	`),
	Run: func(cmd *cobra.Command, args []string) {
		_, err := dir.Key()

		// if there is a valid key in the keyfile, dir.Key will return nil
		// error, so the user must be logged in
		if err == nil {
			dir.WriteKey([]byte{})
			fmt.Println("Logged out.")
			return
		}

		// not logged in
		term.Errorln("you are not logged in.")
	},
}
