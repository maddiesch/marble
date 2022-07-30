package ast

import (
	"strconv"

	"github.com/maddiesch/marble/pkg/token"
)

type Entrypoint struct{}

func (Entrypoint) SourceToken() token.Token {
	return token.Token{}
}

func (e Entrypoint) String() string {
	return e.Name()
}

func (Entrypoint) Name() string {
	return "Entrypoint"
}

func (Entrypoint) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Quote("ENTRYPOINT")), nil
}

var _ Node = (*Entrypoint)(nil)
