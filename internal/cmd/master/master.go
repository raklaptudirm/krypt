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

package master

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
	"laptudirm.com/x/krypt/internal/auth"
	"laptudirm.com/x/krypt/internal/cmdutil"
	"laptudirm.com/x/krypt/pkg/crypto"
	"laptudirm.com/x/krypt/pkg/pass"
	"laptudirm.com/x/krypt/pkg/term"
)

func NewCmd(c *cmdutil.Context) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "master",
		Short: "master changes krypt's master password",
		Args:  cobra.NoArgs,
		Long: heredoc.Doc(`
			Reset your master/primary password to the provided password. After the
			reset, the password will be stored and the existing passwords will be re-
			encrypted automatically.
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			return master(c.Creds, c.AuthManager, c.PassManager)
		},
	}

	return cmd
}

func master(c *auth.Creds, authMan auth.Manager, passMan pass.Manager) error {
	passwords, err := pass.Get(passMan, c.Key)
	if err != nil {
		return err
	}

	dataset, err := passMan.Passwords()
	if err != nil {
		return err
	}

	var hashes [][]byte
	for _, data := range dataset {
		hashes = append(hashes, crypto.Checksum(data))
	}

	err = passMan.Delete(hashes...)
	if err != nil {
		return err
	}

	term.Register(authMan)

	nKey, err := authMan.Key()
	if err != nil {
		return err
	}

	for _, password := range passwords {
		err = password.Write(passMan, nKey)
		if err != nil {
			return err
		}
	}
	return nil
}
