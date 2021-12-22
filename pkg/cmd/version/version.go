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

package version

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewCmd(version, buildDate string) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "version",
		Short: "show the krypt software version",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Print(Format(version, buildDate))
		},
	}

	return cmd
}

func Format(version, date string) string {
	if date == "" {
		return version
	}

	return fmt.Sprintf("krypt %v-%v\n", version, date)
}
