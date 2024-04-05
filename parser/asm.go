package parser

import (
	"fmt"
	"gada/asm"
	"github.com/charmbracelet/log"
	"golang.org/x/exp/maps"
	"math/rand"
	"os"
	"slices"
	"strconv"
	"strings"
)

type AssemblyFile struct {
	FileName string
	Text     string
	EndText  string

	WritingAtEnd bool

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
	PC = 15
	LR = 14
)

type Condition string

const (
	EQ = Condition("EQ")
	NE = Condition("NE")
	CS = Condition("CS")
	CC = Condition("CC")
	MI = Condition("MI")
	PL = Condition("PL")
	VS = Condition("VS")
	VC = Condition("VC")
	HI = Condition("HI")
	LS = Condition("LS")
	GE = Condition("GE")
	LT = Condition("LT")
	GT = Condition("GT")
	LE = Condition("LE")
	AL = Condition("AL")
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
	if a.WritingAtEnd {
		a.EndText += "STMFD SP!, {" + register.String() + "}\n"
	} else {
		a.Text += "STMFD SP!, {" + register.String() + "}\n"
	}
}

func (a *AssemblyFile) Ldmfd(register Register) {
	if a.WritingAtEnd {
		a.EndText += "LDMFD SP!, {" + register.String() + "}\n"
	} else {
		a.Text += "LDMFD SP!, {" + register.String() + "}\n"
	}
}

func (a *AssemblyFile) Ldr(register Register, offset int) {
	str := strconv.Itoa(offset)
	if a.WritingAtEnd {
		a.EndText += "LDR " + register.String() + ", [SP, #" + str + "]\n"
	} else {
		a.Text += "LDR " + register.String() + ", [SP, #" + str + "]\n"
	}
}

func (a *AssemblyFile) Str(register Register) {
	if a.WritingAtEnd {
		a.EndText += "STR " + register.String() + ", [SP]\n"
	} else {
		a.Text += "STR " + register.String() + ", [SP]\n"
	}
}

func (a *AssemblyFile) StrWithOffset(register Register, offset int) {
	str := strconv.Itoa(offset)
	if a.WritingAtEnd {
		a.EndText += "STR " + register.String() + ", [SP, #" + str + "]\n"
	} else {
		a.Text += "STR " + register.String() + ", [SP, #" + str + "]\n"
	}
}

func (a *AssemblyFile) Mov(register Register, value int) {
	str := strconv.Itoa(value)
	if a.WritingAtEnd {
		a.EndText += "MOV " + register.String() + ", #" + str + "\n"
	} else {
		a.Text += "MOV " + register.String() + ", #" + str + "\n"
	}
}

func (a *AssemblyFile) MovCond(register Register, value int, condition Condition) {
	str := strconv.Itoa(value)
	if a.WritingAtEnd {
		a.EndText += "MOV" + string(condition) + " " + register.String() + ", #" + str + "\n"
	} else {
		a.Text += "MOV" + string(condition) + " " + register.String() + ", #" + str + "\n"
	}
}

func (a *AssemblyFile) MovToStackPointer(register Register) {
	// use str
	if a.WritingAtEnd {
		a.EndText += "STR " + register.String() + ", [SP]\n"
	} else {
		a.Text += "STR " + register.String() + ", [SP]\n"
	}
}

func (a *AssemblyFile) MovRegister(dest Register, source Register) {
	if a.WritingAtEnd {
		a.EndText += "MOV " + dest.String() + ", " + source.String() + "\n"
	} else {
		a.Text += "MOV " + dest.String() + ", " + source.String() + "\n"
	}
}

func (a *AssemblyFile) And(register1 Register, register2 Register) {
	if a.WritingAtEnd {
		a.EndText += "AND " + R0.String() + ", " + register1.String() + ", " + register2.String() + "\n"
	} else {
		a.Text += "AND " + R0.String() + ", " + register1.String() + ", " + register2.String() + "\n"
	}
}

func (a *AssemblyFile) Or(register1 Register, register2 Register) {
	if a.WritingAtEnd {
		a.EndText += "OR " + register1.String() + ", " + register1.String() + ", " + register2.String() + "\n"
	} else {
		a.Text += "OR " + register1.String() + ", " + register1.String() + ", " + register2.String() + "\n"
	}
}

func (a *AssemblyFile) Add(register Register, value int) {
	str := strconv.Itoa(value)
	if a.WritingAtEnd {
		a.EndText += "ADD " + register.String() + ", " + register.String() + ", #" + str + "\n"
	} else {
		a.Text += "ADD " + register.String() + ", " + register.String() + ", #" + str + "\n"
	}
}

func (a *AssemblyFile) AddFromStackPointer(register Register, intermediateRegister Register) {
	if a.WritingAtEnd {
		a.EndText += "LDR " + intermediateRegister.String() + ", [SP]\n"
		a.EndText += "ADD " + register.String() + ", " + register.String() + ", " + intermediateRegister.String() + "\n"
	} else {
		a.Text += "LDR " + intermediateRegister.String() + ", [SP]\n"
		a.Text += "ADD " + register.String() + ", " + register.String() + ", " + intermediateRegister.String() + "\n"
	}
}

func (a *AssemblyFile) AddWithOffset(register Register, intermediateRegister Register, offset int) {
	if a.WritingAtEnd {
		a.EndText += "LDR " + intermediateRegister.String() + ", [SP, #" + strconv.Itoa(offset) + "]\n"
		a.EndText += "ADD " + register.String() + ", " + intermediateRegister.String() + ", " + register.String() + "\n"
	} else {
		a.Text += "LDR " + intermediateRegister.String() + ", [SP, #" + strconv.Itoa(offset) + "]\n"
		a.Text += "ADD " + register.String() + ", " + intermediateRegister.String() + ", " + register.String() + "\n"
	}
}

func (a *AssemblyFile) Sub(register Register, value int) {
	str := strconv.Itoa(value)
	if a.WritingAtEnd {
		a.EndText += "SUB " + register.String() + ", " + register.String() + ", #" + str + "\n"
	} else {
		a.Text += "SUB " + register.String() + ", " + register.String() + ", #" + str + "\n"
	}
}

func (a *AssemblyFile) SubFromStackPointer(register Register, intermediateRegister Register) {
	if a.WritingAtEnd {
		a.EndText += "LDR " + intermediateRegister.String() + ", [SP]\n"
		a.EndText += "SUB " + register.String() + ", " + intermediateRegister.String() + ", " + register.String() + "\n"
	} else {
		a.Text += "LDR " + intermediateRegister.String() + ", [SP]\n"
		a.Text += "SUB " + register.String() + ", " + intermediateRegister.String() + ", " + register.String() + "\n"
	}
}

func (a *AssemblyFile) SubWithOffset(register Register, intermediateRegister Register, offset int) {
	if a.WritingAtEnd {
		a.EndText += "LDR " + intermediateRegister.String() + ", [SP, #" + strconv.Itoa(offset) + "]\n"
		a.EndText += "SUB " + register.String() + ", " + intermediateRegister.String() + ", " + register.String() + "\n"
	} else {
		a.Text += "LDR " + intermediateRegister.String() + ", [SP, #" + strconv.Itoa(offset) + "]\n"
		a.Text += "SUB " + register.String() + ", " + intermediateRegister.String() + ", " + register.String() + "\n"
	}
}

func (a *AssemblyFile) Negate(register Register) {
	if a.WritingAtEnd {
		a.EndText += "; Negate " + register.String() + "\n"
		a.EndText += "RSB " + register.String() + ", " + register.String() + ", #0\n"
	} else {
		a.Text += "; Negate " + register.String() + "\n"
		a.Text += "RSB " + register.String() + ", " + register.String() + ", #0\n"
	}
}

func (a *AssemblyFile) Positive(register Register) {
	if a.WritingAtEnd {
		a.EndText += fmt.Sprintf(`; Make %[1]v positive
CMP     %[1]v, #0 ; Compare %[1]v with zero
MOVGE   %[1]v, %[1]v ; If %[1]v is greater than or equal to zero, keep its value (no change)
RSBLT   %[1]v, %[1]v, #0 ; If %[1]v is less than zero, negate it

`, register.String())
	} else {
		a.Text += fmt.Sprintf(`; Make %[1]v positive
CMP     %[1]v, #0 ; Compare %[1]v with zero
MOVGE   %[1]v, %[1]v ; If %[1]v is greater than or equal to zero, keep its value (no change)
RSBLT   %[1]v, %[1]v, #0 ; If %[1]v is less than zero, negate it

`, register.String())
	}
}

func (a *AssemblyFile) Cmp(register Register, value int) {
	str := strconv.Itoa(value)
	if a.WritingAtEnd {
		a.EndText += "CMP " + register.String() + ", #" + str + "\n"
	} else {
		a.Text += "CMP " + register.String() + ", #" + str + "\n"
	}
}

func (a *AssemblyFile) CmpRegisters(register1 Register, register2 Register) {
	if a.WritingAtEnd {
		a.EndText += "CMP " + register1.String() + ", " + register2.String() + "\n"
	} else {
		a.Text += "CMP " + register1.String() + ", " + register2.String() + "\n"
	}
}

func (a *AssemblyFile) AddLabel(label string) {
	if a.WritingAtEnd {
		a.EndText += label + "\n"
	} else {
		a.Text += label + "\n"
	}
}

func (a *AssemblyFile) AddComment(comment string) {
	if a.WritingAtEnd {
		a.EndText += "; " + comment + "\n"
	} else {
		a.Text += "; " + comment + "\n"
	}
}

func (a *AssemblyFile) BranchToLabel(label string) {
	if a.WritingAtEnd {
		a.EndText += "B " + label + "\n"
	} else {
		a.Text += "B " + label + "\n"
	}
}

func (a *AssemblyFile) BranchToLabelWithCondition(label string, condition Condition) {
	if a.WritingAtEnd {
		a.EndText += "B" + string(condition) + " " + label + "\n"
	} else {
		a.Text += "B" + string(condition) + " " + label + "\n"
	}
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

	file.Text += "end\n\n"

	file.Text += file.EndText

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
       LDMFD   SP!, {PC} ; Restore registers and return
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
	if a.WritingAtEnd {
		a.EndText += "STMFD SP!, {LR}\n"
		a.EndText += "BL " + name + "\n"
	} else {
		a.Text += "STMFD SP!, {LR}\n"
		a.Text += "BL " + name + "\n"
	}
}

func (a *AssemblyFile) CallProcedureWithParameters(name string, scope *Scope, removedOffset int) {
	if a.WritingAtEnd {
		a.EndText += "STMFD SP!, {LR}\n"
		a.EndText += "BL " + name + "\n"

		// clear the stack
		a.Add(SP, removedOffset)
	} else {
		a.Text += "STMFD SP!, {LR}\n"
		a.Text += "BL " + name + "\n"

		// clear the stack
		a.Add(SP, removedOffset)
	}
}

func (a *AssemblyFile) ReadFile(graph Graph, node int) {
	// Read all the children
	children := graph.GetChildren(node)

	var declNode int
	var bodyNode int

	for _, child := range children {
		if graph.GetNode(child) == "decl" {
			declNode = child
		} else if graph.GetNode(child) == "body" {
			bodyNode = child
		}
	}

	a.ReadDecl(graph, declNode)
	a.ReadBody(graph, bodyNode)
}

func (a *AssemblyFile) ReadIf(graph Graph, node int) {
	a.AddComment("If statement")

	// Read condition
	condition := graph.GetChildren(node)[0]
	a.ReadOperand(graph, condition, 0)

	a.Ldr(R0, 4)
	a.Add(SP, 4)

	randomLabel := rand.Int()

	a.Cmp(R0, 0)
	a.BranchToLabelWithCondition("else"+strconv.Itoa(randomLabel), "EQ")

	// Read body
	a.ReadBody(graph, graph.GetChildren(node)[1])

	a.AddLabel("else" + strconv.Itoa(randomLabel))

	a.AddComment("End of if statement")
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
			endScope, offset := goUpScope(scope, graph.GetNode(left))

			baseOffset := scope.getCurrentOffset()
			var realOffset int
			if scope == endScope {
				realOffset = baseOffset - offset + 4
			} else {
				fmt.Println(":=", baseOffset, offset)
				realOffset = offset
			}

			// Save the result in stack
			//a.StrWithOffset(R0, offset)
			a.StrWithOffset(R0, realOffset)

			//a.Add(SP, 4)
		case "for":
			a.ReadFor(graph, child)
		case "call":
			name := graph.GetChildren(child)[0]
			args := graph.GetChildren(child)[1]

			if graph.GetNode(name) == "Put" {
				a.ReadOperand(graph, args, 0)

				// Move the result to R0
				a.Ldr(R0, 4)
				a.Add(SP, 4)

				a.CallProcedure("put")
				continue
			}

			for _, arg := range graph.GetChildren(args) {
				a.ReadOperand(graph, arg, 0)
			}

			a.CallProcedureWithParameters(graph.GetNode(name), graph.getScope(node), len(graph.GetChildren(args))*4)
		case "if":
			a.ReadIf(graph, child)
		}
	}
}

func (a *AssemblyFile) ReadFor(graph Graph, node int) {
	goodCounter := a.ForCounter
	a.ForCounter++
	a.CurrentFor++

	// Reserve space for the index
	a.Sub(SP, 4)

	a.AddComment("Loop #" + strconv.Itoa(goodCounter) + " start")

	children := graph.GetChildren(node)
	counterStart, err := strconv.Atoi(graph.GetNode(children[2]))
	if err != nil {
		a.ReadOperand(graph, children[2], 0)
		a.Ldr(R0, 4)
		a.Add(SP, 4)
	} else {
		a.Mov(R0, counterStart)
	}

	a.StrWithOffset(R0, 4)

	a.AddLabel("for" + strconv.Itoa(goodCounter))

	a.Ldr(R0, 4)
	counterEnd, err := strconv.Atoi(graph.GetNode(children[3]))
	if err != nil {
		// FIXME: But it should be read as an operand
		scope := graph.getScope(node)
		endScope, offset := goUpScope(scope, graph.GetNode(children[3]))

		baseOffset := scope.getCurrentOffset()
		var realOffset int
		if scope == endScope {
			fmt.Println(graph.GetNode(children[3]), "loop", baseOffset, offset, 0)
			//realOffset = baseOffset - offset + 4 + 4
			realOffset = baseOffset + offset
		} else {
			fmt.Println(graph.GetNode(children[3]), "endscope loop differs", baseOffset, offset, 0)
			realOffset = offset + 4
		}

		// Load the value from the stack
		a.Ldr(R1, realOffset)

		a.CmpRegisters(R0, R1)
	} else {
		a.Cmp(R0, counterEnd)
	}
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

	// Extend sort for var nodes
	slices.SortFunc(children, func(a, b int) int {
		nodeA := graph.GetNode(a)
		nodeB := graph.GetNode(b)

		fmt.Println(nodeA, nodeB)

		if nodeA == "var" && nodeB == "var" {
			sortedA := maps.Keys(graph.gmap[a])
			slices.Sort(sortedA)
			sortedB := maps.Keys(graph.gmap[b])
			slices.Sort(sortedB)

			nameA := graph.types[sortedA[0]]
			nameB := graph.types[sortedB[0]]
			fmt.Println("why:", nameA, nameB)
			return strings.Compare(nameA, nameB)
		}
		return a - b
	})

	for _, child := range children {
		switch graph.GetNode(child) {
		case "var":
			a.ReadVar(graph, child)
		case "procedure":
			a.ReadProcedure(graph, child)
		}
	}
}

func (a *AssemblyFile) ReadProcedure(graph Graph, node int) {
	children := graph.GetChildren(node)
	procedureName := graph.GetNode(children[0])

	var bodyNode int
	var declNode int
	for _, child := range children {
		if graph.GetNode(child) == "decl" {
			declNode = child
		} else if graph.GetNode(child) == "body" {
			bodyNode = child
		}
	}

	a.WritingAtEnd = true
	// Note: single character labels are not allowed
	a.AddComment("Procedure " + procedureName)
	a.AddLabel(procedureName)
	a.Sub(SP, 4)
	// Read the body of the procedure
	a.ReadBody(graph, bodyNode)
	a.Add(SP, 4)
	// Return
	//a.Ldmfd(PC)
	a.MovRegister(PC, LR)
	a.AddComment("End of procedure " + procedureName)

	a.ReadDecl(graph, declNode)
	a.WritingAtEnd = false
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

					// Load the int value to r0
					a.Mov(R0, *value)
					a.Str(R0)

					// Move the stack pointer
					a.Sub(SP, offset)
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
			// The operand is an int
			// Load the int value to r0
			a.Mov(R0, intValue)
			a.Str(R0)

			// Move the stack pointer
			a.Sub(SP, 4)
		} else {
			if graph.GetNode(node)[0] == '\'' {
				// The operand is a char
				// Load the char value to r0
				a.Mov(R0, int(graph.GetNode(node)[1]))
				a.Str(R0)

				// Move the stack pointer
				a.Sub(SP, 4)
			} else {
				// The operand is an ident
				// Load the ident value to r0

				// Get the address of the ident using the symbol table
				scope := graph.getScope(node)
				endScope, offset := goUpScope(scope, graph.GetNode(node))

				baseOffset := scope.getCurrentOffset()
				var realOffset int
				if scope == endScope {
					realOffset = baseOffset + offset - 4
				} else {
					fmt.Println(graph.GetNode(node), "endscope differs", baseOffset, offset, fixOffset)
					//fixOffset = 0
					realOffset = offset
				}

				// Load the value from the stack
				//a.Ldr(R0, offset+fixOffset)
				a.Ldr(R0, realOffset+fixOffset)

				a.Str(R0)

				// Move the stack pointer
				a.Sub(SP, 4)
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
		a.AddWithOffset(R0, R1, 8) // same as ldr from offset 8 then add

		// We can move the SP back
		a.Add(SP, 8)

		// Save the result in stack
		a.Str(R0)

		// Move the stack pointer
		a.Sub(SP, 4)
	case "-":
		// Read left operand
		a.ReadOperand(graph, children[0], fixOffset+0)

		// Read right operand
		a.ReadOperand(graph, children[1], fixOffset+4)
		a.Ldr(R0, 4)
		a.SubWithOffset(R0, R1, 8)

		// We can move the SP back
		a.Add(SP, 8)

		// Save the result in stack
		a.Str(R0)

		// Move the stack pointer
		a.Sub(SP, 4)
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

		// Move the stack pointer
		a.Add(SP, 8)

		// Save the result in stack
		a.Str(R0)

		// Move the stack pointer
		a.Sub(SP, 4)
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

		// Apply the sign
		// Move left operand in R1, right operand in R2, result in R3
		a.Ldr(R1, 4)
		a.Ldr(R2, 8)
		a.MovRegister(R3, R0)
		a.CallProcedure("fix_sign")

		a.MovRegister(R0, R3)

		// Move the stack pointer
		a.Add(SP, 8)

		// Save the result in stack
		a.Str(R0)

		// Move the stack pointer
		a.Sub(SP, 4)
	case "and":
		// Read left operand
		a.ReadOperand(graph, children[0], fixOffset+0)

		// Read right operand
		a.ReadOperand(graph, children[1], fixOffset+4)

		// Left operand in R1, right operand in R2
		a.Ldr(R1, 4)
		a.Ldr(R2, 8)

		// Use the AND operation
		a.And(R1, R2)

		a.Add(SP, 8)

		// Save the result in stack
		a.Str(R0)

		// Move the stack pointer
		a.Sub(SP, 4)
	case ">":
		// Read left operand
		a.ReadOperand(graph, children[0], fixOffset+0)

		// Read right operand
		a.ReadOperand(graph, children[1], fixOffset+4)

		// Left operand in R0, right operand in R1
		a.Ldr(R1, 4)
		a.Ldr(R0, 8)

		// Compare the operands
		a.CmpRegisters(R0, R1)
		a.MovCond(R0, 1, GT)
		a.MovCond(R0, 0, LE)

		// Move the stack pointer
		a.Add(SP, 8)

		// Save the result in stack
		a.Str(R0)

		// Move the stack pointer
		a.Sub(SP, 4)
	case "call":
		if graph.GetNode(children[0]) == "-" {
			// Read right operand
			a.ReadOperand(graph, children[1], fixOffset+0)

			a.Ldr(R0, 0)
			a.Negate(R0)

			a.Str(R0)
		}
	}
}
