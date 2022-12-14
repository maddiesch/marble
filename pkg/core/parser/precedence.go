package parser

type Precedence uint8

const (
	_ Precedence = iota
	Lowest
	Equals      // ==
	LessGreater // < or >
	Sum         // +
	Product     // *
	Prefix      // -X or !X
	Call        // function(X)
	Assignment  // =
	Subscript   // foo[1]
)
