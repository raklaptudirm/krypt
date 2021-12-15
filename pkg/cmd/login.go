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
	"reflect"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/raklaptudirm/krypt/pkg/crypto"
	"github.com/raklaptudirm/krypt/pkg/dir"
	"github.com/raklaptudirm/krypt/pkg/term"
)

func init() {
	rootCmd.AddCommand(loginCmd)
}

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to your krypt database.",
	Args:  cobra.NoArgs,
	Long: heredoc.Doc(`
		Login stores your provided password in a file so that
		you do not need to enter your	password multiple times.

		Remember to logout once you are finished,	otherwise the
		file with your key will	remain and other people may get
		access to it.
	`),
	Run: func(cmd *cobra.Command, args []string) {
		pw, err := term.Pass("Enter password: ")

		if err != nil {
			term.Errorln(err)
			return
		}

		checksum, err := dir.Checksum()
		if err != nil {
			term.Errorln(err)
			return
		}

		hash := crypto.Sha256(pw)
		if !reflect.DeepEqual(hash, checksum) {
			term.Errorln("Wrong password.")
			return
		}

		salt, err := dir.Salt()
		if err != nil {
			term.Errorln(err)
			return
		}

		key := crypto.Pbkdf2(pw, salt)

		err = dir.WriteKey(key)
		if err == nil {
			fmt.Println("Logged in to krypt.")
		} else {
			term.Errorln(err)
		}
	},
}
