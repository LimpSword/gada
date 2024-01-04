package parser

import (
	"encoding/json"
	"fmt"
	"gada/lexer"
)

type Graph struct {
	gmap      map[int]map[int]struct{}
	types     map[int]string
	terminals []int
	fathers   map[int]int
	nbNode    int
	lexer     *lexer.Lexer
}

func (g Graph) toJson() string {
	result := make(map[string]interface{})
	result["gmap"] = g.gmap
	result["types"] = g.types
	result["terminals"] = g.terminals
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

func complementType(node Node, lexer lexer.Lexer) string {

	if node.Type == "Ident" {
		fmt.Println(node.Type)
		return "Ident : " + lexer.Lexi[node.Index-1]
	}
	return node.Type
}

func addNodes(node *Node, graph *Graph, lexer lexer.Lexer) {
	fatherId := graph.nbNode
	if len(node.Children) == 0 {
		graph.types[graph.nbNode] = complementType(*node, lexer)
		graph.terminals = append(graph.terminals, graph.nbNode)
	} else {
		graph.types[graph.nbNode] = node.Type
	}
	graph.gmap[graph.nbNode] = make(map[int]struct{})
	for _, child := range node.Children {
		//if len(child.Children) == 0 && strings.Contains(child.Type, "Tail") {
		//	return // we don't want to add the tail node
		//}
		graph.nbNode++
		graph.fathers[graph.nbNode] = fatherId
		graph.gmap[fatherId][graph.nbNode] = struct{}{}
		addNodes(child, graph, lexer)
	}
}

func createGraph(node Node, lexer lexer.Lexer) Graph {

	graph := Graph{}
	graph.gmap = make(map[int]map[int]struct{})
	graph.types = make(map[int]string)
	graph.terminals = make([]int, 0)
	graph.fathers = make(map[int]int)
	graph.nbNode = 0
	addNodes(&node, &graph, lexer)

	return graph
}

func clearchains(g Graph) {
	for _, term := range g.terminals {
		tpTo := term
		for len(g.gmap[g.fathers[tpTo]]) == 1 {
			tpTo = g.fathers[tpTo]
		}
		if tpTo != term {
			delete(g.gmap[g.fathers[tpTo]], tpTo)
			g.gmap[g.fathers[tpTo]][term] = struct{}{}
		}
	}
}

func toAst(node Node, lexer lexer.Lexer) Graph {

	graph := createGraph(node, lexer)
	clearchains(graph)
	return graph

}
