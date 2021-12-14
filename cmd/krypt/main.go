package main

import (
	"fmt"
	"reflect"

	"github.com/MakeNowJust/heredoc"
	"github.com/raklaptudirm/krypt/pkg/cmd"
	"github.com/raklaptudirm/krypt/pkg/crypto"
	"github.com/raklaptudirm/krypt/pkg/dir"
	"github.com/raklaptudirm/krypt/pkg/term"
)

func main() {
	_, err := dir.Checksum()
	if err != nil {
		initKrypt()
		return
	}

	cmd.Execute()
}

func initKrypt() {
	fmt.Println(heredoc.Doc(`
		Initialize krypt with a password, so that your
		information can be secured.
	`))

password:
	pw1, err := term.Pass("Enter new password: ")
	if err != nil {
		fmt.Println(err)
		return
	}

	pw2, err := term.Pass("Re-enter your password: ")
	if err != nil {
		fmt.Println(err)
		return
	}

	for !reflect.DeepEqual(pw1, pw2) {
		pw2, err = term.Pass(heredoc.Doc(`
			Your passwords don't match.
			Re-enter your password, or keep it blank to start over: `))
		if err != nil {
			fmt.Println(err)
			return
		}

		if len(pw2) == 0 {
			goto password
		}
	}

	salt := crypto.RandBytes(8)
	hash := crypto.Sha256(pw1)

	err = dir.WriteSalt(salt)
	if err != nil {
		term.Errorln(err)
		return
	}

	err = dir.WriteChecksum(hash)
	if err != nil {
		term.Errorln(err)
		return
	}
}
