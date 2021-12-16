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
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "krypt [command]",
	Short: "Krypt is a powerful password manager",
	Long:  "A powerful, featured password manager.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(heredoc.Doc(`
			Usage: krypt [command] [flags]
			Try 'krypt --help' for more information.
		`))
	},
}

func Execute() int {
	if err := rootCmd.Execute(); err != nil {
		return 1
	}

	return 0
}
