package ast

import "github.com/mateors/graphqlparser/token"

type Node interface {
	TokenLiteral() string
	GetKind() string
}

type Definition interface { //Statement | Definition | SchemaDefinition | TypeSystems | GraphQLObjectType
	Node
	//defNode()
	GetOperation() string
	GetVariableDefinitions() []*VariableDefinition
	GetSelectionSet() *SelectionSet
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
	OperationType       string //query | mutation | subscription
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
	Kind      string
	Token     token.Token
	Name      *Name
	Arguments []*Argument
}

type Argument struct {
	Kind  string
	Token token.Token
	Name  *Name
	Value Value
}

//Implements Node
type SelectionSet struct {
	Kind       string //
	Token      token.Token
	Selections []Selection
}

type Selection interface {
	Node
	GetSelectionSet() *SelectionSet
}

// Ensure that all definition types implements Selection interface
var _ Selection = (*Field)(nil)
var _ Selection = (*FragmentSpread)(nil)

//var _ Selection = (*InlineFragment)(nil)

type Field struct {
	//Alias[opt] Name Arguments[opt] Directives[opt] SelectionSet[opt]
	Kind         string //FIELD
	Token        token.Token
	Alias        *Name
	Name         *Name
	Arguments    []*Argument
	Directives   []*Directive
	SelectionSet *SelectionSet
}

func (f *Field) TokenLiteral() string {
	return f.Token.Literal
}
func (f *Field) GetKind() string {
	return f.Kind
}
func (f *Field) GetSelectionSet() *SelectionSet {
	return f.SelectionSet
}

type FragmentSpread struct {
	//...FragmentName Directives[opt]
	Kind         string //FRAGMENT_SPREAD
	Token        token.Token
	FragmentName *Name
	Directives   []*Directive
}

func (fs *FragmentSpread) TokenLiteral() string {
	return fs.Token.Literal
}
func (fs *FragmentSpread) GetKind() string {
	return fs.Kind
}
func (fs *FragmentSpread) GetSelectionSet() *SelectionSet {
	return nil
}

type Type interface {
	Node
	String() string
}

// Name implements Node & Type
var _ Type = (*Name)(nil)

func (n *Name) TokenLiteral() string {
	return n.Token.Literal
}

func (n *Name) GetKind() string {
	return n.Kind
}
func (n *Name) String() string {
	return n.Value
}

type TypeDefinition interface {
	Node
	GetOperation() string
	GetVariableDefinitions() []*VariableDefinition
	GetSelectionSet() *SelectionSet
}

//var _ TypeDefinition = (*ScalarDefinition)(nil)
var _ TypeDefinition = (*ObjectDefinition)(nil)

//var _ TypeDefinition = (*InterfaceDefinition)(nil)
//var _ TypeDefinition = (*UnionDefinition)(nil)
//var _ TypeDefinition = (*EnumDefinition)(nil)
//var _ TypeDefinition = (*InputObjectDefinition)(nil)

type ObjectDefinition struct {
	//Description[opt] type Name ImplementsInterfaces[opt] Directives[opt] { FieldsDefinition }
	Kind        string //OBJECT_DEFINITION
	Token       token.Token
	Description string
	Name        *Name
	Interfaces  []*Name //NAMED
	Directives  []*Directive
	Fields      []*FieldDefinition
}

func (ob *ObjectDefinition) TokenLiteral() string {
	return ob.Token.Literal
}

func (ob *ObjectDefinition) GetKind() string {
	return ob.Kind
}

func (ob *ObjectDefinition) GetOperation() string {
	return ""
}

func (ob *ObjectDefinition) GetVariableDefinitions() []*VariableDefinition {
	return []*VariableDefinition{}
}

func (ob *ObjectDefinition) GetSelectionSet() *SelectionSet {
	return &SelectionSet{}
}

type FieldDefinition struct {
	//Description[opt] Name ArgumentsDefinition[opt] : Type Directives[opt]
	//ArgumentsDefinition: ( InputValueDefinition[list] )
	Kind        string //FIELD_DEFINITION
	Token       token.Token
	Description string
	Name        *Name
	Arguments   []*InputValueDefinition
	Type        Type
	Directives  []*Directive
}

type InputValueDefinition struct {
	//Description[opt] Name : Type DefaultValue[opt] Directives[opt]
	Kind         string //INPUT_VALUE_DEFINITION
	Token        token.Token
	Description  string
	Name         *Name
	Type         Type
	DefaultValue Value
	Directives   []*Directive
}

type Value interface {
	Node
	GetValue() interface{}
}

// Ensure that all value types implements Value interface
var _ Value = (*Variable)(nil)

//var _ Value = (*IntValue)(nil)
// var _ Value = (*FloatValue)(nil)
// var _ Value = (*StringValue)(nil)
// var _ Value = (*BooleanValue)(nil)
// var _ Value = (*EnumValue)(nil)
// var _ Value = (*ListValue)(nil)
// var _ Value = (*ObjectValue)(nil)

// Variable implements Node, Value
type Variable struct {
	Kind  string //VARIABLE
	Token token.Token
	Name  *Name
}

func (v *Variable) TokenLiteral() string {
	return v.Token.Literal
}

func (v *Variable) GetKind() string {
	return v.Kind
}

func (v *Variable) GetValue() interface{} {
	return v.Name
}
