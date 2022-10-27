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
	DOLLAR     // $
	AT         // @
	HASH       //#
	AMP        // &
	SEMICOLON  // ;
	LPAREN     // (
	RPAREN     // )
	LBRACE     // {
	RBRACE     // }
	LBRACKET   // [
	RBRACKET   // ]
	UNDERSCORE // _
	PIPE       // |

	IDENT        //VARIABLES,IDENTIFIER
	INT          //IDENTIFIER
	FLOAT        //Float
	STRING       //String
	BLOCK_STRING //string
	BOOLEAN      //Boolean
	ID           //ID
	NULL         //null
	VARIADIC     //...

	TYPE         //KEYWORDS
	RETURN       //KEYWORDS
	QUERY        //KEYWORDS
	MUTATION     //KEYWORDS
	ON           //
	INPUT        //KEYWORDS
	SUBSCRIPTION //KEYWORDS
	SCHEMA       //KEYWORDS
	ALIAS        //KEYWORDS
	FRAGMENT     //KEYWORDS
	OPERATION    //KEYWORDS
	SCALAR       //KEYWORDS
	ENUM         //KEYWORDS
	EXTEND       //KEYWORDS
	INTERFACE    //KEYWORDS
	IMPLEMENTS   //KEYWORDS
	UNION        //KEYWORDS
	DIRECTIVE    //KEYWORDS @

)

type TokenType int

type Token struct {
	Type    TokenType
	Line    int //row
	Start   int //column
	End     int //column
	Literal string
}

var keywords = map[string]TokenType{
	"type":         TYPE,         //
	"mutation":     MUTATION,     //
	"fragment":     FRAGMENT,     //
	"input":        INPUT,        //
	"query":        QUERY,        //
	"subscription": SUBSCRIPTION, //
	"schema":       SCHEMA,       //
	"scalar":       SCALAR,       //
	"enum":         ENUM,         //
	"extend":       EXTEND,       //
	"directive":    DIRECTIVE,    //
	"interface":    INTERFACE,    //
	"union":        UNION,        //
	"implements":   IMPLEMENTS,
	//"return":       RETURN,
	//"on":           ON,
	//"Int":          INT,
	//"Float":        FLOAT,
	//"String":       STRING,//??
	//"Boolean":      BOOLEAN,
	//"ID":           ID,
	//"null": NULL,
	//"...":  VARIADIC, //SPREAD
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}

func IsKeyword(ident string) bool {
	if _, ok := keywords[ident]; ok {
		return true
	}
	return false
}
