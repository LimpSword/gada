package parser

import (
	"gada/asm"
	"github.com/charmbracelet/log"
	"os"
	"strconv"
)

type AssemblyFile struct {
	FileName string
	Text     string
}

type Register int

const (
	R0 Register = iota
	R1
	R2
	R3
	R4
	R5
	R6
	R7
	R8
	R9
	R10
	R11
	R12
	R13
	SP  = 13
	R14 = iota - 1
	R15
)

func (r Register) String() string {
	str := strconv.Itoa(int(r))
	return "R" + str
}

func NewAssemblyFile(fileName string) AssemblyFile {
	assembler := AssemblyFile{FileName: fileName, Text: ""}
	return assembler
}

func (a *AssemblyFile) Name() string {
	return a.FileName
}

func (a *AssemblyFile) Content() string {
	return a.Text
}

func (a *AssemblyFile) Stmfd(register Register) {
	a.Text += "STMFD SP!, {" + register.String() + "}\n"
}

func (a *AssemblyFile) Ldmfd(register Register) {
	a.Text += "LDMFD SP!, {" + register.String() + "}\n"
}

func (a *AssemblyFile) Ldr(register Register, offset int) {
	str := strconv.Itoa(offset)
	a.Text += "LDR " + register.String() + ", [SP, #" + str + "]\n"
}

func (a *AssemblyFile) Str(register Register) {
	a.Text += "STR " + register.String() + ", [SP]\n"
}

func (a *AssemblyFile) StrWithOffset(register Register, offset int) {
	str := strconv.Itoa(offset)
	a.Text += "STR " + register.String() + ", [SP, #" + str + "]\n"
}

func (a *AssemblyFile) Mov(register Register, value int) {
	str := strconv.Itoa(value)
	a.Text += "MOV " + register.String() + ", #" + str + "\n"
}

func (a *AssemblyFile) MovToStackPointer(register Register) {
	// use str
	a.Text += "STR " + register.String() + ", [SP]\n"
}

func (a *AssemblyFile) MovRegister(dest Register, source Register) {
	a.Text += "MOV " + dest.String() + ", " + source.String() + "\n"
}

func (a *AssemblyFile) Add(register Register, value int) {
	str := strconv.Itoa(value)
	a.Text += "ADD " + register.String() + ", " + register.String() + ", #" + str + "\n"
}

func (a *AssemblyFile) AddFromStackPointer(register Register, intermediateRegister Register) {
	a.Text += "LDR " + intermediateRegister.String() + ", [SP]\n"
	a.Text += "ADD " + register.String() + ", " + register.String() + ", " + intermediateRegister.String() + "\n"
}

func (a *AssemblyFile) AddWithOffset(register Register, intermediateRegister Register, offset int) {
	a.Text += "LDR " + intermediateRegister.String() + ", [SP, #" + strconv.Itoa(offset) + "]\n"
	a.Text += "ADD " + register.String() + ", " + intermediateRegister.String() + ", " + register.String() + "\n"
}

func (a *AssemblyFile) Sub(register Register, value int) {
	str := strconv.Itoa(value)
	a.Text += "SUB " + register.String() + ", " + register.String() + ", #" + str + "\n"
}

func (a *AssemblyFile) SubFromStackPointer(register Register, intermediateRegister Register) {
	a.Text += "LDR " + intermediateRegister.String() + ", [SP]\n"
	a.Text += "SUB " + register.String() + ", " + intermediateRegister.String() + ", " + register.String() + "\n"
}

func (a *AssemblyFile) SubWithOffset(register Register, intermediateRegister Register, offset int) {
	a.Text += "LDR " + intermediateRegister.String() + ", [SP, #" + strconv.Itoa(offset) + "]\n"
	a.Text += "SUB " + register.String() + ", " + intermediateRegister.String() + ", " + register.String() + "\n"
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
	output := asm.Execute(a.FileName + ".s")
	return output
}

func ReadASTToASM(graph Graph) {
	log.Info("Reading AST to ASM")
	file := NewAssemblyFile("test")
	file.ReadOperand(graph, 21)

	file.Text += "end\n"

	file.Text += `
mul      MOV     R0, #0
mul_loop LSRS    R2, R2, #1
         ADDCS   R0, R0, R1
         LSL     R1, R1, #1
         TST     R2, R2
         BNE     mul_loop
		 LDMFD   SP!, {PC}
`
	log.Info("\n" + file.Text)
}

func (a *AssemblyFile) CallProcedure(name string) {
	a.Text += "STMFD SP!, {PC}\n"
	a.Text += "BL " + name + "\n"
}

func (a *AssemblyFile) ReadOperand(graph Graph, node int) {
	// Read left and right operands and do the operation
	// If the operands are values, use them
	// Else, save them in stack and use them
	children := graph.GetChildren(node)
	// sort

	if len(children) == 0 {
		// The operand is a value
		// Different cases: int, char or ident
		intValue, err := strconv.Atoi(graph.GetNode(node))
		if err == nil {
			// Move the stack pointer
			a.Sub(SP, 4)

			// The operand is an int
			// Load the int value to r0
			a.Mov(R0, intValue)
			a.StrWithOffset(R0, 4)
		} else {
			if graph.GetNode(node)[0] == '\'' {
				// Move the stack pointer
				a.Sub(SP, 4)

				// The operand is a char
				// Load the char value to r0
				a.Mov(R0, int(graph.GetNode(node)[1]))
				a.StrWithOffset(R0, 4)
			} else {
				// The operand is an ident
				// Load the ident value to r0

				// Get the address of the ident using the symbol table
			}
		}
	}

	switch graph.GetNode(node) {
	case "+":
		// Read left operand
		a.ReadOperand(graph, children[0])

		// Read right operand
		a.ReadOperand(graph, children[1])
		a.Ldr(R0, 4)
		a.AddWithOffset(R0, R1, 8)

		// Move the stack pointer
		a.Add(SP, 4)

		// Save the result in stack
		a.StrWithOffset(R0, 4)
	case "-":
		// Read left operand
		a.ReadOperand(graph, children[0])

		// Read right operand
		a.ReadOperand(graph, children[1])
		a.Ldr(R0, 4)
		a.SubWithOffset(R0, R1, 8)

		// Move the stack pointer
		a.Add(SP, 4)

		// Save the result in stack
		a.StrWithOffset(R0, 4)
	case "*":
		// Read left operand
		a.ReadOperand(graph, children[0])

		// Read right operand
		a.ReadOperand(graph, children[1])

		// Left operand in R1, right operand in R2
		a.Ldr(R1, 4)
		a.Ldr(R2, 8)

		// Use the multiplication algorithm at the label mul
		a.CallProcedure("mul")

		// Clear the stack (move the stack pointer)
		a.Add(SP, 4)

		// Save the result in stack
		a.StrWithOffset(R0, 4)
	}
}
