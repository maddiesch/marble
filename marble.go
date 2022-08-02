package marble

import (
	"io"

	"github.com/maddiesch/marble/pkg/core/evaluator"
	"github.com/maddiesch/marble/pkg/core/lexer"
	"github.com/maddiesch/marble/pkg/core/parser"
	"github.com/pkg/errors"
)

type ExecuteOptions struct {
	ParserTracing bool
}

type NamedReader interface {
	io.Reader

	Name() string
}

func Execute(sources []NamedReader, config ...func(*ExecuteOptions)) (any, error) {
	options := ExecuteOptions{}

	for _, config := range config {
		config(&options)
	}

	sourceLexer := lexer.NewLexer()

	for _, reader := range sources {
		lex, err := lexer.New(reader.Name(), reader)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to create lexer for reader: %s", reader.Name())
		}
		sourceLexer.Add(lex)
	}

	parser.SetTracingEnabled(options.ParserTracing)

	parser := parser.New(sourceLexer)
	program := parser.Run()

	if err := parser.Err(); err != nil {
		return nil, errors.Wrap(err, "failed to parse the source")
	}

	rootBinding := evaluator.NewBinding()

	return evaluator.Evaluate(rootBinding, program)
}
