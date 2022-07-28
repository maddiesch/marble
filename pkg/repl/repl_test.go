package repl_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/maddiesch/marble/pkg/repl"
	"github.com/stretchr/testify/assert"
)

func TestRun(t *testing.T) {
	in := strings.NewReader("1 + 1\n")
	var out bytes.Buffer

	repl.Run(in, &out)

	assert.Equal(t, "> token: INT(1)\ntoken: PLUS\ntoken: INT(1)\n> ", out.String())
}
