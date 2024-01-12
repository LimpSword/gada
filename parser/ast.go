package parser

import (
	"encoding/json"
	"fmt"
	"gada/lexer"
	"sort"
	"strconv"
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
	switch {
	case node.Type == "Ident":
		return "Ident : " + lexer.Lexi[node.Index-1], true
	case node.Type == "PrimaryExprInt":
		return "Int : " + lexer.Lexi[node.Index-1], true
	case node.Type == "PrimaryExprChar":
		return "Char : " + lexer.Lexi[node.Index-1], true
	case node.Type == "PrimaryExprTrue":
		return "True", true
	case node.Type == "PrimaryExprFalse":
		return "False", true
	case node.Type == "PrimaryExprNull":
		return "Null", true

	case node.Type == ":=":
		return ":=", true
	case node.Type == "EqualityExpr":
		return "=", true
	case node.Type == "EqualityExprTailEql":
		return "=", true
	case node.Type == "AdditiveExpr":
		for _, child := range node.Children {
			if child.Type == "AdditiveExprTailAdd" {
				return "+", true
			} else if child.Type == "AdditiveExprTailSub" {
				return "-", true
			}
		}
		return node.Type, false
	case node.Type == "MultiplicativeExpr":
		for _, child := range node.Children {
			if child.Type == "MultiplicativeExprTailMul" {
				return "*", true
			} else if child.Type == "MultiplicativeExprTailDiv" {
				return "/", true
			}
		}
		return node.Type, false
	case node.Type == "RelationalExpr":
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
	case node.Type == "InstrIdent":
		for _, child := range node.Children {
			if child.Type == "Instr2Lparen" {
				return "call", true
			}
		}
	case node.Type == "InstrIf":
		return "if", true
	case node.Type == "ElseInstr":
		return "else", true
	case node.Type == "ElseIf":
		return "elif", true
	case node.Type == "DeclFunction":
		return "function", true
	case node.Type == "DeclVar":
		return "var", true
	default:
		return node.Type, false
	}
	return "", false
}

func meaningfulNode(node Node) bool {
	// check if a node is important on the graph
	return !(strings.HasSuffix(node.Type, "Tail")) // || node.Type == "Access2")
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
	fmt.Println(strconv.FormatInt(int64(node), 10) + " was cleaned")
	// remove a node from the graph
	delete(g.gmap[g.fathers[node]], node)
	delete(g.fathers, node)
	delete(g.terminals, node)
	delete(g.meaningful, node)
	delete(g.types, node)
	delete(g.gmap, node)
}

func goUpChilds(g *Graph, node int) {
	fmt.Println("upChilds"+strconv.FormatInt(int64(node), 10), g.fathers[node], len(g.gmap[g.fathers[node]]))
	dadNode := g.fathers[node]
	for child, _ := range g.gmap[node] {
		g.gmap[dadNode][child] = struct{}{}
		g.fathers[child] = dadNode
	}
	cleanNode(g, node)
}

func goUpReplaceNode(g *Graph, node int, name string) {
	// make a node replace his father keeping father childs can also change name
	fmt.Println("goUp for " + strconv.FormatInt(int64(node), 10))
	dadNode := g.fathers[node]
	delete(g.gmap[g.fathers[dadNode]], dadNode)
	delete(g.gmap[dadNode], node)
	fmt.Println(dadNode, node, g.gmap[dadNode])
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
	fmt.Println(g.fathers[node], node, g.gmap[dadNode])
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
	uselessKeywords := []string{"Access2", "InstrPlus2", "DeclStarBegin", "Instr2Semicolon", "ExprPlusComma2Rparen", "", "ElseIfStar", "IdentPlusComma2Colon"}

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
		if g.types[g.fathers[node]] == "InstrPlus2" {
			goUpChilds(g, g.fathers[node])
		}
	case "elif":
		if g.types[g.fathers[node]] == "ElseIfStarElsif" {
			goUpChilds(g, g.fathers[node])
		}
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
