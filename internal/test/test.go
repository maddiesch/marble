package test

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type TestingT interface {
	assert.TestingT
	require.TestingT

	Name() string
}
