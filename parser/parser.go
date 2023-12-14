package parser

import (
	"encoding/json"
	"fmt"
	"gada/lexer"
	"gada/token"
)

type Parser struct {
	lexer *lexer.Lexer
	index int
}

type Node struct {
	Type  string
	Index int
	child []Node
}

func (n Node) addChild(child Node) {
	n.child = append(n.child, child)
}

func (n Node) toJson() string {
	b, err := json.MarshalIndent(n, "", "  ")
	if err != nil {
		fmt.Println(err)
	}
	return string(b)
}

func (p *Parser) readToken() token.Token {
	if p.index >= len(p.lexer.Tokens) {
		return token.EOF
	}
	p.index++
	return token.Token(p.lexer.Tokens[p.index-1].Value)
}

func (p *Parser) readFullToken() (token.Token, int) {
	if p.index >= len(p.lexer.Tokens) {
		return token.EOF, -1
	}
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
	node := readFichier(&parser)

	fmt.Println("Compilation successful")
	fmt.Println("AST:")
	fmt.Println(node.toJson())
}

func expectToken(parser *Parser, tkn token.Token) {
	if parser.peekToken() != tkn {
		panic(fmt.Sprintf("Expected %s, got %s", tkn, parser.peekToken()))
	}
	parser.readToken()
}

func peekExpectToken(parser *Parser, tkn token.Token) {
	if parser.peekToken() != tkn {
		panic(fmt.Sprintf("Expected %s, got %s", tkn, parser.peekToken()))
	}
}

func expectTokenIdent(parser *Parser, ident string) string {
	if parser.peekToken() != token.IDENT {
		panic(fmt.Sprintf("Expected IDENT, got %s", parser.peekToken()))
	}
	_, index := parser.readFullToken()
	if parser.lexer.Lexi[index-1] != ident {
		panic(fmt.Sprintf("Expected IDENT %s, got %s", ident, parser.lexer.Lexi[index-1]))
	}
	return parser.lexer.Lexi[index-1]
}

func expectTokens(parser *Parser, tkns []any) {
	for _, tkn := range tkns {
		if t, ok := tkn.(int); ok {
			expectToken(parser, token.Token(t))
		} else {
			// expect identifier with name tkn
			expectTokenIdent(parser, tkn.(string))
		}
	}
}

func readFichier(parser *Parser) Node {
	node := Node{Type: "Fichier"}

	expectTokens(parser, []any{token.WITH, "Ada", token.PERIOD, "Text_IO", token.SEMICOLON, token.USE, "Ada", token.PERIOD, "Text_IO", token.SEMICOLON, token.PROCEDURE})

	node.addChild(readIdent(parser))
	expectTokens(parser, []any{token.IS, token.BEGIN})
	node.addChild(readInstr_plus(parser))
	expectTokens(parser, []any{token.END})
	node.addChild(readIdent_opt(parser))
	expectTokens(parser, []any{token.SEMICOLON, token.EOF})
	return node
}

func readDecl(parser *Parser) Node {
	var node Node
	switch parser.readToken() {
	case token.PROCEDURE:
		node = Node{Type: "DeclProcedure"}
		node.addChild(readIdent(parser))
		node.addChild(readParams_opt(parser))
		expectTokens(parser, []any{token.IS})
		node.addChild(readDeclStar(parser))
		expectTokens(parser, []any{token.BEGIN})
		node.addChild(readInstr_plus(parser))
		expectTokens(parser, []any{token.END})
		node.addChild(readIdent_opt(parser))
	case token.TYPE:
		node = Node{Type: "DeclType"}
		node.addChild(readIdent(parser))
		node.addChild(readDecl2(parser))
	case token.FUNCTION:
		node = Node{Type: "DeclFunction"}
		node.addChild(readIdent(parser))
		node.addChild(readParams_opt(parser))
		expectTokens(parser, []any{token.RETURN, token.TYPE, token.IS})
		node.addChild(readDeclStar(parser))
		expectTokens(parser, []any{token.BEGIN})
		node.addChild(readInstr_plus(parser))
		expectTokens(parser, []any{token.END})
		node.addChild(readIdent_opt(parser))
		expectTokens(parser, []any{token.SEMICOLON})
	case token.IDENT:
		node = Node{Type: "DeclVar"}
		node.addChild(readIdent_plus_comma(parser))
		expectTokens(parser, []any{token.COLON, token.TYPE})
		node.addChild(readInit(parser))
		expectTokens(parser, []any{token.SEMICOLON})
	default:
		panic(fmt.Sprintf("Expected PROCEDURE, TYPE, FUNCTION or IDENT, got %s", parser.peekToken()))
	}
	return node
}

func readDecl2(parser *Parser) Node {
	var node Node
	switch parser.readToken() {
	case token.IS:
		node = Node{Type: "DeclTypeIs"}
		node.addChild(readDecl3(parser))
	case token.SEMICOLON:
		node = Node{Type: "DeclTypeSemicolon"}
	default:
		panic(fmt.Sprintf("Expected IS or SEMICOLON, got %s", parser.peekToken()))
	}
	return node
}

func readDecl3(parser *Parser) Node {
	var node Node
	switch parser.readToken() {
	case token.ACCESS:
		node = Node{Type: "DeclTypeAccess"}
		node.addChild(readIdent(parser))
		expectTokens(parser, []any{token.SEMICOLON})
	case token.RECORD:
		node = Node{Type: "DeclTypeRecord"}
		node.addChild(readChampsPlus(parser))
		expectTokens(parser, []any{token.END, token.RECORD, token.SEMICOLON})
	default:
		panic(fmt.Sprintf("Expected ACCESS or RECORD, got %s", parser.peekToken()))
	}
	return node
}

func readInit(parser *Parser) Node {
	var node Node
	switch parser.readToken() {
	case token.SEMICOLON:
		node = Node{Type: "InitSemicolon"}
	case token.COLON:
		expectTokens(parser, []any{token.EQL})
		node = Node{Type: "Init"}
		node.addChild(readExpr(parser))
	default:
		panic(fmt.Sprintf("Expected SEMICOLON or COLON, got %s", parser.peekToken()))
	}
	return node
}

func readDeclStar(parser *Parser) Node {
	var node Node
	switch parser.readToken() {
	case token.PROCEDURE, token.IDENT, token.TYPE, token.FUNCTION:
		node = Node{Type: "DeclStarProcedure"}
		node.addChild(readDecl(parser))
		node.addChild(readDeclStar(parser))
	case token.BEGIN:
		node = Node{Type: "DeclStarBegin"}
	default:
		panic(fmt.Sprintf("Expected PROCEDURE, IDENT, TYPE, FUNCTION or BEGIN, got %s", parser.peekToken()))
	}
	return node
}

func readChamps(parser *Parser) Node {
	peekExpectToken(parser, token.IDENT)

	node := Node{Type: "Champs"}
	node.addChild(readIdent_plus_comma(parser))

	expectTokens(parser, []any{token.COLON, token.TYPE, token.SEMICOLON})
	return node
}

func readChampsPlus(parser *Parser) Node {
	peekExpectToken(parser, token.IDENT)

	node := Node{Type: "ChampsPlus"}
	node.addChild(readChamps(parser))
	node.addChild(readChampsPlus2(parser))
	return node
}

func readChampsPlus2(parser *Parser) Node {
	var node Node
	switch parser.readToken() {
	case token.IDENT:
		node = Node{Type: "ChampsPlus2"}
		node.addChild(readChamps(parser))
		node.addChild(readChampsPlus2(parser))
	case token.END:
		node = Node{Type: "ChampsPlus2End"}
	default:
		panic(fmt.Sprintf("Expected IDENT or END, got %s", parser.peekToken()))
	}
	return node
}

func readType_r(parser *Parser) Node {
	var node Node
	switch parser.peekToken() {
	case token.IDENT:
		node = Node{Type: "TypeRIdent"}
		node.addChild(readIdent(parser))
	case token.ACCESS:
		parser.readToken()

		node = Node{Type: "TypeRAccess"}
		node.addChild(readIdent(parser))
	default:
		panic(fmt.Sprintf("Expected IDENT or ACCESS, got %s", parser.peekToken()))
	}
	return node
}

func readParams(parser *Parser) Node {
	expectToken(parser, token.LPAREN)

	node := Node{Type: "Params"}
	node.addChild(readParamPlusSemicolon(parser))

	expectTokens(parser, []any{token.RPAREN, token.SEMICOLON})
	return node
}

func readParams_opt(parser *Parser) Node {
	var node Node
	switch parser.peekToken() {
	case token.IS, token.RETURN:
		node = Node{Type: "ParamsOpt"}
	case token.LPAREN:
		node = Node{Type: "ParamsOptParams"}
		node.addChild(readParams(parser))
	default:
		panic(fmt.Sprintf("Expected IS, RETURN or LPAREN, got %s", parser.peekToken()))
	}
	return node
}

func readParam(parser *Parser) Node {
	peekExpectToken(parser, token.IDENT)

	node := Node{Type: "Param"}
	node.addChild(readIdent_plus_comma(parser))

	expectTokens(parser, []any{token.COLON})

	node.addChild(readModeOpt(parser))
	node.addChild(readType_r(parser))
	return node
}

func readParamPlusSemicolon(parser *Parser) Node {
	peekExpectToken(parser, token.IDENT)

	node := Node{Type: "ParamPlusSemicolon"}
	node.addChild(readParam(parser))
	node.addChild(readParamPlusSemicolon2(parser))

	return node
}

func readParamPlusSemicolon2(parser *Parser) Node {
	var node Node
	switch parser.peekToken() {
	case token.SEMICOLON:
		parser.readToken()
		node = Node{Type: "ParamPlusSemicolon2"}
		node.addChild(readParam(parser))
		node.addChild(readParamPlusSemicolon2(parser))
	case token.RPAREN:
		node = Node{Type: "ParamPlusSemicolon2RParen"}
	default:
		panic(fmt.Sprintf("Expected SEMICOLON or RPAREN, got %s", parser.peekToken()))
	}
	return node
}

func readMode(parser *Parser) Node {
	expectToken(parser, token.IN)
	node := Node{Type: "ModeIn"}
	node.addChild(readMode2(parser))
	return node
}

func readMode2(parser *Parser) Node {
	var node Node
	switch parser.readToken() {
	case token.IDENT, token.ACCESS:
		node = Node{Type: "Mode2Ident"}
	case token.OUT:
		node = Node{Type: "Mode2Out"}
	default:
		panic(fmt.Sprintf("Expected IDENT, ACCESS or OUT, got %s", parser.peekToken()))
	}
	return node
}

func readModeOpt(parser *Parser) Node {
	var node Node
	switch parser.peekToken() {
	case token.IDENT, token.ACCESS:
		node = Node{Type: "ModeOpt"}
	case token.IN:
		node = Node{Type: "ModeOptMode"}
		node.addChild(readMode(parser))
	default:
		panic(fmt.Sprintf("Expected IDENT, ACCESS or IN, got %s", parser.peekToken()))
	}
	return node
}

func readExpr(parser *Parser) Node {
	var node Node
	switch parser.peekToken() {
	case token.IDENT, token.LPAREN, token.NOT, token.SUB, token.INT, token.CHAR, token.TRUE, token.FALSE, token.NULL, token.NEW, token.CHAR_TOK:
		node = Node{Type: "ExprIdent"}
		node.addChild(readOr_expr(parser))
	default:
		panic(fmt.Sprintf("Expected IDENT, LPAREN, NOT, SUB, INT, CHAR, TRUE, FALSE, NULL, NEW or CHAR_TOK, got %s", parser.peekToken()))
	}
	return node
}

func readOr_expr(parser *Parser) Node {
	return Node{}
}

func readOr_expr_tail(parser *Parser) Node {
	return Node{}
}

func readOr_expr_tail2(parser *Parser) Node {
	return Node{}
}

func readAnd_expr(parser *Parser) Node {
	return Node{}
}

func readAnd_expr_tail(parser *Parser) Node {
	return Node{}
}

func readAnd_expr_tail2(parser *Parser) Node {
	return Node{}
}

func readNot_expr(parser *Parser) Node {
	return Node{}
}

func readNot_expr_tail(parser *Parser) Node {
	return Node{}
}

func readEquality_expr(parser *Parser) Node {
	return Node{}
}

func readEquality_expr_tail(parser *Parser) Node {
	return Node{}
}

func readRelational_expr(parser *Parser) Node {
	return Node{}
}

func readRelational_expr_tail(parser *Parser) Node {
	return Node{}
}

func readAdditive_expr(parser *Parser) Node {
	return Node{}
}

func readAdditive_expr_tail(parser *Parser) Node {
	return Node{}
}

func readMultiplicative_expr(parser *Parser) Node {
	return Node{}
}

func readMultiplicative_expr_tail(parser *Parser) Node {
	return Node{}
}

func readUnary_expr(parser *Parser) Node {
	return Node{}
}

func readPrimary_expr(parser *Parser) Node {
	return Node{}
}

func readPrimary_expr2(parser *Parser) Node {
	return Node{}
}

func readPrimary_expr3(parser *Parser) Node {
	return Node{}
}

func readAccess2(parser *Parser) Node {
	return Node{}
}

func readExpr_plus_comma(parser *Parser) Node {
	return Node{}
}

func readExpr_plus_comma2(parser *Parser) Node {
	return Node{}
}

func readExpr_opt(parser *Parser) Node {
	return Node{}
}

func readInstr(parser *Parser) Node {
	node := Node{Type: "Instr"}
	switch parser.peekToken() {
	case token.BEGIN, token.RETURN, token.ACCESS, token.IF, token.FOR, token.WHILE, token.IDENT:
		node.addChild(readIdent(parser))
		node.addChild(readInstr2(parser))
	default:
		panic(fmt.Sprintf("Expected BEGIN, RETURN, ACCESS, IF, FOR, WHILE or IDENT, got %s", parser.peekToken()))
	}
	return node
}

func readInstr2(parser *Parser) Node {
	return Node{}
}

func readInstr_plus(parser *Parser) Node {
	node := Node{Type: "InstrPlus"}
	switch parser.peekToken() {
	case token.BEGIN, token.RETURN, token.ACCESS, token.IF, token.FOR, token.WHILE, token.IDENT:
		node.addChild(readInstr(parser))
		node.addChild(readInstr_plus2(parser))
	default:
		panic(fmt.Sprintf("Expected BEGIN, RETURN, ACCESS, IF, FOR, WHILE or IDENT, got %s", parser.peekToken()))
	}
	return Node{}
}

func readInstr_plus2(parser *Parser) Node {
	return Node{}
}

func readElse_if(parser *Parser) Node {
	return Node{}
}

func readElse_if_star(parser *Parser) Node {
	return Node{}
}

func readElse_instr(parser *Parser) Node {
	return Node{}
}

func readElse_instr_opt(parser *Parser) Node {
	return Node{}
}

func readReverse_instr(parser *Parser) Node {
	return Node{}
}

func readIdent(parser *Parser) Node {
	peekExpectToken(parser, token.IDENT)

	node := Node{Type: "Ident"}
	_, index := parser.readFullToken()
	node.Index = index
	return node
}

func readIdent_opt(parser *Parser) Node {
	node := Node{Type: "IdentOpt"}
	switch parser.peekToken() {
	case token.SEMICOLON:
		return node
	case token.IDENT:
		node.addChild(readIdent(parser))
	default:
		panic(fmt.Sprintf("Expected SEMICOLON or IDENT, got %s", parser.peekToken()))
	}
	return node
}

func readIdent_plus_comma(parser *Parser) Node {
	return Node{}
}

func readIdent_plus_comma2(parser *Parser) Node {
	return Node{}
}
