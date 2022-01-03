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
	"github.com/raklaptudirm/krypt/internal/cmd/add"
	"github.com/raklaptudirm/krypt/internal/cmd/edit"
	"github.com/raklaptudirm/krypt/internal/cmd/list"
	"github.com/raklaptudirm/krypt/internal/cmd/login"
	"github.com/raklaptudirm/krypt/internal/cmd/logout"
	"github.com/raklaptudirm/krypt/internal/cmd/rm"
	"github.com/raklaptudirm/krypt/internal/cmd/version"
	"github.com/raklaptudirm/krypt/internal/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmd(c *cmdutil.Context) *cobra.Command {
	var cmd = &cobra.Command{
		Use:  "krypt command",
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

	// child commands
	cmd.AddCommand(rm.NewCmd(c))
	cmd.AddCommand(add.NewCmd(c))
	cmd.AddCommand(edit.NewCmd(c))
	cmd.AddCommand(list.NewCmd(c))
	cmd.AddCommand(login.NewCmd(c))
	cmd.AddCommand(logout.NewCmd(c))
	cmd.AddCommand(version.NewCmd(c))

	return cmd
}
