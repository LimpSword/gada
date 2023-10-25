package main

import lexer2 "gada/lexer"
import tokens "gada/token"

func main() {
	lexer := lexer2.NewLexer("2 rem 2 + b + 'a'")
	foundTokens, lexicon := lexer.Read()
	for _, token := range foundTokens {
		println(token.Type, tokens.Tokens[token.Value], token.Position)
	}
	for _, lex := range lexicon {
		println(lex.(string))
	}
}
