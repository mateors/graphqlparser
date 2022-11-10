package lexer

import (
	"fmt"

	"github.com/mateors/graphqlparser/token"
)

type Lexer struct {
	input        string
	position     int
	readPosition int
	line         int
	ch           byte
}

func New(input string) *Lexer {

	lex := &Lexer{}
	lex.input = input
	lex.line = 1
	lex.readChar()
	return lex
}

func (l *Lexer) NextToken() token.Token {

	var tok token.Token
	//tok.Line = l.line
	l.skipWhitespace()

	switch l.ch {

	case '=':
		tok = newToken(token.ASSIGN, l.ch, l.line, l.position, l.position+1)

	case '-':
		tok = newToken(token.MINUS, l.ch, l.line, l.position, l.position+1)

	case '{':
		tok = newToken(token.LBRACE, l.ch, l.line, l.position, l.position+1)

	case '}':
		tok = newToken(token.RBRACE, l.ch, l.line, l.position, l.position+1)

	case ':':
		tok = newToken(token.COLON, l.ch, l.line, l.position, l.position+1)

	case '!':
		tok = newToken(token.BANG, l.ch, l.line, l.position, l.position+1)

	case '(':
		tok = newToken(token.LPAREN, l.ch, l.line, l.position, l.position+1)

	case ')':
		tok = newToken(token.RPAREN, l.ch, l.line, l.position, l.position+1)

	case '[':
		tok = newToken(token.LBRACKET, l.ch, l.line, l.position, l.position+1)

	case ']':
		tok = newToken(token.RBRACKET, l.ch, l.line, l.position, l.position+1)

	case '|':
		tok = newToken(token.PIPE, l.ch, l.line, l.position, l.position+1)

	case '"':

		ch := l.ch
		nchar := l.input[l.readPosition+1]
		if ch == '"' && l.peekChar() == '"' && nchar == '"' {
			l.readChar()
			l.readChar()
			tok.Line = l.line
			tok.Start = l.position + 1
			tok.Literal = l.blockString()
			tok.Type = token.BLOCK_STRING
			tok.End = l.position
			l.readChar()
			l.readChar()

		} else {
			tok.Line = l.line
			tok.Start = l.position + 1
			tok.Type = token.STRING
			tok.Literal = l.readString()
			tok.End = tok.Start + len(tok.Literal)
		}

	case ',':
		tok = newToken(token.COMMA, l.ch, l.line, l.position, l.position+1)

	case '$':
		tok = newToken(token.DOLLAR, l.ch, l.line, l.position, l.position+1)

	case '&':
		tok = newToken(token.AMP, l.ch, l.line, l.position, l.position+1)

	case '@':
		tok = newToken(token.AT, l.ch, l.line, l.position, l.position+1)

	case '#':
		tok.Start = 0
		for {
			l.readChar()
			if isLetter(l.ch) && tok.Start == 0 {
				tok.Start = l.position + 0
			}
			if l.ch == '\n' { //or check for next token
				tok.End = l.position
				break
			}
		}
		//fmt.Println("2START:", tok.Start)
		tok.Line = l.line
		tok.Type = token.HASH
		//text := l.readString()
		//tok.End = tok.Start + len(text) + 1
		tok.Literal = l.input[tok.Start:tok.End]
		return tok

	case 0:
		tok.Line = l.line
		tok.Start = l.position
		tok.Type = token.EOF
		tok.Literal = ""
		tok.End = l.position

	default:

		ch := l.ch
		if isLetter(l.ch) {
			tok.Line = l.line
			tok.Start = l.position
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			tok.End = tok.Start + len(tok.Literal)
			return tok

		} else if isDigit(ch) || ch == '.' && isDigit(l.peekChar()) {
			tok.Line = l.line
			tok.Start = l.position
			tok.Type, tok.Literal = l.readNumber()
			tok.End = tok.Start + len(tok.Literal)
			return tok

		} else if ch == '.' {

			if l.peekChar() == '.' {
				tok.Line = l.line
				tok.Start = l.position
				tok.Literal = l.readVariadic()
				tok.Type = token.LookupIdent(tok.Literal)
				tok.End = tok.Start + len(tok.Literal)
				return tok
			}

		} else {
			tok = newToken(token.ILLEGAL, ch, l.line, l.position, l.position+1)
		}
	}
	l.readChar()
	return tok
}

func (l *Lexer) digits(base int, invalid *int) (digsep int) {

	if base <= 10 {
		max := byte('0' + base)
		for isDigit(l.ch) {
			ds := 1
			if l.ch >= max && *invalid < 0 {
				*invalid = l.position //s.offset // record invalid rune offset
			}
			digsep |= ds
			l.readChar()
		}
	}
	return
}

func (l *Lexer) blockString() string {

	position := l.position + 1
	for {
		l.readChar()
		if l.ch == '"' && l.input[l.position+1] == '"' && l.input[l.position+2] == '"' {
			break
		}
	}
	return l.input[position:l.position]
}

func (l *Lexer) readNumber() (token.TokenType, string) {

	position := l.position
	tok := token.ILLEGAL
	digsep := 0
	base := 10
	invalid := -1

	//integer part
	if l.ch != '.' {
		tok = token.INT
		digsep |= l.digits(base, &invalid)
	}

	//fractional
	if l.ch == '.' {
		tok = token.FLOAT
		l.readChar()
		digsep |= l.digits(base, &invalid)
	}

	//exponent
	if e := lower(l.ch); e == 'e' {
		l.readChar()
		tok = token.FLOAT
		if l.ch == '+' || l.ch == '-' {
			l.readChar()
		}
		ds := l.digits(10, nil)
		digsep |= ds
		if ds&1 == 0 {
			fmt.Println(l.position, "exponent has no digits")
		}
	}

	lit := l.input[position:l.position]
	if tok == token.INT && invalid >= 0 {
		fmt.Printf("invalid digit %q", lit[invalid-position])
	}
	return tok, lit
}

func newToken(tokenType token.TokenType, ch byte, line, start, end int) token.Token {
	return token.Token{Type: tokenType, Line: line, Start: start, End: end, Literal: string(ch)}
}

func (l *Lexer) readChar() {

	if len(l.input) <= l.readPosition {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
		if l.ch == '\n' {
			l.line += 1
		}
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
	for isLetter(l.ch) || isDigit(l.ch) {
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
	if len(l.input) > l.readPosition {
		return l.input[l.readPosition]
	}
	return 0
}

func (l *Lexer) readVariadic() string {
	position := l.position
	for l.ch == '.' {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readString() string {
	position := l.position + 1
	for {
		l.readChar()
		if l.ch == '"' || l.ch == 0 {
			break
		}
	}
	return l.input[position:l.position]
}

// returns lower-case ch iff ch is ASCII letter
func lower(ch byte) byte {
	return ('a' - 'A') | ch
}

func isHex(ch byte) bool {
	return '0' <= ch && ch <= '9' || 'a' <= lower(ch) && lower(ch) <= 'f'
}
