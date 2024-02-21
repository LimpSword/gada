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
	Table         map[string]Symbol
	regionCounter *int
}

type Type int

const (
	Int Type = iota
	Char
	Bool
	Float
	Rec
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
	Name   string
	Type   Type
	Fields []Variable
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

func getSymbolType(symbol string) string {
	return strings.ToLower(symbol)
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

func (scope *Scope) String() string {
	return fmt.Sprintf("Region: %d, Nested: %d, Table: %v", scope.Region, scope.Nested, scope.Table)
}

func (scope *Scope) addSymbol(symbol Symbol) {
	scope.Table[symbol.Name()] = symbol
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

var test = false

func dfsSymbols(graph Graph, node int, currentScope *Scope) {
	children := maps.Keys(graph.gmap[node])
	slices.Sort(children)
	for _, child := range children {
		scope := *currentScope

		sorted := maps.Keys(graph.gmap[child])
		slices.Sort(sorted)

		switch graph.types[child] {
		// A node := is not a declaration, it's an assignment
		case "var":
			scope.addSymbol(Variable{VName: graph.types[sorted[0]], SType: getSymbolType(graph.types[sorted[1]])})
		case "param":
			scope.addSymbol(Variable{VName: graph.types[sorted[0]], SType: getSymbolType(graph.types[sorted[1]]), IsParam: true})

			// Add param to parent function or procedure
			for _, symbol := range scope.parent.Table {
				if symbol.Type() == Func && slices.Contains(symbol.(Function).children, child) {
					function := symbol.(Function)
					*function.Params = append(*function.Params, &Variable{VName: graph.types[sorted[0]], SType: getSymbolType(graph.types[sorted[1]]), IsParam: true})
					function.ParamCount++
				} else if symbol.Type() == Proc && slices.Contains(symbol.(Procedure).children, child) {
					procedure := symbol.(Procedure)
					*procedure.Params = append(*procedure.Params, &Variable{VName: graph.types[sorted[0]], SType: getSymbolType(graph.types[sorted[1]]), IsParam: true})
					procedure.ParamCount++
				}
			}
		case "function":
			scope.addSymbol(Function{FName: graph.types[sorted[0]], SType: Func, children: sorted, Params: &[]*Variable{}, ReturnType: getSymbolType(graph.types[sorted[2]])})
		case "procedure":
			scope.addSymbol(Procedure{PName: graph.types[sorted[0]], PType: Proc, children: sorted, Params: &[]*Variable{}})
		}
		graph.scopes[child] = &scope

		if isNodeNewScope(graph.types[child]) {
			if test {
				scope = *newScope(currentScope)
				scope.Nested = currentScope.Nested + 1
				*scope.regionCounter = *currentScope.regionCounter + 1
				scope.Region = *currentScope.regionCounter

				*currentScope.Children = append(*currentScope.Children, &scope)

				fmt.Println("new scope", currentScope.Region, len(*currentScope.Children))
			} else {
				test = true
			}

		}

		dfsSymbols(graph, child, &scope)
	}
}
