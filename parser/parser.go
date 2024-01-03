package parser

import (
	"encoding/json"
	"fmt"
	"gada/lexer"
	"gada/token"
	"github.com/charmbracelet/log"
	"os"
	"strconv"
)

type Parser struct {
	lexer *lexer.Lexer
	index int
}

type Node struct {
	Type     string
	Index    int
	Children []*Node
}

var logger *log.Logger

func init() {
	logger = log.New(os.Stderr)
}

func (n *Node) addChild(child Node) {
	if child.Type != "" || child.Children != nil {
		n.Children = append(n.Children, &child)
	}
}

func (n *Node) addTerminalChild(child string) {
	n.Children = append(n.Children, &Node{Type: child})
}

func (n *Node) addTerminalChilds(childs []string) {
	for _, child := range childs {
		n.Children = append(n.Children, &Node{Type: child})
	}
}

func (n *Node) addTerminalChildFromTok(p *Parser, tok token.Token) {
	n.Children = append(n.Children, &Node{Type: token.Token(p.lexer.Tokens[p.index].Value).String()})
}

func (n *Node) toJson() string {
	b, err := json.MarshalIndent(n, "", "  ")
	if err != nil {
		fmt.Println(err)
	}
	return string(b)
}

func (p *Parser) unreadToken() {
	if p.index <= 0 {
		panic("Cannot unread token")
	}
	p.index--
}

func (p *Parser) unreadTokens(nb int) {
	if p.index-nb < 0 {
		panic("Cannot unread token")
	}
	p.index -= nb
}

func (p *Parser) readToken() token.Token {
	if p.index >= len(p.lexer.Tokens) {
		return token.EOF
	}
	p.index++
	return token.Token(p.lexer.Tokens[p.index-1].Value)
}

func (p *Parser) readFullToken() (token.Token, int, string) {
	if p.index >= len(p.lexer.Tokens) {
		return token.EOF, -1, ""
	}
	p.index++
	return token.Token(p.lexer.Tokens[p.index-1].Value), p.lexer.Tokens[p.index-1].Position, p.lexer.Lexi[p.lexer.Tokens[p.index-1].Position-1]
}

func (p *Parser) peekToken() token.Token {
	if p.index >= len(p.lexer.Tokens) {
		return token.EOF
	}
	return token.Token(p.lexer.Tokens[p.index].Value)
}

func (p *Parser) printTokensBefore(i int) {
	for j := i - 1; j >= 0; j-- {
		t := token.Token(p.lexer.Tokens[p.index-j].Value)
		if t == token.IDENT {
			fmt.Print(p.lexer.Lexi[p.lexer.Tokens[p.index-j].Position-1], " ")
			continue
		}
		fmt.Print(token.Token(p.lexer.Tokens[p.index-j].Value), " ")
	}
	fmt.Println()
}

func Parse(lexer *lexer.Lexer, printAst bool) {
	parser := Parser{lexer: lexer, index: 0}
	node := readFichier(&parser)
	graph := toAst(node)
	logger.Info("Compilation successful")
	if printAst {
		logger.Info("Compilation output", "ast", graph.toJson())
	}
}

func expectToken(parser *Parser, tkn token.Token) {
	if parser.peekToken() != tkn {
		// Wait for the peek since EOF is added by this method
		red := "\x1b[0;31m"
		reset := "\x1b[0m"
		line := parser.lexer.Tokens[parser.index].Beginning.Line
		column := parser.lexer.Tokens[parser.index].Beginning.Column
		file := parser.lexer.FileName + ":" + strconv.Itoa(line) + ":" + strconv.Itoa(column)
		if parser.peekToken() == token.IDENT {
			logger.Fatal("Unexpected token", "expected", tkn, "got", parser.lexer.Lexi[parser.lexer.Tokens[parser.index].Position-1])
		}
		logger.Fatal(file+" "+"Unexpected token: "+parser.lexer.GetLineUpToToken(parser.lexer.Tokens[parser.index])+red+parser.lexer.GetToken(parser.lexer.Tokens[parser.index])+reset, "expected", tkn, "got", parser.peekToken())
	}
	parser.readToken()
}

func peekExpectToken(parser *Parser, tkn token.Token) {
	if parser.peekToken() != tkn {
		logger.Fatal("Unexpected token", "expected", tkn, "got", parser.peekToken())
	}
}

func expectTokenIdent(parser *Parser, ident string) string {
	if parser.peekToken() != token.IDENT {
		logger.Fatal("Unexpected token", "expected", token.IDENT, "got", parser.peekToken())
	}
	_, index, _ := parser.readFullToken()
	if parser.lexer.Lexi[index-1] != ident {
		logger.Fatal("Unexpected token", "expected", ident, "got", parser.lexer.Lexi[index-1])
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

// start parser
func readFichier(parser *Parser) Node {
	node := Node{}

	expectTokens(parser, []any{token.WITH, "Ada", token.PERIOD, "Text_IO", token.SEMICOLON, token.USE, "Ada", token.PERIOD, "Text_IO", token.SEMICOLON, token.PROCEDURE})
	node.addTerminalChilds([]string{"with", "Ada", ".", "Text_IO", ";", "use", "Ada", ".", "Text_IO", ";", "procedure"})
	node.addChild(readIdent(parser))
	expectTokens(parser, []any{token.IS})
	node.addTerminalChild("is")
	node.addChild(readDeclStar(parser))
	expectTokens(parser, []any{token.BEGIN})
	node.addTerminalChild("begin")
	node.addChild(readInstr_plus(parser))
	expectTokens(parser, []any{token.END})
	node.addTerminalChild("end")
	node.addChild(readIdent_opt(parser))
	expectTokens(parser, []any{token.SEMICOLON, token.EOF})
	node.addTerminalChilds([]string{";", "EOF"})
	return node
}

func readDecl(parser *Parser) Node {
	node := Node{}
	switch parser.peekToken() {
	case token.PROCEDURE:
		parser.readToken()
		node.addTerminalChild("procedure")
		node.addChild(readIdent(parser))
		node.addChild(readParams_opt(parser))
		expectTokens(parser, []any{token.IS})
		node.addTerminalChild("is")
		node.addChild(readDeclStar(parser))
		expectTokens(parser, []any{token.BEGIN})
		node.addTerminalChild("begin")
		node.addChild(readInstr_plus(parser))
		expectTokens(parser, []any{token.END})
		node.addTerminalChild("end")
		node.addChild(readIdent_opt(parser))
		expectTokens(parser, []any{token.SEMICOLON})
		node.addTerminalChild(";")
	case token.TYPE:
		parser.readToken()
		node.addTerminalChild("type")
		node.addChild(readIdent(parser))
		node.addChild(readDecl2(parser))
	case token.FUNCTION:
		parser.readToken()
		node.addTerminalChild("function")
		node.addChild(readIdent(parser))
		node.addChild(readParams_opt(parser))
		expectTokens(parser, []any{token.RETURN})
		node.addTerminalChild("return")
		node.addChild(readType_r(parser))
		expectTokens(parser, []any{token.IS})
		node.addTerminalChild("is")
		node.addChild(readDeclStar(parser))
		expectTokens(parser, []any{token.BEGIN})
		node.addTerminalChild("begin")
		node.addChild(readInstr_plus(parser))
		expectTokens(parser, []any{token.END})
		node.addTerminalChild("end")
		node.addChild(readIdent_opt(parser))
		expectTokens(parser, []any{token.SEMICOLON})
		node.addTerminalChild(";")
	case token.IDENT:
		node.addChild(readIdent_plus_comma(parser))
		expectTokens(parser, []any{token.COLON})
		node.addTerminalChild(":")
		node.addChild(readType_r(parser))
		node.addChild(readInit(parser))
		expectTokens(parser, []any{token.SEMICOLON})
		node.addTerminalChild(";")
	default:
		logger.Fatal("Unexpected token", "possible", "procedure type function ident", "got", parser.peekToken())
	}
	return node
}

func readDecl2(parser *Parser) Node {
	var node Node
	switch parser.readToken() {
	case token.IS:
		node.addTerminalChild("is")
		node.addChild(readDecl3(parser))
	case token.SEMICOLON:
		node.addTerminalChild(";")
	default:
		logger.Fatal("Unexpected token", "possible", "is ;", "got", parser.peekToken())
	}
	return node
}

func readDecl3(parser *Parser) Node {
	var node Node
	switch parser.readToken() {
	case token.ACCESS:
		node.addTerminalChild("access")
		node.addChild(readIdent(parser))
		expectTokens(parser, []any{token.SEMICOLON})
		node.addTerminalChild(";")
	case token.RECORD:
		node.addTerminalChild("record")
		node.addChild(readChampsPlus(parser))
		expectTokens(parser, []any{token.END, token.RECORD, token.SEMICOLON})
		node.addTerminalChilds([]string{"end", "record", ";"})
	default:
		logger.Fatal("Unexpected token", "possible", "access record", "got", parser.peekToken())
	}
	return node
}

func readInit(parser *Parser) Node {
	var node Node
	switch parser.peekToken() {
	case token.SEMICOLON:
	case token.COLON:
		expectTokens(parser, []any{token.COLON, token.EQL})
		node.addTerminalChilds([]string{":", "="})
		node.addChild(readExpr(parser))
	default:
		logger.Fatal("Unexpected token", "possible", "; :", "got", parser.peekToken())
	}
	return node
}

func readDeclStar(parser *Parser) Node {
	var node Node
	switch parser.peekToken() {
	case token.PROCEDURE, token.IDENT, token.TYPE, token.FUNCTION:
		node.addChild(readDecl(parser))
		node.addChild(readDeclStar(parser))
	case token.BEGIN:
	default:
		logger.Fatal("Unexpected token", "possible", "procedure ident type function begin", "got", parser.peekToken())
	}
	return node
}

func readChamps(parser *Parser) Node {
	peekExpectToken(parser, token.IDENT)
	var node Node
	node.addChild(readIdent_plus_comma(parser))

	expectTokens(parser, []any{token.COLON})
	node.addTerminalChild(":")
	node.addChild(readType_r(parser))
	expectTokens(parser, []any{token.SEMICOLON})
	node.addTerminalChild(";")
	return node
}

func readChampsPlus(parser *Parser) Node {
	peekExpectToken(parser, token.IDENT)
	var node Node
	node.addChild(readChamps(parser))
	node.addChild(readChampsPlus2(parser))
	return node
}

func readChampsPlus2(parser *Parser) Node {
	var node Node
	var tkn = parser.readToken()
	switch tkn {
	case token.IDENT:
		node.addChild(readChamps(parser))
		node.addChild(readChampsPlus2(parser))
	case token.END:
	default:
		logger.Fatal("Unexpected token", "possible", "ident end", "got", tkn)
	}
	return node
}

func readType_r(parser *Parser) Node {
	var node Node
	switch parser.peekToken() {
	case token.IDENT:
		node.addChild(readIdent(parser))
	case token.ACCESS:
		parser.readToken()
		node.addTerminalChild("access")
		node.addChild(readIdent(parser))
	default:
		logger.Fatal("Unexpected token", "possible", "ident access", "got", parser.peekToken())
	}
	return node
}

func readParams(parser *Parser) Node {
	var node Node
	expectToken(parser, token.LPAREN)
	node.addTerminalChild("(")
	node.addChild(readParamPlusSemicolon(parser))
	expectTokens(parser, []any{token.RPAREN})
	node.addTerminalChild(")")
	return node
}

func readParams_opt(parser *Parser) Node {
	var node Node
	switch parser.peekToken() {
	case token.IS, token.RETURN:
	case token.LPAREN:
		node.addChild(readParams(parser))
	default:
		logger.Fatal("Unexpected token", "possible", "is return (", "got", parser.peekToken())
	}
	return node
}

func readParam(parser *Parser) Node {
	var node Node
	switch parser.peekToken() {
	case token.IDENT:
		node.addChild(readIdent_plus_comma(parser))
		expectTokens(parser, []any{token.COLON})
		node.addTerminalChild(":")
		node.addChild(readModeOpt(parser))
		node.addChild(readType_r(parser))
	default:
		panic(fmt.Sprintf("Expected IDENT, got %s", parser.peekToken()))
	}
	return node
}

func readParamPlusSemicolon(parser *Parser) Node {
	var node Node
	switch parser.peekToken() {
	case token.IDENT:
		node.addChild(readParam(parser))
		node.addChild(readParamPlusSemicolon2(parser))
	default:
		panic(fmt.Sprintf("Expected IDENT, got %s", parser.peekToken()))
	}
	return node
}

func readParamPlusSemicolon2(parser *Parser) Node {
	var node Node
	switch parser.peekToken() {
	case token.SEMICOLON:
		parser.readToken()
		node.addTerminalChild(";")
		node.addChild(readParam(parser))
		node.addChild(readParamPlusSemicolon2(parser))
	case token.RPAREN:
	default:
		logger.Fatal("Unexpected token", "possible", "; )", "got", parser.peekToken())
	}
	return node
}

func readMode(parser *Parser) Node {
	expectToken(parser, token.IN)
	var node Node
	node.addTerminalChild("in")
	node.addChild(readMode2(parser))
	return node
}

func readMode2(parser *Parser) Node {
	var node Node
	switch parser.readToken() {
	case token.IDENT:
		node.addTerminalChild("ident")
	case token.ACCESS:
		node.addTerminalChild("access")
	case token.OUT:
	default:
		logger.Fatal("Unexpected token", "possible", "ident access out", "got", parser.peekToken())
	}
	return node
}

func readModeOpt(parser *Parser) Node {
	var node Node
	switch parser.peekToken() {
	case token.IDENT, token.ACCESS:
	case token.IN:
		node.addChild(readMode(parser))
	default:
		logger.Fatal("Unexpected token", "possible", "ident access in", "got", parser.peekToken())
	}
	return node
}

func readExpr(parser *Parser) Node {
	var node Node
	switch parser.peekToken() {
	case token.IDENT, token.LPAREN, token.NOT, token.SUB, token.INT, token.CHAR, token.TRUE, token.FALSE, token.NULL, token.NEW, token.CHAR_TOK:
		node.addChild(readOr_expr(parser))
	default:
		logger.Fatal("Unexpected token", "possible", "ident ( not - int char true false null new char", "got", parser.peekToken())
	}
	return node
}

func readOr_expr(parser *Parser) Node {
	var node Node
	switch parser.peekToken() {
	case token.IDENT, token.LPAREN, token.NOT, token.SUB, token.INT, token.CHAR, token.TRUE, token.FALSE, token.NULL, token.NEW, token.CHAR_TOK:
		node.addChild(readAnd_expr(parser))
		node.addChild(readOr_expr_tail(parser))
	default:
		logger.Fatal("Unexpected token", "possible", "ident ( not - int char true false null new char", "got", parser.peekToken())
	}
	return node
}

func readOr_expr_tail(parser *Parser) Node {
	var node Node
	switch parser.peekToken() {
	case token.OR:
		parser.readToken()
		node.addTerminalChild("or")
		node.addChild(readOr_expr_tail2(parser))
	case token.SEMICOLON, token.RPAREN, token.THEN, token.COMMA, token.LOOP:
	case token.PERIOD:
		parser.readToken()
		expectTokens(parser, []any{token.PERIOD})
		parser.readToken()
		node.addTerminalChild("..")
	default:
		logger.Fatal("Unexpected token", "possible", "or ; ) then , loop .", "got", parser.peekToken())
	}
	return node
}

func readOr_expr_tail2(parser *Parser) Node {
	var node Node
	switch parser.peekToken() {
	case token.ELSE:
		parser.readToken()
		node.addTerminalChild("else")
		node.addChild(readAnd_expr(parser))
		node.addChild(readOr_expr_tail(parser))
	case token.IDENT, token.LPAREN, token.NOT, token.SUB, token.INT, token.CHAR, token.TRUE, token.FALSE, token.NULL, token.NEW, token.CHAR_TOK:
		node.addChild(readAnd_expr(parser))
		node.addChild(readOr_expr_tail(parser))
	default:
		logger.Fatal("Unexpected token", "possible", "else ident ( not - int char true false null new char", "got", parser.peekToken())
	}
	return node
}

func readAnd_expr(parser *Parser) Node {
	var node Node
	switch parser.peekToken() {
	case token.IDENT, token.LPAREN, token.NOT, token.SUB, token.INT, token.CHAR, token.TRUE, token.FALSE, token.NULL, token.NEW, token.CHAR_TOK:
		node.addChild(readEquality_expr(parser))
		node.addChild(readAnd_expr_tail(parser))
	default:
		logger.Fatal("Unexpected token", "possible", "ident ( not - int char true false null new char", "got", parser.peekToken())
	}
	return node
}

func readAnd_expr_tail(parser *Parser) Node {
	var node Node
	switch parser.peekToken() {
	case token.AND:
		parser.readToken()
		node.addTerminalChild("and")
		node.addChild(readAnd_expr_tail2(parser))
	case token.SEMICOLON, token.RPAREN, token.OR, token.THEN, token.COMMA, token.LOOP:
	case token.PERIOD:
		parser.readToken()
		expectTokens(parser, []any{token.PERIOD})
		parser.readToken()
		node.addTerminalChild("..")
	default:
		logger.Fatal("Unexpected token", "possible", "and ; ) or then , loop .", "got", parser.peekToken())
	}
	return node
}

func readAnd_expr_tail2(parser *Parser) Node {
	var node Node
	switch parser.peekToken() {
	case token.THEN:
		parser.readToken()
		node.addTerminalChild("then")
		node.addChild(readEquality_expr(parser))
		node.addChild(readAnd_expr_tail(parser))
	case token.IDENT, token.LPAREN, token.NOT, token.SUB, token.INT, token.CHAR, token.TRUE, token.FALSE, token.NULL, token.NEW, token.CHAR_TOK:
		node.addChild(readNot_expr(parser))
		node.addChild(readAnd_expr_tail(parser))
	default:
		logger.Fatal("Unexpected token", "possible", "then ident ( not - int char true false null new char", "got", parser.peekToken())
	}
	return node
}

func readNot_expr(parser *Parser) Node {
	var node Node
	switch parser.peekToken() {
	case token.IDENT, token.LPAREN, token.NOT, token.SUB, token.INT, token.CHAR, token.TRUE, token.FALSE, token.NULL, token.NEW, token.CHAR_TOK:
		node.addChild(readEquality_expr(parser))
		node.addChild(readNot_expr_tail(parser))
	default:
		logger.Fatal("Unexpected token", "possible", "ident ( not - int char true false null new char", "got", parser.peekToken())
	}
	return node
}

func readNot_expr_tail(parser *Parser) Node {
	var node Node
	switch parser.peekToken() {
	case token.NOT:
		parser.readToken()
		node.addTerminalChild("not")
		node.addChild(readEquality_expr(parser))
		node.addChild(readNot_expr_tail(parser))
	case token.SEMICOLON, token.RPAREN, token.OR, token.AND, token.THEN, token.COMMA, token.LOOP:
	case token.PERIOD:
		parser.readToken()
		expectTokens(parser, []any{token.PERIOD})
		parser.readToken()
		node.addTerminalChild("..")
	default:
		logger.Fatal("Unexpected token", "possible", "not ; ) or and then , loop .", "got", parser.peekToken())
	}
	return node
}

func readEquality_expr(parser *Parser) Node {
	var node Node
	switch parser.peekToken() {
	case token.IDENT, token.LPAREN, token.NOT, token.SUB, token.INT, token.CHAR, token.TRUE, token.FALSE, token.NULL, token.NEW, token.CHAR_TOK:
		node.addChild(readRelational_expr(parser))
		node.addChild(readEquality_expr_tail(parser))
	default:
		logger.Fatal("Unexpected token", "possible", "ident ( not - int char true false null new char", "got", parser.peekToken())
	}
	return node
}

func readEquality_expr_tail(parser *Parser) Node {
	var node Node
	switch parser.peekToken() {
	case token.EQL:
		parser.readToken()
		node.addTerminalChild("=")
		node.addChild(readRelational_expr(parser))
		node.addChild(readEquality_expr_tail(parser))
	case token.NEQ:
		parser.readToken()
		node.addTerminalChild("/=")
		node.addChild(readRelational_expr(parser))
		node.addChild(readEquality_expr_tail(parser))
	case token.SEMICOLON, token.RPAREN, token.OR, token.AND, token.THEN, token.NOT, token.COMMA, token.LOOP:
		return node
	case token.PERIOD:
		parser.readToken()
		expectTokens(parser, []any{token.PERIOD})
		parser.readToken()
		node.addTerminalChild("..")
	default:
		logger.Fatal("Unexpected token", "possible", "= /= ; ) or and then not , loop .", "got", parser.peekToken())
	}
	return node
}

func readRelational_expr(parser *Parser) Node {
	var node Node
	switch parser.peekToken() {
	case token.IDENT, token.LPAREN, token.NOT, token.SUB, token.INT, token.CHAR, token.TRUE, token.FALSE, token.NULL, token.NEW, token.CHAR_TOK:
		node.addChild(readAdditive_expr(parser))
		node.addChild(readRelational_expr_tail(parser))
	default:
		logger.Fatal("Unexpected token", "possible", "ident ( not - int char true false null new char", "got", parser.peekToken())
	}
	return node
}

func readRelational_expr_tail(parser *Parser) Node {
	var node Node
	switch parser.peekToken() {
	case token.LSS:
		parser.readToken()
		node.addTerminalChild("<")
		node.addChild(readAdditive_expr(parser))
		node.addChild(readRelational_expr_tail(parser))
	case token.LEQ:
		parser.readToken()
		node.addTerminalChild("<=")
		node.addChild(readAdditive_expr(parser))
		node.addChild(readRelational_expr_tail(parser))
	case token.GTR:
		parser.readToken()
		node.addTerminalChild(">")
		node.addChild(readAdditive_expr(parser))
		node.addChild(readRelational_expr_tail(parser))
	case token.GEQ:
		parser.readToken()
		node.addTerminalChild(">=")
		node.addChild(readAdditive_expr(parser))
		node.addChild(readRelational_expr_tail(parser))
	case token.SEMICOLON, token.RPAREN, token.OR, token.AND, token.THEN, token.NOT, token.EQL, token.NEQ, token.COMMA, token.LOOP:
	case token.PERIOD:
		parser.readToken()
		expectTokens(parser, []any{token.PERIOD})
		parser.readToken()
		node.addTerminalChild("..")
	default:
		logger.Fatal("Unexpected token", "possible", "< <= > >= ; ) or and then not = /= , loop .", "got", parser.peekToken())
	}
	return node
}

func readAdditive_expr(parser *Parser) Node {
	var node Node
	switch parser.peekToken() {
	case token.IDENT, token.LPAREN, token.NOT, token.SUB, token.INT, token.CHAR, token.TRUE, token.FALSE, token.NULL, token.NEW, token.CHAR_TOK:
		node.addChild(readMultiplicative_expr(parser))
		node.addChild(readAdditive_expr_tail(parser))
	default:
		logger.Fatal("Unexpected token", "possible", "ident ( not - int char true false null new char", "got", parser.peekToken())
	}
	return node
}

func readAdditive_expr_tail(parser *Parser) Node {
	var node Node
	switch parser.peekToken() {
	case token.ADD:
		parser.readToken()
		node.addTerminalChild("+")
		node.addChild(readMultiplicative_expr(parser))
		node.addChild(readAdditive_expr_tail(parser))
	case token.SUB:
		parser.readToken()
		node.addTerminalChild("-")
		node.addChild(readMultiplicative_expr(parser))
		node.addChild(readAdditive_expr_tail(parser))
	case token.SEMICOLON, token.RPAREN, token.OR, token.AND, token.THEN, token.NOT, token.EQL, token.NEQ, token.LSS, token.LEQ, token.GTR, token.GEQ, token.COMMA, token.LOOP:
	case token.PERIOD:
		parser.readToken()
		expectTokens(parser, []any{token.PERIOD})
		parser.readToken()
		node.addTerminalChild("..")
	default:
		logger.Fatal("Unexpected token", "possible", "+ - ; ) or and then not = /= < <= > >= , loop .", "got", parser.peekToken())
	}
	return node
}

func readMultiplicative_expr(parser *Parser) Node {
	var node Node
	switch parser.peekToken() {
	case token.IDENT, token.LPAREN, token.NOT, token.SUB, token.INT, token.CHAR, token.TRUE, token.FALSE, token.NULL, token.NEW, token.CHAR_TOK:
		node.addChild(readUnary_expr(parser))
		node.addChild(readMultiplicative_expr_tail(parser))
	default:
		logger.Fatal("Unexpected token", "possible", "ident ( not - int char true false null new char", "got", parser.peekToken())
	}
	return node
}

func readMultiplicative_expr_tail(parser *Parser) Node {
	var node Node
	switch parser.peekToken() {
	case token.MUL:
		parser.readToken()
		node.addTerminalChild("*")
		node.addChild(readUnary_expr(parser))
		node.addChild(readMultiplicative_expr_tail(parser))
	case token.QUO:
		parser.readToken()
		node.addTerminalChild("/")
		node.addChild(readUnary_expr(parser))
		node.addChild(readMultiplicative_expr_tail(parser))
	case token.REM:
		parser.readToken()
		node.addTerminalChild("rem")
		node.addChild(readUnary_expr(parser))
		node.addChild(readMultiplicative_expr_tail(parser))
	case token.SEMICOLON, token.RPAREN, token.OR, token.AND, token.THEN, token.NOT, token.EQL, token.NEQ, token.LSS, token.LEQ, token.GTR, token.GEQ, token.ADD, token.SUB, token.COMMA, token.LOOP:
	case token.PERIOD:
		parser.readToken()
		expectTokens(parser, []any{token.PERIOD})
		parser.readToken()
		node.addTerminalChild("..")
	default:
		logger.Fatal("Unexpected token", "possible", "* / rem ; ) or and then not = /= < <= > >= + - , loop .", "got", parser.peekToken())
	}
	return node
}

func readUnary_expr(parser *Parser) Node {
	var node Node
	switch parser.peekToken() {
	case token.SUB:
		parser.readToken()
		node.addTerminalChild("-")
		node.addChild(readUnary_expr(parser))
	case token.IDENT, token.LPAREN, token.NOT, token.INT, token.CHAR, token.TRUE, token.FALSE, token.NULL, token.NEW, token.CHAR_TOK:
		node.addChild(readPrimary_expr(parser))
	default:
		logger.Fatal("Unexpected token", "possible", "- ident ( not int char true false null new char", "got", parser.peekToken())
	}
	return node
}

func readPrimary_expr(parser *Parser) Node {
	var node Node
	switch parser.peekToken() {
	case token.INT:
		node.addChild(readInt(parser))
	case token.CHAR:
		node.addChild(readChar(parser))
	case token.TRUE:
		parser.readToken()
		node.addTerminalChild("true")
	case token.FALSE:
		parser.readToken()
		node.addTerminalChild("false")
	case token.NULL:
		parser.readToken()
		node.addTerminalChild("null")
	case token.LPAREN:
		parser.readToken()
		node.addTerminalChild("(")
		node.addChild(readExpr(parser))
		expectTokens(parser, []any{token.RPAREN})
		node.addTerminalChild(")")
	case token.NOT:
		parser.readToken()
		node.addTerminalChild("not")
	case token.NEW:
		parser.readToken()
		node.addTerminalChild("new")
		node.addChild(readIdent(parser))
	case token.IDENT:
		node.addChild(readIdent(parser))
		node.addChild(readPrimary_expr2(parser))
	case token.CHAR_TOK:
		parser.readToken()
		node.addTerminalChild("character")
		expectTokens(parser, []any{token.CAST, token.VAL, token.LPAREN})
		node.addTerminalChilds([]string{"'", "val", "("})
		node.addChild(readExpr(parser))
		expectTokens(parser, []any{token.RPAREN})
		node.addTerminalChild(")")
	default:
		logger.Fatal("Unexpected token", "possible", "int char true false null ( not new ident char", "got", parser.peekToken())
	}
	return node
}

func readPrimary_expr2(parser *Parser) Node {
	var node Node
	switch parser.peekToken() {
	case token.LPAREN:
		parser.readToken()
		node.addTerminalChild("(")
		node.addChild(readExpr_plus_comma(parser))
		expectTokens(parser, []any{token.RPAREN})
		node.addTerminalChild(")")
		node.addChild(readPrimary_expr3(parser))
	case token.SEMICOLON, token.RPAREN, token.OR, token.AND, token.THEN, token.NOT, token.EQL, token.NEQ, token.LSS, token.LEQ, token.GTR, token.GEQ, token.ADD, token.SUB, token.MUL, token.QUO, token.REM, token.COMMA, token.LOOP:
		node.addChild(readAccess2(parser))
	case token.PERIOD:
		parser.readToken()
		if parser.peekToken() == token.PERIOD {
			parser.readToken()
			node.addTerminalChild("..")
			node.addChild(readAccess2(parser))
		} else {
			node.addTerminalChild(".")
			node.addChild(readAccess2(parser))
		}
	default:
		logger.Fatal("Unexpected token", "possible", "( ; ) or and then not = /= < <= > >= + - * / rem , loop .", "got", parser.peekToken())
	}
	return node
}

func readPrimary_expr3(parser *Parser) Node {
	var node Node
	switch parser.peekToken() {
	case token.PERIOD:
		parser.readToken()
		if parser.peekToken() == token.PERIOD {
			parser.readToken()
			node.addTerminalChild("..")
		} else {
			node.addTerminalChild(".")
			node.addChild(readIdent(parser))
			node.addChild(readAccess2(parser))
		}
	case token.SEMICOLON, token.RPAREN, token.OR, token.AND, token.THEN, token.NOT, token.EQL, token.NEQ, token.LSS, token.LEQ, token.GTR, token.GEQ, token.ADD, token.SUB, token.MUL, token.QUO, token.REM, token.COMMA, token.LOOP:
	default:
		logger.Fatal("Unexpected token", "possible", ". ; ) or and then not = /= < <= > >= + - * / rem , loop .", "got", parser.peekToken())
	}
	return node
}

func readAccess2(parser *Parser) Node {
	var node Node
	switch parser.peekToken() {
	case token.PERIOD:
		parser.readToken()
		if parser.peekToken() == token.PERIOD {
			parser.unreadToken()
		} else {
			node.addTerminalChild(".")
			node.addChild(readIdent(parser))
			node.addChild(readAccess2(parser))
		}
	case token.SEMICOLON, token.RPAREN, token.OR, token.AND, token.THEN, token.NOT, token.EQL, token.NEQ, token.LSS, token.LEQ, token.GTR, token.GEQ, token.ADD, token.SUB, token.MUL, token.QUO, token.REM, token.COMMA, token.LOOP:
	default:
		logger.Fatal("Unexpected token", "possible", ". ; ) or and then not = /= < <= > >= + - * / rem , loop .", "got", parser.peekToken())
	}
	return node
}

func readExpr_plus_comma(parser *Parser) Node {
	var node Node
	switch parser.peekToken() {
	case token.IDENT, token.LPAREN, token.NOT, token.SUB, token.INT, token.CHAR, token.TRUE, token.FALSE, token.NULL, token.NEW, token.CHAR_TOK:
		node.addChild(readExpr(parser))
		node.addChild(readExpr_plus_comma2(parser))
	default:
		logger.Fatal("Unexpected token", "possible", "ident ( not - int char true false null new char", "got", parser.peekToken())
	}
	return node
}

func readExpr_plus_comma2(parser *Parser) Node {
	var node Node
	switch parser.peekToken() {
	case token.COMMA:
		parser.readToken()
		node.addTerminalChild(",")
		node.addChild(readExpr(parser))
		node.addChild(readExpr_plus_comma2(parser))
	case token.RPAREN:
	default:
		logger.Fatal("Unexpected token", "possible", ", )", "got", parser.peekToken())
	}
	return node
}

func readExpr_opt(parser *Parser) Node {
	var node Node
	switch parser.peekToken() {
	case token.IDENT, token.LPAREN, token.NOT, token.SUB, token.INT, token.CHAR, token.TRUE, token.FALSE, token.NULL, token.NEW, token.CHAR_TOK:
		node.addChild(readExpr(parser))
	case token.SEMICOLON:
	default:
		logger.Fatal("Unexpected token", "possible", "ident ( not - int char true false null new char ;", "got", parser.peekToken())
	}
	return node
}

func readInstr(parser *Parser) Node {
	var node Node
	switch parser.peekToken() {
	case token.ACCESS:
		parser.readToken()
		node.addTerminalChild("access")
		expectTokens(parser, []any{token.COLON, token.EQL})
		node.addTerminalChilds([]string{":", "="})
		node.addChild(readExpr(parser))
		expectTokens(parser, []any{token.SEMICOLON})
		node.addTerminalChild(";")
	case token.IDENT:
		node.addChild(readIdent(parser))
		node.addChild(readInstr2(parser))
	case token.RETURN:
		parser.readToken()
		node.addTerminalChild("return")
		node.addChild(readExpr_opt(parser))
		expectTokens(parser, []any{token.SEMICOLON})
		node.addTerminalChild(";")
	case token.BEGIN:
		parser.readToken()
		node.addTerminalChild("begin")
		node.addChild(readInstr_plus(parser))
		expectTokens(parser, []any{token.END, token.SEMICOLON})
		node.addTerminalChilds([]string{"end", ";"})
	case token.IF:
		parser.readToken()
		node.addTerminalChild("if")
		node.addChild(readExpr(parser))
		expectTokens(parser, []any{token.THEN})
		node.addTerminalChild("then")
		node.addChild(readInstr_plus(parser))
		node.addChild(readElse_if_star(parser))
		node.addChild(readElse_instr_opt(parser))
		expectTokens(parser, []any{token.END, token.IF, token.SEMICOLON})
		node.addTerminalChilds([]string{"end", "if", ";"})
	case token.FOR:
		parser.readToken()
		node.addTerminalChild("for")
		node.addChild(readIdent(parser))
		expectTokens(parser, []any{token.IN})
		node.addTerminalChild("in")
		node.addChild(readReverse_instr(parser))
		node.addChild(readExpr(parser))
		expectTokens(parser, []any{token.PERIOD, token.PERIOD})
		node.addTerminalChild("..")
		node.addChild(readExpr(parser))
		expectTokens(parser, []any{token.LOOP})
		node.addTerminalChild("loop")
		node.addChild(readInstr_plus(parser))
		expectTokens(parser, []any{token.END, token.LOOP, token.SEMICOLON})
		node.addTerminalChilds([]string{"end", "loop", ";"})
	case token.WHILE:
		parser.readToken()
		node.addTerminalChild("while")
		node.addChild(readExpr(parser))
		expectTokens(parser, []any{token.LOOP})
		node.addTerminalChild("loop")
		node.addChild(readInstr_plus(parser))
		expectTokens(parser, []any{token.END, token.LOOP, token.SEMICOLON})
		node.addTerminalChilds([]string{"end", "loop", ";"})
	default:
		logger.Fatal("Unexpected token", "possible", "access ident return begin if for while", "got", parser.peekToken())
	}
	return node
}

func readInstr2(parser *Parser) Node {
	var node Node
	switch parser.peekToken() {
	case token.IDENT, token.BEGIN, token.END, token.RETURN, token.ACCESS, token.COLON, token.ELSE, token.PERIOD, token.IF, token.FOR, token.WHILE, token.ELSIF:
		if parser.peekToken() != token.COLON {
			expectTokens(parser, []any{token.COLON})
		}
		node.addChild(readInstr3(parser))
		expectTokens(parser, []any{token.COLON, token.EQL})
		node.addTerminalChilds([]string{":", "="})
		node.addChild(readExpr(parser))
		expectTokens(parser, []any{token.SEMICOLON})
		node.addTerminalChild(";")
	case token.SEMICOLON:
		parser.readToken()
		node.addTerminalChild(";")
	case token.LPAREN:
		parser.readToken()
		node.addTerminalChild("(")
		node.addChild(readExpr_plus_comma(parser))
		expectTokens(parser, []any{token.RPAREN})
		node.addTerminalChild(")")
		node.addChild(readInstr4(parser))
		expectTokens(parser, []any{token.SEMICOLON})
		node.addTerminalChild(";")
	default:
		logger.Fatal("Unexpected token", "possible", "; (", "got", parser.peekToken())
	}
	return node
}

func readInstr3(parser *Parser) Node {
	var node Node
	switch parser.peekToken() {
	case token.COLON:
		expectTokens(parser, []any{token.COLON, token.EQL})
		parser.unreadTokens(2)
	case token.PERIOD:
		parser.readToken()
		node.addTerminalChild(".")
		node.addChild(readIdent(parser))
		node.addChild(readInstr3(parser))
	default:
		logger.Fatal("Unexpected token", "possible", ": .", "got", parser.peekToken())
	}
	return node
}

func readInstr4(parser *Parser) Node {
	var node Node
	switch parser.peekToken() {
	case token.SEMICOLON:
	case token.COLON:
		expectTokens(parser, []any{token.COLON, token.EQL})
		node.addTerminalChilds([]string{":", "="})
		node.addChild(readExpr(parser))
	default:
		logger.Fatal("Unexpected token", "possible", "; :", "got", parser.peekToken())
	}
	return node
}

func readInstr_plus(parser *Parser) Node {
	var node Node
	switch parser.peekToken() {
	case token.BEGIN, token.RETURN, token.ACCESS, token.IF, token.FOR, token.WHILE, token.IDENT:
		node.addChild(readInstr(parser))
		node.addChild(readInstr_plus2(parser))
	default:
		logger.Fatal("Unexpected token", "possible", "begin return access if for while ident", "got", parser.peekToken())
	}
	return node
}

func readInstr_plus2(parser *Parser) Node {
	var node Node
	switch parser.peekToken() {
	case token.BEGIN, token.RETURN, token.ACCESS, token.IF, token.FOR, token.ELSIF /*, token.IF, token.FOR*/, token.WHILE, token.IDENT:
		node.addChild(readInstr(parser))
		node.addChild(readInstr_plus2(parser))
	case token.END, token.ELSE:
	default:
		logger.Fatal("Unexpected token", "possible", "begin return access if for while ident", "got", parser.peekToken())
	}
	return node
}

func readElse_if(parser *Parser) Node {
	var node Node
	switch parser.peekToken() {
	case token.ELSIF:
		parser.readToken()
		node.addTerminalChild("elsif")
		node.addChild(readExpr(parser))
		expectTokens(parser, []any{token.THEN})
		node.addTerminalChild("then")
		node.addChild(readInstr_plus(parser))
	default:
		logger.Fatal("Unexpected token", "possible", "elsif", "got", parser.peekToken())
	}
	return node
}

func readElse_if_star(parser *Parser) Node {
	var node Node
	switch parser.peekToken() {
	case token.ELSIF:
		node.addChild(readElse_if(parser))
		node.addChild(readElse_if_star(parser))
	case token.ELSE, token.END, token.BEGIN, token.RETURN, token.ACCESS, token.IF, token.FOR, token.WHILE, token.IDENT:
		node = Node{Type: "ElseIfStar"}
	default:
		logger.Fatal("Unexpected token", "possible", "begin return access if for while end else ident", "got", parser.peekToken())
	}
	return node
}

func readElse_instr(parser *Parser) Node {
	var node Node
	switch parser.peekToken() {
	case token.ELSE:
		parser.readToken()
		node.addTerminalChild("else")
		node.addChild(readInstr_plus(parser))
	default:
		logger.Fatal("Unexpected token", "possible", "else", "got", parser.peekToken())
	}
	return node
}

func readElse_instr_opt(parser *Parser) Node {
	var node Node
	switch parser.peekToken() {
	case token.ELSE:
		node.addChild(readElse_instr(parser))
	case token.END:
	default:
		logger.Fatal("Unexpected token", "possible", "else end", "got", parser.peekToken())
	}
	return node
}

func readReverse_instr(parser *Parser) Node {
	var node Node
	switch parser.peekToken() {
	case token.IDENT, token.LPAREN, token.NOT, token.SUB, token.INT, token.CHAR, token.TRUE, token.FALSE, token.NULL, token.NEW, token.CHAR_TOK:
	case token.REVERSE:
		parser.readToken()
		node.addTerminalChild("reverse")
	default:
		logger.Fatal("Unexpected token", "possible", "ident ( not - int char true false null new char reverse", "got", parser.peekToken())
	}
	return node
}

func readIdent(parser *Parser) Node {
	peekExpectToken(parser, token.IDENT)
	_, index, value := parser.readFullToken()
	node := Node{Type: value}
	node.Index = index
	return node
}

func readInt(parser *Parser) Node {
	peekExpectToken(parser, token.INT)
	_, index, value := parser.readFullToken()
	node := Node{Type: value}
	node.Index = index
	return node
}

func readChar(parser *Parser) Node {
	peekExpectToken(parser, token.CHAR)
	_, index, value := parser.readFullToken()
	node := Node{Type: value}
	node.Index = index
	return node
}

func readIdent_opt(parser *Parser) Node {
	var node Node
	switch parser.peekToken() {
	case token.SEMICOLON:
	case token.IDENT:
		node.addChild(readIdent(parser))
	default:
		logger.Fatal("Unexpected token", "possible", "ident ;", "got", parser.peekToken())
	}
	return node
}

func readIdent_plus_comma(parser *Parser) Node {
	var node Node
	switch parser.peekToken() {
	case token.IDENT:
		node.addChild(readIdent(parser))
		node.addChild(readIdent_plus_comma2(parser))
	default:
		logger.Fatal("Unexpected token", "possible", "ident", "got", parser.peekToken())
	}
	return node
}

func readIdent_plus_comma2(parser *Parser) Node {
	var node Node
	switch parser.peekToken() {
	case token.SEMICOLON:
	case token.COMMA:
		parser.readToken()
		node.addTerminalChild(",")
		node.addChild(readIdent(parser))
		node.addChild(readIdent_plus_comma2(parser))
	case token.COLON:
	default:
		logger.Fatal("Unexpected token", "possible", "; ,", "got", parser.peekToken())
	}
	return node

}
