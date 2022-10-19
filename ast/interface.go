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

// var _ TypeSystemDefinition = (*SchemaDefinition)(nil)
var _ TypeSystemDefinition = (TypeDefinition)(nil)

//var _ TypeSystemDefinition = (*TypeExtensionDefinition)(nil)
//var _ TypeSystemDefinition = (*DirectiveDefinition)(nil)

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
