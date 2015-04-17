package main

import (
	"fmt"
	"strconv"
)

// Expression is an interface to wrap objects from the parser.
type Expression interface {
	String() string
	Evaluate() Value
}

// Value is an interface to handle different types.
type Value interface {
	String() string
	Evaluate() Value
}

// Int is a type to handle integers
type Int int64

// String returns the string representation of an integer.
func (i Int) String() string {
	return fmt.Sprintf("%d", i)
}

// Evaluate returns the value of the given integer.
func (i Int) Evaluate() Value {
	return i
}

func tryIntString(s string) (Value, error) {
	i, err := strconv.ParseInt(s, 10, 64)
	return Int(i), err
}

// Vector is a type to handle vectors
type Vector []Value

// String returns the string representation of a vector
func (v Vector) String() string {
	ret := ""
	for i := range v {
		ret += fmt.Sprintf("%v ", v[i])
	}
	ret += fmt.Sprintf("\n")
	return ret
}

// Evaluate returns the value of a given vector.
func (v Vector) Evaluate() Value {
	return v
}

// Variable represents a variable.
type Variable struct {
	name string
}

// String returns the string representation of a variable.
func (v Variable) String() string {
	return v.name
}

// Evaluate returns the value holded by the variable v
func (v Variable) Evaluate() Value {
	if val, ok := stack[v.name]; ok {
		return val
	}
	return nil
}

// Unary represents an unary statement
// example +/ 1 2 3
// example +\ 1 2 3
// todo(santiaago): should add +2 -4 a unary objects
type Unary struct {
	Val      Value
	Operator string
}

// String returns the string representation of a unary type
func (u Unary) String() string {
	return fmt.Sprintf("%v %v", u.Operator, u.Val)
}

// Evaluate returns the return of the operator computed with the value of the
// unary type
func (u Unary) Evaluate() Value {
	if u.Operator == "+/" {
		return sum(u.Val)
	} else if u.Operator == "+\\" {
		return scanSum(u.Val)
	} else if u.Operator == "*/" {
		return multiply(u.Val)
	} else if u.Operator == "*\\" {
		return scanMultiply(u.Val)
	}
	return nil
}

// Binary represents a binary statement
// example: 12 + 3
type Binary struct {
	Left     Value
	Right    Value
	Operator string
}

// String returns the string of the number
func (b Binary) String() string {
	return fmt.Sprintf("%v %v %v", b.Left, b.Operator, b.Right)
}

// Evaluate returns the number value
func (b Binary) Evaluate() Value {
	if b.Operator == "+" {
		return add(b.Left, b.Right)
	} else if b.Operator == "-" {
		return minus(b.Left, b.Right)
	} else if b.Operator == "*" {
		return times(b.Left, b.Right)
	} else if b.Operator == "**" {
		return pow(b.Left, b.Right)
	} else if b.Operator == "max" {
		return max(b.Left, b.Right)
	} else if b.Operator == "min" {
		return min(b.Left, b.Right)
	}
	return nil
}

// ValueParse parse the string in the proper value
// todo(santiaago): will have to try different types here...
func ValueParse(s string) Value {
	v, err := tryIntString(s)
	if err != nil {
		fmt.Printf("ERROR failed to parse %v got error: %v\n", s, err)
		return nil
	}
	return v
}
