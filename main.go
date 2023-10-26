package main

import (
	"fmt"
	"gada/reader"
	tokens "gada/token"
)

func main() {
	lexer := reader.FileLexer("examples/test.ada")
	foundTokens, lexicon := lexer.Read()
	line := -1
	for _, token := range foundTokens {
		if token.Line != line {
			line = token.Line
			if token.Position != 0 {
				fmt.Printf("\n(%s:%s %s %d %d )", token.Type, tokens.Tokens[token.Value], lexicon[token.Position-1], token.Line, token.Column)
			} else {
				fmt.Printf("\n(%s:%s %d %d )", token.Type, tokens.Tokens[token.Value], token.Line, token.Column)
			}
		} else {
			if token.Position != 0 {
				fmt.Printf("(%s:%s %s %d %d )", token.Type, tokens.Tokens[token.Value], lexicon[token.Position-1], token.Line, token.Column)
			} else {
				fmt.Printf("(%s:%s %d %d )", token.Type, tokens.Tokens[token.Value], token.Line, token.Column)
			}
		}
	}
	for _, lex := range lexicon {
		fmt.Println(lex.(string))
	}
}
