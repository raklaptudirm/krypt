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

package term

import (
	"fmt"
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/MakeNowJust/heredoc"
	"github.com/raklaptudirm/krypt/internal/auth"
	"github.com/raklaptudirm/krypt/pkg/crypto"
)

// Error acts as fmt.Print in stderr.
func Error(a ...interface{}) (int, error) {
	return fmt.Fprint(os.Stderr, a...)
}

// Errorf acts as fmt.Printf in stderr.
func Errorf(format string, a ...interface{}) (int, error) {
	return fmt.Fprintf(os.Stderr, format, a...)
}

// Errorln acts as fmt.Println in stderr.
func Errorln(a ...interface{}) (int, error) {
	return fmt.Fprintln(os.Stderr, a...)
}

// Register registers the user and stores the credentials in the provided
// authentication manager.
func Register(man auth.Manager) (err error) {
	fmt.Println(heredoc.Doc(`
		Initialize krypt with a password, so that your information
		can be secured. Make sure the password is strong, as it will
		be protecting all your data.
	`))

password: // main password loop
	var pw1, pw2 string

	password := &survey.Password{Message: "Password"}
	err = survey.AskOne(password, &pw1)
	if err != nil {
		return
	}

confirm: // confirmation loop
	confirm := &survey.Password{Message: "Confirm Password"}
	err = survey.AskOne(confirm, &pw2)
	if err != nil {
		return
	}

	// repeat until the two passwords are equal
	for pw1 != pw2 {
		// if user enters nothing, restart password loop
		if len(pw2) == 0 {
			goto password
		}

		goto confirm
	}

	// passwords are equal

	salt := crypto.RandBytes(8)
	hash := crypto.PassChecksum([]byte(pw1))

	err = man.SetSalt(salt)
	if err != nil {
		return
	}

	err = man.SetChecksum(hash)
	if err != nil {
		return
	}

	fmt.Print("Your password has been registered.\n\n")
	return
}
