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
	PIPE       // |

	IDENT    //VARIABLES,IDENTIFIER
	INT      //IDENTIFIER
	FLOAT    //Float
	STRING   //String
	BOOLEAN  //Boolean
	ID       //ID
	NULL     //null
	VARIADIC //...

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
	IMPLEMENTS   //KEYWORDS
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
	"implements":   IMPLEMENTS,
	"union":        UNION,
	"on":           ON,
	"Int":          INT,
	"Float":        FLOAT,
	"String":       STRING,
	"Boolean":      BOOLEAN,
	"ID":           ID,
	"null":         NULL,
	"...":          VARIADIC,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
