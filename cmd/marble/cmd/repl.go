package cmd

import (
	"fmt"

	"github.com/maddiesch/marble/pkg/repl"
	"github.com/spf13/cobra"
)

func init() {
	Marble.AddCommand(Repl)
}

var Repl = &cobra.Command{
	Use:   "repl",
	Short: "runs the marble repl.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Fprintf(cmd.ErrOrStderr(), "Marble R.E.P.L.\nUse '%s' to quit, and '%s' for help\n", repl.ExitCommand, repl.HelpCommand)
		repl.Run(cmd.InOrStdin(), cmd.OutOrStdout())
	},
}
