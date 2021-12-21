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

package main

import (
	"os"

	"github.com/raklaptudirm/krypt/internal/auth"
	"github.com/raklaptudirm/krypt/pkg/cmd/root"
	"github.com/raklaptudirm/krypt/pkg/cmdutil"
	"github.com/raklaptudirm/krypt/pkg/term"
	"github.com/spf13/cobra"
)

type exitCode int

const (
	exitOkay  exitCode = 0
	exitError exitCode = 1
)

func main() {
	exit := kryptMain()
	os.Exit(int(exit))
}

func kryptMain() exitCode {
	factory := cmdutil.NewFactory()

	credentials, err := auth.Get()
	if err == nil {
		factory.Auth = credentials
	}

	rootCmd := root.NewCmd(factory)
	rootCmd.SetArgs(os.Args[1:])

	if cmd, err := rootCmd.ExecuteC(); err != nil {
		printError(err, cmd)
		return exitError
	}
	return exitOkay
}

func printError(err error, cmd *cobra.Command) {
	term.Errorln(err)
	term.Errorln() // a line gap
	term.Errorln(cmd.UsageString())
}
