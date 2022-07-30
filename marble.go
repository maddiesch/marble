package marble

import (
	"fmt"
	"io"
	"os"

	"github.com/maddiesch/marble/pkg/env"
	"github.com/maddiesch/marble/pkg/evaluator"
	"github.com/maddiesch/marble/pkg/lexer"
	"github.com/maddiesch/marble/pkg/parser"
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

	environment := env.New(options.Stdout, options.Stdout)

	result, err := evaluator.Evaluate(environment, program)
	if err != nil {
		return nil, err
	}

	return result.GoValue(), err
}
