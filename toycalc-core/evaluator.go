// evaluator.go
package toycalc_core

import (
	"fmt"
	"math"
	"math/cmplx"
	"strconv"
	"strings" // For ToLower on function names
)

var OutputFormatMode string = "auto" // "auto", "fixed", "sci" (as before)
var OutputDisplayPrecision int = 9   // Default number of decimal places to round to for display

// CalculateExpression orchestrates Lex, Parse, EvaluateRPN, and formatComplexOutput
func CalculateExpression(expressionString string) (string, error) {
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

// Helper to round a float64 to a specific number of decimal places for display
// This helps in making numbers like 0.89999999991 display as 0.9 if precision is, say, 8-10
/*func roundForDisplay(val float64, precision int) float64 {
	if math.IsNaN(val) || math.IsInf(val, 0) {
		return val // Don't round NaN or Inf
	}
	scale := math.Pow(10, float64(precision))
	return math.Round(val*scale) / scale
}*/

//const displayPrecision = 10 // Or a configurable value, for rounding before formatting

// Helper to check if a float64 (potentially after rounding for display) is effectively an integer
func isEffectivelyInteger(f float64, comparisonEpsilon float64) bool {
	if math.IsNaN(f) || math.IsInf(f, 0) {
		return false
	}
	// After f has been rounded to displayPrecision, check if it's an integer.
	// A small epsilon is still useful here to account for tiny residuals
	// if math.Round(f) itself isn't bit-perfectly the same as f when f is an integer.
	return math.Abs(f-math.Round(f)) < comparisonEpsilon
}

// Ternary helper for cleaner sign logic if needed elsewhere
func Ternary(condition bool, trueVal, falseVal string) string {
	if condition {
		return trueVal
	}
	return falseVal
}

// isEffectivelyZero checks if f (already rounded for display) is zero.
func isEffectivelyZero(f float64, comparisonEpsilon float64) bool {
	// After f has been rounded to displayPrecision, check if it's zero.
	return math.Abs(f) < comparisonEpsilon
}

func roundToDecimalPlaces(val float64, places int) float64 {
	if math.IsNaN(val) || math.IsInf(val, 0) {
		return val
	}
	if places < 0 {
		places = 0
	} // Or handle error
	scale := math.Pow(10, float64(places))
	return math.Round(val*scale) / scale
}

func formatComplexOutput(c complex128) string {
	realRaw := real(c)
	imagRaw := imag(c)

	// 1. Handle NaN/Inf (no rounding for these)
	if cmplx.IsNaN(c) {
		// ... (NaN formatting as before) ...
		return "NaN" // Simplified for brevity
	}
	if cmplx.IsInf(c) {
		return fmt.Sprintf("%v", c) // ... (Inf formatting as before) ...
	}

	// 2. Apply display rounding based on global/user setting outputDisplayPrecision
	realVal := roundToDecimalPlaces(realRaw, OutputDisplayPrecision)
	imagVal := roundToDecimalPlaces(imagRaw, OutputDisplayPrecision)

	// 3. Determine characteristics based on these "display-ready" values using Epsilon
	//    Epsilon here is for comparing these already-rounded numbers to perfect zero or integer.
	realIsZero := isEffectivelyZero(realVal, Epsilon)
	imagIsZero := isEffectivelyZero(imagVal, Epsilon)
	realIsInt := isEffectivelyInteger(realVal, Epsilon)
	imagIsInt := isEffectivelyInteger(imagVal, Epsilon)

	// 4. Format based on outputFormatMode (auto, fixed, sci) and outputDisplayPrecision
	var realStr /*imagStr,*/, imagSignStr, imagMagStr string

	// --- Format Real Part ---
	switch OutputFormatMode {
	case "fixed":
		realStr = fmt.Sprintf("%.*f", OutputDisplayPrecision, realVal)
	case "sci":
		realStr = fmt.Sprintf("%.*e", OutputDisplayPrecision, realVal)
	default: // "auto"
		if realIsZero && imagIsZero {
			return "0"
		} // Handle 0+0i early
		if realIsInt {
			realStr = fmt.Sprintf("%.0f", realVal)
		} else {
			realStr = fmt.Sprintf("%g", realVal)
		}
	}

	// --- Format Imaginary Part (if not negligible) ---
	if !imagIsZero {
		absImagVal := math.Abs(imagVal)
		isImagMagOne := isEffectivelyZero(absImagVal-1.0, Epsilon)

		switch OutputFormatMode {
		case "fixed":
			imagMagStr = fmt.Sprintf("%.*f", OutputDisplayPrecision, absImagVal)
		case "sci":
			imagMagStr = fmt.Sprintf("%.*e", OutputDisplayPrecision, absImagVal)
		default: // "auto"
			if isImagMagOne {
				imagMagStr = "" // For "i"
			} else if imagIsInt {
				imagMagStr = fmt.Sprintf("%.0f", absImagVal)
			} else {
				imagMagStr = fmt.Sprintf("%g", absImagVal)
			}
		}

		if imagVal < -Epsilon { // Check against -Epsilon for robust sign
			imagSignStr = "-"
		} else {
			imagSignStr = "+"
		}
	}

	// --- Assemble Output ---
	if imagIsZero {
		return realStr // Purely real (or 0+0i was handled)
	}

	if realIsZero { // Purely imaginary
		if isEffectivelyZero(math.Abs(imagVal)-1.0, Epsilon) { // Is magnitude 1?
			return Ternary(imagVal < 0, "-i", "i")
		}
		return fmt.Sprintf("%s%si", Ternary(imagVal < 0, "-", ""), imagMagStr)
	}

	// Full complex number
	if isEffectivelyZero(math.Abs(imagVal)-1.0, Epsilon) { // Is magnitude of imag 1?
		return fmt.Sprintf("%s %s i", realStr, imagSignStr)
	}
	return fmt.Sprintf("%s %s %si", realStr, imagSignStr, imagMagStr)
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
				"log10", "log2", "sqrt", "real", "imag", "abs", "phase",
				"conj", "degtorad", "radtodeg", "floor", "ceil", "round", "trunc":
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
				case "real":
					result = complex(real(arg1), 0.0)
				case "imag":
					result = complex(imag(arg1), 0.0)
				case "abs":
					result = complex(cmplx.Abs(arg1), 0.0)
				case "phase":
					result = complex(cmplx.Phase(arg1), 0.0)
				case "conj":
					result = cmplx.Conj(arg1)
				case "degtorad":
					result = arg1 * complex(math.Pi/180.0, 0.0)
				case "radtodeg":
					result = arg1 * complex(180.0/math.Pi, 0.0)
				case "floor":
					result = complex(math.Floor(real(arg1)), math.Floor(imag(arg1)))
				case "ceil":
					result = complex(math.Ceil(real(arg1)), math.Ceil(imag(arg1)))
				case "trunc":
					result = complex(math.Trunc(real(arg1)), math.Trunc(imag(arg1)))
				case "round":
					result = complex(math.Round(real(arg1)), math.Round(imag(arg1)))
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
