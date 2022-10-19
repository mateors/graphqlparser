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

func TestObjectDefinition(t *testing.T) {

	field := &FieldDefinition{}
	field.Kind = FIELD_DEFINITION
	field.Name = &Name{Kind: NAME, Value: "name"}
	//field.Type = &ListType{Kind: LIST_TYPE, Type: &NonNullType{Kind: NONNULL_TYPE, Type: &NamedType{Kind: NAMED_TYPE, Name: &Name{Kind: NAME, Value: "String"}}}}
	field.Type = &NonNullType{Kind: NONNULL_TYPE, Type: &ListType{Kind: LIST_TYPE, Type: &NonNullType{Kind: NONNULL_TYPE, Type: &NamedType{Kind: NAMED_TYPE, Name: &Name{Kind: NAME, Value: "String"}}}}}

	dfields := []*FieldDefinition{}
	dfields = append(dfields, field)

	field2 := &FieldDefinition{}
	field2.Kind = FIELD_DEFINITION
	field2.Name = &Name{Kind: NAME, Value: "age"}
	//field2.Type = &NonNullType{Kind: NAMED_TYPE, Type: &Name{Kind: NAME, Value: "Int"} }
	field2.Type = &NonNullType{Kind: NONNULL_TYPE, Type: &NamedType{Kind: NAMED_TYPE, Name: &Name{Kind: NAME, Value: "Int"}}}
	dfields = append(dfields, field2)

	infcs := []*NamedType{}
	infcs = append(infcs, &NamedType{Kind: NAMED_TYPE, Name: &Name{Kind: NAME, Value: "Abs"}})
	infcs = append(infcs, &NamedType{Kind: NAMED_TYPE, Name: &Name{Kind: NAME, Value: "Book"}})

	args := []*Argument{}
	args = append(args, &Argument{Kind: ARGUMENT, Name: &Name{Kind: NAME, Value: "name"}, Value: &StringValue{Kind: STRING_VALUE, Value: "photo"}})
	args = append(args, &Argument{Kind: ARGUMENT, Name: &Name{Kind: NAME, Value: "caching"}, Value: &BooleanValue{Kind: STRING_VALUE, Value: true}})

	fields := []*ObjectField{}
	fields = append(fields, &ObjectField{Kind: OBJECT_FIELD, Name: &Name{Kind: NAME, Value: "lat"}, Value: &FloatValue{Kind: FLOAT_VALUE, Value: "12.43"}})
	fields = append(fields, &ObjectField{Kind: OBJECT_FIELD, Name: &Name{Kind: NAME, Value: "long"}, Value: &IntValue{Kind: INT_VALUE, Value: "212"}})
	args = append(args, &Argument{Kind: ARGUMENT, Name: &Name{Kind: NAME, Value: "location"}, Value: &ObjectValue{Kind: OBJECT_VALUE, Fields: fields}})

	directives := []*Directive{}
	directives = append(directives, &Directive{
		Kind:      DIRECTIVE,
		Name:      &Name{Kind: NAME, Value: "excludeField"},
		Arguments: args,
	})

	var obj ObjectDefinition
	obj.Kind = OBJECT_DEFINITION
	obj.Name = &Name{Kind: NAME, Value: "Lift"}
	obj.Interfaces = infcs
	obj.Directives = directives
	obj.Fields = dfields

	expectedOutput := `type Lift implements Abs & Book @excludeField(name: "photo", caching: true, location: {lat: 12.43, long: 212}) {
name: [String!]!
age: Int!
}`
	if obj.String() != expectedOutput {
		t.Errorf("wrong output,expected=%q, got=%q %d/%d", expectedOutput, obj.String(), len(expectedOutput), len(obj.String()))
	}
}

func TestObjectDefinitionFieldDirective(t *testing.T) {

	args := []*Argument{}
	args = append(args, &Argument{Kind: ARGUMENT, Name: &Name{Kind: NAME, Value: "name"}, Value: &StringValue{Kind: STRING_VALUE, Value: "photo"}})
	args = append(args, &Argument{Kind: ARGUMENT, Name: &Name{Kind: NAME, Value: "caching"}, Value: &BooleanValue{Kind: STRING_VALUE, Value: true}})

	fields := []*ObjectField{}
	fields = append(fields, &ObjectField{Kind: OBJECT_FIELD, Name: &Name{Kind: NAME, Value: "lat"}, Value: &FloatValue{Kind: FLOAT_VALUE, Value: "12.43"}})
	fields = append(fields, &ObjectField{Kind: OBJECT_FIELD, Name: &Name{Kind: NAME, Value: "long"}, Value: &IntValue{Kind: INT_VALUE, Value: "212"}})
	args = append(args, &Argument{Kind: ARGUMENT, Name: &Name{Kind: NAME, Value: "location"}, Value: &ObjectValue{Kind: OBJECT_VALUE, Fields: fields}})

	directives := []*Directive{}
	directives = append(directives, &Directive{
		Kind:      DIRECTIVE,
		Name:      &Name{Kind: NAME, Value: "excludeField"},
		Arguments: args,
	})

	field := &FieldDefinition{}
	field.Kind = FIELD_DEFINITION
	field.Description = "Field comments or description"
	field.Name = &Name{Kind: NAME, Value: "name"}
	field.Type = &NonNullType{Kind: NONNULL_TYPE, Type: &ListType{Kind: LIST_TYPE, Type: &NonNullType{Kind: NONNULL_TYPE, Type: &NamedType{Kind: NAMED_TYPE, Name: &Name{Kind: NAME, Value: "String"}}}}}
	field.Directives = nil

	dfields := []*FieldDefinition{}
	dfields = append(dfields, field)

	field2 := &FieldDefinition{}
	field2.Kind = FIELD_DEFINITION
	field2.Name = &Name{Kind: NAME, Value: "age"}
	field2.Type = &NonNullType{Kind: NONNULL_TYPE, Type: &NamedType{Kind: NAMED_TYPE, Name: &Name{Kind: NAME, Value: "Int"}}}
	field2.Directives = directives
	dfields = append(dfields, field2)

	infcs := []*NamedType{}
	infcs = append(infcs, &NamedType{Kind: NAMED_TYPE, Name: &Name{Kind: NAME, Value: "Abs"}})
	infcs = append(infcs, &NamedType{Kind: NAMED_TYPE, Name: &Name{Kind: NAME, Value: "Book"}})

	var obj ObjectDefinition
	obj.Kind = OBJECT_DEFINITION
	obj.Description = "Description for the type"
	obj.Name = &Name{Kind: NAME, Value: "Lift"}
	obj.Interfaces = infcs
	obj.Directives = nil
	obj.Fields = dfields

	expectedOutput := `"""
Description for the type
"""
type Lift implements Abs & Book {
"Field comments or description"
name: [String!]!
age: Int! @excludeField(name: "photo", caching: true, location: {lat: 12.43, long: 212})
}`
	if obj.String() != expectedOutput {
		t.Errorf("wrong output,expected=%q, got=%q %d/%d", expectedOutput, obj.String(), len(expectedOutput), len(obj.String()))
	}
}

func TestInterfaceDefinition(t *testing.T) {

	id := &InterfaceDefinition{}
	id.Description = &StringValue{Kind: STRING_VALUE, Value: ""}
	id.Name = &Name{Kind: NAME, Value: "Image"}

	infcs1 := []*NamedType{}
	infcs1 = append(infcs1, &NamedType{Kind: NAMED_TYPE, Name: &Name{Kind: NAME, Value: "Resource"}})
	infcs1 = append(infcs1, &NamedType{Kind: NAMED_TYPE, Name: &Name{Kind: NAME, Value: "Node"}})
	id.Interfaces = infcs1

	// idirectives := []*Directive{}
	// idirectives = append(idirectives, &Directive{
	// 	Kind:  DIRECTIVE,
	// 	Token: token.Token{},
	// 	Name:  &Name{Kind: NAME, Value: "addExternalFields"},
	// 	Arguments: []*Argument{
	// 		{
	// 			Kind:  ARGUMENT,
	// 			Token: token.Token{},
	// 			Name:  &Name{Kind: NAME, Value: "name"},
	// 			Value: &StringValue{Kind: STRING_VALUE, Value: "photo"},
	// 		}, {
	// 			Kind:  ARGUMENT,
	// 			Token: token.Token{},
	// 			Name:  &Name{Kind: NAME, Value: "cache"},
	// 			Value: &BooleanValue{Kind: BOOLEAN_VALUE, Value: true},
	// 		}},
	// })
	id.Directives = nil

	fieldi := &FieldDefinition{}
	fieldi.Kind = FIELD_DEFINITION
	fieldi.Name = &Name{Kind: NAME, Value: "name"}
	fieldi.Type = &NamedType{Kind: NAMED_TYPE, Name: &Name{Kind: NAME, Value: "String"}}

	fieldi2 := &FieldDefinition{}
	fieldi2.Kind = FIELD_DEFINITION
	fieldi2.Name = &Name{Kind: NAME, Value: "value"}
	fieldi2.Type = &NamedType{Kind: NAMED_TYPE, Name: &Name{Kind: NAME, Value: "Int"}}

	id.Fields = []*FieldDefinition{
		{Kind: FIELD_DEFINITION, Description: "", Name: &Name{Kind: NAME, Value: "id"}, Arguments: nil, Type: &NamedType{Kind: NAMED_TYPE, Name: &Name{Kind: NAME, Value: "ID"}}, Directives: nil},
		{Kind: FIELD_DEFINITION, Description: "", Name: &Name{Kind: NAME, Value: "url"}, Arguments: nil, Type: &NamedType{Kind: NAMED_TYPE, Name: &Name{Kind: NAME, Value: "String"}}, Directives: nil},
		{Kind: FIELD_DEFINITION, Description: "", Name: &Name{Kind: NAME, Value: "thumbnail"}, Arguments: nil, Type: &NamedType{Kind: NAMED_TYPE, Name: &Name{Kind: NAME, Value: "String"}}, Directives: nil},
	}

	expectedOutput := `interface Image implements Resource & Node {
id: ID
url: String
thumbnail: String
}`

	if id.String() != expectedOutput {
		t.Errorf("wrong output,expected=%q, got=%q %d/%d", expectedOutput, id.String(), len(expectedOutput), len(id.String()))
	}

}
