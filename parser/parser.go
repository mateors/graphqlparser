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
		fmt.Println("...", p.errors)
		fmt.Println("!!!", p.curToken, p.peekToken)
		//p.nextToken()
	}
	return document
}

func (p *Parser) analyzeWhichDefinition() string {

	curToken := p.curToken.Type
	peekToken := p.peekToken.Type

	if p.isDescription() && peekToken == token.TYPE {
		return ast.OBJECT_DEFINITION

	} else if curToken == token.TYPE {
		return ast.OBJECT_DEFINITION

	} else if p.isDescription() && peekToken == token.INTERFACE {
		return ast.INTERFACE_DEFINITION

	} else if curToken == token.INTERFACE {
		return ast.INTERFACE_DEFINITION

	} else if p.isDescription() && peekToken == token.UNION {
		return ast.UNION_DEFINITION

	} else if curToken == token.UNION {
		return ast.UNION_DEFINITION

	} else if p.isDescription() && peekToken == token.ENUM {
		return ast.ENUM_DEFINITION

	} else if curToken == token.ENUM {
		return ast.ENUM_DEFINITION

	} else if p.isDescription() && peekToken == token.INPUT {
		return ast.INPUT_OBJECT_DEFINITION

	} else if curToken == token.INPUT {
		return ast.INPUT_OBJECT_DEFINITION

	} else if p.isDescription() && peekToken == token.SCALAR {
		return ast.SCALAR_DEFINITION

	} else if curToken == token.SCALAR {
		return ast.SCALAR_DEFINITION

	} else if p.isDescription() && peekToken == token.DIRECTIVE {
		return ast.DIRECTIVE_DEFINITION

	} else if curToken == token.DIRECTIVE {
		return ast.DIRECTIVE_DEFINITION

	} else if p.isDescription() && peekToken == token.SCHEMA {
		return ast.SCHEMA_DEFINITION

	} else if curToken == token.SCHEMA {
		return ast.SCHEMA_DEFINITION

	} else if curToken == token.QUERY {
		return ast.OPERATION_DEFINITION

	} else if curToken == token.MUTATION {
		return ast.OPERATION_DEFINITION

	} else if curToken == token.SUBSCRIPTION {
		return ast.OPERATION_DEFINITION

	} else if curToken == token.LBRACE {
		return ast.OPERATION_DEFINITION

	} else if curToken == token.FRAGMENT {
		return ast.FRAGMENT_DEFINITION

	}

	return ast.UNKNOWN
}

func (p *Parser) parseDocument() ast.Node { //ast.Definition

	//fmt.Println("parseDocument>", p.curToken.Type, p.analyzeWhichDefinition())

	switch p.analyzeWhichDefinition() {

	case ast.OBJECT_DEFINITION:
		return p.parseObjectDefinition()

	case ast.INTERFACE_DEFINITION:
		return p.parseInterfaceDefinition()

	case ast.UNION_DEFINITION:
		return p.parseUnionDefinition()

	case ast.ENUM_DEFINITION:
		return p.parseEnumDefinition()

	case ast.INPUT_OBJECT_DEFINITION:
		return p.parseInputObjectDefinition()

	case ast.SCALAR_DEFINITION:
		return p.parseScalarDefinition()

	case ast.OPERATION_DEFINITION:
		return p.parseOperationDefinition() //

	case ast.FRAGMENT_DEFINITION:
		return p.parseFragmentDefinition()

	case ast.DIRECTIVE_DEFINITION:
		return p.parseDirectiveDefinition()

	case ast.SCHEMA_DEFINITION:
		return p.parseSchemaDefinition()

	// 	fmt.Println("tokenDefinitionFns->", p.curToken.Type)
	// 	parseFunc := p.tokenDefinitionFns[p.curToken.Type]
	// 	return parseFunc()

	default:
		fmt.Println("unexpected", p.curToken.Type, p.curToken.Literal)
		return nil //&ast.OperationDefinition{}
	}

}

func (p *Parser) parseSchemaDefinition() ast.Node {

	//fmt.Println("parseSchemaDefinition:", p.curToken, p.peekToken)
	schema := &ast.SchemaDefinition{Kind: ast.SCHEMA_DEFINITION}
	schema.Token = p.curToken
	schema.Description = p.parseDescription()

	if !p.expectToken(token.SCHEMA) {
		return nil
	}
	schema.Directives = p.parseDirectives()
	schema.OperationTypes = p.parseRootOperationTypes()
	return schema
}

func (p *Parser) parseRootOperationTypes() []*ast.RootOperationTypeDefinition {

	//fmt.Println("parseRootOperationTypes", p.curToken, p.peekToken)
	roTypes := []*ast.RootOperationTypeDefinition{}

	if !p.expectToken(token.LBRACE) {
		return nil
	}

	for !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF) {

		roType := p.parseRootOperationTypeDefinition()
		if roType == nil {
			break
		}
		if roType != nil {
			roTypes = append(roTypes, roType)
		}
	}

	if p.curTokenIs(token.RBRACE) {
		p.nextToken()
	}
	return roTypes
}

func (p *Parser) parseRootOperationTypeDefinition() *ast.RootOperationTypeDefinition {

	//fmt.Println("parseRootOperationTypeDefinition", p.curToken, p.peekToken)
	cToken := p.curToken
	ivot := isValidOperationType(cToken.Literal)
	if !ivot {
		p.addError("rootOperationTypeDefinition operationType error!")
		//no return because we want to proceed but return nil at the end
	}

	roType := &ast.RootOperationTypeDefinition{Kind: ast.ROOT_OPERATION_TYPE_DEFINITION}
	roType.OperationType = tokenToOperationType(cToken)

	p.nextToken() //query | mutation | subscription
	roType.Token = cToken

	if !p.expectToken(token.COLON) {
		p.addError("rootOperationTypeDefinition colon missing error!")
		return nil
	}
	roType.NamedType = p.parseNamed()
	if !ivot {
		return nil //follow line 237 reason
	}
	return roType
}

func tokenToOperationType(cToken token.Token) string {

	switch cToken.Literal {
	case ast.QUERY:
		return ast.QUERY
	case ast.MUTATION:
		return ast.MUTATION
	case ast.SUBSCRIPTION:
		return ast.SUBSCRIPTION
	default:
		return "unknown"
	}
}

func isValidOperationType(operationType string) bool {
	if operationType == ast.QUERY || operationType == ast.MUTATION || operationType == ast.SUBSCRIPTION {
		return true
	}
	return false
}

func (p *Parser) parseDirectiveDefinition() ast.Node {

	//fmt.Println("parseDirectiveDefinition", p.curToken, p.peekToken)
	dird := &ast.DirectiveDefinition{Kind: ast.DIRECTIVE_DEFINITION}
	dird.Token = p.curToken
	dird.Description = p.parseDescription()

	if !p.expectToken(token.DIRECTIVE) {
		return nil
	}
	if !p.expectToken(token.AT) {
		return nil
	}

	dird.Name = p.parseName()
	dird.Arguments = p.parseArgumentDefinition() //?
	//fmt.Println(">>", dird.Name, dird.Arguments, p.curToken, p.peekToken)

	if p.curTokenIs(token.REPEATABLE) {
		p.nextToken()
	}

	if !p.expectToken(token.ON) {
		return nil
	}
	dird.Locations = p.parseLocations()
	return dird
}

func (p *Parser) parseLocations() []*ast.Name {

	if !p.curTokenIs(token.IDENT) {
		return nil
	}
	locations := []*ast.Name{}
	for {
		if p.curTokenIs(token.PIPE) {
			p.nextToken()
		}
		name := p.parseName()
		if name != nil {
			locations = append(locations, name)
		}
		if p.curTokenIs(token.PIPE) {
			p.nextToken()
		}
		if name == nil {
			break
		}
	}
	return locations
}

func (p *Parser) parseFragmentDefinition() ast.Node {

	fmt.Println("parseFragmentDefinitionSTART", p.curToken, p.peekToken)
	if !p.expectToken(token.FRAGMENT) {
		return nil
	}
	frg := &ast.FragmentDefinition{Kind: ast.FRAGMENT_DEFINITION}
	frg.Token = p.curToken
	frg.Operation = ""

	frg.FragmentName = p.parseName()
	frg.TypeCondition = p.parseTypeCondition()
	frg.Directives = p.parseDirectives()
	frg.SelectionSet = p.parseSelectionSet()
	fmt.Println("END-->", p.curToken, p.peekToken)
	return frg
}

func (p *Parser) parseOperationDefinition() ast.Node {

	fmt.Println("parseOperationDefinition START", p.curToken, p.peekToken)
	opDef := &ast.OperationDefinition{Kind: ast.OPERATION_DEFINITION}
	opDef.Token = p.curToken
	if p.curTokenIs(token.QUERY) {
		opDef.OperationType = ast.QUERY
		p.nextToken()
	}
	if p.curTokenIs(token.MUTATION) {
		opDef.OperationType = ast.MUTATION
		p.nextToken()
	}
	if p.curTokenIs(token.SUBSCRIPTION) {
		opDef.OperationType = ast.SUBSCRIPTION
		p.nextToken()
	}

	name := p.parseName()
	if name == nil {
		//if opDef.OperationType != "" {
		//p.addError("operationDefinition name error!")
		//}
		//log.Println("optional*-> operation name is missing")
	}
	opDef.Name = name
	opDef.VariablesDefinition = p.parseVariablesDefinition()
	opDef.Directives = p.parseDirectives()
	opDef.SelectionSet = p.parseSelectionSet()
	fmt.Println("END", p.curToken, p.peekToken)
	return opDef
}

func (p *Parser) parseSelectionSet() *ast.SelectionSet {

	if !p.curTokenIs(token.LBRACE) {
		return nil
	}
	selSet := &ast.SelectionSet{Kind: ast.SELECTION_SET}
	selSet.Token = p.curToken
	selSet.Selections = p.parseSelection()
	return selSet
}

func (p *Parser) parseSelection() []ast.Selection { //?

	//TODO
	if !p.expectToken(token.LBRACE) {
		return nil
	}
	selections := []ast.Selection{}

	for !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF) {

		field := p.parseField()
		if field != nil {
			selections = append(selections, field)
		}

		fragSpd := p.parseFragmentSpread()
		if fragSpd != nil {
			selections = append(selections, fragSpd)
		}

		inlineFrg := p.parseInlineFragment()
		if inlineFrg != nil {
			selections = append(selections, inlineFrg)
		}

		if field == nil && fragSpd == nil && inlineFrg == nil {
			fmt.Println("ALL ARE NIL SO BREAK NOW")
			break
		}
	}

	if p.curTokenIs(token.RBRACE) {
		p.nextToken() // }
	}
	return selections
}

func (p *Parser) parseInlineFragment() *ast.InlineFragment {

	cToken := p.curToken
	if p.curTokenIs(token.SPREAD) { //&& p.peekTokenIs(token.ON)
		if !p.expectToken(token.SPREAD) {
			return nil
		}
		inlineFrg := &ast.InlineFragment{Kind: ast.INLINE_FRAGMENT}
		inlineFrg.Token = cToken
		inlineFrg.TypeCondition = p.parseTypeCondition()
		inlineFrg.Directives = p.parseDirectives()
		inlineFrg.SelectionSet = p.parseSelectionSet()
		return inlineFrg
	}
	return nil
}

func (p *Parser) parseTypeCondition() *ast.NamedType {

	//fmt.Println("parseTypeCondition", p.curToken, p.peekToken)
	cToken := p.curToken
	if !p.expectToken(token.ON) {
		return nil
	}
	tc := &ast.NamedType{Kind: ast.NAMED_TYPE}
	tc.Token = cToken
	name := p.parseName()
	if name == nil {
		p.addError("parseTypeCondition name error!")
		return nil
	}
	tc.Name = name
	return tc
}

func (p *Parser) parseFragmentSpread() *ast.FragmentSpread {

	//fmt.Println("parseFragmentSpread", p.curToken, p.peekToken)
	if p.curTokenIs(token.SPREAD) && p.peekTokenIs(token.IDENT) {

		p.nextToken() //...
		frags := &ast.FragmentSpread{Kind: ast.FRAGMENT_SPREAD}
		frags.Token = p.curToken
		name := p.parseName()
		if name == nil {
			p.addError("parseFragmentSpread FragmentName error!")
			return nil //
		}
		frags.FragmentName = name
		frags.Directives = p.parseDirectives()
		return frags
	}
	return nil
}

func (p *Parser) parseField() *ast.Field {

	if !p.curTokenIs(token.IDENT) {
		return nil
	}
	field := &ast.Field{Kind: ast.FIELD}
	field.Token = p.curToken
	field.Alias = p.parseAlias()

	name := p.parseName()
	if name == nil {
		p.addError("parseField name error!")
		return nil
	}
	field.Name = name                    //mandatory
	field.Arguments = p.parseArguments() //?

	//fmt.Println("field.Arguments:", field.Arguments)
	field.Directives = p.parseDirectives()
	field.SelectionSet = p.parseSelectionSet()
	return field
}

func (p *Parser) parseAlias() *ast.Name {

	//fmt.Println("====?ALIAS", p.curToken, p.peekToken)
	if (p.curTokenIs(token.IDENT) || p.curTokenIsKeyword()) && p.peekTokenIs(token.COLON) {
		name := &ast.Name{Kind: ast.NAME, Token: p.curToken, Value: p.curToken.Literal}
		p.nextToken() //IDENT
		p.nextToken() //:
		return name
	}
	return nil
}

func (p *Parser) parseVariablesDefinition() []*ast.VariableDefinition {

	if !p.expectToken(token.LPAREN) {
		return nil
	}
	vars := []*ast.VariableDefinition{}

	for !p.curTokenIs(token.RPAREN) && !p.curTokenIs(token.EOF) {

		vard := p.parseVariableDefinition()
		if vard == nil {
			break
		}
		if vard.Variable == nil {
			vard = nil
		}
		if vard.Type == nil {
			vard = nil
		}
		if vard != nil {
			vars = append(vars, vard)
		}
		if p.curTokenIs(token.COMMA) {
			p.nextToken()
		}
	}
	if p.curTokenIs(token.RPAREN) {
		p.nextToken()
	}
	return vars
}

func (p *Parser) parseVariableDefinition() *ast.VariableDefinition {

	if !p.curTokenIs(token.DOLLAR) {
		return nil
	}
	svar := &ast.VariableDefinition{Kind: ast.VARIABLE_DEFINITION}
	svar.Token = p.curToken
	pvar := p.parseVariable()
	if pvar == nil {
		p.addError("parseVariableDefinition variable error!")
		return nil
	}
	svar.Variable = pvar

	if !p.expectToken(token.COLON) {
		p.addError("parseVariableDefinition colon missing error!")
		return nil
	}
	ttype := p.parseType()
	if ttype == nil {
		p.addError("parseVariableDefinition Type error!")
		return nil
		//p.nextToken() //if we still want to proceed
	}
	svar.Type = ttype
	svar.DefaultValue = p.parseDefaultValue()
	svar.Directives = p.parseDirectives()
	return svar
}

func (p *Parser) parseVariable() *ast.Variable {

	if !p.expectToken(token.DOLLAR) {
		return nil
	}
	v := &ast.Variable{Kind: ast.VARIABLE}
	v.Token = p.curToken
	name := p.parseName()
	if name == nil {
		return nil
	}
	v.Name = name
	return v
}

func (p *Parser) parseScalarDefinition() ast.Node {

	scd := &ast.ScalarDefinition{Kind: ast.SCALAR_DEFINITION}
	scd.Token = p.curToken
	scd.Description = p.parseDescription()

	if !p.expectToken(token.SCALAR) {
		return nil
	}
	name := p.parseName()
	if name == nil {
		p.addError("scalarDefinition name error!")
		return nil
	}
	scd.Name = name
	scd.Directives = p.parseDirectives()
	return scd
}

func (p *Parser) parseInputObjectDefinition() ast.Node {

	inObj := &ast.InputObjectDefinition{Kind: ast.INPUT_OBJECT_DEFINITION}
	inObj.Token = p.curToken
	inObj.Description = p.parseDescription()

	if !p.expectToken(token.INPUT) {
		return nil
	}
	name := p.parseName()
	if name == nil {
		p.addError("inputObjectDefinition name error!")
		return nil
	}
	inObj.Name = name
	inObj.Directives = p.parseDirectives()
	inObj.Fields = p.parseInputFieldsDefinition() //?
	return inObj
}

func (p *Parser) parseInputFieldsDefinition() []*ast.InputValueDefinition {

	fields := []*ast.InputValueDefinition{}
	if !p.curTokenIs(token.LBRACE) {
		return nil
	}

	p.nextToken()
	for !p.curTokenIs(token.RBRACE) {

		//starting with token.IDENT
		ivd := p.parseInputValueDefinition()
		if ivd == nil {
			break
		}
		fields = append(fields, ivd)
	}
	//last current token is token.RPAREN so next
	p.nextToken()
	return fields
}

func (p *Parser) parseEnumDefinition() ast.Node {

	ed := &ast.EnumDefinition{Kind: ast.ENUM_DEFINITION}
	ed.Token = p.curToken
	ed.Description = p.parseDescription()

	if !p.expectToken(token.ENUM) {
		return nil
	}

	name := p.parseName()
	if name == nil {
		p.addError("enumDefinition name error!")
	}
	ed.Name = name
	ed.Directives = p.parseDirectives()
	ed.Values = p.parseEnumValuesDefinition()
	return ed
}

func (p *Parser) parseEnumValuesDefinition() []*ast.EnumValueDefinition {

	evals := []*ast.EnumValueDefinition{}
	if !p.expectToken(token.LBRACE) {
		return nil
	}
	for !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF) {
		evd := p.parseEnumValueDefinition()
		if evd == nil {
			break
		}
		evals = append(evals, evd)
	}
	if p.curTokenIs(token.RBRACE) {
		p.nextToken()
	}
	return evals
}

func (p *Parser) parseEnumValueDefinition() *ast.EnumValueDefinition {

	evd := &ast.EnumValueDefinition{Kind: ast.ENUMVALUE_DEFINITION}
	evd.Token = p.curToken
	evd.Description = p.parseDescription()

	name := p.parseName()
	if name == nil {
		p.addError("enumValueDefinition name error!")
		return nil
	}
	evd.Name = name
	evd.Directives = p.parseDirectives()
	return evd
}

func (p *Parser) parseUnionDefinition() ast.Node { //?

	ud := &ast.UnionDefinition{Kind: ast.UNION_DEFINITION}
	ud.Token = p.curToken
	ud.Description = p.parseDescription()

	if !p.expectToken(token.UNION) {
		return nil
	}
	name := p.parseName()
	if name == nil {
		p.addError("unionDefinition name error!")
	}
	ud.Name = name
	ud.Directives = p.parseDirectives()
	ud.UnionMemberTypes = p.parseUnionMemberTypes()
	return ud
}

func (p *Parser) parseUnionMemberTypes() []*ast.NamedType {

	if !p.expectToken(token.ASSIGN) {
		return nil
	}

	namedSlc := []*ast.NamedType{}

	for {
		named := p.parseNamed()
		if named != nil {
			namedSlc = append(namedSlc, named)
		}
		if named == nil {
			break
		}
		if p.curTokenIs(token.PIPE) {
			p.nextToken()
		}
	}
	return namedSlc
}

func (p *Parser) parseInterfaceDefinition() ast.Node {

	//fmt.Println("parseInterfaceDefinition:", p.curToken)
	id := &ast.InterfaceDefinition{Kind: ast.INTERFACE_DEFINITION}
	id.Token = p.curToken
	id.Description = p.parseDescription()

	if !p.expectToken(token.INTERFACE) {
		return nil
	}

	name := p.parseName()
	if name == nil {
		p.addError("interfaceDefinition name error!")
	}
	id.Name = name

	id.Interfaces = p.parseImplementInterfaces()

	id.Directives = p.parseDirectives()

	fields := p.parseFieldsDefinition()
	if fields == nil {
		p.addError("interfaceDefinition fields parse error")
	}
	id.Fields = fields
	return id
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
		p.addError("objectDefinition name error!")
	}
	od.Name = name

	od.Interfaces = p.parseImplementInterfaces()
	//fmt.Println("@@", p.curToken) //if everything okay then current token is token.AT or token.LBRACE

	od.Directives = p.parseDirectives()

	fields := p.parseFieldsDefinition()
	if fields == nil {
		p.addError("objecDefinition fields parse error")
	}
	od.Fields = fields
	//fmt.Println("parseObjectDefinition->DONE", p.errors, len(p.errors))
	return od
}

func (p *Parser) parseFieldsDefinition() []*ast.FieldDefinition { //???? working not finished yet

	if !p.expectToken(token.LBRACE) { //expecting {
		return nil
	}
	fields := []*ast.FieldDefinition{}

	for !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF) {

		//starting with token.IDENT
		fd := p.parseFieldDefinition() //?

		if fd != nil {
			fields = append(fields, fd)
			if fd.Type == nil {
				fd = nil
			}
		}
		// if fd == nil {
		// 	break
		// }
		//fmt.Println(">>", fd, p.curToken, p.peekToken)
	}
	if p.curTokenIs(token.RBRACE) {
		p.nextToken()
	}
	return fields
}

func (p *Parser) parseDirectives() []*ast.Directive {

	if !p.curTokenIs(token.AT) {
		return nil
	}
	dirs := make([]*ast.Directive, 0)
	for !p.curTokenIs(token.LBRACE) && !p.curTokenIs(token.EOF) {

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
	return directive
}

func (p *Parser) parseArguments() []*ast.Argument { //?

	//fmt.Println("parseArguments>>", p.curToken, p.peekToken)
	if !p.curTokenIs(token.LPAREN) {
		return nil
	}
	args := []*ast.Argument{}
	p.nextToken() //-> (

	for !p.curTokenIs(token.RPAREN) {

		arg := p.parseArgument() //???
		if arg == nil {
			break
		}
		if arg.Name == nil {
			p.addError("parseArgument Name missing")
			arg = nil
		}
		if arg.Value == nil {
			p.addError("parseArgument Value missing")
			arg = nil
		}
		if arg != nil {
			args = append(args, arg)
		}
	}

	if p.curTokenIs(token.RPAREN) { //updated
		p.nextToken()
	}
	return args
}

func (p *Parser) parseArgument() *ast.Argument {

	if !p.curTokenIs(token.IDENT) && !p.curTokenIsKeyword() {
		return nil
	}
	arg := &ast.Argument{Kind: ast.ARGUMENT, Token: p.curToken}
	arg.Name = p.parseName()
	if p.curTokenIs(token.COLON) {
		p.nextToken()
	}

	val := p.parseValueLiteral()
	if val == nil {
		p.addError("parseArgument value error!")
		return nil
	}
	arg.Value = val
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

func (p *Parser) parseFieldDefinition() *ast.FieldDefinition {

	//fmt.Println("fieldDefinitionSTART", p.curToken, p.peekToken) //starting with token.IDENT
	if !p.curTokenIs(token.IDENT) && !p.curTokenIsKeyword() && !p.isDescription() {
		return nil
	}
	fd := &ast.FieldDefinition{}
	fd.Kind = ast.FIELD_DEFINITION
	fd.Token = p.curToken
	fd.Description = p.parseDescription()

	name := p.parseName()
	if name == nil {
		p.addError("parseFieldDefinition.parseName fd.Name missing!")
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
	fd.Directives = p.parseDirectives()
	return fd
}

func (p *Parser) parseArgumentDefinition() []*ast.InputValueDefinition {

	if !p.curTokenIs(token.LPAREN) {
		return nil
	}
	args := []*ast.InputValueDefinition{}
	p.nextToken() // (
	for !p.curTokenIs(token.RPAREN) && !p.curTokenIs(token.EOF) {

		//starting with token.IDENT
		ivd := p.parseInputValueDefinition()
		if ivd != nil {
			args = append(args, ivd)
		}
		if p.curTokenIs(token.COMMA) {
			p.nextToken()
		}
		if ivd == nil {
			break
		}
	}
	//last current token is token.RPAREN so next
	if p.curTokenIs(token.RPAREN) {
		p.nextToken() // )
	}
	return args
}

func (p *Parser) parseInputValueDefinition() *ast.InputValueDefinition {

	//fmt.Println("parseInputValueDefinition", p.curToken, p.peekToken)
	if !p.curTokenIs(token.IDENT) && !p.curTokenIsKeyword() && !p.isDescription() {
		//fmt.Println("**cnil", p.curToken, p.peekToken)
		return nil
	}

	inv := &ast.InputValueDefinition{Kind: ast.INPUT_VALUE_DEFINITION}
	inv.Token = p.curToken
	inv.Description = p.parseDescription()

	//current token.IDENT or keyword
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
	inv.Directives = p.parseDirectives()
	//last token is token.RPAREN = )
	//fmt.Println("*adir", p.curToken, p.peekToken)
	return inv
}

func (p *Parser) parseDefaultValue() ast.Value {
	if !p.expectToken(token.ASSIGN) {
		return nil
	}
	return p.parseValueLiteral()
}

func (p *Parser) parseValueLiteral() ast.Value {

	//TODO
	cToken := p.curToken
	var value ast.Value

	if cToken.Type == token.MINUS {
		p.nextToken() //-
		if p.curTokenIs(token.FLOAT) {
			value = &ast.FloatValue{Kind: ast.FLOAT_VALUE, Token: cToken, Value: "-" + p.curToken.Literal}
		}
		if p.curTokenIs(token.INT) {
			value = &ast.IntValue{Kind: ast.INT_VALUE, Token: cToken, Value: "-" + p.curToken.Literal}
		}

	} else if cToken.Type == token.IDENT {

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

	} else if cToken.Type == token.LBRACKET {
		value = p.parseList()

	} else if cToken.Type == token.DOLLAR {
		value = p.parseVariable()

	} else if cToken.Type == token.LBRACE {
		value = p.parseObjectValue()
	}
	p.nextToken()
	return value
}

func (p *Parser) parseObjectValue() *ast.ObjectValue {

	cToken := p.curToken
	if !p.expectToken(token.LBRACE) {
		return nil
	}
	obj := &ast.ObjectValue{Kind: ast.OBJECT_VALUE}
	obj.Token = cToken
	obj.Fields = p.parseObjectFields()
	return obj
}

func (p *Parser) parseObjectFields() []*ast.ObjectField {

	if !p.curTokenIs(token.IDENT) {
		return nil
	}

	objFlds := []*ast.ObjectField{}
	for !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF) {

		objfld := p.parseObjectField()
		if objfld != nil {
			objFlds = append(objFlds, objfld)
		}
		if p.curTokenIs(token.COMMA) {
			p.nextToken()
		}
		if objfld == nil {
			break
		}
	}
	return objFlds
}

func (p *Parser) parseObjectField() *ast.ObjectField {

	cToken := p.curToken
	if !p.curTokenIs(token.IDENT) {
		return nil
	}
	objfld := &ast.ObjectField{Kind: ast.OBJECT_FIELD}
	objfld.Token = cToken
	objfld.Name = p.parseName()
	if !p.expectToken(token.COLON) {
		p.addError("parseObjectField colon missing!")
		return nil
	}
	objfld.Value = p.parseValueLiteral()
	if objfld.Name == nil || objfld.Value == nil {
		p.addError("parseObjectField name or value nil")
		return nil
	}
	return objfld
}

func (p *Parser) parseList() *ast.ListValue {

	cToken := p.curToken
	if !p.expectToken(token.LBRACKET) {
		return nil
	}
	//fmt.Println("parseList", p.curToken, p.peekToken)
	listVal := &ast.ListValue{Kind: ast.LIST_VALUE}
	listVal.Token = cToken

	vals := []ast.Value{}

	for !p.curTokenIs(token.RBRACKET) && !p.curTokenIs(token.EOF) {

		val := p.parseValueLiteral()
		if val != nil {
			vals = append(vals, val)
		}
		if p.curTokenIs(token.COMMA) {
			//fmt.Println("COMMA FOUND")
			p.nextToken()
		}
	}

	if p.curTokenIs(token.RBRACKET) {
		//fmt.Println("RBRACKET FOUND")
		p.nextToken()
	}
	//fmt.Println("###", len(vals), p.curToken, p.peekToken)
	listVal.Values = vals
	return listVal
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
func (p *Parser) curTokenIsKeyword() bool {
	return token.IsKeyword(p.curToken.Literal)
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
	if p.curTokenIs(token.BLOCK_STRING) || p.curTokenIs(token.STRING) || p.curTokenIs(token.HASH) {
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

	if p.curTokenIs(token.IDENT) || p.curTokenIsKeyword() {
		name := &ast.Name{Kind: ast.NAME, Token: p.curToken, Value: p.curToken.Literal}
		p.nextToken()
		return name
	}
	return nil
	//fmt.Println(">", p.curToken, p.peekToken)
	// if !p.curTokenIs(token.IDENT) {
	// 	p.addError("parseName identifier missing")
	// 	return nil
	// }
	// name := &ast.Name{Kind: ast.NAME, Token: p.curToken, Value: p.curToken.Literal}
	// p.nextToken()
	// return name
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
