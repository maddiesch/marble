package main

import (
	"fmt"
	"os"

	"github.com/maddiesch/marble/pkg/repl"
	"github.com/urfave/cli/v2"
)

var CommandRepl = &cli.Command{
	Name: "repl",
	Action: func(c *cli.Context) error {
		fmt.Fprintf(os.Stderr, "Marble R.E.P.L.\nUse '%s' to quit, and '%s' for help\n", repl.ExitCommand, repl.HelpCommand)
		repl.Run(os.Stdin, os.Stdout)

		return nil
	},
}
