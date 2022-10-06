package main

import (
	"fmt"

	"github.com/mateors/graphqlparser/lexer"
	"github.com/mateors/graphqlparser/token"
)

func main() {

	input := `type Student{ }`
	lex := lexer.New(input)

	for {

		tok := lex.NextToken()

		if tok.Type == token.EOF {
			break
		}
		fmt.Println(tok, tok.Type, tok.Literal)

	}

}
