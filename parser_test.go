package main

import (
	"reflect"
	"strings"
	"testing"
)

func TestParser_Scan(t *testing.T) {
	var tests = []struct {
		s    string
		expr Expression
		err  string
	}{
		{
			s:    `a = 1`,
			expr: Int(1),
		},
		{
			s: `1 + 2`,
			expr: Binary{
				Left:     Int(1),
				Right:    Int(2),
				Operator: "+",
			},
		},
		{
			s: `1 - 2`,
			expr: Binary{
				Left:     Int(1),
				Right:    Int(2),
				Operator: "-",
			},
		},
		{
			s: `1* 2`,
			expr: Binary{
				Left:     Int(1),
				Right:    Int(2),
				Operator: "*",
			},
		},
		{
			s: `a + 2`,
			expr: Binary{
				Left:     Int(1),
				Right:    Int(2),
				Operator: "+",
			},
		},
		{
			s: `2 + a`,
			expr: Binary{
				Left:     Int(2),
				Right:    Int(1),
				Operator: "+",
			},
		},
		{
			s: `a + a`,
			expr: Binary{
				Left:     Int(1),
				Right:    Int(1),
				Operator: "+",
			},
		},
		{
			s:   `b`,
			err: `ERROR`,
		},
		{
			s:    `a + a + a + a`,
			expr: Int(4),
		},
		{
			s:    `a + a - a + a`,
			expr: Int(2),
		},
		{
			s:    `b = 42`,
			expr: Int(42),
		},
		{
			s:    `a = b`,
			expr: Int(42),
		},
	}

	for i, tt := range tests {
		expr, err := NewParser(strings.NewReader(tt.s)).Parse()
		if !strings.Contains(errstring(err), tt.err) {
			t.Errorf("%d. %q: error mismatch:\n  exp=%s\n  got=%s\n\n", i, tt.s, tt.err, err)
		} else if tt.err == "" && !reflect.DeepEqual(tt.expr.Evaluate(), (*expr).Evaluate()) {
			t.Errorf("%d. %q\n\nstmt mismatch:\n\nexp=%#v\n\ngot=%#v\n\n", i, tt.s, tt.expr.Evaluate(), (*expr).Evaluate())
		}
	}
}

// errstring returns the string representation of an error.
func errstring(err error) string {
	if err != nil {
		return err.Error()
	}
	return ""
}
