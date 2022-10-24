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
	field.Description = &StringValue{Kind: STRING_VALUE, Value: "Field comments or description"}
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
	obj.Description = &StringValue{Kind: STRING_VALUE, Value: "Description for the type"}
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
	id.Kind = INTERFACE_DEFINITION
	id.Description = &StringValue{Kind: STRING_VALUE, Value: ""}
	id.Name = &Name{Kind: NAME, Value: "Image"}

	infcs1 := []*NamedType{
		{Kind: NAMED_TYPE, Name: &Name{Kind: NAME, Value: "Resource"}},
		{Kind: NAMED_TYPE, Name: &Name{Kind: NAME, Value: "Node"}},
	}
	id.Interfaces = infcs1

	// idirectives := []*Directive{
	// 	{Kind: DIRECTIVE, Name: &Name{Kind: NAME, Value: "addExternalFields"}, Arguments: []*Argument{
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
	// 		}}},

	// 	{Kind: DIRECTIVE, Name: &Name{Kind: NAME, Value: "skip"}, Arguments: []*Argument{
	// 		{
	// 			Kind:  ARGUMENT,
	// 			Token: token.Token{},
	// 			Name:  &Name{Kind: NAME, Value: "name"},
	// 			Value: &StringValue{Kind: ENUM_VALUE, Value: "id"},
	// 		}, {
	// 			Kind:  ARGUMENT,
	// 			Token: token.Token{},
	// 			Name:  &Name{Kind: NAME, Value: "cache"},
	// 			Value: &BooleanValue{Kind: BOOLEAN_VALUE, Value: true},
	// 		}}},
	// }
	// id.Directives = idirectives

	fieldi := &FieldDefinition{}
	fieldi.Kind = FIELD_DEFINITION
	fieldi.Name = &Name{Kind: NAME, Value: "name"}
	fieldi.Type = &NamedType{Kind: NAMED_TYPE, Name: &Name{Kind: NAME, Value: "String"}}

	fieldi2 := &FieldDefinition{}
	fieldi2.Kind = FIELD_DEFINITION
	fieldi2.Name = &Name{Kind: NAME, Value: "value"}
	fieldi2.Type = &NamedType{Kind: NAMED_TYPE, Name: &Name{Kind: NAME, Value: "Int"}}

	id.Fields = []*FieldDefinition{
		{Kind: FIELD_DEFINITION, Description: nil, Name: &Name{Kind: NAME, Value: "id"}, Arguments: nil, Type: &NamedType{Kind: NAMED_TYPE, Name: &Name{Kind: NAME, Value: "ID"}}, Directives: nil},
		{Kind: FIELD_DEFINITION, Description: nil, Name: &Name{Kind: NAME, Value: "url"}, Arguments: nil, Type: &NamedType{Kind: NAMED_TYPE, Name: &Name{Kind: NAME, Value: "String"}}, Directives: nil},
		{Kind: FIELD_DEFINITION, Description: nil, Name: &Name{Kind: NAME, Value: "thumbnail"}, Arguments: nil, Type: &NamedType{Kind: NAMED_TYPE, Name: &Name{Kind: NAME, Value: "String"}}, Directives: nil},
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

func TestUnionDefinition(t *testing.T) {

	ud := UnionDefinition{}
	ud.Kind = UNION_DEFINITION
	ud.Description = nil //&StringValue{Kind: STRING_VALUE, Value: "Test des"}
	ud.Name = &Name{Kind: NAME, Value: "SearchResult"}
	// ud.Directives = nil []*Directive{
	// 	{Kind: DIRECTIVE, Name: &Name{Kind: NAME, Value: "skip"}, Arguments: []*Argument{
	// 		{Kind: ARGUMENT, Name: &Name{Kind: NAME, Value: "caching"}, Value: &BooleanValue{Kind: BOOLEAN_VALUE, Value: true}},
	// 	}},
	// }
	ud.UnionMemberTypes = []*NamedType{
		{Kind: NAMED_TYPE, Name: &Name{Kind: NAME, Value: "Photo"}},
		{Kind: NAMED_TYPE, Name: &Name{Kind: NAME, Value: "Person"}},
	}
	expectedOutput := `union SearchResult = Photo | Person`
	if ud.String() != expectedOutput {
		t.Errorf("wrong output,expected=%q, got=%q %d/%d", expectedOutput, ud.String(), len(expectedOutput), len(ud.String()))
	}
}

func TestEnumDefinition(t *testing.T) {

	ed := EnumDefinition{}
	ed.Kind = ENUM_DEFINITION
	ed.Description = &StringValue{Kind: STRING_VALUE, Value: ""}
	ed.Name = &Name{Kind: NAME, Value: "Country"}
	// ed.Directives = []*Directive{
	// 	{Kind: DIRECTIVE, Name: &Name{Kind: NAME, Value: "skip"}, Arguments: []*Argument{
	// 		{Kind: ARGUMENT, Name: &Name{Kind: NAME, Value: "caching"}, Value: &BooleanValue{Kind: BOOLEAN_VALUE, Value: true}},
	// 	}},
	// }
	ed.Values = []*EnumValueDefinition{
		{Kind: ENUMVALUE_DEFINITION, Name: &Name{Kind: NAME, Value: "BANGLADESH"}},
		{Kind: ENUMVALUE_DEFINITION, Name: &Name{Kind: NAME, Value: "INDIA"}},
	}
	expectedOutput := `enum Country {
  BANGLADESH
  INDIA
}`
	if ed.String() != expectedOutput {
		t.Errorf("wrong output,expected=%q, got=%q %d/%d", expectedOutput, ed.String(), len(expectedOutput), len(ed.String()))
	}
}

func TestInputObjectDefinition(t *testing.T) {

	iod := InputObjectDefinition{}
	iod.Kind = INPUT_OBJECT_DEFINITION
	//iod.Description = &StringValue{Kind: STRING_VALUE, Value: "test"}
	iod.Name = &Name{Kind: NAME, Value: "Client"}
	// iod.Directives = []*Directive{
	// 	{Kind: DIRECTIVE, Name: &Name{Kind: NAME, Value: "skip"}, Arguments: []*Argument{
	// 		{Kind: ARGUMENT, Name: &Name{Kind: NAME, Value: "caching"}, Value: &BooleanValue{Kind: BOOLEAN_VALUE, Value: true}},
	// 	}},
	// }
	iod.Fields = []*InputValueDefinition{
		{Kind: INPUT_VALUE_DEFINITION, Name: &Name{Kind: NAME, Value: "name"}, Type: &NonNullType{Kind: STRING_VALUE, Type: &Name{Kind: NAME, Value: "String"}}},
		{DefaultValue: &IntValue{Kind: INT_VALUE, Value: "10"}, Kind: INPUT_VALUE_DEFINITION, Name: &Name{Kind: NAME, Value: "age"}, Type: &NamedType{Kind: STRING_VALUE, Name: &Name{Kind: NAME, Value: "Int"}}},
	}

	expectedOutput := `input Client {
  name: String!
  age: Int = 10
}`
	if iod.String() != expectedOutput {
		t.Errorf("wrong output,expected=%q, got=%q %d/%d", expectedOutput, iod.String(), len(expectedOutput), len(iod.String()))
	}
}

func TestScalarDefiniton(t *testing.T) {

	scd := ScalarDefinition{Kind: SCALAR_DEFINITION}
	//scd.Description = &StringValue{Kind: STRING_VALUE, Value: "test"}
	scd.Name = &Name{Kind: NAME, Value: "Upload"}
	// scd.Directives = []*Directive{
	// 	{Kind: DIRECTIVE, Name: &Name{Kind: NAME, Value: "skip"}, Arguments: []*Argument{
	// 		{Kind: ARGUMENT, Name: &Name{Kind: NAME, Value: "caching"}, Value: &BooleanValue{Kind: BOOLEAN_VALUE, Value: true}},
	// 	}},
	// }
	expectedOutput := `scalar Upload`
	if scd.String() != expectedOutput {
		t.Errorf("wrong output,expected=%q, got=%q %d/%d", expectedOutput, scd.String(), len(expectedOutput), len(scd.String()))
	}
}

func TestDirectiveDefinition(t *testing.T) {

	dd := DirectiveDefinition{}
	dd.Kind = DIRECTIVE_DEFINITION
	//dd.Description = &StringValue{Kind: STRING_VALUE, Value: "test"}
	dd.Name = &Name{Kind: NAME, Value: "cacheControl"}
	dd.Arguments = []*InputValueDefinition{
		{Kind: INPUT_VALUE_DEFINITION, Name: &Name{Kind: NAME, Value: "maxAge"}, Type: &NamedType{Kind: NAMED_TYPE, Name: &Name{Kind: NAME, Value: "Int"}}},
		{Kind: INPUT_VALUE_DEFINITION, Name: &Name{Kind: NAME, Value: "scope"}, Type: &NamedType{Kind: NAMED_TYPE, Name: &Name{Kind: NAME, Value: "CacheControlScope"}}},
	}
	dd.Locations = []*Name{
		{Kind: NAME, Value: "FIELD_DEFINITION"},
		{Kind: NAME, Value: "OBJECT"},
		{Kind: NAME, Value: "INTERFACE"},
	}

	expectedOutput := `directive @cacheControl(maxAge: Int, scope: CacheControlScope) on FIELD_DEFINITION | OBJECT | INTERFACE`
	if dd.String() != expectedOutput {
		t.Errorf("wrong output,expected=%q, got=%q %d/%d", expectedOutput, dd.String(), len(expectedOutput), len(dd.String()))
	}
}

func TestSelectionSet(t *testing.T) {

	ss := SelectionSet{}
	ss.Kind = SELECTION_SET
	ss.Selections = []Selection{
		&Field{
			Kind:  FIELD,
			Alias: &Name{Kind: NAME, Value: "test1"},
			Name:  &Name{Kind: NAME, Value: "Test"},
			Arguments: []*Argument{
				{Kind: ARGUMENT, Name: &Name{Kind: NAME, Value: "cache"}, Value: &BooleanValue{Kind: BOOLEAN_VALUE, Value: true}},
			},
			Directives: []*Directive{
				{
					Kind: DIRECTIVE,
					Name: &Name{Kind: NAME, Value: "skip"},
					Arguments: []*Argument{
						{Kind: ARGUMENT, Name: &Name{Kind: NAME, Value: "caching"}, Value: &BooleanValue{Kind: BOOLEAN_VALUE, Value: true}},
					},
				},
			},
		},
	}

	expectedOutput := `{
  test1: Test(cache: true)  @skip(caching: true)
}`

	if ss.String() != expectedOutput {
		t.Errorf("wrong output,expected=%q, got=%q %d/%d", expectedOutput, ss.String(), len(expectedOutput), len(ss.String()))
	}

}

func TestOperationDefinition(t *testing.T) {

	od := OperationDefinition{}
	od.Kind = OPERATION_DEFINITION
	od.OperationType = OperationTypeQuery
	od.Name = &Name{Kind: NAME, Value: "Test"}
	od.VariablesDefinition = nil
	od.VariablesDefinition = []*VariableDefinition{
		{Kind: VARIABLE, Variable: &Variable{Kind: VARIABLE, Name: &Name{Kind: NAME, Value: "status"}}, Type: &NamedType{Kind: NAMED_TYPE, Name: &Name{Kind: NAME, Value: "String"}}, DefaultValue: &StringValue{Kind: STRING_VALUE, Value: "Active"}},
		{
			Kind:         VARIABLE,
			Variable:     &Variable{Kind: VARIABLE, Name: &Name{Kind: NAME, Value: "point"}},
			Type:         &NamedType{Kind: NAMED_TYPE, Name: &Name{Kind: NAME, Value: "Int"}},
			DefaultValue: &IntValue{Kind: STRING_VALUE, Value: "0"},
			//Directives:   []*Directive{{Kind: DIRECTIVE, Name: &Name{Kind: NAME, Value: "skip"}, Arguments: []*Argument{{Kind: ARGUMENT, Name: &Name{Kind: NAME, Value: "caching"}, Value: &BooleanValue{Kind: BOOLEAN_VALUE, Value: true}}}}},
		},
	}

	od.Directives = []*Directive{
		{Kind: DIRECTIVE, Name: &Name{Kind: NAME, Value: "skip"}, Arguments: []*Argument{{Kind: ARGUMENT, Name: &Name{Kind: NAME, Value: "caching"}, Value: &BooleanValue{Kind: BOOLEAN_VALUE, Value: true}}}},
	}
	od.SelectionSet = &SelectionSet{Kind: SELECTION_SET, Selections: []Selection{
		&Field{
			Kind: FIELD,
			Name: &Name{Kind: NAME, Value: "status"},
			Arguments: []*Argument{
				{Kind: ARGUMENT, Name: &Name{Kind: NAME, Value: "status"}, Value: &StringValue{Kind: NAME, Value: "LiftStatus"}},
			},
		},
	}}

	expectedOutput := `query Test($status: String = "Active", $point: Int = 0)  @skip(caching: true) {
  status(status: LiftStatus)
}`

	if od.String() != expectedOutput {
		t.Errorf("wrong output,expected=%q, got=%q %d/%d", expectedOutput, od.String(), len(expectedOutput), len(od.String()))
	}
}

func TestFragmentDefinition(t *testing.T) {

	frgd := FragmentDefinition{}
	frgd.Kind = FRAGMENT_DEFINITION
	frgd.Operation = ""
	frgd.FragmentName = &Name{Kind: NAME, Value: "friendFields"}
	frgd.TypeCondition = &NamedType{Kind: NAMED_TYPE, Name: &Name{Kind: NAME, Value: "User"}}
	// frgd.Directives = []*Directive{
	// 	{Kind: DIRECTIVE, Name: &Name{Kind: NAME, Value: "skip"}, Arguments: []*Argument{{Kind: ARGUMENT, Name: &Name{Kind: NAME, Value: "caching"}, Value: &BooleanValue{Kind: BOOLEAN_VALUE, Value: true}}}},
	// }
	frgd.SelectionSet = &SelectionSet{Kind: SELECTION_SET, Selections: []Selection{
		&Field{Kind: FIELD, Name: &Name{Kind: NAME, Value: "id"}},
		&Field{Kind: FIELD, Name: &Name{Kind: NAME, Value: "name"}},
		&Field{
			Kind: FIELD,
			Name: &Name{Kind: NAME, Value: "profilePic"},
			Arguments: []*Argument{
				{Kind: ARGUMENT, Name: &Name{Kind: NAME, Value: "size"}, Value: &IntValue{Kind: INT_VALUE, Value: "50"}},
			},
		},
	}}
	expectedOutput := `fragment friendFields on User {
  id
  name
  profilePic(size: 50)
}`

	if frgd.String() != expectedOutput {
		t.Errorf("wrong output,expected=%q, got=%q %d/%d", expectedOutput, frgd.String(), len(expectedOutput), len(frgd.String()))
	}
}

func TestOperationFragmentSpread(t *testing.T) {

	od := OperationDefinition{}
	od.Kind = OPERATION_DEFINITION
	od.OperationType = OperationTypeQuery
	od.Name = &Name{Kind: NAME, Value: "withFragments"}
	od.VariablesDefinition = nil

	od.SelectionSet = &SelectionSet{Kind: SELECTION_SET, Selections: []Selection{
		&Field{
			Kind: FIELD,
			Name: &Name{Kind: NAME, Value: "user"},
			Arguments: []*Argument{
				{Kind: ARGUMENT, Name: &Name{Kind: NAME, Value: "id"}, Value: &IntValue{Kind: NAME, Value: "4"}},
			},

			SelectionSet: &SelectionSet{Kind: SELECTION_SET, Selections: []Selection{
				&Field{
					Kind: FIELD,
					Name: &Name{Kind: NAME, Value: "friends"},
					Arguments: []*Argument{
						{Kind: ARGUMENT, Name: &Name{Kind: NAME, Value: "first"}, Value: &IntValue{Kind: NAME, Value: "10"}},
					},
					SelectionSet: &SelectionSet{Kind: SELECTION_SET, Selections: []Selection{
						&FragmentSpread{Kind: FRAGMENT_SPREAD, FragmentName: &Name{Kind: NAME, Value: "friendFields"}},
					}},
				},
				&Field{
					Kind: FIELD,
					Name: &Name{Kind: NAME, Value: "mutualFriends"},
					Arguments: []*Argument{
						{Kind: ARGUMENT, Name: &Name{Kind: NAME, Value: "first"}, Value: &IntValue{Kind: NAME, Value: "10"}},
					},
					SelectionSet: &SelectionSet{Kind: SELECTION_SET, Selections: []Selection{
						&FragmentSpread{Kind: FRAGMENT_SPREAD, FragmentName: &Name{Kind: NAME, Value: "friendFields"}, Directives: nil},
					}},
				},
			}},
		},
	}}

	expectedOutput := `query withFragments {
  user(id: 4) {
    friends(first: 10) {
      ...friendFields
    }
    mutualFriends(first: 10) {
      ...friendFields
    }
  }
}`

	if od.String() != expectedOutput {
		t.Errorf("wrong output,expected=%q, got=%q %d/%d", expectedOutput, od.String(), len(expectedOutput), len(od.String()))
	}
}

func TestOperationInlineFragmentSpread(t *testing.T) {

	od := OperationDefinition{}
	od.Kind = OPERATION_DEFINITION
	od.OperationType = OperationTypeQuery
	od.Name = &Name{Kind: NAME, Value: "inlineFragmentNoType"}
	od.VariablesDefinition = []*VariableDefinition{
		{
			Kind:     VARIABLE_DEFINITION,
			Variable: &Variable{Kind: VARIABLE, Name: &Name{Kind: NAME, Value: "expandedInfo"}},
			Type:     &NamedType{Kind: NAMED_TYPE, Name: &Name{Kind: NAME, Value: "Boolean"}},
			//DefaultValue: &StringValue{Kind: STRING_VALUE, Value: "Active"}
		},
	}

	od.SelectionSet = &SelectionSet{Kind: SELECTION_SET, Selections: []Selection{
		&Field{
			Kind: FIELD,
			Name: &Name{Kind: NAME, Value: "user"},
			Arguments: []*Argument{
				{Kind: ARGUMENT, Name: &Name{Kind: NAME, Value: "handle"}, Value: &StringValue{Kind: STRING_VALUE, Value: "zuck"}},
			},

			SelectionSet: &SelectionSet{Kind: SELECTION_SET, Selections: []Selection{
				&Field{
					Kind: FIELD,
					Name: &Name{Kind: NAME, Value: "id"},
				},
				&Field{
					Kind: FIELD,
					Name: &Name{Kind: NAME, Value: "name"},
				},
				&InlineFragment{
					Kind:          INLINE_FRAGMENT,
					TypeCondition: nil,
					Directives: []*Directive{
						{Kind: DIRECTIVE, Name: &Name{Kind: NAME, Value: "include"}, Arguments: []*Argument{
							{Kind: ARGUMENT, Name: &Name{Kind: NAME, Value: "if"}, Value: &Variable{Kind: VARIABLE, Name: &Name{Kind: NAME, Value: "expandedInfo"}}},
						}}},
					SelectionSet: &SelectionSet{Kind: SELECTION_SET, Selections: []Selection{

						&Field{Kind: FIELD, Name: &Name{Kind: NAME, Value: "firstName"}},
						&Field{Kind: FIELD, Name: &Name{Kind: NAME, Value: "lastName"}},
						&Field{Kind: FIELD, Name: &Name{Kind: NAME, Value: "birthday"}},
					}},
				},
			}},
		},
	}}

	expectedOutput := `query inlineFragmentNoType($expandedInfo: Boolean) {
  user(handle: "zuck") {
    id
    name
    ...  @include(if: $expandedInfo) {
      firstName
      lastName
      birthday
    }
  }
}`

	if od.String() != expectedOutput {
		t.Errorf("wrong output,expected=%q, got=%q %d/%d", expectedOutput, od.String(), len(expectedOutput), len(od.String()))
	}
}

func TestOperationInlineFragmentSpread2(t *testing.T) {

	od := OperationDefinition{}
	od.Kind = OPERATION_DEFINITION
	od.OperationType = OperationTypeQuery
	od.Name = &Name{Kind: NAME, Value: "inlineFragmentTyping"}

	od.SelectionSet = &SelectionSet{Kind: SELECTION_SET, Selections: []Selection{
		&Field{
			Kind: FIELD,
			Name: &Name{Kind: NAME, Value: "profiles"},
			Arguments: []*Argument{
				{Kind: ARGUMENT, Name: &Name{Kind: NAME, Value: "handles"}, Value: &ListValue{Kind: LIST_VALUE, Values: []Value{
					&StringValue{Kind: STRING_VALUE, Value: "zuck"},
					&StringValue{Kind: STRING_VALUE, Value: "coca-cola"},
				}}},
			},

			SelectionSet: &SelectionSet{Kind: SELECTION_SET, Selections: []Selection{
				&Field{
					Kind: FIELD,
					Name: &Name{Kind: NAME, Value: "handle"},
				},

				&InlineFragment{
					Kind:          INLINE_FRAGMENT,
					TypeCondition: &NamedType{Kind: NAMED_TYPE, Name: &Name{Kind: NAME, Value: "User"}},
					Directives:    nil,
					SelectionSet: &SelectionSet{Kind: SELECTION_SET, Selections: []Selection{

						&Field{
							Kind: FIELD,
							Name: &Name{Kind: NAME, Value: "friends"},
							SelectionSet: &SelectionSet{Kind: SELECTION_SET, Selections: []Selection{
								&Field{Kind: FIELD, Name: &Name{Kind: NAME, Value: "count"}},
							}},
						},
					}},
				},
				&InlineFragment{
					Kind:          INLINE_FRAGMENT,
					TypeCondition: &NamedType{Kind: NAMED_TYPE, Name: &Name{Kind: NAME, Value: "Page"}},
					Directives:    nil,
					SelectionSet: &SelectionSet{Kind: SELECTION_SET, Selections: []Selection{

						&Field{
							Kind: FIELD,
							Name: &Name{Kind: NAME, Value: "likers"},
							SelectionSet: &SelectionSet{Kind: SELECTION_SET, Selections: []Selection{
								&Field{Kind: FIELD, Name: &Name{Kind: NAME, Value: "count"}},
							}},
						},
					}},
				},
			}},
		},
	}}

	expectedOutput := `query inlineFragmentTyping {
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

	if od.String() != expectedOutput {
		t.Errorf("wrong output,expected=%q, got=%q %d/%d", expectedOutput, od.String(), len(expectedOutput), len(od.String()))
	}
}

func TestSchemaDefinition(t *testing.T) {

	scmd := SchemaDefinition{}
	scmd.Kind = SCHEMA_DEFINITION
	scmd.Description = nil //&StringValue{Kind: STRING_VALUE, Value: "test"}
	// scmd.Directives = []*Directive{
	// 	{Kind: DIRECTIVE, Name: &Name{Kind: NAME, Value: "include"}, Arguments: []*Argument{
	// 		{Kind: ARGUMENT, Name: &Name{Kind: NAME, Value: "if"}, Value: &Variable{Kind: VARIABLE, Name: &Name{Kind: NAME, Value: "expandedInfo"}}},
	// 	}},
	// }
	scmd.OperationTypes = []*RootOperationTypeDefinition{
		{Kind: ROOT_OPERATION_TYPE_DEFINITION, OperationType: OperationTypeQuery, NamedType: &NamedType{Kind: NAMED_TYPE, Name: &Name{Kind: NAME, Value: "MyQueryRootType"}}},
		{Kind: ROOT_OPERATION_TYPE_DEFINITION, OperationType: OperationTypeMutation, NamedType: &NamedType{Kind: NAMED_TYPE, Name: &Name{Kind: NAME, Value: "MyMutationRootType"}}},
	}

	expectedOutput := `schema {
  query: MyQueryRootType
  mutation: MyMutationRootType
}`

	if scmd.String() != expectedOutput {
		t.Errorf("wrong output,expected=%q, got=%q %d/%d", expectedOutput, scmd.String(), len(expectedOutput), len(scmd.String()))
	}
}
