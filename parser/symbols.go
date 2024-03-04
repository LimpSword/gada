package parser

import (
	"encoding/json"
	"fmt"
	"golang.org/x/exp/maps"
	"slices"
	"strings"
)

type Scope struct {
	Region int
	Nested int

	parent        *Scope
	Children      *[]*Scope
	Table         map[string][]Symbol
	regionCounter *int
}

type Type int

const (
	Int Type = iota
	Char
	Bool
	Float
	Rec     = "rec"
	Func    = "func"
	Proc    = "proc"
	Unknown = "unknown"
)

type Symbol interface {
	// Name returns the name of the symbol
	Name() string
	// Type returns the type of the symbol
	Type() string
	// Offset returns the offset of the symbol
	Offset() int
}

type Variable struct {
	VName   string
	SType   string
	IsParam bool
}

type Function struct {
	FName      string
	SType      string
	ParamCount int
	Params     *[]*Variable
	ReturnType string
	children   []int
}

type Procedure struct {
	PName      string
	PType      string
	ParamCount int
	Params     *[]*Variable
	children   []int
}

type Record struct {
	RName  string
	SType  string
	Fields map[string]string
}

func (v Variable) Name() string {
	return v.VName
}

func (v Variable) Type() string {
	return v.SType
}

func (v Variable) Offset() int {
	return 0
}

func (f Function) Name() string {
	return f.FName
}

func (f Function) Type() string {
	return f.SType
}

func (f Function) Offset() int {
	return 0
}

func (p Procedure) Name() string {
	return p.PName
}

func (p Procedure) Type() string {
	return p.PType
}

func (p Procedure) Offset() int {
	return 0
}

func (r Record) Name() string {
	return r.RName
}

func (r Record) Type() string {
	return r.SType
}

func (r Record) Offset() int {
	return 0
}

func getSymbolType(symbol string) string {
	return strings.ToLower(symbol)
}

func newScope(parent *Scope) *Scope {
	var regionCounter *int
	var scope *Scope
	if parent == nil {
		regionCounter = new(int)
		*regionCounter = 0

		scope = &Scope{parent: nil, Table: make(map[string][]Symbol), Children: &[]*Scope{}, regionCounter: regionCounter, Region: 0, Nested: 0}
	} else {
		*parent.regionCounter++
		scope = &Scope{parent: parent, Table: make(map[string][]Symbol), Children: &[]*Scope{}, regionCounter: parent.regionCounter, Region: *parent.regionCounter, Nested: parent.Nested + 1}
		*parent.Children = append(*parent.Children, scope)
		fmt.Println("new scope", scope.Region, len(*parent.Children))
	}

	return scope
}

func (scope *Scope) String() string {
	return fmt.Sprintf("Region: %d, Nested: %d, Table: %v", scope.Region, scope.Nested, scope.Table)
}

func (scope *Scope) addSymbol(symbol Symbol) {
	name := symbol.Name()
	if existingSymbols, ok := scope.Table[name]; ok {
		// Array already exists, append the symbol to it
		scope.Table[name] = append(existingSymbols, symbol)
	} else {
		// Array doesn't exist, create a new array with the symbol
		scope.Table[name] = []Symbol{symbol}
	}
}

func ReadAST(graph Graph) (*Scope, error) {
	fileScope := newScope(nil)
	currentScope := *fileScope
	fileNodeIndex := 0

	dfsSymbols(graph, fileNodeIndex, &currentScope)

	// fileScope to json
	b, err := json.MarshalIndent(fileScope, "", "  ")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	fmt.Println(string(b))
	return fileScope, nil
}

func addParam(graph Graph, node int, currentFunc *Function) {
	if graph.types[node] == "param" {
		children := maps.Keys(graph.gmap[node])
		slices.Sort(children)
		if graph.types[children[0]] == "sameType" {
			for _, child := range maps.Keys(graph.gmap[children[0]]) {
				*currentFunc.Params = append(*currentFunc.Params, &Variable{VName: graph.types[child], SType: getSymbolType(graph.types[children[1]]), IsParam: true})
				currentFunc.ParamCount++
			}
		} else {
			*currentFunc.Params = append(*currentFunc.Params, &Variable{VName: graph.types[children[0]], SType: getSymbolType(graph.types[children[1]]), IsParam: true})
			currentFunc.ParamCount++
		}
	}
}

func addParamProc(graph Graph, node int, currentProc *Procedure) {
	if graph.types[node] == "param" {
		children := maps.Keys(graph.gmap[node])
		slices.Sort(children)
		if graph.types[children[0]] == "sameType" {
			for _, child := range maps.Keys(graph.gmap[children[0]]) {
				*currentProc.Params = append(*currentProc.Params, &Variable{VName: graph.types[child], SType: getSymbolType(graph.types[children[1]]), IsParam: true})
				currentProc.ParamCount++
			}
		} else {
			*currentProc.Params = append(*currentProc.Params, &Variable{VName: graph.types[children[0]], SType: getSymbolType(graph.types[children[1]]), IsParam: true})
			currentProc.ParamCount++
		}
	}
}

func dfsSymbols(graph Graph, node int, currentScope *Scope) {
	sorted := maps.Keys(graph.gmap[node])
	slices.Sort(sorted)
	scope := *currentScope

	switch graph.types[node] {
	case "file":
		shift := 0
		if graph.types[sorted[1]] == "decl" {
			children := maps.Keys(graph.gmap[sorted[1]])
			for _, child := range children {
				dfsSymbols(graph, child, currentScope)
			}
			shift++
		}
		dfsSymbols(graph, sorted[1+shift], currentScope)
	case "function":
		funcElem := Function{FName: graph.types[sorted[0]], SType: Func, children: sorted, Params: &[]*Variable{}}
		shift := 0
		if graph.types[sorted[1]] == "params" {

			for _, param := range maps.Keys(graph.gmap[sorted[1]]) {
				addParam(graph, param, &funcElem)
			}
			funcElem.ReturnType = getSymbolType(graph.types[sorted[2]])
			scope.addSymbol(funcElem)
			shift = 1
		}
		if graph.types[sorted[2+shift]] == "decl" {
			funcScope := newScope(&scope)
			children := maps.Keys(graph.gmap[sorted[2+shift]])
			for _, child := range children {
				dfsSymbols(graph, child, funcScope)
			}
			shift++
		}
		dfsSymbols(graph, sorted[2+shift], currentScope)
	case "procedure":
		procElem := Procedure{PName: graph.types[sorted[0]], PType: Proc, children: sorted, Params: &[]*Variable{}}
		shift := 0
		if graph.types[sorted[1]] == "params" {
			for _, param := range maps.Keys(graph.gmap[sorted[1]]) {
				addParamProc(graph, param, &procElem)
			}
			scope.addSymbol(procElem)
			shift = 1
		}
		if graph.types[sorted[1+shift]] == "decl" {
			procScope := newScope(&scope)
			children := maps.Keys(graph.gmap[sorted[1+shift]])
			for _, child := range children {
				dfsSymbols(graph, child, procScope)
			}
			shift++
		}
		dfsSymbols(graph, sorted[1+shift], currentScope)
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
		fmt.Println(graph.types[node])
		for _, node := range sorted {
			dfsSymbols(graph, node, currentScope)
		}
	}
}
