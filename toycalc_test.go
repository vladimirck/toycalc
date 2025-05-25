// toycalc_test.go
package main

import (
	"math"
	"math/cmplx"
	"reflect"
	"testing"
)

// --- Test Helper Functions (to be developed as needed) ---

func checkError(t *testing.T, expectedErrorMsg string, actualError error) {
	t.Helper()
	if actualError == nil {
		if expectedErrorMsg != "" {
			t.Errorf("Expected error '%s', but got nil.", expectedErrorMsg)
		}
		return
	}
	if expectedErrorMsg == "" {
		t.Errorf("Expected no error, but got: '%s'", actualError.Error())
		return
	}
	// Basic string comparison, might need more robust error type checking later
	if actualError.Error() != expectedErrorMsg {
		// If using CalculationError, could check calcErr.Message
		t.Errorf("Expected error message '%s', got '%s'", expectedErrorMsg, actualError.Error())
	}
}

func compareComplex(t *testing.T, expected, actual complex128, description string) {
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
		if !cmplx.IsInf(actual) || (real(expected)*real(actual) < 0 && real(expected) != 0) || (imag(expected)*imag(actual) < 0 && imag(expected) != 0) {
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
		t.Errorf("%s: token slices do not match.\nExpected: %v\nGot:      %v", description, expected, actual)
	}
}

// --- Initial Setup Tests ---
func TestFrameworkInitialSetup(t *testing.T) {
	t.Log("Test framework is running. Add more tests as features are implemented.")
	if 1+1 != 2 {
		t.Error("Something is fundamentally wrong if 1+1 is not 2.")
	}
}

// --- Lexer Tests (To be developed in Stage 1) ---
// func TestLexerNumbersAndOperators(t *testing.T) { /* ... */ }

// --- Parser Tests (To be developed in Stage 1) ---
// func TestParserSimpleArithmetic(t *testing.T) { /* ... */ }

// --- Evaluator Tests (To be developed in Stage 1) ---
// func TestEvaluatorAddition(t *testing.T) { /* ... */ }

// --- Output Formatting Tests (To be developed in Stage 1) ---
// func TestFormatOutput(t *testing.T) { /* ... */ }
