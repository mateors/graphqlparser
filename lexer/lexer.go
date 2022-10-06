package lexer

import (
	"github.com/mateors/graphqlparser/token"
)

type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           byte
}

func New(input string) *Lexer {

	lex := &Lexer{}
	lex.input = input
	lex.readChar()
	return lex
}

func (l *Lexer) NextToken() token.Token {

	var tok token.Token
	l.skipWhitespace()

	switch l.ch {

	case '=':
		tok = newToken(token.ASSIGN, l.ch)

	case '{':
		tok = newToken(token.LBRACE, l.ch)

	case '}':
		tok = newToken(token.RBRACE, l.ch)

	case ':':
		tok = newToken(token.COLON, l.ch)

	case '!':
		tok = newToken(token.BANG, l.ch)

	case '(':
		tok = newToken(token.LPAREN, l.ch)

	case ')':
		tok = newToken(token.RPAREN, l.ch)

	case '[':
		tok = newToken(token.LBRACKET, l.ch)

	case ']':
		tok = newToken(token.RBRACKET, l.ch)

	case '|':
		tok = newToken(token.PIPE, l.ch)

	case '.':

		if l.peekChar() == '.' {
			tok.Literal = l.readVariadic()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok

		} else {

			//number or anything else
		}

	case 0:
		tok.Type = token.EOF
		tok.Literal = ""

	default:

		//fmt.Println(l.ch)
		if isLetter(l.ch) {

			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok

		} else if isDigit(l.ch) {
			tok.Type = token.INT
			tok.Literal = l.readNumber()
			return tok

		} else {
			tok = newToken(token.ILLEGAL, l.ch)
		}
	}

	l.readChar()
	return tok
}

func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}

func (l *Lexer) readChar() {

	if len(l.input) <= l.readPosition {
		l.ch = 0

	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *Lexer) readIdentifier() string {

	position := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}
}

func (l *Lexer) readVariadic() string {

	position := l.position
	for l.ch == '.' {
		l.readChar()
	}
	return l.input[position:l.position]
}
