package main

import (
	"reflect"
	"strings"
	"testing"
)

func TestParser_Scan(t *testing.T) {
	var tests = []struct {
		s    string
		stmt *Statement
		err  string
	}{
		{
			s: `a = b`,
			stmt: &Statement{
				Left:  "a",
				Right: "b",
			},
		},
		{
			s:   `a`,
			err: `found "", expected '='`,
		},
	}

	for i, tt := range tests {
		stmt, err := NewParser(strings.NewReader(tt.s)).Parse()
		if !strings.Contains(errstring(err), tt.err) {
			t.Errorf("%d. %q: error mismatch:\n  exp=%s\n  got=%s\n\n", i, tt.s, tt.err, err)
		} else if tt.err == "" && !reflect.DeepEqual(tt.stmt, stmt) {
			t.Errorf("%d. %q\n\nstmt mismatch:\n\nexp=%#v\n\ngot=%#v\n\n", i, tt.s, tt.stmt, stmt)
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
