package parser

import (
	"golang.org/x/exp/maps"
	"slices"
	"strconv"
)

func CheckSemantics(graph Graph) {
	//dfsSemantics(graph, 0)
	semCheck(graph, 0)
}

func matchFunc(scope *Scope, name string, args map[int]string) string {
	if symbol, ok := scope.Table[name]; ok {
		for _, f := range symbol {
			if f.Type() == Func {
				if f.(Function).ParamCount == len(args) {
					for i := 1; i <= len(args); i++ {
						if f.(Function).Params[i].SType != args[i] {
							continue
						}
					}
					return f.(Function).ReturnType
					// TODO: check return type overloadding and return the correct one
				}
				continue
			} else {
				logger.Error(name + " is a " + symbol[0].Type())
			}
		}

	} else {
		if scope.parent == nil {
			logger.Error(name + " type is undefined")
		} else {
			matchFunc(scope.parent, name, args)
		}
	}
	return Unknown
}

func getReturnType(graph Graph, scope *Scope, node int) string {
	children := maps.Keys(graph.gmap[node])
	slices.Sort(children)
	if len(children) == 0 {
		val := graph.types[node]
		if val == "true" || val == "false" {
			return "boolean"
		}
		if val[0] == '\'' {
			return "character"
		}
		_, err := strconv.Atoi(val)
		if err == nil {
			return "integer"
		} else {
			return findIdentifierType(graph, scope, node)
		}
	}

	switch graph.types[node] {
	case "+", "-", "*", "/":
		if getReturnType(graph, scope, children[0]) == "integer" && getReturnType(graph, scope, children[1]) == "integer" {
			return "integer"
		} else {
			logger.Error("Operator " + graph.types[node] + " should have integer operands")
		}
	case "and", "or":
		if getReturnType(graph, scope, children[0]) == "boolean" && getReturnType(graph, scope, children[1]) == "boolean" {
			return "boolean"
		} else {
			logger.Error("Operator " + graph.types[node] + " should have boolean operands")
		}
	case "call":
		args := make(map[int]string)
		sorted := maps.Keys(graph.gmap[children[1]])
		slices.Sort(sorted)
		for ind, val := range sorted {
			args[ind] = getReturnType(graph, scope, val)
		}
		return matchFunc(scope, graph.types[children[0]], args)
	}
	return Unknown
}

func findIdentifierType(graph Graph, scope *Scope, node int) string {
	name := graph.types[node]
	if symbol, ok := scope.Table[name]; ok {
		if symbol[0].Type() == "integer" || symbol[0].Type() == "character" || symbol[0].Type() == "boolean" {
			return symbol[0].Type()
		} else {
			if symbol[0].Type() == Func {
				return symbol[0].(Function).ReturnType
			} else {
				findIdentifierType(graph, scope.parent, node)
			}
		}
	} else {
		if scope.parent == nil {
			logger.Error(name + " type is undefined")
		} else {
			findIdentifierType(graph, scope.parent, node)
		}
	}
	return Unknown
}

func findType(scope *Scope, name string) string {
	if name == "integer" || name == "character" || name == "boolean" {
		return name
	}
	if symbol, ok := scope.Table[name]; ok {
		if symbol[0].Type() == Rec {
			return symbol[0].Name()
		} else {
			logger.Error(name + " is a " + symbol[0].Type())
		}
	} else {
		if scope.parent == nil {
			logger.Error(name + " type is undefined")
		} else {
			findType(scope.parent, name)
		}
	}
	return Unknown
}

func compareFunc(f1 Function, f2 Function) bool {
	if f1.ParamCount == f2.ParamCount && f1.ReturnType == f2.ReturnType {
		for i := 1; i <= f1.ParamCount; i++ {
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
		for i := 1; i <= f1.ParamCount; i++ {
			if f1.Params[i].SType != f2.Params[i].SType {
				return false
			}
		}
		return true
	}
	return false
}

func checkParam(graph Graph, node int, funcScope *Scope) {
	children := maps.Keys(graph.gmap[node])
	slices.Sort(children)
	paramType := getSymbolType(graph.types[children[1]])
	findType(funcScope, paramType)
}

func semCheck(graph Graph, node int) {
	sorted := maps.Keys(graph.gmap[node])
	slices.Sort(sorted)
	scope := graph.scopes[node]
	trashScope := newScope(nil)
	switch graph.types[node] {
	case "file":
		shift := 0
		if graph.types[sorted[1]] == "decl" {
			children := maps.Keys(graph.gmap[sorted[1]])
			for _, child := range children {
				semCheck(graph, child)
			}
			shift++
		}
		semCheck(graph, sorted[1+shift])

	case "function":
		funcParam := make(map[int]*Variable)
		funcElem := Function{FName: graph.types[sorted[0]], SType: Func, children: sorted, Params: funcParam}
		shift := 0
		if graph.types[sorted[1]] == "params" {
			child := maps.Keys(graph.gmap[sorted[1]])
			slices.Sort(child)
			for _, param := range child {
				addParam(graph, param, &funcElem, trashScope)
				checkParam(graph, param, scope)
			}
			shift = 1
		}
		funcElem.ReturnType = getSymbolType(graph.types[sorted[1+shift]])
		findType(scope, funcElem.ReturnType)
		countSame := 0
		for _, fun := range scope.Table[funcElem.FName] {
			if fun.Type() == Func {
				if compareFunc(fun.(Function), funcElem) {
					countSame++
					if countSame > 1 {
						logger.Error(funcElem.FName + " function redeclared with same parameters and return type")
						//break is we stop at first conflict
					}
				}
			} else {
				logger.Error(funcElem.FName + " is already declared in this scope")
				//break is we stop at first conflict
			}
		}
		if graph.types[sorted[2+shift]] == "decl" {
			children := maps.Keys(graph.gmap[sorted[2+shift]])
			for _, child := range children {
				semCheck(graph, child)
			}
			shift++
		}
		semCheck(graph, sorted[2+shift])
	case "procedure":
		procParam := make(map[int]*Variable)
		procElem := Procedure{PName: graph.types[sorted[0]], PType: Proc, children: sorted, Params: procParam}
		shift := 0
		if graph.types[sorted[1]] == "params" {
			child := maps.Keys(graph.gmap[sorted[1]])
			slices.Sort(child)
			for _, param := range child {
				addParamProc(graph, param, &procElem, trashScope)
				checkParam(graph, param, scope)
			}
			shift = 1
		}
		countSame := 0
		for _, proc := range scope.Table[procElem.PName] {
			if proc.Type() == Proc {
				if compareProc(proc.(Procedure), procElem) {
					countSame++
					if countSame > 1 {
						logger.Error("Procedure redeclared with same parameters")
					}
				}
			} else {
				logger.Error(procElem.PName + " is already declared in this scope")
				//break
			}
		}

		if graph.types[sorted[1+shift]] == "decl" {
			children := maps.Keys(graph.gmap[sorted[1+shift]])
			for _, child := range children {
				semCheck(graph, child)
			}
			shift++
		}
		semCheck(graph, sorted[1+shift])
	case "for":
		// todo stop variable assignation
	case "var":
		// check if something is already declared with the same name
		if graph.types[sorted[0]] == "sameType" {
			for _, child := range maps.Keys(graph.gmap[sorted[0]]) {
				if r, ok := scope.Table[graph.types[child]]; ok {
					if len(r) > 1 {
						logger.Error(graph.types[child] + " is already declared in this scope")
					}
				}
			}
		} else {
			if r, ok := scope.Table[graph.types[sorted[0]]]; ok {
				if len(r) > 1 {
					logger.Error(graph.types[sorted[0]] + " is already declared in this scope")
				}
			}
		}
		// check if the type exists
		declType := getSymbolType(graph.types[sorted[1]])
		findType(scope, declType)

	case "type":
		if r, ok := scope.Table[graph.types[node]]; ok {
			if len(r) > 1 {
				logger.Error(graph.types[node] + " is already declared in this scope")
			}
		}
		recordElem := Record{RName: graph.types[sorted[0]], SType: Rec, Fields: make(map[string]string)}
		for _, child := range maps.Keys(graph.gmap[sorted[1]]) {
			childChild := maps.Keys(graph.gmap[child])
			slices.Sort(childChild)
			recordElem.Fields[graph.types[childChild[0]]] = getSymbolType(graph.types[childChild[1]])
		}
	case ":=":
		varType := getReturnType(graph, scope, sorted[0])
		assignType := getReturnType(graph, scope, sorted[1])
		if varType != assignType {
			logger.Error("Type mismatch for variable: " + graph.types[sorted[0]] + " is " + varType + " and was assigned to " + assignType)
		}
	default:
		for _, child := range sorted {
			semCheck(graph, child)
		}
	}
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
