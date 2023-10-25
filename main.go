package main

import (
	lexer2 "gada/lexer"
	"gada/reader"
	tokens "gada/token"
)

func main() {
	lexer := reader.FileLexer("examples/test.ada")
	foundTokens, lexicon := lexer.Read()
	for _, token := range foundTokens {
		println(token.Type, tokens.Tokens[token.Value], token.Position)
	}
	for _, lex := range lexicon {
		println(lex.(string))
	}
}
