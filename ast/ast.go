package ast

import (
	"bytes"
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/mateors/graphqlparser/token"
)

const (
	OperationTypeQuery        = "query"
	OperationTypeMutation     = "mutation"
	OperationTypeSubscription = "subscription"
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
	Definitions []Node //[]Definition
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

// Ensure that all definition types implements Definition interface
var _ Definition = (*OperationDefinition)(nil)
var _ Definition = (*FragmentDefinition)(nil)
var _ Definition = (TypeSystemDefinition)(nil) // experimental non-spec addition.

type FragmentDefinition struct {
	//fragment FragmentName TypeCondition Directives[opt] SelectionSet
	Kind          string //FRAGMENT_DEFINITION
	Token         token.Token
	Operation     string
	FragmentName  *Name
	TypeCondition *NamedType
	Directives    []*Directive
	SelectionSet  *SelectionSet
}

func (fd *FragmentDefinition) TokenLiteral() string {
	return fd.Token.Literal
}
func (fd *FragmentDefinition) GetKind() string {
	return fd.Kind
}
func (fd *FragmentDefinition) GetOperation() string {
	return fd.Operation
}
func (fd *FragmentDefinition) GetVariableDefinitions() []*VariableDefinition {
	return []*VariableDefinition{}
}
func (fd *FragmentDefinition) GetSelectionSet() *SelectionSet {
	return &SelectionSet{}
}
func (fd *FragmentDefinition) String() string {
	name := fmt.Sprintf("%v", fd.FragmentName)
	typeCondition := fd.TypeCondition.String()
	directives := toSliceString(fd.Directives)

	selectionSet := fd.SelectionSet.String()
	return "fragment " + name + " on " + typeCondition + " " + wrap("", join(directives, " "), " ") + selectionSet
}

type OperationDefinition struct {
	//OperationType Name[opt] VariablesDefinition[opt] Directives[opt] SelectionSet
	Kind                string //OPERATION_DEFINITION
	Token               token.Token
	OperationType       string //query | mutation | subscription
	Name                *Name
	VariablesDefinition []*VariableDefinition
	Directives          []*Directive
	SelectionSet        *SelectionSet
}

func (od *OperationDefinition) TokenLiteral() string {
	return od.Token.Literal
}
func (dd *OperationDefinition) GetKind() string {
	return dd.Kind
}
func (dd *OperationDefinition) GetOperation() string {
	return dd.OperationType
}
func (dd *OperationDefinition) GetVariableDefinitions() []*VariableDefinition {
	return dd.VariablesDefinition
}
func (dd *OperationDefinition) GetSelectionSet() *SelectionSet {
	return dd.SelectionSet
}
func (dd *OperationDefinition) String() string {

	op := dd.OperationType
	name := fmt.Sprintf("%v", dd.Name)
	varDefs := wrap("(", join(toSliceString(dd.VariablesDefinition), ", "), ")")
	directives := join(toSliceString(dd.Directives), " ")
	selectionSet := fmt.Sprintf("%v", dd.SelectionSet)
	// Anonymous queries with no directives or variable definitions can use
	// the query short form.
	str := ""
	if name == "" && directives == "" && varDefs == "" && op == OperationTypeQuery {
		str = selectionSet
	} else {
		str = join([]string{
			op,
			join([]string{name, varDefs}, ""),
			directives,
			selectionSet,
		}, " ")
	}
	return str
}

var _ Node = (*Name)(nil)
var _ Type = (*Name)(nil)

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

var _ Node = (*VariableDefinition)(nil)

type VariableDefinition struct {
	//Variable : Type DefaultValue[opt] Directives[opt]
	Kind         string
	Token        token.Token
	Variable     *Variable
	Type         Type
	DefaultValue Value
	Directives   []*Directive
}

func (vd *VariableDefinition) TokenLiteral() string {
	return vd.Token.Literal
}
func (vd *VariableDefinition) GetKind() string {
	return vd.Kind
}

func (vd *VariableDefinition) String() string {

	var defaultValue string
	variable := fmt.Sprintf("%v", vd.Variable)
	ttype := fmt.Sprintf("%v", vd.Type)
	if vd.DefaultValue != nil {
		defaultValue = fmt.Sprintf("%v", vd.DefaultValue)
	}
	directives := join(toSliceString(vd.Directives), " ")
	return variable + ": " + ttype + wrap(" = ", defaultValue, "") + wrap("", directives, "")
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
	out.WriteString(" @" + d.Name.String() + "(" + args + ")")
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

var _ Node = (*SelectionSet)(nil)

func (ss *SelectionSet) TokenLiteral() string {
	return ss.Token.Literal
}
func (ss *SelectionSet) GetKind() string {
	return ss.Kind
}
func (ss *SelectionSet) String() string {

	//fmt.Println(len(ss.Selections), ss.Selections)
	return block(ss.Selections)
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
func (f *Field) String() string {

	var alias, name, selectionSet string
	name = f.Name.String()
	if f.Alias != nil {
		alias = f.Alias.String()
	}
	args := toSliceString(f.Arguments)
	directives := toSliceString(f.Directives)

	if f.SelectionSet != nil {
		selectionSet = f.SelectionSet.String()
	}

	str := join(
		[]string{
			wrap("", alias, ": ") + name + wrap("(", join(args, ", "), ")"),
			join(directives, " "),
			selectionSet,
		},
		" ",
	)

	return str
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
func (fs *FragmentSpread) String() string {
	name := fs.FragmentName.String()
	directives := toSliceString(fs.Directives)
	return "..." + name + wrap(" ", join(directives, " "), "")
}

type InlineFragment struct {
	//...TypeCondition[opt] Directives[opt] SelectionSet
	Kind          string //INLINE_FRAGMENT
	Token         token.Token
	TypeCondition *NamedType //on NamedType
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
func (f *InlineFragment) String() string {

	var typeCondition string
	if f.TypeCondition != nil {
		typeCondition = " on " + f.TypeCondition.String()
	}
	directives := toSliceString(f.Directives)
	selectionSet := f.SelectionSet.String()
	return "..." + typeCondition + " " + wrap("", join(directives, " "), " ") + selectionSet
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

var _ TypeDefinition = (*ScalarDefinition)(nil)
var _ TypeDefinition = (*ObjectDefinition)(nil)
var _ TypeDefinition = (*InterfaceDefinition)(nil)
var _ TypeDefinition = (*UnionDefinition)(nil)
var _ TypeDefinition = (*EnumDefinition)(nil)
var _ TypeDefinition = (*InputObjectDefinition)(nil)

type ScalarDefinition struct {
	//Description[opt] scalar Name Directives[opt]
	Kind        string //INPUT_OBJECT_DEFINITION
	Token       token.Token
	Description *StringValue
	Name        *Name
	Directives  []*Directive
}

func (sd *ScalarDefinition) TokenLiteral() string {
	return sd.Token.Literal
}
func (sd *ScalarDefinition) GetKind() string {
	return sd.Kind
}
func (sd *ScalarDefinition) GetOperation() string {
	return ""
}
func (sd *ScalarDefinition) GetVariableDefinitions() []*VariableDefinition {
	return []*VariableDefinition{}
}
func (sd *ScalarDefinition) GetSelectionSet() *SelectionSet {
	return &SelectionSet{}
}
func (sd *ScalarDefinition) String() string {

	name := sd.Name.String()
	directives := toSliceString(sd.Directives)
	str := join([]string{
		"scalar",
		name,
		join(directives, " "),
	}, " ")

	if sd.Description != nil {
		desc := sd.Description.Value
		if desc != "" {
			desc = join([]string{`"""`, desc, `"""`}, "\n")
			str = fmt.Sprintf("%s\n%s", desc, str)
		}
	}
	return str
}

type InputObjectDefinition struct {
	//Description[opt] input Name Directives[opt] InputFieldsDefinition
	//Description[opt] input Name Directives[opt]
	Kind        string //INPUT_OBJECT_DEFINITION
	Token       token.Token
	Description *StringValue
	Name        *Name
	Directives  []*Directive
	Fields      []*InputValueDefinition
}

func (iod *InputObjectDefinition) TokenLiteral() string {
	return iod.Token.Literal
}
func (iod *InputObjectDefinition) GetKind() string {
	return iod.Kind
}
func (iod *InputObjectDefinition) GetOperation() string {
	return ""
}
func (iod *InputObjectDefinition) GetVariableDefinitions() []*VariableDefinition {
	return []*VariableDefinition{}
}
func (iod *InputObjectDefinition) GetSelectionSet() *SelectionSet {
	return &SelectionSet{}
}
func (iod *InputObjectDefinition) String() string {
	name := fmt.Sprintf("%v", iod.Name)
	fields := iod.Fields
	directives := toSliceString(iod.Directives)
	str := join([]string{
		"input",
		name,
		join(directives, " "),
		block(fields),
	}, " ")

	if iod.Description != nil {
		desc := iod.Description.Value
		if desc != "" {
			desc = join([]string{`"""`, desc, `"""`}, "\n")
			str = fmt.Sprintf("%s\n%s", desc, str)
		}
	}
	return str
}

type EnumDefinition struct {
	//Description[opt] enum Name Directives[opt] EnumValuesDefinition
	//Description[opt] enum Name Directives[opt]
	Kind        string //ENUM_DEFINITION
	Token       token.Token
	Description *StringValue
	Name        *Name
	Directives  []*Directive
	Values      []*EnumValueDefinition
}

type EnumValueDefinition struct {
	//Description[opt] EnumValue Directives[opt]
	Kind        string //ENUMVALUE_DEFINITION
	Token       token.Token
	Description *StringValue
	Name        *Name
	Directives  []*Directive
}

func (evd *EnumValueDefinition) TokenLiteral() string {
	return evd.Token.Literal
}
func (evd *EnumValueDefinition) GetKind() string {
	return evd.Kind
}
func (evd *EnumValueDefinition) String() string {

	name := fmt.Sprintf("%v", evd.Name)
	directives := toSliceString(evd.Directives)

	str := join([]string{
		name,
		join(directives, " "),
	}, " ")

	if evd.Description != nil {
		desc := evd.Description.Value
		if desc != "" {
			desc = join([]string{`"""`, desc, `"""`}, "\n")
			str = fmt.Sprintf("%s\n%s", desc, str)
		}
	}
	return str
}

func (ed *EnumDefinition) TokenLiteral() string {
	return ed.Token.Literal
}
func (ed *EnumDefinition) GetKind() string {
	return ed.Kind
}
func (ed *EnumDefinition) GetOperation() string {
	return ""
}
func (ed *EnumDefinition) GetVariableDefinitions() []*VariableDefinition {
	return []*VariableDefinition{}
}
func (ed *EnumDefinition) GetSelectionSet() *SelectionSet {
	return &SelectionSet{}
}
func (ed *EnumDefinition) String() string {

	name := fmt.Sprintf("%v", ed.Name)
	values := toSliceString(ed.Values) //?

	directives := toSliceString(ed.Directives)
	str := join([]string{
		"enum",
		name,
		join(directives, " "),
		block(values),
	}, " ")

	if ed.Description != nil {
		desc := ed.Description.Value
		if desc != "" {
			desc = join([]string{`"""`, desc, `"""`}, "\n")
			str = fmt.Sprintf("%s\n%s", desc, str)
		}
	}
	return str
}

func block(Slice interface{}) string {
	s := toSliceString(Slice)
	if len(s) == 0 {
		return "{}"
	}
	//return "{\n" + join(s, "\n") + "\n}"
	return indent("{\n"+join(s, "\n")) + "\n}"
}
func indent(istring interface{}) string {
	if istring == nil {
		return ""
	}
	switch str := istring.(type) {
	case string:
		return strings.Replace(str, "\n", "\n  ", -1)
	}
	return ""
}

type UnionDefinition struct {
	//Description[opt] union Name Directives[opt] UnionMemberTypes[opt]
	Kind             string //UNION_DEFINITION
	Token            token.Token
	Description      *StringValue
	Name             *Name
	Directives       []*Directive
	UnionMemberTypes []*NamedType
}

func (ud *UnionDefinition) TokenLiteral() string {
	return ud.Token.Literal
}
func (ud *UnionDefinition) GetKind() string {
	return ud.Kind
}
func (ud *UnionDefinition) GetOperation() string {
	return ""
}
func (ud *UnionDefinition) GetVariableDefinitions() []*VariableDefinition {
	return []*VariableDefinition{}
}
func (ud *UnionDefinition) GetSelectionSet() *SelectionSet {
	return &SelectionSet{}
}
func (ud *UnionDefinition) String() string {

	name := fmt.Sprintf("%v", ud.Name)
	types := toSliceString(ud.UnionMemberTypes)

	directives := toSliceString(ud.Directives)
	str := join([]string{
		"union",
		name,
		join(directives, " "),
		"= " + join(types, " | "),
	}, " ")

	if ud.Description != nil {
		desc := ud.Description.Value
		if desc != "" {
			desc = join([]string{`"""`, desc, `"""`}, "\n")
			str = fmt.Sprintf("%s\n%s", desc, str)
		}
	}
	return str
}

func join(str []string, sep string) string {
	slc := []string{}
	for _, s := range str {
		if s == "" {
			continue
		}
		slc = append(slc, s)
	}
	return strings.Join(slc, sep)
}

func toSliceString(slice interface{}) []string {

	defer func() {
		if rec := recover(); rec != nil {
			log.Println("activity panicking with value >>", rec)
		}
	}()
	//fmt.Println("--toSliceString")
	if slice == nil {
		return []string{}
	}
	res := []string{}
	switch reflect.TypeOf(slice).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(slice)
		for i := 0; i < s.Len(); i++ {
			elem := s.Index(i)
			reflect.ValueOf(&elem).MethodByName("String").Call([]reflect.Value{})
			elemv := fmt.Sprintf("%v", elem)
			res = append(res, elemv)
		}
		return res
	default:
		return res
	}
}

type ObjectDefinition struct {
	//Description[opt] type Name ImplementsInterfaces[opt] Directives[opt] { FieldsDefinition }
	Kind        string //OBJECT_DEFINITION
	Token       token.Token
	Description *StringValue
	Name        *Name
	Interfaces  []*NamedType //NAMED
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

	if ob.Description != nil {
		if len(ob.Description.Value) > 0 {
			//desc := join([]string{`"""`, desc, `"""`}, sep)
			out.WriteString(fmt.Sprintf("\"\"\"\n%s\n\"\"\"", ob.Description.Value) + "\n")
		}
	}
	name := ob.Name.Value
	out.WriteString("type" + " " + name)

	if len(ob.Interfaces) > 0 {
		out.WriteString(" implements" + " ")
		var infcs string
		for _, inf := range ob.Interfaces {
			//iname := inf.Value
			infcs += fmt.Sprintf("%s & ", inf)
		}
		infcs = strings.TrimRight(infcs, " & ")
		out.WriteString(infcs)
	}

	directives := []string{}
	for _, directive := range ob.Directives {
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

	out.WriteString(" {")

	for _, field := range ob.Fields {
		out.WriteString("\n" + field.String())
	}

	out.WriteString("\n}")
	return out.String()
}

var _ Node = (*FieldDefinition)(nil)

type FieldDefinition struct {
	//Description[opt] Name ArgumentsDefinition[opt] : Type Directives[opt]
	//ArgumentsDefinition: ( InputValueDefinition[list] )
	Kind        string //FIELD_DEFINITION
	Token       token.Token
	Description *StringValue
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
	if fd.Description != nil {
		if len(fd.Description.Value) > 0 {
			out.WriteString(fmt.Sprintf("\"%s\"\n", fd.Description.Value))
		}
	}
	out.WriteString(fd.Name.String())
	if len(fd.Arguments) > 0 {
		out.WriteString("(")
		var vals string
		for _, arg := range fd.Arguments {
			vals += fmt.Sprintf("%s, ", arg.String())
		}
		vals = strings.TrimRight(vals, ", ")
		out.WriteString(vals)
		out.WriteString(")")
		//out.WriteString(": " + fd.Type.String())
	}
	out.WriteString(": " + fd.Type.String()) //?

	directives := []string{}
	for _, directive := range fd.Directives {
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
var _ Type = (*StringValue)(nil)

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
	var vals string = s.Value
	if s.Kind == STRING_VALUE {
		vals = fmt.Sprintf(`"%s"`, s.Value)
	}
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
