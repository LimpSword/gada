package parser

import "gada/lexer"

type Parser struct {
	lexer *lexer.Lexer
	index int
}

func Parse(lexer *lexer.Lexer) {
	parser := Parser{lexer: lexer, index: 0}
	readFile(&parser)
}

func readFile(parser *Parser) {

}
