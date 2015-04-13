package idm

import (
	"strings"
	"testing"
)

func TestScanWhitespace(t *testing.T) {
	var tests = []struct {
		s   string
		tok Token
		lit string
	}{
		{s: ``, tok: Space, lit: "\x00"}, // describes the byte with the value 0
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
