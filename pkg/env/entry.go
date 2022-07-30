package env

import (
	"github.com/maddiesch/marble/pkg/object"
	"github.com/maddiesch/marble/pkg/token"
)

type Entry struct {
	Value object.Object

	location  token.Location
	frame     uint64
	pointer   uint64
	mutable   bool
	protected bool
}
