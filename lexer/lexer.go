package lexer

import (
	"bufio"
	"gada/token"
	"strconv"
	"strings"
	"unicode"
)

type Lexer struct {
	line   int
	column int
	reader *bufio.Reader
}

type Token struct {
	Type     string
	Position int
	Value    int
	Line     int
	Column   int
}

func NewLexer(text string) *Lexer {
	reader := bufio.NewReader(strings.NewReader(text))
	return &Lexer{reader: reader}
}

// Read reads the text and returns the list of tokens and the associated lexicon.
func (l *Lexer) Read() ([]Token, []any) {
	tokens := make([]Token, 0)
	lexi := make([]interface{}, 0)
	position := 1
	for {
		l.column++

		r, _, err := l.reader.ReadRune()
		if err == nil {
			switch r {
			case '\n':
				l.line++
				l.column = 0
			case '+':
				tokens = append(tokens, Token{Type: "Operator", Value: token.ADD, Line: l.line, Column: l.column})
			case '-':
				// comments are --
				r, _, err := l.reader.ReadRune()
				if err == nil {
					// Check if it's a comment.
					if r == '-' {
						// Skip until the end of the line.
						for {
							r, _, err := l.reader.ReadRune()
							if err == nil {
								if r == '\n' {
									l.reader.UnreadRune()
									break
								}
							} else {
								break
							}
						}
					} else {
						tokens = append(tokens, Token{Type: "Operator", Value: token.SUB, Line: l.line, Column: l.column})
						l.reader.UnreadRune()
					}
				}
			case '*':
				tokens = append(tokens, Token{Type: "Operator", Value: token.MUL, Line: l.line, Column: l.column})
			case '/':
				r, _, err := l.reader.ReadRune()
				if err == nil {
					if r == '=' {
						tokens = append(tokens, Token{Type: "Operator", Value: token.NEQ, Line: l.line, Column: l.column})
					} else {
						tokens = append(tokens, Token{Type: "Operator", Value: token.QUO, Line: l.line, Column: l.column})
						l.reader.UnreadRune()
					}
				}
			case '=':
				tokens = append(tokens, Token{Type: "Operator", Value: token.EQL, Line: l.line, Column: l.column})
			case '.':
				tokens = append(tokens, Token{Type: "Operator", Value: token.PERIOD, Line: l.line, Column: l.column})
			case ';':
				tokens = append(tokens, Token{Type: "Separator", Value: token.SEMICOLON, Line: l.line, Column: l.column})
			case ',':
				tokens = append(tokens, Token{Type: "Separator", Value: token.COMMA, Line: l.line, Column: l.column})
			case ':':
				tokens = append(tokens, Token{Type: "Separator", Value: token.COLON, Line: l.line, Column: l.column})
			case '(':
				tokens = append(tokens, Token{Type: "Separator", Value: token.LPAREN, Line: l.line, Column: l.column})
			case ')':
				tokens = append(tokens, Token{Type: "Separator", Value: token.RPAREN, Line: l.line, Column: l.column})
			case '>':
				if r, _, err := l.reader.ReadRune(); err == nil {
					if r == '=' {
						tokens = append(tokens, Token{Type: "Operator", Value: token.GEQ, Line: l.line, Column: l.column})
					} else {
						tokens = append(tokens, Token{Type: "Operator", Value: token.GTR, Line: l.line, Column: l.column})
						l.reader.UnreadRune()
					}
				}
			case '<':
				if r, _, err := l.reader.ReadRune(); err == nil {
					if r == '=' {
						tokens = append(tokens, Token{Type: "Operator", Value: token.LEQ, Line: l.line, Column: l.column})
					} else {
						tokens = append(tokens, Token{Type: "Operator", Value: token.LSS, Line: l.line, Column: l.column})
						l.reader.UnreadRune()
					}
				}
			case '\'':
				// A char is a single character surrounded by single quotes.
				r, _, err := l.reader.ReadRune()
				if err == nil {
					char := string(r)
					r, _, err := l.reader.ReadRune()
					if err == nil {
						if r == '\'' {
							tokens = append(tokens, Token{Type: "Literal", Position: position, Value: token.CHAR, Line: l.line, Column: l.column})
							lexi = append(lexi, char)
							position++
						} else {
							// Send an error
							unexpected := string(r)
							for {
								r, _, err := l.reader.ReadRune()
								if err == nil {
									if r == '\'' {
										break
									} else {
										unexpected += string(r)
									}
								} else {
									println("Lexical error: unexpected end of file at line " + strconv.FormatInt(int64(l.line), 10) + " and column " + strconv.FormatInt(int64(l.column), 10) + ".")
									break
								}
							}
							println("Lexical error: unexpected character '" + char + unexpected + "' at line " + strconv.FormatInt(int64(l.line), 10) + " and column " + strconv.FormatInt(int64(l.column), 10) + ".")
						}
					}
				}
			case '"':
				// A string is a sequence of characters surrounded by double quotes.
				str := ""
				for {
					r, _, err := l.reader.ReadRune()
					if err == nil {
						if r == '"' {
							break
						} else {
							str += string(r)
						}
					} else {
						println("Lexical error: unexpected end of file at line " + strconv.FormatInt(int64(l.line), 10) + " and column " + strconv.FormatInt(int64(l.column), 10) + ".")
						break
					}
				}
				tokens = append(tokens, Token{Type: "Literal", Position: position, Value: token.STRING, Line: l.line, Column: l.column})
				lexi = append(lexi, str)
				position++
			default:
				if unicode.IsSpace(r) {
					continue
				} else if unicode.IsDigit(r) {
					number := string(r)
					for {
						r, _, err := l.reader.ReadRune()
						if err == nil {
							if unicode.IsDigit(r) {
								number += string(r)
								continue
							} else {
								l.reader.UnreadRune()
								break
							}
						} else {
							break
						}
					}
					tokens = append(tokens, Token{Type: "Literal", Position: position, Value: token.INT, Line: l.line, Column: l.column})
					lexi = append(lexi, number)
					position++
				} else if unicode.IsLetter(r) {
					name := string(r)
					for {
						r, _, err := l.reader.ReadRune()
						if err == nil {
							if token.CanBeIdentifier(r) {
								name += string(r)
								continue
							} else {
								// Check if we have a lexical error.
								if !unicode.IsSpace(r) && !token.IsOperatorString(string(r)) {
									println("Lexical error: unexpected character '" + string(r) + "' at line " + strconv.FormatInt(int64(l.line), 10) + " and column " + strconv.FormatInt(int64(l.column), 10) + ".")
								}
								l.reader.UnreadRune()
								break
							}
						} else {
							break
						}
					}
					if token.IsKeywordString(name) {
						// Check if it is them rem operator
						if name == "rem" {
							tokens = append(tokens, Token{Type: "Operator", Value: token.REM, Line: l.line, Column: l.column})
						} else {
							tokens = append(tokens, Token{Type: "Keyword", Value: int(token.LookupIdent(name)), Line: l.line, Column: l.column})
						}
					} else {
						tokens = append(tokens, Token{Type: "Identifier", Position: position, Value: token.IDENT, Line: l.line, Column: l.column})
						lexi = append(lexi, name)
						position++
					}
				} else {
					// Check if we have a lexical error.
					if !unicode.IsSpace(r) {
						println("Lexical error: unexpected character '" + string(r) + "' at line " + strconv.FormatInt(int64(l.line), 10) + " and column " + strconv.FormatInt(int64(l.column), 10) + ".")
					}
				}
			}
		} else {
			if err.Error() != "EOF" {
				panic(err)
			}
			break
		}
	}
	return tokens, lexi
}
