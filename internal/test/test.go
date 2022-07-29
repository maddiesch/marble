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

type Tuple2[One any, Two any] struct {
	One One
	Two Two
}

type Tuple3[One any, Two any, Three any] struct {
	One   One
	Two   Two
	Three Three
}
