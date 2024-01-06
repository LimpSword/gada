package parser

import (
	"encoding/json"
	"fmt"
	"gada/lexer"
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
	// this change node Types deppending of his current type and childs
	switch {
	//case strings.HasSuffix(node.Type, "Tail"):
	//	return ""
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
			tpTo = g.fathers[tpTo]
		}
		if tpTo != term {
			fmt.Println(tpTo, term)
			delete(g.gmap[g.fathers[tpTo]], tpTo)
			delete(g.gmap[g.fathers[term]], term)
			g.gmap[g.fathers[tpTo]][term] = struct{}{}
			g.fathers[term] = g.fathers[tpTo]
			g.depth[term] = g.depth[tpTo]
		}
	}
}

func goUpNode(g *Graph, node int, name string) {
	// make a node replace his father keeping father childs
	dadNode := g.fathers[node]
	delete(g.gmap[g.fathers[dadNode]], dadNode)
	delete(g.gmap[dadNode], node)
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
	upTheNode(g, node)
}

func removeUselessTerminals(g *Graph) {

	//remove Terminals that are not usefull to the graph

	//uselessKeywords := make(map[string]struct{})
	//uselessKeywords["Access2"] = struct{}{}
	//uselessKeywords["InstrPlus2"] = struct{}{}

	for term, _ := range g.terminals {
		if g.types[term] == "Access2" || g.types[term] == "InstrPlus2" || g.types[term] == "DeclStarBegin" {
			delete(g.gmap[g.fathers[term]], term)
			delete(g.terminals, term)
			delete(g.meaningful, term)
		}
	}
}

func upTheNode(g *Graph, node int) {
	switch g.types[node] {
	case ":=":
		if g.types[g.fathers[node]] == "Instr2Ident" {
			goUpNode(g, node, ":=")
		}
		if g.types[g.fathers[node]] == "InstrIdent" {
			goUpNode(g, node, ":=")
		}
	}
}

func compactNodes(g *Graph) {
	for term, _ := range g.meaningful {
		upTheNode(g, term)
	}
}

func toAst(node Node, lexer lexer.Lexer) Graph {
	// return the ast as a graph structure (similar to a tree but not recursive)
	graph := createGraph(node, lexer)
	compactNodes(graph)
	//clearchains(graph)
	//removeUselessTerminals(graph)
	//clearchains(graph)

	return *graph

}
