package parser

import (
	"fmt"

	"github.com/mateors/graphqlparser/ast"
	"github.com/mateors/graphqlparser/lexer"
	"github.com/mateors/graphqlparser/token"
)

type Parser struct {
	l         *lexer.Lexer
	errors    []string
	curToken  token.Token //position
	peekToken token.Token //read position
	//prefixParseFns map[token.TokenType]prefixParseFn
	//infixParseFns  map[token.TokenType]infixParseFn
}

func New(l *lexer.Lexer) *Parser {

	p := &Parser{l: l, errors: []string{}}

	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseDocument() *ast.Document {

	document := &ast.Document{} //root node
	document.Definitions = []ast.Node{}

	for p.curToken.Type != token.EOF {
		doc := p.parseDocument()
		if doc != nil {
			document.Definitions = append(document.Definitions, doc)
		}
		p.nextToken()
	}
	return document
}

func (p *Parser) parseDocument() ast.Node { //ast.Definition

	switch p.curToken.Type {
	case token.TYPE:
		return p.parseTypeSystemDefinition()

	default:
		return &ast.OperationDefinition{}
	}
}

func (p *Parser) parseTypeSystemDefinition() ast.Node { //ast.TypeSystemDefinition

	switch p.curToken.Type {
	case token.TYPE:
		return p.parseObjectDefinition()
	case token.IDENT:
		return nil
	default:
		fmt.Println(">>", p.curToken)
		return nil
	}
}

func (p *Parser) parseObjectDefinition() ast.Node {

	od := &ast.ObjectDefinition{}
	return od
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

// Converts a name lex token into a name parse node.
func (p *Parser) parseName() *ast.Name {
	if !p.expectPeek(token.IDENT) {
		return nil
	}
	return &ast.Name{Kind: ast.NAME, Token: p.curToken, Value: p.curToken.Literal}
}

/**
 * NamedType : Name
 */
func (p *Parser) parseNamed() *ast.NamedType {
	name := p.parseName()
	return &ast.NamedType{Kind: ast.NAMED_TYPE, Token: p.curToken, Name: name}
}

/**
 * Type :
 *   - NamedType
 *   - ListType
 *   - NonNullType
 */
func (p *Parser) parseType() (ttype ast.Type) {

	// [ String! ]!
	switch p.curToken.Type {
	case token.LBRACKET: //[
		p.nextToken()
		ttype = p.parseType()
		fallthrough

	case token.RBRACKET: //]
		p.nextToken()
		ttype = &ast.ListType{Kind: ast.LIST_TYPE, Token: p.curToken, Type: ttype}

	case token.IDENT:
		ttype = p.parseNamed()
	}

	// BANG must be executed
	if p.expectPeek(token.BANG) {
		p.nextToken()
		ttype = &ast.NonNullType{Kind: ast.NONNULL_TYPE, Token: p.curToken, Type: ttype}
	}
	return ttype
}
