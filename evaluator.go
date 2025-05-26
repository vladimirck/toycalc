// evaluator.go
package main

import (
	"fmt"
	"math"
	"math/cmplx"
	"strconv"
	"strings" // For ToLower on function names
)

// calculateExpression orchestrates Lex, Parse, EvaluateRPN, and formatComplexOutput
func calculateExpression(expressionString string) (string, error) {
	tokens, err := Lex(expressionString)
	if err != nil {
		return "", err
	}

	rpnQueue, err := Parse(tokens)
	if err != nil {
		return "", err
	}
	// For debugging RPN output from parser:
	// fmt.Printf("RPN Debug: ")
	// for _, t := range rpnQueue {
	// 	fmt.Printf("{%s '%s' p%d} ", t.Type, t.Literal, t.Position)
	// }
	// fmt.Println()

	resultComplex, err := EvaluateRPN(rpnQueue)
	if err != nil {
		return "", err
	}

	return formatComplexOutput(resultComplex), nil
}

func formatComplexOutput(c complex128) string { // Renamed back if you prefer
	realPart := real(c)
	imagPart := imag(c)

	if cmplx.IsNaN(c) {
		return "NaN"
	}
	if cmplx.IsInf(c) {
		return fmt.Sprintf("%v", c) // Default Go formatting for Inf
	}

	// Check if the imaginary part is negligible
	if math.Abs(imagPart) < Epsilon { // Epsilon defined in core.go
		// Imaginary part is negligible, display only the real part
		// Check if the real part is essentially an integer
		if math.Abs(realPart-math.Round(realPart)) < Epsilon { // Check if realPart is very close to an integer
			return fmt.Sprintf("%.0f", realPart) // Format as integer
		}
		return fmt.Sprintf("%0.10g", realPart) // Format as float (concise)
	}

	if math.Abs(realPart) < Epsilon { // Epsilon defined in core.go
		// Imaginary part is negligible, display only the real part
		// Check if the real part is essentially an integer
		if math.Abs(imagPart-math.Round(imagPart)) < Epsilon { // Check if realPart is very close to an integer
			return fmt.Sprintf("%.0fi", imagPart) // Format as integer
		}
		return fmt.Sprintf("%.10gi", imagPart) // Format as float (concise)
	}

	// Display the full complex number
	// We need to format real and imaginary parts potentially as integers if they are whole numbers
	realStr := ""
	if math.Abs(realPart-math.Round(realPart)) < Epsilon {
		realStr = fmt.Sprintf("%.0f", realPart)
	} else {
		realStr = fmt.Sprintf("%.10g", realPart)
	}

	imagStr := ""
	// For the imaginary part, we always want the sign if it's not part of "0i"
	if math.Abs(imagPart-math.Round(imagPart)) < Epsilon {
		// %.0f for imagPart, then add "i"
		// We need to handle the sign explicitly for the imaginary part if it's formatted as an integer
		if imagPart >= 0 {
			imagStr = fmt.Sprintf("+%.0fi", imagPart)
		} else {
			imagStr = fmt.Sprintf("%.0fi", imagPart) // Negative sign will be included by %.0f
		}
	} else {
		imagStr = fmt.Sprintf("%+.10gi", imagPart) // %+g includes sign and is concise
	}

	// Special case: if real part is exactly 0 and imag part is not negligible
	if math.Abs(realPart) < Epsilon {
		// If imagStr already starts with "+", remove it if we are only printing the imaginary part
		if imagStr[0] == '+' {
			return imagStr[1:] // e.g. "2i" instead of "+2i"
		}
		return imagStr // e.g. "-2i"
	}

	return fmt.Sprintf("%s%s", realStr, imagStr)
}

// Helper for Modulo operation based on Gaussian integer remainder:
// r = a - x*b, where x is the complex integer closest to a/b.
func calculateModulo(a, b complex128, operatorToken Token) (complex128, error) {
	if b == complex(0, 0) {
		// Using cmplx.Abs to catch very small numbers that might behave like zero
		// } else if cmplx.Abs(b) < Epsilon*Epsilon { // Avoid Epsilon itself if b could be Epsilon
		return complex(math.NaN(), math.NaN()), NewCalculationError(
			fmt.Sprintf("divisor is zero for modulo operator at position %d", operatorToken.Position),
		)
	}

	divisionResult := a / b

	// Quotient x is the complex integer closest to a/b
	// math.Round rounds to the nearest even number for .5 cases (e.g., Round(2.5)=2, Round(3.5)=4)
	// This is a standard rounding method.
	roundedQuotient := complex(math.Round(real(divisionResult)), math.Round(imag(divisionResult)))

	remainder := a - (roundedQuotient * b)
	return remainder, nil
}

// EvaluateRPN evaluates a token queue in Reverse Polish Notation
func EvaluateRPN(rpnQueue []Token) (complex128, error) {
	operandStack := []complex128{}

	for _, token := range rpnQueue {
		switch token.Type {
		case NUMBER:
			// The lexer ensures number literals are in a format ParseFloat can handle (incl. scientific)
			val, err := strconv.ParseFloat(token.Literal, 64)
			if err != nil {
				// This error should ideally be caught by the lexer if the number format is truly bad,
				// but strconv.ParseFloat is the ultimate validator.
				return complex(math.NaN(), math.NaN()), NewCalculationError(
					fmt.Sprintf("invalid number format '%s' at position %d", token.Literal, token.Position),
				)
			}
			operandStack = append(operandStack, complex(val, 0))

		case IDENT:
			var result complex128
			processed := false // To track if the IDENT was handled

			lowerLiteral := strings.ToLower(token.Literal)
			switch lowerLiteral {
			// Constants
			case "i":
				result = complex(0, 1)
				operandStack = append(operandStack, result)
				processed = true
			case "pi":
				result = complex(math.Pi, 0)
				operandStack = append(operandStack, result)
				processed = true
			case "e":
				result = complex(math.E, 0)
				operandStack = append(operandStack, result)
				processed = true

			// Stage 1 & 2 Functions (all unary for now)
			case "log", "exp", "sin", "cos", "tan", "asin", "acos", "atan",
				"sinh", "cosh", "tanh", "asinh", "acosh", "atanh",
				"log10", "log2", "sqrt":
				if len(operandStack) < 1 {
					return complex(math.NaN(), math.NaN()), NewCalculationError(
						fmt.Sprintf("insufficient operands for function '%s' at position %d (expected 1)",
							token.Literal, token.Position),
					)
				}
				arg1 := operandStack[len(operandStack)-1]
				operandStack = operandStack[:len(operandStack)-1] // Pop one argument

				switch lowerLiteral { // Inner switch for function logic
				case "log":
					result = cmplx.Log(arg1)
				case "exp":
					result = cmplx.Exp(arg1)
				case "sin":
					result = cmplx.Sin(arg1)
				case "cos":
					result = cmplx.Cos(arg1)
				case "tan":
					result = cmplx.Tan(arg1)
				case "asin":
					result = cmplx.Asin(arg1)
				case "acos":
					result = cmplx.Acos(arg1)
				case "atan":
					result = cmplx.Atan(arg1)
				case "sinh":
					result = cmplx.Sinh(arg1)
				case "cosh":
					result = cmplx.Cosh(arg1)
				case "tanh":
					result = cmplx.Tanh(arg1)
				case "asinh":
					result = cmplx.Asinh(arg1)
				case "acosh":
					result = cmplx.Acosh(arg1)
				case "atanh":
					result = cmplx.Atanh(arg1)
				case "log10":
					result = cmplx.Log10(arg1)
				case "log2":
					result = cmplx.Log(arg1) / cmplx.Log(complex(2, 0))
				case "sqrt":
					result = cmplx.Sqrt(arg1)
				}
				operandStack = append(operandStack, result)
				processed = true
			} // End inner switch for function/constant names

			if !processed { // If IDENT was not a known constant or function
				return complex(math.NaN(), math.NaN()), NewCalculationError(
					fmt.Sprintf("unknown identifier '%s' encountered during evaluation at position %d", token.Literal, token.Position),
				)
			}

		case PLUS, MINUS, ASTERISK, SLASH, PERCENT, CARET, UNARY_MINUS: // Add UNARY_MINUS
			var op1, op2 complex128 // op1 is not used for unary
			var numOperandsNeeded int

			if token.Type == UNARY_MINUS {
				numOperandsNeeded = 1
			} else {
				numOperandsNeeded = 2
			}

			if len(operandStack) < numOperandsNeeded {
				return complex(math.NaN(), math.NaN()), NewCalculationError(
					fmt.Sprintf("insufficient operands for operator '%s' (type %s) at position %d", token.Literal, token.Type, token.Position),
				)
			}

			if numOperandsNeeded == 2 {
				op2 = operandStack[len(operandStack)-1]
				op1 = operandStack[len(operandStack)-2]
				operandStack = operandStack[:len(operandStack)-2]
			} else { // Unary
				op2 = operandStack[len(operandStack)-1] // Unary op acts on op2
				operandStack = operandStack[:len(operandStack)-1]
			}

			var result complex128
			var opErr error

			switch token.Type {
			case PLUS:
				result = op1 + op2
			case MINUS:
				result = op1 - op2
			case ASTERISK:
				result = op1 * op2
			case SLASH:
				result = op1 / op2
			case PERCENT:
				result, opErr = calculateModulo(op1, op2, token)
				if opErr != nil {
					return complex(math.NaN(), math.NaN()), opErr
				}
			case CARET:
				result = cmplx.Pow(op1, op2)
				//fmt.Printf("(%v)^(%v) = %v\n", op1, op2, result)
			case UNARY_MINUS:
				// op2 is the single operand for unary minus (e.g., the '4' in '-4')
				tempRes := -op2 // Perform the negation, e.g., -(4+0i) -> (-4-0i)

				// Normalize signed zeros in the result of this UNARY_MINUS operation.
				// This ensures that if the user types "-N" (N positive real),
				// it's treated as complex(-N, +0.0) for subsequent operations
				// like Pow, aligning with the standard branch cut convention for Log.
				r := real(tempRes)
				i := imag(tempRes)

				if r == 0.0 { // This normalizes -0.0 real to +0.0 real
					r = 0.0 // Assigning 0.0 defaults to +0.0
				}
				if i == 0.0 { // This normalizes -0.0 imag to +0.0 imag
					i = 0.0 // Assigning 0.0 defaults to +0.0
				}
				result = complex(r, i)
			}
			operandStack = append(operandStack, result)

		default:
			// This should not be reached if the RPN queue is well-formed by the parser
			// and contains only known token types for evaluation.
			return complex(math.NaN(), math.NaN()), NewCalculationError(
				fmt.Sprintf("unexpected token type '%s' in RPN queue (token: '%s' at pos %d)", token.Type, token.Literal, token.Position),
			)
		}
	}

	// After processing all tokens, the stack should contain exactly one result.
	if len(operandStack) == 1 {
		return operandStack[0], nil
	} else if len(operandStack) == 0 {
		// This could happen if the RPN queue was empty (e.g. empty input string,
		// though Parse should catch this) or an operator consumed all operands
		// without producing a result (which shouldn't happen with correct logic).
		return complex(math.NaN(), math.NaN()), NewCalculationError("invalid expression: no result on stack (empty RPN or malformed expression)")
	} else {
		// More than one value on the stack means the expression was malformed,
		// typically too many numbers or too few operators.
		return complex(math.NaN(), math.NaN()), NewCalculationError(
			fmt.Sprintf("invalid expression: %d values left on stack, expected 1 (check operators and operands)", len(operandStack)),
		)
	}
}
