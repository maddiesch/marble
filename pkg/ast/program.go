package ast

import (
	"strings"

	"github.com/maddiesch/marble/pkg/token"
)

type Program struct {
	StatementList []Statement
}

func (p *Program) SourceToken() token.Token {
	if len(p.StatementList) > 0 {
		return p.StatementList[0].SourceToken()
	}

	return token.ILLEGAL
}

func (p *Program) String() string {
	var builder strings.Builder

	builder.WriteString("Program Statements:")

	for _, s := range p.StatementList {
		builder.WriteString("\n  " + s.String())
	}

	return builder.String()
}

func (b *Program) MarshalJSON() ([]byte, error) {
	return marshalNode(b)
}

func (*Program) Name() string {
	return "Program"
}

var _ Node = (*Program)(nil)
