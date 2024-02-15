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
	regionCounter *int
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
	Params     *[]*Variable
	ReturnType Type
	children   []int
}

type Procedure struct {
	PName      string
	PType      Type
	ParamCount int
	Params     *[]*Variable
	children   []int
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

func (p Procedure) Name() string {
	return p.PName
}

func (p Procedure) Type() Type {
	return p.PType
}

func (p Procedure) Offset() int {
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
		regionCounter = parent.regionCounter
	}
	return &Scope{parent: parent, Table: make(map[string]Symbol), Children: &[]*Scope{}, regionCounter: regionCounter}
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

			// Add param to parent function or procedure
			for _, symbol := range scope.parent.Table {
				if symbol.Type() == Func && slices.Contains(symbol.(Function).children, child) {
					function := symbol.(Function)
					*function.Params = append(*function.Params, &Variable{VName: graph.types[sorted[0]], SType: Int, IsParam: true})
					function.ParamCount++
				} else if symbol.Type() == Proc && slices.Contains(symbol.(Procedure).children, child) {
					procedure := symbol.(Procedure)
					*procedure.Params = append(*procedure.Params, &Variable{VName: graph.types[sorted[0]], SType: Int, IsParam: true})
					procedure.ParamCount++
				}
			}
		case "function":
			sorted := maps.Keys(graph.gmap[child])
			slices.Sort(sorted)
			scope.addSymbol(Function{FName: graph.types[sorted[0]], SType: Func, children: sorted, Params: &[]*Variable{}})
		}

		if isNodeNewScope(graph.types[child]) {
			scope = *newScope(currentScope)
			scope.Nested = currentScope.Nested + 1
			*scope.regionCounter = *currentScope.regionCounter + 1
			scope.Region = *currentScope.regionCounter
			fmt.Println("new scope", currentScope.Region, len(*currentScope.Children))

			*currentScope.Children = append(*currentScope.Children, &scope)

			fmt.Println("new scope", currentScope.Region, len(*currentScope.Children))
		}

		dfs(graph, child, &scope)
	}
}
