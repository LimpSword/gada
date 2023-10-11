package main

import lexer2 "gada/lexer"

func main() {
	lexer := lexer2.NewLexer("56rem5rem8 + 8 * 5 --a comment\na+b")
	tokens, lexicon := lexer.Read()
	for _, token := range tokens {
		println(token.Type, token.Value, token.Position)
	}
	for _, lex := range lexicon {
		println(lex.(string))
	}
}
