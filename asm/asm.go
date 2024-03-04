package asm

import (
	"github.com/charmbracelet/log"
	"os"
)

type AssemblyFile struct {
	FileName string
	Text     string
}

func NewAssemblyFile(fileName string) AssemblyFile {
	return AssemblyFile{FileName: fileName, Text: ""}
}

func (a AssemblyFile) Name() string {
	return a.FileName
}

func (a AssemblyFile) Content() string {
	return a.Text
}

func (a AssemblyFile) Write() {
	file, err := os.Create(a.FileName + ".s")
	if err != nil {
		log.Fatal("Error while creating file")
	}
	defer file.Close()
	_, err = file.WriteString(a.Text)
	if err != nil {
		log.Fatal("Error while writing to file")
	}
}

func (a AssemblyFile) Execute() []string {
	output := Execute(a.FileName + ".s")
	return output
}
