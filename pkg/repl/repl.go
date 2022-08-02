package repl

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/maddiesch/marble/pkg/core/binding"
	"github.com/maddiesch/marble/pkg/core/evaluator"
	"github.com/maddiesch/marble/pkg/core/lexer"
	"github.com/maddiesch/marble/pkg/core/object"
	"github.com/maddiesch/marble/pkg/core/parser"
	"github.com/maddiesch/marble/pkg/version"
)

var (
	Prompt = ">   "
	Indent = "\t"
)

const (
	ExitCommand = ":exit"
	HelpCommand = ":help"
)

var Builtin = map[string]func(*binding.Binding[object.Object], io.Writer) bool{
	ExitCommand: func(*binding.Binding[object.Object], io.Writer) bool {
		return false
	},
	HelpCommand: func(_ *binding.Binding[object.Object], out io.Writer) bool {
		io.WriteString(out, fmt.Sprintf("Marble R.E.P.L. Help (%s)", version.Current))
		io.WriteString(out, "\n")
		return true
	},
	":dump": func(e *binding.Binding[object.Object], out io.Writer) bool {
		fmt.Fprintf(out, e.DebugString())
		return true
	},
}

func Run(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	count := 0

	b := evaluator.NewBinding()

	for {
		count += 1

		io.WriteString(out, fmt.Sprintf("%d%s", count, Prompt))

		var buf bytes.Buffer
		cont := processLine(b, count, scanner, &buf)
		buf.WriteTo(out)

		if !cont {
			return
		}
	}
}

func processLine(b *binding.Binding[object.Object], i int, scanner *bufio.Scanner, buf *bytes.Buffer) bool {
	scanned := scanner.Scan()
	if !scanned {
		return false
	}

	line := scanner.Text()

	if fn, ok := Builtin[strings.TrimSpace(line)]; ok {
		return fn(b, buf)
	}

	l, err := lexer.New("[REPL]", strings.NewReader(line))
	if err != nil {
		buf.WriteString("ERR: ")
		buf.WriteString(err.Error())
		return true
	}

	p := parser.New(l)
	prog := p.Run()

	if err := p.Err(); err != nil {
		if parseErr, ok := err.(*parser.ParseError); ok {
			fmt.Fprintf(os.Stderr, "Parse Errors: %d\n", len(parseErr.Children))
			for _, err := range parseErr.Children {
				fmt.Fprintln(os.Stderr, err.Error())
			}
		} else {
			fmt.Fprintf(os.Stderr, "%+v", err)
		}
		return true
	}

	out, err := evaluator.Evaluate(b, prog)
	if err != nil {
		spew.Dump(err)
		// TODO: Handler Error
		return true
	}

	io.WriteString(buf, fmt.Sprintf("$R%d: ", i))
	io.WriteString(buf, out.DebugString())
	io.WriteString(buf, "\n")

	return true
}
