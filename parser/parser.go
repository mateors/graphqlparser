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
