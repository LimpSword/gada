package parser

import (
	"encoding/json"
	"fmt"
	"os"
)

type Graph struct {
	gmap      map[int]map[int]struct{}
	types     map[int]string
	terminals []int
}

func (g Graph) toJson() string {
	result := make(map[string]interface{})
	result["gmap"] = g.gmap
	result["types"] = g.types
	result["terminals"] = g.terminals

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
	fmt.Printf("Processing node: %v, Index: %v, Type: %v, Children: %v\n", node, index, node.Type, node.Children)
	graph.gmap[index] = make(map[int]struct{})
	graph.types[index] = node.Type
	if len(node.Children) == 0 {
		fmt.Printf("Adding terminal: %v\n", index)
		graph.terminals = append(graph.terminals, index)
		fmt.Printf("Terminals: %v\n", graph.terminals)
	}
	for _, child := range node.Children {
		index++
		graph.gmap[fatherId][index] = struct{}{}
		addNodes(child, graph)
	}
}

var index int = 0

func createGraph(filePath string) (Graph, error) {
	jsonData, err := os.ReadFile(filePath)
	if err != nil {
		return Graph{}, err
	}

	node, err := fromJSON(string(jsonData))

	graph := Graph{}
	graph.gmap = make(map[int]map[int]struct{})
	graph.types = make(map[int]string)
	graph.terminals = make([]int, 0)
	addNodes(node, &graph)

	return graph, nil
}

func toAst(filePath string) (Graph, error) {

	graph, err := createGraph(filePath)

	return graph, err

}
