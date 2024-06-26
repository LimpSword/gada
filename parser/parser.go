package parser

import (
	"encoding/json"
	"fmt"
	"gada/lexer"
	"gada/token"
	"github.com/charmbracelet/log"
	"os"
	"os/exec"
	"strconv"
)

type Parser struct {
	lexer     *lexer.Lexer
	index     int
	exprError bool
	hadError  bool
}

type Node struct {
	Type         string
	Index        int
	Children     []*Node
	Line, Column int
}

var logger *log.Logger

func init() {
	logger = log.New(os.Stderr)
}

func (n *Node) addChild(child Node) {
	n.Children = append(n.Children, &child)
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

func (p *Parser) peekTokenToString() string {
	if p.index >= len(p.lexer.Tokens) {
		return "EOF"
	}
	if p.lexer.Tokens[p.index].Value == token.IDENT {
		return p.lexer.Lexi[p.lexer.Tokens[p.index].Position-1]
	}
	if p.lexer.Tokens[p.index].Value == token.CHAR_TOK {
		return "'" + p.lexer.Lexi[p.lexer.Tokens[p.index].Position-1] + "'"
	}
	if p.lexer.Tokens[p.index].Value == token.INT {
		return p.lexer.Lexi[p.lexer.Tokens[p.index].Position-1]
	}
	return token.Token(p.lexer.Tokens[p.index].Value).String()
}

func (p *Parser) peekTokenFurther(i int) token.Token {
	if p.index+i >= len(p.lexer.Tokens) {
		return token.EOF
	}
	return token.Token(p.lexer.Tokens[p.index+i].Value)
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

func Parse(lex *lexer.Lexer, printAst bool, pythonExecutable string) {
	parser := Parser{lexer: lex, index: 0, exprError: false, hadError: false}
	lex.Tokens = append(lex.Tokens, lexer.Token{Value: token.EOF, Beginning: lexer.Position{Line: lex.Tokens[len(lex.Tokens)-1].End.Line, Column: lex.Tokens[len(lex.Tokens)-1].End.Column}, End: lexer.Position{Line: lex.Tokens[len(lex.Tokens)-1].End.Line, Column: lex.Tokens[len(lex.Tokens)-1].End.Column}})
	node := readFichier(&parser)
	os.WriteFile("./test/parser/parsetree.json", []byte(node.toJson()), 0644)
	graph := toAst(node, *lex)
	os.WriteFile("./test/parser/ast.json", []byte(graph.toJson()), 0644)
	_, err := ReadAST(&graph, true)
	if err != nil {
		logger.Error("Error while reading AST", "error", err)
		return
	}
	if printAst {
		logger.Info("Rendering AST...")
		cmd := exec.Command(pythonExecutable, "./test/parser/json_to_image.py")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Start()
		if err != nil {
			logger.Error("Error while running python script", "error", err)
			return
		}
		err = cmd.Wait()
		if err != nil {
			logger.Error("Error while running python script", "error", err)
			return
		}
		logger.Info("AST rendered")
	}
	logger.Info("Checking semantics...")
	CheckSemantics(graph)

	logger.Info("Compiling to ASM...")
	os.WriteFile("./test/parser/astSem.json", []byte(graph.toJson()), 0644)
	ReadASTToASM(graph)
	if parser.hadError {
		// no crash for now
		logger.Error("Compilation failed")
		return
	} else {
		logger.Info("Compilation successful")
	}
}

func (parser *Parser) advanceExpr(tokens []token.Token) {
	for parser.peekToken() != token.EOF {
		for _, tkn := range tokens {
			if parser.peekToken() == tkn {
				return
			} else if parser.peekToken() == token.IDENT {
				if parser.peekTokenFurther(1) == token.COLON || parser.peekTokenFurther(1) == token.EQL {
					return
				}
			}
		}
		parser.readToken()
	}
}

func (parser *Parser) advance(tokens []token.Token) {
	for parser.peekToken() != token.EOF {
		for _, tkn := range tokens {
			if parser.peekToken() == tkn {
				return
			}
		}
		parser.readToken()
	}
}

func (parser *Parser) advance2(tokens ...token.Token) {
	for parser.peekToken() != token.EOF {
		for _, tkn := range tokens {
			if parser.peekToken() == tkn {
				return
			}
		}
		parser.readToken()
	}
}

func customError(parser *Parser, error string) {
	line := parser.lexer.Tokens[parser.index].Beginning.Line
	column := parser.lexer.Tokens[parser.index].Beginning.Column
	file := parser.lexer.FileName + ":" + strconv.Itoa(line) + ":" + strconv.Itoa(column)
	logger.Error(file + " " + error)

	parser.hadError = true
}

func unexpectedToken(parser *Parser, possible, got string) {
	red := "\x1b[0;31m"
	reset := "\x1b[0m"
	line := parser.lexer.Tokens[parser.index].Beginning.Line
	column := parser.lexer.Tokens[parser.index].Beginning.Column
	file := parser.lexer.FileName + ":" + strconv.Itoa(line) + ":" + strconv.Itoa(column)
	logger.Error(file+" "+"Unexpected token: "+red+parser.lexer.GetToken(parser.lexer.Tokens[parser.index])+reset, "possible", possible, "got", got)

	parser.hadError = true
}

func (node *Node) setLineColumn(parser Parser) {
	if parser.index == 0 {
		node.Line = parser.lexer.Tokens[parser.index].Beginning.Line
		node.Column = parser.lexer.Tokens[parser.index].Beginning.Column
	} else {
		node.Line = parser.lexer.Tokens[parser.index-1].Beginning.Line
		node.Column = parser.lexer.Tokens[parser.index-1].Beginning.Column
	}
}

func expectToken(parser *Parser, tkn token.Token) bool {
	if parser.peekToken() != tkn {
		// Wait for the peek since EOF is added by this method
		parser.hadError = true
		if parser.index >= len(parser.lexer.Tokens) {
			parser.index--
		}

		red := "\x1b[0;31m"
		reset := "\x1b[0m"
		line := parser.lexer.Tokens[parser.index].Beginning.Line
		column := parser.lexer.Tokens[parser.index].Beginning.Column
		file := parser.lexer.FileName + ":" + strconv.Itoa(line) + ":" + strconv.Itoa(column)
		if tkn == token.SEMICOLON {
			// There is a missing semicolon, specific message and line/column
			// We can just continue parsing
			parser.unreadToken()
			line := parser.lexer.Tokens[parser.index].Beginning.Line
			column := parser.lexer.Tokens[parser.index].End.Column
			file := parser.lexer.FileName + ":" + strconv.Itoa(line) + ":" + strconv.Itoa(column)
			logger.Error(file + " " + "Missing semicolon after: " + parser.lexer.GetLineUpToTokenIncluded(parser.lexer.Tokens[parser.index]))
			parser.readToken()
		} else if parser.peekToken() == token.IDENT {
			logger.Error(file+" "+"Unexpected token: "+parser.lexer.GetLineUpToToken(parser.lexer.Tokens[parser.index])+red+parser.lexer.GetToken(parser.lexer.Tokens[parser.index])+reset, "expected", tkn, "got", parser.lexer.Lexi[parser.lexer.Tokens[parser.index].Position-1])
			// no read to continue parsing
		} else {
			logger.Error(file+" "+"Unexpected token: "+parser.lexer.GetLineUpToToken(parser.lexer.Tokens[parser.index])+red+parser.lexer.GetToken(parser.lexer.Tokens[parser.index])+reset, "expected", tkn, "got", parser.peekToken())
			// no read to continue parsing
		}
		return false
	}
	parser.readToken()
	return true
}

func peekExpectToken(parser *Parser, tkn token.Token) {
	if parser.peekToken() != tkn {
		red := "\x1b[0;31m"
		reset := "\x1b[0m"
		line := parser.lexer.Tokens[parser.index].Beginning.Line
		column := parser.lexer.Tokens[parser.index].Beginning.Column
		file := parser.lexer.FileName + ":" + strconv.Itoa(line) + ":" + strconv.Itoa(column)
		logger.Error(file+" "+"Unexpected token: "+parser.lexer.GetLineUpToToken(parser.lexer.Tokens[parser.index])+red+parser.lexer.GetToken(parser.lexer.Tokens[parser.index])+reset, "expected", tkn, "got", parser.peekToken())
	}
}

func expectTokenIdent(parser *Parser, ident string, recovery []any) string {
	red := "\x1b[0;31m"
	reset := "\x1b[0m"
	line := parser.lexer.Tokens[parser.index].Beginning.Line
	column := parser.lexer.Tokens[parser.index].Beginning.Column
	file := parser.lexer.FileName + ":" + strconv.Itoa(line) + ":" + strconv.Itoa(column)
	if parser.peekToken() != token.IDENT {
		// don't read, just assume it's there and raise the error
		logger.Error(file+" "+"Unexpected token: "+parser.lexer.GetLineUpToToken(parser.lexer.Tokens[parser.index])+red+parser.lexer.GetToken(parser.lexer.Tokens[parser.index])+reset, "expected", ident, "got", parser.peekToken())
		// if next token is the right one (in recovery), assume the current token is right to continue parsing
		for _, r := range recovery {
			if parser.peekTokenFurther(1) == token.Token(r.(int)) {
				parser.readToken()
				return ""
			}
		}
		return ""
	}
	_, index := parser.readFullToken()
	if parser.lexer.Lexi[index-1] != ident {
		logger.Error(file+" "+"Unexpected token: "+parser.lexer.GetLineUpToToken(parser.lexer.Tokens[parser.index-1])+red+parser.lexer.GetToken(parser.lexer.Tokens[parser.index-1])+reset, "expected", ident, "got", parser.lexer.Lexi[index-1])
	}
	return parser.lexer.Lexi[index-1]
}

func expectTokens(parser *Parser, tkns []any) bool {
	found := false
	for i, tkn := range tkns {
		if t, ok := tkn.(int); ok {
			found = expectToken(parser, token.Token(t))
		} else {
			// expect identifier with name tkn
			if i < len(tkns)-1 {
				expectTokenIdent(parser, tkn.(string), []any{tkns[i+1]})
			} else {
				expectTokenIdent(parser, tkn.(string), []any{})
			}
		}
	}
	return found
}

func readFichier(parser *Parser) Node {
	node := Node{Type: "Fichier"}
	node.setLineColumn(*parser)
	expectTokens(parser, []any{token.WITH, "Ada", token.PERIOD, "Text_IO", token.SEMICOLON, token.USE, "Ada", token.PERIOD, "Text_IO", token.SEMICOLON, token.PROCEDURE})

	node.addChild(readIdent(parser))
	expectTokens(parser, []any{token.IS})
	node.addChild(readDeclStar(parser))
	expectTokens(parser, []any{token.BEGIN})
	node.addChild(readInstr_plus(parser))
	expectTokens(parser, []any{token.END})
	node.addChild(readIdent_opt(parser))
	expectTokens(parser, []any{token.SEMICOLON, token.EOF})
	return node
}

func readDecl(parser *Parser) Node {
	var node Node
	switch parser.peekToken() {
	case token.PROCEDURE:
		parser.readToken()
		node = Node{Type: "DeclProcedure"}
		node.setLineColumn(*parser)
		node.addChild(readIdent(parser))
		node.addChild(readParams_opt(parser))
		expectTokens(parser, []any{token.IS})
		node.addChild(readDeclStar(parser))
		expectTokens(parser, []any{token.BEGIN})
		node.addChild(readInstr_plus(parser))
		expectTokens(parser, []any{token.END})
		node.addChild(readIdent_opt(parser))
		expectTokens(parser, []any{token.SEMICOLON})
	case token.TYPE:
		parser.readToken()
		node = Node{Type: "DeclType"}
		node.setLineColumn(*parser)
		node.addChild(readIdent(parser))
		node.addChild(readDecl2(parser))
	case token.FUNCTION:
		parser.readToken()
		node = Node{Type: "DeclFunction"}
		node.setLineColumn(*parser)
		node.addChild(readIdent(parser))
		node.addChild(readParams_opt(parser))
		expectTokens(parser, []any{token.RETURN})
		node.addChild(readType_r(parser))
		expectTokens(parser, []any{token.IS})
		node.addChild(readDeclStar(parser))
		expectTokens(parser, []any{token.BEGIN})
		node.addChild(readInstr_plus(parser))
		expectTokens(parser, []any{token.END})
		node.addChild(readIdent_opt(parser))
		expectTokens(parser, []any{token.SEMICOLON})
	case token.IDENT:
		node = Node{Type: "DeclVar"}
		node.setLineColumn(*parser)
		node.addChild(readIdent_plus_comma(parser))
		expectTokens(parser, []any{token.COLON})
		node.addChild(readType_r(parser))
		node.addChild(readInit(parser))
		expectTokens(parser, []any{token.SEMICOLON})
	default:
		parser.advance([]token.Token{token.PROCEDURE, token.IDENT, token.TYPE, token.FUNCTION, token.SEMICOLON})
		unexpectedToken(parser, "procedure type function ident", parser.peekTokenToString())
	}
	return node
}

func readDecl2(parser *Parser) Node {
	var node Node
	switch parser.readToken() {
	case token.IS:
		node = Node{Type: "DeclTypeIs"}
		node.setLineColumn(*parser)
		node.addChild(readDecl3(parser))
	case token.SEMICOLON:
		node = Node{Type: "DeclTypeSemicolon"}
		node.setLineColumn(*parser)
	default:
		node = Node{Type: "DeclTypeSemicolon"}
		node.setLineColumn(*parser)
		unexpectedToken(parser, "is ;", parser.peekTokenToString())
	}
	return node
}

func readDecl3(parser *Parser) Node {
	var node Node
	switch parser.readToken() {
	case token.ACCESS:
		node = Node{Type: "DeclTypeAccess"}
		node.setLineColumn(*parser)
		node.addChild(readIdent(parser))
		expectTokens(parser, []any{token.SEMICOLON})
	case token.RECORD:
		node = Node{Type: "DeclTypeRecord"}
		node.setLineColumn(*parser)
		node.addChild(readChampsPlus(parser))
		expectTokens(parser, []any{token.END, token.RECORD, token.SEMICOLON})
	default:
		parser.advance([]token.Token{token.ACCESS, token.RECORD, token.SEMICOLON})
		unexpectedToken(parser, "access record", parser.peekTokenToString())
	}
	return node
}

func readInit(parser *Parser) Node {
	var node Node
	switch parser.peekToken() {
	case token.SEMICOLON:
		node = Node{Type: "InitSemicolon"}
		node.setLineColumn(*parser)
	case token.COLON:
		expectTokens(parser, []any{token.COLON, token.EQL})
		node = Node{Type: "Init"}
		node.setLineColumn(*parser)
		node.addChild(readExpr(parser))
	default:
		// Error recovery, if the next token is a valid expression start, assume there is a missing semicolon
		if parser.peekToken() == token.BEGIN {
			return node
		}
		if parser.peekToken() == token.EQL {
			// assume that this is a mistake and that the user meant to write := instead of =
			customError(parser, "Malformed assignment statement. Did you mean to use := instead of =?")

			parser.readToken()
			node = Node{Type: "Init"}
			node.setLineColumn(*parser)
			node.addChild(readExpr(parser))
			return node
		}
		unexpectedToken(parser, "; :", parser.peekTokenToString())
	}
	return node
}

func readDeclStar(parser *Parser) Node {
	var node Node
fix:
	switch parser.peekToken() {
	case token.PROCEDURE, token.IDENT, token.TYPE, token.FUNCTION:
		node = Node{Type: "DeclStarProcedure"}
		node.setLineColumn(*parser)
		node.addChild(readDecl(parser))
		node.addChild(readDeclStar(parser))
	case token.BEGIN:
		node = Node{Type: "DeclStarBegin"}
		node.setLineColumn(*parser)
	default:
		unexpectedToken(parser, "procedure ident type function begin", parser.peekTokenToString())
		parser.advance([]token.Token{token.PROCEDURE, token.IDENT, token.TYPE, token.FUNCTION, token.BEGIN})
		if parser.peekToken() == token.EOF {
			return node
		}
		goto fix
	}
	return node
}

func readChamps(parser *Parser) Node {
	peekExpectToken(parser, token.IDENT)

	node := Node{Type: "Champs"}
	node.setLineColumn(*parser)
	node.addChild(readIdent_plus_comma(parser))

	expectTokens(parser, []any{token.COLON})
	node.addChild(readType_r(parser))
	expectTokens(parser, []any{token.SEMICOLON})
	return node
}

func readChampsPlus(parser *Parser) Node {
	peekExpectToken(parser, token.IDENT)

	node := Node{Type: "ChampsPlus"}
	node.setLineColumn(*parser)
	node.addChild(readChamps(parser))
	node.addChild(readChampsPlus2(parser))
	return node
}

func readChampsPlus2(parser *Parser) Node {
	var node Node
	var tkn = parser.peekToken()
	switch tkn {
	case token.IDENT:
		node = Node{Type: "ChampsPlus2"}
		node.setLineColumn(*parser)
		node.addChild(readChamps(parser))
		node.addChild(readChampsPlus2(parser))
	case token.END:
		node = Node{Type: "ChampsPlus2End"}
		node.setLineColumn(*parser)
	default:
		// TODO
		node = Node{Type: "ChampsPlus2End"}
		node.setLineColumn(*parser)
		parser.advance([]token.Token{token.END})
		logger.Error("Unexpected token", "possible", "ident end", "got", tkn)
	}
	return node
}

func readType_r(parser *Parser) Node {
	var node Node
	switch parser.peekToken() {
	case token.IDENT:
		node = Node{Type: "TypeRIdent"}
		node.setLineColumn(*parser)
		node.addChild(readIdent(parser))
	case token.ACCESS:
		parser.readToken()

		node = Node{Type: "TypeRAccess"}
		node.setLineColumn(*parser)
		node.addChild(readIdent(parser))
	default:
		if parser.peekToken() == token.SEMICOLON {
			// Error recovery, just continue
			return node
		}
		unexpectedToken(parser, "ident access", parser.peekTokenToString())
	}
	return node
}

func readParams(parser *Parser) Node {
	expectToken(parser, token.LPAREN)

	node := Node{Type: "Params"}
	node.setLineColumn(*parser)
	node.addChild(readParamPlusSemicolon(parser))

	expectTokens(parser, []any{token.RPAREN})
	return node
}

func readParams_opt(parser *Parser) Node {
	var node Node
	switch parser.peekToken() {
	case token.IS, token.RETURN:
		node = Node{Type: "ParamsOpt"}
		node.setLineColumn(*parser)
	case token.LPAREN:
		node = Node{Type: "ParamsOptParams"}
		node.setLineColumn(*parser)
		node.addChild(readParams(parser))
	default:
		// error recovery, next token to match a successful expression should be is or return
		parser.advance([]token.Token{token.IS, token.RETURN})
		unexpectedToken(parser, "is return (", parser.peekTokenToString())
	}
	return node
}

func readParam(parser *Parser) Node {
	peekExpectToken(parser, token.IDENT)

	node := Node{Type: "Param"}
	node.setLineColumn(*parser)
	node.addChild(readIdent_plus_comma(parser))

	expectTokens(parser, []any{token.COLON})

	node.addChild(readModeOpt(parser))
	node.addChild(readType_r(parser))
	return node
}

func readParamPlusSemicolon(parser *Parser) Node {
	peekExpectToken(parser, token.IDENT)

	node := Node{Type: "ParamPlusSemicolon"}
	node.setLineColumn(*parser)
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
		node.setLineColumn(*parser)
		node.addChild(readParam(parser))
		node.addChild(readParamPlusSemicolon2(parser))
	case token.RPAREN:
		node = Node{Type: "ParamPlusSemicolon2RParen"}
		node.setLineColumn(*parser)
	default:
		// error recovery, assume the params are over
		node = Node{Type: "ParamPlusSemicolon2RParen"}
		node.setLineColumn(*parser)
		unexpectedToken(parser, "; )", parser.peekTokenToString())
	}
	return node
}

func readMode(parser *Parser) Node {
	expectToken(parser, token.IN)
	node := Node{Type: "ModeIn"}
	node.setLineColumn(*parser)
	node.addChild(readMode2(parser))
	return node
}

func readMode2(parser *Parser) Node {
	var node Node
	switch parser.peekToken() {
	case token.IDENT, token.ACCESS:
		node = Node{Type: "Mode2Ident"}
		node.setLineColumn(*parser)
	case token.OUT:
		parser.readToken()
		node = Node{Type: "Mode2Out"}
		node.setLineColumn(*parser)
	default:
		parser.advance([]token.Token{token.SEMICOLON, token.RPAREN})
		unexpectedToken(parser, "ident access out", parser.peekTokenToString())
	}
	return node
}

func readModeOpt(parser *Parser) Node {
	var node Node
	switch parser.peekToken() {
	case token.IDENT, token.ACCESS:
		node = Node{Type: "ModeOpt"}
		node.setLineColumn(*parser)
	case token.IN:
		node = Node{Type: "ModeOptMode"}
		node.setLineColumn(*parser)
		node.addChild(readMode(parser))
	default:
		if parser.peekToken() == token.OUT {
			customError(parser, "Malformed parameter. Did you mean to use in out instead?")
			parser.readToken()
			return node
		}
		parser.advance([]token.Token{token.SEMICOLON, token.RPAREN})
		unexpectedToken(parser, "ident access in", parser.peekTokenToString())
	}
	return node
}

func readExpr(parser *Parser) Node {
	var node Node
	switch parser.peekToken() {
	case token.IDENT, token.LPAREN, token.NOT, token.SUB, token.INT, token.CHAR, token.TRUE, token.FALSE, token.NULL, token.NEW, token.CHAR_TOK:
		node = Node{Type: "ExprIdent"}
		node.setLineColumn(*parser)
		node.addChild(readOr_expr(parser))
	default:
		unexpectedToken(parser, "ident ( not - int char true false null new char", parser.peekTokenToString())
		parser.advanceExpr([]token.Token{token.SEMICOLON, token.RPAREN, token.COLON, token.COMMA, token.RETURN, token.END})
		parser.exprError = false
	}
	return node
}

func readOr_expr(parser *Parser) Node {
	var node Node
	switch parser.peekToken() {
	case token.IDENT, token.LPAREN, token.NOT, token.SUB, token.INT, token.CHAR, token.TRUE, token.FALSE, token.NULL, token.NEW, token.CHAR_TOK:
		node = Node{Type: "OrExpr"}
		node.setLineColumn(*parser)
		node.addChild(readAnd_expr(parser))
		node = readOr_expr_tail(parser, &node)
	default:
		unexpectedToken(parser, "ident ( not - int char true false null new char", parser.peekTokenToString())
		parser.advanceExpr([]token.Token{token.SEMICOLON, token.RPAREN, token.COLON, token.COMMA, token.RETURN, token.END, token.BEGIN})
		parser.exprError = false
	}
	return node
}

func readOr_expr_tail(parser *Parser, nd *Node) Node {
	var node Node
	if nd == nil {
		node = Node{}
		node.setLineColumn(*parser)
	} else {
		node = *nd
	}
	for parser.peekToken() == token.OR {
		parser.readToken()
		prev := node
		node = Node{Type: "OrExprTailOr"}
		node.setLineColumn(*parser)
		if nd != nil {
			node.addChild(prev)
		}
		node = readOr_expr_tail2(parser, &node)
	}
	switch parser.peekToken() {
	case token.SEMICOLON, token.RPAREN, token.THEN, token.COMMA, token.LOOP:
		//node.Type = "OrExprTail"
		return node
	case token.PERIOD:
		if parser.peekTokenFurther(1) == token.PERIOD {
			node.Type = "OrExprTail"
			return node
		}
		node.Type = "OrExprTailPeriod"
		parser.readToken()
		expectTokens(parser, []any{token.PERIOD})
		parser.readToken()
	default:
		parser.advanceExpr([]token.Token{token.SEMICOLON, token.RPAREN, token.COLON, token.COMMA, token.RETURN, token.END, token.BEGIN})
		parser.exprError = false
		//logger.Fatal("Unexpected token", "possible", "or ; ) then , loop .", "got", parser.peekToken())
	}
	return node
}

func readOr_expr_tail2(parser *Parser, nd *Node) Node {
	var node Node
	switch parser.peekToken() {
	case token.ELSE:
		parser.readToken()
		prev := *nd
		node = Node{Type: "OrExprTail2Else"}
		node.setLineColumn(*parser)
		node.addChild(prev)
		node.addChild(readAnd_expr(parser))
	case token.IDENT, token.LPAREN, token.NOT, token.SUB, token.INT, token.CHAR, token.TRUE, token.FALSE, token.NULL, token.NEW, token.CHAR_TOK:
		prev := *nd
		node = Node{Type: "OrExprTail2"}
		node.setLineColumn(*parser)
		node.addChild(prev)
		node.addChild(readAnd_expr(parser))
	default:
		parser.advanceExpr([]token.Token{token.SEMICOLON, token.RPAREN, token.COLON, token.COMMA, token.RETURN, token.END, token.BEGIN})
		parser.exprError = false
		//logger.Fatal("Unexpected token", "possible", "else ident ( not - int char true false null new char", "got", parser.peekToken())
	}
	return node
}

func readAnd_expr(parser *Parser) Node {
	var node Node
	switch parser.peekToken() {
	case token.IDENT, token.LPAREN, token.NOT, token.SUB, token.INT, token.CHAR, token.TRUE, token.FALSE, token.NULL, token.NEW, token.CHAR_TOK:
		node = Node{Type: "AndExpr"}
		node.setLineColumn(*parser)
		node.addChild(readEquality_expr(parser))
		node = readAnd_expr_tail(parser, &node)
	default:
		unexpectedToken(parser, "ident ( not - int char true false null new char", parser.peekTokenToString())
		parser.advanceExpr([]token.Token{token.SEMICOLON, token.RPAREN, token.COLON, token.COMMA, token.RETURN, token.END, token.BEGIN})
		parser.exprError = false
	}
	return node
}

func readAnd_expr_tail(parser *Parser, nd *Node) Node {
	var node Node
	if nd == nil {
		node = Node{}
		node.setLineColumn(*parser)
	} else {
		node = *nd
	}
	for parser.peekToken() == token.AND {
		parser.readToken()
		prev := node
		node = Node{Type: "AndExprTailAnd"}
		node.setLineColumn(*parser)
		if nd != nil {
			node.addChild(prev)
		}
		node = readAnd_expr_tail2(parser, &node)
	}
	switch parser.peekToken() {
	case token.SEMICOLON, token.RPAREN, token.OR, token.THEN, token.COMMA, token.LOOP:
		//node.Type = "AndExprTail"
		return node
	case token.PERIOD:
		if parser.peekTokenFurther(1) == token.PERIOD {
			node.Type = "OrExprTail"
			return node
		}
		node.Type = "AndExprTailPeriod"
		parser.readToken()
		expectTokens(parser, []any{token.PERIOD})
		parser.readToken()
	default:
		if !parser.exprError {
			unexpectedToken(parser, "and ; ) or then , loop .", parser.peekTokenToString())
			parser.exprError = true
		}
	}
	return node
}

func readAnd_expr_tail2(parser *Parser, nd *Node) Node {
	var node Node
	switch parser.peekToken() {
	case token.THEN:
		parser.readToken()
		prev := *nd
		node = Node{Type: "AndExprTail2Then"}
		node.setLineColumn(*parser)
		node.addChild(prev)
		node.addChild(readEquality_expr(parser))
	case token.IDENT, token.LPAREN, token.NOT, token.SUB, token.INT, token.CHAR, token.TRUE, token.FALSE, token.NULL, token.NEW, token.CHAR_TOK:
		node = Node{Type: "AndExprTail2"}
		node.setLineColumn(*parser)
		prev := *nd
		node.addChild(prev)
		node.addChild(readEquality_expr(parser))
	default:
		if !parser.exprError {
			unexpectedToken(parser, "then ident ( not - int char true false null new char", parser.peekTokenToString())
			parser.exprError = true
		}
	}
	return node
}

func readEquality_expr(parser *Parser) Node {
	var node Node
	switch parser.peekToken() {
	case token.IDENT, token.LPAREN, token.NOT, token.SUB, token.INT, token.CHAR, token.TRUE, token.FALSE, token.NULL, token.NEW, token.CHAR_TOK:
		node = Node{Type: "EqualityExpr"}
		node.setLineColumn(*parser)
		node.addChild(readRelational_expr(parser))
		node.addChild(readEquality_expr_tail(parser))
	default:
		unexpectedToken(parser, "ident ( not - int char true false null new char", parser.peekTokenToString())
		parser.advanceExpr([]token.Token{token.SEMICOLON, token.RPAREN, token.COLON, token.COMMA, token.RETURN, token.END, token.BEGIN})
		parser.exprError = false
	}
	return node
}

func readEquality_expr_tail(parser *Parser) Node {
	var node Node
	switch parser.peekToken() {
	case token.EQL:
		parser.readToken()
		node = Node{Type: "EqualityExprTailEql"}
		node.setLineColumn(*parser)
		node.addChild(readRelational_expr(parser))
		node.addChild(readEquality_expr_tail(parser))
	case token.NEQ:
		parser.readToken()
		node = Node{Type: "EqualityExprTailNeq"}
		node.setLineColumn(*parser)
		node.addChild(readRelational_expr(parser))
		node.addChild(readEquality_expr_tail(parser))
	case token.SEMICOLON, token.RPAREN, token.OR, token.AND, token.THEN, token.NOT, token.COMMA, token.LOOP:
		node = Node{Type: "EqualityExprTail"}
		node.setLineColumn(*parser)
	case token.PERIOD:
		if parser.peekTokenFurther(1) == token.PERIOD {
			node = Node{Type: "OrExprTail"}
			node.setLineColumn(*parser)
			return node
		}
		node = Node{Type: "EqualityExprTailPeriod"}
		node.setLineColumn(*parser)
		parser.readToken()
		expectTokens(parser, []any{token.PERIOD})
		parser.readToken()
	default:
		if !parser.exprError {
			unexpectedToken(parser, "= /= ; ) or and then not , loop .", parser.peekTokenToString())
			parser.exprError = true
		}
	}
	return node
}

func readRelational_expr(parser *Parser) Node {
	var node Node
	switch parser.peekToken() {
	case token.IDENT, token.LPAREN, token.NOT, token.SUB, token.INT, token.CHAR, token.TRUE, token.FALSE, token.NULL, token.NEW, token.CHAR_TOK:
		node = Node{Type: "RelationalExpr"}
		node.setLineColumn(*parser)
		node.addChild(readAdditive_expr(parser))
		node.addChild(readRelational_expr_tail(parser))
	default:
		unexpectedToken(parser, "ident ( not - int char true false null new char", parser.peekTokenToString())
		parser.advanceExpr([]token.Token{token.SEMICOLON, token.RPAREN, token.COLON, token.COMMA, token.RETURN, token.END, token.BEGIN})
		parser.exprError = false
	}
	return node
}

func readRelational_expr_tail(parser *Parser) Node {
	var node Node
	switch parser.peekToken() {
	case token.LSS:
		parser.readToken()
		node = Node{Type: "RelationalExprTailLss"}
		node.setLineColumn(*parser)
		node.addChild(readAdditive_expr(parser))
		node.addChild(readRelational_expr_tail(parser))
	case token.LEQ:
		parser.readToken()
		node = Node{Type: "RelationalExprTailLeq"}
		node.setLineColumn(*parser)
		node.addChild(readAdditive_expr(parser))
		node.addChild(readRelational_expr_tail(parser))
	case token.GTR:
		parser.readToken()
		node = Node{Type: "RelationalExprTailGtr"}
		node.setLineColumn(*parser)
		node.addChild(readAdditive_expr(parser))
		node.addChild(readRelational_expr_tail(parser))
	case token.GEQ:
		parser.readToken()
		node = Node{Type: "RelationalExprTailGeq"}
		node.setLineColumn(*parser)
		node.addChild(readAdditive_expr(parser))
		node.addChild(readRelational_expr_tail(parser))
	case token.SEMICOLON, token.RPAREN, token.OR, token.AND, token.THEN, token.NOT, token.EQL, token.NEQ, token.COMMA, token.LOOP:
		node = Node{Type: "RelationalExprTail"}
		node.setLineColumn(*parser)
	case token.PERIOD:
		if parser.peekTokenFurther(1) == token.PERIOD {
			node = Node{Type: "OrExprTail"}
			node.setLineColumn(*parser)
			return node
		}
		node = Node{Type: "RelationalExprTailPeriod"}
		node.setLineColumn(*parser)
		parser.readToken()
		expectTokens(parser, []any{token.PERIOD})
		parser.readToken()
	default:
		if !parser.exprError {
			unexpectedToken(parser, "< <= > >= ; ) or and then not = /= , loop .", parser.peekTokenToString())
			parser.exprError = true
		}
	}
	return node
}

func readAdditive_expr(parser *Parser) Node {
	var node Node
	switch parser.peekToken() {
	case token.IDENT, token.LPAREN, token.NOT, token.SUB, token.INT, token.CHAR, token.TRUE, token.FALSE, token.NULL, token.NEW, token.CHAR_TOK:
		node = Node{Type: "AdditiveExpr"}
		node.setLineColumn(*parser)
		node.addChild(readMultiplicative_expr(parser))
		node = readAdditive_expr_tail(parser, &node)
	default:
		unexpectedToken(parser, "ident ( not - int char true false null new char", parser.peekTokenToString())
		parser.advanceExpr([]token.Token{token.SEMICOLON, token.RPAREN, token.COLON, token.COMMA, token.RETURN, token.END, token.BEGIN})
		parser.exprError = false
	}
	return node
}

func readAdditive_expr_tail(parser *Parser, nd *Node) Node {
	node := *nd
	for parser.peekToken() == token.ADD || parser.peekToken() == token.SUB {
		switch parser.peekToken() {
		case token.ADD:
			parser.readToken()
			prev := node
			node = Node{Type: "AdditiveExprTailAdd"}
			node.setLineColumn(*parser)
			node.addChild(prev)
			node.addChild(readMultiplicative_expr(parser))
		case token.SUB:
			parser.readToken()
			prev := node
			node = Node{Type: "AdditiveExprTailSub"}
			node.setLineColumn(*parser)
			node.addChild(prev)
			node.addChild(readMultiplicative_expr(parser))
		}
	}
	if parser.peekToken() == token.SEMICOLON || parser.peekToken() == token.RPAREN || parser.peekToken() == token.OR || parser.peekToken() == token.AND || parser.peekToken() == token.THEN || parser.peekToken() == token.NOT || parser.peekToken() == token.EQL || parser.peekToken() == token.NEQ || parser.peekToken() == token.LSS || parser.peekToken() == token.LEQ || parser.peekToken() == token.GTR || parser.peekToken() == token.GEQ || parser.peekToken() == token.COMMA || parser.peekToken() == token.LOOP {
		return node
	} else if parser.peekToken() == token.PERIOD {
		if parser.peekTokenFurther(1) == token.PERIOD {
			node.Type = "OrExprTail"
			return node
		}
		node.Type = "AdditiveExprTailPeriod"
		parser.readToken()
		expectTokens(parser, []any{token.PERIOD})
		parser.readToken()
		return node
	} else {
		if !parser.exprError {
			unexpectedToken(parser, "+ - ; ) or and then not = /= < <= > >= , loop .", parser.peekTokenToString())
			parser.exprError = true
		}
		return node
	}
}

func readMultiplicative_expr(parser *Parser) Node {
	var node Node
	switch parser.peekToken() {
	case token.IDENT, token.LPAREN, token.NOT, token.SUB, token.INT, token.CHAR, token.TRUE, token.FALSE, token.NULL, token.NEW, token.CHAR_TOK:
		node = Node{Type: "MultiplicativeExpr"}
		node.setLineColumn(*parser)
		node.addChild(readUnary_expr(parser))
		node = readMultiplicative_expr_tail(parser, &node)
	default:
		unexpectedToken(parser, "ident ( not - int char true false null new char", parser.peekTokenToString())
		parser.advanceExpr([]token.Token{token.SEMICOLON, token.RPAREN, token.COLON, token.COMMA, token.RETURN, token.END, token.BEGIN})
		parser.exprError = false
	}
	return node
}

func readMultiplicative_expr_tail(parser *Parser, nd *Node) Node {
	node := *nd
	for parser.peekToken() == token.MUL || parser.peekToken() == token.QUO || parser.peekToken() == token.REM {
		switch parser.peekToken() {
		case token.MUL:
			parser.readToken()
			prev := node
			node = Node{Type: "MultiplicativeExprTailMul"}
			node.setLineColumn(*parser)
			node.addChild(prev)
			node.addChild(readUnary_expr(parser))
		case token.QUO:
			parser.readToken()
			prev := node
			node = Node{Type: "MultiplicativeExprTailQuo"}
			node.setLineColumn(*parser)
			node.addChild(prev)
			node.addChild(readUnary_expr(parser))
		case token.REM:
			parser.readToken()
			prev := node
			node = Node{Type: "MultiplicativeExprTailRem"}
			node.setLineColumn(*parser)
			node.addChild(prev)
			node.addChild(readUnary_expr(parser))
		}
	}
	if parser.peekToken() == token.SEMICOLON || parser.peekToken() == token.RPAREN || parser.peekToken() == token.OR || parser.peekToken() == token.AND || parser.peekToken() == token.THEN || parser.peekToken() == token.NOT || parser.peekToken() == token.EQL || parser.peekToken() == token.NEQ || parser.peekToken() == token.LSS || parser.peekToken() == token.LEQ || parser.peekToken() == token.GTR || parser.peekToken() == token.GEQ || parser.peekToken() == token.ADD || parser.peekToken() == token.SUB || parser.peekToken() == token.COMMA || parser.peekToken() == token.LOOP {
		return node
	} else if parser.peekToken() == token.PERIOD {
		if parser.peekTokenFurther(1) == token.PERIOD {
			node.Type = "OrExprTail"
			return node
		}
		node.Type = "MultiplicativeExprTailPeriod"
		parser.readToken()
		expectTokens(parser, []any{token.PERIOD})
		parser.readToken()
		return node
	} else {
		if !parser.exprError {
			unexpectedToken(parser, "* / rem ; ) or and then not = /= < <= > >= + - , loop .", parser.peekTokenToString())
			parser.exprError = true
		}
		return node
	}
}

func readUnary_expr(parser *Parser) Node {
	var node Node
	switch parser.peekToken() {
	case token.SUB:
		parser.readToken()
		node = Node{Type: "UnaryExprSub"}
		node.setLineColumn(*parser)
		node.addChild(readUnary_expr(parser))
	case token.NOT:
		parser.readToken()
		node = Node{Type: "UnaryExprNot"}
		node.setLineColumn(*parser)
		node.addChild(readUnary_expr(parser))
	case token.IDENT, token.LPAREN, token.INT, token.CHAR, token.TRUE, token.FALSE, token.NULL, token.NEW, token.CHAR_TOK:
		node = Node{Type: "UnaryExpr"}
		node.setLineColumn(*parser)
		node.addChild(readPrimary_expr(parser))
	default:
		unexpectedToken(parser, "- not ident ( int char true false null new char", parser.peekTokenToString())
		parser.advanceExpr([]token.Token{token.SEMICOLON, token.RPAREN, token.COLON, token.COMMA, token.RETURN, token.END, token.BEGIN})
		parser.exprError = false
	}
	return node
}

func readPrimary_expr(parser *Parser) Node {
	var node Node
	switch parser.peekToken() {
	case token.INT:
		_, index := parser.readFullToken()
		node = Node{Type: "PrimaryExprInt", Index: index}
		node.setLineColumn(*parser)
	case token.CHAR:
		_, index := parser.readFullToken()
		node = Node{Type: "PrimaryExprChar", Index: index}
		node.setLineColumn(*parser)
	case token.TRUE:
		parser.readToken()
		node = Node{Type: "PrimaryExprTrue"}
		node.setLineColumn(*parser)
	case token.FALSE:
		parser.readToken()
		node = Node{Type: "PrimaryExprFalse"}
		node.setLineColumn(*parser)
	case token.NULL:
		parser.readToken()
		node = Node{Type: "PrimaryExprNull"}
		node.setLineColumn(*parser)
	case token.LPAREN:
		parser.readToken()
		node = Node{Type: "PrimaryExprLparen"}
		node.setLineColumn(*parser)
		node.addChild(readExpr(parser))
		expectTokens(parser, []any{token.RPAREN})
	case token.NOT:
		parser.readToken()
		node = Node{Type: "PrimaryExprNot"}
		node.setLineColumn(*parser)
	case token.NEW:
		parser.readToken()
		node = Node{Type: "PrimaryExprNew"}
		node.setLineColumn(*parser)
		node.addChild(readIdent(parser))
	case token.IDENT:
		node = Node{Type: "PrimaryExprIdent"}
		node.setLineColumn(*parser)
		node.addChild(readIdent(parser))
		node.addChild(readPrimary_expr2(parser))
	case token.CHAR_TOK:
		parser.readToken()
		node = Node{Type: "PrimaryExprCharTok"}
		node.setLineColumn(*parser)
		expectTokens(parser, []any{token.CAST, token.VAL, token.LPAREN})
		node.addChild(readExpr(parser))
		expectTokens(parser, []any{token.RPAREN})
	default:
		unexpectedToken(parser, "int char true false null ( not new ident char", parser.peekTokenToString())
		parser.advanceExpr([]token.Token{token.SEMICOLON, token.RPAREN, token.COLON, token.COMMA, token.RETURN, token.END, token.BEGIN})
		parser.exprError = false
	}
	return node
}

func readPrimary_expr2(parser *Parser) Node {
	var node Node
	switch parser.peekToken() {
	case token.LPAREN:
		parser.readToken()
		node = Node{Type: "PrimaryExpr2Lparen"}
		node.setLineColumn(*parser)
		node.addChild(readExpr_plus_comma(parser))
		expectTokens(parser, []any{token.RPAREN})
		node.addChild(readPrimary_expr3(parser))
	case token.SEMICOLON, token.RPAREN, token.OR, token.AND, token.THEN, token.NOT, token.EQL, token.NEQ, token.LSS, token.LEQ, token.GTR, token.GEQ, token.ADD, token.SUB, token.MUL, token.QUO, token.REM, token.COMMA, token.LOOP:
		node = Node{Type: "PrimaryExpr2"}
		node.setLineColumn(*parser)
		node.addChild(readAccess2(parser))
	case token.PERIOD:
		if parser.peekTokenFurther(1) == token.PERIOD {
			node = Node{Type: "OrExprTail"}
			node.setLineColumn(*parser)
			return node
		} else {
			node = Node{Type: "PrimaryExpr2Period"}
			node.setLineColumn(*parser)
			node.addChild(readAccess2(parser))
		}
	default:
		if !parser.exprError {
			unexpectedToken(parser, "( ; ) or and then not = /= < <= > >= + - * / rem , loop .", parser.peekTokenToString())
			parser.exprError = true
		}
	}
	return node
}

func readPrimary_expr3(parser *Parser) Node {
	var node Node
	switch parser.peekToken() {
	case token.PERIOD:
		if parser.peekTokenFurther(1) == token.PERIOD {
			node = Node{Type: "OrExprTail"}
			node.setLineColumn(*parser)
			return node
		}
		parser.readToken()
		if parser.peekToken() == token.PERIOD {
			node = Node{Type: "PrimaryExpr3DoublePeriod"}
			node.setLineColumn(*parser)
			parser.readToken()
		} else {
			node = Node{Type: "PrimaryExpr3Period"}
			node.setLineColumn(*parser)
			node.addChild(readIdent(parser))
			node.addChild(readAccess2(parser))
		}
	case token.SEMICOLON, token.RPAREN, token.OR, token.AND, token.THEN, token.NOT, token.EQL, token.NEQ, token.LSS, token.LEQ, token.GTR, token.GEQ, token.ADD, token.SUB, token.MUL, token.QUO, token.REM, token.COMMA, token.LOOP:
		node = Node{Type: "PrimaryExpr3"}
		node.setLineColumn(*parser)
	default:
		unexpectedToken(parser, ". ; ) or and then not = /= < <= > >= + - * / rem , loop .", parser.peekTokenToString())
	}
	return node
}

func readAccess2(parser *Parser) Node {
	var node Node
	switch parser.peekToken() {
	case token.PERIOD:
		if parser.peekTokenFurther(1) == token.PERIOD {
			node = Node{Type: "OrExprTail"}
			node.setLineColumn(*parser)
			return node
		}
		parser.readToken()
		if parser.peekToken() == token.PERIOD {
			node = Node{Type: "Access2DoublePeriod"}
			node.setLineColumn(*parser)
		} else {
			node = Node{Type: "Access2Period"}
			node.setLineColumn(*parser)
			node.addChild(readIdent(parser))
			node.addChild(readAccess2(parser))
		}
	case token.SEMICOLON, token.RPAREN, token.OR, token.AND, token.THEN, token.NOT, token.EQL, token.NEQ, token.LSS, token.LEQ, token.GTR, token.GEQ, token.ADD, token.SUB, token.MUL, token.QUO, token.REM, token.COMMA, token.LOOP:
		node = Node{Type: "Access2"}
		node.setLineColumn(*parser)
	default:
		node = Node{Type: "Access2"}
		node.setLineColumn(*parser)
		unexpectedToken(parser, ". ; ) or and then not = /= < <= > >= + - * / rem , loop", parser.peekTokenToString())
	}
	return node
}

func readExpr_plus_comma(parser *Parser) Node {
	var node Node
	switch parser.peekToken() {
	case token.IDENT, token.LPAREN, token.NOT, token.SUB, token.INT, token.CHAR, token.TRUE, token.FALSE, token.NULL, token.NEW, token.CHAR_TOK:
		node = Node{Type: "ExprPlusComma"}
		node.setLineColumn(*parser)
		node.addChild(readExpr(parser))
		node.addChild(readExpr_plus_comma2(parser))
	default:
		// TODO look at this
		unexpectedToken(parser, "ident ( not - int char true false null new char", parser.peekTokenToString())
		parser.advance2(token.RPAREN)
	}
	return node
}

func readExpr_plus_comma2(parser *Parser) Node {
	var node Node
	switch parser.peekToken() {
	case token.COMMA:
		parser.readToken()
		node = Node{Type: "ExprPlusComma2Comma"}
		node.setLineColumn(*parser)
		node.addChild(readExpr(parser))
		node.addChild(readExpr_plus_comma2(parser))
	case token.RPAREN:
		node = Node{Type: "ExprPlusComma2Rparen"}
		node.setLineColumn(*parser)
	default:
		node = Node{Type: "ExprPlusComma2Rparen"}
		node.setLineColumn(*parser)
		unexpectedToken(parser, ", )", parser.peekTokenToString())
	}
	return node
}

func readExpr_opt(parser *Parser) Node {
	var node Node
	switch parser.peekToken() {
	case token.IDENT, token.LPAREN, token.NOT, token.SUB, token.INT, token.CHAR, token.TRUE, token.FALSE, token.NULL, token.NEW, token.CHAR_TOK:
		node = Node{Type: "ExprOpt"}
		node.setLineColumn(*parser)
		node.addChild(readExpr(parser))
	case token.SEMICOLON:
		node = Node{Type: "ExprOptSemicolon"}
		node.setLineColumn(*parser)
	default:
		node = Node{Type: "ExprOptSemicolon"}
		node.setLineColumn(*parser)
		unexpectedToken(parser, "ident ( not - int char true false null new char ;", parser.peekTokenToString())
	}
	return node
}

func readInstr(parser *Parser) Node {
	node := Node{Type: "Instr"}
	node.setLineColumn(*parser)
	switch parser.peekToken() {
	case token.ACCESS:
		parser.readToken()
		node = Node{Type: "InstrAccess"}
		node.setLineColumn(*parser)
		expectTokens(parser, []any{token.COLON, token.EQL})
		node.addChild(readExpr(parser))
		expectTokens(parser, []any{token.SEMICOLON})
	case token.IDENT:
		node = Node{Type: "InstrIdent"}
		node.setLineColumn(*parser)
		node.addChild(readIdent(parser))
		node.addChild(readInstr2(parser))
	case token.RETURN:
		parser.readToken()
		node = Node{Type: "InstrReturn"}
		node.setLineColumn(*parser)
		node.addChild(readExpr_opt(parser))
		expectTokens(parser, []any{token.SEMICOLON})
	case token.BEGIN:
		parser.readToken()
		node = Node{Type: "InstrBegin"}
		node.setLineColumn(*parser)
		node.addChild(readInstr_plus(parser))
		expectTokens(parser, []any{token.END, token.SEMICOLON})
	case token.IF:
		parser.readToken()
		node = Node{Type: "InstrIf"}
		node.setLineColumn(*parser)
		node.addChild(readExpr(parser))
		expectTokens(parser, []any{token.THEN})
		node.addChild(readInstr_plus(parser))
		node.addChild(readElse_if_star(parser))
		node.addChild(readElse_instr_opt(parser))
		expectTokens(parser, []any{token.END, token.IF, token.SEMICOLON})
	case token.FOR:
		parser.readToken()
		node = Node{Type: "InstrFor"}
		node.setLineColumn(*parser)
		node.addChild(readIdent(parser))
		expectTokens(parser, []any{token.IN})
		node.addChild(readReverse_instr(parser))
		node.addChild(readExpr(parser))
		expectTokens(parser, []any{token.PERIOD, token.PERIOD})
		node.addChild(readExpr(parser))
		expectTokens(parser, []any{token.LOOP})
		node.addChild(readInstr_plus(parser))
		expectTokens(parser, []any{token.END, token.LOOP, token.SEMICOLON})
	case token.WHILE:
		parser.readToken()
		node = Node{Type: "InstrWhile"}
		node.setLineColumn(*parser)
		node.addChild(readExpr(parser))
		expectTokens(parser, []any{token.LOOP})
		node.addChild(readInstr_plus(parser))
		expectTokens(parser, []any{token.END, token.LOOP, token.SEMICOLON})
	default:
		// error recovery
		unexpectedToken(parser, "access ident return begin if for while", parser.peekTokenToString())
		parser.advance2(token.END)
	}
	return node
}

func readInstr2(parser *Parser) Node {
	var node Node
	switch parser.peekToken() {
	case token.IDENT, token.BEGIN, token.END, token.RETURN, token.ACCESS, token.COLON, token.ELSE, token.PERIOD, token.IF, token.FOR, token.WHILE, token.ELSIF:
		if parser.peekToken() != token.COLON && parser.peekToken() != token.PERIOD {
			parser.readToken()
		}
		node = Node{Type: "Instr2Ident"}
		node.setLineColumn(*parser)
		node.addChild(readInstr3(parser))
		if parser.peekToken() == token.EQL {
			parser.readToken()
			customError(parser, "Malformed assignment statement. Did you mean to use := instead of =?")
		} else {
			if parser.peekToken() == token.COLON && parser.peekTokenFurther(1) != token.EQL {
				parser.readToken()
				customError(parser, "Malformed assignment statement. Did you mean to use := instead of :?")
			} else {
				expectTokens(parser, []any{token.COLON, token.EQL})
			}
		}
		node.addChild(readExpr(parser))
		expectTokens(parser, []any{token.SEMICOLON})
	case token.SEMICOLON:
		parser.readToken()
		node = Node{Type: "Instr2Semicolon"}
		node.setLineColumn(*parser)
	case token.LPAREN:
		parser.readToken()
		node = Node{Type: "Instr2Lparen"}
		node.setLineColumn(*parser)
		node.addChild(readExpr_plus_comma(parser))
		expectTokens(parser, []any{token.RPAREN})
		node.addChild(readInstr4(parser))
		expectTokens(parser, []any{token.SEMICOLON})
	default:
		// TODO: look at this
		if parser.peekToken() == token.EQL {
			// assume that this is a mistake and that the user meant to write := instead of =
			customError(parser, "Malformed assignment statement. Did you mean to use := instead of =?")

			parser.readToken()
			node = Node{Type: "Instr2Ident"}
			node.setLineColumn(*parser)
			node.addChild(readExpr(parser))
			expectTokens(parser, []any{token.SEMICOLON})
			return node
		}
		unexpectedToken(parser, "ident begin end return access : else . if for while elsif ; (", parser.peekTokenToString())
	}
	return node
}

func readInstr3(parser *Parser) Node {
	var node Node
	switch parser.peekToken() {
	case token.COLON:
		//
		//expectTokens(parser, []any{token.COLON, token.EQL})
		node = Node{Type: ":="}
		node.setLineColumn(*parser)
		//parser.unreadTokens(2)
	case token.PERIOD:
		if parser.peekTokenFurther(1) == token.PERIOD {
			node = Node{Type: "OrExprTail"}
			node.setLineColumn(*parser)
			return node
		}
		parser.readToken()
		node = Node{Type: "Instr3Period"}
		node.setLineColumn(*parser)
		node.addChild(readIdent(parser))
		node.addChild(readInstr3(parser))
	default:
		// error recovery
		if parser.peekTokenFurther(1) == token.END {
			return node
		}
		if parser.peekToken() == token.EQL {
			// Deal with error later
			return node
		}
		parser.advance2(token.COLON, token.PERIOD)
		unexpectedToken(parser, ": .", parser.peekTokenToString())
	}
	return node
}

func readInstr4(parser *Parser) Node {
	var node Node
	switch parser.peekToken() {
	case token.SEMICOLON:
	case token.COLON:
		expectTokens(parser, []any{token.COLON, token.EQL})
		node = Node{Type: "Instr4Colon"}
		node.setLineColumn(*parser)
		node.addChild(readExpr(parser))
	default:
		// error recovery
		if parser.peekToken() == token.END || parser.peekToken() == token.SEMICOLON {
			return node
		}
		unexpectedToken(parser, "; :", parser.peekTokenToString())
		parser.advance([]token.Token{token.END, token.SEMICOLON, token.IF})
	}
	return node
}

func readInstr_plus(parser *Parser) Node {
	var node Node
	switch parser.peekToken() {
	case token.BEGIN, token.RETURN, token.ACCESS, token.IF, token.FOR, token.WHILE, token.IDENT:
		node = Node{Type: "InstrPlus"}
		node.setLineColumn(*parser)
		node.addChild(readInstr(parser))
		node.addChild(readInstr_plus2(parser))
	default:
		if parser.peekToken() == token.END || parser.peekToken() == token.ELSE || parser.peekToken() == token.ELSIF {
			return node
		}
		parser.advance([]token.Token{token.END, token.ELSE, token.ELSIF})
		unexpectedToken(parser, "begin return access if for while ident", parser.peekTokenToString())
	}
	return node
}

func readInstr_plus2(parser *Parser) Node {
	node := Node{Type: "InstrPlus2"}
	node.setLineColumn(*parser)
	switch parser.peekToken() {
	case token.BEGIN, token.RETURN, token.ACCESS, token.IF, token.FOR, token.WHILE, token.IDENT:
		node.addChild(readInstr(parser))
		node.addChild(readInstr_plus2(parser))
	case token.END, token.ELSE, token.ELSIF:
	default:
		parser.advance([]token.Token{token.RPAREN, token.SEMICOLON})
		unexpectedToken(parser, "begin return access if for while ident", parser.peekTokenToString())
	}
	return node
}

func readElse_if(parser *Parser) Node {
	var node Node
	switch parser.peekToken() {
	case token.ELSIF:
		parser.readToken()
		node = Node{Type: "ElseIf"}
		node.setLineColumn(*parser)
		node.addChild(readExpr(parser))
		expectTokens(parser, []any{token.THEN})
		node.addChild(readInstr_plus(parser))
	default:
		// Impossible since this is called if there is an elsif
		unexpectedToken(parser, "elsif", parser.peekTokenToString())
	}
	return node
}

func readElse_if_star(parser *Parser) Node {
	var node Node
	switch parser.peekToken() {
	case token.ELSIF:
		node = Node{Type: "ElseIfStarElsif"}
		node.setLineColumn(*parser)
		node.addChild(readElse_if(parser))
		node.addChild(readElse_if_star(parser))
	case token.ELSE, token.END, token.BEGIN, token.RETURN, token.ACCESS, token.IF, token.FOR, token.WHILE, token.IDENT:
		node = Node{Type: "ElseIfStar"}
		node.setLineColumn(*parser)
	default:
		node = Node{Type: "ElseIfStar"}
		node.setLineColumn(*parser)
		unexpectedToken(parser, "begin return access if for while end else ident", parser.peekTokenToString())
	}
	return node
}

func readElse_instr(parser *Parser) Node {
	node := Node{Type: "ElseInstr"}
	node.setLineColumn(*parser)
	switch parser.peekToken() {
	case token.ELSE:
		parser.readToken()
		node.addChild(readInstr_plus(parser))
	default:
		// Impossible since this is called if there is an else
		unexpectedToken(parser, "else", parser.peekTokenToString())
	}
	return node
}

func readElse_instr_opt(parser *Parser) Node {
	var node Node
	switch parser.peekToken() {
	case token.ELSE:
		node = Node{Type: "ElseInstrOptElse"}
		node.setLineColumn(*parser)
		node.addChild(readElse_instr(parser))
	case token.END:
		node = Node{Type: "ElseInstrOptEnd"}
		node.setLineColumn(*parser)
	default:
		unexpectedToken(parser, "else end", parser.peekTokenToString())

		// Go to end of if statement
		parser.advance([]token.Token{token.END, token.IF, token.SEMICOLON})
	}
	return node
}

func readReverse_instr(parser *Parser) Node {
	var node Node
	switch parser.peekToken() {
	case token.IDENT, token.LPAREN, token.NOT, token.SUB, token.INT, token.CHAR, token.TRUE, token.FALSE, token.NULL, token.NEW, token.CHAR_TOK:
		node = Node{Type: "ReverseInstr"}
		node.setLineColumn(*parser)
	case token.REVERSE:
		node = Node{Type: "ReverseInstrReverse"}
		node.setLineColumn(*parser)
		parser.readToken()
	default:
		unexpectedToken(parser, "ident ( not - int char true false null new char reverse", parser.peekTokenToString())
		parser.advance2(token.IDENT, token.LPAREN, token.NOT, token.SUB, token.INT, token.CHAR, token.TRUE, token.FALSE, token.NULL, token.NEW, token.CHAR_TOK, token.REVERSE)
	}
	return node
}

func readIdent(parser *Parser) Node {
	peekExpectToken(parser, token.IDENT)

	//node := Node{Type: "Ident : " }
	_, index := parser.readFullToken()
	node := Node{Type: "Ident"}
	node.setLineColumn(*parser)
	node.Index = index

	node.setLineColumn(*parser)

	return node
}

func readIdent_opt(parser *Parser) Node {
	node := Node{Type: "IdentOpt"}
	node.setLineColumn(*parser)
	switch parser.peekToken() {
	case token.SEMICOLON:
		return node
	case token.IDENT:
		node.addChild(readIdent(parser))
	default:
		unexpectedToken(parser, "ident ;", parser.peekTokenToString())
		return node
	}
	return node
}

func readIdent_plus_comma(parser *Parser) Node {
	var node Node
	switch parser.peekToken() {
	case token.IDENT:
		node = Node{Type: "IdentPlusComma"}
		node.setLineColumn(*parser)
		node.addChild(readIdent(parser))
		node.addChild(readIdent_plus_comma2(parser))
	default:
		unexpectedToken(parser, "ident", parser.peekTokenToString())
	}
	return node
}

func readIdent_plus_comma2(parser *Parser) Node {
	var node Node
	switch parser.peekToken() {
	case token.SEMICOLON:
		node = Node{Type: "IdentPlusComma2Semicolon"}
		node.setLineColumn(*parser)
	case token.COMMA:
		node = Node{Type: "IdentPlusComma2Comma"}
		node.setLineColumn(*parser)
		parser.readToken()
		node.addChild(readIdent(parser))
		node.addChild(readIdent_plus_comma2(parser))
	case token.COLON:
		node = Node{Type: "IdentPlusComma2Colon"}
		node.setLineColumn(*parser)
	default:
		// If there is an ident after, it might just be a missing colon
		if parser.peekToken() == token.IDENT {
			return node
		}
		unexpectedToken(parser, "; , :", parser.peekTokenToString())
	}
	return node

}
