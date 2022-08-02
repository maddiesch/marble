package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/maddiesch/marble"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var Run = &cobra.Command{
	Use:     "run",
	Short:   "Execute the script",
	Example: "marble run ./example.marble",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("must provide at least one argument")
		}

		for i, arg := range args {
			path, err := filepath.Abs(arg)
			if err != nil {
				return errors.Wrap(err, "failed to get absolute path for argument "+fmt.Sprint(i))
			}
			_, err = os.Stat(path)
			if os.IsNotExist(err) {
				return errors.New("unable to find file at path " + path)
			}
			args[i] = path
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		namedReaderSlice := make([]marble.NamedReader, len(args))

		for i, sourcePath := range args {
			file, err := os.Open(sourcePath)
			if err != nil {
				return errors.Wrapf(err, "failed to open source file at path %s", sourcePath)
			}
			defer file.Close() // Called when function goes out of scope, not the loop

			namedReaderSlice[i] = marble.NamedReader(file)
		}

		_, err := marble.Execute(namedReaderSlice, func(o *marble.ExecuteOptions) {

		})

		return err
	},
}
