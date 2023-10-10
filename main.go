package main

import lexer2 "gada/lexer"

func main() {
	lexer := lexer2.NewLexer("type point is record\nabcisse : integer ;\nordonnee : integer ;\nend record;")
	tokens := lexer.Read()
	for _, token := range tokens {
		println(token.Type, token.Value)
	}
}
