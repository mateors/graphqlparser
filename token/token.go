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
	SPREAD       //...

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

var tokenDescription = map[TokenType]string{
	ILLEGAL:      "ILLEGAL",
	ASSIGN:       "=",
	EOF:          "EOF",
	BANG:         "!",     //DELIMETERS
	COMMA:        ",",     //DELIMETERS
	COLON:        ":",     //DELIMETERS
	DOT:          ".",     //DELIMETERS
	DOLLAR:       "$",     //DELIMETERS
	AT:           "@",     //DELIMETERS
	HASH:         "#",     //DELIMETERS
	AMP:          "&",     //DELIMETERS
	SEMICOLON:    ";",     //DELIMETERS
	LPAREN:       "(",     //DELIMETERS
	RPAREN:       ")",     //DELIMETERS
	LBRACE:       "{",     //DELIMETERS
	RBRACE:       "}",     //DELIMETERS
	LBRACKET:     "[",     //DELIMETERS
	RBRACKET:     "]",     //DELIMETERS
	UNDERSCORE:   "_",     //DELIMETERS
	PIPE:         "|",     //DELIMETERS
	IDENT:        "IDENT", //VARIABLES,IDENTIFIER
	INT:          "INT",
	FLOAT:        "FLOAT",
	STRING:       "STRING",
	BLOCK_STRING: "BLOCK_STRING",
	BOOLEAN:      "BOOLEAN", //Boolean
	ID:           "ID",
	NULL:         "NULL",
	SPREAD:       "...",          //VARIADIC|SPREAD
	TYPE:         "TYPE",         //KEYWORDS
	RETURN:       "RETURN",       //KEYWORDS
	QUERY:        "QUERY",        //KEYWORDS
	MUTATION:     "MUTATION",     //KEYWORDS
	ON:           "ON",           //
	INPUT:        "INPUT",        //KEYWORDS
	SUBSCRIPTION: "SUBSCRIPTION", //KEYWORDS
	SCHEMA:       "SCHEMA",       //KEYWORDS
	ALIAS:        "ALIAS",        //KEYWORDS
	FRAGMENT:     "FRAGMENT",     //KEYWORDS
	OPERATION:    "OPERATION",    //KEYWORDS
	SCALAR:       "SCALAR",       //KEYWORDS
	ENUM:         "ENUM",         //KEYWORDS
	EXTEND:       "EXTEND",       //KEYWORDS
	INTERFACE:    "INTERFACE",    //KEYWORDS
	IMPLEMENTS:   "IMPLEMENTS",   //KEYWORDS
	UNION:        "UNION",        //KEYWORDS
	DIRECTIVE:    "DIRECTIVE",    //KEYWORDS @
}

func Name(tokType TokenType) string {
	if name, isFound := tokenDescription[tokType]; isFound {
		return name
	}
	return "UNKNOWN"
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
	"...": SPREAD, //VARIADIC
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
