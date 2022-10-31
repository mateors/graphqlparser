package parser

import (
	"errors"
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

	name := p.parseName()
	if name == nil {
		//fmt.Println("name nil so continue")
		p.addError("objectDefinition name error!")
	}
	od.Name = name

	od.Interfaces = p.parseImplementInterfaces()
	//fmt.Println("@@", p.curToken) //if everything okay then current token is token.AT or token.LBRACE
	//current token is token.LBRACE

	od.Directives = p.parseDirectives()
	//loop current token is token.LPAREN
	p.nextToken()

	fields := p.parseFieldsDefinition()
	if fields == nil {
		p.addError("objecDefinition fields parse error")
	}
	od.Fields = fields
	fmt.Println("parseObjectDefinition->DONE", p.errors, len(p.errors))
	return od
}

func (p *Parser) parseFieldsDefinition() []*ast.FieldDefinition { //???? working not finished yet

	//check if return nil
	fields := []*ast.FieldDefinition{}

	for !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF) {

		//starting with token.IDENT
		fd := p.parseFieldDefinition()

		if fd != nil {

			fields = append(fields, fd)
			//fmt.Println("<2><2>", fd, p.curToken, p.peekToken) //fd.Type
			if fd.Type == nil {
				fd = nil
			}
		}

	}
	//fmt.Println("@@", p.curToken, p.peekToken, len(fields))
	return fields
}

func (p *Parser) parseDirectives() []*ast.Directive {

	if !p.curTokenIs(token.AT) {
		return nil
	}
	dirs := make([]*ast.Directive, 0)
	for !p.curTokenIs(token.LBRACE) {

		directive := p.parseDirective()
		if directive != nil {
			dirs = append(dirs, directive)
		}
		if directive == nil {
			break
		}
	}
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
		p.tokenError(token.AT)
		return nil
	}

	directive := &ast.Directive{Kind: ast.DIRECTIVE, Token: p.curToken}
	directive.Name = p.parseName()
	directive.Arguments = p.parseArguments()

	if !p.curTokenIs(token.LBRACE) {
		p.nextToken() //--> )
	}
	return directive
}

func (p *Parser) parseArguments() []*ast.Argument {

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
	if p.curTokenIs(token.COMMA) {
		p.nextToken()
	}
	return arg
}

func (p *Parser) parseImplementInterfaces() []*ast.NamedType {

	if !p.curTokenIs(token.IMPLEMENTS) {
		return nil
	}

	namedSlc := []*ast.NamedType{}
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
			break
		}
	}
	return namedSlc
}

func (p *Parser) parseNamed() *ast.NamedType {

	//expecting current token is token.IDENT
	if !p.curTokenIs(token.IDENT) {
		return nil
	}
	named := &ast.NamedType{Kind: ast.NAMED_TYPE}
	named.Token = p.curToken
	name := p.parseName()
	if name == nil {
		return nil
	}
	named.Name = name
	return named
}

func (p *Parser) parseFieldDefinition() *ast.FieldDefinition { //??

	//fmt.Println("fieldDefinition", p.curToken) //starting with token.IDENT
	fd := &ast.FieldDefinition{}
	fd.Kind = ast.FIELD_DEFINITION
	fd.Token = p.curToken
	fd.Description = p.parseDescription()

	name := p.parseName()
	if name == nil {
		p.addError("parseFieldDefinition.parseName type missing")
		return nil
	}

	fd.Name = name
	fd.Arguments = p.parseArgumentDefinition()

	if !p.expectToken(token.COLON) {
		p.tokenError(token.COLON)
		return nil
	}

	ptype := p.parseType()
	if ptype == nil {
		p.addError("parseFieldDefinition.parseType error")
		return nil
	}
	fd.Type = ptype
	fd.Directives = nil
	return fd
}

func (p *Parser) parseArgumentDefinition() []*ast.InputValueDefinition {

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

	inv := &ast.InputValueDefinition{Kind: ast.INPUT_VALUE_DEFINITION}
	inv.Token = p.curToken
	inv.Description = p.parseDescription()

	//current token.IDENT
	name := p.parseName()
	if name == nil {
		return nil
	}
	inv.Name = name

	if !p.expectToken(token.COLON) {
		p.tokenError(token.COLON)
		return nil
	}

	ptype := p.parseType()
	if ptype == nil {
		p.addError("parseInputValueDefinition parseType error")
		return nil
	}

	inv.Type = ptype
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

func (p *Parser) tokenError(tokType token.TokenType) error {
	msg := fmt.Sprintf("expected token %v got %v at line %d column %d", token.Name(tokType), token.Name(p.curToken.Type), p.curToken.Line, p.curToken.Start)
	p.errors = append(p.errors, msg)
	return errors.New(msg)
}

func (p *Parser) addError(msg string) error {
	amsg := fmt.Sprintf("%s token: %s,%s", msg, token.Name(p.curToken.Type), token.Name(p.peekToken.Type))
	p.errors = append(p.errors, amsg)
	return errors.New(msg)
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

	if !p.curTokenIs(token.IDENT) {
		p.addError("parseName identifier missing")
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
//Instead of error we return nil
func (p *Parser) parseType() (ttype ast.Type) { //????

	cToken := p.curToken
	switch p.curToken.Type {
	case token.LBRACKET: //[
		p.nextToken()
		//if p.curToken.Type == token.RBRACKET{}
		ttype = p.parseType()
		fallthrough

	case token.RBRACKET: //]
		if ttype != nil {
			p.nextToken()
			ttype = &ast.ListType{Kind: ast.LIST_TYPE, Token: cToken, Type: ttype}
		}

	case token.IDENT, token.STRING:
		ttype = p.parseNamed()
	}

	// BANG must be executed
	if ttype == nil {
		p.nextToken()
	}

	if p.curTokenIs(token.BANG) && ttype != nil {
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
