package parser

import (
	"encoding/json"
	"fmt"
	"gada/lexer"
	"sort"
	"strings"
)

type Graph struct {
	gmap       map[int]map[int]struct{}
	types      map[int]string
	terminals  map[int]struct{}
	meaningful map[int]struct{}
	fathers    map[int]int
	depth      map[int]int
	nbNode     int
	lexer      *lexer.Lexer
}

func (g Graph) toJson() string {
	result := make(map[string]interface{})
	result["gmap"] = g.gmap
	result["types"] = g.types
	result["terminals"] = g.terminals
	result["meaningful"] = g.meaningful
	result["fathers"] = g.fathers

	b, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		fmt.Println(err)
	}
	return string(b)
}

func fromJSON(jsonStr string) (*Node, error) {
	var node Node
	err := json.Unmarshal([]byte(jsonStr), &node)
	if err != nil {
		return nil, err
	}

	return &node, nil
}

func nodeManagement(node Node, lexer lexer.Lexer) (string, bool) {
	// this change node Types depending on his current type and childs
	// this is the function choosing if a node has some interest and change their name
	//return node.Type, false
	switch node.Type {
	case "Fichier":
		return "file", true
		// ident
	case "Ident":
		return "Ident : " + lexer.Lexi[node.Index-1], true
		// Int
	case "PrimaryExprInt":
		return "Int : " + lexer.Lexi[node.Index-1], true
		// Char
	case "PrimaryExprChar":
		return "Char : " + lexer.Lexi[node.Index-1], true
		// True
	case "PrimaryExprTrue":
		return "True", true
		// False
	case "PrimaryExprFalse":
		return "False", true
		//Null
	case "PrimaryExprNull":
		return "Null", true
		// assignation
	case ":=":
		return ":=", true
		// equality
	case "EqualityExpr":
		for _, child := range node.Children {
			if child.Type == "EqualityExprTailEql" {
				return "=", true
			} else if child.Type == "EqualityExprTailNeq" {
				return "!=", true
			}
		}
		return node.Type, false
	case "EqualityExprTailEql":
		return "=", true
	case "EqualityExprTailNeq":
		return "!=", true
	// and or
	case "OrExpr":
		for _, child := range node.Children {
			if child.Type == "OrExprTailOr" {
				return "or", true
			}
		}
		return node.Type, false
	case "AndExprTail2Then": // always after the node and
		for _, child := range node.Children {
			if child.Type == "AndExprTailAnd" {
				return "and", true
			}
		}
		return node.Type, false
	case "AndExpr":
		for _, child := range node.Children {
			if child.Type == "AndExprTailAnd" {
				return "and", true
			}
		}
		return node.Type, false
		// operators
	case "AdditiveExpr":
		for _, child := range node.Children {
			if child.Type == "AdditiveExprTailAdd" {
				return "+", true
			} else if child.Type == "AdditiveExprTailSub" {
				return "-", true
			}
		}
		return node.Type, false
	case "AdditiveExprTailAdd":
		return "+", true
	case "AdditiveExprTailSub":
		return "-", true
	case "MultiplicativeExpr":
		for _, child := range node.Children {
			if child.Type == "MultiplicativeExprTailMul" {
				return "*", true
			} else if child.Type == "MultiplicativeExprTailQuo" {
				return "/", true
			} else if child.Type == "MultiplicativeExprTailRem" {
				return "rem", true
			}
		}
		return node.Type, false
	case "MultiplicativeExprTailRem":
		return "rem", true
	case "MultiplicativeExprTailQuo":
		return "/", true
	case "MultiplicativeExprTailMul":
		return "*", true
		// relational expr
	case "RelationalExpr":
		for _, child := range node.Children {
			if child.Type == "RelationalExprTailLss" {
				return "<", true
			} else if child.Type == "RelationalExprTailLeq" {
				return "<=", true
			} else if child.Type == "RelationalExprTailGtr" {
				return ">", true
			} else if child.Type == "RelationalExprTailGeq" {
				return ">=", true
			}
		}
		return node.Type, false
		// procedure call
	case "InstrIdent":
		for _, child := range node.Children {
			if child.Type == "Instr2Lparen" {
				return "call", true
			}
		}
		return node.Type, false
	case "PrimaryExprIdent":
		for _, child := range node.Children {
			if child.Type == "PrimaryExpr2Lparen" {
				return "call", true
			} else if child.Type == "PrimaryExpr2Period" { // call ident.ident
				return "call", true
			}
		}
		return node.Type, false
	case "Access2Period":
		return "call", true
	case "Instr3Period":
		return "call", true
		// unaryExpr
	case "UnaryExprNot":
		return "callNot", true
	case "UnaryExprSub":
		return "callSub", true
	// call multiple args
	case "ExprPlusComma":
		return "args", true
		// procedure
	case "DeclStarBegin":
		return "decl", true
	case "DeclProcedure":
		return "procedure", true
	case "InstrPlus":
		return "body", true
		// if else
	case "InstrIf":
		return "if", true
	case "ElseInstr":
		return "else", true
	case "ElseIf":
		return "elif", true
		// function
	case "DeclFunction":
		return "function", true
	case "Param":
		return "param", true
	case "ParamPlusSemicolon": // always after Params node easier way to handle
		return "params", true
	case "IdentPlusComma":
		return "sameType", true
		// variable declaration
	case "DeclVar":
		return "var", true
		// for loop
	case "InstrFor":
		return "for", true
		// while loop
	case "InstrWhile":
		return "while", true

	default:
		return node.Type, false
	}
	return "", false
}

func meaningfulNode(node Node) bool {
	// check if a node is important on the graph
	return !(strings.HasSuffix(node.Type, "Tail"))
}

func addNodes(node *Node, graph *Graph, lexer lexer.Lexer, depth int) {
	// add a tree recursively
	fatherId := graph.nbNode

	newType, meaningfull := nodeManagement(*node, lexer)

	graph.gmap[graph.nbNode] = make(map[int]struct{})
	graph.types[graph.nbNode] = newType
	graph.depth[graph.nbNode] = depth

	if len(node.Children) == 0 {
		meaningfull = true
		graph.terminals[graph.nbNode] = struct{}{}
	}

	if meaningfull {
		graph.meaningful[graph.nbNode] = struct{}{}
	}

	for _, child := range node.Children {
		if meaningfulNode(*child) {
			graph.nbNode++
			graph.fathers[graph.nbNode] = fatherId
			graph.gmap[fatherId][graph.nbNode] = struct{}{}
			addNodes(child, graph, lexer, depth+1)
		}
	}
}

func createGraph(node Node, lexer lexer.Lexer) *Graph {
	// initialyze the graph with the parsetree
	graph := Graph{}
	graph.gmap = make(map[int]map[int]struct{})
	graph.types = make(map[int]string)
	graph.terminals = make(map[int]struct{})
	graph.meaningful = make(map[int]struct{})
	graph.fathers = make(map[int]int)
	graph.depth = make(map[int]int)
	graph.nbNode = 0
	addNodes(&node, &graph, lexer, 1)

	return &graph
}

func clearchains(g *Graph) {
	// remove chains of single node link to each other
	for term, _ := range g.meaningful {
		tpTo := term
		for len(g.gmap[g.fathers[tpTo]]) == 1 {
			pastNode := tpTo
			tpTo = g.fathers[tpTo]
			if pastNode != term {
				cleanNode(g, pastNode)
			}
		}
		if tpTo != term {
			delete(g.gmap[g.fathers[tpTo]], tpTo)
			delete(g.gmap[g.fathers[term]], term)
			g.gmap[g.fathers[tpTo]][term] = struct{}{}
			g.fathers[term] = g.fathers[tpTo]
			g.depth[term] = g.depth[tpTo]
			cleanNode(g, tpTo)
		}
	}
}

func cleanNode(g *Graph, node int) {
	//fmt.Println(strconv.FormatInt(int64(node), 10) + " was cleaned")
	// remove a node from the graph
	delete(g.gmap[g.fathers[node]], node)
	delete(g.fathers, node)
	delete(g.terminals, node)
	delete(g.meaningful, node)
	delete(g.types, node)
	delete(g.gmap, node)
}

func goUpChilds(g *Graph, node int) {
	//fmt.Println("upChilds"+strconv.FormatInt(int64(node), 10), g.fathers[node], len(g.gmap[g.fathers[node]]))
	dadNode := g.fathers[node]
	for child, _ := range g.gmap[node] {
		g.gmap[dadNode][child] = struct{}{}
		g.fathers[child] = dadNode
	}
	cleanNode(g, node)
}

func fromChildToFather(g *Graph, node int) {
	//fmt.Println("upChilds"+strconv.FormatInt(int64(node), 10), g.fathers[node], len(g.gmap[g.fathers[node]]))
	dadNode := g.fathers[node]
	daddaddyNode := g.fathers[dadNode]
	// removing previous link
	delete(g.gmap[dadNode], node)
	delete(g.gmap[daddaddyNode], dadNode)
	// linking to father
	g.gmap[daddaddyNode][node] = struct{}{}
	g.fathers[node] = daddaddyNode
	// moving childs
	for child, _ := range g.gmap[node] {
		g.gmap[dadNode][child] = struct{}{}
		g.fathers[child] = dadNode
		delete(g.gmap[node], child)
	}
	// link to previous father
	g.gmap[node][dadNode] = struct{}{}
	g.fathers[dadNode] = node
	upTheNode(g, node)
}

func moveDown(g *Graph, node int) { // manage Instr3Period
	dadNode := g.fathers[node]
	smallestChild := -1
	for child := range g.gmap[dadNode] {
		if smallestChild == -1 || child < smallestChild {
			smallestChild = child
		}
	}
	if smallestChild == -1 || smallestChild == node {
		return
	}
	delete(g.gmap[dadNode], smallestChild)
	g.gmap[node][smallestChild] = struct{}{}
	g.fathers[smallestChild] = node
}

func handleUnary(g *Graph, node int, exp string, newExpr string) {
	g.nbNode++
	newNode := g.nbNode
	g.gmap[newNode] = make(map[int]struct{})
	g.types[newNode] = newExpr
	g.fathers[newNode] = node
	g.gmap[node][newNode] = struct{}{}
	g.depth[newNode] = g.depth[node] + 1
	g.types[node] = exp
	g.terminals[newNode] = struct{}{}
	g.meaningful[newNode] = struct{}{}
}

func goUpReplaceNode(g *Graph, node int, name string) {
	// make a node replace his father keeping father childs can also change name
	//fmt.Println("goUp for " + strconv.FormatInt(int64(node), 10))
	dadNode := g.fathers[node]
	delete(g.gmap[g.fathers[dadNode]], dadNode)
	delete(g.gmap[dadNode], node)
	//fmt.Println(dadNode, node, g.gmap[dadNode])
	g.gmap[g.fathers[dadNode]][node] = struct{}{}
	for child, _ := range g.gmap[dadNode] {
		if child != node {
			delete(g.gmap[dadNode], child)
			g.gmap[node][child] = struct{}{}
			g.fathers[child] = node
		}
	}
	g.meaningful[dadNode] = struct{}{}
	g.fathers[node] = g.fathers[dadNode]
	cleanNode(g, dadNode)
	//fmt.Println(g.fathers[node], node, g.gmap[dadNode])
	if !checkTerminal(g, node) {
		delete(g.terminals, node)
	}
	upTheNode(g, node)
}

func checkTerminal(g *Graph, node int) bool {
	return len(g.gmap[node]) == 0
}

func Contains(slice []string, term string) bool {
	for _, value := range slice {
		if term == value {
			return true
		}
	}
	return false
}

func removeUselessTerminals(g *Graph) {
	uselessKeywords := []string{"Access2", "InstrPlus2", "DeclStarBegin", "Instr2Semicolon", "ExprPlusComma2Rparen", "",
		"ElseIfStar", "IdentPlusComma2Colon", "ParamPlusSemicolon2RParen", "PrimaryExpr3", "InitSemicolon", "ParamsOpt",
		"ModeOpt", "ReverseInstr"}

	for term := range g.terminals {
		if Contains(uselessKeywords, g.types[term]) {
			cleanNode(g, term)
		}
	}
}

func upTheNode(g *Graph, node int) {
	switch g.types[node] {
	case ":=":
		if g.types[g.fathers[node]] == "Instr2Ident" {
			goUpReplaceNode(g, node, ":=")
		}
		if g.types[g.fathers[node]] == "InstrIdent" {
			goUpReplaceNode(g, node, ":=")
		}
		if g.types[g.fathers[node]] == "call" {
			fromChildToFather(g, node)
		}
		if g.types[g.fathers[node]] == "InstrPlus2" {
			goUpChilds(g, g.fathers[node])
		}
	case "elif":
		if g.types[g.fathers[node]] == "ElseIfStarElsif" {
			goUpChilds(g, g.fathers[node])
		}
	case "decl":
		if g.types[g.fathers[node]] == "DeclStarProcedure" {
			goUpReplaceNode(g, node, "decl")
		}
	case "sameType":
		for child, _ := range g.gmap[node] {
			if g.types[child] == "IdentPlusComma2Comma" {
				goUpChilds(g, child)
			}
		}
	case "args":
		for child, _ := range g.gmap[node] {
			if g.types[child] == "ExprPlusComma2Comma" {
				goUpChilds(g, child)
			}
		}
	case "for":
		if g.types[g.fathers[node]] == "InstrPlus2" {
			goUpReplaceNode(g, node, "for")
		}
	case "call":
		if g.types[g.fathers[node]] == "InstrPlus2" {
			goUpReplaceNode(g, node, "call")
		} else if g.types[g.fathers[node]] == ":=" || g.types[g.fathers[node]] == "call" {
			moveDown(g, node)
		}
	case "callNot":
		handleUnary(g, node, "call", "not")
	case "callSub":
		handleUnary(g, node, "call", "-")
	}
}

func compactNodes(g *Graph) {
	// Create a slice of keys from g.meaningful
	keys := make([]int, 0, len(g.meaningful))
	for key := range g.meaningful {
		keys = append(keys, key)
	}

	// Define a custom sorting function based on g.depth
	sort.Slice(keys, func(i, j int) bool {
		return g.depth[keys[i]] > g.depth[keys[j]]
	})

	// Iterate through sorted keys
	for _, term := range keys {
		upTheNode(g, term)
	}
}

func toAst(node Node, lexer lexer.Lexer) Graph {
	// return the ast as a graph structure (similar to a tree but not recursive)
	graph := createGraph(node, lexer)
	compactNodes(graph)
	clearchains(graph)
	removeUselessTerminals(graph)
	clearchains(graph)

	return *graph

}
