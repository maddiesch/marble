package token

import (
	"fmt"
	"strings"
)

type Kind string

const (
	// Keywords
	Let      Kind = "LET"
	Const    Kind = "CONST"
	Function Kind = "FUNC"
	Return   Kind = "RETURN"
	If       Kind = "IF"
	Else     Kind = "ELSE"
	Not      Kind = "NOT"
	Mutate   Kind = "MUTATE"
	Do       Kind = "DO"

	// Identifiers & Literals
	Identifier   Kind = "IDENT"
	Integer      Kind = "INT"
	Float        Kind = "FLOAT"
	Comment      Kind = "COMMENT"
	String       Kind = "STRING"
	LiteralTrue  Kind = "TRUE"
	LiteralFalse Kind = "FALSE"
	LiteralNull  Kind = "NULL"

	// Operators
	Assign      Kind = "ASSIGN"
	Plus        Kind = "PLUS"
	Minus       Kind = "MINUS"
	Asterisk    Kind = "ASTER"
	Slash       Kind = "SLASH"
	LessThan    Kind = "LESSTHAN"
	GreaterThan Kind = "GREATERTHAN"
	Equate      Kind = "EQUATE"
	NotEquate   Kind = "NOTEQUATE"
	Bang        Kind = "BANG"
	DoubleColon Kind = "DOUBLECOLON"

	// Delimiters
	Comma     Kind = "COMMA"
	Question  Kind = "QUESTION"
	Semicolon Kind = "SEMICOLON"
	Colon     Kind = "COLON"
	LParen    Kind = "LPAREN"
	RParen    Kind = "RPAREN"
	LBrace    Kind = "LBRACE"
	RBrace    Kind = "RBRACE"
	LBracket  Kind = "LBRACKET"
	RBracket  Kind = "RBRACKET"
	Tilde     Kind = "TILDE"
	Carrot    Kind = "CARROT"
	Percent   Kind = "PERCENT"
	Period    Kind = "PERIOD"

	EndOfInput Kind = "EOF"
	Invalid    Kind = "INVALID"
)

// Defines a shared illegal token
var ILLEGAL = Token{Kind: Invalid}

type Token struct {
	Kind     Kind
	Value    string
	Location Location
}

func (t Token) Is(k Kind) bool {
	return t.Kind == k
}

func (t Token) IsNot(k Kind) bool {
	return !t.Is(k)
}

func (t Token) String() string {
	var b strings.Builder
	b.WriteString(string(t.Kind))
	if t.Value != "" {
		b.WriteRune('(')
		b.WriteString(t.Value)
		b.WriteRune(')')
	}
	return b.String()
}

type Location struct {
	Filename string
	Offset   int
	Line     int
	Column   int
}

func (l Location) String() string {
	return fmt.Sprintf("%s:%d+%d", l.Filename, l.Line, l.Column)
}
