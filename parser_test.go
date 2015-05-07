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
		{s: `a = 1`, expr: Variable{name: "a"}},
		{s: `1 + 2`, expr: Binary{Left: Int(1), Right: Int(2), Operator: "+"}},
		{s: `1 / 2`, expr: Binary{Left: Int(1), Right: Int(2), Operator: "/"}},
		{s: `1 - 2`, expr: Binary{Left: Int(1), Right: Int(2), Operator: "-"}},
		{s: `1* 2`, expr: Binary{Left: Int(1), Right: Int(2), Operator: "*"}},
		{s: `2 ** 2`, expr: Binary{Left: Int(2), Right: Int(2), Operator: "**"}},
		{s: `a + 2`, expr: Binary{Left: Int(1), Right: Int(2), Operator: "+"}},
		{s: `2 + a`, expr: Binary{Left: Int(2), Right: Int(1), Operator: "+"}},
		{s: `a + a`, expr: Binary{Left: Int(1), Right: Int(1), Operator: "+"}},
		{s: `a + a - a + a`, expr: Int(2)},
		{s: `a + a + a + a`, expr: Int(4)},
		{s: `b = 42`, expr: Variable{name: "b"}},
		{s: `a = b`, expr: Variable{name: "a"}},
		{s: `1 2 3 4`, expr: Vector([]Value{Int(1), Int(2), Int(3), Int(4)})},
		{
			s: `1 2 3 4 + 1 2 3 4`,
			expr: Binary{
				Left:     Vector([]Value{Int(1), Int(2), Int(3), Int(4)}),
				Right:    Vector([]Value{Int(1), Int(2), Int(3), Int(4)}),
				Operator: "+"},
		},
		{
			s: `1 2 3 4 - 1 2 3 4`,
			expr: Binary{
				Left:     Vector([]Value{Int(1), Int(2), Int(3), Int(4)}),
				Right:    Vector([]Value{Int(1), Int(2), Int(3), Int(4)}),
				Operator: "-"},
		},
		{
			s: `1 2 3 4 * 1 2 3 4`,
			expr: Binary{
				Left:     Vector([]Value{Int(1), Int(2), Int(3), Int(4)}),
				Right:    Vector([]Value{Int(1), Int(2), Int(3), Int(4)}),
				Operator: "*"},
		},
		{s: `c`, err: `ERROR`},
	}

	for i, tt := range tests {
		expr, err := NewParser(strings.NewReader(tt.s)).Parse()
		if !strings.Contains(errstring(err), tt.err) {
			t.Errorf("%d. %q: error mismatch:\n  exp=%s\n  got=%s\n\n", i, tt.s, tt.err, err)
		} else if tt.err == "" &&
			!reflect.DeepEqual(tt.expr.Evaluate(), (*expr).Evaluate()) {
			t.Errorf("%d. %q\n\nstmt mismatch:\n\nexp=%#v\n\ngot=%#v", i, tt.s, tt.expr.Evaluate(), (*expr).Evaluate())
		}
	}
}

func TestParser_NumberValues(t *testing.T) {
	var tests = []struct {
		s    string
		expr Expression
		err  string
	}{
		{s: `1`, expr: Int(1)},
		{s: `-1`, expr: Int(-1)},
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

func TestParser_Errors(t *testing.T) {
	var tests = []struct {
		s    string
		expr Expression
		err  string
	}{
		{s: `- 1`, err: `ERROR`},
		{s: `? 1`, err: `ERROR`},
		{s: `-`, err: `ERROR`},
		{s: `2 ?`, err: `ERROR`},
		{s: `2 = 2`, err: `ERROR`},
		{s: `2 = a`, err: `ERROR`},
		{s: `a = 1`, expr: Int(1)},
		{s: `a = ?`, err: `ERROR`},
		{s: `a = c`, err: `ERROR`},
		{s: `a + ?`, err: `ERROR`},
		{s: `a + c`, err: `ERROR`},
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

func TestParser_VariableValues(t *testing.T) {
	var tests = []struct {
		s    string
		expr Expression
		err  string
	}{
		{s: `a = 1`, expr: Int(1)},
		{s: `a`, expr: Int(1)},
		{s: `b = 42`, expr: Int(42)},
		{s: `a = b`, expr: Int(42)},
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

func TestParser_NumberArithmeticValues(t *testing.T) {
	var tests = []struct {
		s    string
		expr Expression
		err  string
	}{
		{s: `1 + 2`, expr: Int(3)},
		{s: `1 - 2`, expr: Int(-1)},
		{s: `-1 + 2`, expr: Int(1)},
		{s: `-1 + -2`, expr: Int(-3)},
		{s: `1 / 2`, expr: Int(0)},
		{s: `-1 - -2 + -10`, expr: Int(-9)},
		{s: `1* 2`, expr: Int(2)},
		{s: `2 ** 2`, expr: Int(4)},
		{s: `2 max 1`, expr: Int(2)},
		{s: `2 min 1`, expr: Int(1)},
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

func TestParser_VariablerArithmeticValues(t *testing.T) {
	var tests = []struct {
		s    string
		expr Expression
		err  string
	}{
		{s: `a = 1`, expr: Int(1)},
		{s: `2 + a`, expr: Int(3)},
		{s: `a + a`, expr: Int(2)},
		{s: `a + a - a + a`, expr: Int(2)},
		{s: `a + a + a + a`, expr: Int(4)},
		{s: `b = 42`, expr: Int(42)},
		{s: `a = b`, expr: Int(42)},
		{s: `a + b`, expr: Int(84)},
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

func TestParser_VectorValues(t *testing.T) {
	var tests = []struct {
		s    string
		expr Expression
		err  string
	}{
		{s: `1 2 3 4`, expr: Vector([]Value{Int(1), Int(2), Int(3), Int(4)})},
		{s: `-1 -2 -3 -4`, expr: Vector([]Value{Int(-1), Int(-2), Int(-3), Int(-4)})},
		{s: `-1 -2 3 4`, expr: Vector([]Value{Int(-1), Int(-2), Int(3), Int(4)})},
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

func TestParser_VectorArithmeticValues(t *testing.T) {
	var tests = []struct {
		s    string
		expr Expression
		err  string
	}{
		{
			s:    `1 2 3 4 + 1 2 3 4`,
			expr: Vector([]Value{Int(2), Int(4), Int(6), Int(8)}),
		},
		{
			s:    `1 2 3 4 - 1 2 3 4`,
			expr: Vector([]Value{Int(0), Int(0), Int(0), Int(0)}),
		},
		{
			s: `1 2 3 4 * 1 2 3 4`,

			expr: Vector([]Value{Int(1), Int(4), Int(9), Int(16)}),
		},
		{
			s:    `1 2 3 4 ** 2 2 2 2`,
			expr: Vector([]Value{Int(1), Int(4), Int(9), Int(16)}),
		},
		{
			s:    `1 2 3 4 min 2 2 2 2`,
			expr: Vector([]Value{Int(1), Int(2), Int(2), Int(2)}),
		},
		{
			s:    `1 2 3 4 max 2 2 2 2`,
			expr: Vector([]Value{Int(2), Int(2), Int(3), Int(4)}),
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

func TestParser_ScanOperationsValues(t *testing.T) {
	var tests = []struct {
		s    string
		expr Expression
		err  string
	}{
		{s: `+\ 1`, expr: Int(1)},
		{s: `+/ 1`, expr: Int(1)},
		{s: `*\ 1`, expr: Int(1)},
		{s: `*/ 1`, expr: Int(1)},
		{
			s:    `+\ 1 2 3 4`,
			expr: Vector([]Value{Int(1), Int(3), Int(6), Int(10)}),
		},
		{s: `+/ 1 2 3 4`, expr: Int(10)},
		{
			s:    `*\ 1 2 3 4`,
			expr: Vector([]Value{Int(1), Int(2), Int(6), Int(24)}),
		},
		{s: `*/ 1 2 3 4`, expr: Int(24)},
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
