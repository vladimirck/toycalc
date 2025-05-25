// core.go
package main

import "fmt"

// Epsilon constant for floating-point comparisons
const Epsilon = 1e-10

// TokenType identifies the type of a token
type TokenType string

// TokenTypes for Stage 1
const (
	// Special tokens
	ILLEGAL TokenType = "ILLEGAL" // Unrecognized token/character
	EOF     TokenType = "EOF"     // End Of Input

	// Literals
	NUMBER TokenType = "NUMBER" // e.g., 123, 3.14
	IDENT  TokenType = "IDENT"  // For function names like log, exp

	// Operators
	PLUS        TokenType = "+"
	MINUS       TokenType = "-"
	ASTERISK    TokenType = "*"
	SLASH       TokenType = "/"
	PERCENT     TokenType = "%"           // Modulo
	CARET       TokenType = "^"           // Power
	UNARY_MINUS TokenType = "UNARY_MINUS" // Or UMINUS
	UNARY_PLUS  TokenType = "UNARY_PLUS"  // Or UMINUS

	// Delimiters
	LPAREN   TokenType = "(" // Left Parenthesis
	RPAREN   TokenType = ")" // Right Parenthesis
	LBRACKET TokenType = "[" // Left Bracket
	RBRACKET TokenType = "]" // Right Bracket
	LBRACE   TokenType = "{" // Left Brace
	RBRACE   TokenType = "}" // Right Brace
	COMMA    TokenType = "," // For function arguments (though not heavily used in Stage 1 funcs)
)

// Token represents a lexical unit
type Token struct {
	Type     TokenType
	Literal  string // The literal value of the token
	Position int    // for detailed error reporting
}

// CalculationError (as defined in Stage 0)
type CalculationError struct {
	Message string
}

func (e *CalculationError) Error() string {
	return fmt.Sprintf("Calculation error: %s", e.Message)
}

func NewCalculationError(message string) error {
	return &CalculationError{Message: message}
}
