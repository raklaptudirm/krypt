package cmd

import (
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/raklaptudirm/krypt/pkg/dir"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(logoutCmd)
}

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Log off your krypt database.",
	Args:  cobra.NoArgs,
	Long: heredoc.Doc(`
		Logout clears the file which stores your database
		key, so that accessing the passwords requires logging
		in with the master password.
	`),
	Run: func(cmd *cobra.Command, args []string) {
		_, err := dir.Key()
		if err == nil {
			dir.WriteKey([]byte{})
			fmt.Println("Logged out of krypt.")
			return
		}

		fmt.Println("Already logged out of krypt.")
	},
}
