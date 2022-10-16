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

	directives := []*ast.Directive{}
	directives = append(directives, &ast.Directive{
		Kind:      ast.DIRECTIVE,
		Name:      &ast.Name{Kind: ast.NAME, Value: "excludeField"},
		Arguments: args,
	})

	iv1 := &ast.InputValueDefinition{
		Name:         &ast.Name{Kind: ast.NAME, Token: token.Token{}, Value: "unit"},
		Type:         &ast.NamedType{Kind: ast.NAMED_TYPE, Token: token.Token{}, Name: &ast.Name{Kind: "Name", Token: token.Token{}, Value: "LengthUnit"}},
		DefaultValue: &ast.StringValue{Kind: ast.STRING_VALUE, Token: token.Token{}, Value: "METER"},
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
