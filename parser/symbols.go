package parser

import (
	"fmt"
	"golang.org/x/exp/maps"
	"slices"
)

type Scope struct {
	Region int
	Nested int

	Parent *Scope
	Table  map[string]Symbol
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
	name    string
	sType   Type
	IsParam bool
}

type Function struct {
	name       string
	sType      Type
	paramCount int
	params     []Variable
	returnType Type
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
	return v.name
}

func (v Variable) Type() Type {
	return v.sType
}

func (v Variable) Offset() int {
	return 0
}

func (f Function) Name() string {
	return f.name
}

func (f Function) Type() Type {
	return f.sType
}

func (f Function) Offset() int {
	return 0
}

func getSymbolType(symbol string) Type {
	return Int
}

func newScope(parent *Scope) *Scope {
	return &Scope{Parent: parent, Table: make(map[string]Symbol)}
}

func isNodeNewScope(node string) bool {
	switch node {
	case "function", "procedure", "record", "for", "while", "if", "else", "elif":
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

	fmt.Println(fileScope)
}

func dfs(graph Graph, node int, currentScope *Scope) {
	children := maps.Keys(graph.gmap[node])
	slices.Sort(children)
	for _, child := range children {
		scope := *currentScope
		if isNodeNewScope(graph.types[child]) {
			scope = *newScope(currentScope)
			scope.Nested++
			scope.Region++
			fmt.Println("new scope")
		} else {
			switch graph.types[child] {
			case ":=":
				sorted := maps.Keys(graph.gmap[child])
				slices.Sort(sorted)
				scope.addSymbol(Variable{name: graph.types[sorted[0]], sType: Int})
			case "var":
				sorted := maps.Keys(graph.gmap[child])
				slices.Sort(sorted)
				scope.addSymbol(Variable{name: graph.types[sorted[0]], sType: Int})
			case "param":
				sorted := maps.Keys(graph.gmap[child])
				slices.Sort(sorted)
				scope.addSymbol(Variable{name: graph.types[sorted[0]], sType: Int, IsParam: true})
			case "function":
				sorted := maps.Keys(graph.gmap[child])
				slices.Sort(sorted)
				scope.addSymbol(Function{name: graph.types[sorted[0]], sType: Func})
			}

		}

		dfs(graph, child, &scope)
	}
}
