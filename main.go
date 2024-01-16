package main

import (
	"gada/reader"
	"os"
	"strings"
)

func main() {

	argsWithoutProg := os.Args[1:]
	if len(argsWithoutProg) > 0 {
		compileConfig := reader.CompileConfig{Path: getProgramName()}

		// Arguments
		printAst, _ := containsArgument(argsWithoutProg, "--print-ast")
		compileConfig.PrintAst = printAst
		pythonExecutable, pythonExecutableValue := containsArgument(argsWithoutProg, "--python-executable")
		if pythonExecutable {
			compileConfig.PythonExecutable = pythonExecutableValue
		} else {
			compileConfig.PythonExecutable = "python3"
		}

		reader.CompileFile(compileConfig)
		return
	}
	reader.CompileFile(reader.CompileConfig{Path: "examples/expressions/helloWorld.ada", PrintAst: true})
}

func getProgramName() string {
	var programName string
	for _, arg := range os.Args[1:] {
		if !strings.HasPrefix(arg, "--") {
			programName = arg
			break
		}
	}
	return programName
}

// containsArgument checks if the given argument is present in the list of arguments and returns the value of the argument if present, "" otherwise
func containsArgument(args []string, arg string) (bool, string) {
	for _, a := range args {
		if strings.HasPrefix(a, arg) && (len(a) == len(arg) || a[len(arg)] == '=') {
			if len(a) == len(arg) {
				return true, ""
			}
			argValue := a[len(arg)+1:]
			return true, argValue
		}
	}
	return false, ""
}
