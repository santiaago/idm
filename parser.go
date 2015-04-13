package idm

import (
	"fmt"
	"io"
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

// Statement represents a code statement a = 2.
type Statement struct {
	Left  string
	Right string
}

// Parse parse a assign statement a = b
func (p *Parser) Parse() (*Statement, error) {
	stmt := &Statement{}

	tok, lit := p.scanIgnoreWhitespace()
	if tok != Identifier {
		return nil, fmt.Errorf("found %q, expected left", lit)
	}
	stmt.Left = lit

	if tok, lit := p.scanIgnoreWhitespace(); tok != Assign {
		return nil, fmt.Errorf("found %q, expected '=' with tok: %v, expected %v", lit, tok, Assign)
	}

	tok, lit = p.scanIgnoreWhitespace()
	if tok != Identifier {
		return nil, fmt.Errorf("found %q, expected identifier name", lit)
	}
	stmt.Right = lit

	return stmt, nil
}
