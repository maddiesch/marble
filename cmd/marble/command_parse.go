package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/davecgh/go-spew/spew"
	"github.com/maddiesch/marble/pkg/lexer"
	"github.com/maddiesch/marble/pkg/parser"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
)

var CommandParse = &cli.Command{
	Name: "parse",
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

		lex, err := lexer.New(f.Name(), f)
		if err != nil {
			return errors.Wrap(err, "failed to create lexer")
		}

		parse := parser.New(lex)

		program := parse.Run()

		if err := parse.Err(); err != nil {
			if parseErr, ok := err.(*parser.ParseError); ok {
				for _, child := range parseErr.Children {
					// TODO: Print child error info
					spew.Dump(child)
				}
				return errors.New("parser error")
			} else {
				return err
			}
		}

		fmt.Fprintln(os.Stdout, program.String())

		return nil
	},
}
