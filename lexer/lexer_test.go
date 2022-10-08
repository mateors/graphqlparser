package lexer

import (
	"testing"

	"github.com/mateors/graphqlparser/token"
)

func TestNextToken(t *testing.T) {

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
	  `

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{

		{token.TYPE, "type"},
		{token.IDENT, "Person"},
		{token.LBRACE, "{"},
		{token.IDENT, "id"},
		{token.COLON, ":"},
		{token.ID, "ID"},
		{token.BANG, "!"},

		{token.IDENT, "adult"},
		{token.COLON, ":"},
		{token.BOOLEAN, "Boolean"},
		{token.BANG, "!"},

		{token.IDENT, "name"},
		{token.COLON, ":"},
		{token.STRING, "String"},
		{token.BANG, "!"},

		{token.IDENT, "age"},
		{token.COLON, ":"},
		{token.INT, "Int"},
		{token.BANG, "!"},

		{token.IDENT, "salary"},
		{token.COLON, ":"},
		{token.FLOAT, "Float"},
		{token.BANG, "!"},

		{token.IDENT, "length"},
		{token.LPAREN, "("},
		{token.IDENT, "unit"},
		{token.COLON, ":"},
		{token.IDENT, "LengthUnit"},
		{token.ASSIGN, "="},
		{token.IDENT, "METER"},
		{token.RPAREN, ")"},
		{token.COLON, ":"},
		{token.FLOAT, "Float"}, //appearsIn: [Episode]!

		{token.IDENT, "appearsIn"},
		{token.COLON, ":"},
		{token.LBRACKET, "["},
		{token.IDENT, "Episode"},
		{token.RBRACKET, "]"},
		{token.BANG, "!"},
		{token.RBRACE, "}"},

		{token.VARIADIC, "..."},
		{token.ON, "on"},
		{token.IDENT, "Droid"},
		{token.LBRACE, "{"},
		{token.IDENT, "primaryFunction"},
		{token.RBRACE, "}"},

		//union SearchResult = Human | Droid | Starship
		{token.UNION, "union"},
		{token.IDENT, "SearchResult"},
		{token.ASSIGN, "="},
		{token.IDENT, "Human"},
		{token.PIPE, "|"},
		{token.IDENT, "Droid"},
		{token.PIPE, "|"},
		{token.IDENT, "Starship"},

		{token.IDENT, "search"},
		{token.LPAREN, "("},
		{token.IDENT, "text"},
		{token.COLON, ":"},
		{token.STRING, "an"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.IDENT, "__typename"},
		{token.VARIADIC, "..."},
		{token.ON, "on"},
		{token.IDENT, "Human"},
		{token.LBRACE, "{"},
		{token.IDENT, "name"},
		{token.IDENT, "height"},
		{token.RBRACE, "}"},
		{token.RBRACE, "}"},
		{token.RBRACE, "}"},

		// type Query {
		// 	shop(owner: String!, name: String!, location: Location): Shop!
		//  }

		{token.TYPE, "type"},
		{token.IDENT, "Query"},
		{token.LBRACE, "{"},
		{token.IDENT, "shop"},
		{token.LPAREN, "("},
		{token.IDENT, "owner"},
		{token.COLON, ":"},
		{token.STRING, "String"},
		{token.BANG, "!"},
		{token.COMMA, ","},
		{token.IDENT, "name"},
		{token.COLON, ":"},
		{token.STRING, "String"},
		{token.BANG, "!"},
		{token.COMMA, ","},
		{token.IDENT, "location"},
		{token.COLON, ":"},
		{token.IDENT, "Location"},
		{token.RPAREN, ")"},
		{token.COLON, ":"},
		{token.IDENT, "Shop"},
		{token.BANG, "!"},
		{token.RBRACE, "}"},

		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {

		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q", i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)
		}
	}

}
