package parser

import (
	"fmt"
	"os"
	"strings"
)

func SetTracingEnabled(e bool) {
	tracer.enabled = e
}

var tracer = &_tracer{
	enabled: false,
	indent:  "  ",
}

type _tracer struct {
	enabled bool
	level   int
	indent  string
}

func (a *_tracer) traceIndent() string {
	return strings.Repeat(a.indent, a.level-1)
}

func (a *_tracer) print(s string) {
	if !a.enabled {
		return
	}
	fmt.Fprintf(os.Stderr, "%s%s\n", a.traceIndent(), s)
}

func (a *_tracer) incr() {
	a.level += 1
}

func (a *_tracer) decr() {
	a.level -= 1
}

func trace(m string) string {
	tracer.incr()
	tracer.print("BEGIN " + m)
	return m
}

func untrace(m string) {
	tracer.print("END " + m)
	tracer.decr()
}
