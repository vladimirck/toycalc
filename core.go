// core.go
package main

import "fmt"

// Epsilon constant for floating-point comparisons (important for complex128)
const Epsilon = 1e-14 // A common value, adjust if necessary

// TokenType identifies the type of a token (string for easy debugging initially)
type TokenType string

// Initial TokenType definitions (will be expanded in Stage 1)
const (
	ILLEGAL TokenType = "ILLEGAL" // Unrecognized token/character
	EOF     TokenType = "EOF"     // End Of Input

	NUMBER TokenType = "NUMBER" // e.g., 123, 3.14
	// More types like OPERATOR, LPAREN, IDENT, etc., will be added in Stage 1
)

// Token represents a lexical unit
type Token struct {
	Type    TokenType
	Literal string // The literal value of the token (e.g., "123", "+")
	// Position int    // Optional: Starting position in the input string for detailed errors
}

// CalculationError is a custom error type for the calculator
type CalculationError struct {
	Message string
	// Pos     int // Optional: position of the error
}

func (e *CalculationError) Error() string {
	// In the future: return fmt.Sprintf("Calculation error (pos %d): %s", e.Pos, e.Message)
	return fmt.Sprintf("Calculation error: %s", e.Message)
}

// NewCalculationError creates a new CalculationError
func NewCalculationError(message string /*, pos int*/) error {
	return &CalculationError{Message: message /*, Pos: pos*/}
}
