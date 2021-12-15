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

	"github.com/raklaptudirm/krypt/pkg/cmd"
	"github.com/raklaptudirm/krypt/pkg/dir"
	"github.com/raklaptudirm/krypt/pkg/term"
)

func main() {
	_, err := dir.Checksum()

	// if there is no valid checksum, user has not registered a
	// password yet, so request it
	if err != nil {
		err = term.Register()
		if err != nil {
			term.Errorln(err)
			return
		}
	}

	os.Exit(cmd.Execute())
}