package marble_test

import (
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/maddiesch/marble"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExampleScript(t *testing.T) {
	fileNames := []string{"example.marble"}
	dir, _ := os.Getwd()

	for _, name := range fileNames {
		path := filepath.Join(dir, name)
		file, err := os.Open(path)
		require.NoError(t, err)
		defer file.Close()

		t.Run(name, func(t *testing.T) {
			_, err := marble.Execute(name, file, func(o *marble.ExecuteOptions) {
				o.Stdout = io.Discard
				o.Stderr = io.Discard
			})

			assert.NoError(t, err)
		})
	}
}
