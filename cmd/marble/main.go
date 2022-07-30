package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	err := App.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

var App = cli.App{
	Name: "marble",
	Commands: []*cli.Command{
		CommandRepl,
		CommandRun,
	},
}
