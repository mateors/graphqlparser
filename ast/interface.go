package ast

import "github.com/mateors/graphqlparser/token"

type InterfaceDefinition struct {
	//Description[opt] interface Name ImplementsInterfaces[opt] Directives[opt] FieldsDefinition
	//Description[opt] interface Name ImplementsInterfaces[opt] Directives[opt]
	Kind        string //INTERFACE_DEFINITION
	Token       token.Token
	Description *StringValue
	Name        *Name
	Interfaces  []*NamedType
	Directives  []*Directive
	Fields      []*FieldDefinition
}

var _ TypeDefinition = (*InterfaceDefinition)(nil)

type TypeSystemDefinition interface {
	Node
	GetOperation() string
	GetVariableDefinitions() []*VariableDefinition
	GetSelectionSet() *SelectionSet
}

//var _ TypeSystemDefinition = (*SchemaDefinition)(nil)
var _ TypeSystemDefinition = (TypeDefinition)(nil)

//var _ TypeSystemDefinition = (*TypeExtensionDefinition)(nil)
//var _ TypeSystemDefinition = (*DirectiveDefinition)(nil)

func (i *InterfaceDefinition) TokenLiteral() string {
	return i.Token.Literal
}
func (i *InterfaceDefinition) GetKind() string {
	return i.Kind
}
func (ob *InterfaceDefinition) GetOperation() string {
	return ""
}
func (i *InterfaceDefinition) GetVariableDefinitions() []*VariableDefinition {
	return []*VariableDefinition{}
}
func (i *InterfaceDefinition) GetSelectionSet() *SelectionSet {
	return &SelectionSet{}
}
