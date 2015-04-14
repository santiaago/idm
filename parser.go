package main

import (
	"fmt"
	"io"
	"strconv"
)

// Parser represents a parser.
type Parser struct {
	s   *Scanner
	buf struct {
		t   Token  // last read token
		lit string // last read literal
		n   int    // buffer size (max=1)
	}
}

// NewParser returns a new instance of Parser.
func NewParser(r io.Reader) *Parser {
	return &Parser{s: NewScanner(r)}
}

// scan returns the next token from the underlying scanner.
// if a token has been unscanned then read that instead.
func (p *Parser) scan() (t Token, lit string) {

	if p.buf.n != 0 {
		p.buf.n = 0
		return p.buf.t, p.buf.lit
	}

	t, lit = p.s.Scan()

	p.buf.t, p.buf.lit = t, lit
	return
}

// unscan pushes the previously read token back onto the buffer.
func (p *Parser) unscan() { p.buf.n = 1 }

// scanIgnoreWhitespace scans the next non-whitespace token.
func (p *Parser) scanIgnoreWhitespace() (t Token, lit string) {
	t, lit = p.scan()
	if t == Space {
		t, lit = p.scan()
	}
	return
}

// Expression is the interface for an expression.
type Expression interface {
	String() string
	Evaluate() string
}

// Statement represents a code statement a = 2.
type Statement struct {
	Left  string
	Right string
}

// String retuns the statement as a string.
func (s Statement) String() string {
	return fmt.Sprintf("%v%v", s.Left, s.Right)
}

// Evaluate evaluates the given statement.
func (s Statement) Evaluate() string {
	return fmt.Sprintf("%v%v", s.Left, s.Right)
}

// Num represents a number statement.
type Num struct {
	n string
}

// String returns the string of the number
func (n Num) String() string {
	return fmt.Sprintf("%v", n.n)
}

// Evaluate returns the number value
func (n Num) Evaluate() string {
	return n.String()
}

// Binary represents a binary statement
// example: 12 + 3
type Binary struct {
	Left     string
	Right    string
	Operator string
}

// String returns the string of the number
func (b Binary) String() string {
	return fmt.Sprintf("%v %v %v", b.Left, b.Operator, b.Right)
}

// Evaluate returns the number value
func (b Binary) Evaluate() string {
	if b.Operator == "+" {
		return add(b.Left, b.Right)
	}
	return ""
}

func add(a, b string) string {
	var ia, ib int64
	var err error
	ia, err = strconv.ParseInt(a, 10, 64)
	if err != nil {
		fmt.Println(err)
	}
	ib, err = strconv.ParseInt(b, 10, 64)
	return fmt.Sprintf("%d", ia+ib)
}

// Parse parse a assign statement a = b
func (p *Parser) Parse() (*Expression, error) {
	var left, right, operator string
	tok, lit := p.scanIgnoreWhitespace()
	var lastTok Token
	if tok == Identifier {
		left = lit
	} else if tok == Number {
		left = lit
		lastTok = Number
	} else {
		return nil, fmt.Errorf("found %q, expected left", lit)
	}

	tok, lit = p.scanIgnoreWhitespace()
	if tok == EOF && lastTok == Number {
		e := Expression(Num{left})
		return &e, nil
	}

	isAssign := tok == Assign
	if tok != Assign && tok != Operator {
		return nil, fmt.Errorf("found %q, expected '=' with tok: %v, expected %v", lit, tok, Assign)
	}
	// get token
	if tok == Operator {
		operator = lit
	}

	tok, lit = p.scanIgnoreWhitespace()
	if tok == Identifier {
		right = lit
	} else if tok == Number && !isAssign {
		right = lit
	} else {
		return nil, fmt.Errorf("found %q, expected identifier name", lit)
	}
	var expr Expression
	if isAssign {
		expr = Expression(Statement{Left: left, Right: right})
	} else {
		expr = Expression(Binary{Left: left, Right: right, Operator: operator})
	}
	return &expr, nil
}
