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

package root

import (
	"github.com/spf13/cobra"
	"laptudirm.com/x/krypt/internal/cmd/add"
	"laptudirm.com/x/krypt/internal/cmd/completion"
	"laptudirm.com/x/krypt/internal/cmd/edit"
	"laptudirm.com/x/krypt/internal/cmd/help"
	"laptudirm.com/x/krypt/internal/cmd/list"
	"laptudirm.com/x/krypt/internal/cmd/login"
	"laptudirm.com/x/krypt/internal/cmd/logout"
	"laptudirm.com/x/krypt/internal/cmd/master"
	"laptudirm.com/x/krypt/internal/cmd/rm"
	"laptudirm.com/x/krypt/internal/cmd/version"
	"laptudirm.com/x/krypt/internal/cmdutil"
)

func NewCmd(c *cmdutil.Context) *cobra.Command {
	var cmd = &cobra.Command{
		Use:  "krypt",
		Args: cobra.NoArgs,

		SilenceErrors: true,
		SilenceUsage:  true,
	}

	// global flags
	cmd.PersistentFlags().BoolP("help", "h", false, "show help for command")
	cmd.PersistentFlags().BoolP("version", "v", false, "show software version")

	versionStr := c.Version.String()
	cmd.SetVersionTemplate(versionStr)
	cmd.Version = versionStr

	// set custom help
	hc := help.NewCmd()
	cmd.SetHelpCommand(hc)
	cmd.SetHelpFunc(hc.Run)

	// child commands
	cmd.AddCommand(rm.NewCmd(c))
	cmd.AddCommand(add.NewCmd(c))
	cmd.AddCommand(edit.NewCmd(c))
	cmd.AddCommand(list.NewCmd(c))
	cmd.AddCommand(login.NewCmd(c))
	cmd.AddCommand(logout.NewCmd(c))
	cmd.AddCommand(master.NewCmd(c))
	cmd.AddCommand(version.NewCmd(c))
	cmd.AddCommand(completion.NewCmd())

	return cmd
}
