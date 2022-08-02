package cmd

import (
	"fmt"

	"github.com/maddiesch/marble/pkg/version"
	"github.com/spf13/cobra"
)

const (
	CurrentVersion = "1.0.0.pre-alpha"
)

var Version = &cobra.Command{
	Use:   "version",
	Short: "Print the current Marble version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Fprintf(cmd.OutOrStdout(), "Marble CLI %s\nMarble Language %s\n", CurrentVersion, version.Current)
	},
}
