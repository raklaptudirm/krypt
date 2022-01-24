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

package zsh

import (
	"os"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "zsh",
		Short: "generate the autocompletion script for zsh",
		Args:  cobra.NoArgs,
		Long: heredoc.Doc(`
			Generate the autocompletion script for the zsh shell.

			If shell completion is not already enabled in your environment you will need
			to enable it.  You can execute the following once:
			
			$ echo "autoload -U compinit; compinit" >> ~/.zshrc
			
			To load completions for every new session, execute once:
			# Linux:
			$ krypt completion zsh > "${fpath[1]}/_krypt"
			# macOS:
			$ krypt completion zsh > /usr/local/share/zsh/site-functions/_krypt
			
			You will need to start a new shell for this setup to take effect.
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			return zsh(cmd.Root())
		},
	}

	return cmd
}

func zsh(cmd *cobra.Command) error {
	return cmd.GenZshCompletion(os.Stdout)
}
