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
	l.column++
	l.line++
	position := 1
	for {
		startLine, startColumn := l.line, l.column
		r, _, err := l.reader.ReadRune()
		l.column++
		if err == nil {
			switch r {
			case '\n':
				l.line++
				l.column = 1
			case '+':
				tokens = append(tokens, Token{Type: "Operator", Value: token.ADD, Line: startLine, Column: startColumn})
			case '-':
				// comments are --
				r, _, err := l.reader.ReadRune()
				l.column++
				if err == nil {
					// Check if it's a comment.
					if r == '-' {
						// Skip until the end of the line.
						for {
							r, _, err := l.reader.ReadRune()
							l.column++
							if err == nil {
								if r == '\n' {
									l.reader.UnreadRune()
									l.column--
									break
								}
							} else {
								break
							}
						}
					} else {
						tokens = append(tokens, Token{Type: "Operator", Value: token.SUB, Line: startLine, Column: startColumn})
						l.reader.UnreadRune()
						l.column--
					}
				}
			case '*':
				tokens = append(tokens, Token{Type: "Operator", Value: token.MUL, Line: startLine, Column: startColumn})
			case '/':
				r, _, err := l.reader.ReadRune()
				l.column++
				if err == nil {
					if r == '=' {
						tokens = append(tokens, Token{Type: "Operator", Value: token.NEQ, Line: startLine, Column: startColumn})
					} else {
						tokens = append(tokens, Token{Type: "Operator", Value: token.QUO, Line: startLine, Column: startColumn})
						l.reader.UnreadRune()
						l.column--
					}
				}
			case '=':
				tokens = append(tokens, Token{Type: "Operator", Value: token.EQL, Line: startLine, Column: startColumn})
			case '.':
				tokens = append(tokens, Token{Type: "Operator", Value: token.PERIOD, Line: startLine, Column: startColumn})
			case ';':
				tokens = append(tokens, Token{Type: "Separator", Value: token.SEMICOLON, Line: startLine, Column: startColumn})
			case ',':
				tokens = append(tokens, Token{Type: "Separator", Value: token.COMMA, Line: startLine, Column: startColumn})
			case ':':
				tokens = append(tokens, Token{Type: "Separator", Value: token.COLON, Line: startLine, Column: startColumn})
			case '(':
				tokens = append(tokens, Token{Type: "Separator", Value: token.LPAREN, Line: startLine, Column: startColumn})
			case ')':
				tokens = append(tokens, Token{Type: "Separator", Value: token.RPAREN, Line: startLine, Column: startColumn})
			case '>':
				r, _, err := l.reader.ReadRune()
				l.column++
				if err == nil {
					if r == '=' {
						tokens = append(tokens, Token{Type: "Operator", Value: token.GEQ, Line: startLine, Column: startColumn})
					} else {
						tokens = append(tokens, Token{Type: "Operator", Value: token.GTR, Line: startLine, Column: startColumn})
						l.reader.UnreadRune()
						l.column--
					}
				}
			case '<':
				r, _, err := l.reader.ReadRune()
				l.column++
				if err == nil {
					if r == '=' {
						tokens = append(tokens, Token{Type: "Operator", Value: token.LEQ, Line: startLine, Column: startColumn})
					} else {
						tokens = append(tokens, Token{Type: "Operator", Value: token.LSS, Line: startLine, Column: startColumn})
						l.reader.UnreadRune()
						l.column--
					}
				}
			case '\'':
				// A char is a single character surrounded by single quotes.
				r, _, err := l.reader.ReadRune()
				l.column++
				if err == nil {
					char := string(r)
					r, _, err := l.reader.ReadRune()
					l.column++
					if err == nil {
						if r == '\'' {
							tokens = append(tokens, Token{Type: "Literal", Position: position, Value: token.CHAR, Line: startLine, Column: startColumn})
							lexi = append(lexi, char)
							position++
						} else {
							// Send an error
							unexpected := string(r)
							for {
								r, _, err := l.reader.ReadRune()
								l.column++
								//l.column
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
					l.column++
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
				tokens = append(tokens, Token{Type: "Literal", Position: position, Value: token.STRING, Line: startLine, Column: startColumn})
				lexi = append(lexi, str)
				position++
			default:
				if unicode.IsSpace(r) {
					continue
				} else if unicode.IsDigit(r) {
					number := string(r)
					for {
						r, _, err := l.reader.ReadRune()
						l.column++
						if err == nil {
							if unicode.IsDigit(r) {
								number += string(r)
								continue
							} else {
								l.reader.UnreadRune()
								l.column--
								break
							}
						} else {
							break
						}
					}
					tokens = append(tokens, Token{Type: "Literal", Position: position, Value: token.INT, Line: startLine, Column: startColumn})
					lexi = append(lexi, number)
					position++
				} else if unicode.IsLetter(r) {
					name := string(r)
					for {
						r, _, err := l.reader.ReadRune()
						l.column++
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
								l.column--
								break
							}
						} else {
							break
						}
					}
					if token.IsKeywordString(name) {
						// Check if it is them rem operator
						if name == "rem" {
							tokens = append(tokens, Token{Type: "Operator", Value: token.REM, Line: startLine, Column: startColumn})
						} else {
							tokens = append(tokens, Token{Type: "Keyword", Value: int(token.LookupIdent(name)), Line: startLine, Column: startColumn})
						}
					} else {
						tokens = append(tokens, Token{Type: "Identifier", Position: position, Value: token.IDENT, Line: startLine, Column: startColumn})
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
