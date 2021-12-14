package cmd

import (
	"fmt"
	"os"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "krypt [command]",
	Short: "Krypt is a powerful password manager",
	Long:  "A powerful, featured password manager.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(heredoc.Doc(`
			Usage: krypt [command] [flags]
			Try 'krypt --help' for more information.
		`))
		os.Exit(2)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
