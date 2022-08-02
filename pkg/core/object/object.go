package object

import (
	"github.com/maddiesch/marble/pkg/core/visitor"
	"github.com/maddiesch/marble/pkg/debug"
)

type Object interface {
	debug.Description

	Accept(visitor.Visitor[Object])

	Type() ObjectType

	GoValue() any
}
