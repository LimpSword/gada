package parser

import (
	"encoding/json"
	"fmt"
)

type Graph struct {
	gmap      map[int]map[int]struct{}
	types     map[int]string
	terminals []int
	fathers   map[int]int
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

func addNodes(node *Node, graph *Graph) {
	fatherId := index
	graph.gmap[index] = make(map[int]struct{})
	graph.types[index] = node.Type
	if len(node.Children) == 0 {
		graph.terminals = append(graph.terminals, index)
	}
	for _, child := range node.Children {
		index++
		graph.fathers[index] = fatherId
		graph.gmap[fatherId][index] = struct{}{}
		addNodes(child, graph)
	}
}

var index int = 0

func createGraph(node Node) Graph {

	graph := Graph{}
	graph.gmap = make(map[int]map[int]struct{})
	graph.types = make(map[int]string)
	graph.terminals = make([]int, 0)
	graph.fathers = make(map[int]int)
	addNodes(&node, &graph)

	return graph
}

func clearchains(g Graph) {
	for _, term := range g.terminals {
		fmt.Println(term)
		tpTo := term
		for len(g.gmap[g.fathers[tpTo]]) == 1 {
			fmt.Println(term, tpTo)
			tpTo = g.fathers[tpTo]
		}
		if tpTo != term {
			delete(g.gmap[g.fathers[tpTo]], tpTo)
			g.gmap[g.fathers[tpTo]][term] = struct{}{}
		}
	}
}

func toAst(node Node) Graph {

	graph := createGraph(node)
	clearchains(graph)
	return graph

}
