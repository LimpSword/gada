package parser

import (
	"golang.org/x/exp/maps"
	"slices"
	"strconv"
)

func CheckSemantics(graph Graph) {
	dfsSemantics(graph, 0)
}

func dfsSemantics(graph Graph, node int) {
	children := maps.Keys(graph.gmap[node])
	slices.Sort(children)
	for _, child := range children {
		scope := graph.scopes[child]
		if scope != nil {
			//fmt.Println(scope.String())

			switch graph.types[child] {
			case ":=":
				sorted := maps.Keys(graph.gmap[child])
				slices.Sort(sorted)
				var valueType = scope.getValueType(graph.types[sorted[1]])

				// TODO: check operations

				//fmt.Println(graph.types[sorted[0]], ":", valueType)

				// check if the variable is already declared with the same type
				checkingScope := scope
				for checkingScope != nil {
					//fmt.Println(checkingScope.Table)
					if symbol, ok := checkingScope.Table[graph.types[sorted[0]]]; ok {
						//fmt.Println("found")
						founded := false
						for _, symb := range symbol {
							if symb.Type() == valueType {
								founded = true
							}
						}
						if !founded {
							logger.Error("Type mismatch for variable: " + graph.types[sorted[0]])
						}
						break
					}
					checkingScope = checkingScope.parent
				}
			case "var":
			}
		}
		dfsSemantics(graph, child)
	}
}

func (scope *Scope) getValueType(val string) string {
	//fmt.Println("val: ", val)
	if val == "true" || val == "false" {
		return "bool"
	}
	if val[0] == '\'' {
		return "char"
	}
	_, err := strconv.Atoi(val)
	if err == nil {
		return "int"
	}
	// might be identifier
	var t string
	var currentScope = scope
	for currentScope != nil {
		if symbol, ok := currentScope.Table[val]; ok {
			t = symbol[0].Type()
			break
		}
		currentScope = currentScope.parent
	}
	if t == "" {
		t = Unknown
		logger.Warn("Unknown type for value: ", val)
	}
	return t
}
