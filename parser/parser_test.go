package parser

import (
	"testing"

	"github.com/mateors/graphqlparser/lexer"
)

func TestObjectTypeDefinition(t *testing.T) {

	input := `
	"""
	Test description
	"""
	type Person implements Human @skip(name:true, age:false) {
		id: ID!
		age: []!
		length("Yes" unit: LengthUnit = METER, "No" corner: Int = 50): Float
		oldField: String @deprecated(reason: "Use newField.")
	}`

	lex := lexer.New(input)
	p := New(lex) //parser.New(lex)
	doc := p.ParseDocument()
	// for i, def := range doc.Definitions {
	// 	fmt.Println("*", i, def.GetKind(), def)
	// }

	def := doc.Definitions[0]

	expectedOutput := `"""
Test description
"""
type Person implements Human @skip(name: true, age: false) {
id: ID!
length("Yes" unit: LengthUnit = METER, "No" corner: Int = 50): Float
oldField: String @deprecated(reason: "Use newField.")
}`

	if def.String() != expectedOutput {
		t.Errorf("wrong output,expected=%q, got=%q", expectedOutput, def.String())
	}

}

func TestObjectTypeDefinition2(t *testing.T) {

	input := `"""
Test description
"""
type Person implements Human @skip(name: true, age: false) {
id: ID!
length("Yes" unit: LengthUnit = METER, "No" corner: Int = 50): Float
oldField: String @deprecated(reason: "Use newField.")
}`

	lex := lexer.New(input)
	p := New(lex)
	doc := p.ParseDocument()
	def := doc.Definitions[0]
	if def.String() != input {
		t.Errorf("wrong output,expected=%q, got=%q", input, def.String())
	}

}

func TestInterfaceTypeDefinition(t *testing.T) {

	input := `"""
test
"""
interface Image implements Resource & Node {
id: ID!
url: String
thumbnail: String
}`

	lex := lexer.New(input)
	p := New(lex)
	doc := p.ParseDocument()
	def := doc.Definitions[0]
	
	if def.String() != input {
		t.Errorf("wrong output,expected=%q, got=%q", input, def.String())
	}

}
