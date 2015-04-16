package main

import (
	"fmt"
	"io"
	"strconv"
)

var (
	// stack is used to store variable names and values.
	stack map[string]Value
)

func init() {
	stack = make(map[string]Value)
}

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

// numberOrVector returns a number or a vector.
func (p *Parser) numberOrVector() Value {
	tok, lit := p.scanIgnoreWhitespace()
	var vector Vector
	if tok == Number {
		vector = append(vector, ValueParse(lit))
	} else if tok == Operator && lit == "-" {
		// we use scan here because the '-' sign number must be
		// right next to number, no space in between
		tok, lit = p.scan()
		if tok == Number {
			vector = append(vector, ValueParse("-"+lit))
		} else {
			fmt.Printf("ERROR function does not support token %v.\n", tok)
			return nil
		}
	} else {
		fmt.Println("ERROR function does not support token %v yet.", tok)
	}

	for {
		// Read a field.
		tok, lit := p.scanIgnoreWhitespace()
		if tok != Number {
			p.unscan()
			break
		}
		v := ValueParse(lit)
		vector = append(vector, v)
	}

	// todo(santiaago) do we need this?
	if len(vector) == 1 {
		return vector[0]
	}
	return vector
}

// Expression is the interface for an expression.
type Expression interface {
	String() string
	Evaluate() Value
}

// Statement represents a code statement a = 2.
type Statement struct {
	Left  Value
	Right Value
}

// String retuns the statement as a string.
func (s Statement) String() string {
	return fmt.Sprintf("%v%v", s.Left, s.Right)
}

// Evaluate evaluates the given statement.
func (s Statement) Evaluate() Value {
	// todo(santiaago): need to think what this should really do..
	return s.Left
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
func (n Num) Evaluate() Value {
	return n
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
	}
	return nil
}

// times performs a 'a' + 'b' operation and returns it.
func add(a, b Value) Value {
	// todo(santiaago): will have to check types at some point
	if _, ok := a.(Int); ok {
		return Int(a.(Int) + b.(Int))
	}
	if _, ok := a.(Vector); ok {
		var v Vector
		for i := 0; i < len(a.(Vector)); i++ {
			v = append(v, add(a.(Vector)[i], b.(Vector)[i]))
		}
		return v
	}
	fmt.Println("ERROR add: case not supported")
	return nil
}

// times performs a 'a' - 'b' operation and returns it.
func minus(a, b Value) Value {
	if _, ok := a.(Int); ok {
		return Int(a.(Int) - b.(Int))
	}
	if _, ok := a.(Vector); ok {
		var v Vector
		for i := 0; i < len(a.(Vector)); i++ {
			v = append(v, minus(a.(Vector)[i], b.(Vector)[i]))
		}
		return v
	}
	return nil
}

// times performs a 'a' * 'b' operation and returns it.
func times(a, b Value) Value {
	if _, ok := a.(Int); ok {
		return Int(a.(Int) * b.(Int))
	}
	if _, ok := a.(Vector); ok {
		var v Vector
		for i := 0; i < len(a.(Vector)); i++ {
			v = append(v, times(a.(Vector)[i], b.(Vector)[i]))
		}
		return v
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

// Parse parse a assign statement a = b
func (p *Parser) Parse() (*Expression, error) {

	// First token can be an identifier or number(a number can start with a '-' sign)
	var left Value
	var right, operator string
	tok, lit := p.scanIgnoreWhitespace()
	lastTok := tok
	if tok == Identifier {
		left = Variable{name: lit}
	} else if tok == Number {
		p.unscan()
		left = p.numberOrVector()
	} else if tok == Operator {
		p.unscan()
		left = p.numberOrVector()
		lastTok = Number
		// todo(santiaago): need to handle negative identifiers and vectors.
	} else {
		return nil, fmt.Errorf("ERROR found %q, expected left", lit)
	}

	// Next it could be EOF, an operator or an assignment (for now)
	tok, lit = p.scanIgnoreWhitespace()
	if tok == EOF {
		if lastTok == Number {
			expr := Expression(left)
			return &expr, nil
		} else if lastTok == Identifier {
			// todo(santiaago): do we need this error check?
			if left.Evaluate() == nil {
				return nil, fmt.Errorf("ERROR")
			}
			expr := Expression(left)
			return &expr, nil
		} else {
			fmt.Println("ERROR not a number or identifier ", lastTok)
		}
	}

	isOperator := tok == Operator
	isAssign := tok == Assign
	if !isAssign && !isOperator {
		return nil, fmt.Errorf("ERROR found %q, expected '=' with tok: %v, expected %v", lit, tok, Assign)
	}

	if isOperator {
		operator = lit
	}

	// Next: Take care of assign case.
	if isAssign {
		tok, lit = p.scanIgnoreWhitespace()
		if tok == Number {
			// todo(santiaago):
			// we should print an error here if left is a number and
			// not a variable as the asignment 1 = 2 doesn't make sense.
			var expr Expression
			if v, ok := left.(Variable); ok {
				stack[v.name] = ValueParse(lit)
				expr = Expression(v)
			} else {
				fmt.Println("ERROR left hand side should be a variable.")
			}
			return &expr, nil
		}
		// identifier case
		// todo(santiaago): should add an token check here.
		right = lit
		var expr Expression
		var r Value
		if val, ok := stack[right]; ok {
			r = val
		} else {
			r = ValueParse(right)
		}

		if v, ok := left.(Variable); ok {
			stack[v.name] = r
			expr = Expression(v)
		}
		return &expr, nil
	}

	// Next: Take care of the operator case.
	// We should loop over all our operators.
	var terms []Value
	var operators []string

	// Initialize arrays with first term and first operator.
	terms = append(terms, left)
	operators = append(operators, operator)

	for {
		// Read a field.
		tok, lit = p.scanIgnoreWhitespace()
		if tok != Identifier && tok != Number {
			return nil, fmt.Errorf("found %q, expected number or identifier", lit)
		}
		// todo(santiaago): should check here cases (number, identifier)
		if tok == Number {
			p.unscan()
			term := p.numberOrVector()
			terms = append(terms, term)
		} else if tok == Identifier {
			if _, ok := stack[lit]; ok {
				terms = append(terms, stack[lit])
			} else {
				return nil, fmt.Errorf("ERROR variable %v not found", lit)
			}
		} else {
			return nil, fmt.Errorf("ERROR unexpected token %v", tok)
		}
		// Read operator
		tok, lit = p.scanIgnoreWhitespace()
		// If the next token is not an operator then break the loop
		if tok != Operator {

			p.unscan()
			break
		}
		operators = append(operators, lit)
	}

	// At this point we have terms and, operators.
	// We now need to process all of this.
	return buildOperatorExpression(terms, operators)
}

// At this point we have the following
// left, operator, terms, operators
// we need now to process all of this.
func buildOperatorExpression(terms []Value, operators []string) (*Expression, error) {

	if len(terms)-1 != len(operators) {
		return nil, fmt.Errorf("ERROR terms and operators size mismatch")
	}

	var cumulExpr Expression
	first := terms[0]
	// todo(santiaago):
	// you should be able to do:
	//     a = 1
	// 1
	//     1 2 a
	// 1 2 1
	// This means that a vector can also contain identifiers.
	// we will ignore this case for now, we will only work with vector of numbers.
	// if _, ok := stack[first]; ok {
	// 	cumulExpr = Expression(Variable{name: first})
	// } else {
	cumulExpr = first
	//}

	for i := 0; i < len(operators); i++ {

		right := terms[i+1]
		op := operators[i]

		var r Value
		// if val, ok := stack[right]; ok {
		// 	r = val
		// } else {
		r = right
		//}
		if cumulExpr == nil {
			fmt.Println("ERROR nil expression.")
		}
		b := Binary{Left: cumulExpr.Evaluate(), Right: r, Operator: op}
		cumulExpr = Expression(b)
	}
	return &cumulExpr, nil
}
