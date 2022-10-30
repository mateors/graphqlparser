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

	name, err := p.parseName()
	if err != nil {
		fmt.Println(err)
	}
	od.Name = name

	infs, err := p.parseImplementInterfaces()
	if err != nil {
		fmt.Println(err)
	}
	od.Interfaces = infs
	//fmt.Println("@@", p.curToken) //if everything okay then current token is token.AT or token.LBRACE
	//current token is token.LBRACE

	dirs, err := p.parseDirectives()
	if err != nil {
		fmt.Println(err)
	}
	od.Directives = dirs
	//loop current token is token.LPAREN
	p.nextToken()

	fields, err := p.parseFieldsDefinition()
	if err != nil {
		fmt.Println("parseFieldsDefinition->", err)
	}
	od.Fields = fields
	fmt.Println("parseObjectDefinition->DONE", p.errors, od.Fields)
	return od
}

func (p *Parser) parseFieldsDefinition() ([]*ast.FieldDefinition, error) { //???? working not finished yet

	fields := []*ast.FieldDefinition{}
	for !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF) {
		//starting with token.IDENT
		fd, err := p.parseFieldDefinition()

		fmt.Println("<2><2>", fd.Type, p.curToken, p.peekToken)
		if fd.Type == nil {
			fmt.Println("EXIT", p.curToken, p.peekToken)
			//os.Exit(2)
			//return nil, errors.New("<<type nil error<<")
			fd = nil
		}
		if err != nil {
			fmt.Println(">>>", err)
			//return nil, err
		}

		//fmt.Println(">>>>>>>>>>", fd.Name, err)
		if fd != nil {
			fields = append(fields, fd)
		}
		//if fd == nil {
		//fmt.Println("BREAK...", p.curToken.Literal)
		//break
		//}
	}
	fmt.Println("@@", p.curToken, p.peekToken)
	return fields, nil
}

func (p *Parser) parseDirectives() ([]*ast.Directive, error) {

	if !p.curTokenIs(token.AT) {
		err := p.tokenError(token.AT)
		return nil, err
	}
	dirs := make([]*ast.Directive, 0)
	for !p.curTokenIs(token.LBRACE) {
		directive, err := p.parseDirective()
		if err != nil {
			return nil, err
		}
		if directive != nil {
			dirs = append(dirs, directive)
		}
		if directive == nil {
			break
		}
	}
	return dirs, nil
}

func (p *Parser) parseDirective() (*ast.Directive, error) {

	if !p.expectToken(token.AT) {
		err := p.tokenError(token.AT)
		return nil, err
	}
	// if p.curTokenIs(token.LBRACE) {
	// 	return nil
	// }
	if !p.curTokenIs(token.IDENT) {
		err := p.tokenError(token.AT)
		return nil, err
	}
	directive := &ast.Directive{Kind: ast.DIRECTIVE, Token: p.curToken}

	name, err := p.parseName()
	if err != nil {
		return nil, err
	}
	directive.Name = name

	args, err := p.parseArguments()
	if err != nil {
		return nil, err
	}

	directive.Arguments = args
	if !p.curTokenIs(token.LBRACE) {
		p.nextToken() //--> )
	}
	return directive, nil
}

func (p *Parser) parseArguments() ([]*ast.Argument, error) {

	if !p.curTokenIs(token.LPAREN) {
		err := p.tokenError(token.LPAREN)
		return nil, err
	}
	args := []*ast.Argument{}
	p.nextToken() //-> (
	for !p.curTokenIs(token.RPAREN) {

		arg, err := p.parseArgument()
		if err != nil {
			return nil, err
		}

		if arg != nil {
			args = append(args, arg)
		}
		if arg == nil {
			break
		}
	}
	return args, nil
}

func (p *Parser) parseArgument() (*ast.Argument, error) {

	if !p.curTokenIs(token.IDENT) {
		err := p.tokenError(token.IDENT)
		return nil, err
	}
	arg := &ast.Argument{Kind: ast.ARGUMENT, Token: p.curToken}
	name, err := p.parseName()
	if err != nil {
		return nil, err
	}
	arg.Name = name
	if p.curTokenIs(token.COLON) {
		p.nextToken()
	}
	arg.Value = p.parseValueLiteral()
	if p.curTokenIs(token.COMMA) {
		p.nextToken()
	}
	return arg, nil
}

func (p *Parser) parseImplementInterfaces() ([]*ast.NamedType, error) {

	if !p.curTokenIs(token.IMPLEMENTS) {
		err := p.tokenError(token.IMPLEMENTS)
		return nil, err
	}
	namedSlc := []*ast.NamedType{}
	//fmt.Println("implements ==>", p.curToken.Literal)
	p.nextToken()
	for !p.curTokenIs(token.LBRACE) {
		named, err := p.parseNamed()
		if err != nil {
			return nil, err
		}
		//if named != nil {
		namedSlc = append(namedSlc, named)
		//}

		if p.curTokenIs(token.AMP) {
			p.nextToken()
		}
		if p.curTokenIs(token.AT) {
			break
		}
	}
	//fmt.Println("!!!", p.curToken, p.peekToken, namedSlc, len(namedSlc))
	return namedSlc, nil
}

func (p *Parser) parseNamed() (*ast.NamedType, error) {

	//expecting current token is token.IDENT
	if !p.curTokenIs(token.IDENT) {
		err := p.tokenError(token.IDENT)
		return nil, err
	}
	named := &ast.NamedType{Kind: ast.NAMED_TYPE}
	named.Token = p.curToken
	name, err := p.parseName()
	if err != nil {
		return nil, err
	}
	named.Name = name
	return named, nil
}

func (p *Parser) parseFieldDefinition() (*ast.FieldDefinition, error) { //??

	//fmt.Println("fieldDefinition", p.curToken) //starting with token.IDENT
	var err error
	fd := &ast.FieldDefinition{}
	fd.Kind = ast.FIELD_DEFINITION
	fd.Token = p.curToken

	fd.Description = p.parseDescription()

	name, err := p.parseName()
	if err != nil {
		err = p.addError("parseFieldDefinition.parseName type missing")
		return nil, err
	}
	fd.Name = name
	//fmt.Println("field-->", name)

	argd, err := p.parseArgumentDefinition() //optional
	if err != nil {
		fmt.Println("parseArgumentDefinition>", err)
	}
	fd.Arguments = argd

	if !p.expectToken(token.COLON) {
		err = p.tokenError(token.COLON)
		return nil, err
	}

	ptype, err := p.parseType()
	if err != nil {
		err = p.addError("parseFieldDefinition.parseType error")
		return nil, err
	}
	fmt.Println("<1><1>", ptype, p.curToken, p.peekToken)
	// if ptype == nil {
	// 	fmt.Println(fd.Name, "nil so return", p.curToken, p.peekToken)
	// 	return nil, errors.New("type nil error")
	// }

	fd.Type = ptype
	fd.Directives = nil
	//fmt.Println("------>", err, fd.Name, ptype, "**", p.curToken, p.peekToken)
	return fd, nil
}

func (p *Parser) parseArgumentDefinition() ([]*ast.InputValueDefinition, error) {

	//fmt.Println("parseArgumentDefinition", p.curToken, p.peekToken)
	args := []*ast.InputValueDefinition{}
	//var err error
	if !p.curTokenIs(token.LPAREN) {
		return nil, nil
	}

	p.nextToken()
	for !p.curTokenIs(token.RPAREN) {

		//starting with token.IDENT
		ivd, err := p.parseInputValueDefinition()
		if err != nil {
			return nil, err
		}
		//if ivd != nil {
		args = append(args, ivd)
		//}

		if p.curTokenIs(token.COMMA) {
			p.nextToken()
		}
	}
	//last current token is token.RPAREN so next
	p.nextToken()
	return args, nil
}

func (p *Parser) parseInputValueDefinition() (*ast.InputValueDefinition, error) {

	//fmt.Println("parseInputValueDefinition", p.curToken)
	inv := &ast.InputValueDefinition{Kind: ast.INPUT_VALUE_DEFINITION}
	inv.Token = p.curToken
	inv.Description = p.parseDescription()

	//current token.IDENT
	name, err := p.parseName()
	if err != nil {
		return nil, err
	}
	inv.Name = name

	if !p.expectToken(token.COLON) {
		err := p.tokenError(token.COLON)
		return nil, err
	}

	ptype, err := p.parseType()
	if err != nil {
		err = p.addError("parseInputValueDefinition parseType error")
		return nil, err
	}

	inv.Type = ptype
	inv.DefaultValue = p.parseDefaultValue()
	inv.Directives = nil
	//last token is token.RPAREN = )
	return inv, nil
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
func (p *Parser) parseName() (*ast.Name, error) {

	if !p.curTokenIs(token.IDENT) {
		//err := p.tokenError(token.IDENT)
		err := p.addError("parseName identifier missing")
		return nil, err
	}
	name := &ast.Name{Kind: ast.NAME, Token: p.curToken, Value: p.curToken.Literal}
	p.nextToken()
	return name, nil
}

/**
 * Type :
 *   - NamedType
 *   - ListType
 *   - NonNullType
 */
func (p *Parser) parseType() (ttype ast.Type, err error) { //????

	//fmt.Println("parseType", p.curToken, p.peekToken)
	cToken := p.curToken

	switch p.curToken.Type {
	case token.LBRACKET: //[
		p.nextToken()
		if p.curToken.Type == token.RBRACKET {
			//fmt.Println("###", p.curToken.Literal) //expecting IDENT
			//return nil, errors.New("--missing type--")
		}
		ttype, err = p.parseType()
		//fmt.Println("2###", p.curToken.Literal, ttype, err)
		if err != nil {
			//fmt.Println(">>", err)
			return nil, err
		}

		fallthrough

	case token.RBRACKET: //]

		fmt.Println("~~", p.curToken, p.peekToken, ttype)
		if ttype != nil {
			p.nextToken()
			ttype = &ast.ListType{Kind: ast.LIST_TYPE, Token: cToken, Type: ttype}
		}
		if ttype == nil {
			fmt.Println("==========NIL========")
		}

	case token.IDENT, token.STRING:
		ttype, err = p.parseNamed()
		if err != nil {
			msg := "type identifier missing"
			fmt.Println(">>", msg)
			p.addError(msg)
			return nil, errors.New(msg)
		}
	}

	// BANG must be executed
	//fmt.Println("1~~~~", p.curToken, p.peekToken)

	if ttype == nil {
		//fmt.Println("nil so next", p.curToken, p.peekToken)
		p.nextToken()
	}

	if p.curTokenIs(token.BANG) && ttype != nil {
		fmt.Println("INSIDE BANG")
		ttype = &ast.NonNullType{Kind: ast.NONNULL_TYPE, Token: p.curToken, Type: ttype}
		p.nextToken()
	}
	//fmt.Println("2~~~~~~~~~~", p.curToken, p.peekToken, ttype)
	//fmt.Println()
	return ttype, nil
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
