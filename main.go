package main

import (
	"fmt"
	"os"

	"github.com/mateors/graphqlparser/ast"
	"github.com/mateors/graphqlparser/lexer"
	"github.com/mateors/graphqlparser/parser"
	"github.com/mateors/graphqlparser/token"
)

// returns lower-case ch iff ch is ASCII letter
func lower(ch byte) byte {
	return ('a' - 'A') | ch
}

type Lift struct {
	Name string
	Size int
}

func manualTest() {
	// fmt.Println(lower('S'))
	// fmt.Println(32 | 83)

	// vals := []string{"NamedEntity", "ValuedEntity"}
	// var infcs string
	// for _, iname := range vals {
	// 	//iname := inf.Value
	// 	infcs += fmt.Sprintf("%s & ", iname)
	// }
	// infcs = strings.TrimRight(infcs, " & ")
	// fmt.Println(infcs)

	// maps := map[string]string{
	// 	{"key": "name"},
	// }

	// evd := ast.EnumValueDefinition{}
	// evd.Kind = ast.ENUMVALUE_DEFINITION
	// evd.Description = nil
	// evd.Name = &ast.Name{Kind: ast.NAME, Value: "NORTH"}
	// evd.Directives = nil
	// fmt.Println(evd.String())

	// ss := ast.SelectionSet{}
	// ss.Kind = ast.SELECTION_SET
	// ss.Selections = []ast.Selection{
	// 	&ast.Field{
	// 		Kind:  ast.FIELD,
	// 		Alias: &ast.Name{Kind: ast.NAME, Value: "test1"},
	// 		Name:  &ast.Name{Kind: ast.NAME, Value: "Test"},
	// 		Arguments: []*ast.Argument{
	// 			{Kind: ast.ARGUMENT, Name: &ast.Name{Kind: ast.NAME, Value: "cache"}, Value: &ast.BooleanValue{Kind: ast.BOOLEAN_VALUE, Value: true}},
	// 		},
	// 		Directives: []*ast.Directive{
	// 			{
	// 				Kind: ast.DIRECTIVE,
	// 				Name: &ast.Name{Kind: ast.NAME, Value: "skip"},
	// 				Arguments: []*ast.Argument{
	// 					{Kind: ast.ARGUMENT, Name: &ast.Name{Kind: ast.NAME, Value: "caching"}, Value: &ast.BooleanValue{Kind: ast.BOOLEAN_VALUE, Value: true}},
	// 				},
	// 			},
	// 		},
	// 	},
	// }

	// frgd := ast.FragmentDefinition{}
	// frgd.Kind = ast.FRAGMENT_DEFINITION
	// frgd.Operation = ""
	// frgd.FragmentName = &ast.Name{Kind: ast.NAME, Value: "friendFields"}
	// frgd.TypeCondition = &ast.NamedType{Kind: ast.NAMED_TYPE, Name: &ast.Name{Kind: ast.NAME, Value: "User"}}
	// frgd.Directives = []*ast.Directive{
	// 	{Kind: ast.DIRECTIVE, Name: &ast.Name{Kind: ast.NAME, Value: "skip"}, Arguments: []*ast.Argument{{Kind: ast.ARGUMENT, Name: &ast.Name{Kind: ast.NAME, Value: "caching"}, Value: &ast.BooleanValue{Kind: ast.BOOLEAN_VALUE, Value: true}}}},
	// }
	// frgd.SelectionSet = &ast.SelectionSet{Kind: ast.SELECTION_SET, Selections: []ast.Selection{
	// 	&ast.Field{Kind: ast.FIELD, Name: &ast.Name{Kind: ast.NAME, Value: "id"}},
	// 	&ast.Field{Kind: ast.FIELD, Name: &ast.Name{Kind: ast.NAME, Value: "name"}},
	// 	&ast.Field{
	// 		Kind: ast.FIELD,
	// 		Name: &ast.Name{Kind: ast.NAME, Value: "profilePic"},
	// 		Arguments: []*ast.Argument{
	// 			{Kind: ast.ARGUMENT, Name: &ast.Name{Kind: ast.NAME, Value: "size"}, Value: &ast.IntValue{Kind: ast.INT_VALUE, Value: "50"}},
	// 		},
	// 	},
	// }}

	// fmt.Println(frgd.String())
	// os.Exit(1)

	scmd := ast.SchemaDefinition{}
	scmd.Kind = ast.SCHEMA_DEFINITION
	scmd.Description = &ast.StringValue{Kind: ast.STRING_VALUE, Value: "test"}
	scmd.Directives = []*ast.Directive{
		{Kind: ast.DIRECTIVE, Name: &ast.Name{Kind: ast.NAME, Value: "include"}, Arguments: []*ast.Argument{
			{Kind: ast.ARGUMENT, Name: &ast.Name{Kind: ast.NAME, Value: "if"}, Value: &ast.Variable{Kind: ast.VARIABLE, Name: &ast.Name{Kind: ast.NAME, Value: "expandedInfo"}}},
		}},
	}
	scmd.OperationTypes = []*ast.RootOperationTypeDefinition{
		{Kind: ast.ROOT_OPERATION_TYPE_DEFINITION, OperationType: ast.OperationTypeQuery, NamedType: &ast.NamedType{Kind: ast.NAMED_TYPE, Name: &ast.Name{Kind: ast.NAME, Value: "MyQueryRootType"}}},
		{Kind: ast.ROOT_OPERATION_TYPE_DEFINITION, OperationType: ast.OperationTypeMutation, NamedType: &ast.NamedType{Kind: ast.NAMED_TYPE, Name: &ast.Name{Kind: ast.NAME, Value: "MyMutationRootType"}}},
	}
	fmt.Println(scmd.String())
	os.Exit(1)

	inf := ast.InlineFragment{}
	inf.Kind = ast.INLINE_FRAGMENT
	//inf.TypeCondition = &ast.NamedType{Kind: ast.NAMED_TYPE, Name: &ast.Name{Kind: ast.NAME, Value: "User"}}
	inf.Directives = []*ast.Directive{
		{Kind: ast.DIRECTIVE, Name: &ast.Name{Kind: ast.NAME, Value: "include"}, Arguments: []*ast.Argument{
			{Kind: ast.ARGUMENT, Name: &ast.Name{Kind: ast.NAME, Value: "if"}, Value: &ast.Variable{Kind: ast.VARIABLE, Name: &ast.Name{Kind: ast.NAME, Value: "expandedInfo"}}},
		}},
	}
	inf.SelectionSet = &ast.SelectionSet{Kind: ast.SELECTION_SET, Selections: []ast.Selection{
		&ast.Field{
			Kind:      ast.FIELD,
			Name:      &ast.Name{Kind: ast.NAME, Value: "friends"},
			Arguments: nil,
			SelectionSet: &ast.SelectionSet{Kind: ast.SELECTION_SET, Selections: []ast.Selection{
				&ast.Field{Kind: ast.FIELD, Name: &ast.Name{Kind: ast.NAME, Value: "count"}}},
			}},
	},
	}

	fmt.Println(inf.String())
	os.Exit(1)

	od := ast.OperationDefinition{}
	od.Kind = ast.OPERATION_DEFINITION
	od.OperationType = ast.OperationTypeQuery
	od.Name = &ast.Name{Kind: ast.NAME, Value: "withFragments"}
	od.VariablesDefinition = nil
	// od.VariablesDefinition = []*ast.VariableDefinition{
	// 	{Kind: ast.VARIABLE_DEFINITION, Variable: &ast.Variable{Kind: ast.VARIABLE, Name: &ast.Name{Kind: ast.NAME, Value: "status"}}, Type: &ast.NamedType{Kind: ast.NAMED_TYPE, Name: &ast.Name{Kind: ast.NAME, Value: "String"}}, DefaultValue: &ast.StringValue{Kind: ast.STRING_VALUE, Value: "Active"}},
	// 	{
	// 		Kind:         ast.VARIABLE,
	// 		Variable:     &ast.Variable{Kind: ast.VARIABLE, Name: &ast.Name{Kind: ast.NAME, Value: "point"}},
	// 		Type:         &ast.NamedType{Kind: ast.NAMED_TYPE, Name: &ast.Name{Kind: ast.NAME, Value: "Int"}},
	// 		DefaultValue: &ast.IntValue{Kind: ast.STRING_VALUE, Value: "0"},
	// 		//Directives:   []*ast.Directive{{Kind: ast.DIRECTIVE, Name: &ast.Name{Kind: ast.NAME, Value: "skip"}, Arguments: []*ast.Argument{{Kind: ast.ARGUMENT, Name: &ast.Name{Kind: ast.NAME, Value: "caching"}, Value: &ast.BooleanValue{Kind: ast.BOOLEAN_VALUE, Value: true}}}}},
	// 	},
	// }

	// od.Directives = []*ast.Directive{
	// 	{Kind: ast.DIRECTIVE, Name: &ast.Name{Kind: ast.NAME, Value: "skip"}, Arguments: []*ast.Argument{{Kind: ast.ARGUMENT, Name: &ast.Name{Kind: ast.NAME, Value: "caching"}, Value: &ast.BooleanValue{Kind: ast.BOOLEAN_VALUE, Value: true}}}},
	// }
	od.SelectionSet = &ast.SelectionSet{Kind: ast.SELECTION_SET, Selections: []ast.Selection{
		&ast.Field{
			Kind: ast.FIELD,
			Name: &ast.Name{Kind: ast.NAME, Value: "user"},
			Arguments: []*ast.Argument{
				{Kind: ast.ARGUMENT, Name: &ast.Name{Kind: ast.NAME, Value: "id"}, Value: &ast.IntValue{Kind: ast.NAME, Value: "4"}},
			},

			SelectionSet: &ast.SelectionSet{Kind: ast.SELECTION_SET, Selections: []ast.Selection{
				&ast.Field{
					Kind: ast.FIELD,
					Name: &ast.Name{Kind: ast.NAME, Value: "friends"},
					Arguments: []*ast.Argument{
						{Kind: ast.ARGUMENT, Name: &ast.Name{Kind: ast.NAME, Value: "first"}, Value: &ast.IntValue{Kind: ast.NAME, Value: "10"}},
					},
					SelectionSet: &ast.SelectionSet{Kind: ast.SELECTION_SET, Selections: []ast.Selection{
						&ast.FragmentSpread{Kind: ast.FRAGMENT_SPREAD, FragmentName: &ast.Name{Kind: ast.NAME, Value: "friendFields"}},
					}},
				},
				&ast.Field{
					Kind: ast.FIELD,
					Name: &ast.Name{Kind: ast.NAME, Value: "mutualFriends"},
					Arguments: []*ast.Argument{
						{Kind: ast.ARGUMENT, Name: &ast.Name{Kind: ast.NAME, Value: "first"}, Value: &ast.IntValue{Kind: ast.NAME, Value: "10"}},
					},
					SelectionSet: &ast.SelectionSet{Kind: ast.SELECTION_SET, Selections: []ast.Selection{
						&ast.FragmentSpread{Kind: ast.FRAGMENT_SPREAD, FragmentName: &ast.Name{Kind: ast.NAME, Value: "friendFields"}, Directives: nil},
					}},
				},
			}},
		},
	}}

	fmt.Println(od.String())
	os.Exit(1)

	dd := ast.DirectiveDefinition{}
	dd.Kind = ast.DIRECTIVE_DEFINITION
	dd.Description = &ast.StringValue{Kind: ast.STRING_VALUE, Value: "test"}
	dd.Name = &ast.Name{Kind: ast.NAME, Value: "cacheControl"}
	dd.Arguments = []*ast.InputValueDefinition{
		{Kind: ast.INPUT_VALUE_DEFINITION, Name: &ast.Name{Kind: ast.NAME, Value: "maxAge"}, Type: &ast.NamedType{Kind: ast.NAMED_TYPE, Name: &ast.Name{Kind: ast.NAME, Value: "Int"}}},
		{Kind: ast.INPUT_VALUE_DEFINITION, Name: &ast.Name{Kind: ast.NAME, Value: "scope"}, Type: &ast.NamedType{Kind: ast.NAMED_TYPE, Name: &ast.Name{Kind: ast.NAME, Value: "CacheControlScope"}}},
	}
	dd.Locations = []*ast.Name{
		{Kind: ast.NAME, Value: "FIELD_DEFINITION"},
		{Kind: ast.NAME, Value: "OBJECT"},
		{Kind: ast.NAME, Value: "INTERFACE"},
	}

	// &ast.InputValueDefinition{
	// 	Name:         &ast.Name{Kind: ast.NAME, Token: token.Token{}, Value: "unit"},
	// 	Type:         &ast.NamedType{Kind: ast.NAMED_TYPE, Token: token.Token{}, Name: &ast.Name{Kind: "Name", Token: token.Token{}, Value: "LengthUnit"}},
	// 	DefaultValue: &ast.StringValue{Kind: ast.ENUM_VALUE, Token: token.Token{}, Value: "METER"},
	// 	Directives:   directives,
	// }
	fmt.Println(dd.String())
	os.Exit(1)

	scd := ast.ScalarDefinition{Kind: ast.SCALAR_DEFINITION}
	scd.Description = &ast.StringValue{Kind: ast.STRING_VALUE, Value: "test"}
	scd.Name = &ast.Name{Kind: ast.NAME, Value: "Upload"}
	scd.Directives = []*ast.Directive{
		{Kind: ast.DIRECTIVE, Name: &ast.Name{Kind: ast.NAME, Value: "skip"}, Arguments: []*ast.Argument{
			{Kind: ast.ARGUMENT, Name: &ast.Name{Kind: ast.NAME, Value: "caching"}, Value: &ast.BooleanValue{Kind: ast.BOOLEAN_VALUE, Value: true}},
		}},
	}

	fmt.Println(scd.String())
	os.Exit(1)

	iod := ast.InputObjectDefinition{}
	iod.Kind = ast.INPUT_OBJECT_DEFINITION
	//iod.Description = &ast.StringValue{Kind: ast.STRING_VALUE, Value: "test"}
	iod.Name = &ast.Name{Kind: ast.NAME, Value: "Client"}
	iod.Directives = []*ast.Directive{
		{Kind: ast.DIRECTIVE, Name: &ast.Name{Kind: ast.NAME, Value: "skip"}, Arguments: []*ast.Argument{
			{Kind: ast.ARGUMENT, Name: &ast.Name{Kind: ast.NAME, Value: "caching"}, Value: &ast.BooleanValue{Kind: ast.BOOLEAN_VALUE, Value: true}},
		}},
	}
	iod.Fields = []*ast.InputValueDefinition{
		{Kind: ast.INPUT_VALUE_DEFINITION, Name: &ast.Name{Kind: ast.NAME, Value: "name"}, Type: &ast.NonNullType{Kind: ast.STRING_VALUE, Type: &ast.Name{Kind: ast.NAME, Value: "String"}}},
		{DefaultValue: &ast.IntValue{Kind: ast.INT_VALUE, Value: "10"}, Kind: ast.INPUT_VALUE_DEFINITION, Name: &ast.Name{Kind: ast.NAME, Value: "age"}, Type: &ast.NamedType{Kind: ast.STRING_VALUE, Name: &ast.Name{Kind: ast.NAME, Value: "Int"}}},
	}

	fmt.Println(iod.String())
	os.Exit(1)

	ed := ast.EnumDefinition{}
	ed.Kind = ast.ENUM_DEFINITION
	ed.Description = &ast.StringValue{Kind: ast.STRING_VALUE, Value: ""}
	ed.Name = &ast.Name{Kind: ast.NAME, Value: "Country"}
	ed.Directives = []*ast.Directive{
		{Kind: ast.DIRECTIVE, Name: &ast.Name{Kind: ast.NAME, Value: "skip"}, Arguments: []*ast.Argument{
			{Kind: ast.ARGUMENT, Name: &ast.Name{Kind: ast.NAME, Value: "caching"}, Value: &ast.BooleanValue{Kind: ast.BOOLEAN_VALUE, Value: true}},
		}},
	}
	ed.Values = []*ast.EnumValueDefinition{
		{Kind: ast.ENUMVALUE_DEFINITION, Name: &ast.Name{Kind: ast.NAME, Value: "BANGLADESH"}},
		{Kind: ast.ENUMVALUE_DEFINITION, Name: &ast.Name{Kind: ast.NAME, Value: "INDIA"}},
	}
	fmt.Println(ed.String())
	os.Exit(1)

	ud := ast.UnionDefinition{}
	ud.Kind = ast.UNION_DEFINITION
	ud.Description = nil //&ast.StringValue{Kind: ast.STRING_VALUE, Value: "Test des"}
	ud.Name = &ast.Name{Kind: ast.NAME, Value: "SearchResult"}
	// ud.Directives = nil []*ast.Directive{
	// 	{Kind: ast.DIRECTIVE, Name: &ast.Name{Kind: ast.NAME, Value: "skip"}, Arguments: []*ast.Argument{
	// 		{Kind: ast.ARGUMENT, Name: &ast.Name{Kind: ast.NAME, Value: "caching"}, Value: &ast.BooleanValue{Kind: ast.BOOLEAN_VALUE, Value: true}},
	// 	}},
	// }
	ud.UnionMemberTypes = []*ast.NamedType{
		{Kind: ast.NAMED_TYPE, Name: &ast.Name{Kind: ast.NAME, Value: "Photo"}},
		{Kind: ast.NAMED_TYPE, Name: &ast.Name{Kind: ast.NAME, Value: "Person"}},
	}

	//union SearchResult = Photo | Person

	fmt.Println(ud.String())

	os.Exit(1)

	id := &ast.InterfaceDefinition{}
	id.Description = &ast.StringValue{Kind: ast.STRING_VALUE, Value: "Description test"}
	id.Name = &ast.Name{Kind: ast.NAME, Value: "NamedEntity"}

	infcs1 := []*ast.NamedType{}
	infcs1 = append(infcs1, &ast.NamedType{Kind: ast.NAMED_TYPE, Name: &ast.Name{Kind: ast.NAME, Value: "Abs"}})
	infcs1 = append(infcs1, &ast.NamedType{Kind: ast.NAMED_TYPE, Name: &ast.Name{Kind: ast.NAME, Value: "Book"}})
	id.Interfaces = nil

	idirectives := []*ast.Directive{}
	idirectives = append(idirectives, &ast.Directive{
		Kind:  ast.DIRECTIVE,
		Token: token.Token{},
		Name:  &ast.Name{Kind: ast.NAME, Value: "addExternalFields"},
		Arguments: []*ast.Argument{
			{
				Kind:  ast.ARGUMENT,
				Token: token.Token{},
				Name:  &ast.Name{Kind: ast.NAME, Value: "name"},
				Value: &ast.StringValue{Kind: ast.STRING_VALUE, Value: "photo"},
			}, {
				Kind:  ast.ARGUMENT,
				Token: token.Token{},
				Name:  &ast.Name{Kind: ast.NAME, Value: "cache"},
				Value: &ast.BooleanValue{Kind: ast.BOOLEAN_VALUE, Value: true},
			}},
	})
	id.Directives = idirectives

	fieldi := &ast.FieldDefinition{}
	fieldi.Kind = ast.FIELD_DEFINITION
	fieldi.Name = &ast.Name{Kind: ast.NAME, Value: "name"}
	fieldi.Type = &ast.NamedType{Kind: ast.NAMED_TYPE, Name: &ast.Name{Kind: ast.NAME, Value: "String"}}

	fieldi2 := &ast.FieldDefinition{}
	fieldi2.Kind = ast.FIELD_DEFINITION
	fieldi2.Name = &ast.Name{Kind: ast.NAME, Value: "value"}
	fieldi2.Type = &ast.NamedType{Kind: ast.NAMED_TYPE, Name: &ast.Name{Kind: ast.NAME, Value: "Int"}}

	ifields := []*ast.FieldDefinition{}
	ifields = append(ifields, fieldi)
	ifields = append(ifields, fieldi2)

	id.Fields = ifields
	fmt.Println(id.String())
	os.Exit(1)

	field := &ast.FieldDefinition{}
	field.Kind = ast.FIELD_DEFINITION
	field.Name = &ast.Name{Kind: ast.NAME, Value: "name"}
	field.Type = &ast.NamedType{Kind: ast.NAMED_TYPE, Name: &ast.Name{Kind: ast.NAME, Value: "String"}}

	dfields := []*ast.FieldDefinition{}
	dfields = append(dfields, field)

	field2 := &ast.FieldDefinition{}
	field2.Kind = ast.FIELD_DEFINITION
	field2.Name = &ast.Name{Kind: ast.NAME, Value: "age"}
	//field2.Type = &ast.NonNullType{Kind: ast.NAMED_TYPE, Type: &ast.Name{Kind: ast.NAME, Value: "Int"} }
	field2.Type = &ast.NonNullType{Kind: ast.NONNULL_TYPE, Type: &ast.NamedType{Kind: ast.NAMED_TYPE, Name: &ast.Name{Kind: ast.NAME, Value: "Int"}}}
	dfields = append(dfields, field2)

	infcs := []*ast.NamedType{}
	infcs = append(infcs, &ast.NamedType{Kind: ast.NAMED_TYPE, Name: &ast.Name{Kind: ast.NAME, Value: "Abs"}})
	infcs = append(infcs, &ast.NamedType{Kind: ast.NAMED_TYPE, Name: &ast.Name{Kind: ast.NAME, Value: "Book"}})

	var obj ast.ObjectDefinition
	obj.Kind = ast.OBJECT_DEFINITION
	obj.Name = &ast.Name{Kind: ast.NAME, Value: "Lift"}
	obj.Fields = dfields
	obj.Interfaces = infcs
	fmt.Println(">>", obj.String())
	os.Exit(1)

	//var ttype ast.Type
	iv := ast.InputValueDefinition{}
	iv.Name = &ast.Name{Kind: "Name", Token: token.Token{}, Value: "name"}
	//aa := &ast.NamedType{Kind: ast.NAMED_TYPE, Token: token.Token{}, Name: &ast.Name{Kind: "Name", Token: token.Token{}, Value: "String"}}
	//aa := &ast.NonNullType{Kind: ast.NONNULL_TYPE, Token: token.Token{}, Type: &ast.NamedType{Kind: "Name", Token: token.Token{}, Name: &ast.Name{Kind: "Name", Token: token.Token{}, Value: "String"}}}
	//aa := &ast.ListType{Kind: ast.LIST_TYPE, Token: token.Token{}, Type: &ast.NamedType{Kind: "Name", Token: token.Token{}, Name: &ast.Name{Kind: "Name", Token: token.Token{}, Value: "String"}}}
	aa := &ast.NonNullType{Kind: ast.NONNULL_TYPE, Token: token.Token{}, Type: &ast.ListType{Kind: ast.LIST_TYPE, Token: token.Token{}, Type: &ast.NonNullType{Kind: ast.NONNULL_TYPE, Token: token.Token{}, Type: &ast.NamedType{Kind: "Name", Token: token.Token{}, Name: &ast.Name{Kind: "Name", Token: token.Token{}, Value: "String"}}}}}

	iv.Type = aa // &ast.StringValue{Kind: ast.STRING_VALUE, Token: token.Token{}, Value: "test"}
	//ast.NonNullType

	//fmt.Println(iv.String())

	fd := ast.FieldDefinition{}
	fd.Kind = ast.FIELD_DEFINITION
	fd.Token = token.Token{}
	fd.Name = &ast.Name{Kind: ast.NAME, Token: token.Token{}, Value: "length"}
	fd.Type = &ast.NamedType{Kind: ast.NAMED_TYPE, Token: token.Token{}, Name: &ast.Name{Kind: ast.NAME, Token: token.Token{}, Value: "Float"}}
	ivd := []*ast.InputValueDefinition{}

	args := []*ast.Argument{}
	args = append(args, &ast.Argument{Kind: ast.ARGUMENT, Name: &ast.Name{Kind: ast.NAME, Value: "name"}, Value: &ast.StringValue{Kind: ast.STRING_VALUE, Value: "photo"}})

	args = append(args, &ast.Argument{Kind: ast.ARGUMENT, Name: &ast.Name{Kind: ast.NAME, Value: "caching"}, Value: &ast.BooleanValue{Kind: ast.STRING_VALUE, Value: true}})
	//args = append(args, &ast.Argument{Kind: ast.ARGUMENT, Name: &ast.Name{Kind: ast.NAME, Value: "caching"}, Value: &ast.Variable{Kind: ast.VARIABLE, Name: &ast.Name{Kind: ast.NAME, Value: "isCaching"}}})

	//args = append(args, &ast.Argument{Kind: ast.ARGUMENT, Name: &ast.Name{Kind: ast.NAME, Value: "age"}, Value: &ast.IntValue{Kind: ast.INT_VALUE, Value: "50"}})
	//args = append(args, &ast.Argument{Kind: ast.ARGUMENT, Name: &ast.Name{Kind: ast.NAME, Value: "point"}, Value: &ast.FloatValue{Kind: ast.FLOAT_VALUE, Value: "500.45"}})

	vals := []ast.Value{}
	vals = append(vals, &ast.StringValue{Kind: ast.STRING_VALUE, Value: "Wania"})
	vals = append(vals, &ast.StringValue{Kind: ast.STRING_VALUE, Value: "Arisha"})
	args = append(args, &ast.Argument{Kind: ast.ARGUMENT, Name: &ast.Name{Kind: ast.NAME, Value: "siblings"}, Value: &ast.ListValue{Kind: ast.LIST_VALUE, Values: vals}})

	fields := []*ast.ObjectField{}
	fields = append(fields, &ast.ObjectField{Kind: ast.OBJECT_FIELD, Name: &ast.Name{Kind: ast.NAME, Value: "lat"}, Value: &ast.FloatValue{Kind: ast.FLOAT_VALUE, Value: "12.43"}})
	fields = append(fields, &ast.ObjectField{Kind: ast.OBJECT_FIELD, Name: &ast.Name{Kind: ast.NAME, Value: "long"}, Value: &ast.IntValue{Kind: ast.INT_VALUE, Value: "212"}})

	//vals = append(vals, &ast.ObjectValue{Kind: ast.OBJECT_VALUE, Fields: fields})

	//args = append(args, &ast.Argument{Kind: ast.ARGUMENT, Name: &ast.Name{Kind: ast.NAME, Value: "siblings"}, Value: &ast.ListValue{Kind: ast.LIST_VALUE, Values: vals}})
	args = append(args, &ast.Argument{Kind: ast.ARGUMENT, Name: &ast.Name{Kind: ast.NAME, Value: "location"}, Value: &ast.ObjectValue{Kind: ast.OBJECT_VALUE, Fields: fields}})

	directives := []*ast.Directive{}
	directives = append(directives, &ast.Directive{
		Kind:      ast.DIRECTIVE,
		Name:      &ast.Name{Kind: ast.NAME, Value: "excludeField"},
		Arguments: args,
	})
	// directives = append(directives, &ast.Directive{
	// 	Kind:      ast.DIRECTIVE,
	// 	Name:      &ast.Name{Kind: ast.NAME, Value: "skip"},
	// 	Arguments: args,
	// })

	iv1 := &ast.InputValueDefinition{
		Name:         &ast.Name{Kind: ast.NAME, Token: token.Token{}, Value: "unit"},
		Type:         &ast.NamedType{Kind: ast.NAMED_TYPE, Token: token.Token{}, Name: &ast.Name{Kind: "Name", Token: token.Token{}, Value: "LengthUnit"}},
		DefaultValue: &ast.StringValue{Kind: ast.ENUM_VALUE, Token: token.Token{}, Value: "METER"},
		Directives:   directives,
	}
	ivd = append(ivd, iv1)
	fd.Arguments = ivd

	fmt.Println(">>", fd.String())
	os.Exit(1)

	input := `type Person {
		id: ID!
		adult: Boolean!
		name: String!
		age: Int!
		salary: Float!
		length(unit: LengthUnit = METER): Float
		appearsIn: [Episode]!
	}

	interface Book {
		title: String!
		author: Author!
	}
	  
	type Textbook implements Book {
		title: String! # Must be present
		author: Author! # Must be present
		courses: [Course!]!
	}
	`

	//fmt.Println(input[5:14], len(input[5:14]))

	lex := lexer.New(input)

	for {
		tok := lex.NextToken()
		if tok.Type == token.EOF {
			break
		}
		if tok.Literal == input[tok.Start:tok.End] {
			fmt.Println(tok.Line, tok.Literal, tok.Start, tok.End)
		} else {
			fmt.Println(tok.Type, tok.Literal, len(tok.Literal), ">>", tok.Start, tok.End, "=", input[tok.Start:tok.End])
		}
	}

}

func manualParseObjectDefinition() {

	input := `
	"""test desc"""
	type Person {
		id: ID!
		name: String!
		age: Int
		subject: [String!]!
	}`
	lex := lexer.New(input)
	p := parser.New(lex)
	doc := p.ParseDocument()
	for i, def := range doc.Definitions {
		fmt.Println("*", i, def.GetKind(), def)
	}
}

func main() {

	// input := `
	// """test"""
	// type Person {
	// 	id: ID!
	// 	adult: Boolean!
	// 	name: String!
	// 	age: Int!
	// 	salary: Float!
	// 	length(unit: LengthUnit = METER): Float
	// 	appearsIn: [Episode]!
	// }
	// `

	input := `
	query GetBooksAndAuthors {

	 books(id: 4) {
	  title
	 }
	  
	 authors {
	  name
	 }

	}
	`

	lex := lexer.New(input)
	p := parser.New(lex)
	doc := p.ParseDocument()

	def := doc.Definitions[0]
	fmt.Println(def.String())
	fmt.Println("---->", len(doc.Definitions))

	for i, def := range doc.Definitions {

		fmt.Println("*", i, def.GetKind())

	}

	// for {
	// 	tok := lex.NextToken()
	// 	if tok.Type == token.EOF {
	// 		//fmt.Println("eof")
	// 		break
	// 	}
	// 	if tok.Literal == input[tok.Start:tok.End] {
	// 		fmt.Println(tok.Line, tok.Literal, tok.Type, tok.Start, tok.End)
	// 	} else {
	// 		fmt.Println("ERR", tok.Type, tok.Literal, len(tok.Literal), ">>", tok.Start, tok.End, "=", input[tok.Start:tok.End])
	// 	}
	// }

}
