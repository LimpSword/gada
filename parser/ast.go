package parser

import (
	"encoding/json"
	"fmt"
	"os"
)

type Graph map[int]map[int]struct{}

func (n Graph) toJson() string {
	b, err := json.MarshalIndent(n, "", "  ")
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

func addNodes(node *Node, graph Graph, father int) {
	father_id := index
	graph[index] = make(map[int]struct{})
	fmt.Println(node.Index)
	for _, child := range node.Children {
		index++
		graph[father_id][index] = struct{}{}
		addNodes(child, graph, index)
	}
}

var index int = 0

func createGraph(filePath string) (Graph, error) {
	jsonData, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	node, err := fromJSON(string(jsonData))

	graph := make(Graph)
	addNodes(node, graph, index)

	return graph, nil
}

func toAst(filePath string) (Graph, error) {

	graph, err := createGraph(filePath)

	fmt.Println(graph.toJson())

	return graph, err

}
