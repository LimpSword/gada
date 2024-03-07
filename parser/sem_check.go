package parser

import (
	"fmt"
	"golang.org/x/exp/maps"
	"slices"
	"strconv"
)

func CheckSemantics(graph Graph) {
	//dfsSemantics(graph, 0)
	semCheck(graph, 0)
}

func getTypeSize(t string, scope Scope) int {
	switch t {
	case "integer":
		return 4
	case "character":
		return 1
	case "boolean":
		return 1
	default:
		// Is it a record?
		if symbol, ok := scope.Table[t]; ok {
			if symbol[0].Type() == Rec {
				size := 0
				for _, field := range symbol[0].(Record).Fields {
					size += getTypeSize(field, scope)
				}
				return size
			}
		}
		return 0
	}
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

func matchFunc(graph Graph, scope *Scope, name string, args []int) string {
	argstype := make(map[int]string)
	slices.Sort(args)
	for ind, val := range args {
		argstype[ind+1] = getReturnType(graph, scope, val)
	}
	if symbol, ok := scope.Table[name]; ok {
		for _, f := range symbol {
			if f.Type() == Func {
				fun := f.(Function)
				if fun.ParamCount == len(argstype) {
					buffer := []string{}
					breaked := false
					for i := 1; i <= len(argstype); i++ {
						if fun.Params[i].SType != argstype[i] {
							breaked = true
							break
						} else if fun.Params[i].IsParamOut {
							if whichFinal(graph, args[i-1]) != "identifier" || findStruct(graph, scope, args[i-1], false) == nil {
								buffer = append(buffer, "Parameter in out "+fun.Params[i].VName+" should be a variable currently is "+graph.types[args[i-1]])
							}
						}
					}

					if breaked {
						continue
					}
					for _, val := range buffer {
						logger.Error(val)
					}
					return f.(Function).ReturnType
					// TODO: check return type overloadding and return the correct one
				}
				continue
			} else {
				logger.Error(name + " is a " + symbol[0].Type() + " and not a function")
			}
		}

	} else {
		if scope.parent == nil {
			logger.Error(name + " function is undefined")
		} else {
			matchFunc(graph, scope.parent, name, args)
		}
	}
	return Unknown
}

func matchProc(graph Graph, scope *Scope, name string, args []int) string {
	argstype := make(map[int]string)
	slices.Sort(args)
	for ind, val := range args {
		argstype[ind+1] = getReturnType(graph, scope, val)
	}
	if symbol, ok := scope.Table[name]; ok {
		for _, f := range symbol {
			if f.Type() == Proc {
				fun := f.(Procedure)
				if fun.ParamCount == len(argstype) {
					buffer := []string{}
					breaked := false
					for i := 1; i <= len(argstype); i++ {
						if fun.Params[i].SType != argstype[i] {
							breaked = true
							break
						} else if fun.Params[i].IsParamOut {
							if whichFinal(graph, args[i-1]) != "identifier" || findStruct(graph, scope, args[i-1], false) == nil {
								buffer = append(buffer, "Parameter in out "+fun.Params[i].VName+" should be a variable currently is "+graph.types[args[i-1]])
							}
						}
					}

					if breaked {
						continue
					}
					for _, val := range buffer {
						logger.Error(val)
					}
				}
				continue
			} else {
				logger.Error(name + " is a " + symbol[0].Type() + " and not a procedure")
			}
		}

	} else {
		if scope.parent == nil {
			logger.Error(name + " procedure is undefined")
		} else {
			matchProc(graph, scope.parent, name, args)
		}
	}
	return Unknown
}

func whichFinal(graph Graph, node int) string {
	val := graph.types[node]
	if val == "True" || val == "False" {
		return "boolean"
	}
	if val[0] == '\'' {
		return "character"
	}
	_, err := strconv.Atoi(val)
	if err == nil {
		return "integer"
	} else {
		return "identifier"
	}
}

func getSymbol(graph Graph, scope *Scope, node int) string {
	// give the symbol type of the identifier
	name := graph.types[node]
	if symbol, ok := scope.Table[name]; ok {
		return symbol[0].Type()
	} else {
		if scope.parent == nil {
			logger.Error(name + " ident is undefined")
		} else {
			return findIdentifierType(graph, scope.parent, node)
		}
	}
	return Unknown
}

func getReturnType(graph Graph, scope *Scope, node int) string {
	// give the return type of the node
	children := maps.Keys(graph.gmap[node])
	slices.Sort(children)
	if len(children) == 0 {
		theType := whichFinal(graph, node)
		if theType == "identifier" {
			return findIdentifierType(graph, scope, node)
		}
		return theType
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
		return matchFunc(graph, scope, graph.types[children[0]], maps.Keys(graph.gmap[children[1]]))
	case "access":
		mainType := findIdentifierType(graph, scope, children[0])
		finalType := findAccessType(graph, scope, children[1], mainType)
		return finalType
	}
	return Unknown
}

func findIdentifierType(graph Graph, scope *Scope, node int) string {
	// give the return type of the identifier
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
			logger.Error(name + " ident is undefined")
		} else {
			return findIdentifierType(graph, scope.parent, node)
		}
	}
	return Unknown
}

func findStruct(graph Graph, scope *Scope, node int, log bool) *Variable {
	name := graph.types[node]
	if name == "access" {
		children := maps.Keys(graph.gmap[node])
		slices.Sort(children)
		return findStruct(graph, scope, children[0], log)
	}
	if symbol, ok := scope.Table[name]; ok {
		if variable, ok := symbol[0].(Variable); ok {
			return &variable
		} else {
			if log {
				logger.Error("left side of assignment " + name + " is not a variable")
			}
		}
	} else {
		if scope.parent == nil {
			if log {
				logger.Error("left side of assignment " + name + " is undefined")
			}
		} else {
			return findStruct(graph, scope.parent, node, log)
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
			logger.Error(name + " is a " + symbol[0].Type() + " and not a type")
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

// goUpScope: get the scope containing the variable and the total offset to reach it
func goUpScope(scope *Scope, name string) (*Scope, int) {
	if symbol, ok := scope.Table[name]; ok {
		for _, s := range symbol {
			if variable, ok := s.(Variable); ok {
				return scope, variable.Offset
			}
		}
	}

	totalOffset := scope.getCurrentOffset()
	if scope.parent == nil {
		// should never happen
		logger.Error(name + " variable is undefined")
	} else {
		parentScope, offset := goUpScope(scope.parent, name)
		return parentScope, totalOffset + offset
	}
	return nil, 0
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
		varStruct := findStruct(graph, scope, sorted[0], true)
		if varStruct != nil {
			if !varStruct.IsParamOut && varStruct.IsParamIn {
				logger.Error("Variable " + varStruct.VName + " is an in parameter and cannot be assigned")
			}
		}
	case "return":
		//todo check if the return type is the same as the function return type
		//todo check if the return is in a function

	case "call":
		symbolType := getSymbol(graph, scope, sorted[0])
		fmt.Println("symbolType", symbolType, graph.types[sorted[0]])
		if symbolType == Func { //todo handle after return only
			logger.Error("Cannot use call to function " + graph.types[sorted[0]] + " as a statement")
		} else if symbolType == Proc {
			//fmt.Println("Proc", graph.types[sorted[0]], maps.Keys(graph.gmap[sorted[1]]))
			matchProc(graph, scope, graph.types[sorted[0]], maps.Keys(graph.gmap[sorted[1]]))
		} else if symbolType == Rec {
			logger.Error("Cannot use call to type " + graph.types[sorted[0]] + " as a statement")
		} else if symbolType == Unknown {
			logger.Error("Cannot use call to " + graph.types[sorted[0]] + " as a statement")
		} else {
			logger.Error("Cannot use call to variable" + graph.types[sorted[0]] + " as a statement")
		}
	default:
		for _, child := range sorted {
			semCheck(graph, child)
		}
	}
}
