package token

const (
	ILLEGAL TokenType = iota
	ASSIGN            // =
	EOF

	//DELIMETERS
	BANG       // !
	COMMA      // ,
	COLON      // :
	DOT        // .
	SEMICOLON  // ;
	LPAREN     // (
	RPAREN     // )
	LBRACE     // {
	RBRACE     // }
	LBRACKET   // [
	RBRACKET   // ]
	UNDERSCORE // _

	IDENT //VARIABLES,IDENTIFIER
	INT   //IDENTIFIER

	TYPE         //KEYWORDS
	RETURN       //KEYWORDS
	QUERY        //KEYWORDS
	MUTATION     //KEYWORDS
	ON           //
	INPUT        //KEYWORDS
	SUBSCRIPTION //KEYWORDS
	ALIAS        //KEYWORDS
	FRAGMENT     //KEYWORDS
	OPERATION    //KEYWORDS
	SCALAR       //KEYWORDS
	ENUM         //KEYWORDS
	INTERFACE    //KEYWORDS
	UNION        //KEYWORDS
	DIRECTIVE    //@

)

type TokenType int

type Token struct {
	Type    TokenType
	Literal string
}

var keywords = map[string]TokenType{
	"type":         TYPE,
	"mutation":     MUTATION,
	"fragment":     FRAGMENT,
	"input":        INPUT,
	"query":        QUERY,
	"subscription": SUBSCRIPTION,
	"return":       RETURN,
	"scalar":       SCALAR,
	"enum":         ENUM,
	"interface":    INTERFACE,
	"union":        UNION,
	"on":           ON,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
