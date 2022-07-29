package token_test

import (
	"testing"

	"github.com/maddiesch/marble/pkg/token"
	"github.com/stretchr/testify/assert"
)

func TestLocationStringer(t *testing.T) {
	l := token.Location{
		Filename: "testing_file",
		Line:     42,
		Offset:   129,
		Column:   16,
	}

	assert.Equal(t, "testing_file:42+16", l.String())
}
