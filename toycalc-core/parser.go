// parser.go
package toycalc_core

import (
	"fmt"
	"strings"
	// "strconv"
)

// Constants for known identifiers (used for implied multiplication logic and IDENT handling)
// These maps should be kept in sync with what the evaluator can handle.
var (
	knownConstants = map[string]bool{
		"i":  true,
		"pi": true,
		"e":  true,
	}
	knownFunctions = map[string]bool{ // Stage 1 & 2 functions
		"log": true, "exp": true,
		"sin": true, "cos": true, "tan": true,
		"asin": true, "acos": true, "atan": true,
		"sinh": true, "cosh": true, "tanh": true,
		"asinh": true, "acosh": true, "atanh": true,
		"log10": true, "log2": true, "sqrt": true,
		"real": true, "imag": true, "abs": true, "phase": true, "conj": true,
		"degtorad": true, "radtodeg": true,
		"floor": true, "ceil": true, "round": true, "trunc": true,
	}
)

type Parser struct {
	tokens       []Token
	currentIndex int
	// currentToken is now managed internally by nextToken/peekToken, not a struct field directly always up-to-date
	// This avoids confusion as nextToken() now effectively means "consume and advance".

	outputQueue   []Token
	operatorStack []Token

	precedence      map[TokenType]int
	leftAssociative map[TokenType]bool

	// Tracks if the previous token suggests that the next token should be an operand (or a prefix unary operator)
	// This is true at the start, after '(', '[', '{', ',', or after another operator.
	expectOperand bool
}

func NewParser(tokens []Token) *Parser {
	p := &Parser{
		tokens: tokens, // Including EOF
		precedence: map[TokenType]int{
			// We might introduce UNARY_MINUS here with higher precedence if we take that path
			PLUS:        2,
			MINUS:       2,
			ASTERISK:    3,
			SLASH:       3,
			PERCENT:     3,
			CARET:       5,
			UNARY_MINUS: 4,
		},
		leftAssociative: map[TokenType]bool{
			PLUS:     true,
			MINUS:    true,
			ASTERISK: true,
			SLASH:    true,
			PERCENT:  true,
			CARET:    false,
		},
		expectOperand: true, // At the start of an expression, we expect an operand or unary prefix
	}
	return p
}

// peekCurrentToken looks at the token at currentIndex without advancing
/*func (p *Parser) peekCurrentToken() Token {
	if p.currentIndex < len(p.tokens) {
		return p.tokens[p.currentIndex]
	}
	// Should ideally not be called if currentIndex is already at EOF,
	// but return EOF if it happens. The last token from lexer is EOF.
	return p.tokens[len(p.tokens)-1] // This will be EOF
}*/

// consumeToken advances currentIndex and returns the token that was just consumed.
func (p *Parser) consumeToken() Token {
	if p.currentIndex < len(p.tokens) {
		tok := p.tokens[p.currentIndex]
		p.currentIndex++
		return tok
	}
	return p.tokens[len(p.tokens)-1] // Return EOF if called past the end
}

// Helper methods for operator stack (pushOperator, popOperator, peekOperator) - same as before

func (p *Parser) pushOperator(op Token) {
	p.operatorStack = append(p.operatorStack, op)
}

func (p *Parser) popOperator() (Token, bool) {
	if len(p.operatorStack) == 0 {
		return Token{Type: ILLEGAL}, false
	}
	op := p.operatorStack[len(p.operatorStack)-1]
	p.operatorStack = p.operatorStack[:len(p.operatorStack)-1]
	return op, true
}

func (p *Parser) peekOperator() (Token, bool) {
	if len(p.operatorStack) == 0 {
		return Token{Type: ILLEGAL}, false
	}
	return p.operatorStack[len(p.operatorStack)-1], true
}

// Type check helpers (isOperator, isFunction, isLeftParen, isRightParen, getMatchingLeftParen) - same as before
func isOperator(tokenType TokenType) bool { // Checks for binary operators for Shunting-Yard logic
	switch tokenType {
	case PLUS, MINUS, ASTERISK, SLASH, PERCENT, CARET, UNARY_MINUS:
		return true
	}
	return false
}

func isFunction(tokenType TokenType) bool {
	return tokenType == IDENT
}

func isLeftParen(tokenType TokenType) bool {
	switch tokenType {
	case LPAREN, LBRACKET, LBRACE:
		return true
	}
	return false
}

/*func isRightParen(tokenType TokenType) bool {
	switch tokenType {
	case RPAREN, RBRACKET, RBRACE:
		return true
	}
	return false
}*/

func getMatchingLeftParen(rightParenType TokenType) TokenType {
	switch rightParenType {
	case RPAREN:
		return LPAREN
	case RBRACKET:
		return LBRACKET
	case RBRACE:
		return LBRACE
	}
	return ILLEGAL
}

// ParseToRPN converts infix token stream to RPN (postfix)
func (p *Parser) ParseToRPN() ([]Token, error) {
	p.outputQueue = []Token{}
	p.operatorStack = []Token{}
	p.currentIndex = 0
	p.expectOperand = true // Reset at the start of parsing

	currentToken := p.consumeToken() // Get the first token

	for currentToken.Type != EOF {

		// --- Start of Implied Multiplication Logic ---
		if !p.expectOperand { // An operator is expected
			isOperandStarter := false
			switch currentToken.Type {
			case NUMBER, LPAREN, LBRACKET, LBRACE:
				isOperandStarter = true
			case IDENT:
				// An IDENT can start an operand if it's a constant or a function call
				lowerLiteral := strings.ToLower(currentToken.Literal)
				if _, isConst := knownConstants[lowerLiteral]; isConst {
					isOperandStarter = true
				} else if _, isFunc := knownFunctions[lowerLiteral]; isFunc {
					isOperandStarter = true // e.g. (1+2)log(x)
				}
			}

			// If an operand starter appears where an operator is expected,
			// AND it's not a PLUS/MINUS (which have their own unary handling),
			// then insert an implicit multiplication.
			if isOperandStarter && currentToken.Type != PLUS && currentToken.Type != MINUS {
				implicitAsterisk := Token{Type: ASTERISK, Literal: "*", Position: currentToken.Position} // Use pos of the token that implies mult

				// Process this virtual ASTERISK token using Shunting-Yard logic
				op1Implicit := implicitAsterisk
				for {
					op2Stack, ok := p.peekOperator()
					if !ok || isLeftParen(op2Stack.Type) {
						break
					}
					if (p.leftAssociative[op1Implicit.Type] && p.precedence[op1Implicit.Type] <= p.precedence[op2Stack.Type]) ||
						(!p.leftAssociative[op1Implicit.Type] && p.precedence[op1Implicit.Type] < p.precedence[op2Stack.Type]) {
						poppedOp, _ := p.popOperator()
						p.outputQueue = append(p.outputQueue, poppedOp)
					} else {
						break
					}
				}
				p.pushOperator(op1Implicit)
				p.expectOperand = true // After the implicit '*', we now expect an operand

				// DO NOT advance p.currentIndex here. The currentToken (e.g., LPAREN or IDENT)
				// still needs to be processed by the main switch in the *same* loop iteration
				// but now p.expectOperand is true.
				// The loop will re-evaluate the same currentToken. This needs careful loop control.

				// To achieve "re-evaluation" of currentToken with the new expectOperand state:
				// We simply let the loop continue. The currentToken for the switch below
				// will be the one that triggered the implicit multiplication.
			}
		}

		switch currentToken.Type {
		case NUMBER:
			if !p.expectOperand {
				// If we were not expecting an operand, it means an operator was missing
				// e.g., "2 3" or ") 3" or "x 3" (if x is a var/constant)
				return nil, NewCalculationError(
					fmt.Sprintf("unexpected number '%s' at position %d; an operator may be missing", currentToken.Literal, currentToken.Position),
				)
			}
			p.outputQueue = append(p.outputQueue, currentToken)
			p.expectOperand = false // After an operand, we expect an operator or closing paren

		case IDENT:
			isConstant := false
			isKnownFunction := false // We'll use this to differentiate known functions from unknown idents

			lowerLiteral := strings.ToLower(currentToken.Literal)
			switch lowerLiteral {
			case "i", "pi", "e": // Added "pi" and "e"
				isConstant = true
			case "log", "exp", // Stage 1 functions
				"sin", "cos", "tan", "asin", "acos", "atan", // Stage 2 trig
				"sinh", "cosh", "tanh", "asinh", "acosh", "atanh", // Stage 2 hyperbolic
				"log10", "log2", "sqrt", "real", "imag", "abs", "phase",
				"conj", "degtorad", "radtodeg", "floor", "ceil", "round", "trunc": // Stage 2 other
				isKnownFunction = true
			}

			if isConstant {
				if !p.expectOperand {
					return nil, NewCalculationError(
						fmt.Sprintf("unexpected constant '%s' at position %d; an operator may be missing", currentToken.Literal, currentToken.Position),
					)
				}
				p.outputQueue = append(p.outputQueue, currentToken) // Token is {IDENT, "pi", pos}, etc.
				p.expectOperand = false                             // After an operand/constant, we expect an operator
			} else if isKnownFunction {
				p.pushOperator(currentToken) // Function name goes to operator stack
				// expectOperand state is managed by LPAREN that should follow a function
			} else {
				// Unknown identifier
				return nil, NewCalculationError(
					fmt.Sprintf("unknown identifier or function '%s' at position %d", currentToken.Literal, currentToken.Position),
				)
			}
		case PLUS, MINUS:
			operatorToken := currentToken
			if p.expectOperand { // Potential unary operator
				if operatorToken.Type == MINUS {
					// Convert to a UNARY_MINUS token
					// This UNARY_MINUS will be pushed onto the operator stack with its own (higher) precedence.
					operatorToken = Token{Type: UNARY_MINUS, Literal: "-", Position: operatorToken.Position}
				} else { // Unary PLUS
					// Unary PLUS can be ignored; it doesn't change the value.
					// We simply consume it and expect an operand next.
					p.expectOperand = true // We still need an operand after an ignored unary plus
					currentToken = p.consumeToken()
					continue // Move to the next token
				}
			}
			// Now operatorToken is either the original binary PLUS/MINUS, or the new UNARY_MINUS.
			// Proceed with Shunting-Yard logic for this operator.
			for {
				op2, ok := p.peekOperator()
				if !ok || isLeftParen(op2.Type) {
					break
				}
				// For UNARY_MINUS (which is right-associative), only pop op2 if op2 has strictly greater precedence.
				// For binary ops, use the existing left-associativity check for equal precedence.
				if (p.leftAssociative[operatorToken.Type] && p.precedence[op2.Type] >= p.precedence[operatorToken.Type]) ||
					(!p.leftAssociative[operatorToken.Type] && p.precedence[op2.Type] > p.precedence[operatorToken.Type]) {
					poppedOp, _ := p.popOperator()
					p.outputQueue = append(p.outputQueue, poppedOp)
				} else {
					break
				}
			}
			p.pushOperator(operatorToken)
			p.expectOperand = true // After any operator (unary or binary), we expect an operand

		case ASTERISK, SLASH, PERCENT, CARET: // These are always binary in this context
			if p.expectOperand {
				// This means an operator like '*' appeared where an operand was expected, e.g., "* 5" or "( * 5)"
				return nil, NewCalculationError(fmt.Sprintf("unexpected operator '%s' at position %d; operand expected", currentToken.Literal, currentToken.Position))
			}
			op1 := currentToken
			for {
				op2, ok := p.peekOperator()
				if !ok || isLeftParen(op2.Type) {
					break
				}
				if (p.precedence[op2.Type] > p.precedence[op1.Type]) ||
					(p.precedence[op2.Type] == p.precedence[op1.Type] && p.leftAssociative[op1.Type]) {
					p.popOperator()
					p.outputQueue = append(p.outputQueue, op2)
				} else {
					break
				}
			}
			p.pushOperator(op1)
			p.expectOperand = true // After a binary operator, we expect an operand

		case COMMA:
			if p.expectOperand { // Comma should not appear where an operand is expected right before it
				return nil, NewCalculationError(fmt.Sprintf("unexpected comma at position %d; operand expected before comma", currentToken.Position))
			}
			foundLeftParen := false
			for len(p.operatorStack) > 0 {
				op, _ := p.peekOperator()
				if isLeftParen(op.Type) {
					foundLeftParen = true
					break
				}
				poppedOp, _ := p.popOperator()
				p.outputQueue = append(p.outputQueue, poppedOp)
			}
			if !foundLeftParen {
				return nil, NewCalculationError(fmt.Sprintf("mismatched comma or parentheses at position %d", currentToken.Position))
			}
			p.expectOperand = true // After a comma, we expect another argument (operand)

		case LPAREN, LBRACKET, LBRACE:
			// If IDENT (function name) was the previous token pushed to opStack, this LPAREN confirms it's a function call.
			// Check if previous token pushed to opStack was IDENT to confirm function call.
			if op, ok := p.peekOperator(); ok && isFunction(op.Type) {
				// It's a function call. The IDENT is already on stack.
				// Push the LPAREN.
			} else if !p.expectOperand {
				// We have something like "5(" or ")(" which implies multiplication.
				// This is for Stage 4 (implied multiplication). For now, it's an error.
				return nil, NewCalculationError(fmt.Sprintf("unexpected parenthesis '%s' at position %d; operator expected or implied multiplication not supported", currentToken.Literal, currentToken.Position))
			}
			p.pushOperator(currentToken)
			p.expectOperand = true // After '(', we expect an operand (or unary operator)

		case RPAREN, RBRACKET, RBRACE:
			/*if p.expectOperand && len(p.outputQueue) > 0 && p.outputQueue[len(p.outputQueue)-1].Type != COMMA {
				// Case like "( )" or "(,". Check if the output queue's state makes sense.
				// If we expect an operand but see a ')' and the output queue isn't empty due to an argument,
				// it might be an empty pair of parentheses like `()` in `f()`.
				// This needs careful thought for empty function calls if allowed.
				// For now, `log()` would fail as `log` expects an arg.
				// A simple `()` might be invalid.
				// If the parser encounters something like `(`, `EOF` without an operand for `log(`, this check is too late.
				// The check `if p.expectOperand` implies nothing was pushed to outputQueue since last operator/LPAREN/comma
			}*/
			if p.expectOperand && (len(p.outputQueue) == 0 || isOperator(p.outputQueue[len(p.outputQueue)-1].Type) || isLeftParen(p.outputQueue[len(p.outputQueue)-1].Type) || p.outputQueue[len(p.outputQueue)-1].Type == COMMA) {
				// This means something like `()` or `(,)` or `(*)` which is an error if an operand was expected
				// but the part before `)` is not a valid operand.
				// Example: `log()` - `log` is on opStack, `(` is on opStack. `)` comes. `expectOperand` is true.
				// This situation would mean no argument was provided for the function.
				return nil, NewCalculationError(fmt.Sprintf("missing operand before closing parenthesis '%s' at position %d", currentToken.Literal, currentToken.Position))
			}

			expectedLeftParen := getMatchingLeftParen(currentToken.Type)
			foundMatchingParen := false
			for len(p.operatorStack) > 0 {
				op, _ := p.peekOperator()
				if op.Type == expectedLeftParen {
					_, _ = p.popOperator() // Discard the left parenthesis
					foundMatchingParen = true
					break
				}
				poppedOp, _ := p.popOperator()
				p.outputQueue = append(p.outputQueue, poppedOp)
			}
			if !foundMatchingParen {
				return nil, NewCalculationError(fmt.Sprintf("mismatched parentheses/brackets/braces for '%s' at position %d", currentToken.Literal, currentToken.Position))
			}
			// If token at top of stack is a function name, pop it to output.
			if op, ok := p.peekOperator(); ok && isFunction(op.Type) {
				poppedFunc, _ := p.popOperator()
				p.outputQueue = append(p.outputQueue, poppedFunc)
			}
			p.expectOperand = false // After ')', we expect an operator

		default: // Should be unreachable if lexer is correct
			return nil, NewCalculationError(fmt.Sprintf("parser encountered unexpected token '%s' (type %s) at position %d", currentToken.Literal, currentToken.Type, currentToken.Position))
		}
		currentToken = p.consumeToken() // Consume current token and advance to the next
	}

	// After loop, pop all remaining operators from stack to output queue
	for len(p.operatorStack) > 0 {
		op, _ := p.popOperator()
		if isLeftParen(op.Type) {
			return nil, NewCalculationError(fmt.Sprintf("mismatched parentheses/brackets/braces at end (unclosed '%s' at pos %d)", op.Literal, op.Position))
		}
		p.outputQueue = append(p.outputQueue, op)
	}

	// Final check: if the output queue is empty, it means no valid expression was formed.
	if len(p.outputQueue) == 0 && (p.currentIndex > 0 && p.tokens[0].Type != EOF) { // Make sure it wasn't just an empty input that Parse() should have caught
		return nil, NewCalculationError("parsed expression resulted in empty RPN queue (invalid expression structure)")
	}

	return p.outputQueue, nil
}

// Main Parse function (entry point) - from before
func Parse(tokens []Token) ([]Token, error) {
	if len(tokens) == 0 {
		return nil, NewCalculationError("no tokens provided to parse (empty token slice)")
	}
	if len(tokens) == 1 && tokens[0].Type == EOF {
		return nil, NewCalculationError("no expression provided to parse (only EOF token found)")
	}
	parser := NewParser(tokens)
	return parser.ParseToRPN()
}
