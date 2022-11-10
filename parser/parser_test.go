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

func TestUnionTypeDefinition(t *testing.T) {

	input := `"""
test
"""
union SearchResult = Photo | Person`

	lex := lexer.New(input)
	p := New(lex)
	doc := p.ParseDocument()
	def := doc.Definitions[0]

	if def.String() != input {
		t.Errorf("wrong output,expected=%q, got=%q", input, def.String())
	}

}

func TestEnumTypeDefinition(t *testing.T) {

	input := `"""
description test
"""
enum Direction  @skip(name: true, age: false) {
  NORTH
  EAST
  SOUTH
  WEST
}`
	lex := lexer.New(input)
	p := New(lex)
	doc := p.ParseDocument()
	def := doc.Definitions[0]
	if def.String() != input {
		t.Errorf("wrong output,expected=%q, got=%q", input, def.String())
	}
}

func TestInputObjectTypeDefinition(t *testing.T) {

	input := `"""
test description
"""
input Example  @skip(name: true, age: false) {
  self: [Example!]!
  picture: Url = "https://mateors.com"
}`
	lex := lexer.New(input)
	p := New(lex)
	doc := p.ParseDocument()
	def := doc.Definitions[0]
	if def.String() != input {
		t.Errorf("wrong output,expected=%q, got=%q", input, def.String())
	}
}

func TestScalarTypeDefinition(t *testing.T) {

	input := `"""
test description
"""
scalar UUID  @specifiedBy(url: "https://tools.ietf.org/html/rfc4122")`
	lex := lexer.New(input)
	p := New(lex)
	doc := p.ParseDocument()
	def := doc.Definitions[0]
	if def.String() != input {
		t.Errorf("wrong output,expected=%q, got=%q", input, def.String())
	}
}

func TestOperationTypeDefinition(t *testing.T) {

	input := `query GetBooksAndAuthors($name: String = "Mostain")  @skip(cache: true) {
  books(id: 4) {
    title
  }
  authors {
    name
  }
}`
	lex := lexer.New(input)
	p := New(lex)
	doc := p.ParseDocument()
	def := doc.Definitions[0]
	if def.String() != input {
		t.Errorf("wrong output,expected=%q, got=%q", input, def.String())
	}
}

func TestFragmentSpread(t *testing.T) {

	input := `query withFragments {
  user(id: 4) {
    friends(first: 10) {
      ...friendFields
    }
    mutualFriends(first: 10) {
      ...friendFields
    }
  }
}`
	lex := lexer.New(input)
	p := New(lex)
	doc := p.ParseDocument()
	def := doc.Definitions[0]
	if def.String() != input {
		t.Errorf("wrong output,expected=%q, got=%q", input, def.String())
	}
}

func TestInlineFragment(t *testing.T) {

	input := `query inlineFragmentTyping  @skip(cache: true) {
  profiles(handles: ["zuck", "coca-cola"]) {
    handle
    ... on User {
      friends {
        count
      }
    }
    ... on Page {
      likers {
        count
      }
    }
  }
}`
	lex := lexer.New(input)
	p := New(lex)
	doc := p.ParseDocument()
	def := doc.Definitions[0]
	if def.String() != input {
		t.Errorf("wrong output,expected=%q, got=%q", input, def.String())
	}
}

func TestQueryShortHand(t *testing.T) { //OperationTypeDefinition

	input := `{
  user(id: 4) {
    id
    name
    smallPic: profilePic(size: 64)
    bigPic: profilePic(size: 1024)
  }
}`
	lex := lexer.New(input)
	p := New(lex)
	doc := p.ParseDocument()
	def := doc.Definitions[0]
	if def.String() != input {
		t.Errorf("wrong output,expected=%q, got=%q", input, def.String())
	}
}

func TestFields(t *testing.T) { //OperationTypeDefinition

	input := `{
  me {
    id
    firstName
    lastName
    birthday {
      month
      day
    }
    friends {
      name
    }
  }
}`
	lex := lexer.New(input)
	p := New(lex)
	doc := p.ParseDocument()
	def := doc.Definitions[0]
	if def.String() != input {
		t.Errorf("wrong output,expected=%q, got=%q", input, def.String())
	}
}
