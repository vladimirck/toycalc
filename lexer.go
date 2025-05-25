package main

// Lex (stub) - Analyzes the input string and returns a list of Tokens
func Lex(input string) ([]Token, error) {
	// Actual implementation will come in Stage 1
	// For now, a simple stub for basic testing:
	if input == "" {
		return []Token{{Type: EOF, Literal: ""}}, nil
	}
	// Simulate the whole input as a single number for the stub evaluator test
	return []Token{{Type: NUMBER, Literal: input}, {Type: EOF, Literal: ""}}, nil
}
