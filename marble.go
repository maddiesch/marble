package marble

import (
	"fmt"
	"io"
	"os"

	"github.com/maddiesch/marble/pkg/core/evaluator"
	"github.com/maddiesch/marble/pkg/core/lexer"
	"github.com/maddiesch/marble/pkg/core/parser"
	"github.com/pkg/errors"
)

type ExecuteOptions struct {
	ParserTracing bool
	PrintAST      bool
	Stdout        io.Writer
	Stderr        io.Writer
}

func Execute(programName string, source io.Reader, config ...func(*ExecuteOptions)) (any, error) {
	options := ExecuteOptions{
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}

	for _, config := range config {
		config(&options)
	}

	lex, err := lexer.New(programName, source)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create lexer")
	}

	parser.SetTracingEnabled(options.ParserTracing)

	parse := parser.New(lex)

	program := parse.Run()

	if err := parse.Err(); err != nil {
		return nil, err
	}

	if options.PrintAST {
		fmt.Fprintf(os.Stderr, program.String())
	}

	bind := evaluator.NewBinding()

	result, err := evaluator.Evaluate(bind, program)
	if err != nil {
		return nil, err
	}

	return result.GoValue(), err
}
