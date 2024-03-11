package parser

import (
	"fmt"
	"gada/asm"
	"github.com/charmbracelet/log"
	"os"
	"strconv"
	"strings"
)

type AssemblyFile struct {
	FileName string
	Text     string

	ForCounter int
	CurrentFor int
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

const (
	EQ = "EQ"
	NE = "NE"
	CS = "CS"
	CC = "CC"
	MI = "MI"
	PL = "PL"
	VS = "VS"
	VC = "VC"
	HI = "HI"
	LS = "LS"
	GE = "GE"
	LT = "LT"
	GT = "GT"
	LE = "LE"
	AL = "AL"
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

func (a *AssemblyFile) Negate(register Register) {
	a.Text += "; Negate " + register.String() + "\n"
	a.Text += "RSB " + register.String() + ", " + register.String() + ", #0\n\n"
}

func (a *AssemblyFile) Positive(register Register) {
	a.Text += fmt.Sprintf(`; Make %[1]v positive
CMP     %[1]v, #0 ; Compare %[1]v with zero
MOVGE   %[1]v, %[1]v ; If %[1]v is greater than or equal to zero, keep its value (no change)
RSBLT   %[1]v, %[1]v, #0 ; If %[1]v is less than zero, negate it

`, register.String())
}

func (a *AssemblyFile) Cmp(register Register, value int) {
	str := strconv.Itoa(value)
	a.Text += "CMP " + register.String() + ", #" + str + "\n"
}

func (a *AssemblyFile) AddLabel(label string) {
	a.Text += label + "\n"
}

func (a *AssemblyFile) AddComment(comment string) {
	a.Text += "; " + comment + "\n"
}

func (a *AssemblyFile) BranchToLabel(label string) {
	a.Text += "B " + label + "\n"
}

func (a *AssemblyFile) BranchToLabelWithCondition(label string, condition string) {
	a.Text += "B" + condition + " " + label + "\n"
}

func (a AssemblyFile) Write() {
	file, err := os.Create(a.FileName)
	if err != nil {
		log.Fatal("Error while creating file", err)
	}
	defer file.Close()
	_, err = file.WriteString(a.Text)
	if err != nil {
		log.Fatal("Error while writing to file")
	}
}

func (a AssemblyFile) Execute() []string {
	output := asm.Execute(a.FileName)
	return output
}

func ReadASTToASM(graph Graph) {
	log.Info("Reading AST to ASM")
	file := NewAssemblyFile(strings.Replace(graph.fileName, ".ada", ".s", -1))
	file.ReadFile(graph, 0)

	file.Text += "end\n"

	// Multiplication algorithm
	file.Text += `
;       Multiplication algorithm
;       R0 = result, R1 = multiplicand, R2 = multiplier
mul      MOV     R0, #0
mul_loop LSRS    R2, R2, #1
         ADDCS   R0, R0, R1
         LSL     R1, R1, #1
         TST     R2, R2
         BNE     mul_loop
		 LDMFD   SP!, {PC}
`

	// Integer division algorithm
	file.Text += `
;       Integer division routine
;       Arguments:
;       R0 = Dividend
;       R1 = Divisor
;       Returns:
;       R0 = Quotient
;       R1 = Remainder
div32  
       STMFD   SP!, {R2-R4, LR} ; Save registers on the stack
       MOV     R4, #1 ; Bit position = 1
       MOV     R2, #0 ; Quotient = 0
       MOV     R3, R0 ; Remainder = Dividend

loop   
       CMP     R3, R1 ; Compare remainder and divisor
       BCC     shift ; If remainder < divisor, shift
       SUB     R3, R3, R1 ; Remainder = Remainder - Divisor
       ADD     R2, R2, R4 ; Quotient = Quotient + Bit position
       B       loop

shift  
       MOV     R0, R2 ; R0 = Quotient
       MOV     R1, R3 ; R1 = Remainder
       LDMFD   SP!, {R2-R4, PC} ; Restore registers and return
`

	// Fix sign for division
	file.Text += `
fix_sign   
       stmfd   sp!, {PC}
       bl      mul
       subs    r0, r0, #0
       blt     minus_sign
       LDMFD   SP!, {PC}
minus_sign 
       rsb     r3, r3, #0
       LDMFD   SP!, {PC}
`

	log.Info("\n" + file.Text)

	file.Write()
}

func (a *AssemblyFile) CallProcedure(name string) {
	a.Text += "STMFD SP!, {PC}\n"
	a.Text += "BL " + name + "\n"
}

func (a *AssemblyFile) ReadFile(graph Graph, node int) {
	// Read all the children
	children := graph.GetChildren(node)
	for _, child := range children {
		if graph.GetNode(child) == "decl" {
			a.ReadDecl(graph, child)
		} else if graph.GetNode(child) == "body" {
			a.ReadBody(graph, child)
		}
	}
}

func (a *AssemblyFile) ReadBody(graph Graph, node int) {
	// Read all the children
	children := graph.GetChildren(node)
	for _, child := range children {
		switch graph.GetNode(child) {
		case ":=":
			left := graph.GetChildren(child)[0]
			right := graph.GetChildren(child)[1]

			a.ReadOperand(graph, right, 0)

			a.Add(SP, 4)

			// Move the result to the left operand
			scope := graph.getScope(node)
			_, offset := goUpScope(scope, graph.GetNode(left))

			// Save the result in stack
			a.StrWithOffset(R0, offset)
		case "for":
			a.ReadFor(graph, child)
		}
	}
}

func (a *AssemblyFile) ReadFor(graph Graph, node int) {
	goodCounter := a.ForCounter
	a.ForCounter++
	a.CurrentFor++

	// Reserve space for the loop variable and the current for index
	a.AddComment("Loop #" + strconv.Itoa(goodCounter) + " start")
	a.Sub(SP, 4)

	children := graph.GetChildren(node)
	counterStart, _ := strconv.Atoi(graph.GetNode(children[2]))
	counterEnd, _ := strconv.Atoi(graph.GetNode(children[3]))

	a.Mov(R0, counterStart)
	a.StrWithOffset(R0, 4)

	a.AddLabel("for" + strconv.Itoa(goodCounter))

	a.Ldr(R0, 4)
	a.Cmp(R0, counterEnd)
	if graph.GetNode(children[1]) == "not reverse" {
		a.BranchToLabelWithCondition("endfor"+strconv.Itoa(goodCounter), "GT")
	} else {
		a.BranchToLabelWithCondition("endfor"+strconv.Itoa(goodCounter), "LT")
	}

	// Read the body of the for loop
	a.ReadBody(graph, children[4])

	// Increment the counter
	a.Ldr(R0, 4)
	if graph.GetNode(children[1]) == "not reverse" {
		a.Add(R0, 1)
	} else {
		a.Sub(R0, 1)
	}
	a.StrWithOffset(R0, 4)

	// Go to the beginning of the loop
	a.BranchToLabel("for" + strconv.Itoa(goodCounter))

	// End of the for loop
	a.AddLabel("endfor" + strconv.Itoa(goodCounter))

	// Clear the stack
	a.CurrentFor--
	a.Add(SP, 4)

	a.AddComment("Loop #" + strconv.Itoa(goodCounter) + " end")
}

func (a *AssemblyFile) ReadDecl(graph Graph, node int) {
	// Read all the children
	children := graph.GetChildren(node)
	for _, child := range children {
		if graph.GetNode(child) == "var" {
			a.ReadVar(graph, child)
		}
	}
}

func (a *AssemblyFile) ReadVar(graph Graph, node int) {
	var nameList []string
	var value *int

	name := graph.GetChildren(node)[0]

	if graph.GetNode(name) == "sameType" {
		for _, child := range graph.GetChildren(name) {
			nameList = append(nameList, graph.GetNode(child))
		}
	} else {
		nameList = append(nameList, graph.GetNode(name))
	}
	if len(graph.GetChildren(node)) > 2 {
		v, _ := strconv.Atoi(graph.GetNode(graph.GetChildren(node)[2]))
		value = &v
	}

	// use values from the symbol table
	scope := graph.getScope(node)

	for _, name := range nameList {
		for _, symbol := range scope.Table[name] {
			if variable, ok := symbol.(Variable); ok {
				if value != nil {
					offset := getTypeSize(variable.SType, *scope)

					// Move the stack pointer
					a.Sub(SP, offset)

					// Load the int value to r0
					a.Mov(R0, *value)
					a.StrWithOffset(R0, offset)
				} else {
					offset := getTypeSize(variable.SType, *scope)

					// Move the stack pointer
					a.Sub(SP, offset)
				}
			}
		}
	}
}

func (a *AssemblyFile) ReadOperand(graph Graph, node int, fixOffset int) {
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
				scope := graph.getScope(node)
				_, offset := goUpScope(scope, graph.GetNode(node))

				// Load the value from the stack
				a.Ldr(R0, offset+fixOffset)

				// Move the stack pointer
				a.Sub(SP, 4)
				a.StrWithOffset(R0, 4)
			}
		}
	}

	switch graph.GetNode(node) {
	case "+":
		// Read left operand
		a.ReadOperand(graph, children[0], fixOffset+0)

		// Read right operand
		a.ReadOperand(graph, children[1], fixOffset+4)
		a.Ldr(R0, 4)
		a.AddWithOffset(R0, R1, 8)

		// Move the stack pointer
		a.Add(SP, 4)

		// Save the result in stack
		a.StrWithOffset(R0, 4)
	case "-":
		// Read left operand
		a.ReadOperand(graph, children[0], fixOffset+0)

		// Read right operand
		a.ReadOperand(graph, children[1], fixOffset+4)
		a.Ldr(R0, 4)
		a.SubWithOffset(R0, R1, 8)

		// Move the stack pointer
		a.Add(SP, 4)

		// Save the result in stack
		a.StrWithOffset(R0, 4)
	case "*":
		// Read left operand
		a.ReadOperand(graph, children[0], fixOffset+0)

		// Read right operand
		a.ReadOperand(graph, children[1], fixOffset+4)

		// Left operand in R1, right operand in R2
		a.Ldr(R1, 4)
		a.Ldr(R2, 8)

		// Use the multiplication algorithm at the label mul
		a.CallProcedure("mul")

		// Clear the stack (move the stack pointer)
		a.Add(SP, 4)

		// Save the result in stack
		a.StrWithOffset(R0, 4)
	case "/":
		// Read left operand
		a.ReadOperand(graph, children[0], fixOffset+0)

		// Read right operand
		a.ReadOperand(graph, children[1], fixOffset+4)

		// Left operand in R0, right operand in R1
		a.Ldr(R1, 4)
		a.Ldr(R0, 8)

		// Make R0 and R1 positive
		a.Positive(R0)
		a.Positive(R1)

		// Use the division algorithm at the label div32
		a.CallProcedure("div32")

		// Clear the stack (move the stack pointer)
		a.Add(SP, 8)

		// Apply the sign
		// Move left operand in R1, right operand in R2, result in R3
		a.Ldr(R1, 0)
		a.Ldr(R2, 4)
		a.MovRegister(R3, R0)
		a.CallProcedure("fix_sign")

		// Save the result in stack
		a.StrWithOffset(R3, 4)
	case "call":
		if graph.GetNode(children[0]) == "-" {
			// Read right operand
			a.ReadOperand(graph, children[1], fixOffset+0)

			a.Ldr(R0, 4)
			a.Negate(R0)

			a.StrWithOffset(R0, 4)
		}
	}
}
