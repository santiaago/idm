package main

import "fmt"

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
func main() {
	fmt.Println("hello, world")
}
