// toycalc_test.go
package main

import (
	"fmt"
	"math"
	"math/cmplx"
	"reflect"
	"strings"
	"testing"
)

// --- Test Helper Functions (ensure these are present and correct) ---

func checkError(t *testing.T, expectedErrorMsg string, actualError error) {
	t.Helper()
	if actualError == nil {
		if expectedErrorMsg != "" {
			t.Errorf("Expected error containing '%s', but got nil.", expectedErrorMsg)
		}
		return
	}
	if expectedErrorMsg == "" {
		t.Errorf("Expected no error, but got: '%s'", actualError.Error())
		return
	}
	if !strings.Contains(actualError.Error(), expectedErrorMsg) {
		t.Errorf("Expected error message to contain '%s', got '%s'", expectedErrorMsg, actualError.Error())
	}
}

func compareComplex(t *testing.T, expected, actual complex128, description string) {
	// ... (implementation from before, ensuring Epsilon is used)
	t.Helper()
	if cmplx.IsNaN(expected) {
		if !cmplx.IsNaN(actual) {
			t.Errorf("%s: expected NaN, got %v", description, actual)
		}
		return
	}
	if cmplx.IsNaN(actual) {
		t.Errorf("%s: expected %v, got NaN", description, expected)
		return
	}
	if cmplx.IsInf(expected) {
		expectedSignReal := math.Signbit(real(expected))
		expectedSignImag := math.Signbit(imag(expected))
		actualSignReal := math.Signbit(real(actual))
		actualSignImag := math.Signbit(imag(actual))

		isExpectedRealInf := math.IsInf(real(expected), 0)
		isExpectedImagInf := math.IsInf(imag(expected), 0)
		isActualRealInf := math.IsInf(real(actual), 0)
		isActualImagInf := math.IsInf(imag(actual), 0)

		if !(isExpectedRealInf == isActualRealInf && isExpectedImagInf == isActualImagInf &&
			(expectedSignReal == actualSignReal || !isExpectedRealInf) &&
			(expectedSignImag == actualSignImag || !isExpectedImagInf)) {
			t.Errorf("%s: expected Inf (%v), got %v", description, expected, actual)
		}
		return
	}
	if cmplx.IsInf(actual) {
		t.Errorf("%s: expected %v, got Inf (%v)", description, expected, actual)
		return
	}

	realDiff := math.Abs(real(expected) - real(actual))
	imagDiff := math.Abs(imag(expected) - imag(actual))

	if realDiff > Epsilon || imagDiff > Epsilon {
		t.Errorf("%s: expected %v (real: %g, imag: %g), got %v (real: %g, imag: %g). Diffs: real=%g, imag=%g",
			description,
			expected, real(expected), imag(expected),
			actual, real(actual), imag(actual),
			realDiff, imagDiff)
	}
}

func compareTokenSlices(t *testing.T, expected, actual []Token, description string) {
	t.Helper()
	if !reflect.DeepEqual(expected, actual) {
		// Enhanced error reporting for token slice mismatches
		msg := fmt.Sprintf("%s: token slices do not match.\nExpected (%d tokens):\n", description, len(expected))
		for i, tok := range expected {
			msg += fmt.Sprintf("  [%d] %+v\n", i, tok)
		}
		msg += fmt.Sprintf("Got (%d tokens):\n", len(actual))
		for i, tok := range actual {
			msg += fmt.Sprintf("  [%d] %+v\n", i, tok)
		}

		// Check for length difference first
		if len(expected) != len(actual) {
			msg += fmt.Sprintf("\nLength mismatch: expected %d, got %d", len(expected), len(actual))
			t.Errorf("%s", msg)
			return
		}
		// Find the first differing token
		for i := 0; i < len(expected); i++ {
			if !reflect.DeepEqual(expected[i], actual[i]) {
				msg += fmt.Sprintf("\nFirst difference at index %d:\nExpected: %+v\nGot:      %+v", i, expected[i], actual[i])
				break
			}
		}
		t.Errorf("%s", msg)
	}
}

// --- Initial Setup Tests ---
func TestFrameworkInitialSetup(t *testing.T) {
	t.Log("Test framework is running. Add more tests as features are implemented.")
	if 1+1 != 2 {
		t.Error("Something is fundamentally wrong if 1+1 is not 2.")
	}
}

// --- Lexer Tests (Stage 1) ---

func TestLexer(t *testing.T) {
	testCases := []struct {
		name           string
		input          string
		expectedTokens []Token
		expectedError  string // Substring of expected error message, or "" for no error
	}{
		// Basic Numbers
		{
			name:  "Integer",
			input: "123",
			expectedTokens: []Token{
				{Type: NUMBER, Literal: "123", Position: 0},
				{Type: EOF, Literal: "", Position: 3},
			},
		},
		{
			name:  "Float",
			input: "3.14159",
			expectedTokens: []Token{
				{Type: NUMBER, Literal: "3.14159", Position: 0},
				{Type: EOF, Literal: "", Position: 7},
			},
		},
		{
			name:  "Number scientific notation lowercase e",
			input: "1.23e-4",
			expectedTokens: []Token{
				{Type: NUMBER, Literal: "1.23e-4", Position: 0},
				{Type: EOF, Literal: "", Position: 7},
			},
		},
		{
			name:  "Number scientific notation uppercase E",
			input: "5E+10",
			expectedTokens: []Token{
				{Type: NUMBER, Literal: "5E+10", Position: 0},
				{Type: EOF, Literal: "", Position: 5},
			},
		},
		{
			name:  "Number ending with decimal point",
			input: "123.",
			expectedTokens: []Token{
				{Type: NUMBER, Literal: "123.", Position: 0},
				{Type: EOF, Literal: "", Position: 4},
			},
		},
		// Operators
		{
			name:  "All Single Char Operators",
			input: "+-*/%^",
			expectedTokens: []Token{
				{Type: PLUS, Literal: "+", Position: 0},
				{Type: MINUS, Literal: "-", Position: 1},
				{Type: ASTERISK, Literal: "*", Position: 2},
				{Type: SLASH, Literal: "/", Position: 3},
				{Type: PERCENT, Literal: "%", Position: 4},
				{Type: CARET, Literal: "^", Position: 5},
				{Type: EOF, Literal: "", Position: 6},
			},
		},
		// Delimiters
		{
			name:  "All Delimiters",
			input: "()[]{},",
			expectedTokens: []Token{
				{Type: LPAREN, Literal: "(", Position: 0},
				{Type: RPAREN, Literal: ")", Position: 1},
				{Type: LBRACKET, Literal: "[", Position: 2},
				{Type: RBRACKET, Literal: "]", Position: 3},
				{Type: LBRACE, Literal: "{", Position: 4},
				{Type: RBRACE, Literal: "}", Position: 5},
				{Type: COMMA, Literal: ",", Position: 6},
				{Type: EOF, Literal: "", Position: 7},
			},
		},
		// Identifiers (Function Names for Stage 1)
		{
			name:  "Log function",
			input: "log",
			expectedTokens: []Token{
				{Type: IDENT, Literal: "log", Position: 0},
				{Type: EOF, Literal: "", Position: 3},
			},
		},
		{
			name:  "Identifier with underscore and numbers",
			input: "var_1_test",
			expectedTokens: []Token{
				{Type: IDENT, Literal: "var_1_test", Position: 0},
				{Type: EOF, Literal: "", Position: 10},
			},
		},
		// Combined expressions
		{
			name:  "Simple addition with spaces",
			input: "1 + 2", // Positions: 0, 2, 4. EOF at 5.
			expectedTokens: []Token{
				{Type: NUMBER, Literal: "1", Position: 0},
				{Type: PLUS, Literal: "+", Position: 2},
				{Type: NUMBER, Literal: "2", Position: 4},
				{Type: EOF, Literal: "", Position: 5},
			},
		},
		{
			name:  "Expression with mixed operators and grouping",
			input: "log(10.5 * (varA - 3)) % arr[idx]^2",
			// Positions:
			// log:0, (:3, 10.5:4, *:9, (:11, varA:12, -:17, 3:19, ):20, ):21, %:23, arr:25, [:28, idx:29, ]:32, ^:33, 2:34, EOF:35
			expectedTokens: []Token{
				{Type: IDENT, Literal: "log", Position: 0}, {Type: LPAREN, Literal: "(", Position: 3},
				{Type: NUMBER, Literal: "10.5", Position: 4}, {Type: ASTERISK, Literal: "*", Position: 9}, {Type: LPAREN, Literal: "(", Position: 11},
				{Type: IDENT, Literal: "varA", Position: 12}, {Type: MINUS, Literal: "-", Position: 17}, {Type: NUMBER, Literal: "3", Position: 19}, {Type: RPAREN, Literal: ")", Position: 20},
				{Type: RPAREN, Literal: ")", Position: 21}, {Type: PERCENT, Literal: "%", Position: 23}, {Type: IDENT, Literal: "arr", Position: 25},
				{Type: LBRACKET, Literal: "[", Position: 28}, {Type: IDENT, Literal: "idx", Position: 29}, {Type: RBRACKET, Literal: "]", Position: 32},
				{Type: CARET, Literal: "^", Position: 33}, {Type: NUMBER, Literal: "2", Position: 34},
				{Type: EOF, Literal: "", Position: 35},
			},
		},
		// Whitespace handling
		{
			name:  "Whitespace abundant",
			input: "  1   \t + \n 2  ", // "1" at 2, "+" at 8, "2" at 12. Len 15. EOF at 15.
			expectedTokens: []Token{
				{Type: NUMBER, Literal: "1", Position: 2},
				{Type: PLUS, Literal: "+", Position: 8},
				{Type: NUMBER, Literal: "2", Position: 12},
				{Type: EOF, Literal: "", Position: 15}, // Position of EOF is after the last char consumed/skipped
			},
		},
		// Edge cases for numbers - lexer behavior might make these tricky for position if they are malformed
		{
			name:  "Number scientific notation missing exponent digits",
			input: "1.2e", // "1.2e" at 0. EOF at 4.
			expectedTokens: []Token{
				{Type: NUMBER, Literal: "1.2e", Position: 0}, // Lexer tokenizes this, strconv.ParseFloat will fail later
				{Type: EOF, Literal: "", Position: 4},
			},
		},
		{
			name:  "Number scientific notation with sign but no digits",
			input: "1.2e-", // "1.2e-" at 0. EOF at 5.
			expectedTokens: []Token{
				{Type: NUMBER, Literal: "1.2e-", Position: 0}, // Lexer tokenizes this
				{Type: EOF, Literal: "", Position: 5},
			},
		},
		// Illegal Characters
		{
			name:  "Illegal character in expression",
			input: "1 @ 2", // "1" at 0, "@" at 2 (ILLEGAL). EOF at 3.
			expectedTokens: []Token{
				{Type: NUMBER, Literal: "1", Position: 0},
				{Type: ILLEGAL, Literal: "@", Position: 2},
				// Lex() should return these two tokens and an error.
				// The EOF token might or might not be appended by Lex() before returning error.
				// The current Lex() appends ILLEGAL then breaks, so EOF isn't in the success path.
				// Let's assume the tokens slice returned with the error contains the ILLEGAL token.
			},
			expectedError: "illegal character '@' found at position 2",
		},
		{
			name:  "Illegal character at start",
			input: "$1 + 2", // "$" at 0 (ILLEGAL)
			expectedTokens: []Token{
				{Type: ILLEGAL, Literal: "$", Position: 0},
			},
			expectedError: "illegal character '$' found at position 0",
		},
		{
			name:  "Empty input",
			input: "",
			expectedTokens: []Token{
				{Type: EOF, Literal: "", Position: 0},
			},
		},
		{
			name:  "Only whitespace",
			input: "   \t\n   ", // Length 8. EOF at 8.
			expectedTokens: []Token{
				{Type: EOF, Literal: "", Position: 8},
			},
		},
		{
			name:  "Imaginary unit i",
			input: "i",
			expectedTokens: []Token{
				{Type: IDENT, Literal: "i", Position: 0},
				{Type: EOF, Literal: "", Position: 1},
			},
		},
		{
			name:  "Expression with i",
			input: "2*i + 3", // i at pos 2
			expectedTokens: []Token{
				{Type: NUMBER, Literal: "2", Position: 0},
				{Type: ASTERISK, Literal: "*", Position: 1},
				{Type: IDENT, Literal: "i", Position: 2},
				{Type: PLUS, Literal: "+", Position: 4},
				{Type: NUMBER, Literal: "3", Position: 6},
				{Type: EOF, Literal: "", Position: 7},
			},
		},
		{
			name:  "Identifier log10",
			input: "log10",
			expectedTokens: []Token{
				{Type: IDENT, Literal: "log10", Position: 0},
				{Type: EOF, Literal: "", Position: 5},
			},
		},
		{
			name:  "Identifier log2",
			input: "log2",
			expectedTokens: []Token{
				{Type: IDENT, Literal: "log2", Position: 0},
				{Type: EOF, Literal: "", Position: 4},
			},
		},
		{
			name:  "Identifier var2test",
			input: "var2test",
			expectedTokens: []Token{
				{Type: IDENT, Literal: "var2test", Position: 0},
				{Type: EOF, Literal: "", Position: 8},
			},
		},
		{
			name:  "Identifier with leading letter then numbers func123",
			input: "func123",
			expectedTokens: []Token{
				{Type: IDENT, Literal: "func123", Position: 0},
				{Type: EOF, Literal: "", Position: 7},
			},
		},
		{
			name:  "Identifier starting with underscore _val1",
			input: "_val1",
			expectedTokens: []Token{
				{Type: IDENT, Literal: "_val1", Position: 0},
				{Type: EOF, Literal: "", Position: 5},
			},
		},
		{
			name:  "Identifier mixed case with number Val10Test",
			input: "Val10Test",
			expectedTokens: []Token{
				{Type: IDENT, Literal: "Val10Test", Position: 0},
				{Type: EOF, Literal: "", Position: 9},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actualTokens, err := Lex(tc.input)

			checkError(t, tc.expectedError, err)

			if tc.expectedError != "" {
				// If an error is expected and occurs, we should also check if the partial token stream
				// matches what's expected up to (and including) the ILLEGAL token.
				// The Lex function, as currently designed, returns the tokens accumulated so far *including* the ILLEGAL one.
				compareTokenSlices(t, tc.expectedTokens, actualTokens, "Token stream up to error")
			} else {
				// No error expected, compare the full token slice.
				compareTokenSlices(t, tc.expectedTokens, actualTokens, "Token stream")
			}
		})
	}
}

// toycalc_test.go
// ... (add to existing file)

func TestParserToRPN(t *testing.T) {
	testCases := []struct {
		name          string
		input         string
		expectedRPN   []Token // We'll compare Type and Literal, ignore Position for RPN comparison simplicity
		expectedError string
	}{
		{
			name:  "Simple addition",
			input: "1 + 2",
			expectedRPN: []Token{
				{Type: NUMBER, Literal: "1", Position: 0},
				{Type: NUMBER, Literal: "2", Position: 4},
				{Type: PLUS, Literal: "+", Position: 2},
			},
		},
		{
			name:  "Precedence 1+2*3",
			input: "1 + 2 * 3", // Positions: 1(0), +(2), 2(4), *(6), 3(8)
			expectedRPN: []Token{
				{Type: NUMBER, Literal: "1", Position: 0},
				{Type: NUMBER, Literal: "2", Position: 4},
				{Type: NUMBER, Literal: "3", Position: 8},
				{Type: ASTERISK, Literal: "*", Position: 6},
				{Type: PLUS, Literal: "+", Position: 2},
			},
		},
		{
			name:  "Parentheses (1+2)*3",
			input: "(1 + 2) * 3", // (0, 1(1), +(3), 2(5), )(6), *(8), 3(10)
			expectedRPN: []Token{
				{Type: NUMBER, Literal: "1", Position: 1},
				{Type: NUMBER, Literal: "2", Position: 5},
				{Type: PLUS, Literal: "+", Position: 3},
				{Type: NUMBER, Literal: "3", Position: 10},
				{Type: ASTERISK, Literal: "*", Position: 8},
			},
		},
		{
			name:  "Power right associative 2^3^2",
			input: "2^3^2", // 2(0), ^(1), 3(2), ^(3), 2(4)
			expectedRPN: []Token{
				{Type: NUMBER, Literal: "2", Position: 0},
				{Type: NUMBER, Literal: "3", Position: 2},
				{Type: NUMBER, Literal: "2", Position: 4},
				{Type: CARET, Literal: "^", Position: 3}, // Inner ^
				{Type: CARET, Literal: "^", Position: 1}, // Outer ^
			},
		},
		{
			name:  "Unary minus -5",
			input: "-5", // -(0), 5(1)
			expectedRPN: []Token{
				{Type: NUMBER, Literal: "5", Position: 1},
				{Type: UNARY_MINUS, Literal: "-", Position: 0},
			},
		},
		{
			name:  "Unary plus +5 (ignored by parser, only operand remains)",
			input: "+5", // +(0), 5(1)
			expectedRPN: []Token{
				{Type: NUMBER, Literal: "5", Position: 1},
			},
		},
		{
			name:  "Function call log(10)",
			input: "log(10)", // log(0), ((3), 10(4), )(6)
			expectedRPN: []Token{
				{Type: NUMBER, Literal: "10", Position: 4},
				{Type: IDENT, Literal: "log", Position: 0},
			},
		},
		{
			name:  "Function call with expression exp(1+2)",
			input: "exp(1+2)", // exp(0), ((3), 1(4), +(5), 2(6), )(7)
			expectedRPN: []Token{
				{Type: NUMBER, Literal: "1", Position: 4},
				{Type: NUMBER, Literal: "2", Position: 6},
				{Type: PLUS, Literal: "+", Position: 5},
				{Type: IDENT, Literal: "exp", Position: 0},
			},
		},

		// --- Test Cases for Implied Multiplication ---
		{
			name:  "Implied mult num-paren: 2(3+4)",
			input: "2(3+4)", // 2(0), ((1), 3(2), +(3), 4(4), )(5) -> implicit * at pos 1
			expectedRPN: []Token{
				{Type: NUMBER, Literal: "2", Position: 0},
				{Type: NUMBER, Literal: "3", Position: 2},
				{Type: NUMBER, Literal: "4", Position: 4},
				{Type: PLUS, Literal: "+", Position: 3},
				{Type: ASTERISK, Literal: "*", Position: 1}, // Implicit * takes position of '('
			},
		},
		{
			name:  "Implied mult paren-paren: (1)(2)",
			input: "(1)(2)", // (0,1(1),)(2), ((3),2(4),)(5) -> implicit * at pos 3
			expectedRPN: []Token{
				{Type: NUMBER, Literal: "1", Position: 1},
				{Type: NUMBER, Literal: "2", Position: 4},
				{Type: ASTERISK, Literal: "*", Position: 3}, // Implicit * takes position of second '('
			},
		},
		{
			name:  "Implied mult num-ident(const): 3i",
			input: "3i", // 3(0), i(1) -> implicit * at pos 1
			expectedRPN: []Token{
				{Type: NUMBER, Literal: "3", Position: 0},
				{Type: IDENT, Literal: "i", Position: 1},    // 'i' is an operand
				{Type: ASTERISK, Literal: "*", Position: 1}, // Implicit * takes position of 'i'
			},
		},
		{
			name:  "Implied mult num-ident(func): 3log(10)",
			input: "3log(10)", // 3(0), log(1), ((4), 10(5), )(7) -> implicit * at pos 1
			expectedRPN: []Token{
				{Type: NUMBER, Literal: "3", Position: 0},
				{Type: NUMBER, Literal: "10", Position: 5},
				{Type: IDENT, Literal: "log", Position: 1},
				{Type: ASTERISK, Literal: "*", Position: 1}, // Implicit * takes position of 'log'
			},
		},
		{
			name:  "User case: sin(1)cos(2)",
			input: "sin(1)cos(2)", // sin(0),(,1,),cos(6),(,2,) -> implicit * at pos 6 (start of 'cos')
			// RPN: 1 sin 2 cos *
			expectedRPN: []Token{
				{Type: NUMBER, Literal: "1", Position: 4}, {Type: IDENT, Literal: "sin", Position: 0},
				{Type: NUMBER, Literal: "2", Position: 10}, {Type: IDENT, Literal: "cos", Position: 6},
				{Type: ASTERISK, Literal: "*", Position: 6},
			},
		},
		{
			name:  "User case: i pi",
			input: "i pi", // i(0), pi(2) -> implicit * at pos 2
			// RPN: i pi *
			expectedRPN: []Token{
				{Type: IDENT, Literal: "i", Position: 0}, {Type: IDENT, Literal: "pi", Position: 2},
				{Type: ASTERISK, Literal: "*", Position: 2},
			},
		},
		{
			name:  "No implied mult for binary minus: (2+3) - 5",
			input: "(2+3) - 5", // -(5) is binary minus
			// RPN: 2 3 + 5 -
			expectedRPN: []Token{
				{Type: NUMBER, Literal: "2", Position: 1}, {Type: NUMBER, Literal: "3", Position: 3}, {Type: PLUS, Literal: "+", Position: 2},
				{Type: NUMBER, Literal: "5", Position: 9}, {Type: MINUS, Literal: "-", Position: 7},
			},
		},
		{
			name:          "Mismatched parenthesis (open)",
			input:         "(1+2", // ( at 0
			expectedError: "mismatched parentheses/brackets/braces",
		},
		/*{
			name:          "Syntax error from parser for adjacent numbers",
			input:         "2 3 + 4", // 3 at pos 2
			expectedError: "unexpected number '3' at position 2; an operator may be missing",
		},*/
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tokens, lexErr := Lex(tc.input)
			if lexErr != nil {
				// If lexing fails, and an error was not primarily expected at parser stage for this input
				if tc.expectedError == "" || !strings.Contains(lexErr.Error(), tc.expectedError) {
					t.Fatalf("Lexing failed unexpectedly: %v", lexErr)
				}
				// If lexer error was the one expected, then test passes for error stage.
				if tc.expectedError != "" && strings.Contains(lexErr.Error(), tc.expectedError) {
					return
				}
			}

			// Create a simplified version of expectedRPN for comparison if Positions are not checked in RPN
			simpleExpectedRPN := make([]Token, len(tc.expectedRPN))
			for i, tok := range tc.expectedRPN {
				simpleExpectedRPN[i] = Token{Type: tok.Type, Literal: tok.Literal} // Position ignored for RPN output check
			}

			actualRPN, parseErr := Parse(tokens)
			checkError(t, tc.expectedError, parseErr)

			if tc.expectedError == "" && parseErr == nil {
				// Create a simplified version of actualRPN for comparison
				simpleActualRPN := make([]Token, len(actualRPN))
				for i, tok := range actualRPN {
					simpleActualRPN[i] = Token{Type: tok.Type, Literal: tok.Literal}
				}
				compareTokenSlices(t, simpleExpectedRPN, simpleActualRPN, "RPN output for "+tc.input)
			}
		})
	}
}

// toycalc_test.go
// ... (add to existing file)

func TestEvaluateRPN(t *testing.T) {
	testCases := []struct {
		name           string
		input          string
		expectedResult complex128
		expectedError  string // Substring of error message
	}{
		// Basic Arithmetic
		{"add", "1 + 2", complex(3, 0), ""},
		{"subtract", "5 - 1.5", complex(3.5, 0), ""},
		{"multiply", "2.5 * 4", complex(10, 0), ""},
		{"divide", "10 / 4", complex(2.5, 0), ""},
		{"divide by zero", "1 / 0", complex(math.Inf(1), math.NaN()), ""}, // Go's cmplx default for 1/(0+0i) is (+Inf+NaNi)
		{"zero by zero", "0 / 0", complex(math.NaN(), math.NaN()), ""},    // (NaN+NaNi)

		// Precedence & Parentheses (implicitly tested via full pipeline)
		{"precedence", "2 + 3 * 4", complex(14, 0), ""},
		{"parentheses", "(2 + 3) * 4", complex(20, 0), ""},
		{"all grouping", "{[ (1) + 2 ] - 3} / 4", complex(0, 0), ""}, // (1+2-3)/4 = 0/4 = 0

		// Power
		{"power simple", "2^3", complex(8, 0), ""},
		{"power fractional", "16^0.5", complex(4, 0), ""},
		{"power negative base fractional exp", "(-4)^0.5", complex(0, 2), ""}, // sqrt(-4) = 2i
		{"power complex", "(1+1*i)^2", complex(0, 2), ""},                     // (1+i)^2 = 1 + 2i + i^2 = 1 + 2i - 1 = 2i

		// Modulo
		{"modulo simple", "10 % 3", complex(1, 0), ""},                                                                          // 10 = 3*3 + 1
		{"modulo negative", "10 % -3", complex(1, 0), ""},                                                                       // 10 = (-3)*(-3) + 1. (10/-3 = -3.33, round to -3. -3*-3=9. 10-9=1)
		{"modulo result negative", "-10 % 3", complex(-1, 0), ""},                                                               // -10 = (-3)*3 - 1
		{"modulo float", "10.5 % 3.2", complex(10.5-3*3.2, 0), ""},                                                              // 10.5/3.2 = 3.28_ ; round(3.28) = 3. 10.5 - 3*3.2 = 10.5 - 9.6 = 0.9
		{"modulo by zero error", "5 % 0", complex(math.NaN(), math.NaN()), "divisor is zero for modulo operator at position 2"}, // Error from calculateModulo

		// Functions
		{"log of 1", "log(1)", complex(0, 0), ""},
		{"log of e", fmt.Sprintf("log(%.20f)", math.E), complex(1, 0), ""}, // log(e)
		{"log of -1", "log(-1)", complex(0, math.Pi), ""},                  // ln(-1) = i*pi
		{"exp of 0", "exp(0)", complex(1, 0), ""},
		{"exp of 1", "exp(1)", complex(math.E, 0), ""},
		{"log(exp(2.5))", "log(exp(2.5))", complex(2.5, 0), ""},
		{"log of zero", "log(0)", complex(math.Inf(-1), math.NaN()), ""}, // cmplx.Log(0) behavior

		// Error cases
		{"insufficient ops for plus", "1 +", 0, "insufficient operands for operator '+'"}, // Parser might catch this first
		{"unknown function", "unknown(5)", 0, "unknown identifier or function 'unknown'"},
		// imaginary unit
		{"imaginary unit", "i", complex(0, 1), ""},
		{"simple complex", "2+3*i", complex(2, 3), ""},
		{"i squared", "i^2", complex(-1, 0), ""},
		{"log of i", "log(i)", cmplx.Log(complex(0, 1)), ""},                      // log(i) = i*pi/2
		{"exp of i*pi", fmt.Sprintf("exp(i*%.20f)", math.Pi), complex(-1, 0), ""}, // e^(i*pi) = -1 (Euler's identity)
		{"-i", "-i", complex(0, -1), ""},                                          // Test with unary minus

		// Complex example
		{"complex full", "(1+2*i)/(1-1*i) + log(-1)^2", cmplx.Log(complex(-1, 0))*cmplx.Log(complex(-1, 0)) + (complex(1, 2) / complex(1, -1)), ""},
		// (1+2i)/(1-1i) = (1+2i)(1+i) / ((1-i)(1+i)) = (1+i+2i-2)/2 = (-1+3i)/2 = -0.5 + 1.5i
		// log(-1) = pi*i
		// (pi*i)^2 = pi^2 * i^2 = -pi^2
		// Result: -0.5 + 1.5i - pi^2
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// We use calculateExpression for an end-to-end test for the evaluator
			_, err := calculateExpression(tc.input)

			checkError(t, tc.expectedError, err)

			if tc.expectedError == "" && err == nil {
				// Need to parse actualStr back to complex128 if we want to use compareComplex.
				// Or, pre-calculate the expected string output.
				// For now, let's compare the complex numbers directly.
				// This requires calculateExpression to return complex128 for testing,
				// or we have a separate test for EvaluateRPN directly.
				// Let's assume we are testing the full pipeline and comparing the final result.

				// To properly test, we need the expected complex result, not just string.
				// The test cases above have `expectedResult complex128`.
				// So, we need to get the complex result from calculateExpression,
				// which means we need to call EvaluateRPN directly after Lex and Parse for this test.

				tokens, lexErr := Lex(tc.input)
				if lexErr != nil {
					t.Fatalf("Lexing failed for input '%s': %v", tc.input, lexErr)
				}
				rpn, parseErr := Parse(tokens)
				if parseErr != nil {
					t.Fatalf("Parsing failed for input '%s': %v", tc.input, parseErr)
				}
				actualResult, evalErr := EvaluateRPN(rpn)

				// Re-check error, this time from EvaluateRPN directly
				checkError(t, tc.expectedError, evalErr)
				if tc.expectedError == "" && evalErr == nil {
					compareComplex(t, tc.expectedResult, actualResult, "Result for '"+tc.input+"'")
				}
			} else if tc.expectedError != "" && err != nil {
				// Error was expected and received, checkError already did the message comparison.
			} else if tc.expectedError != "" && err == nil {
				t.Errorf("Expected error '%s' but got no error for input '%s'", tc.expectedError, tc.input)
			} else if tc.expectedError == "" && err != nil {
				t.Errorf("Expected no error but got '%s' for input '%s'", err.Error(), tc.input)
			}
		})
	}
}

// --- End-to-End Tests (Stage 1 - testing calculateExpression) ---
func TestCalculateExpressionEndToEnd(t *testing.T) {
	// For comparing float strings, use a helper or be mindful of precision.
	// fmt.Sprintf("%g", ...) is used by formatComplexOutput, so we can use it for expected values.
	//piStr := fmt.Sprintf("%g", math.Pi) // For use in expected strings if needed
	eStr := fmt.Sprintf("%.10g", math.E) // For use in expected strings if needed

	testCases := []struct {
		name                   string
		input                  string
		expectedOutput         string // Expected string from formatComplexOutput
		expectedErrorSubstring string // Substring of error message, or "" for no error
	}{
		// Basic Arithmetic
		{"E2E Add", "1 + 2", "3", ""},
		{"E2E Subtract", "5 - 1.5", "3.5", ""},
		{"E2E Multiply", "2.5 * 4", "10", ""},
		{"E2E Divide", "10 / 4", "2.5", ""},
		{"E2E Divide leading to float", "1 / 2", "0.5", ""},
		{"E2E Multiple ops", "10 - 2 * 3 + 8 / 4", "6", ""}, // 10 - 6 + 2 = 6

		// Unary Operators
		{"E2E Unary Minus Start", "-5", "-5", ""},
		{"E2E Unary Plus Start", "+5", "5", ""},
		{"E2E Unary Minus After Operator", "10 * -2", "-20", ""},
		{"E2E Unary Plus After Operator", "10 * +2", "20", ""},
		{"E2E Unary Minus Parenthesized", "(-3)", "-3", ""},
		{"E2E Unary Plus Parenthesized", "(+3)", "3", ""},
		{"E2E Double Unary Minus", "--3", "3", ""},
		{"E2E Unary Plus Minus", "+-3", "-3", ""},
		{"E2E Unary Minus Plus", "-+3", "-3", ""}, // - (+3)
		//{"E2E Unary Minus Complex", "-(1+2i)", "-1-2i", ""},

		// Power
		{"E2E Power Simple", "2^3", "8", ""},
		{"E2E Power Fractional", "16^0.5", "4", ""},
		{"E2E Power Negative Base Sqrt (Principal Value)", "(-4)^0.5", "2i", ""}, // Due to UNARY_MINUS normalization
		{"E2E Power with Unary Minus Base", "-2^4", "-16", ""},                   // -(2^4) because ^ has higher precedence than UNARY_MINUS
		{"E2E Power with Parenthesized Unary Minus Base", "(-2)^4", "16", ""},
		{"E2E Power Complex Base", "(1+i)^2", "2i", ""},
		//{"E2E Power Complex Exponent", "(1+i)^(1-i)", "2.807879+1.317865i", ""}, // Approximate, from WolframAlpha for (1+i)^(1-i)

		// Modulo (Gaussian Remainder)
		{"E2E Modulo Simple", "10 % 3", "1", ""},
		{"E2E Modulo Negative Divisor", "10 % -3", "1", ""},   // 10 = (-3)*(-3) + 1
		{"E2E Modulo Negative Dividend", "-10 % 3", "-1", ""}, // -10 = 3*(-3) - 1
		//{"E2E Modulo Float", "10.5 % 3.2", "0.9", ""},           // 10.5/3.2 = 3.28... round(3.28)=3. 10.5 - 3*3.2 = 0.9
		{"E2E Modulo Complex 1", "(5+3*i) % (2+1*i)", "-1", ""}, // (5+3i)/(2+1i) = (13/5 + 1/5 i) = 2.6+0.2i. Round: 3+0i. (5+3i) - (3)(2+1i) = 5+3i - (6+3i) = -1+0i.
		// Let's recheck (5+3i)/(2+1i) = ((5+3i)(2-i)) / ((2+i)(2-i)) = (10-5i+6i+3)/5 = (13+i)/5 = 2.6+0.2i.
		// Rounded quotient x = 3+0i.
		// Remainder r = (5+3i) - (3+0i)*(2+1i) = (5+3i) - (6+3i) = -1.
		// Expected: -1
		{"E2E Modulo Complex 2", "(1+i) % (1-i)", "0", ""}, // (1+i)/(1-i) = i. Rounded quotient = i. (1+i) - i*(1-i) = 1+i - (i+1) = 0.

		// Functions
		{"E2E Log of 1", "log(1)", "0", ""},
		{"E2E Log of E", fmt.Sprintf("log(%v)", math.E), "1", ""},
		{"E2E Log of -1 (Principal)", "log(-1)", fmt.Sprintf("%.10gi", math.Pi), ""},
		{"E2E Exp of 0", "exp(0)", "1", ""},
		{"E2E Exp of 1", "exp(1)", eStr, ""},
		{"E2E Log(Exp(2.5))", "log(exp(2.5))", "2.5", ""},
		{"E2E Exp(Log(1+i))", "exp(log(1+i))", "1+1i", ""},

		// Constant i
		{"E2E Constant i", "i", "1i", ""},
		{"E2E 2*i", "2*i", "2i", ""},
		{"E2E i*i", "i*i", "-1", ""},
		{"E2E i^2", "i^2", "-1", ""},
		{"E2E -i", "-i", "-1i", ""},
		{"constant pi", "pi", fmt.Sprintf("%.10g", math.Pi), ""},
		{"constant e", "e", fmt.Sprintf("%.10g", math.E), ""},
		{"expr with pi", "2*pi", fmt.Sprintf("%.10g", 2*math.Pi), ""},
		{"expr with e", "e^2", fmt.Sprintf("%.10g", math.E*math.E), ""},
		{"Euler's identity", "exp(i*pi)", "-1", ""}, // This is a classic!
		{"log of e", "log(e)", "1", ""},
		{"sin of pi/2", "sin(pi/2)", "1", ""},

		// Update previous examples if they used hardcoded Pi or E values
		// Example: from Stage 1 log tests
		// {"E2E Log of E", fmt.Sprintf("log(%s)", eStr), "1", ""}, // Old
		{"E2E Log of E new", "log(e)", "1", ""}, // New

		// Example: from Stage 1 exp tests
		// {"E2E Exp of 1", "exp(1)", eStr, ""}, // Old
		{"E2E Exp of 1 new", "exp(1)", fmt.Sprintf("%.10g", math.E), ""}, // Still useful to test exp(1)

		// Example: from Stage 2 trig tests
		// {"sin(pi/2)", fmt.Sprintf("sin(%g)", math.Pi/2), "1", ""}, // Old
		{"sin(pi/2) new", "sin(pi/2)", "1", ""}, // New

		// Example: from Stage 2 hyperbolic tests
		// {"cosh(i*pi)", fmt.Sprintf("cosh(i*%g)", math.Pi), "-1", ""}, // Old
		{"cosh(i*pi) new", "cosh(i*pi)", "-1", ""}, // New

		// Grouping Symbols
		{"E2E Parentheses", "(1+2)*3", "9", ""},
		{"E2E Brackets", "[1+2]*3", "9", ""},
		{"E2E Braces", "{1+2}*3", "9", ""},
		{"E2E Mixed Grouping", "{[ (10.5) ] / (1 + 2) } - 0.5", "3", ""}, // (10.5 / 3) - 0.5 = 3.5 - 0.5 = 3

		// Error Conditions
		{"E2E Mismatched Paren Open", "(1+2", "", "mismatched parentheses"},
		{"E2E Mismatched Paren Close", "1+2)", "", "mismatched parentheses"},
		{"E2E Unexpected Operator", "*2+3", "", "unexpected operator '*'"}, // Parser should catch this with expectOperand
		//{"E2E Missing Operator", "2 3 + 4", "", "unexpected number '3'"},                        // Parser should catch this
		{"E2E Div by Zero Error", "1/0", "(+Inf+NaNi)", ""},                                     // cmplx default, formatComplexOutput uses fmt.Sprintf %v
		{"E2E Mod by Zero Error", "5%0", "", "divisor is zero for modulo operator at position"}, // Error from calculateModulo
		{"E2E Unknown Function", "unknown(5)", "", "unknown identifier or function 'unknown'"},
		{"E2E Insufficient Ops for Plus", "1+", "", "insufficient operands for operator '+'"},        // This might be caught by parser due to EOF
		{"E2E Insufficient Ops for Func", "log()", "", "missing operand before closing parenthesis"}, // Parser error

		// More complex expressions
		{"E2E Complex 1", "-(2+3*i) * (1-i) + exp(log(5)) / ( (1+i)^2 % (1+0.5*i) )", "-5+9i", ""}, // This one needs careful manual calculation for expectedOutput
		// -(2+3i)*(1-i) = - (2-2i+3i+3) = -(5+i) = -5-i
		// exp(log(5)) = 5
		// (1+i)^2 = 2i
		// (1+0.5i)
		// 2i / (1+0.5i) = 2i(1-0.5i) / (1+0.25) = (2i+1)/1.25 = (1+2i)/1.25 = 0.8+1.6i. Round: 1+2i
		// 2i % (1+0.5i) -> r = 2i - (1+2i)*(1+0.5i) = 2i - (1+0.5i+2i-1) = 2i - 2.5i = -0.5i
		// Result: (-5-i) + 5 / (-0.5i) = -5-i + (5 * 2i) / (-0.5i * 2i) = -5-i + 10i / 1 = -5+9i
	}
	// Add the expected output for the complex one above
	//testCases[len(testCases)-1].expectedOutput = "-5+9i"

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actualOutput, err := calculateExpression(tc.input)
			checkError(t, tc.expectedErrorSubstring, err)

			if tc.expectedErrorSubstring == "" && err == nil {
				if actualOutput != tc.expectedOutput {
					// For complex float comparisons, direct string might be tricky due to precision.
					// However, %g format helps a lot. If this fails, it might be a genuine difference
					// or a very minor precision formatting difference not handled by %g.
					t.Errorf("Input '%s': Expected output '%s', but got '%s'",
						tc.input, tc.expectedOutput, actualOutput)
				}
			}
			// If an error was expected and received, checkError already verified it.
			// If an error was expected but not received, checkError handles it.
			// If an error was not expected but received, checkError handles it.
		})
	}
}

// --- Output Formatting Tests (Stage 1) ---
func TestFormatComplexOutput(t *testing.T) {
	tests := []struct {
		name     string
		input    complex128
		expected string
	}{
		{"Real number", complex(5, 0), "5"},
		{"Real number with tiny imag", complex(5, Epsilon/10), "5"},
		{"Real number with just above Epsilon imag", complex(5, Epsilon*2), fmt.Sprintf("%g%+gi", 5.0, Epsilon*2)},
		{"Negative real number", complex(-5, 0), "-5"},
		{"Pure imaginary", complex(0, 5), "5i"}, // %g for 0 is "0"
		{"Pure negative imaginary", complex(0, -5), "-5i"},
		{"Complex number", complex(3, -2), "3-2i"},
		{"Complex number positive imag", complex(3, 2), "3+2i"},
		{"Zero", complex(0, 0), "0"},
		{"NaN real", complex(math.NaN(), 0), "NaN"},
		{"NaN imag", complex(0, math.NaN()), "NaN"},
		{"NaN both", complex(math.NaN(), math.NaN()), "NaN"},
		{"Inf real", complex(math.Inf(1), 0), "(+Inf+0i)"}, // Based on current formatComplexOutput
		{"Inf imag", complex(0, math.Inf(-1)), "(0-Infi)"}, // Based on current formatComplexOutput
		{"Complex Inf", complex(math.Inf(1), math.Inf(-1)), "(+Inf-Infi)"},
		{"Small real part", complex(1.23e-10, 5), "1.23e-10+5i"},
		{"Small imag part (not negligible)", complex(5, 1.23e-10), "5+1.23e-10i"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := formatComplexOutput(tt.input); got != tt.expected {
				t.Errorf("formatComplexOutput(%v) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

// This struct can be used for tests that specifically check the complex128 result.
type rpnEvalTestCase struct {
	name           string
	input          string     // Input expression string
	expectedResult complex128 // Expected complex128 result
	expectedError  string     // Substring of error message, or "" if no error
	skipStrCompare bool       // Set to true if string output is too hard to predict exactly for this test
}

// runRPNEvalTests is a helper to run a table of rpnEvalTestCase
func runRPNEvalTests(t *testing.T, testCases []rpnEvalTestCase) {
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tokens, lexErr := Lex(tc.input)
			if lexErr != nil {
				// If lexing itself fails, check if an error was expected at this stage
				if tc.expectedError != "" && strings.Contains(lexErr.Error(), tc.expectedError) {
					return // Lexer error was expected and matched
				}
				t.Fatalf("Lexing failed unexpectedly for input '%s': %v", tc.input, lexErr)
			}

			rpn, parseErr := Parse(tokens)
			if parseErr != nil {
				// If parsing fails, check if an error was expected at this stage
				if tc.expectedError != "" && strings.Contains(parseErr.Error(), tc.expectedError) {
					return // Parser error was expected and matched
				}
				t.Fatalf("Parsing failed unexpectedly for input '%s': %v", tc.input, parseErr)
			}

			actualResult, evalErr := EvaluateRPN(rpn)
			checkError(t, tc.expectedError, evalErr) // Check for evaluation errors

			if tc.expectedError == "" && evalErr == nil {
				compareComplex(t, tc.expectedResult, actualResult, "Result for '"+tc.input+"'")
			}
			// The following can be added if you also want to spot-check string output for these,
			// but the primary focus is the complex128 result.
			// if !tc.skipStrCompare && tc.expectedError == "" && evalErr == nil {
			//     actualStrOutput := formatComplexOutput(actualResult)
			//     expectedStrOutput := formatComplexOutput(tc.expectedResult) // Format expected for consistency
			//     if actualStrOutput != expectedStrOutput {
			//         t.Errorf("String output for '%s': Expected '%s', got '%s'", tc.input, expectedStrOutput, actualStrOutput)
			//     }
			// }
		})
	}
}

func TestStage2FunctionsEvaluator(t *testing.T) {
	/*z1 := complex(3, 4) // |z1|=5, phase(z1) approx 0.927
	z2 := complex(-1, 2)
	zReal := complex(5.7, 0)
	zImag := complex(0, -2.5)
	zZero := complex(0, 0)
	zNegReal := complex(-9.0, 0.0)    */ // will become (-9+0i) after unary minus normalization
	//zNegRealWithNegZeroImag := complex(-9.0, math.Copysign(0.0, -1.0)) // For specific tests if needed

	testCases := []rpnEvalTestCase{
		// Core Complex Operations
		{name: "real(3+4i)", input: "real(3+4i)", expectedResult: complex(3, 0), expectedError: ""},
		{name: "imag(3+4i)", input: "imag(3+4i)", expectedResult: complex(4, 0), expectedError: ""},
		{name: "abs(3+4i)", input: "abs(3+4i)", expectedResult: complex(5, 0), expectedError: ""},
		{name: "abs(-5)", input: "abs(-5)", expectedResult: complex(5, 0), expectedError: ""},
		{name: "phase(1+i)", input: "phase(1+i)", expectedResult: complex(math.Pi/4, 0), expectedError: ""},
		{name: "phase(1)", input: "phase(1)", expectedResult: complex(0, 0), expectedError: ""},
		{name: "phase(-1)", input: "phase(-1)", expectedResult: complex(math.Pi, 0), expectedError: ""}, // After unary minus normalization for -1
		{name: "phase(i)", input: "phase(i)", expectedResult: complex(math.Pi/2, 0), expectedError: ""},
		{name: "phase(-i)", input: "phase(-i)", expectedResult: complex(-math.Pi/2, 0), expectedError: ""},
		{name: "phase(0)", input: "phase(0)", expectedResult: complex(0, 0), expectedError: ""}, // cmplx.Phase(0) is 0
		{name: "conj(3+4i)", input: "conj(3+4i)", expectedResult: complex(3, -4), expectedError: ""},
		{name: "conj(5)", input: "conj(5)", expectedResult: complex(5, 0), expectedError: ""},
		{name: "conj(2*i)", input: "conj(2*i)", expectedResult: complex(0, -2), expectedError: ""},

		// Logarithmic & Sqrt (some might overlap with existing E2E tests but good for specific complex check)
		{name: "log10(100+0i)", input: "log10(100)", expectedResult: complex(2, 0), expectedError: ""},
		{name: "log10(i)", input: "log10(i)", expectedResult: cmplx.Log10(complex(0, 1)), expectedError: ""},
		{name: "log2(8+0i)", input: "log2(8)", expectedResult: complex(3, 0), expectedError: ""},
		{name: "log2(i)", input: "log2(i)", expectedResult: cmplx.Log(complex(0, 1)) / cmplx.Log(complex(2, 0)), expectedError: ""},
		{name: "sqrt(2i)", input: "sqrt(2*i)", expectedResult: complex(1, 1), expectedError: ""},
		{name: "sqrt(-9)", input: "sqrt(-9)", expectedResult: complex(0, 3), expectedError: ""},                                             // After unary minus fix for -9
		{name: "sqrt(complex(-9, -0.0))", input: "sqrt(-9-0i)", expectedResult: complex(0.0, 3.0), expectedError: "", skipStrCompare: true}, // specific test if needed

		// Trigonometric
		{name: "sin(pi/2)", input: "sin(pi/2)", expectedResult: complex(1, 0), expectedError: ""},
		{name: "sin(i)", input: "sin(i)", expectedResult: complex(0, math.Sinh(1)), expectedError: ""},
		{name: "cos(pi)", input: "cos(pi)", expectedResult: complex(-1, 0), expectedError: ""},
		{name: "cos(i)", input: "cos(i)", expectedResult: complex(math.Cosh(1), 0), expectedError: ""},
		{name: "tan(pi/4)", input: "tan(pi/4)", expectedResult: complex(1, 0), expectedError: ""},
		{name: "tan(i)", input: "tan(i)", expectedResult: complex(0, math.Tanh(1)), expectedError: ""},
		{name: "tan(pi/2)", input: "tan(pi/2)", expectedResult: cmplx.Tan(complex(math.Pi/2, 0)), expectedError: "", skipStrCompare: true}, // Result is large or Inf

		// Inverse Trigonometric
		{name: "asin(0.5)", input: "asin(0.5)", expectedResult: cmplx.Asin(complex(0.5, 0)), expectedError: ""},
		{name: "asin(2)", input: "asin(2)", expectedResult: cmplx.Asin(complex(2, 0)), expectedError: ""}, // Complex result
		{name: "acos(0.5)", input: "acos(0.5)", expectedResult: cmplx.Acos(complex(0.5, 0)), expectedError: ""},
		{name: "acos(2)", input: "acos(2)", expectedResult: cmplx.Acos(complex(2, 0)), expectedError: ""}, // Complex result
		{name: "atan(1)", input: "atan(1)", expectedResult: cmplx.Atan(complex(1, 0)), expectedError: ""},
		{name: "atan(2*i)", input: "atan(2*i)", expectedResult: cmplx.Atan(complex(0, 2)), expectedError: ""},

		// Hyperbolic
		{name: "sinh(1)", input: "sinh(1)", expectedResult: cmplx.Sinh(complex(1, 0)), expectedError: ""},
		{name: "sinh(i*pi/2)", input: "sinh(i*pi/2)", expectedResult: complex(0, 1), expectedError: ""}, // sinh(ix) = i sin(x)
		{name: "cosh(1)", input: "cosh(1)", expectedResult: cmplx.Cosh(complex(1, 0)), expectedError: ""},
		{name: "cosh(i*pi)", input: "cosh(i*pi)", expectedResult: complex(-1, 0), expectedError: ""}, // cosh(ix) = cos(x)
		{name: "tanh(1)", input: "tanh(1)", expectedResult: cmplx.Tanh(complex(1, 0)), expectedError: ""},
		{name: "tanh(i*pi/4)", input: "tanh(i*pi/4)", expectedResult: complex(0, 1), expectedError: ""}, // tanh(ix) = i tan(x)

		// Inverse Hyperbolic
		{name: "asinh(sinh(2))", input: "asinh(sinh(2))", expectedResult: complex(2, 0), expectedError: ""},
		{name: "acosh(cosh(2))", input: "acosh(cosh(2))", expectedResult: complex(2, 0), expectedError: ""},        // For real x >= 1
		{name: "acosh(0.5)", input: "acosh(0.5)", expectedResult: cmplx.Acosh(complex(0.5, 0)), expectedError: ""}, // Complex result
		{name: "atanh(tanh(0.5))", input: "atanh(tanh(0.5))", expectedResult: complex(0.5, 0), expectedError: ""},
		{name: "atanh(2)", input: "atanh(2)", expectedResult: cmplx.Atanh(complex(2, 0)), expectedError: ""}, // Complex result

		// Angle Conversion
		{name: "degToRad(180)", input: "degToRad(180)", expectedResult: complex(math.Pi, 0), expectedError: ""},
		{name: "degToRad(90+180i)", input: "degToRad(90 + 180i)", expectedResult: complex(math.Pi/2.0, math.Pi), expectedError: ""}, // (90+180i/pi) * (pi/180) = pi/2 + i
		{name: "radToDeg(pi)", input: "radToDeg(pi)", expectedResult: complex(180, 0), expectedError: ""},
		{name: "radToDeg(pi/2 + i)", input: "radToDeg(pi/2 + i)", expectedResult: complex(90, 180/math.Pi), expectedError: ""},

		// Component-wise Integer Functions
		{name: "floor(3.7+2.3i)", input: "floor(3.7+2.3i)", expectedResult: complex(3, 2), expectedError: ""},
		{name: "floor(-3.7-2.3i)", input: "floor(-3.7-2.3i)", expectedResult: complex(-4, -3), expectedError: ""},
		{name: "ceil(3.2+2.8i)", input: "ceil(3.2+2.8i)", expectedResult: complex(4, 3), expectedError: ""},
		{name: "ceil(-3.2-2.8i)", input: "ceil(-3.2-2.8i)", expectedResult: complex(-3, -2), expectedError: ""},
		{name: "round(3.5+2.5i)", input: "round(3.5+2.5i)", expectedResult: complex(4, 3), expectedError: ""}, // Go's math.Round rounds .5 to nearest even
		{name: "round(3.7+2.3i)", input: "round(3.7+2.3i)", expectedResult: complex(4, 2), expectedError: ""},
		{name: "trunc(3.7+2.3i)", input: "trunc(3.7+2.3i)", expectedResult: complex(3, 2), expectedError: ""},
		{name: "trunc(-3.7-2.3i)", input: "trunc(-3.7-2.3i)", expectedResult: complex(-3, -2), expectedError: ""},
	}
	runRPNEvalTests(t, testCases)
}

// It might be good to also have the TestCalculateExpressionEndToEnd function from before,
// which checks the final string output. The runRPNEvalTests helper above focuses on
// the complex128 result from EvaluateRPN. You can have both.
// For TestCalculateExpressionEndToEnd, you would define `expectedOutput` as a string.

// Example of how you might keep TestCalculateExpressionEndToEnd for string checks:
func TestCalculateExpressionEndToEnd_Stage2(t *testing.T) {
	testCases := []struct {
		name                   string
		input                  string
		expectedOutput         string
		expectedErrorSubstring string
	}{
		// Add selected cases here for string output verification
		{name: "E2E real(3+4i)", input: "real(3+4i)", expectedOutput: "3", expectedErrorSubstring: ""},
		{name: "E2E imag(3+4i)", input: "imag(3+4i)", expectedOutput: "4", expectedErrorSubstring: ""},
		{name: "E2E abs(3+4i)", input: "abs(3+4i)", expectedOutput: "5", expectedErrorSubstring: ""},
		{name: "E2E phase(1+i)", input: "phase(1+i)", expectedOutput: fmt.Sprintf("%0.10g", math.Pi/4), expectedErrorSubstring: ""},
		{name: "E2E conj(3+4i)", input: "conj(3+4i)", expectedOutput: "3-4i", expectedErrorSubstring: ""},
		{name: "E2E sqrt(-9)", input: "sqrt(-9)", expectedOutput: "3i", expectedErrorSubstring: ""},
		{name: "E2E degToRad(180)", input: "degToRad(180)", expectedOutput: fmt.Sprintf("%0.10g", math.Pi), expectedErrorSubstring: ""},
		{name: "E2E radToDeg(pi)", input: "radToDeg(pi)", expectedOutput: "180", expectedErrorSubstring: ""},
		{name: "E2E floor(3.7+2.3i)", input: "floor(3.7+2.3i)", expectedOutput: "3+2i", expectedErrorSubstring: ""},

		// Test NaN/Inf propagation to string output
		{name: "E2E log(0) string", input: "log(0)", expectedOutput: "(-Inf+0i)", expectedErrorSubstring: ""}, // Check your formatComplexOutput
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actualOutput, err := calculateExpression(tc.input)
			checkError(t, tc.expectedErrorSubstring, err)

			if tc.expectedErrorSubstring == "" && err == nil {
				if actualOutput != tc.expectedOutput {
					t.Errorf("Input '%s': Expected string output '%s', but got '%s'",
						tc.input, tc.expectedOutput, actualOutput)
				}
			}
		})
	}
}
