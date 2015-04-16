package main

import (
	"bufio"
	"bytes"
	"io"
)

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

// Scan returns the next token and literal value.
func (s *Scanner) Scan() (t Token, lit string) {
	// read the next rune.
	r := s.read()
	sr := string(r)

	// if we see whitespace then consume all contiguous whitespace.
	// if we see a letter then consume as an identifier keyword word.
	if isWhitespace(r) {
		s.unread()
		return s.scanWhitespace()
	} else if isLetter(r) {
		s.unread()
		t, lit = s.scanIdentifier()
		if !isKeyword(lit) {
			return
		}
		sr = lit
	} else if isDigit(r) {
		s.unread()
		return s.scanDigit()
	}

	// rune cases
	switch r {
	case eof:
		return EOF, ""
	case '+':
		// check +\ and +/ cases
		r1 := s.read()
		if r1 == '\\' {
			return Operator, "+\\"
		} else if r1 == '/' {
			return Operator, "+/"
		}
		s.unread()
		return Operator, string(r)
	case '-':
		return Operator, string(r)
	case '/':
		return Operator, string(r)
	case '*':
		// check power case
		r1 := s.read()
		if r1 == '*' {
			return Operator, "**"
		}
		s.unread()
		return Operator, string(r)
	case '=':
		return Assign, string(r)
	}

	// keyword cases
	switch sr {
	case "max":
		return Operator, sr
	case "min":
		return Operator, sr
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

// scanIdentifier consumes the current rune and all contiguous identifier runes.
func (s *Scanner) scanIdentifier() (t Token, lit string) {
	// Create a buffer and read the current character into it.
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	// Read every subsequent identifier into the buffer.
	// Non indentifier characters and EOF will cause the loop to exit.
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

// scanNumber consumes the current rune and all contiguous digit runes.
func (s *Scanner) scanDigit() (tok Token, lit string) {
	// Create a buffer and rtead the current character into it.
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	// Read every subsequent digit into the buffer.
	// Non digit characters and EOF will cause the loop to exit.
	for {
		if r := s.read(); r == eof {
			break
		} else if !isDigit(r) {
			s.unread()
			break
		} else {
			_, _ = buf.WriteRune(r)
		}
	}
	return Number, buf.String()
}

// eof rune to treat EOF like any other character
var eof = rune(0)

// isWhitespace determines if the rune passed as param corresponds to a whitespace.
// whitespace are ' ', '\t' and '\n'
func isWhitespace(r rune) bool {
	return r == ' ' || r == '\t' || r == '\n'
}

// isLetter determines if the rune passed as param corresponds to a letter.
// whitespace are a-zA-Z
func isLetter(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')
}

// isDigit determines if the rune passed as param corresponds to a digit.
// whitespace are 0-9
func isDigit(r rune) bool {
	return (r >= '0' && r <= '9')
}

func isKeyword(s string) bool {
	return (s == "max") || (s == "min")
}

func isUnary(s string) bool {
	return (s == "+\\") || (s == "+/")
}
