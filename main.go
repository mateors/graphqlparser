package main

import (
	"fmt"
	"os"

	"github.com/mateors/graphqlparser/ast"
	"github.com/mateors/graphqlparser/lexer"
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

func main() {

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
