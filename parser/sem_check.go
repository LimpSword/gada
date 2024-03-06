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

func findAccessType(graph Graph, scope *Scope, node int, curType string) string {
	children := maps.Keys(graph.gmap[node])
	slices.Sort(children)
	if symbol, ok := scope.Table[curType]; ok {
		if symbol[0].Type() == Rec {
			if graph.types[node] == "access" {
				if _, ok1 := symbol[0].(Record).Fields[graph.types[children[0]]]; ok1 {
					newType := symbol[0].(Record).Fields[graph.types[children[0]]]
					return findAccessType(graph, scope, children[1], newType)
				} else {
					logger.Error(graph.types[children[0]] + " is not a field of " + curType)
				}
			} else {

				if _, ok1 := symbol[0].(Record).Fields[graph.types[node]]; ok1 {
					newType := symbol[0].(Record).Fields[graph.types[node]]
					return newType
				} else {
					logger.Error(graph.types[node] + " is not a field of " + curType)
				}
			}
		} else {
			logger.Error(curType + " is a " + symbol[0].Type() + " and not a record")
		}
	} else {
		if scope.parent == nil {
			logger.Error(curType + " type is undefined")
		} else {
			return findAccessType(graph, scope.parent, node, curType)
		}
	}
	return Unknown
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
	case "access":
		mainType := findIdentifierType(graph, scope, children[0])
		finalType := findAccessType(graph, scope, children[1], mainType)
		return finalType
	}
	return Unknown
}

func findIdentifierType(graph Graph, scope *Scope, node int) string {
	name := graph.types[node]
	if symbol, ok := scope.Table[name]; ok {
		if symbol[0].Type() == "integer" || symbol[0].Type() == "character" || symbol[0].Type() == "boolean" {
			return symbol[0].Type()
		} else {
			if symbol[0].Type() == Func { //it mean it's a function without arguments
				return symbol[0].(Function).ReturnType
			} else {
				return symbol[0].Type()
			}
		}
	} else {
		if scope.parent == nil {
			logger.Error(name + " variable is undefined")
		} else {
			return findIdentifierType(graph, scope.parent, node)
		}
	}
	return Unknown
}

func findStruct(graph Graph, scope *Scope, node int) *Variable {
	name := graph.types[node]
	if symbol, ok := scope.Table[name]; ok {
		if variable, ok := symbol[0].(Variable); ok {
			return &variable
		} else {
			logger.Error("left side of assignment " + name + " is not a variable")
		}
	} else {
		if scope.parent == nil {
			logger.Error("left side of assignment " + name + " is undefined")
		} else {
			return findStruct(graph, scope.parent, node)
		}
	}
	return nil
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

func findAccessName(graph Graph, node int, buffer string) string {
	if graph.types[node] == "access" {
		children := maps.Keys(graph.gmap[node])
		slices.Sort(children)
		buffer = buffer + graph.types[children[0]] + "."
		return findAccessName(graph, children[1], buffer)
	} else {
		return buffer + graph.types[node]
	}
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
	paramType := getSymbolType(graph.types[children[len(children)-1]])
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
		if graph.types[sorted[0]] != graph.types[sorted[len(sorted)-1]] {
			if graph.types[sorted[len(sorted)-1]] != "end" {
				logger.Error("Function " + graph.types[sorted[0]] + " end name do not match")
			}
		}
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
			if _, ok := recordElem.Fields[graph.types[childChild[0]]]; ok {
				logger.Error("Field " + graph.types[childChild[0]] + " is duplicate in record " + graph.types[sorted[0]] + " declaration")
			}
			recordElem.Fields[graph.types[childChild[0]]] = getSymbolType(graph.types[childChild[1]])
			findType(scope, getSymbolType(graph.types[childChild[1]]))
		}
	case ":=":
		varType := getReturnType(graph, scope, sorted[0])
		assignType := getReturnType(graph, scope, sorted[1])
		if varType != assignType {
			logger.Error("Type mismatch for variable: " + findAccessName(graph, sorted[0], "") + " is " + varType + " and was assigned to " + assignType)
		}
		varStruct := findStruct(graph, scope, sorted[0])
		if varStruct != nil {
			if !varStruct.IsParamOut && varStruct.IsParamIn {
				logger.Error("Variable " + varStruct.VName + " is an in parameter and cannot be assigned")
			}
		}
	default:
		for _, child := range sorted {
			semCheck(graph, child)
		}
	}
}
