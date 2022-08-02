package cmd

import (
	"context"

	"github.com/spf13/cobra"
)

func NewRootCommand() *cobra.Command {
	root := &cobra.Command{
		Use:     "marble",
		Short:   "The Marble Language CLI",
		Example: "marble run ./example.marble",
	}

	root.AddCommand(Repl)
	root.AddCommand(Run)
	root.AddCommand(Version)

	return root
}

func Execute(ctx context.Context) {
	cobra.CheckErr(NewRootCommand().ExecuteContext(ctx))
}
