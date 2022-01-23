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

package bash

import (
	"os"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "bash",
		Short: "generate the autocompletion script for bash",
		Args:  cobra.NoArgs,
		Long: heredoc.Doc(`
			Generate the autocompletion script for the bash shell.

			This script depends on the 'bash-completion' package.
			If it is not installed already, you can install it via your OS's package manager.
			
			To load completions in your current shell session:
			$ source <(krypt completion bash)
			
			To load completions for every new session, execute once:
			Linux:
				$ krypt completion bash > /etc/bash_completion.d/krypt
			MacOS:
				$ krypt completion bash > /usr/local/etc/bash_completion.d/krypt
			
			You will need to start a new shell for this setup to take effect.
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			return bash(cmd.Root())
		},
	}

	return cmd
}

func bash(cmd *cobra.Command) error {
	return cmd.GenBashCompletion(os.Stdout)
}
