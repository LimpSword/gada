package lexer

import (
	"gada/lexer"
	"gada/token"
)

func getExpected() map[string]testlexer {
	expected := make(map[string]testlexer)
	// helloWorld
	tokens := make([]lexer.Token, 0)
	lexi := make([]interface{}, 0)
	tokens = append(tokens, lexer.Token{"", 0, token.WITH, lexer.Position{1, 1}, lexer.Position{1, 5}})
	tokens = append(tokens, lexer.Token{"", 1, token.IDENT, lexer.Position{1, 6}, lexer.Position{1, 13}})
	tokens = append(tokens, lexer.Token{"", 0, token.SEMICOLON, lexer.Position{1, 13}, lexer.Position{1, 14}})
	lexi = append(lexi, "Text_IO")
	expected["helloWorld"] = testlexer{
		tokens:  tokens,
		lexidic: lexi}
	return expected
}
