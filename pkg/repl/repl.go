package repl

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/maddiesch/marble/pkg/evaluator"
	"github.com/maddiesch/marble/pkg/lexer"
	"github.com/maddiesch/marble/pkg/parser"
)

var (
	Prompt      = ">   "
	ExitCommand = ":exit"
	Indent      = "\t"
)

var Builtin = map[string]func() bool{}

func Run(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	count := 0

	for {
		count += 1

		io.WriteString(out, fmt.Sprintf("%d%s", count, Prompt))

		var buf bytes.Buffer
		cont := processLine(count, scanner, &buf)
		buf.WriteTo(out)

		if !cont {
			return
		}
	}
}

func processLine(i int, scanner *bufio.Scanner, buf *bytes.Buffer) bool {
	scanned := scanner.Scan()
	if !scanned {
		return false
	}

	line := scanner.Text()

	if line == ExitCommand {
		return false
	}

	if fn, ok := Builtin[line]; ok {
		return fn()
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

	out, err := evaluator.Evaluate(prog)
	if err != nil {
		spew.Dump(err)
		// TODO: Handler Error
		return true
	}

	io.WriteString(buf, fmt.Sprintf("$R%d: ", i))
	io.WriteString(buf, out.Description())
	io.WriteString(buf, "\n")

	return true
}
