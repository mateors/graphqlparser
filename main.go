package main

import (
	"fmt"

	"github.com/mateors/graphqlparser/lexer"
	"github.com/mateors/graphqlparser/token"
)

// returns lower-case ch iff ch is ASCII letter
func lower(ch byte) byte {
	return ('a' - 'A') | ch
}

func main() {

	// fmt.Println(lower('S'))
	// fmt.Println(32 | 83)
	// os.Exit(1)

	input := `type Person {
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
	  }

	  union SearchResult = Human | Droid | Starship

	 search(text: "an") {
		__typename
		... on Human {
		  name
		  height
		}
	   }

	 }

	 type Query {
		shop(owner: String!, name: String!, location: Location): Shop!
	 }

	 0.123
	 123e4
	 2.71828e-1000
	 123E4
	 123e-4
	 123e+4
	 100.500
	 213
	 07801234567.
	 bar1
	 """Hello block""" 
	 #comment
	`

	//fmt.Println(input[5:14], len(input[5:14]))

	lex := lexer.New(input)

	for {

		tok := lex.NextToken()

		if tok.Type == token.EOF {
			break
		}
		if tok.Literal == input[tok.Start:tok.End] {
			//fmt.Println("ok")
		} else {
			fmt.Println(tok.Type, tok.Literal, len(tok.Literal), ">>", tok.Start, tok.End, "=", input[tok.Start:tok.End])
		}

	}

}
