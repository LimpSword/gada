package main

import (
	"gada/reader"
)

func main() {
	lexer := reader.FileLexer("examples/test.ada")
	tokens, lexicon := lexer.Read()
	for _, token := range tokens {
		println(token.Type, token.Value, token.Position)
	}
	for _, lex := range lexicon {
		println(lex.(string))
	}
}
