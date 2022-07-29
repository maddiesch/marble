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

type TestingTuple2[One any, Two any] []Tuple2[One, Two]

func (t TestingTuple2[One, Two]) Each(fn func(One, Two)) {
	for _, tt := range t {
		fn(tt.One, tt.Two)
	}
}

type Tuple3[One any, Two any, Three any] struct {
	One   One
	Two   Two
	Three Three
}
