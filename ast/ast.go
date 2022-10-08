package ast

type Node interface {
	TokenLiteral() string
}

type Definition interface { //Statement | Definition | SchemaDefinition | TypeSystems | GraphQLObjectType
	Node
	defNode()
}

type Expression interface {
	Node
	expNode()
}

type Document struct { //Program | Document
	Definitions []Definition
}

func (d *Document) TokenLiteral() string {

	if len(d.Definitions) > 0 {
		return d.Definitions[0].TokenLiteral()
	}
	return ""
}
