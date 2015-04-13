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
	// Space represents space separation between tokens
	Space
	// Identifier represent an identifier such as a var name
)

// eof rune to treat EOF like any other character
var eof = rune(0)

func isWhitespace(r rune) bool {
	return r == ' ' || r == '\t' || r == 'n'
}

func isLetter(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')
}

func main() {
	fmt.Println("hello, world")
}
