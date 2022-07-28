package repl

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/maddiesch/marble/pkg/lexer"
	"github.com/maddiesch/marble/pkg/token"
)

var (
	Prompt      = "> "
	ExitCommand = "_exit"
)

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

	l, err := lexer.New("[REPL]", strings.NewReader(line))
	if err != nil {
		buf.WriteString("ERR: ")
		buf.WriteString(err.Error())
		return true
	}

	for {
		t := l.NextToken()
		if t.Kind == token.EndOfInput {
			return true
		}

		json.NewEncoder(buf).Encode(t)
	}
}
