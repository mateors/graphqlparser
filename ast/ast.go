package ast

type Node interface {
	TokenLiteral() string
}

type Statement interface { //Definition | SchemaDefinition | TypeSystems | GraphQLObjectType
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct { //Document
	Statements []Statement
}
