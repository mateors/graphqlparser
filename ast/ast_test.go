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
