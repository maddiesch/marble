package cmd

import (
	"fmt"

	"github.com/maddiesch/marble/pkg/build"
	"github.com/spf13/cobra"
)

func init() {
	Marble.AddCommand(Version)
}

var Version = &cobra.Command{
	Use:   "version",
	Short: "Print the current Marble version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Fprintf(cmd.OutOrStdout(), "Marble %s\n", build.Version)
	},
}
