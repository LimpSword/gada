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
				tokens = append(tokens, Token{Type: "Operator", Value: token.ADD})
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
									break
								}
							} else {
								break
							}
						}
					} else {
						tokens = append(tokens, Token{Type: "Operator", Value: token.SUB})
					}
				}
			case '*':
				tokens = append(tokens, Token{Type: "Operator", Value: token.MUL})
			case '/':
				r, _, err := l.reader.ReadRune()
				if err == nil {
					if r == '=' {
						tokens = append(tokens, Token{Type: "Operator", Value: token.NEQ})
					} else {
						tokens = append(tokens, Token{Type: "Operator", Value: token.QUO})
						l.reader.UnreadRune()
					}
				}
			case '=':
				tokens = append(tokens, Token{Type: "Operator", Value: token.EQL})
			case '.':
				tokens = append(tokens, Token{Type: "Operator", Value: token.PERIOD})
			case ';':
				tokens = append(tokens, Token{Type: "Separator", Value: token.SEMICOLON})
			case ',':
				tokens = append(tokens, Token{Type: "Separator", Value: token.COMMA})
			case ':':
				tokens = append(tokens, Token{Type: "Separator", Value: token.COLON})
			case '(':
				tokens = append(tokens, Token{Type: "Separator", Value: token.LPAREN})
			case ')':
				tokens = append(tokens, Token{Type: "Separator", Value: token.RPAREN})
			case '>':
				if r, _, err := l.reader.ReadRune(); err == nil {
					if r == '=' {
						tokens = append(tokens, Token{Type: "Operator", Value: token.GEQ})
					} else {
						tokens = append(tokens, Token{Type: "Operator", Value: token.GTR})
						l.reader.UnreadRune()
					}
				}
			case '<':
				if r, _, err := l.reader.ReadRune(); err == nil {
					if r == '=' {
						tokens = append(tokens, Token{Type: "Operator", Value: token.LEQ})
					} else {
						tokens = append(tokens, Token{Type: "Operator", Value: token.LSS})
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
							tokens = append(tokens, Token{Type: "Literal", Position: position, Value: token.CHAR})
						}
					}
					// FIXME: Check if we have a lexical error.
					lexi = append(lexi, char)
					position++
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
				// FIXME: Check if we have a lexical error, ie if the next rune is not an operator or whitespace.
				tokens = append(tokens, Token{Type: "Literal", Position: position, Value: token.STRING})
				lexi = append(lexi, str)
				position++
			default:
				if unicode.IsSpace(r) {
					continue
				} else if unicode.IsDigit(r) {
					// FIXME: what about if we have a lot of rem operators?

					number := string(r)
					for {
						r, _, err := l.reader.ReadRune()
						if err == nil {
							if unicode.IsDigit(r) {
								number += string(r)
								continue
							} else {
								// TODO

								// Check if we have a lexical error.
								unexpected := string(r)
								for {
									r, _, err := l.reader.ReadRune()
									if err == nil {
										if !unicode.IsSpace(r) && !token.IsOperatorString(string(r)) {
											unexpected += string(r)
										} else {
											l.reader.UnreadRune()
											break
										}
									} else {
										break
									}
								}
								if len(unexpected) > 0 {
									// Check for the rem operator
									println(unexpected)
									if strings.HasPrefix(unexpected, "rem") && (len(unexpected) == 3 || len(unexpected) > 3 && !unicode.IsDigit(rune(unexpected[3]))) {
										println("Lexical error: unexpected character '" + unexpected + "' at line " + strconv.FormatInt(int64(l.line), 10) + " and column " + strconv.FormatInt(int64(l.column), 10) + ".")
									} else {
										// Unread the unexpected characters.
										println(len(unexpected))
										for i := len(unexpected) - 1; i >= 0; i-- {
											l.reader.UnreadRune()
										}
										r, _, _ := l.reader.ReadRune()
										println(string(r))
										l.reader.UnreadRune()
									}
								}
								l.reader.UnreadRune()
								break
							}
						} else {
							break
						}
					}
					tokens = append(tokens, Token{Type: "Literal", Position: position, Value: token.INT})
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
							tokens = append(tokens, Token{Type: "Operator", Value: token.REM})
						} else {
							tokens = append(tokens, Token{Type: "Keyword", Value: int(token.LookupIdent(name))})
						}
					} else {
						tokens = append(tokens, Token{Type: "Identifier", Position: position, Value: token.IDENT})
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
