package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"
)

// Token identifies the type of lex items.
type Token int

const (
	// EOF represents the end of file
	EOF Token = iota
	// Error represents an error
	Error
	// Assign represents the assignment '='
	Assign
	// Number represents a simple number
	Number
	// Operator an operator such as '+' '-' '*'
	Operator
	// Asterix the multiplication operator
	Asterix
	// Space represents space separation between tokens
	Space
	// Identifier represent an identifier such as a var name
	Identifier
)

// eof rune to treat EOF like any other character
var eof = rune(0)

func isWhitespace(r rune) bool {
	return r == ' ' || r == '\t' || r == 'n'
}

func isLetter(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')
}

func isDigit(r rune) bool {
	return (r >= '0' && r <= '9')
}

// Scanner represents a lexical scanner
type Scanner struct {
	r *bufio.Reader
}

// NewScanner returns a new instance of Scanner.
func NewScanner(r io.Reader) *Scanner {
	return &Scanner{r: bufio.NewReader(r)}
}

// read reads the next rune from the bufferred reader.
// Returns the rune(0) if an error occurs (or io.EOF is returned)
func (s *Scanner) read() rune {
	r, _, err := s.r.ReadRune()
	if err != nil {
		return eof
	}
	return r
}

func (s *Scanner) unread() { _ = s.r.UnreadRune() }

// Scan retusn the next token and literal value.
func (s *Scanner) Scan() (t Token, lit string) {
	// read the next rune.
	r := s.read()

	// if we see whitespace then consume all contiguous whitespace.
	// if we see a letter then consume as an identifier keyword word.
	if isWhitespace(r) {
		s.unread()
		return s.scanWhitespace()
	} else if isLetter(r) {
		s.unread()
		return s.scanIdentifier()
	}

	switch r {
	case eof:
		return EOF, ""
	case '*':
		return Asterix, string(r)
	case '=':
		return Assign, string(r)
	}
	return Error, string(r)
}

// scanWhitespace consumes the current rune and all contiguous whitespace.
func (s *Scanner) scanWhitespace() (t Token, lit string) {
	// Create a buffer and read the current character into it.
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	// Read every subsequent whitespace character into the bufer.
	// non whitespace characters and EOF will cause the loop to exit.
	for {
		if r := s.read(); r == eof {
			break
		} else if !isWhitespace(r) {
			s.unread()
			break
		} else {
			buf.WriteRune(r)
		}
	}
	return Space, buf.String()
}

func (s *Scanner) scanIdentifier() (t Token, lit string) {
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	for {
		if r := s.read(); r == eof {
			break
		} else if !isLetter(r) && !isDigit(r) && r != '_' {
			s.unread()
			break
		} else {
			_, _ = buf.WriteRune(r)
		}
	}
	return Identifier, buf.String()
}

// Parser represents a parser
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
	// if we have a token on the buffer, then return it
	if p.buf.n != 0 {
		p.buf.n = 0
		return p.buf.t, p.buf.lit
	}
	// otherwise read the next token from the scanner
	t, lit = p.s.Scan()

	// save it to the buffer in case we unscan later.
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
	Rigth string
}

// Parse parse a assign statement a = b
func (p *Parser) Parse() (*Statement, error) {
	stmt := &Statement{}

	// Next we should loop over all our comma-delimited fields.
	// Read a field.
	tok, lit := p.scanIgnoreWhitespace()
	if tok != Identifier {
		return nil, fmt.Errorf("found %q, expected left", lit)
	}
	stmt.Left = lit

	// Next we should see the "=" keyword.
	if tok, lit := p.scanIgnoreWhitespace(); tok != Assign {
		return nil, fmt.Errorf("found %q, expected '=' with tok: %v, expected %v", lit, tok, Assign)
	}

	// Finally we should read the left var/number name.
	tok, lit = p.scanIgnoreWhitespace()
	if tok != Identifier {
		return nil, fmt.Errorf("found %q, expected identifier name", lit)
	}
	stmt.Rigth = lit

	// Return the successfully parsed statement.
	return stmt, nil
}

func main() {
	fmt.Println("hello, world")
	s := "a = b"
	stmt, err := NewParser(strings.NewReader(s)).Parse()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%+v", stmt)
}
