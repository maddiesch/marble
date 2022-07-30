package ast

type AssignmentStatement interface {
	Statement

	Label() string
	Mutable() bool
	AssignmentExpression() Expression
	CurrentFrameOnly() bool
	RequireUndefined() bool
	RequireDefined() bool
}
