// evaluator.go
package main

import (
	"fmt"
	"math" // Needed for Abs, Floor, Round in the future
	// Needed for complex operations in Stage 1
	"strconv" // To convert string to float64
)

// calculateExpression is the entry point for calculation logic
// In Stage 0, it's a stub that attempts a simple conversion.
func calculateExpression(expressionString string) (string, error) {
	// Full flow (Later Stages):
	// 1. tokens, err := Lex(expressionString)
	// 2. rpnQueue, err := Parse(tokens)
	// 3. resultComplex, err := EvaluateRPN(rpnQueue)
	// 4. return formatComplexOutput(resultComplex), nil

	// Stub for Stage 0: Try to convert the expression to float64 and display it as complex.
	// This is to verify the basic flow of main.go.
	f, err := strconv.ParseFloat(expressionString, 64)
	if err != nil {
		// If not a simple number, indicate that full functionality is not implemented.
		return "", NewCalculationError(fmt.Sprintf("direct evaluation of '%s' not supported in Stage 0. Full implementation in Stage 1.", expressionString))
	}

	// Simulate a complex128 result (imaginary part zero for now)
	resultComplex := complex(f, 0)

	// Output formatting (very simplified version for Stage 0)
	// The full version with Epsilon will come in Stage 1.
	return formatComplexOutput(resultComplex) // Use the stub formatter
}

// EvaluateRPN (stub) - Evaluates the RPN queue
// Will be fully implemented in Stage 1
func EvaluateRPN(rpnQueue []Token) (complex128, error) {
	// RPN evaluator logic with complex128 stack
	return complex(0, 0), NewCalculationError("RPN evaluator not yet implemented.")
}

// formatComplexOutput (stub) - Formats the complex128 output
// Will be implemented with Epsilon logic in Stage 1
func formatComplexOutput(c complex128) (string, error) {
	realPart := real(c)
	imagPart := imag(c)

	if math.Abs(imagPart) < Epsilon { // Epsilon defined in core.go
		return fmt.Sprintf("%g", realPart), nil
	}
	// Format for complex numbers, e.g., "3+2i" or "3-2i"
	// Using %+g for imagPart to ensure the sign is always present.
	return fmt.Sprintf("%g%+gi", realPart, imagPart), nil
}
