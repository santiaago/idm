package idm

import (
	"strings"
	"testing"
)

func TestScan(t *testing.T) {
	var tests = []struct {
		s   string
		tok Token
		lit string
	}{
		{s: ``, tok: EOF},
		{s: `#`, tok: Error, lit: `#`},
		{s: ` `, tok: Space, lit: ` `},
		{s: "\t", tok: Space, lit: "\t"},
		{s: "\n", tok: Space, lit: "\n"},
		{s: `*`, tok: Asterix, lit: `*`},
		{s: `=`, tok: Assign, lit: `=`},
		{s: `a`, tok: Identifier, lit: `a`},
		{s: `a42`, tok: Identifier, lit: `a42`},
		{s: `a_42`, tok: Identifier, lit: `a_42`},
	}
	for i, tt := range tests {
		s := NewScanner(strings.NewReader(tt.s))
		tok, lit := s.Scan()
		if tt.tok != tok {
			t.Errorf("%d. %q token mismatch: exp=%q got=%q <%q>", i, tt.s, tt.tok, tok, lit)
		} else if tt.lit != lit {
			t.Errorf("%d. %q literal mismatch: exp=%q got=%q", i, tt.s, tt.lit, lit)
		}
	}
}

func TestScanWhitespace(t *testing.T) {
	var tests = []struct {
		s   string
		tok Token
		lit string
	}{
		{s: ``, tok: Space, lit: "\x00"}, // a byte with the value 0
		{s: ` #`, tok: Space, lit: ` `},
		{s: `  `, tok: Space, lit: "  "},
		{s: "\t", tok: Space, lit: "\t"},
		{s: "\n", tok: Space, lit: "\n"},
		{s: "a", tok: Space, lit: "a"},
		{s: " a", tok: Space, lit: " "},
		{s: "   a", tok: Space, lit: "   "},
		{s: "a ", tok: Space, lit: "a "},
		{s: "a  ", tok: Space, lit: "a  "},
		{s: "a  = b", tok: Space, lit: "a  "},
		{s: "  a  ", tok: Space, lit: "  "},
	}

	for i, tt := range tests {
		s := NewScanner(strings.NewReader(tt.s))

		tok, lit := s.scanWhitespace()
		if tt.tok != tok {
			t.Errorf("%d. %q token mismatch: exp=%q got=%q <%q>", i, tt.s, tt.tok, tok, lit)
		} else if tt.lit != lit {
			t.Errorf("%d. %q literal mismatch: exp=%q got=%q", i, tt.s, tt.lit, lit)
		}
	}

}

func TestScanIdentifier(t *testing.T) {
	var tests = []struct {
		s   string
		tok Token
		lit string
	}{
		{s: ``, tok: Identifier, lit: "\x00"}, // a byte with the value 0
		{s: ` #`, tok: Identifier, lit: ` `},
		{s: `  `, tok: Identifier, lit: " "},
		{s: "\t", tok: Identifier, lit: "\t"},
		{s: "\n", tok: Identifier, lit: "\n"},
		{s: "a", tok: Identifier, lit: "a"},
		{s: " a", tok: Identifier, lit: " a"},
		{s: "   a", tok: Identifier, lit: " "},
		{s: "a ", tok: Identifier, lit: "a"},
		{s: "a  ", tok: Identifier, lit: "a"},
		{s: "a  = b", tok: Identifier, lit: "a"},
		{s: "  a  ", tok: Identifier, lit: " "},
		{s: "aaa  ", tok: Identifier, lit: "aaa"},
		{s: " aaa", tok: Identifier, lit: " aaa"},
		{s: "ab123", tok: Identifier, lit: "ab123"},
		{s: "ab_123", tok: Identifier, lit: "ab_123"},
	}

	for i, tt := range tests {
		s := NewScanner(strings.NewReader(tt.s))

		tok, lit := s.scanIdentifier()
		if tt.tok != tok {
			t.Errorf("%d. %q token mismatch: exp=%q got=%q <%q>", i, tt.s, tt.tok, tok, lit)
		} else if tt.lit != lit {
			t.Errorf("%d. %q literal mismatch: exp=%q got=%q", i, tt.s, tt.lit, lit)
		}
	}
}

func TestIsWhitespace(t *testing.T) {
	var tests = []struct {
		r        rune
		expected bool
	}{
		{r: ' ', expected: true},
		{r: '\t', expected: true},
		{r: '\n', expected: true},
		{r: 'a', expected: false},
	}
	for i, tt := range tests {
		got := isWhitespace(tt.r)
		if got != tt.expected {
			t.Errorf("%d. %t expected, got %t", i, tt.expected, got)
		}
	}
}

func TestIsLetter(t *testing.T) {
	var tests = []struct {
		r        rune
		expected bool
	}{
		{r: ' ', expected: false},
		{r: '\t', expected: false},
		{r: '0', expected: false},
		{r: 'a', expected: true},
		{r: 'A', expected: true},
	}
	for i, tt := range tests {
		got := isLetter(tt.r)
		if got != tt.expected {
			t.Errorf("%d. %t expected, got %t", i, tt.expected, got)
		}
	}
}

func TestIsDigit(t *testing.T) {
	var tests = []struct {
		r        rune
		expected bool
	}{
		{r: ' ', expected: false},
		{r: '\t', expected: false},
		{r: '0', expected: true},
		{r: 'a', expected: false},
		{r: 'A', expected: false},
	}
	for i, tt := range tests {
		got := isDigit(tt.r)
		if got != tt.expected {
			t.Errorf("%d. %t expected, got %t", i, tt.expected, got)
		}
	}
}
