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

package pwsh

import (
	"os"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "pwsh",
		Short: "generate the autocompletion script for powershell",
		Args:  cobra.NoArgs,
		Long: heredoc.Doc(`
			Generate the autocompletion script for powershell.

			To load completions in your current shell session:
			PS C:\> krypt completion powershell | Out-String | Invoke-Expression
			
			To load completions for every new session, add the output of the above command
			to your powershell profile.
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			return pwsh(cmd.Root())
		},
	}

	return cmd
}

func pwsh(cmd *cobra.Command) error {
	return cmd.GenPowerShellCompletion(os.Stdout)
}
