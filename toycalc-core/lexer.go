package toycalc_core

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

type Lexer struct {
	input        string
	position     int  // current position in input (points to current char)
	readPosition int  // current reading position in input (after current char)
	ch           rune // current char under examination
}

func NewLexer(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar() // Initialize l.ch, l.position, and l.readPosition
	return l
}

// readChar gives us the next character and advances our position in the input string.
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0                  // 0 is ASCII code for "NUL", signifies EOF or not read yet
		l.position = len(l.input) // Position for EOF is at the very end
	} else {
		r, size := utf8.DecodeRuneInString(l.input[l.readPosition:])
		l.ch = r
		l.position = l.readPosition
		l.readPosition += size
	}
}

// peekChar looks ahead in the input without consuming the character.
/*func (l *Lexer) peekChar() rune {
	if l.readPosition >= len(l.input) {
		return 0
	}
	r, _ := utf8.DecodeRuneInString(l.input[l.readPosition:])
	return r
}*/

func (l *Lexer) NextToken() Token {
	var tok Token

	l.skipWhitespace()

	tokenStartPosition := l.position // Capture start position before consuming the token

	switch l.ch {
	case '+':
		tok = Token{Type: PLUS, Literal: "+", Position: tokenStartPosition}
	case '-':
		tok = Token{Type: MINUS, Literal: "-", Position: tokenStartPosition}
	case '*':
		tok = Token{Type: ASTERISK, Literal: "*", Position: tokenStartPosition}
	case '/':
		tok = Token{Type: SLASH, Literal: "/", Position: tokenStartPosition}
	case '%':
		tok = Token{Type: PERCENT, Literal: "%", Position: tokenStartPosition}
	case '^':
		tok = Token{Type: CARET, Literal: "^", Position: tokenStartPosition}
	case '(':
		tok = Token{Type: LPAREN, Literal: "(", Position: tokenStartPosition}
	case ')':
		tok = Token{Type: RPAREN, Literal: ")", Position: tokenStartPosition}
	case '[':
		tok = Token{Type: LBRACKET, Literal: "[", Position: tokenStartPosition}
	case ']':
		tok = Token{Type: RBRACKET, Literal: "]", Position: tokenStartPosition}
	case '{':
		tok = Token{Type: LBRACE, Literal: "{", Position: tokenStartPosition}
	case '}':
		tok = Token{Type: RBRACE, Literal: "}", Position: tokenStartPosition}
	case ',':
		tok = Token{Type: COMMA, Literal: ",", Position: tokenStartPosition}
	case 0: // EOF
		tok = Token{Type: EOF, Literal: "", Position: tokenStartPosition}
	default:
		if isLetter(l.ch) {
			literal := l.readIdentifier() // readIdentifier consumes chars & updates l.ch, l.position
			tok = Token{Type: IDENT, Literal: literal, Position: tokenStartPosition}
			return tok // Return directly; readIdentifier already advanced past the token
		} else if isDigit(l.ch) {
			literal := l.readNumber() // readNumber consumes chars & updates l.ch, l.position
			tok = Token{Type: NUMBER, Literal: literal, Position: tokenStartPosition}
			return tok // Return directly; readNumber already advanced past the token
		} else {
			// For an illegal character, the literal is just that one character.
			tok = Token{Type: ILLEGAL, Literal: string(l.ch), Position: tokenStartPosition}
			// Fall through to call l.readChar() to advance past this illegal char
		}
	}

	// For single-character tokens (operators, delimiters) or an ILLEGAL char tokenized above,
	// advance the lexer to prepare for the next token.
	// This is NOT called if readIdentifier or readNumber returned early.
	l.readChar()
	return tok
}

func (l *Lexer) skipWhitespace() {
	for unicode.IsSpace(l.ch) {
		l.readChar()
	}
}

// readIdentifier reads in an identifier and advances the lexer's position until it
// encounters a non-letter character.
func (l *Lexer) readIdentifier() string {
	startPosition := l.position
	for isLetter(l.ch) || (l.position != startPosition && isDigit(l.ch)) {
		l.readChar()
	}
	return l.input[startPosition:l.position]
}

// readNumber reads in a number (integer or float) and advances the lexer's position.
func (l *Lexer) readNumber() string {
	position := l.position
	hasDecimal := false
	for isDigit(l.ch) || (l.ch == '.' && !hasDecimal) {
		if l.ch == '.' {
			hasDecimal = true
		}
		l.readChar()
	}
	// Basic scientific notation (e.g., 1.2e-3 or 3E+2)
	if strings.ToLower(string(l.ch)) == "e" {
		l.readChar() // consume 'e' or 'E'
		if l.ch == '+' || l.ch == '-' {
			l.readChar() // consume sign
		}
		if !isDigit(l.ch) { // Must have at least one digit after 'e' or 'e[sign]'
			// This indicates an incomplete scientific notation, could be an error
			// For simplicity, we'll stop here, but a full lexer would validate more.
			return l.input[position:l.position]
		}
		for isDigit(l.ch) {
			l.readChar()
		}
	}
	return l.input[position:l.position]
}

// Helper functions for character types
func isLetter(ch rune) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
	// Allow underscore for identifiers, common in function names
}

func isDigit(ch rune) bool {
	return '0' <= ch && ch <= '9'
}

// Lex function (adjust error message for new Token.Position)
func Lex(input string) ([]Token, error) {
	l := NewLexer(input)
	var tokens []Token
	for {
		tok := l.NextToken()
		tokens = append(tokens, tok)
		if tok.Type == EOF {
			break
		}
		if tok.Type == ILLEGAL {
			return tokens, NewCalculationError(fmt.Sprintf("illegal character '%s' found at position %d", tok.Literal, tok.Position))
		}
	}
	return tokens, nil
}
