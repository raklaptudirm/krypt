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
	"errors"
	"io"
	"os"

	"github.com/raklaptudirm/krypt/internal/auth"
	"github.com/raklaptudirm/krypt/internal/build"
	"github.com/raklaptudirm/krypt/internal/cmd/root"
	"github.com/raklaptudirm/krypt/internal/cmdutil"
	"github.com/raklaptudirm/krypt/pkg/term"
	"github.com/spf13/cobra"
)

type exitCode int

const (
	exitOkay   exitCode = 0
	exitError  exitCode = 1
	exitIntrpt exitCode = 2
)

func main() {
	exit := m(os.Args)
	os.Exit(int(exit))
}

func m(args []string) exitCode {
	context := cmdutil.NewContext()
	context.Creds = auth.Get()
	context.Version = cmdutil.NewVersion(build.Version, build.Date)

	// register user if not already
	if !context.Creds.Registered() {
		err := term.Register()
		if err != nil {
			term.Errorln(err)
			return exitError
		}
	}

	rootCmd := root.NewCmd(context)
	rootCmd.SetArgs(args[1:])

	return handleError(rootCmd.ExecuteC())
}

func handleError(cmd *cobra.Command, err error) exitCode {
	if err == nil {
		return exitOkay
	}

	switch {
	case errors.Is(err, io.EOF):
		term.Errorln("[interrupted]")
		return exitIntrpt
	}

	printError(cmd, err)
	return exitError
}

func printError(cmd *cobra.Command, err error) {
	term.Errorln(err)
	term.Errorln() // a line gap
	term.Errorln(cmd.UsageString())
}
