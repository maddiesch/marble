package cmd_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/maddiesch/marble/cmd/marble/cmd"
	"github.com/maddiesch/marble/pkg/version"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVersionCommand(t *testing.T) {
	command := cmd.NewRootCommand()

	var stdout strings.Builder
	var stderr strings.Builder

	command.SetOut(&stdout)
	command.SetErr(&stderr)
	command.SetArgs([]string{"version"})

	require.NoError(t, command.Execute())

	assert.Equal(t, "", stderr.String())
	assert.Equal(t, fmt.Sprintf("Marble CLI %s\nMarble Language %s\n", cmd.CurrentVersion, version.Current), stdout.String())
}
