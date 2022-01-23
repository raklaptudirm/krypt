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

package fish

import (
	"os"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	var noDesc *bool
	var cmd = &cobra.Command{
		Use:   "fish",
		Short: "generate the autocompletion script for fish",
		Args:  cobra.NoArgs,
		Long: heredoc.Doc(`
			Generate the autocompletion script for the fish shell.

			To load completions in your current shell session:
			$ krypt completion fish | source
			
			To load completions for every new session, execute once:
			$ krypt completion fish > ~/.config/fish/completions/krypt.fish  
			
			You will need to start a new shell for this setup to take effect.
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			return fish(cmd.Root(), *noDesc)
		},
	}

	noDesc = cmd.Flags().BoolP("no-descriptions", "n", false, "disable completion descriptions")

	return cmd
}

func fish(cmd *cobra.Command, noDesc bool) error {
	return cmd.GenFishCompletion(os.Stdout, !noDesc)
}
