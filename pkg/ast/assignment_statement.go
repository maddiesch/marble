package ast

type AssignmentStatement interface {
	Statement

	Label() string
	Mutable() bool
	AssignmentExpression() Expression
}
