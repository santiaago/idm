package main

// Token represents a lexical token
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
	Identifier
)
