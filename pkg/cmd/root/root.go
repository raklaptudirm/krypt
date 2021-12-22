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
	"github.com/raklaptudirm/krypt/pkg/cmd/add"
	"github.com/raklaptudirm/krypt/pkg/cmd/edit"
	"github.com/raklaptudirm/krypt/pkg/cmd/list"
	"github.com/raklaptudirm/krypt/pkg/cmd/login"
	"github.com/raklaptudirm/krypt/pkg/cmd/logout"
	"github.com/raklaptudirm/krypt/pkg/cmd/rm"
	"github.com/raklaptudirm/krypt/pkg/cmd/version"
	"github.com/raklaptudirm/krypt/pkg/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmd(f *cmdutil.Factory, versionNum, buildDate string) *cobra.Command {
	var cmd = &cobra.Command{
		Use:  "krypt command",
		Args: cobra.NoArgs,

		SilenceErrors: true,
		SilenceUsage:  true,
	}

	// global flags
	cmd.PersistentFlags().BoolP("help", "h", false, "show help for command")
	cmd.PersistentFlags().BoolP("version", "v", false, "show software version")

	versionStr := version.Format(versionNum, buildDate)
	cmd.SetVersionTemplate(versionStr)
	cmd.Version = versionStr

	// child commands
	cmd.AddCommand(rm.NewCmd(f))
	cmd.AddCommand(add.NewCmd(f))
	cmd.AddCommand(edit.NewCmd(f))
	cmd.AddCommand(list.NewCmd(f))
	cmd.AddCommand(login.NewCmd(f))
	cmd.AddCommand(logout.NewCmd(f))
	cmd.AddCommand(version.NewCmd(versionNum, buildDate))

	return cmd
}
