package parser

import (
	"fmt"

	"github.com/mateors/graphqlparser/ast"
	"github.com/mateors/graphqlparser/lexer"
	"github.com/mateors/graphqlparser/token"
)

type parseDefinitionFn func() ast.Node

//var tokenDefinitionFn map[string]parseDefinitionFn

type Parser struct {
	l                  *lexer.Lexer
	errors             []string
	curToken           token.Token //position
	peekToken          token.Token //read position
	tokenDefinitionFns map[token.TokenType]parseDefinitionFn
	//prefixParseFns map[token.TokenType]prefixParseFn
	//infixParseFns  map[token.TokenType]infixParseFn
}

func New(l *lexer.Lexer) *Parser {

	p := &Parser{l: l, errors: []string{}}
	p.tokenDefinitionFns = make(map[token.TokenType]parseDefinitionFn)

	//p.registerTokenDefinitionFns(token.LBRACE, parseOperationDefinition)
	p.registerTokenDefinitionFns(token.STRING, p.parseTypeSystemDefinition)
	p.registerTokenDefinitionFns(token.BLOCK_STRING, p.parseTypeSystemDefinition)

	p.registerTokenDefinitionFns(token.IDENT, p.parseTypeSystemDefinition)
	p.registerTokenDefinitionFns(token.TYPE, p.parseTypeSystemDefinition)

	//p.registerTokenDefinitionFns(token.TYPE, p.parseObjectDefinition)
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) registerTokenDefinitionFns(tokenType token.TokenType, fn parseDefinitionFn) {

	p.tokenDefinitionFns[tokenType] = fn
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseDocument() *ast.Document {

	document := &ast.Document{Kind: ast.DOCUMENT} //root node
	document.Definitions = []ast.Node{}

	for p.curToken.Type != token.EOF {
		doc := p.parseDocument()
		if doc != nil {
			document.Definitions = append(document.Definitions, doc)
		}
		fmt.Println("...")
		p.nextToken()
	}
	return document
}

func (p *Parser) parseDocument() ast.Node { //ast.Definition

	fmt.Println("parseDocument>", p.curToken.Type)
	switch p.curToken.Type {

	//case token.TYPE:
	//return p.parseObjectDefinition()

	case token.BLOCK_STRING, token.STRING, token.HASH:
		return p.parseObjectDefinition()

	// case token.IDENT: //,token.LBRACE, token.STRING

	// 	fmt.Println("tokenDefinitionFns->", p.curToken.Type)
	// 	parseFunc := p.tokenDefinitionFns[p.curToken.Type]
	// 	return parseFunc()

	//case token.IDENT:
	//return p.parseFieldDefinition()

	default:
		fmt.Println("unexpected", p.curToken.Type)
		return nil //&ast.OperationDefinition{}
	}

}

func (p *Parser) parseTypeSystemDefinition() ast.Node { //ast.TypeSystemDefinition

	fmt.Println("parseTypeSystemDefinition>", p.curToken)

	var keyWordToken token.Token
	if p.isDescription() {
		keyWordToken = p.peekToken
	}

	if !p.peekTokenIsKeyword() {
		p.peekError(token.IDENT)
		return nil
	}

	//fmt.Println("c2>", p.curToken, keyWordToken)
	item, ok := p.tokenDefinitionFns[keyWordToken.Type]
	if !ok {
		return nil
	}
	return item()
}

func (p *Parser) isDescription() bool {
	if p.curTokenIs(token.BLOCK_STRING) || p.curTokenIs(token.STRING) {
		return true
	}
	return false
}

func (p *Parser) parseFieldDefinition() *ast.FieldDefinition {

	fmt.Println("fieldDefinition", p.curToken) //starting with token.IDENT
	fd := &ast.FieldDefinition{}
	fd.Kind = ast.FIELD_DEFINITION
	fd.Token = p.curToken
	fd.Name = p.parseName()

	if !p.expectPeek(token.COLON) {
		return nil
	}

	p.nextToken()
	fd.Type = p.parseType()
	//fmt.Println("-->", fd.Type, p.curToken)
	return fd
}

func (p *Parser) expectToken(t token.TokenType) bool {
	fmt.Println("expectToken", t)
	if p.curTokenIs(t) {
		p.nextToken()
		return true
	}
	return false
}

func (p *Parser) parseObjectDefinition() ast.Node {

	fmt.Println("parseObjectDefinition", p.curToken) //starting from first token
	od := &ast.ObjectDefinition{Kind: ast.OBJECT_DEFINITION}
	od.Token = p.curToken
	od.Description = p.parseDescription()

	//fmt.Println("AfterDescription:", p.curToken)
	if !p.expectToken(token.TYPE) {
		//fmt.Println("*nil*")
		return nil
	}

	od.Name = p.parseName()
	//fmt.Println("od.Name", od.Name, p.curToken)

	//loop
	p.nextToken()
	p.nextToken()
	od.Fields = []*ast.FieldDefinition{}
	for !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF) {

		fd := p.parseFieldDefinition()
		//fmt.Println("fieldD--->", fd)
		if fd != nil {
			od.Fields = append(od.Fields, fd)
		}
		p.nextToken()
	}
	fmt.Println("DONE")
	return od
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) peekTokenIsKeyword() bool {
	return token.IsKeyword(p.peekToken.Literal)
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %v, got %v instead", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		//fmt.Println("peekTokenIspeekTokenIs")
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

// Converts a name lex token into a name parse node.
func (p *Parser) parseName() *ast.Name {

	if !p.curTokenIs(token.IDENT) {
		return nil
	}
	return &ast.Name{Kind: ast.NAME, Token: p.curToken, Value: p.curToken.Literal}
}

/**
 * NamedType : Name
 */
func (p *Parser) parseNamed() *ast.NamedType {

	//fmt.Println("parseNamed()", p.curToken)
	name := p.parseName()
	//fmt.Println("name:", name)
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

	case token.IDENT, token.STRING:
		ttype = p.parseNamed()
	}

	// BANG must be executed
	if p.expectPeek(token.BANG) {
		ttype = &ast.NonNullType{Kind: ast.NONNULL_TYPE, Token: p.curToken, Type: ttype}
	}
	return ttype
}

func (p *Parser) parseDescription() *ast.StringValue {

	fmt.Println("parseDescription", p.curToken, p.peekToken)
	if p.curTokenIs(token.STRING) || p.curTokenIs(token.BLOCK_STRING) || p.curTokenIs(token.HASH) {
		p.nextToken()
		return p.parseStringLiteral()
	}
	return nil
}

func (p *Parser) parseStringLiteral() *ast.StringValue {

	fmt.Println("parseStringLiteral", p.curToken)
	cToken := p.curToken
	p.nextToken()
	return &ast.StringValue{Kind: ast.STRING_VALUE, Token: cToken, Value: cToken.Literal}

}
