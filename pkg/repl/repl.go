package repl

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/maddiesch/marble/pkg/lexer"
	"github.com/maddiesch/marble/pkg/parser"
)

var (
	Prompt      = "> "
	ExitCommand = "_exit"
	Indent      = "\t"
)

var Builtin = map[string]func() bool{}

func Run(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Fprint(out, Prompt)

		var buf bytes.Buffer
		cont := processLine(scanner, &buf)
		buf.WriteTo(out)

		if !cont {
			return
		}
	}
}

func processLine(scanner *bufio.Scanner, buf *bytes.Buffer) bool {
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

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "\t")

	for _, stmt := range prog.StatementList {
		encoder.Encode(stmt)
	}

	return true
}
