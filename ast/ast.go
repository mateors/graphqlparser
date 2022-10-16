package ast

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/mateors/graphqlparser/token"
)

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

// Implements Node
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
var _ Selection = (*InlineFragment)(nil)

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

type InlineFragment struct {
	//...TypeCondition[opt] Directives[opt] SelectionSet
	Kind          string //INLINE_FRAGMENT
	Token         token.Token
	TypeCondition *Name //on NamedType
	Directives    []*Directive
	SelectionSet  *SelectionSet
}

func (f *InlineFragment) TokenLiteral() string {
	return f.Token.Literal
}
func (f *InlineFragment) GetKind() string {
	return f.Kind
}
func (f *InlineFragment) GetSelectionSet() *SelectionSet {
	return f.SelectionSet
}

type Type interface {
	Node
	String() string
}

// Name implements Node & Type
var _ Type = (*NamedType)(nil)
var _ Type = (*ListType)(nil)
var _ Type = (*NonNullType)(nil)

// Named implements Node, Type
type NamedType struct {
	Kind  string //NAMED_TYPE
	Token token.Token
	Name  *Name
}

func (n *NamedType) TokenLiteral() string {
	return n.Name.Token.Literal
}

func (n *NamedType) GetKind() string {
	return n.Kind
}
func (n *NamedType) String() string {
	return n.GetKind()
}

type ListType struct {
	Kind  string //LIST_TYPE
	Token token.Token
	Type  Type
}

func (l *ListType) TokenLiteral() string {
	return l.Token.Literal
}
func (l *ListType) GetKind() string {
	return l.Kind
}
func (l *ListType) String() string {
	return l.GetKind()
}

type NonNullType struct {
	Kind  string //NONNULL_TYPE
	Token token.Token
	Type  Type
}

func (n *NonNullType) TokenLiteral() string {
	return n.Token.Literal
}
func (n *NonNullType) GetKind() string {
	return n.Kind
}
func (n *NonNullType) String() string {
	return n.GetKind()
}

type TypeDefinition interface {
	Node
	GetOperation() string
	GetVariableDefinitions() []*VariableDefinition
	GetSelectionSet() *SelectionSet
}

// var _ TypeDefinition = (*ScalarDefinition)(nil)
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

func (ob *ObjectDefinition) String() string {
	var out bytes.Buffer

	name := ob.Name.Value
	out.WriteString("type" + " ")
	out.WriteString(name + " " + "implements" + " ")

	var infcs string
	for _, inf := range ob.Interfaces {
		iname := inf.Value
		infcs += fmt.Sprintf("%s & ", iname)
	}
	infcs = strings.TrimRight(infcs, " & ")
	out.WriteString(infcs + " ")
	out.WriteString("{")

	//ob.Fields.String()
	for _, field := range ob.Fields {
		field.String() //??
	}

	out.WriteString("}")
	return out.String()
}

var _ Node = (*FieldDefinition)(nil)

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

func (fd *FieldDefinition) TokenLiteral() string {
	return fd.Token.Literal
}
func (fd *FieldDefinition) GetKind() string {
	return fd.Kind
}

// printer
func (fd *FieldDefinition) String() string {
	var out bytes.Buffer
	//Description[opt] Name ArgumentsDefinition[opt] : Type Directives[opt]
	//name: String
	//name: String!
	out.WriteString(fd.Name.Value)
	if len(fd.Arguments) > 0 {
		out.WriteString("(")
		var vals string
		for _, arg := range fd.Arguments {
			vals += fmt.Sprintf("%s, ", arg.String())
		}
		vals = strings.TrimRight(vals, ", ")
		out.WriteString(vals)
		out.WriteString(")")
	}
	return out.String()
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

func (iv *InputValueDefinition) TokenLiteral() string {
	return iv.Token.Literal
}
func (iv *InputValueDefinition) GetKind() string {
	return iv.Kind
}
func (iv *InputValueDefinition) String() string {
	var out bytes.Buffer
	name := fmt.Sprintf("%v", iv.Name.Value)
	ttype := fmt.Sprintf("%v", iv.Type)

	out.WriteString(name + ": " + ttype)
	return out.String()
}

type Value interface {
	Node
	GetValue() interface{}
}

// Ensure that all value types implements Value interface
var _ Value = (*Variable)(nil)

// var _ Value = (*IntValue)(nil)
// var _ Value = (*FloatValue)(nil)
var _ Value = (*StringValue)(nil)

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

type StringValue struct {
	Kind  string //VARIABLE
	Token token.Token
	Value string
}

func (s *StringValue) TokenLiteral() string {
	return s.Token.Literal
}

func (s *StringValue) GetKind() string {
	return s.Kind
}

func (s *StringValue) GetValue() interface{} {
	return s.Value
}
