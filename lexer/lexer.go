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

	Tokens []Token
	Lexi   []string
}

type Position struct {
	Line   int
	Column int
}

type Token struct {
	Type      string
	Position  int
	Value     int
	Beginning Position
	End       Position
}

func NewLexer(text string) *Lexer {
	reader := bufio.NewReader(strings.NewReader(text))
	return &Lexer{reader: reader}
}

// Read reads the text and returns the list of Tokens and the associated lexicon.
func (l *Lexer) Read() ([]Token, []string) {
	tokens := make([]Token, 0)
	lexi := make([]string, 0)
	l.column++
	l.line++
	position := 1
	for {
		beginPos := Position{l.line, l.column}
		r, _, err := l.reader.ReadRune()
		l.column++
		if err == nil {
			switch r {
			case '\n':
				l.line++
				l.column = 1
			case '+':
				tokens = append(tokens, Token{Type: "Operator", Value: token.ADD, Beginning: beginPos, End: Position{l.line, l.column}})
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
						tokens = append(tokens, Token{Type: "Operator", Value: token.SUB, Beginning: beginPos, End: Position{l.line, l.column}})
						l.reader.UnreadRune()
						l.column--
					}
				}
			case '*':
				tokens = append(tokens, Token{Type: "Operator", Value: token.MUL, Beginning: beginPos, End: Position{l.line, l.column}})
			case '/':
				r, _, err := l.reader.ReadRune()
				l.column++
				if err == nil {
					if r == '=' {
						tokens = append(tokens, Token{Type: "Operator", Value: token.NEQ, Beginning: beginPos, End: Position{l.line, l.column}})
					} else {
						tokens = append(tokens, Token{Type: "Operator", Value: token.QUO, Beginning: beginPos, End: Position{l.line, l.column}})
						l.reader.UnreadRune()
						l.column--
					}
				}
			case '=':
				tokens = append(tokens, Token{Type: "Operator", Value: token.EQL, Beginning: beginPos, End: Position{l.line, l.column}})
			case '.':
				tokens = append(tokens, Token{Type: "Operator", Value: token.PERIOD, Beginning: beginPos, End: Position{l.line, l.column}})
			case ';':
				tokens = append(tokens, Token{Type: "Separator", Value: token.SEMICOLON, Beginning: beginPos, End: Position{l.line, l.column}})
			case ',':
				tokens = append(tokens, Token{Type: "Separator", Value: token.COMMA, Beginning: beginPos, End: Position{l.line, l.column}})
			case ':':
				tokens = append(tokens, Token{Type: "Separator", Value: token.COLON, Beginning: beginPos, End: Position{l.line, l.column}})
			case '(':
				tokens = append(tokens, Token{Type: "Separator", Value: token.LPAREN, Beginning: beginPos, End: Position{l.line, l.column}})
			case ')':
				tokens = append(tokens, Token{Type: "Separator", Value: token.RPAREN, Beginning: beginPos, End: Position{l.line, l.column}})
			case '>':
				r, _, err := l.reader.ReadRune()
				l.column++
				if err == nil {
					if r == '=' {
						tokens = append(tokens, Token{Type: "Operator", Value: token.GEQ, Beginning: beginPos, End: Position{l.line, l.column}})
					} else {
						tokens = append(tokens, Token{Type: "Operator", Value: token.GTR, Beginning: beginPos, End: Position{l.line, l.column}})
						l.reader.UnreadRune()
						l.column--
					}
				}
			case '<':
				r, _, err := l.reader.ReadRune()
				l.column++
				if err == nil {
					if r == '=' {
						tokens = append(tokens, Token{Type: "Operator", Value: token.LEQ, Beginning: beginPos, End: Position{l.line, l.column}})
					} else {
						tokens = append(tokens, Token{Type: "Operator", Value: token.LSS, Beginning: beginPos, End: Position{l.line, l.column}})
						l.reader.UnreadRune()
						l.column--
					}
				}
			case '\'':
				// Check if the keyword 'character' is present before.
				if len(tokens) > 0 {
					if tokens[len(tokens)-1].Value == token.IDENT && lexi[len(lexi)-1] == "character" {
						tokens = append(tokens, Token{Type: "Operator", Value: token.CAST, Beginning: beginPos, End: Position{l.line, l.column}})
						break
					}
				}
				// A char is a single character surrounded by single quotes.
				r, _, err := l.reader.ReadRune()
				l.column++
				if err == nil {
					char := string(r)
					r, _, err := l.reader.ReadRune()
					l.column++
					if err == nil {
						if r == '\'' {
							tokens = append(tokens, Token{Type: "Literal", Position: position, Value: token.CHAR, Beginning: beginPos, End: Position{l.line, l.column}})
							lexi = append(lexi, char)
							position++
						} else {
							// Send an error
							unexpected := string(r)
							eofBreaked := true
							for {
								r, _, err := l.reader.ReadRune()
								l.column++
								//l.column
								if err == nil {
									if r == '\'' {
										break
									} else if r == '\n' {
										println("Lexical error: new line in rune at line" + strconv.FormatInt(int64(l.line), 10) + " and column " + strconv.FormatInt(int64(l.column)-1, 10) + ".")
										tokens = append(tokens, Token{Type: "ILLEGAL", Position: position, Value: token.ILLEGAL, Beginning: beginPos, End: Position{l.line, l.column}})
										lexi = append(lexi, "Lexical error: new line in rune at line "+strconv.FormatInt(int64(l.line), 10)+" and column "+strconv.FormatInt(int64(l.column)-1, 10)+".")
										position++
										l.reader.UnreadRune()
										l.column--
										break
									} else {
										unexpected += string(r)
									}
								} else {
									println("Lexical error: unexpected end of file at line " + strconv.FormatInt(int64(l.line), 10) + " and column " + strconv.FormatInt(int64(l.column)-1, 10) + ".")
									tokens = append(tokens, Token{Type: "ILLEGAL", Position: position, Value: token.ILLEGAL, Beginning: beginPos, End: Position{l.line, l.column}})
									lexi = append(lexi, "Lexical error: unexpected end of file at line "+strconv.FormatInt(int64(l.line), 10)+" and column "+strconv.FormatInt(int64(l.column)-1, 10)+".")
									position++
									eofBreaked = true
									break
								}
							}
							if !eofBreaked {
								println("Lexical error: unexpected character '" + char + unexpected + "' at line " + strconv.FormatInt(int64(l.line), 10) + " between column " + strconv.FormatInt(int64(beginPos.Column), 10) + " and " + strconv.FormatInt(int64(l.column), 10) + ".")
								tokens = append(tokens, Token{Type: "ILLEGAL", Position: position, Value: token.ILLEGAL, Beginning: beginPos, End: Position{l.line, l.column}})
								lexi = append(lexi, "Lexical error: unexpected character '"+char+unexpected+"' at line "+strconv.FormatInt(int64(l.line), 10)+" between column "+strconv.FormatInt(int64(beginPos.Column), 10)+" and "+strconv.FormatInt(int64(l.column), 10)+".")
								position++
							}
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
						println("Lexical error: unexpected end of file at line " + strconv.FormatInt(int64(l.line), 10) + " and column " + strconv.FormatInt(int64(l.column)-1, 10) + ".")
						tokens = append(tokens, Token{Type: "ILLEGAL", Position: position, Value: token.ILLEGAL, Beginning: beginPos, End: Position{l.line, l.column}})
						lexi = append(lexi, "Lexical error: unexpected end of file at line "+strconv.FormatInt(int64(l.line), 10)+" and column "+strconv.FormatInt(int64(l.column)-1, 10)+".")
						position++
						break
					}
				}
				tokens = append(tokens, Token{Type: "Literal", Position: position, Value: token.STRING, Beginning: beginPos, End: Position{l.line, l.column}})
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
					tokens = append(tokens, Token{Type: "Literal", Position: position, Value: token.INT, Beginning: beginPos, End: Position{l.line, l.column}})
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
								//// Check if we have a lexical error.
								//if !unicode.IsSpace(r) && !token.IsOperatorString(string(r)) {
								//	println("Lexical error: unexpected character '" + string(r) + "' at line " + strconv.FormatInt(int64(l.line), 10) + " and column " + strconv.FormatInt(int64(l.column), 10) + ".")
								//}
								l.reader.UnreadRune()
								l.column--
								break
							}
						} else {
							l.column--
							break
						}
					}
					if token.IsKeywordString(name) {
						// Check if it is them rem operator
						if name == "rem" {
							tokens = append(tokens, Token{Type: "Operator", Value: token.REM, Beginning: beginPos, End: Position{l.line, l.column}})
						} else {
							tokens = append(tokens, Token{Type: "Keyword", Value: int(token.LookupIdent(name)), Beginning: beginPos, End: Position{l.line, l.column}})
						}
					} else {
						tokens = append(tokens, Token{Type: "Literals", Position: position, Value: token.IDENT, Beginning: beginPos, End: Position{l.line, l.column}})
						lexi = append(lexi, name)
						position++
					}
				} else {
					// Check if we have a lexical error.
					if !unicode.IsSpace(r) {
						println("Lexical error: unexpected character '" + string(r) + "' at line " + strconv.FormatInt(int64(l.line), 10) + " and column " + strconv.FormatInt(int64(l.column)-1, 10) + ".")
						tokens = append(tokens, Token{Type: "ILLEGAL", Position: position, Value: token.ILLEGAL, Beginning: beginPos, End: Position{l.line, l.column}})
						lexi = append(lexi, "Lexical error: unexpected character '"+string(r)+"' at line "+strconv.FormatInt(int64(l.line), 10)+" and column "+strconv.FormatInt(int64(l.column)-1, 10)+".")
						position++
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
	l.Tokens = tokens
	l.Lexi = lexi
	return tokens, lexi
}
