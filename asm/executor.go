package asm

import (
	"fmt"
	"github.com/charmbracelet/log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	start = "---- PROGRAM OUTPUT ----"
	end   = "---- END PROGRAM OUTPUT ----"
)

func Execute(relativePath string) []string {
	file, err := os.Open(relativePath)
	if err != nil {
		log.Fatal("Error while opening file")
	}
	defer file.Close()
	absolutePath, err := filepath.Abs(relativePath)
	cmd := exec.Command("java", "-jar", "pcl.jar", absolutePath)

	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal("Error while executing file", err)
	}
	text := string(out)
	lines := strings.Split(text, "\n")

	programOutput := make([]string, 0)

	prnt := false
	for _, line := range lines {
		if strings.Contains(line, end) {
			prnt = false
		}
		if prnt {
			programOutput = append(programOutput, line)
		}
		if strings.Contains(line, start) {
			prnt = true
		}
	}

	fmt.Println(string(out))
	return programOutput
}
