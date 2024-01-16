package reader

import (
	"fmt"
	"testing"
)

func TestParseCorrectSyntax(t *testing.T) {
	folder := "../examples/correctsyntax"
	for _, file := range ListFiles(folder) {
		fmt.Println("Parsing file", file)
		CompileFile(CompileConfig{Path: file, PythonExecutable: "python"})
	}
}

func TestParseExec(t *testing.T) {
	folder := "../examples/exec"
	for _, file := range ListFiles(folder) {
		fmt.Println("Parsing file", file)
		CompileFile(CompileConfig{Path: file, PythonExecutable: "python"})
	}
}
