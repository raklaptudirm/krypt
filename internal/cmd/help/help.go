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

package help

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/raklaptudirm/krypt/pkg/term"
	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "help",
		Short: "help prints information about the provided command",
		Long: heredoc.Doc(`
			Help prints, if present, the long help line, otherwise the short help line
			along with the usage information for the command.

			It can also be used as the --help or -h flag.
		`),
		Run: func(cmd *cobra.Command, args []string) {
			root := cmd.Root()
			target, extra, _ := root.Traverse(args)

			help(target, extra)
		},
	}

	return cmd
}

func help(cmd *cobra.Command, extra []string) {
	if len(extra) > 0 {
		term.Errorf("unknown command \"%v\" for \"%v\"\n", extra[0], getFullName(cmd))
		return
	}

	if cmd.Long != "" {
		term.Errorln(cmd.Long)
	} else {
		term.Errorln(cmd.Short)
	}

	cmd.Usage()
	term.Errorln()
	return
}

func getFullName(cmd *cobra.Command) string {
	name := cmd.Name()

	if cmd.HasParent() {
		name = getFullName(cmd.Parent()) + " " + name
	}

	return name
}
