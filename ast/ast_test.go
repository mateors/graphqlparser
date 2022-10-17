package ast

import (
	"testing"

	"github.com/mateors/graphqlparser/token"
)

func TestInputValueDefinition(t *testing.T) {

	ivd := InputValueDefinition{}
	ivd.Name = &Name{Kind: "Name", Token: token.Token{}, Value: "name"}

	nType := &NamedType{Kind: NAMED_TYPE, Token: token.Token{}, Name: &Name{Kind: "Name", Token: token.Token{}, Value: "String"}}
	ivd.Type = nType
	if ivd.String() != "name: String" {
		t.Errorf("wrong output, got=%q", ivd.String())
	}

	nnType := &NonNullType{Kind: NONNULL_TYPE, Token: token.Token{}, Type: &NamedType{Kind: "Name", Token: token.Token{}, Name: &Name{Kind: "Name", Token: token.Token{}, Value: "String"}}}
	ivd.Type = nnType
	if ivd.String() != "name: String!" {
		t.Errorf("wrong output, got=%q", ivd.String())
	}

	lt := &ListType{Kind: LIST_TYPE, Token: token.Token{}, Type: &NamedType{Kind: "Name", Token: token.Token{}, Name: &Name{Kind: "Name", Token: token.Token{}, Value: "String"}}}
	ivd.Type = lt
	if ivd.String() != "name: [String]" {
		t.Errorf("wrong output, got=%q", ivd.String())
	}

	nnTypeL := &NonNullType{Kind: NONNULL_TYPE, Token: token.Token{}, Type: &ListType{Kind: LIST_TYPE, Token: token.Token{}, Type: &NonNullType{Kind: NONNULL_TYPE, Token: token.Token{}, Type: &NamedType{Kind: "Name", Token: token.Token{}, Name: &Name{Kind: "Name", Token: token.Token{}, Value: "String"}}}}}
	ivd.Type = nnTypeL
	if ivd.String() != "name: [String!]!" {
		t.Errorf("wrong output, got=%q", ivd.String())
	}

}

func TestInputValueDefinition2(t *testing.T) {

	test := []struct {
		ivd            InputValueDefinition
		expectedOutput string
	}{
		{
			InputValueDefinition{Name: &Name{Kind: NAME, Token: token.Token{}, Value: "name"}, Type: &NamedType{Kind: NAMED_TYPE, Token: token.Token{}, Name: &Name{Kind: NAME, Token: token.Token{}, Value: "String"}}},
			"name: String",
		},
		{
			InputValueDefinition{Name: &Name{Kind: NAME, Token: token.Token{}, Value: "name"}, Type: &NonNullType{Kind: NONNULL_TYPE, Token: token.Token{}, Type: &NamedType{Kind: NAME, Token: token.Token{}, Name: &Name{Kind: NAME, Token: token.Token{}, Value: "String"}}}},
			"name: String!",
		},

		{
			InputValueDefinition{Name: &Name{Kind: NAME, Token: token.Token{}, Value: "name"}, Type: &ListType{Kind: LIST_TYPE, Token: token.Token{}, Type: &NamedType{Kind: NAME, Token: token.Token{}, Name: &Name{Kind: NAME, Token: token.Token{}, Value: "String"}}}},
			"name: [String]",
		},
		{
			InputValueDefinition{Name: &Name{Kind: NAME, Token: token.Token{}, Value: "name"}, Type: &ListType{Kind: LIST_TYPE, Token: token.Token{}, Type: &NonNullType{Kind: NAME, Token: token.Token{}, Type: &NamedType{Kind: NAME, Token: token.Token{}, Name: &Name{Kind: NAME, Token: token.Token{}, Value: "String"}}}}},
			"name: [String!]",
		},
		{
			InputValueDefinition{Name: &Name{Kind: NAME, Token: token.Token{}, Value: "name"}, Type: &NonNullType{Kind: NONNULL_TYPE, Token: token.Token{}, Type: &ListType{Kind: LIST_TYPE, Token: token.Token{}, Type: &NonNullType{Kind: NONNULL_TYPE, Token: token.Token{}, Type: &NamedType{Kind: NAME, Token: token.Token{}, Name: &Name{Kind: NAME, Token: token.Token{}, Value: "String"}}}}}},
			"name: [String!]!",
		},
	}

	for i, obj := range test {

		if obj.ivd.String() != obj.expectedOutput {
			t.Errorf("%d wrong output,expected=%q, got=%q", i, obj.expectedOutput, obj.ivd.String())
		}
	}

}

func TestFieldDefinition(t *testing.T) {

	ivd := []*InputValueDefinition{}
	//args := []*Argument{}
	//args = append(args, &Argument{Kind: ARGUMENT, Name: &Name{Kind: NAME, Value: "name"}, Value: &StringValue{Kind: STRING_VALUE, Value: "photo"}})
	// directives := []*Directive{}
	// directives = append(directives, &Directive{
	// 	Kind:      DIRECTIVE,
	// 	Name:      &Name{Kind: NAME, Value: "excludeField"},
	// 	Arguments: args,
	// })

	iv1 := &InputValueDefinition{
		Name:         &Name{Kind: NAME, Token: token.Token{}, Value: "unit"},
		Type:         &NamedType{Kind: NAMED_TYPE, Token: token.Token{}, Name: &Name{Kind: "Name", Token: token.Token{}, Value: "LengthUnit"}},
		DefaultValue: &StringValue{Kind: ENUM_VALUE, Token: token.Token{}, Value: "METER"},
		Directives:   nil,
	}
	ivd = append(ivd, iv1)

	field := FieldDefinition{}
	field.Kind = FIELD_DEFINITION
	field.Name = &Name{Kind: NAME, Value: "name"}
	field.Arguments = ivd
	field.Type = &NamedType{Kind: NAMED_TYPE, Name: &Name{Kind: NAME, Value: "String"}}

	//expectedString := `name(unit: LengthUnit = METER @excludeField(name: "photo")): String`
	expectedString := `name(unit: LengthUnit = METER): String`

	if field.String() != expectedString {
		t.Errorf("wrong output,expected=%q, got=%q", expectedString, field.String())
	}

}
