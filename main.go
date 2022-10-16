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

	var ttype ast.Type
	fd := ast.InputValueDefinition{}
	fd.Name = &ast.Name{Kind: "Name", Token: token.Token{}, Value: "name"}
	fd.Type = ttype // &ast.StringValue{Kind: ast.STRING_VALUE, Token: token.Token{}, Value: "test"}
	//ast.NonNullType
	fmt.Println(fd.String())

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
