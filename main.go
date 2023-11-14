package main

import "gada/reader"

func main() {
	//lexer.AllTest() // set parameter true to print lexical analysis

	reader.CompileFile("examples/expressions/expression.ada")
}
