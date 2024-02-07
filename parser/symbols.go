package parser

import (
	"encoding/json"
	"fmt"
	"golang.org/x/exp/maps"
	"slices"
)

type Scope struct {
	Region int
	Nested int

	parent        *Scope
	Children      *[]*Scope
	Table         map[string]Symbol
	RegionCounter *int
}

type Type int

const (
	Int Type = iota
	Char
	Bool
	Float
	Func
	Proc
	Rec
)

type Symbol interface {
	// Name returns the name of the symbol
	Name() string
	// Type returns the type of the symbol
	Type() Type
	// Offset returns the offset of the symbol
	Offset() int
}

type Variable struct {
	VName   string
	SType   Type
	IsParam bool
}

type Function struct {
	FName      string
	SType      Type
	ParamCount int
	Params     []Variable
	ReturnType Type
}

type Procedure struct {
	Name       string
	Type       Type
	ParamCount int
	Params     []Variable
}

type Record struct {
	Name   string
	Type   Type
	Fields []Variable
}

func (v Variable) Name() string {
	return v.VName
}

func (v Variable) Type() Type {
	return v.SType
}

func (v Variable) Offset() int {
	return 0
}

func (f Function) Name() string {
	return f.FName
}

func (f Function) Type() Type {
	return f.SType
}

func (f Function) Offset() int {
	return 0
}

func getSymbolType(symbol string) Type {
	return Int
}

func newScope(parent *Scope) *Scope {
	var regionCounter *int
	if parent == nil {
		regionCounter = new(int)
		*regionCounter = 0
	} else {
		regionCounter = parent.RegionCounter
	}
	return &Scope{parent: parent, Table: make(map[string]Symbol), Children: &[]*Scope{}, RegionCounter: regionCounter}
}

func isNodeNewScope(node string) bool {
	switch node {
	case "function", "procedure", "record", "for", "while", "if", "else", "elif", "decl":
		return true
	}
	return false
}

func (scope *Scope) addSymbol(symbol Symbol) {
	scope.Table[symbol.Name()] = symbol
}

func ReadAST(graph Graph) {
	fileScope := newScope(nil)
	currentScope := *fileScope
	fileNodeIndex := 0

	dfs(graph, fileNodeIndex, &currentScope)

	// fileScope to json
	b, err := json.MarshalIndent(fileScope, "", "  ")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(b))
}

func dfs(graph Graph, node int, currentScope *Scope) {
	children := maps.Keys(graph.gmap[node])
	slices.Sort(children)
	for _, child := range children {
		scope := *currentScope

		switch graph.types[child] {
		case ":=":
			sorted := maps.Keys(graph.gmap[child])
			slices.Sort(sorted)
			scope.addSymbol(Variable{VName: graph.types[sorted[0]], SType: Int})
		case "var":
			sorted := maps.Keys(graph.gmap[child])
			slices.Sort(sorted)
			scope.addSymbol(Variable{VName: graph.types[sorted[0]], SType: Int})
		case "param":
			sorted := maps.Keys(graph.gmap[child])
			slices.Sort(sorted)
			scope.addSymbol(Variable{VName: graph.types[sorted[0]], SType: Int, IsParam: true})

			// TODO: add param to parent function or procedure
		case "function":
			sorted := maps.Keys(graph.gmap[child])
			slices.Sort(sorted)
			scope.addSymbol(Function{FName: graph.types[sorted[0]], SType: Func})
		}

		if isNodeNewScope(graph.types[child]) {
			scope = *newScope(currentScope)
			scope.Nested = currentScope.Nested + 1
			*scope.RegionCounter = *currentScope.RegionCounter + 1
			scope.Region = *currentScope.RegionCounter
			fmt.Println("new scope", currentScope.Region, len(*currentScope.Children))

			*currentScope.Children = append(*currentScope.Children, &scope)

			fmt.Println("new scope", currentScope.Region, len(*currentScope.Children))
		}

		dfs(graph, child, &scope)
	}
}
