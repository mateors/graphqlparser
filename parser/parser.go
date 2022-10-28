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
	//p.registerTokenDefinitionFns(token.STRING, p.parseTypeSystemDefinition)
	//p.registerTokenDefinitionFns(token.BLOCK_STRING, p.parseTypeSystemDefinition)
	//p.registerTokenDefinitionFns(token.IDENT, p.parseTypeSystemDefinition)
	//p.registerTokenDefinitionFns(token.TYPE, p.parseTypeSystemDefinition)

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

	case token.BLOCK_STRING, token.STRING, token.HASH, token.TYPE:
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

func (p *Parser) parseObjectDefinition() ast.Node {

	//fmt.Println("parseObjectDefinition->START", p.curToken) //starting from first token
	od := &ast.ObjectDefinition{Kind: ast.OBJECT_DEFINITION}
	od.Token = p.curToken
	od.Description = p.parseDescription()

	if !p.expectToken(token.TYPE) {
		return nil
	}
	od.Name = p.parseName()
	od.Interfaces = p.parseImplementInterfaces()
	//fmt.Println("@@", p.curToken) //if everything okay then current token is token.AT or token.LBRACE
	//current token is token.LBRACE
	od.Directives = p.parseDirectives()
	//loop current token is token.LPAREN
	p.nextToken()
	od.Fields = p.parseFieldsDefinition()
	fmt.Println("parseObjectDefinition->DONE")
	return od
}

func (p *Parser) parseFieldsDefinition() []*ast.FieldDefinition {

	fields := []*ast.FieldDefinition{}
	for !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF) {
		//starting with token.IDENT
		fd := p.parseFieldDefinition()
		if fd != nil {
			fields = append(fields, fd)
		}
	}
	return fields
}

func (p *Parser) parseDirectives() []*ast.Directive {

	fmt.Println("DDDDD", p.curToken, p.peekToken)
	dirs := make([]*ast.Directive, 0)

	for !p.curTokenIs(token.LBRACE) {
		directive := p.parseDirective()
		if directive != nil {
			dirs = append(dirs, directive)
		}
		fmt.Println("parseDirectives2", directive, p.curToken)
		if directive == nil {
			break
		}
	}
	fmt.Println("parseDirectives>>", dirs, len(dirs), "-->", p.curToken.Literal, p.peekToken.Literal)
	return dirs
}

func (p *Parser) parseDirective() *ast.Directive {

	if !p.expectToken(token.AT) {
		return nil
	}
	// if p.curTokenIs(token.LBRACE) {
	// 	return nil
	// }
	if !p.curTokenIs(token.IDENT) {
		return nil
	}
	directive := &ast.Directive{Kind: ast.DIRECTIVE, Token: p.curToken}
	directive.Name = p.parseName()

	directive.Arguments = p.Arguments()
	if !p.curTokenIs(token.LBRACE) {
		p.nextToken() //--> )
	}
	return directive
}

func (p *Parser) Arguments() []*ast.Argument {

	if !p.curTokenIs(token.LPAREN) {
		return nil
	}
	args := []*ast.Argument{}
	p.nextToken() //-> (
	for !p.curTokenIs(token.RPAREN) {

		arg := p.parseArgument()
		if arg != nil {
			args = append(args, arg)
		}
		if arg == nil {
			break
		}
	}
	return args
}

func (p *Parser) parseArgument() *ast.Argument {

	if !p.curTokenIs(token.IDENT) {
		return nil
	}
	arg := &ast.Argument{Kind: ast.ARGUMENT, Token: p.curToken}
	arg.Name = p.parseName()
	if p.curTokenIs(token.COLON) {
		p.nextToken()
	}
	arg.Value = p.parseValueLiteral()
	return arg
}

func (p *Parser) parseImplementInterfaces() []*ast.NamedType {

	namedSlc := []*ast.NamedType{}
	if !p.curTokenIs(token.IMPLEMENTS) {
		return nil
	}
	//fmt.Println("implements ==>", p.curToken.Literal)
	p.nextToken()
	for !p.curTokenIs(token.LBRACE) {
		named := p.parseNamed()
		if named != nil {
			namedSlc = append(namedSlc, named)
		}
		if p.curTokenIs(token.AMP) {
			p.nextToken()
		}

		if p.curTokenIs(token.AT) {
			//fmt.Println("@ FOUND", p.peekToken)
			break
		}
	}
	//fmt.Println("!!!", p.curToken, p.peekToken, namedSlc, len(namedSlc))
	return namedSlc
}

func (p *Parser) parseNamed() *ast.NamedType {

	//expecting current token is token.IDENT
	if !p.curTokenIs(token.IDENT) {
		return nil
	}
	named := &ast.NamedType{Kind: ast.NAMED_TYPE}
	named.Token = p.curToken
	named.Name = p.parseName()
	return named
}

func (p *Parser) parseFieldDefinition() *ast.FieldDefinition {

	//fmt.Println("fieldDefinition", p.curToken) //starting with token.IDENT
	fd := &ast.FieldDefinition{}
	fd.Kind = ast.FIELD_DEFINITION
	fd.Token = p.curToken

	fd.Description = p.parseDescription()
	fd.Name = p.parseName()

	fd.Arguments = p.parseArgumentDefinition()

	if !p.expectToken(token.COLON) {
		return nil
	}

	fd.Type = p.parseType()
	return fd
}

func (p *Parser) parseArgumentDefinition() []*ast.InputValueDefinition {

	//fmt.Println("parseArgumentDefinition", p.curToken, p.peekToken)
	args := []*ast.InputValueDefinition{}

	if !p.curTokenIs(token.LPAREN) {
		return nil
	}

	p.nextToken()
	for !p.curTokenIs(token.RPAREN) {

		//starting with token.IDENT
		ivd := p.parseInputValueDefinition()
		if ivd != nil {
			args = append(args, ivd)
		}
		if p.curTokenIs(token.COMMA) {
			p.nextToken()
		}
	}
	//last current token is token.RPAREN so next
	p.nextToken()
	return args
}

func (p *Parser) parseInputValueDefinition() *ast.InputValueDefinition {

	//fmt.Println("parseInputValueDefinition", p.curToken)
	inv := &ast.InputValueDefinition{Kind: ast.INPUT_VALUE_DEFINITION}
	inv.Token = p.curToken
	inv.Description = p.parseDescription()

	//current token.IDENT
	inv.Name = p.parseName()

	if !p.expectToken(token.COLON) {
		return nil
	}

	inv.Type = p.parseType()
	inv.DefaultValue = p.parseDefaultValue()
	inv.Directives = nil
	//last token is token.RPAREN = )
	return inv
}

func (p *Parser) parseDefaultValue() ast.Value {
	if !p.expectToken(token.ASSIGN) {
		return nil
	}
	return p.parseValueLiteral()
}

func (p *Parser) parseValueLiteral() ast.Value {

	//fmt.Println("parseValueLiteral", p.curToken)
	cToken := p.curToken
	var value ast.Value

	if cToken.Type == token.IDENT {

		if cToken.Literal == "true" {
			value = &ast.BooleanValue{Kind: ast.BOOLEAN_VALUE, Token: cToken, Value: true}

		} else if cToken.Literal == "false" {
			value = &ast.BooleanValue{Kind: ast.BOOLEAN_VALUE, Token: cToken, Value: false}

		} else if cToken.Literal == "null" {
			value = &ast.EnumValue{Kind: ast.ENUM_VALUE, Token: cToken, Value: cToken.Literal}

		} else if cToken.Literal != "null" {
			value = &ast.EnumValue{Kind: ast.ENUM_VALUE, Token: cToken, Value: cToken.Literal}
		}

	} else if cToken.Type == token.INT {
		value = &ast.IntValue{Kind: ast.INT_VALUE, Token: cToken, Value: cToken.Literal}

	} else if cToken.Type == token.FLOAT {
		value = &ast.FloatValue{Kind: ast.FLOAT_VALUE, Token: cToken, Value: cToken.Literal}

	} else if cToken.Type == token.STRING {
		value = &ast.StringValue{Kind: ast.STRING_VALUE, Token: cToken, Value: cToken.Literal}
		//value = p.parseStringLiteral()

	} else if cToken.Type == token.LBRACKET {
		//parseList

	} else if cToken.Type == token.LBRACE {
		//parseObject
	}
	p.nextToken()
	//fmt.Println("???", p.curToken) // )
	return value
}

// func (p *Parser) parseTypeSystemDefinition() ast.Node { //ast.TypeSystemDefinition

// 	fmt.Println("parseTypeSystemDefinition>", p.curToken)
// 	var keyWordToken token.Token
// 	if p.isDescription() {
// 		keyWordToken = p.peekToken
// 	}
// 	if !p.peekTokenIsKeyword() {
// 		p.peekError(token.IDENT)
// 		return nil
// 	}
// 	item, ok := p.tokenDefinitionFns[keyWordToken.Type]
// 	if !ok {
// 		return nil
// 	}
// 	return item()
// }

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
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

func (p *Parser) isDescription() bool {
	if p.curTokenIs(token.BLOCK_STRING) || p.curTokenIs(token.STRING) {
		return true
	}
	return false
}

func (p *Parser) expectToken(t token.TokenType) bool {
	if p.curTokenIs(t) {
		p.nextToken()
		return true
	}
	return false
}

// Converts a name lex token into a name parse node.
func (p *Parser) parseName() *ast.Name {

	//fmt.Println("parseName-->", p.curToken)
	if !p.curTokenIs(token.IDENT) {
		return nil
	}
	name := &ast.Name{Kind: ast.NAME, Token: p.curToken, Value: p.curToken.Literal}
	p.nextToken()
	return name
}

/**
 * Type :
 *   - NamedType
 *   - ListType
 *   - NonNullType
 */
func (p *Parser) parseType() (ttype ast.Type) {

	//fmt.Println("parseType", p.curToken, p.peekToken)
	cToken := p.curToken

	switch p.curToken.Type {
	case token.LBRACKET: //[
		p.nextToken()
		ttype = p.parseType()
		fallthrough

	case token.RBRACKET: //]
		p.nextToken()
		ttype = &ast.ListType{Kind: ast.LIST_TYPE, Token: cToken, Type: ttype}

	case token.IDENT, token.STRING:
		ttype = p.parseNamed()
	}

	// BANG must be executed
	if p.curTokenIs(token.BANG) {
		ttype = &ast.NonNullType{Kind: ast.NONNULL_TYPE, Token: p.curToken, Type: ttype}
		p.nextToken()
	}
	return ttype
}

func (p *Parser) parseDescription() *ast.StringValue {

	//fmt.Println("parseDescription", p.curToken, p.peekToken)
	if p.curTokenIs(token.STRING) || p.curTokenIs(token.BLOCK_STRING) || p.curTokenIs(token.HASH) {
		return p.parseStringLiteral()
	}
	return nil
}

func (p *Parser) parseStringLiteral() *ast.StringValue {
	cToken := p.curToken
	p.nextToken()
	return &ast.StringValue{Kind: ast.STRING_VALUE, Token: cToken, Value: cToken.Literal}
}
