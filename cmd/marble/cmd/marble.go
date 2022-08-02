package cmd

import "github.com/spf13/cobra"

var Marble = &cobra.Command{
	Use:     "marble",
	Short:   "The Marble Language CLI",
	Example: "marble run ./example.marble",
}
