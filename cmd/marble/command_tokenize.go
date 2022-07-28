package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/maddiesch/marble/pkg/lexer"
	"github.com/maddiesch/marble/pkg/token"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
)

var CommandTokenize = &cli.Command{
	Name: "tokenize",
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

		encoder := json.NewEncoder(os.Stdout)

		for {
			t := lex.NextToken()

			switch t.Kind {
			case token.EndOfInput:
				return nil
			case token.Invalid:
				return fmt.Errorf("invalid input %s at %s", t, t.Location)
			default:
				encoder.Encode(t)
			}
		}
	},
}
