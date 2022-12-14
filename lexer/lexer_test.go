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

		{token.SPREAD, "..."},
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
		{token.SPREAD, "..."},
		{token.ON, "on"},
		{token.IDENT, "Human"},
		{token.LBRACE, "{"},
		{token.IDENT, "name"},
		{token.IDENT, "height"},
		{token.RBRACE, "}"},
		{token.RBRACE, "}"},
		{token.RBRACE, "}"},

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

		{token.FLOAT, "0.123"},
		{token.FLOAT, "123e4"},
		{token.FLOAT, "2.71828e-1000"},
		{token.FLOAT, "123E4"},
		{token.FLOAT, "123e-4"},
		{token.FLOAT, "123e+4"},
		{token.FLOAT, "100.500"},
		{token.INT, "213"},
		{token.FLOAT, "07801234567."},
		{token.IDENT, "bar1"},

		{token.BLOCK_STRING, "Hello block"},
		{token.HASH, "#"},
		{token.IDENT, "comment"},

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
