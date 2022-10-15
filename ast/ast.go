package ast

import "github.com/mateors/graphqlparser/token"

type Node interface {
	TokenLiteral() string
	GetKind() string
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
	Kind        string
	Definitions []Definition
}

func (d *Document) TokenLiteral() string {

	if len(d.Definitions) > 0 {
		return d.Definitions[0].TokenLiteral()
	}
	return ""
}

func (d *Document) GetKind() string {
	return d.Kind
}

type OperationDefinition struct {
	//OperationType Name[opt] VariablesDefinition[opt] Directives[opt] SelectionSet
	Kind                string
	OperationType       string
	Name                *Name
	VariablesDefinition []*VariableDefinition
	Directives          []*Directive
	SelectionSet        *SelectionSet
}

type Name struct {
	Kind  string
	Token token.Token
	Value string
}

type VariableDefinition struct {
	//Variable : Type DefaultValue[opt] Directives[opt]
	Kind         string
	Token        token.Token
	Variable     *Variable
	Type         Type
	DefaultValue Value
	Directives   []*Directive
}

type Directive struct {
}

type SelectionSet struct {
}

type Variable struct {
}

type Type interface {
}

type Value interface {
}
