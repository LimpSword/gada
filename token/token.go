package token

import (
	"unicode"
)

type Token int

const (
	eof = rune(0)

	// Special tokens
	ILLEGAL = iota
	EOF
	COMMENT

	// Literals
	literals_beg
	IDENT
	INT    // 12345
	CHAR   // 'a'
	STRING // "abc"
	literals_end

	// Operators
	operator_beg
	ADD    // +
	SUB    // -
	MUL    // *
	QUO    // /
	REM_OP // %

	EQL // =
	LSS // <
	GTR // >

	NEQ // !=
	LEQ // <=
	GEQ // >=

	LPAREN // (
	COMMA  // ,
	PERIOD // .

	RPAREN    // )
	SEMICOLON // ;
	COLON     // :
	operator_end

	// Keywords
	keywords_beg
	ACCESS
	AND
	BEGIN
	ELSE
	ELSIF
	END
	FALSE
	FOR
	FUNCTION
	IF
	IN
	IS
	LOOP
	NEW
	NOT
	NULL
	OR
	OUT
	PROCEDURE
	RECORD
	REM
	RETURN
	REVERSE
	THEN
	TRUE
	TYPE
	USE
	WHILE
	WITH
	keywords_end
)

var tokens = [...]string{
	ILLEGAL: "ILLEGAL",
	EOF:     "EOF",
	COMMENT: "COMMENT",

	IDENT:  "IDENT",
	INT:    "INT",
	CHAR:   "CHAR",
	STRING: "STRING",

	ADD:    "+",
	SUB:    "-",
	MUL:    "*",
	QUO:    "/",
	REM_OP: "%",

	EQL: "=",
	LSS: "<",
	GTR: ">",

	NEQ: "!=",
	LEQ: "<=",
	GEQ: ">=",

	LPAREN: "(",
	COMMA:  ",",
	PERIOD: ".",

	RPAREN:    ")",
	SEMICOLON: ";",
	COLON:     ":",

	ACCESS:    "access",
	AND:       "and",
	BEGIN:     "begin",
	ELSE:      "else",
	ELSIF:     "elsif",
	END:       "end",
	FALSE:     "false",
	FOR:       "for",
	FUNCTION:  "function",
	IF:        "if",
	IN:        "in",
	IS:        "is",
	LOOP:      "loop",
	NEW:       "new",
	NOT:       "not",
	NULL:      "null",
	OR:        "or",
	OUT:       "out",
	PROCEDURE: "procedure",
	RECORD:    "record",
	REM:       "rem",
	RETURN:    "return",
	REVERSE:   "reverse",
	THEN:      "then",
	TRUE:      "true",
	TYPE:      "type",
	USE:       "use",
	WHILE:     "while",
	WITH:      "with",
}

func (t Token) String() string {
	return tokens[t]
}

func (t Token) Precedence() int {
	switch t {
	case PERIOD:
		return 8
	case MUL, QUO, REM_OP:
		return 7
	case ADD, SUB:
		return 6
	case LSS, LEQ, GTR, GEQ:
		return 5
	case EQL, NEQ:
		return 4
	case NOT:
		return 3
	case AND:
		return 2
	case OR:
		return 1
	default:
		return 0
	}
}

var keywords map[string]Token

func init() {
	keywords = make(map[string]Token)
	for i := keywords_beg + 1; i < keywords_end; i++ {
		keywords[tokens[i]] = Token(i)
	}
}

func LookupIdent(ident string) Token {
	if tok, ok := keywords[ident]; ok {
		// The token is a keyword.
		return tok
	}
	return IDENT
}

func IsLiteral(tok Token) bool {
	return tok > literals_beg && tok < literals_end
}

func IsOperator(tok Token) bool {
	return tok > operator_beg && tok < operator_end
}

func IsKeyword(tok Token) bool {
	return tok > keywords_beg && tok < keywords_end
}

func IsKeywordString(s string) bool {
	if _, ok := keywords[s]; ok {
		return true
	}
	return false
}

func IsIdentifier(name string) bool {
	if len(name) == 0 {
		return false
	}
	if _, ok := keywords[name]; ok {
		return false
	}
	for i, c := range name {
		if i == 0 && !unicode.IsLetter(c) {
			return false
		}
		if !unicode.IsLetter(c) && !unicode.IsDigit(c) && c == '_' {
			return false
		}
	}
	return true
}

func CanBeIdentifier(c rune) bool {
	return unicode.IsLetter(c) || c == '_' || unicode.IsDigit(c)
}
