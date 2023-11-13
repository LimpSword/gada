package reader

import (
	"fmt"
	"gada/lexer"
	"gada/parser"
	"os"
)

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

func FileLexer(path string) *lexer.Lexer {
	content, err := ReadFile(path)
	if err != nil {
		fmt.Println(fmt.Errorf("error while reading file %s: %s", path, err))
		return nil
	}
	return lexer.NewLexer(content)
}

func CompileFile(path string) {
	l := FileLexer(path)
	if l == nil {
		return
	}
	l.Read()
	parser.Parse(l)
}
