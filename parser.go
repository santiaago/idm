package main

import (
	"fmt"
	"io"
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
		t    []Token  // stack of last read tokens
		lit  []string // stack of last read literals
		n    int      // buffer size (max=1)
		size int      // stack size for 't' and 'lit'
	}
}

// NewParser returns a new instance of Parser.
func NewParser(r io.Reader) *Parser {
	p := Parser{s: NewScanner(r)}
	p.buf.size = 10
	return &p
}

// scan returns the next token from the underlying scanner.
// if a token has been unscanned then read that instead.
func (p *Parser) scan() (t Token, lit string) {
	if p.buf.n != 0 {
		// todo(santiaago): refactor
		t, p.buf.t = p.buf.t[len(p.buf.t)-1], p.buf.t[:len(p.buf.t)-1]
		lit, p.buf.lit = p.buf.lit[len(p.buf.lit)-1], p.buf.lit[:len(p.buf.lit)-1]
		p.buf.n--
		return
	}

	t, lit = p.s.Scan()
	if len(p.buf.t) < p.buf.size {
		p.buf.t = append(p.buf.t, t)
		p.buf.lit = append(p.buf.lit, lit)
	} else {
		// stack limit reached so shift values and insert new ones
		// todo(santiaago): refactor
		p.buf.t = p.buf.t[1:]
		p.buf.t = append(p.buf.t, t)
		p.buf.lit = p.buf.lit[1:]
		p.buf.lit = append(p.buf.lit, lit)
	}
	return
}

// unscan pushes the previously read token back onto the buffer.
func (p *Parser) unscan() {
	if p.buf.n == p.buf.size {
		fmt.Println("ERROR cannot unscan anymore, stack size limit reached.")
		return
	}
	p.buf.n++
}

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
		// todo(santiaago): handle this as a unary operator.
		// we use scan here because the '-' sign number must be
		// right next to number, no space in between
		tok, lit = p.scan()
		if tok == Number {
			vector = append(vector, ValueParse("-"+lit))
		} else {
			return nil
		}
	} else {
		return nil
	}

	for {
		// Read a field.
		tok, lit := p.scanIgnoreWhitespace()
		if tok == Operator && lit == "-" {
			// we use scan here because the '-' sign number must be
			// right next to number, no space in between
			tok, lit = p.scan()
			if tok == Number {
				lit = "-" + lit
			} else {
				// unscan twice to roll back both scan
				p.unscan()
				p.unscan()
				break
			}
		} else if tok != Number {
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

// Parse parse a assign statement a = b
func (p *Parser) Parse() (*Expression, error) {

	// First token can be an identifier or number(a number can start with a '-' sign)
	// todo(santiaago): refactor first token.
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
		if lit == "-" {
			p.unscan()
			left = p.numberOrVector()
			if left != nil {
				lastTok = Number
			}
		} else if isUnary(lit) {
			// remember operator
			p.unscan()
			operator = lit
		}
		// todo(santiaago): need to handle negative identifiers and vectors.
	} else {
		return nil, fmt.Errorf("ERROR found %q, expected left", lit)
	}

	// Next it could be EOF, an operator or an assignment (for now)
	// todo(santiaago): refactor EOF case.
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
			return nil, fmt.Errorf("ERROR not a number or identifier")
		}
	}

	isOperator := tok == Operator
	isAssign := tok == Assign
	isNumber := tok == Number
	if !isAssign && !isOperator && !isNumber {
		return nil, fmt.Errorf("ERROR found %q, expected assignment, operator or number. Got token %v", lit, tok)
	}

	// if the literal scanned was a number we unscan it to scan it completly.
	// This is to scan all numbers in a vector and don't skip the first one.
	if isNumber {
		p.unscan()
	}

	// if last token read is an operator, remember it.
	if isOperator {
		operator = lit
	}

	// Next: Take care of assign case.
	// todo(santiaago): refactor assign case.
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
				return nil, fmt.Errorf("ERROR left hand side should be a variable")
			}
			return &expr, nil
		} else if lastTok == Number {
			return nil, fmt.Errorf("ERROR left hand side should be a variable")
		} else if tok == Identifier {
			right = lit
			var expr Expression
			var r Value
			if val, ok := stack[right]; ok {
				r = val
			} else {
				return nil, fmt.Errorf("ERROR variable undefined")
			}
			if v, ok := left.(Variable); ok {
				stack[v.name] = r
				expr = Expression(v)
			}
			return &expr, nil
		}
		return nil, fmt.Errorf("ERROR right hand side has unexpected token")
	}

	// Next: Take care of the operator case.
	// todo(santiaago): refactor this.
	// We should loop over all our operators.
	var terms []Value
	var operators []string

	// Initialize arrays with first term and first operator.
	if left != nil {
		terms = append(terms, left)
	}
	operators = append(operators, operator)

	for {
		// Read a field.
		tok, lit = p.scanIgnoreWhitespace()
		if tok != Identifier && tok != Number && tok != Operator {
			return nil, fmt.Errorf("ERROR found %q, expected number or identifier or sign", lit)
		}
		// todo(santiaago): should check here cases (number, identifier, sign)
		if tok == Number || tok == Operator {
			p.unscan()
			term := p.numberOrVector()

			terms = append(terms, term)
		} else if tok == Identifier {
			if _, ok := stack[lit]; ok {
				terms = append(terms, stack[lit])
			} else {
				return nil, fmt.Errorf("ERROR variable %v not found", lit)
			}
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
		op := operators[i]
		// unary case
		// todo(santiaago): how to handle +/ 1 2 3 + +/ 1 2 3 ?
		if isUnary(op) {
			u := Unary{Val: cumulExpr.Evaluate(), Operator: op}
			cumulExpr = Expression(u)
			continue
		}
		// todo(santiaago): need to clean this.
		if len(terms)-1 != len(operators) {
			return nil, fmt.Errorf("ERROR terms and operators size mismatch")
		}
		right := terms[i+1]

		var r Value
		// if val, ok := stack[right]; ok {
		// 	r = val
		// } else {
		r = right
		//}
		if cumulExpr != nil {
			b := Binary{Left: cumulExpr.Evaluate(), Right: r, Operator: op}
			cumulExpr = Expression(b)
		} else {
			return nil, fmt.Errorf("ERROR nil expression")
		}
	}
	return &cumulExpr, nil
}
