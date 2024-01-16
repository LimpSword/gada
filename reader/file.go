package reader

import (
	"fmt"
	"gada/lexer"
	"gada/parser"
	"gada/token"
	"github.com/charmbracelet/log"
	"os"
)

type CompileConfig struct {
	Path             string
	PrintAst         bool
	PythonExecutable string
}

func ReadFile(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println(fmt.Errorf("error while opening file %s: %s", path, err))
		return "", err
	}
	defer file.Close()

	// Check if it is a file
	stat, err := file.Stat()
	if err != nil {
		fmt.Println(fmt.Errorf("error while getting file stats %s: %s", path, err))
		return "", err
	}
	if stat.IsDir() {
		fmt.Println(fmt.Errorf("error: %s is a directory", path))
		return "", err
	}

	// Read file content
	fileContent := make([]byte, stat.Size())
	_, err = file.Read(fileContent)
	if err != nil {
		fmt.Println(fmt.Errorf("error while reading file %s: %s", path, err))
		return "", err
	}
	return string(fileContent), nil
}

func ListFiles(folder string) []string {
	files := []string{}
	file, err := os.Open(folder)
	if err != nil {
		fmt.Println(fmt.Errorf("error while opening folder %s: %s", folder, err))
		return files
	}
	defer file.Close()

	// Check if it is a folder
	stat, err := file.Stat()
	if err != nil {
		fmt.Println(fmt.Errorf("error while getting folder stats %s: %s", folder, err))
		return files
	}
	if !stat.IsDir() {
		fmt.Println(fmt.Errorf("error: %s is not a directory", folder))
		return files
	}

	// Read folder content
	fileInfos, err := file.Readdir(-1)
	if err != nil {
		fmt.Println(fmt.Errorf("error while reading folder %s: %s", folder, err))
		return files
	}
	for _, fileInfo := range fileInfos {
		if !fileInfo.IsDir() {
			files = append(files, folder+"/"+fileInfo.Name())
		}
	}
	return files
}

func FileLexer(path string) *lexer.Lexer {
	content, err := ReadFile(path)
	if err != nil {
		fmt.Println(fmt.Errorf("error while reading file %s: %s", path, err))
		return nil
	}
	return lexer.NewLexer(path, content)
}

func CompileFile(config CompileConfig) {
	l := FileLexer(config.Path)
	if l == nil {
		return
	}
	l.Read()

	if len(l.Tokens) == 0 {
		log.Error("The provided file is empty")
		return
	}

	// remove illegal tokens
	for i := 0; i < len(l.Tokens); i++ {
		if l.Tokens[i].Value == token.ILLEGAL {
			l.Tokens = append(l.Tokens[:i], l.Tokens[i+1:]...)
			i--
		}
	}

	if len(l.Tokens) == 0 {
		log.Error("The provided file has no valid tokens")
		return
	}

	parser.Parse(l, config.PrintAst, config.PythonExecutable)
}
