package marble_test

import (
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
		t.Run(name, func(t *testing.T) {
			path := filepath.Join(dir, name)
			file, err := os.Open(path)
			require.NoError(t, err)
			defer file.Close()

			_, err = marble.Execute([]marble.NamedReader{marble.NamedReader(file)})

			assert.NoError(t, err)
		})
	}
}
