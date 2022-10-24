package ast

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/mateors/graphqlparser/token"
)

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

var _ TypeSystemDefinition = (*SchemaDefinition)(nil)
var _ TypeSystemDefinition = (TypeDefinition)(nil)

// var _ TypeSystemDefinition = (*TypeExtensionDefinition)(nil)
var _ TypeSystemDefinition = (*DirectiveDefinition)(nil)

// TypeSystemExtension
type TypeExtensionDefinition struct {
}

// SchemaDefinition implements Node, Definition
type SchemaDefinition struct {
	//Description[opt] schema Directives[opt] { RootOperationTypeDefinition[list] }
	Kind           string //SCHEMA_DEFINITION
	Token          token.Token
	Description    *StringValue
	Directives     []*Directive
	OperationTypes []*RootOperationTypeDefinition
}

func (sd *SchemaDefinition) TokenLiteral() string {
	return sd.Token.Literal
}
func (sd *SchemaDefinition) GetKind() string {
	return sd.Kind
}
func (sd *SchemaDefinition) GetOperation() string {
	return ""
}
func (sd *SchemaDefinition) GetVariableDefinitions() []*VariableDefinition {
	return []*VariableDefinition{}
}
func (sd *SchemaDefinition) GetSelectionSet() *SelectionSet {
	return &SelectionSet{}
}
func (sd *SchemaDefinition) String() string {

	directives := toSliceString(sd.Directives)
	str := join([]string{
		"schema",
		join(directives, " "),
		block(sd.OperationTypes),
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

// OperationTypeDefinition implements Node, Definition
type RootOperationTypeDefinition struct {
	//OperationType:NamedType
	Kind          string //SCHEMA_DEFINITION
	Token         token.Token
	OperationType string //query | muation | subscription
	NamedType     *NamedType
}

func (rotd *RootOperationTypeDefinition) TokenLiteral() string {
	return rotd.Token.Literal
}
func (rotd *RootOperationTypeDefinition) GetKind() string {
	return rotd.Kind
}
func (rotd *RootOperationTypeDefinition) GetOperation() string {
	return ""
}
func (rotd *RootOperationTypeDefinition) GetVariableDefinitions() []*VariableDefinition {
	return []*VariableDefinition{}
}
func (rotd *RootOperationTypeDefinition) GetSelectionSet() *SelectionSet {
	return &SelectionSet{}
}
func (rotd *RootOperationTypeDefinition) String() string {
	str := fmt.Sprintf("%v: %v", rotd.OperationType, rotd.NamedType.String())
	return str
}

type DirectiveDefinition struct {
	//Description[opt] directive @Name ArgumentsDefinition[opt] repeatable[opt] on DirectiveLocations
	Kind        string //DIRECTIVE_DEFINITION
	Token       token.Token
	Description *StringValue
	Name        *Name
	Arguments   []*InputValueDefinition
	Locations   []*Name
}

func (dd *DirectiveDefinition) TokenLiteral() string {
	return dd.Token.Literal
}
func (dd *DirectiveDefinition) GetKind() string {
	return dd.Kind
}
func (dd *DirectiveDefinition) GetOperation() string {
	return ""
}
func (dd *DirectiveDefinition) GetVariableDefinitions() []*VariableDefinition {
	return []*VariableDefinition{}
}
func (dd *DirectiveDefinition) GetSelectionSet() *SelectionSet {
	return &SelectionSet{}
}
func (dd *DirectiveDefinition) String() string {

	args := toSliceString(dd.Arguments)
	hasArgDesc := false
	for _, arg := range dd.Arguments {
		if arg.Description != "" {
			hasArgDesc = true
			break
		}
	}
	var argsStr string
	if hasArgDesc {
		argsStr = wrap("(", indent("\n"+join(args, "\n")), "\n)")
	} else {
		argsStr = wrap("(", join(args, ", "), ")")
	}
	str := fmt.Sprintf("directive @%v%v on %v", dd.Name, argsStr, join(toSliceString(dd.Locations), " | "))

	if dd.Description != nil {
		desc := dd.Description.Value
		if desc != "" {
			desc = join([]string{`"""`, desc, `"""`}, "\n")
			str = fmt.Sprintf("%s\n%s", desc, str)
		}
	}
	return str
}

func (i *InterfaceDefinition) TokenLiteral() string {
	return i.Token.Literal
}
func (i *InterfaceDefinition) GetKind() string {
	return i.Kind
}
func (i *InterfaceDefinition) GetOperation() string {
	return ""
}
func (i *InterfaceDefinition) GetVariableDefinitions() []*VariableDefinition {
	return []*VariableDefinition{}
}
func (i *InterfaceDefinition) GetSelectionSet() *SelectionSet {
	return &SelectionSet{}
}

func (i *InterfaceDefinition) String() string {

	var out bytes.Buffer
	if len(i.Description.Value) > 0 {
		//desc := join([]string{`"""`, desc, `"""`}, sep)
		out.WriteString(fmt.Sprintf("\"\"\"\n%s\n\"\"\"", i.Description) + "\n")
	}
	name := i.Name.Value
	out.WriteString("interface" + " " + name)

	if len(i.Interfaces) > 0 {
		out.WriteString(" implements" + " ")
		var infcs string
		for _, inf := range i.Interfaces {
			infcs += fmt.Sprintf("%s & ", inf)
		}
		infcs = strings.TrimRight(infcs, " & ")
		out.WriteString(infcs)
	}

	directives := []string{}
	for _, directive := range i.Directives {
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
	for _, field := range i.Fields {
		out.WriteString("\n" + field.String())
	}
	out.WriteString("\n}")

	return out.String()
}

func wrap(start, text, end string) string {
	if text == "" {
		return text
	}
	return start + text + end
}
