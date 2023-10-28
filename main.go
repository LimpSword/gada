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
		if token.Beginning.Line != line {
			line = token.Beginning.Line
			if token.Position != 0 {
				fmt.Printf("\nLine : %d (%s:%s:%s from: %d to :%d )", token.Beginning.Line, token.Type, tokens.Tokens[token.Value], lexicon[token.Position-1], token.Beginning.Column, token.End.Column)
			} else {
				fmt.Printf("\nLine : %d (%s:%s from: %d to :%d )", token.Beginning.Line, token.Type, tokens.Tokens[token.Value], token.Beginning.Column, token.End.Column)
			}
		} else {
			if token.Position != 0 {
				fmt.Printf("(%s:%s:%s from: %d to :%d )", token.Type, tokens.Tokens[token.Value], lexicon[token.Position-1], token.Beginning.Column, token.End.Column)
			} else {
				fmt.Printf("(%s:%s from: %d to :%d )", token.Type, tokens.Tokens[token.Value], token.Beginning.Column, token.End.Column)
			}
		}
	}
	for _, lex := range lexicon {
		fmt.Println(lex.(string))
	}
}
