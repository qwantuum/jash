package token

type TokenType string

const (
	EOF     TokenType = "EOF"
	NEWLINE TokenType = "NEWLINE"
	INDENT  TokenType = "INDENT"
	DEDENT  TokenType = "DEDENT"

	IDENT  TokenType = "IDENT"
	NUMBER TokenType = "NUMBER"
	STRING TokenType = "STRING"

	DEF    TokenType = "DEF"
	RETURN TokenType = "RETURN"
	IF     TokenType = "IF"
	ELSE   TokenType = "ELSE"
	ELIF   TokenType = "ELIF"
	FOR    TokenType = "FOR"
	IN     TokenType = "IN"
	WHILE  TokenType = "WHILE"
	REPEAT TokenType = "REPEAT"
	TRUE   TokenType = "TRUE"
	FALSE  TokenType = "FALSE"
	NULL   TokenType = "NULL"
	AND    TokenType = "AND"
	OR     TokenType = "OR"
	NOT    TokenType = "NOT"
	IMPORT TokenType = "IMPORT"
	BREAK   TokenType = "BREAK"
	CONTINUE TokenType = "CONTINUE"

	LBRACE   TokenType = "LBRACE"
	RBRACE   TokenType = "RBRACE"
	LBRACKET TokenType = "LBRACKET"
	RBRACKET TokenType = "RBRACKET"
	LPAREN   TokenType = "LPAREN"
	RPAREN   TokenType = "RPAREN"
	COLON    TokenType = "COLON"
	COMMA    TokenType = "COMMA"
	DOT      TokenType = "DOT"

	ASSIGN TokenType = "ASSIGN"
	PLUS   TokenType = "PLUS"
	MINUS  TokenType = "MINUS"
	STAR   TokenType = "STAR"
	SLASH  TokenType = "SLASH"
	EQ     TokenType = "EQ"
	NEQ    TokenType = "NEQ"
	LT     TokenType = "LT"
	GT     TokenType = "GT"
	LTE    TokenType = "LTE"
	GTE    TokenType = "GTE"
	MOD    TokenType = "MOD"

	PLUS_ASSIGN  TokenType = "PLUS_ASSIGN"
	MINUS_ASSIGN TokenType = "MINUS_ASSIGN"
	STAR_ASSIGN  TokenType = "STAR_ASSIGN"
	SLASH_ASSIGN TokenType = "SLASH_ASSIGN"
)

type Token struct {
	Type    TokenType
	Literal string
	Line    int
	Column  int
}

var Keywords = map[string]TokenType{
	"def":    DEF,
	"return": RETURN,
	"if":     IF,
	"else":   ELSE,
	"elif":   ELIF,
	"for":    FOR,
	"in":     IN,
	"while":  WHILE,
	"repeat": REPEAT,
	"true":   TRUE,
	"false":  FALSE,
	"null":   NULL,
	"and":    AND,
	"or":     OR,
	"not":    NOT,
	"import":  IMPORT,
	"break":   BREAK,
	"continue": CONTINUE,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := Keywords[ident]; ok {
		return tok
	}
	return IDENT
}
