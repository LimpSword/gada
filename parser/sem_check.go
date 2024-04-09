package parser

import (
	"fmt"
	"golang.org/x/exp/maps"
	"slices"
	"strconv"
)

func CheckSemantics(graph Graph) {
	//dfsSemantics(graph, 0)
	semCheck(&graph, 0)
}

func getTypeSize(t string, scope Scope) int {
	switch t {
	case "integer":
		return 4
	case "character":
		return 4
	case "boolean":
		return 4
	default:
		// Is it a record?
		for {
			if symbol, ok := scope.Table[t]; ok {
				if symbol[0].Type() == Rec {
					size := 0
					for _, field := range symbol[0].(Record).Fields {
						size += getTypeSize(field, scope)
					}
					// The size needs to be a multiple of 4
					if size%4 != 0 {
						size += 4 - size%4
					}
					return size
				}
			}
			if scope.parent == nil {
				return 0
			}
			scope = *scope.parent
		}
		return 0
	}
}

func findAccessType(graph *Graph, scope *Scope, node int, curType string) string {
	children := maps.Keys(graph.gmap[node])
	slices.Sort(children)
	if symbol, ok := scope.Table[curType]; ok {
		if symbol[0].Type() == Rec {
			if graph.types[node] == "access" {
				if _, ok1 := symbol[0].(Record).Fields[graph.types[children[0]]]; ok1 {
					newType := symbol[0].(Record).Fields[graph.types[children[0]]]
					return findAccessType(graph, scope, children[1], newType)
				} else {
					fileName := graph.fileName
					line := strconv.Itoa(graph.line[node])
					column := strconv.Itoa(graph.column[node])
					logger.Error(fileName + ":" + line + ":" + column + " " + graph.types[children[0]] + " is not a field of " + curType)
				}
			} else {

				if _, ok1 := symbol[0].(Record).Fields[graph.types[node]]; ok1 {
					newType := symbol[0].(Record).Fields[graph.types[node]]
					return newType
				} else {
					fileName := graph.fileName
					line := strconv.Itoa(graph.line[node])
					column := strconv.Itoa(graph.column[node])
					logger.Error(fileName + ":" + line + ":" + column + " " + graph.types[node] + " is not a field of " + curType)
				}
			}
		} else {
			fileName := graph.fileName
			line := strconv.Itoa(graph.line[node])
			column := strconv.Itoa(graph.column[node])
			logger.Error(fileName + ":" + line + ":" + column + " " + curType + " is a " + symbol[0].Type() + " and not a record")
		}
	} else {
		if scope.parent == nil {
			if curType != "unknown" {
				fileName := graph.fileName
				line := strconv.Itoa(graph.line[node])
				column := strconv.Itoa(graph.column[node])
				logger.Error(fileName + ":" + line + ":" + column + " " + curType + " type is undefined")
			}
		} else {
			return findAccessType(graph, scope.parent, node, curType)
		}
	}
	return Unknown
}

func matchFuncReturn(graph *Graph, scope *Scope, node int, args []int, argstype map[int]map[string]struct{}, returnType map[string]struct{}) map[string]struct{} {
	// match function with expected return types

	matching := []Function{}
	returnTypes := make(map[string]struct{})

	name := graph.types[node]
	fileName := graph.fileName
	line := strconv.Itoa(graph.line[node])
	column := strconv.Itoa(graph.column[node])

	if symbol, ok := scope.Table[name]; ok {
		for _, f := range symbol {
			if f.Type() == Func {
				fun := f.(Function)
				if fun.ParamCount == len(argstype) && haveType(returnType, fun.ReturnType) {
					buffer := []string{}
					breaked := false
					for i := 1; i <= len(argstype); i++ {
						if !haveType(argstype[i], fun.Params[i].SType) {
							breaked = true
							break
						} else if fun.Params[i].IsParamOut {
							if whichFinal(graph, args[i-1]) != "identifier" || findStruct(graph, scope, args[i-1], false) == nil {
								buffer = append(buffer, fileName+":"+line+":"+column+" "+"Parameter in out "+fun.Params[i].VName+" should be a variable currently is "+graph.types[args[i-1]])
							}
						}
					}

					if breaked {
						continue
					}
					for _, val := range buffer {
						logger.Error(val)
					}
					matching = append(matching, f.(Function))
				}
				continue
			} else {
				logger.Error(fileName + ":" + line + ":" + column + " " + name + " is a " + f.Type() + " and not a function")
			}
		}
		if len(matching) > 1 {
			logger.Error(fileName + ":" + line + ":" + column + " " + name + " call is ambiguous")
			returnTypes[Unknown] = struct{}{}
			return returnTypes
		}
		if len(matching) > 0 {
			returnTypes[matching[0].ReturnType] = struct{}{}
			return returnTypes
		}

	}
	if scope.parent == nil {
		logger.Error(fileName + ":" + line + ":" + column + " " + name + " function is undefined")
		returnTypes[Unknown] = struct{}{}
		return returnTypes
	} else {
		return matchFuncReturn(graph, scope.parent, node, args, argstype, returnType)
	}
}

func matchFunc(graph *Graph, scope *Scope, node int, args []int, argstype map[int]map[string]struct{}) map[string]struct{} {

	matching := []Function{}
	returnTypes := make(map[string]struct{})

	name := graph.types[node]
	fileName := graph.fileName
	line := strconv.Itoa(graph.line[node])
	column := strconv.Itoa(graph.column[node])

	if symbol, ok := scope.Table[name]; ok {
		for _, f := range symbol {
			if f.Type() == Func {
				fun := f.(Function)
				if fun.ParamCount == len(argstype) {
					buffer := []string{}
					breaked := false
					for i := 1; i <= len(argstype); i++ {
						if !haveType(argstype[i], fun.Params[i].SType) {
							breaked = true
							break
						} else if fun.Params[i].IsParamOut {
							if whichFinal(graph, args[i-1]) != "identifier" || findStruct(graph, scope, args[i-1], false) == nil {
								buffer = append(buffer, fileName+":"+line+":"+column+" "+"Parameter in out "+fun.Params[i].VName+" should be a variable currently is "+graph.types[args[i-1]])
							}
						}
					}

					if breaked {
						continue
					}
					for _, val := range buffer {
						logger.Error(val)
					}
					matching = append(matching, f.(Function))
					//return f.(Function).ReturnType
				}
				continue
			} else {
				logger.Error(fileName + ":" + line + ":" + column + " " + name + " is a " + f.Type() + " and not a function")
			}
		}
		for _, f := range matching {
			if _, ok := returnTypes[f.ReturnType]; ok {
				fmt.Println(name + " call have multiple possibilities")
			} else {
				returnTypes[f.ReturnType] = struct{}{}
			}
		}
		if len(matching) > 0 {
			return returnTypes
		}

	}
	if scope.parent == nil {
		logger.Error(fileName + ":" + line + ":" + column + " " + name + " function is undefined")
		returnTypes[Unknown] = struct{}{}
		return returnTypes
	} else {
		return matchFunc(graph, scope.parent, node, args, argstype)
	}
}

func matchProc(graph *Graph, scope *Scope, node int, args []int, argstype map[int]map[string]struct{}) string {

	fmt.Printf("proc %s, args:%v, node : %d\n", graph.types[node], argstype, node)
	matching := []Procedure{}

	name := graph.types[node]
	fileName := graph.fileName
	line := strconv.Itoa(graph.line[node])
	column := strconv.Itoa(graph.column[node])

	if symbol, ok := scope.Table[name]; ok {
		for _, f := range symbol {
			if f.Type() == Proc {

				fun := f.(Procedure)
				if fun.ParamCount == len(argstype) {
					buffer := []string{}
					breaked := false
					for i := 1; i <= len(argstype); i++ {
						if !haveType(argstype[i], fun.Params[i].SType) {
							breaked = true
							break
						} else if fun.Params[i].IsParamOut {
							if whichFinal(graph, args[i-1]) != "identifier" || findStruct(graph, scope, args[i-1], false) == nil {
								buffer = append(buffer, fileName+":"+line+":"+column+" "+"Parameter in out "+fun.Params[i].VName+" should be a variable currently is "+graph.types[args[i-1]])
							}
						}
					}
					if breaked {
						continue
					}
					for _, val := range buffer {

						logger.Error(val)
					}
					matching = append(matching, f.(Procedure))
				}
				continue
			} else {
				logger.Error(fileName + ":" + line + ":" + column + " " + name + " is a " + f.Type() + " and not a procedure")
			}
		}
		if len(matching) > 1 {
			logger.Error(fileName + ":" + line + ":" + column + " " + name + " call is ambiguous")
		} else if len(matching) == 1 {
			return "found"
		}
	}
	if scope.parent == nil {
		logger.Error(fileName + ":" + line + ":" + column + " " + name + " procedure is undefined")
		return Unknown
	} else {
		return matchProc(graph, scope.parent, node, args, argstype)
	}
}

func whichFinal(graph *Graph, node int) string {
	// give the final type of the node
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

func getSymbol(graph *Graph, scope *Scope, node int) string {
	// give the symbol type of the identifier
	name := graph.types[node]

	if name == "access" {
		children := maps.Keys(graph.gmap[node])
		if len(children) != 0 {
			//mainType := findIdentifierType(graph, scope, children[0])
			//finalType := findAccessType(graph, scope, children[1], mainType)
			return "access"
		}

	}
	if symbol, ok := scope.Table[name]; ok {
		return symbol[0].Type()
	} else {
		if scope.parent == nil {
			fileName := graph.fileName
			line := strconv.Itoa(graph.line[node])
			column := strconv.Itoa(graph.column[node])
			logger.Error(fileName + ":" + line + ":" + column + " " + " ident " + name + " is undefined")
		} else {
			return getSymbol(graph, scope.parent, node)
		}
	}
	return Unknown
}

func haveType(types map[string]struct{}, wantedtype string) bool {
	if _, ok := types[wantedtype]; ok {
		return true
	} else {
		return false
	}
}

//func getExpectedTypes(graph *Graph, scope *Scope, node int, ind int) map[string]struct{} {
//	// give the expected types of the node
//	switch graph.types[graph.fathers[node]] {
//	case ":=":
//		return getExpectedTypes(graph, scope, graph.fathers[node], -1)
//	case "args":
//		childs := maps.Keys(graph.gmap[graph.fathers[node]])
//		slices.Sort(childs)
//		for argsInd, child := range childs {
//			if child == node {
//				return getExpectedTypes(graph, scope, graph.fathers[node], argsInd)
//			}
//		}
//	case "call":
//		childs := maps.Keys(graph.gmap[graph.fathers[node]])
//		slices.Sort(childs)
//		return matchFunc(graph, scope, graph.types[childs[0]], maps.Keys(graph.gmap[childs[1]]))
//	}
//
//	returnTypes := make(map[string]struct{})
//	childs := maps.Keys(graph.gmap[node])
//	slices.Sort(childs)
//	switch graph.types[node] {
//	case ":=":
//		return getReturnType(graph, scope, childs[0])
//	case "args":
//
//	}
//	returnTypes[Unknown] = struct{}{}
//	return returnTypes
//}

func getReturnType(graph *Graph, scope *Scope, node int, expectedReturn map[string]struct{}) map[string]struct{} {
	// give the return type of the node
	returnTypes := make(map[string]struct{})
	children := maps.Keys(graph.gmap[node])
	slices.Sort(children)
	if len(children) == 0 {

		theType := whichFinal(graph, node)
		if theType == "identifier" {
			//fmt.Printf("node: %s %s\n", graph.types[node], findIdentifierType(graph, scope, node))
			return findIdentifierType(graph, scope, node)
		}
		returnTypes[theType] = struct{}{}
		return returnTypes
	}
	switch graph.types[node] {
	case "+", "-", "*", "/", "rem":
		if haveType(getReturnType(graph, scope, children[0], expectedReturn), "integer") && haveType(getReturnType(graph, scope, children[1], expectedReturn), "integer") {
			returnTypes["integer"] = struct{}{}
			return returnTypes
		} else {
			fileName := graph.fileName
			line := strconv.Itoa(graph.line[node])
			column := strconv.Itoa(graph.column[node])
			logger.Error(fileName + ":" + line + ":" + column + " Operator " + graph.types[node] + " should have integer operands")
		}
	case "and", "or", "and then", "or else":
		if haveType(getReturnType(graph, scope, children[0], expectedReturn), "boolean") && haveType(getReturnType(graph, scope, children[1], expectedReturn), "boolean") {
			returnTypes["boolean"] = struct{}{}
			return returnTypes
		} else {
			fileName := graph.fileName
			line := strconv.Itoa(graph.line[node])
			column := strconv.Itoa(graph.column[node])
			logger.Error(fileName + ":" + line + ":" + column + " Operator " + graph.types[node] + " should have boolean operands")
		}
	case "not":
		if haveType(getReturnType(graph, scope, children[0], expectedReturn), "boolean") {
			returnTypes["boolean"] = struct{}{}
			return returnTypes
		} else {
			fileName := graph.fileName
			line := strconv.Itoa(graph.line[node])
			column := strconv.Itoa(graph.column[node])
			logger.Error(fileName + ":" + line + ":" + column + " Operator not should have boolean operands")
		}
	case ">", "<", ">=", "<=", "=", "/=":
		if haveType(getReturnType(graph, scope, children[0], expectedReturn), "integer") && haveType(getReturnType(graph, scope, children[1], expectedReturn), "integer") {
			returnTypes["boolean"] = struct{}{}
			return returnTypes
		} else {
			fileName := graph.fileName
			line := strconv.Itoa(graph.line[node])
			column := strconv.Itoa(graph.column[node])
			logger.Error(fileName + ":" + line + ":" + column + " Operator " + graph.types[node] + " should have integer operands")
		}
	case "call":
		if graph.types[children[0]] == "-" {
			if haveType(getReturnType(graph, scope, children[1], expectedReturn), "integer") {
				returnTypes["integer"] = struct{}{}
				return returnTypes
			} else {
				fileName := graph.fileName
				line := strconv.Itoa(graph.line[node])
				column := strconv.Itoa(graph.column[node])
				logger.Error(fileName + ":" + line + ":" + column + " Operator - should have integer operands")
			}
		} else {
			if len(expectedReturn) == 0 {
				return matchFunc(graph, scope, children[0], maps.Keys(graph.gmap[children[1]]), genArgsMap(graph, scope, maps.Keys(graph.gmap[children[1]])))
			} else {
				return matchFuncReturn(graph, scope, children[0], maps.Keys(graph.gmap[children[1]]), genArgsMap(graph, scope, maps.Keys(graph.gmap[children[1]])), expectedReturn)
			}
		}
	case "access":
		mainTypes := findIdentifierType(graph, scope, children[0])
		var mainType string
		for k := range mainTypes {
			mainType = k
			break // Exit the loop after extracting the key
		}
		finalType := findAccessType(graph, scope, children[1], mainType)
		returnTypes[finalType] = struct{}{}
		return returnTypes
	}
	returnTypes[Unknown] = struct{}{}
	return returnTypes
}

func findIdentifierType(graph *Graph, scope *Scope, node int) map[string]struct{} {
	// give the return type of the identifier
	name := graph.types[node]
	returnTypes := make(map[string]struct{})
	//fmt.Printf("scope: %v\n", scope.Table)
	if symbol, ok := scope.Table[name]; ok {
		if symbol[0].Type() == "integer" || symbol[0].Type() == "character" || symbol[0].Type() == "boolean" {
			returnTypes[symbol[0].Type()] = struct{}{}
			return returnTypes
		} else {
			if symbol[0].Type() == Func { //it means it's a function without arguments
				newNode := makeChild2(graph, node, "call", symbol[0].Name())
				return matchFunc(graph, scope, newNode, []int{}, make(map[int]map[string]struct{}))
			} else {
				returnTypes[symbol[0].Type()] = struct{}{}
				return returnTypes
			}
		}
	} else {
		if scope.parent == nil {
			fileName := graph.fileName
			line := strconv.Itoa(graph.line[node])
			column := strconv.Itoa(graph.column[node])
			logger.Error(fileName + ":" + line + ":" + column + " ident " + name + " is undefined")
		} else {
			return findIdentifierType(graph, scope.parent, node)
		}
	}
	returnTypes[Unknown] = struct{}{}
	return returnTypes
}

func findStruct(graph *Graph, scope *Scope, node int, log bool) *Variable {
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
				fileName := graph.fileName
				line := strconv.Itoa(graph.line[node])
				column := strconv.Itoa(graph.column[node])
				logger.Error(fileName + ":" + line + ":" + column + " left side of assignment " + name + " is not a variable")
			}
		}
	} else {
		if scope.parent == nil {
			if log {
				fileName := graph.fileName
				line := strconv.Itoa(graph.line[node])
				column := strconv.Itoa(graph.column[node])
				logger.Error(fileName + ":" + line + ":" + column + " left side of assignment " + name + " is undefined")
			}
		} else {
			return findStruct(graph, scope.parent, node, log)
		}
	}
	return nil
}

func findType(scope *Scope, name string) (string, error) {
	if name == "integer" || name == "character" || name == "boolean" {
		return name, nil
	}
	if symbol, ok := scope.Table[name]; ok {
		if symbol[0].Type() == Rec {
			return symbol[0].Name(), nil
		} else {
			return "", fmt.Errorf(name + " is a " + symbol[0].Type() + " and not a type")
		}
	} else {
		if scope.parent == nil {
			if name != "unknown" {
				return "", fmt.Errorf(name + " type is undefined")
			}
		} else {
			return findType(scope.parent, name)
		}
	}
	return Unknown, nil
}

func getDeclOffset(graph Graph, node int) int {
	scope := graph.scopes[node]
	if _, ok := scope.ScopeSymbol.(Procedure); ok {
		offset := 0
		for _, symbols := range scope.Table {
			for _, symbol := range symbols {
				if variable, ok := symbol.(Variable); ok {
					if !variable.IsParamIn && !variable.IsParamOut {
						offset += 4
					}
				}
			}
		}
		return offset
	}
	return 0
}

// goUpScope: get the scope containing the variable and the total offset to reach it
func goUpScope(scope *Scope, name string) (*Scope, int) {
	//totalOffset := scope.getCurrentOffset()
	if symbol, ok := scope.Table[name]; ok {
		for _, s := range symbol {
			if variable, ok := s.(Variable); ok {
				if variable.IsParamIn || variable.IsParamOut {
					return scope, -(variable.Offset - 4) + 16
				}
				var fixParamOffset = 0
				if _, ok := scope.ScopeSymbol.(Procedure); ok {
					fixParamOffset = 4 * scope.ScopeSymbol.(Procedure).ParamCount
				} else if _, ok := scope.ScopeSymbol.(Function); ok {
					fixParamOffset = 4 * scope.ScopeSymbol.(Function).ParamCount
				}
				return scope, -(variable.Offset - 4) + fixParamOffset
			}
		}
	}

	if scope.parent == nil {
		// should never happen
		logger.Warn(name + " variable is undefined")
	} else {
		parentScope, offset := goUpScope(scope.parent, name)
		/*if _, ok := scope.ScopeSymbol.(Procedure); ok {
			offset += 8
		}*/
		return parentScope, offset
	}
	return nil, 0
}

func getRegion(graph Graph, node int) int {
	scope := graph.scopes[node]
	return scope.Region
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

func findAccessName(graph *Graph, node int, buffer string) string {
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

func checkParam(graph *Graph, node int, funcScope *Scope) {
	children := maps.Keys(graph.gmap[node])
	slices.Sort(children)
	paramType := getSymbolType(graph.types[children[len(children)-1]])
	_, err := findType(funcScope, paramType)
	if err != nil {
		fileName := graph.fileName
		line := strconv.Itoa(graph.line[node])
		column := strconv.Itoa(graph.column[node])
		errorMessage := fileName + ":" + line + ":" + column + " " + err.Error()
		logger.Error(errorMessage)
	}
}

func findMotherFunc(scope *Scope) Symbol {
	if scope.ScopeSymbol != nil {
		return scope.ScopeSymbol
	} else {
		return findMotherFunc(scope.parent)
	}
}

func isHardReturn(graph *Graph, node int) bool {
	hasHardReturn := false
	if graph.types[node] == "return" {
		return true
	} else {
		children := maps.Keys(graph.gmap[node])
		slices.Sort(children)
		for ind, child := range children {
			if isHardReturn(graph, child) {
				hasHardReturn = true
			} else {
				if ind != 0 && (graph.types[node] == "if" || graph.types[node] == "elif") {
					return false
				}
			}
		}
	}
	return hasHardReturn
}

func updateReturn(graph *Graph, node int) {
	if _, ok := graph.hasReturn[node]; !ok {
		graph.hasReturn[node] = struct{}{}
		if (graph.types[node] != "function") && node != 0 {
			updateReturn(graph, graph.fathers[node])
		}
	}
}

func genArgsMap(graph *Graph, scope *Scope, args []int) map[int]map[string]struct{} {
	argstype := make(map[int]map[string]struct{})
	slices.Sort(args)
	for ind, val := range args {
		argstype[ind+1] = getReturnType(graph, scope, val, make(map[string]struct{}))
	}
	return argstype
}

func semCheck(graph *Graph, node int) {
	sorted := maps.Keys(graph.gmap[node])
	slices.Sort(sorted)
	scope := graph.scopes[node]
	trashScope := newScope(nil)
	switch graph.types[node] {
	case "file":
		shift := 0
		if graph.types[sorted[0]] != graph.types[sorted[len(sorted)-1]] {
			if graph.types[sorted[len(sorted)-1]] != "end" {
				fileName := graph.fileName
				line := strconv.Itoa(graph.line[sorted[len(sorted)-1]])
				column := strconv.Itoa(graph.column[sorted[len(sorted)-1]])
				logger.Error(fileName + ":" + line + ":" + column + " " + "Procedure " + graph.types[sorted[0]] + " end name do not match")
			}
		}
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
				fileName := graph.fileName
				line := strconv.Itoa(graph.line[node])
				column := strconv.Itoa(graph.column[node])
				logger.Error(fileName + ":" + line + ":" + column + " Function " + graph.types[sorted[0]] + " end name do not match")
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

		_, err := findType(scope, funcElem.ReturnType)
		if err != nil {
			fileName := graph.fileName
			line := strconv.Itoa(graph.line[node])
			column := strconv.Itoa(graph.column[node])
			errorMessage := fileName + ":" + line + ":" + column + " " + err.Error()
			logger.Error(errorMessage)
		}

		countSame := 0
		for _, fun := range scope.Table[funcElem.FName] {
			if fun.Type() == Func {
				if compareFunc(fun.(Function), funcElem) {
					countSame++
					if countSame > 1 {
						fileName := graph.fileName
						line := strconv.Itoa(graph.line[node])
						column := strconv.Itoa(graph.column[node])
						logger.Error(fileName + ":" + line + ":" + column + " " + funcElem.FName + " function redeclared with same parameters and return type")
						//break is we stop at first conflict
					}
				}
			} else {
				fileName := graph.fileName
				line := strconv.Itoa(graph.line[node])
				column := strconv.Itoa(graph.column[node])
				logger.Error(fileName + ":" + line + ":" + column + " " + funcElem.FName + " is already declared in this scope")
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
		if _, ok := graph.hasReturn[node]; !ok {
			fileName := graph.fileName
			line := strconv.Itoa(graph.line[node])
			column := strconv.Itoa(graph.column[node])
			logger.Error(fileName + ":" + line + ":" + column + " " + "Function " + funcElem.FName + " has no return statement")
		} else {
			if !isHardReturn(graph, sorted[2+shift]) {
				fileName := graph.fileName
				line := strconv.Itoa(graph.line[node])
				column := strconv.Itoa(graph.column[node])
				logger.Error(fileName + ":" + line + ":" + column + " " + "Function " + funcElem.FName + " may miss return statement")
			}
		}
	case "procedure":
		procParam := make(map[int]*Variable)
		procElem := Procedure{PName: graph.types[sorted[0]], PType: Proc, children: sorted, Params: procParam}
		shift := 0
		if graph.types[sorted[0]] != graph.types[sorted[len(sorted)-1]] {
			if graph.types[sorted[len(sorted)-1]] != "end" {
				fileName := graph.fileName
				line := strconv.Itoa(graph.line[node])
				column := strconv.Itoa(graph.column[node])
				logger.Error(fileName + ":" + line + ":" + column + " " + "Procedure " + graph.types[sorted[0]] + " end name do not match")
			}
		}
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
						fileName := graph.fileName
						line := strconv.Itoa(graph.line[node])
						column := strconv.Itoa(graph.column[node])
						logger.Error(fileName + ":" + line + ":" + column + " " + "Procedure redeclared with same parameters")
					}
				}
			} else {
				fileName := graph.fileName
				line := strconv.Itoa(graph.line[node])
				column := strconv.Itoa(graph.column[node])
				logger.Error(fileName + ":" + line + ":" + column + " " + procElem.PName + " is already declared in this scope")
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
		for _, child := range sorted {
			if graph.types[child] == "body" {
				semCheck(graph, child)
			}
		}
	case "var":
		// check if something is already declared with the same name
		if graph.types[sorted[0]] == "sameType" {
			for _, child := range maps.Keys(graph.gmap[sorted[0]]) {
				if r, ok := scope.Table[graph.types[child]]; ok {
					if len(r) > 1 {
						fileName := graph.fileName
						line := strconv.Itoa(graph.line[node])
						column := strconv.Itoa(graph.column[node])
						logger.Error(fileName + ":" + line + ":" + column + " " + graph.types[child] + " is already declared in this scope")
					}
				}
			}
		} else {
			if r, ok := scope.Table[graph.types[sorted[0]]]; ok {
				if len(r) > 1 {
					fileName := graph.fileName
					line := strconv.Itoa(graph.line[node])
					column := strconv.Itoa(graph.column[node])
					logger.Error(fileName + ":" + line + ":" + column + " " + graph.types[sorted[0]] + " is already declared in this scope")
				}
			}
		}
		// check if the type exists
		declType := getSymbolType(graph.types[sorted[1]])

		_, err := findType(scope, declType)
		if err != nil {
			fileName := graph.fileName
			line := strconv.Itoa(graph.line[sorted[1]])
			column := strconv.Itoa(graph.column[sorted[1]])
			errorMessage := fileName + ":" + line + ":" + column + " " + err.Error()
			logger.Error(errorMessage)
		}

	case "type":
		if r, ok := scope.Table[graph.types[node]]; ok {
			if len(r) > 1 {
				fileName := graph.fileName
				line := strconv.Itoa(graph.line[node])
				column := strconv.Itoa(graph.column[node])
				logger.Error(fileName + ":" + line + ":" + column + " " + graph.types[node] + " is already declared in this scope")
			}
		}
		recordElem := Record{RName: graph.types[sorted[0]], SType: Rec, Fields: make(map[string]string)}
		for _, child := range maps.Keys(graph.gmap[sorted[1]]) {
			childChild := maps.Keys(graph.gmap[child])
			slices.Sort(childChild)
			if _, ok := recordElem.Fields[graph.types[childChild[0]]]; ok {
				fileName := graph.fileName
				line := strconv.Itoa(graph.line[node])
				column := strconv.Itoa(graph.column[node])
				logger.Error(fileName + ":" + line + ":" + column + " " + "Field " + graph.types[childChild[0]] + " is duplicate in record " + graph.types[sorted[0]] + " declaration")
			}
			recordElem.Fields[graph.types[childChild[0]]] = getSymbolType(graph.types[childChild[1]])

			_, err := findType(scope, getSymbolType(graph.types[childChild[1]]))
			if err != nil {
				fileName := graph.fileName
				line := strconv.Itoa(graph.line[node])
				column := strconv.Itoa(graph.column[node])
				errorMessage := fileName + ":" + line + ":" + column + " " + err.Error()
				logger.Error(errorMessage)
			}

		}
	case ":=":
		if whichFinal(graph, sorted[0]) != "identifier" {
			fileName := graph.fileName
			line := strconv.Itoa(graph.line[sorted[0]])
			column := strconv.Itoa(graph.column[sorted[0]])
			logger.Error(fileName + ":" + line + ":" + column + " " + "Left side of assignment is not a variable")
		} else {
			if Contains([]string{Func, Proc, Rec}, getSymbol(graph, scope, sorted[0])) {
				fileName := graph.fileName
				line := strconv.Itoa(graph.line[sorted[0]])
				column := strconv.Itoa(graph.column[sorted[0]])
				logger.Error(fileName + ":" + line + ":" + column + " " + "Left side of assignment is not a variable")
			}
		}

		varTypes := getReturnType(graph, scope, sorted[0], make(map[string]struct{}))
		var varType string
		for k := range varTypes {
			varType = k
			break
		}
		assignTypes := getReturnType(graph, scope, sorted[1], varTypes)
		var assignType string
		for k := range assignTypes {
			assignType = k
			break
		}
		if varType != assignType {
			if varType != "unknown" && assignType != "unknown" {
				fileName := graph.fileName
				line := strconv.Itoa(graph.line[node])
				column := strconv.Itoa(graph.column[node])
				logger.Error(fileName + ":" + line + ":" + column + " " + "Type mismatch for variable: " + findAccessName(graph, sorted[0], "") + " is " + varType + " and was assigned to " + assignType)
			}
		}
		varStruct := findStruct(graph, scope, sorted[0], true)
		if varStruct != nil {
			if varStruct.IsLoop {
				fileName := graph.fileName
				line := strconv.Itoa(graph.line[sorted[0]])
				column := strconv.Itoa(graph.column[sorted[0]])
				logger.Error(fileName + ":" + line + ":" + column + " " + "Loop variable " + varStruct.VName + " cannot be assigned")
			}
			if !varStruct.IsParamOut && varStruct.IsParamIn {
				fileName := graph.fileName
				line := strconv.Itoa(graph.line[sorted[0]])
				column := strconv.Itoa(graph.column[sorted[0]])
				logger.Error(fileName + ":" + line + ":" + column + " " + "Variable " + varStruct.VName + " is an in parameter and cannot be assigned")
			}
		}
	case "return":
		scopeSymb := findMotherFunc(scope)
		// return either func or proc symbol
		if _, ok := scopeSymb.(Procedure); ok {
			if len(sorted) != 0 {
				fileName := graph.fileName
				line := strconv.Itoa(graph.line[node])
				column := strconv.Itoa(graph.column[node])
				logger.Error(fileName + ":" + line + ":" + column + " " + "Procedure can't return a value")
			}
		} else {
			if len(sorted) == 0 {
				fileName := graph.fileName
				line := strconv.Itoa(graph.line[node])
				column := strconv.Itoa(graph.column[node])
				logger.Error(fileName + ":" + line + ":" + column + " " + "return can't be standalone in function")
			} else {
				expectedType := make(map[string]struct{})
				expectedType[scopeSymb.(Function).ReturnType] = struct{}{}
				returnType := getReturnType(graph, scope, sorted[0], expectedType)
				//fmt.Println(scopeSymb.(Function).ReturnType, returnType)
				if !haveType(returnType, scopeSymb.(Function).ReturnType) {
					stringTypes := ""
					for k := range returnType {
						stringTypes = stringTypes + ", " + k
					}
					fileName := graph.fileName
					line := strconv.Itoa(graph.line[node])
					column := strconv.Itoa(graph.column[node])
					logger.Error(fileName + ":" + line + ":" + column + " " + "Return types " + stringTypes[2:] + " don't match " + scopeSymb.(Function).FName + " return type " + scopeSymb.(Function).ReturnType)
				}
			}
		}
		updateReturn(graph, node)
	case "call":
		fmt.Printf("call scope: %d %v\n", scope.Region, scope.Table)
		symbolType := getSymbol(graph, scope, sorted[0])
		//fmt.Println("symbolType", symbolType, graph.types[sorted[0]])
		if symbolType == Func {
			fileName := graph.fileName
			line := strconv.Itoa(graph.line[node])
			column := strconv.Itoa(graph.column[node])
			logger.Error(fileName + ":" + line + ":" + column + " " + "Cannot use call to function " + graph.types[sorted[0]] + " as a statement")
		} else if symbolType == Proc {
			//fmt.Println("Proc", graph.types[sorted[0]], maps.Keys(graph.gmap[sorted[1]]))

			matchProc(graph, scope, sorted[0], maps.Keys(graph.gmap[sorted[1]]), genArgsMap(graph, scope, maps.Keys(graph.gmap[sorted[1]])))
		} else if symbolType == Rec {
			fileName := graph.fileName
			line := strconv.Itoa(graph.line[node])
			column := strconv.Itoa(graph.column[node])
			logger.Error(fileName + ":" + line + ":" + column + " " + "Cannot use call to type " + graph.types[sorted[0]] + " as a statement")
		} else if symbolType == Unknown {
			fileName := graph.fileName
			line := strconv.Itoa(graph.line[node])
			column := strconv.Itoa(graph.column[node])
			logger.Error(fileName + ":" + line + ":" + column + " " + "Cannot use call to " + graph.types[sorted[0]] + " as a statement")
		} else {
			fileName := graph.fileName
			line := strconv.Itoa(graph.line[node])
			column := strconv.Itoa(graph.column[node])
			logger.Error(fileName + ":" + line + ":" + column + " " + "Cannot use call to variable " + graph.types[sorted[0]] + " as a statement")
		}
	case "if", "elif":
		if !haveType(getReturnType(graph, scope, sorted[0], make(map[string]struct{})), "boolean") {
			fileName := graph.fileName
			line := strconv.Itoa(graph.line[node])
			column := strconv.Itoa(graph.column[node])
			logger.Error(fileName + ":" + line + ":" + column + " " + "Condition should be boolean")
		}
		for _, child := range sorted[1:] {
			semCheck(graph, child)
		}
	default:
		//is something not accepted
		if len(sorted) == 0 && whichFinal(graph, node) == "identifier" {
			identType := getSymbol(graph, scope, node)
			fileName := graph.fileName
			line := strconv.Itoa(graph.line[node])
			column := strconv.Itoa(graph.column[node])
			if identType == Proc {
				newNode := makeChild2(graph, node, "call", graph.types[node])
				matchProc(graph, scope, newNode, []int{}, genArgsMap(graph, scope, []int{}))
			} else {
				logger.Error(fileName + ":" + line + ":" + column + " " + identType + " " + graph.types[node] + " is not a statement")
			}
		}
		for _, child := range sorted {
			semCheck(graph, child)
		}
	}
}
