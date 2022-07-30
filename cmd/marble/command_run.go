package main

import (
	"os"
	"path/filepath"

	"github.com/maddiesch/marble"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
)

var CommandRun = &cli.Command{
	Name: "run",
	Action: func(c *cli.Context) error {
		path := c.Args().First()
		if path == "" {
			return errors.New("must provide a file path")
		}

		fullPath, err := filepath.Abs(path)
		if err != nil {
			return errors.Wrap(err, "failed to get absolute path to the file")
		}

		f, err := os.Open(fullPath)
		if err != nil {
			return errors.Wrap(err, "failed to open file")
		}
		defer f.Close()

		_, err = marble.Execute(f.Name(), f)

		return err
	},
}
