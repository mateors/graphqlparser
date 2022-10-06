package main

import (
	"fmt"

	"github.com/mateors/graphqlparser/lexer"
	"github.com/mateors/graphqlparser/token"
)

func main() {

	input := `
	type Person {
		id: ID!
		adult: Boolean!
		name: String!
		age: Int!
		salary: Float!
		length(unit: LengthUnit = METER): Float
		appearsIn: [Episode]!
	}
	
	... on Droid {
		primaryFunction
	  }`

	lex := lexer.New(input)

	for {

		tok := lex.NextToken()

		if tok.Type == token.EOF {
			break
		}
		fmt.Println(tok.Type, tok.Literal)

	}

}
