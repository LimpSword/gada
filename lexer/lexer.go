package lexer

import (
	"bufio"
	"gada/token"
	"strings"
	"unicode"
)

type Lexer struct {
	line   int
	column int
	reader *bufio.Reader
}

type Token struct {
	Type  string
	Value string
}

func NewLexer(text string) *Lexer {
	reader := bufio.NewReader(strings.NewReader(text))
	return &Lexer{reader: reader}
}

func (l *Lexer) Read() []Token {
	tokens := make([]Token, 0)
	for {
		l.column++

		r, _, err := l.reader.ReadRune()
		if err == nil {
			switch r {
			case '\n':
				l.line++
				l.column = 0
			case ' ':
				continue
			case '+':
				tokens = append(tokens, Token{Type: "Operator", Value: "+"})
			case '-':
				tokens = append(tokens, Token{Type: "Operator", Value: "-"})
			case '*':
				tokens = append(tokens, Token{Type: "Operator", Value: "*"})
			case '/':
				r, _, err := l.reader.ReadRune()
				if err == nil {
					// Check if it's a comment.
					if r == '/' {
						// Skip until the end of the line.
						for {
							r, _, err := l.reader.ReadRune()
							if err == nil {
								if r == '\n' {
									break
								}
							} else {
								break
							}
						}
					} else {
						// It can be a /= operator.
						r, _, err := l.reader.ReadRune()
						if err == nil {
							if r == '=' {
								tokens = append(tokens, Token{Type: "Operator", Value: "/="})
							} else {
								tokens = append(tokens, Token{Type: "Operator", Value: "/"})
								l.reader.UnreadRune()
							}
						}
					}
				}
			case '=':
				tokens = append(tokens, Token{Type: "Operator", Value: "="})
			case ';':
				tokens = append(tokens, Token{Type: "Separator", Value: ";"})
			case ',':
				tokens = append(tokens, Token{Type: "Separator", Value: ","})
			case ':':
				tokens = append(tokens, Token{Type: "Separator", Value: ":"})
			case '(':
				tokens = append(tokens, Token{Type: "Separator", Value: "("})
			case ')':
				tokens = append(tokens, Token{Type: "Separator", Value: ")"})
			case '>':
				if r, _, err := l.reader.ReadRune(); err == nil {
					if r == '=' {
						tokens = append(tokens, Token{Type: "Operator", Value: ">="})
					} else {
						tokens = append(tokens, Token{Type: "Operator", Value: ">"})
						l.reader.UnreadRune()
					}
				}
			case '<':
				if r, _, err := l.reader.ReadRune(); err == nil {
					if r == '=' {
						tokens = append(tokens, Token{Type: "Operator", Value: "<="})
					} else {
						tokens = append(tokens, Token{Type: "Operator", Value: "<"})
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
							tokens = append(tokens, Token{Type: "Literal", Value: char})
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
						break
					}
				}
				tokens = append(tokens, Token{Type: "Literal", Value: str})
			default:
				if unicode.IsDigit(r) {
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
					tokens = append(tokens, Token{Type: "Literal", Value: number})
				}
				if unicode.IsLetter(r) {
					name := string(r)
					for {
						r, _, err := l.reader.ReadRune()
						if err == nil {
							if token.CanBeIdentifier(r) {
								name += string(r)
								continue
							} else {
								l.reader.UnreadRune()
								break
							}
						} else {
							break
						}
					}
					if token.IsKeywordString(name) {
						tokens = append(tokens, Token{Type: "Keyword", Value: name})
					} else {
						tokens = append(tokens, Token{Type: "Identifier", Value: name})
					}
				}
			}
		} else {
			break
		}
	}
	return tokens
}
