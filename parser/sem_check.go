package parser

import (
	"fmt"
	"golang.org/x/exp/maps"
	"slices"
	"strconv"
)

func CheckSemantics(graph Graph) {
	//dfsSemantics(graph, 0)
}

func findType(scope *Scope, name string) string {
	if symbol, ok := scope.Table[name]; ok {
		return symbol[0].Type()
	}
	return Unknown
}

func compareFunc(f1 Function, f2 Function) bool {
	if f1.ParamCount == f2.ParamCount && f1.ReturnType == f2.ReturnType {
		for i := 0; i < f1.ParamCount; i++ {
			if f1.Params[i].SType != f2.Params[i].SType {
				return false
			}
		}
		return true
	}
	return false
}

func compareProc(f1 Procedure, f2 Procedure) bool {
	if f1.ParamCount == f2.ParamCount {
		for i := 0; i < f1.ParamCount; i++ {
			if f1.Params[i].SType != f2.Params[i].SType {
				return false
			}
		}
		return true
	}
	return false
}

func checkDecl(graph Graph, node int, currentScope *Scope) {
	sorted := maps.Keys(graph.gmap[node])
	slices.Sort(sorted)
	scope := *currentScope

	switch graph.types[node] {
	case "file":
		shift := 0
		if graph.types[sorted[1]] == "decl" {
			children := maps.Keys(graph.gmap[sorted[1]])
			for _, child := range children {
				checkDecl(graph, child, currentScope)
			}
			shift++
		}
		checkDecl(graph, sorted[1+shift], currentScope)
	case "function":
		funcParam := make(map[int]*Variable)
		funcElem := Function{FName: graph.types[sorted[0]], SType: Func, children: sorted, Params: funcParam}
		shift := 0
		if graph.types[sorted[1]] == "params" {

			for _, param := range maps.Keys(graph.gmap[sorted[1]]) {
				addParam(graph, param, &funcElem)
			}
			funcElem.ReturnType = getSymbolType(graph.types[sorted[2]])
			scope.addSymbol(funcElem)
			shift = 1
		}
		for _, fun := range scope.Table[funcElem.FName] {
			if fun.Type() == Func {
				if compareFunc(fun.(Function), funcElem) {
					logger.Error("Function redeclared with same parameters and return type")
				}
			} else {
				logger.Error(funcElem.FName + " is already declared in this scope")
				break
			}
		}
		if graph.types[sorted[2+shift]] == "decl" {
			funcScope := newScope(&scope)
			children := maps.Keys(graph.gmap[sorted[2+shift]])
			for _, child := range children {
				checkDecl(graph, child, funcScope)
			}
			shift++
		}
		checkDecl(graph, sorted[2+shift], currentScope)
	case "procedure":
		procParam := make(map[int]*Variable)
		procElem := Procedure{PName: graph.types[sorted[0]], PType: Proc, children: sorted, Params: procParam}
		shift := 0
		if graph.types[sorted[1]] == "params" {
			for _, param := range maps.Keys(graph.gmap[sorted[1]]) {
				addParamProc(graph, param, &procElem)
			}
			scope.addSymbol(procElem)
			shift = 1
		}

		for _, proc := range scope.Table[procElem.PName] {
			if proc.Type() == Proc {
				if compareProc(proc.(Procedure), procElem) {
					logger.Error("Procedure redeclared with same parameters")
				}
			} else {
				logger.Error(procElem.PName + " is already declared in this scope")
				break
			}
		}

		if graph.types[sorted[1+shift]] == "decl" {
			procScope := newScope(&scope)
			children := maps.Keys(graph.gmap[sorted[1+shift]])
			for _, child := range children {
				checkDecl(graph, child, procScope)
			}
			shift++
		}
		checkDecl(graph, sorted[1+shift], currentScope)
	case "for":
		fmt.Println("for")
		forScope := newScope(&scope)
		forScope.addSymbol(Variable{VName: graph.types[sorted[0]], SType: "integer"})
	case "var":
		if graph.types[sorted[0]] == "sameType" {
			for _, child := range maps.Keys(graph.gmap[sorted[0]]) {
				scope.addSymbol(Variable{VName: graph.types[sorted[child]], SType: getSymbolType(graph.types[sorted[1]])})
			}
		} else {
			scope.addSymbol(Variable{VName: graph.types[sorted[0]], SType: getSymbolType(graph.types[sorted[1]])})
		}
	case "type":
		recordElem := Record{RName: graph.types[sorted[0]], SType: Rec, Fields: make(map[string]string)}
		for _, child := range maps.Keys(graph.gmap[sorted[1]]) {
			childChild := maps.Keys(graph.gmap[child])
			slices.Sort(childChild)
			recordElem.Fields[graph.types[childChild[0]]] = getSymbolType(graph.types[childChild[1]])
		}
	default:
		for _, node := range sorted {
			checkDecl(graph, node, currentScope)
		}
	}
	graph.scopes[node] = &scope
}

//func dfsSemantics(graph Graph, node int) {
//	children := maps.Keys(graph.gmap[node])
//	slices.Sort(children)
//	scope := graph.scopes[node]
//	if scope != nil {
//		//fmt.Println(scope.String())
//
//		switch graph.types[node] {
//		case ":=":
//			//fmt.Println(graph.types[sorted[0]], ":", graph.types[sorted[1]])
//			var valueType = scope.getValueType(graph.types[sorted[1]])
//
//			// TODO: check operations
//
//			//fmt.Println(graph.types[sorted[0]], ":", valueType)
//
//			// check if the variable is already declared with the same type
//			checkingScope := scope
//			for checkingScope != nil {
//				//fmt.Println(checkingScope.Table)
//				if symbol, ok := checkingScope.Table[graph.types[sorted[0]]]; ok {
//					//fmt.Println("found")
//					founded := false
//					for _, symb := range symbol {
//						if symb.Type() == valueType {
//							founded = true
//						}
//					}
//					if !founded {
//						logger.Error("Type mismatch for variable: " + graph.types[sorted[0]])
//					}
//					break
//				}
//				checkingScope = checkingScope.parent
//			}
//		case "var":
//		}
//	}
//}

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
