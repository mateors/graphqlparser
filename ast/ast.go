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

var _ Node = (*Name)(nil)

type Name struct {
	Kind  string
	Token token.Token
	Value string
}

func (n *Name) TokenLiteral() string {
	return n.Token.Literal
}
func (n *Name) GetKind() string {
	return n.Kind
}
func (n *Name) String() string {
	return n.Value
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

var _ Node = (*Directive)(nil)

type Directive struct {
	Kind      string
	Token     token.Token
	Name      *Name
	Arguments []*Argument
}

func (d *Directive) TokenLiteral() string {
	return d.Token.Literal
}
func (d *Directive) GetKind() string {
	return d.Kind
}
func (d *Directive) String() string {
	var out bytes.Buffer
	//@Name Arguments[opt]
	var args string
	for _, arg := range d.Arguments {
		args += fmt.Sprintf("%s, ", arg.String())
	}
	args = strings.TrimRight(args, ", ")
	out.WriteString(" @" + d.Name.Value + "(" + args + ")")
	return out.String()
}

var _ Node = (*Argument)(nil)

type Argument struct {
	Kind  string
	Token token.Token
	Name  *Name
	Value Value
}

func (a *Argument) TokenLiteral() string {
	return a.Token.Literal
}
func (a *Argument) GetKind() string {
	return a.Kind
}
func (a *Argument) String() string {
	return fmt.Sprintf("%s: %v", a.Name.Value, a.Value)
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
	return n.Name.Value
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
	return fmt.Sprintf("[%s]", l.Type.String())
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
	return fmt.Sprintf("%s!", n.Type.String())
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
		out.WriteString(": " + fd.Type.String())
	}
	out.WriteString(": " + fd.Type.String())
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

	if iv.DefaultValue != nil {
		defaultValue := fmt.Sprintf("%v", iv.DefaultValue.GetValue())
		if len(defaultValue) > 0 {
			out.WriteString(fmt.Sprintf(" = %s", defaultValue))
		}
	}

	directives := []string{}
	for _, directive := range iv.Directives {
		directives = append(directives, fmt.Sprintf("%v", directive.String()))
	}

	if len(directives) > 0 {
		var dstr string
		for _, str := range directives {
			dstr += fmt.Sprintf("%s ", str)
		}
		dstr = strings.TrimRight(dstr, " ")
		out.WriteString(dstr)
	}
	return out.String()
}

type Value interface {
	Node
	GetValue() interface{}
}

// Ensure that all value types implements Value interface
var _ Value = (*Variable)(nil)
var _ Value = (*IntValue)(nil)
var _ Value = (*FloatValue)(nil)
var _ Value = (*StringValue)(nil)
var _ Value = (*BooleanValue)(nil)
var _ Value = (*EnumValue)(nil)
var _ Value = (*ListValue)(nil)
var _ Value = (*ObjectValue)(nil)

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
	return v.Name.Value
}
func (v *Variable) String() string {
	return fmt.Sprintf("$%s", v.GetValue())
}

var _ Node = (*StringValue)(nil)

type StringValue struct {
	Kind  string //STRING_VALUE
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
	//fmt.Println("kind:", s.Kind)
	if s.Kind == STRING_VALUE {
		return s.String()
	}
	return s.Value
}
func (s *StringValue) String() string {
	var vals string = fmt.Sprintf(`"%s"`, s.Value)
	return vals
}

var _ Node = (*BooleanValue)(nil)

type BooleanValue struct {
	Kind  string //BOOLEAN_VALUE
	Token token.Token
	Value bool
}

func (b *BooleanValue) TokenLiteral() string {
	return b.Token.Literal
}

func (b *BooleanValue) GetKind() string {
	return b.Kind
}

func (b *BooleanValue) GetValue() interface{} {
	return b.Value
}
func (s *BooleanValue) String() string {
	return fmt.Sprint(s.Value)
}

var _ Node = (*IntValue)(nil)
var _ Value = (*IntValue)(nil)

type IntValue struct {
	Kind  string //IntValue
	Token token.Token
	Value string
}

func (iv *IntValue) TokenLiteral() string {
	return iv.Token.Literal
}

func (iv *IntValue) GetKind() string {
	return iv.Kind
}

func (iv *IntValue) GetValue() interface{} {
	return iv.Value
}
func (iv *IntValue) String() string {
	return fmt.Sprint(iv.Value)
}

var _ Node = (*FloatValue)(nil)
var _ Value = (*FloatValue)(nil)

type FloatValue struct {
	Kind  string //FLOAT_VALUE
	Token token.Token
	Value string
}

func (v *FloatValue) TokenLiteral() string {
	return v.Token.Literal
}
func (v *FloatValue) GetKind() string {
	return v.Kind
}
func (v *FloatValue) GetValue() interface{} {
	return v.Value
}
func (v *FloatValue) String() string {
	return fmt.Sprint(v.Value)
}

var _ Node = (*EnumValue)(nil)
var _ Value = (*EnumValue)(nil)

type EnumValue struct {
	Kind  string //ENUM_VALUE
	Token token.Token
	Value string
}

func (v *EnumValue) TokenLiteral() string {
	return v.Token.Literal
}
func (v *EnumValue) GetKind() string {
	return v.Kind
}
func (v *EnumValue) GetValue() interface{} {
	return v.Value
}
func (v *EnumValue) String() string {
	return fmt.Sprint(v.Value)
}

var _ Node = (*ListValue)(nil)
var _ Value = (*ListValue)(nil)

type ListValue struct {
	Kind   string //LIST_VALUE
	Token  token.Token
	Values []Value
}

func (v *ListValue) TokenLiteral() string {
	return v.Token.Literal
}
func (v *ListValue) GetKind() string {
	return v.Kind
}
func (v *ListValue) GetValue() interface{} {
	return v.Values
}

func (v *ListValue) String() string {
	var vals string
	for _, val := range v.Values {
		vals += fmt.Sprintf("%v, ", val)
	}
	vals = strings.TrimRight(vals, ", ")
	return fmt.Sprintf("[%s]", vals)
}

var _ Node = (*ObjectValue)(nil)
var _ Value = (*ObjectValue)(nil)

type ObjectValue struct {
	Kind   string //OBJECT_VALUE
	Token  token.Token
	Fields []*ObjectField
}

func (v *ObjectValue) TokenLiteral() string {
	return v.Token.Literal
}
func (v *ObjectValue) GetKind() string {
	return v.Kind
}
func (v *ObjectValue) GetValue() interface{} {
	return v.Fields
}
func (v *ObjectValue) String() string {
	var ofields string
	for _, f := range v.Fields {
		ofields += fmt.Sprintf("%s, ", f.String())
	}
	ofields = strings.TrimRight(ofields, ", ")
	return fmt.Sprintf("{%s}", ofields)
}

var _ Node = (*ObjectField)(nil)
var _ Value = (*ObjectField)(nil)

type ObjectField struct {
	Kind  string //OBJECT_FIELD
	Token token.Token
	Name  *Name
	Value Value
}

func (o *ObjectField) TokenLiteral() string {
	return o.Token.Literal
}
func (o *ObjectField) GetKind() string {
	return o.Kind
}
func (o *ObjectField) GetValue() interface{} {
	return o.Value
}

func (o *ObjectField) String() string {
	return fmt.Sprintf("%s: %s", o.Name.String(), o.Value)
}
