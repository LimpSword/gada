package parser

import (
	"encoding/json"
	"gada/lexer"
	"gada/token"
)

type Parser struct {
	lexer *lexer.Lexer
	index int
}

type Node struct {
	Type  token.Token
	Index int
	Child []Node
}

func (p *Parser) readToken() token.Token {
	p.index++
	return token.Token(p.lexer.Tokens[p.index-1].Value)
}

func (p *Parser) readFullToken() (token.Token, int) {
	p.index++
	return token.Token(p.lexer.Tokens[p.index-1].Value), p.lexer.Tokens[p.index-1].Position
}

func (p *Parser) peekToken() token.Token {
	if p.index >= len(p.lexer.Tokens) {
		return token.EOF
	}
	return token.Token(p.lexer.Tokens[p.index].Value)
}

func Parse(lexer *lexer.Lexer) {
	parser := Parser{lexer: lexer, index: 0}
	readFile(&parser)
}

func readFile(parser *Parser) {
	node := readExpression(parser)
	// print node json
	b, err := json.MarshalIndent(node, "", "  ")
	if err != nil {
		panic(err)
	}
	println(string(b))
}

func readExpression(parser *Parser) Node {
	return readOr(parser)
}

func readOr(parser *Parser) Node {
	andExpr := readAnd(parser)
	for parser.peekToken() == token.OR {
		parser.readToken()
		other := readAnd(parser)
		andExpr = Node{Type: token.OR, Child: []Node{andExpr, other}}
	}
	return andExpr
}

func readAnd(parser *Parser) Node {
	equalityExpr := readEquality(parser)
	for parser.peekToken() == token.AND {
		parser.readToken()
		other := readEquality(parser)
		equalityExpr = Node{Type: token.AND, Child: []Node{equalityExpr, other}}
	}
	return equalityExpr
}

func readEquality(parser *Parser) Node {
	relationalExpr := readRelational(parser)
	for parser.peekToken() == token.EQL || parser.peekToken() == token.NEQ {
		tkn := parser.readToken()
		other := readRelational(parser)
		relationalExpr = Node{Type: tkn, Child: []Node{relationalExpr, other}}
	}
	return relationalExpr
}

func readRelational(parser *Parser) Node {
	additiveExpr := readAdditive(parser)
	for parser.peekToken() == token.LSS || parser.peekToken() == token.LEQ ||
		parser.peekToken() == token.GTR || parser.peekToken() == token.GEQ {
		tkn := parser.readToken()
		other := readAdditive(parser)
		additiveExpr = Node{Type: tkn, Child: []Node{additiveExpr, other}}
	}
	return additiveExpr
}

func readAdditive(parser *Parser) Node {
	multiplicativeExpr := readMultiplicative(parser)
	for parser.peekToken() == token.ADD || parser.peekToken() == token.SUB {
		tkn := parser.readToken()
		other := readMultiplicative(parser)
		multiplicativeExpr = Node{Type: tkn, Child: []Node{multiplicativeExpr, other}}
	}
	return multiplicativeExpr
}

func readMultiplicative(parser *Parser) Node {
	unaryExpr := readUnary(parser)
	for parser.peekToken() == token.MUL || parser.peekToken() == token.QUO || parser.peekToken() == token.REM_OP || parser.peekToken() == token.REM {
		tkn := parser.readToken()
		other := readUnary(parser)
		unaryExpr = Node{Type: tkn, Child: []Node{unaryExpr, other}}
	}
	return unaryExpr
}

func readUnary(parser *Parser) Node {
	if parser.peekToken() == token.SUB {
		tkn := parser.readToken()
		other := readUnary(parser)
		return Node{Type: tkn, Child: []Node{other}}
	}
	return readPrimary(parser)
}

func readPrimary(parser *Parser) Node {
	// temp return what is read
	tkn, pos := parser.readFullToken()
	return Node{Type: tkn, Index: pos}
}
