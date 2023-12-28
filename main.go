package main

import (
	"gada/reader"
	"os"
)

func main() {
	argsWithoutProg := os.Args[1:]
	if len(argsWithoutProg) > 0 {
		reader.CompileFile(reader.CompileConfig{Path: argsWithoutProg[0], PrintAst: containsArgument(argsWithoutProg, "--print-ast")})
		return
	}
	reader.CompileFile(reader.CompileConfig{Path: "examples/expressions/expression.ada", PrintAst: true})
}

func containsArgument(args []string, arg string) bool {
	for _, a := range args {
		if a == arg {
			return true
		}
	}
	return false
}
