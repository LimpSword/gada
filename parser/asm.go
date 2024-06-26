package parser

import (
	"fmt"
	"gada/asm"
	"github.com/charmbracelet/log"
	"golang.org/x/exp/maps"
	"math/rand"
	"os"
	"runtime"
	"slices"
	"strconv"
	"strings"
)

type AssemblyFile struct {
	FileName string
	Text     string
	EndText  string

	WritingAtEnd bool

	ForCounter  int
	CurrentAddr int
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

func (a *AssemblyFile) Fill(n int) string {
	addr := "addr" + strconv.Itoa(a.CurrentAddr)
	if a.WritingAtEnd {
		a.EndText += addr + " FILL " + strconv.Itoa(n) + "\n"
	} else {
		a.Text += addr + " FILL " + strconv.Itoa(n) + "\n"
	}
	a.CurrentAddr++
	return addr
}

func (a *AssemblyFile) Stmfd(register Register) {
	if a.WritingAtEnd {
		a.EndText += "STMFD SP!, {" + register.String() + "}\n"
	} else {
		a.Text += "STMFD SP!, {" + register.String() + "}\n"
	}
}

func (a *AssemblyFile) StmfdMultiple(registers []Register) {
	if a.WritingAtEnd {
		a.EndText += "STMFD SP!, {"
		for i, register := range registers {
			if i == len(registers)-1 {
				a.EndText += register.String() + "}\n"
			} else {
				a.EndText += register.String() + ", "
			}
		}
	} else {
		a.Text += "STMFD SP!, {"
		for i, register := range registers {
			if i == len(registers)-1 {
				a.Text += register.String() + "}\n"
			} else {
				a.Text += register.String() + ", "
			}
		}
	}
}

func (a *AssemblyFile) Ldmfd(register Register) {
	if a.WritingAtEnd {
		a.EndText += "LDMFD SP!, {" + register.String() + "}\n"
	} else {
		a.Text += "LDMFD SP!, {" + register.String() + "}\n"
	}
}

func (a *AssemblyFile) LdmfdMultiple(registers []Register) {
	if a.WritingAtEnd {
		a.EndText += "LDMFD SP!, {"
		for i, register := range registers {
			if i == len(registers)-1 {
				a.EndText += register.String() + "}\n"
			} else {
				a.EndText += register.String() + ", "
			}
		}
	} else {
		a.Text += "LDMFD SP!, {"
		for i, register := range registers {
			if i == len(registers)-1 {
				a.Text += register.String() + "}\n"
			} else {
				a.Text += register.String() + ", "
			}
		}
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

func (a *AssemblyFile) LdrAddr(register Register, addr string) {
	if a.WritingAtEnd {
		a.EndText += "LDR " + register.String() + ", =" + addr + "\n"
	} else {
		a.Text += "LDR " + register.String() + ", =" + addr + "\n"
	}
}

func (a *AssemblyFile) LdrFrom(register Register, fromRegister Register, offset int) {
	str := strconv.Itoa(offset)
	if a.WritingAtEnd {
		a.EndText += "LDR " + register.String() + ", [" + fromRegister.String() + ", #" + str + "]\n"
	} else {
		a.Text += "LDR " + register.String() + ", [" + fromRegister.String() + ", #" + str + "]\n"
	}
}

func (a *AssemblyFile) LdrFromFramePointer(register Register, offset int) {
	str := strconv.Itoa(offset)
	if a.WritingAtEnd {
		a.EndText += "LDR " + register.String() + ", [R11, #" + str + "]\n"
	} else {
		a.Text += "LDR " + register.String() + ", [R11, #" + str + "]\n"
	}
}

func (a *AssemblyFile) Str(register Register) {
	if a.WritingAtEnd {
		a.EndText += "STR " + register.String() + ", [SP]\n"
	} else {
		a.Text += "STR " + register.String() + ", [SP]\n"
	}
}

func (a *AssemblyFile) StrFrom(register Register, fromRegister Register, offset int) {
	str := strconv.Itoa(offset)
	if a.WritingAtEnd {
		a.EndText += "STR " + register.String() + ", [" + fromRegister.String() + ", #" + str + "]\n"
	} else {
		a.Text += "STR " + register.String() + ", [" + fromRegister.String() + ", #" + str + "]\n"
	}
}

func (a *AssemblyFile) StrFromFramePointer(register Register, offset int) {
	str := strconv.Itoa(offset)
	if a.WritingAtEnd {
		a.EndText += "STR " + register.String() + ", [R11, #" + str + "]\n"
	} else {
		a.Text += "STR " + register.String() + ", [R11, #" + str + "]\n"
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
		a.EndText += "LDR " + register.String() + ", =" + str + "\n"
	} else {
		a.Text += "LDR " + register.String() + ", =" + str + "\n"
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

func (a *AssemblyFile) Not(register Register) {
	if a.WritingAtEnd {
		a.EndText += "; Not " + register.String() + "\n"
		a.EndText += "MVN " + register.String() + ", " + register.String() + "\n"
	} else {
		a.Text += "; Not " + register.String() + "\n"
		a.Text += "MVN " + register.String() + ", " + register.String() + "\n"
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

func (a *AssemblyFile) CommentPreviousLine(comment string) {
	if a.WritingAtEnd {
		a.EndText = a.EndText[:len(a.EndText)-1] + " ; " + comment + "\n"
	} else {
		a.Text = a.Text[:len(a.Text)-1] + " ; " + comment + "\n"
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

func getStringFromRight(s string) string {
	if runtime.GOOS == "windows" {
		index := strings.LastIndex(s, "\\")
		if index == -1 {
			return s // No '/' found, return the original string
		}
		return s[index+1:] // Return the substring from index+1 to the end
	}
	index := strings.LastIndex(s, "/")
	if index == -1 {
		return s // No '/' found, return the original string
	}
	return s[index+1:] // Return the substring from index+1 to the end
}

func changeOrAddExtension(s string) string {
	// Find the last occurrence of '.'
	index := strings.LastIndex(s, ".")
	if index == -1 {
		// If '.' doesn't exist, just append ".s"
		return s + ".s"
	}
	// Replace the substring from index to end with ".s"
	return s[:index] + ".s"
}

func ReadASTToASM(graph Graph) {
	file := NewAssemblyFile(changeOrAddExtension(fmt.Sprintf("examples/asm/%s", getStringFromRight(graph.fileName))))

	file.Text += "STR_OUT      FILL    0x1000\n"
	file.Text += "MOV R11, SP\n"

	file.ReadFile(graph, 0)

	file.Text += "end\n\n"

	file.Text += file.EndText

	// Multiplication algorithm
	file.Text += `
;       Multiplication algorithm
;       R0 = result, R1 = multiplicand, R2 = multiplier
mul      STMFD   SP!, {LR}
         MOV     R0, #0
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
;       R1 = Dividend
;       R2 = Divisor
;       Returns:
;       R0 = Quotient
;       R1 = Remainder
div32    STMFD   SP!, {LR, R2-R5}
         MOV     R0, #0
         MOV     R3, #0
         CMP     R1, #0
         RSBLT   R1, R1, #0
         EORLT   R3, R3, #1
         CMP     R2, #0
         RSBLT   R2, R2, #0
         EORLT   R3, R3, #1
         MOV     R4, R2
         MOV     R5, #1
div_max  LSL     R4, R4, #1
         LSL     R5, R5, #1
         CMP     R4, R1
         BLE     div_max
div_loop LSR     R4, R4, #1
         LSR     R5, R5, #1
         CMP     R4,R1
         BGT     div_loop
         ADD     R0, R0, R5
         SUB     R1, R1, R4
         CMP     R1, R2
         BGE     div_loop
         CMP     R3, #1
         BNE     div_exit
         CMP     R1, #0
         ADDNE   R0, R0, #1
         RSB     R0, R0, #0
         RSB     R1, R1, #0
         ADDNE   R1, R1, R2
div_exit CMP     R0, #0
         ADDEQ   R1, R1, R4
         LDMFD   SP!, {PC, R2-R5}
`

	// Fix sign for division
	file.Text += `
fix_sign   
       stmfd   sp!, {LR}
       bl      mul
       subs    r0, r0, #0
       blt     minus_sign
       LDMFD   SP!, {PC}
minus_sign 
       rsb     r3, r3, #0
       LDMFD   SP!, {PC}
`

	file.Text += `
println      STMFD   SP!, {LR, R0-R3}
             MOV     R3, R0
             LDR     R1, =STR_OUT ; address of the output buffer
PRINTLN_LOOP LDRB    R2, [R0], #1
             STRB    R2, [R1], #1
             TST     R2, R2
             BNE     PRINTLN_LOOP
             MOV     R2, #0
             STRB    R2, [R1, #-1]
             MOV     R2, #0
             STRB    R2, [R1]


             ;       we need to clear the output buffer
             LDR     R1, =STR_OUT
             MOV     R0, R3
CLEAN        LDRB    R2, [R0], #1
             MOV     R3, #0
             STRB    R3, [R1], #1
             TST     R2, R2
             BNE     CLEAN
             ;       clear 3 more because why not
             STRB    R3, [R1], #1
             STRB    R3, [R1], #1
             STRB    R3, [R1], #1

             LDMFD   SP!, {PC, R0-R3}
`

	file.Text += `to_ascii      STMFD   SP!, {LR, R4-R7}
; make it positive
MOV R7, R0
					 CMP     R0, #0	
					 MOVGE   R6, R0	
					 RSBLT   R6, R0, #0
					 MOV     R0, R6


              MOV     R4, #0 ; Initialize digit counter
              MOV     R5, #10 ; Radix for decimal

to_ascii_loop MOV     R1, R0 ; Save the value in R6
              MOV     R2, #10
              BL      div32 ; R0 = R0 / 10, R1 = R0 % 10
              ADD     R1, R1, #48 ; Convert digit to ASCII
              STRB    R1, [R3, R4] ; Store the ASCII digit
              ADD     R4, R4, #1 ; Increment digit counter
              CMP     R0, #0
              BNE     to_ascii_loop

; add the sign if it was negative
						  CMP     R7, #0
						  MOVGE   R1, #0	
						  MOVLT   R1, #45
						  STRB    R1, [R3, R4]
						  ADD     R4, R4, #1


              LDMFD   SP!, {PC, R4-R7}
`

	file.Write()
}

func (a *AssemblyFile) CallProcedure(name string) {
	if a.WritingAtEnd {
		a.EndText += "BL " + name + "\n"
	} else {
		a.Text += "BL " + name + "\n"
	}
}

func (a *AssemblyFile) CallWithParameters(blName string, scope *Scope, removedOffset int) {
	symbol := scope.ScopeSymbol
	_, isFunction := symbol.(Function)
	if a.WritingAtEnd {
		a.EndText += "BL " + blName + "\n"

		if isFunction {
			a.Add(SP, removedOffset)
			a.Ldr(R0, 0)
			//a.Add(SP, 4)
			return
		}

		// clear the stack
		a.Add(SP, removedOffset)
		a.CommentPreviousLine("Clear the stack of parameters")
	} else {
		a.Text += "BL " + blName + "\n"

		if isFunction {
			a.Add(SP, removedOffset)
			a.Ldr(R0, 0)
			//a.Add(SP, 4)
			return
		}

		// clear the stack
		a.Add(SP, removedOffset)
		a.CommentPreviousLine("Clear the stack of parameters")
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

	a.StmfdMultiple([]Register{R10, R11, LR})
	a.Mov(R10, getRegion(graph, node))
	a.MovRegister(R11, SP)
	a.Sub(R11, 4)

	a.ReadDecl(graph, declNode, All)
	a.ReadBody(graph, bodyNode)
}

func (a *AssemblyFile) ReadIf(graph Graph, node int) {
	a.AddComment("If statement")
	a.AddComment("Start of condition")

	// Read condition
	condition := graph.GetChildren(node)[0]
	a.ReadOperand(graph, condition)

	a.Ldr(R0, 0)
	a.CommentPreviousLine("Load result of condition")
	a.Add(SP, 4)
	a.CommentPreviousLine("Remove the result of the condition from the stack")

	randomLabel := strconv.Itoa(rand.Int())

	a.Cmp(R0, 0)
	a.AddComment("End of condition")
	a.BranchToLabelWithCondition("else"+randomLabel, "EQ")

	// Read body
	a.ReadBody(graph, graph.GetChildren(node)[1])

	a.BranchToLabel("end_if_" + randomLabel)

	if len(graph.GetChildren(node)) >= 3 {
		// Read elsif or else
		if graph.GetNode(graph.GetChildren(node)[2]) == "elif" {
			a.AddComment("Elif statement")
			a.AddLabel("else" + randomLabel)

			// Elif statement
			a.ReadIf(graph, graph.GetChildren(node)[2])

			a.AddComment("End of elif statement")

			if len(graph.GetChildren(node)) >= 4 {
				// Else statement
				a.AddComment("Else statement")
				a.ReadBody(graph, graph.GetChildren(node)[3])
			}
		} else {
			// Else statement
			a.AddLabel("else" + randomLabel)
			a.ReadBody(graph, graph.GetChildren(node)[2])
		}
	} else {
		a.AddLabel("else" + randomLabel)
	}
	a.AddLabel("end_if_" + randomLabel)

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

			a.AddComment("Assignment of " + graph.GetNode(left))
			a.ReadOperand(graph, right)

			// Get the address of the ident using the symbol table
			scope := graph.getScope(node)
			endScope, offset := goUpScope(graph, scope, left, graph.GetNode(left))

			a.MovRegister(R9, R11)
			if scope == endScope {
				a.AddComment(fmt.Sprintf("(S) Store the value of %v", graph.GetNode(left)))
			} else {
				// Loop through dynamic links until we reach the correct region

				random := strconv.Itoa(rand.Int()) + "_" + graph.GetNode(node)

				a.LdrFromFramePointer(R8, 4)
				a.Cmp(R8, endScope.Region)
				a.BranchToLabelWithCondition("notload_"+random, EQ)
				a.AddLabel("load_" + random)
				a.LdrFromFramePointer(R11, 8)
				a.LdrFromFramePointer(R8, 4)
				a.Cmp(R8, endScope.Region)
				a.BranchToLabelWithCondition("load_"+random, NE)

				a.AddLabel("notload_" + random)

				// Go 1 level up
				a.LdrFromFramePointer(R11, 8)

			}
			name := graph.GetNode(left)
			isAccess := false
			if name == "access" {
				// If this is an access we are looking to the right most
				r := graph.GetChildren(left)[1]
				for graph.GetNode(r) == "access" {
					r = graph.GetChildren(r)[1]
				}
				name = graph.GetNode(r)
				isAccess = true
			}
			// What's the size of the type?
			var typeSize int
			if isAccess {
				// TODO: Size of the field
				typeSize = 4
			} else {
				typeSize = getTypeSize(endScope.Table[name][0].Type(), *scope)
			}

			sOffset := 0
			for typeSize > 0 {
				a.Ldr(R0, sOffset)
				a.StrFromFramePointer(R0, offset+sOffset)

				sOffset += 4
				typeSize -= 4
			}

			a.Add(SP, 4)

			// Restore R9
			a.MovRegister(R11, R9)
		case "for":
			a.ReadFor(graph, child)
		case "while":
			a.ReadWhile(graph, child)
		case "call":
			name := graph.GetChildren(child)[0]
			var args int
			if len(graph.GetChildren(child)) > 1 {
				args = graph.GetChildren(child)[1]
			} else {
				args = -1
			}

			a.Call(child, graph, name, args)
		case "return":
			a.AddComment("Return statement")
			if len(graph.GetChildren(child)) == 0 {
				// Leave the procedure
				// Return the in out parameters
				symbol := graph.fullSymbols[node]
				_, isFunction := symbol.(Function)
				_, isProcedure := symbol.(Procedure)
				if isFunction {
					paramOffset := 0
					for _, param := range symbol.(Function).Params {
						if param.Offset > paramOffset {
							paramOffset = param.Offset
						}
					}
					for i := 1; i < len(symbol.(Function).Params)+1; i++ {
						if symbol.(Function).Params[i].IsParamIn && symbol.(Function).Params[i].IsParamOut {
							// Load the address of the in out parameter
							a.LdrFromFramePointer(R0, paramOffset-symbol.(Function).Params[i].Offset+16)
							a.CommentPreviousLine("Load the address of the in out parameter")
							a.LdrFromFramePointer(R1, paramOffset-symbol.(Function).Params[i].Offset+4+16)
							a.CommentPreviousLine("Load the value of the in out parameter")
							a.StrFrom(R1, R0, 0)
							a.CommentPreviousLine("Store the value of the in out parameter")
						}
					}
				} else if isProcedure {
					paramOffset := 0
					for _, param := range symbol.(Procedure).Params {
						if param.Offset > paramOffset {
							paramOffset = param.Offset
						}
					}
					for i := 1; i < len(symbol.(Procedure).Params)+1; i++ {
						if symbol.(Procedure).Params[i].IsParamIn && symbol.(Procedure).Params[i].IsParamOut {
							// Load the address of the in out parameter
							a.LdrFromFramePointer(R0, paramOffset-symbol.(Procedure).Params[i].Offset+16)
							a.CommentPreviousLine("Load the address of the in out parameter")
							a.LdrFromFramePointer(R1, paramOffset-symbol.(Procedure).Params[i].Offset+4+16)
							a.CommentPreviousLine("Load the value of the in out parameter")
							a.StrFrom(R1, R0, 0)
							a.CommentPreviousLine("Store the value of the in out parameter")
						}
					}
				}

				// Need to quit every loop we are in

				a.Add(SP, getDeclOffset(graph, node))
				a.LdmfdMultiple([]Register{R10, R11, PC})
				a.CommentPreviousLine("Return from the procedure")
				return
			}

			// Read the return operand
			operand := graph.GetChildren(child)[0]
			a.ReadOperand(graph, operand)

			// Move the result to R0
			a.Ldr(R0, 0)

			// Save the result at the right place
			// We have to jump the parameters
			scope := graph.getScope(node)
			fnc, isFunction := scope.ScopeSymbol.(Function)
			paramOffset := 0
			if isFunction {
				for _, param := range fnc.Params {
					size := getTypeSize(param.Type(), *scope)
					paramOffset += size
					if param.IsParamIn && param.IsParamOut {
						paramOffset += 4
					}
				}
			}

			symbol := graph.fullSymbols[node]
			if symbol == nil {
				symbol = scope.ScopeSymbol
			}

			/*returnedSize := getTypeSize(fnc.ReturnType, *scope)
			offset := 0
			for returnedSize > 0 {
				a.Ldr(R0, offset)
				a.StrFromFramePointer(R0, 16+paramOffset)
				paramOffset += 4
				returnedSize -= 4
				offset += 4
			}*/

			returnedSize := getTypeSize(fnc.ReturnType, *scope)
			offset := returnedSize - 4
			offset2 := 0
			paramOffset += returnedSize - 4
			for returnedSize > 0 {
				a.Ldr(R0, offset2)
				a.StrFromFramePointer(R0, 16+paramOffset)
				paramOffset -= 4
				returnedSize -= 4
				offset -= 4
				offset2 += 4
			}

			// Leave the procedure
			a.Add(SP, getTypeSize(scope.ScopeSymbol.(Function).ReturnType, *scope))
			a.CommentPreviousLine("Remove the return value from the stack")

			// Return the in out parameters
			_, isProcedure := symbol.(Procedure)
			if isFunction {
				for i := 1; i < len(symbol.(Function).Params)+1; i++ {
					if symbol.(Function).Params[i].IsParamIn && symbol.(Function).Params[i].IsParamOut {
						// Load the address of the in out parameter
						a.LdrFromFramePointer(R0, paramOffset-symbol.(Function).Params[i].Offset+16)
						a.CommentPreviousLine("Load the address of the in out parameter")
						a.LdrFromFramePointer(R1, paramOffset-symbol.(Function).Params[i].Offset+4+16)
						a.CommentPreviousLine("Load the value of the in out parameter")
						a.StrFrom(R1, R0, 0)
						a.CommentPreviousLine("Store the value of the in out parameter")
					}
				}
			} else if isProcedure {
				paramOffset := 0
				for _, param := range symbol.(Procedure).Params {
					if param.Offset > paramOffset {
						paramOffset = param.Offset
					}
				}
				for i := 1; i < len(symbol.(Procedure).Params); i++ {
					if symbol.(Procedure).Params[i].IsParamIn && symbol.(Procedure).Params[i].IsParamOut {
						// Load the address of the in out parameter
						a.LdrFromFramePointer(R0, paramOffset-symbol.(Procedure).Params[i].Offset+16)
						a.CommentPreviousLine("Load the address of the in out parameter")
						a.LdrFromFramePointer(R1, paramOffset-symbol.(Procedure).Params[i].Offset+4+16)
						a.CommentPreviousLine("Load the value of the in out parameter")
						a.StrFrom(R1, R0, 0)
						a.CommentPreviousLine("Store the value of the in out parameter")
					}
				}
			}

			fmt.Println("offset", getDeclOffset(graph, node))
			a.Add(SP, getDeclOffset(graph, node))
			a.CommentPreviousLine("Clear the stack of declarations: " + strconv.Itoa(getDeclOffset(graph, node)))

			a.LdmfdMultiple([]Register{R10, R11, PC})
			a.CommentPreviousLine("Return from the procedure with params")
		case "if":
			a.ReadIf(graph, child)
		}
	}
}

func (a *AssemblyFile) Call(node int, graph Graph, name int, args int) {
	if graph.GetNode(name) == "new_line" {
		a.Sub(SP, 4)
		a.Mov(R0, 10)
		a.Str(R0)
		a.MovRegister(R0, SP)
		a.CallProcedure("println")
		a.Add(SP, 4)
		return
	}
	if graph.GetNode(name) == "put" {
		a.AddComment("Put statement")
		a.ReadOperand(graph, graph.GetChildren(args)[0])

		// TODO: check for variable that could be a char (or char in record)
		isChar := graph.GetNode(graph.GetChildren(args)[0]) == "cast" || graph.GetNode(graph.GetChildren(args)[0])[0] == '\''
		if isChar {
			a.AddComment("Printing char")
			// Move the result to R0
			a.Ldr(R0, 0)

			a.Sub(SP, 4)

			// Store result in SP with a zero first
			a.Mov(R1, 0)
			a.Str(R1)
			a.Sub(SP, 4)
			a.Str(R0)
			a.MovRegister(R0, SP)
		} else {
			// Move the result to R0
			a.Ldr(R0, 0)

			addr := a.Fill(12)
			a.LdrAddr(R3, addr)

			// Cast the result
			a.CallProcedure("to_ascii")

			a.LdrAddr(R0, addr)
		}

		/*
			if to string
			addr                    FILL    12
			                        LDR     R3, =addr
			                        BL      to_ascii
			                        ldr     r0, =addr
		*/

		a.CallProcedure("println")

		if isChar {
			a.Add(SP, 12)
		} else {
			// do nothing
			a.Add(SP, 4) // why wouldn't the value be removed from the stack?
		}
		a.AddComment("End of put statement")
		return
	}

	symbol := graph.fullSymbols[name]
	_, isFunction := symbol.(Function)
	_, isProcedure := symbol.(Procedure)
	if isFunction {
		a.Sub(SP, getTypeSize(symbol.(Function).ReturnType, *graph.getScope(node)))
		a.CommentPreviousLine("Save space for the return value")
	}

	a.AddComment("Read the arguments")
	removedOffset := 0
	for k, arg := range graph.GetChildren(args) {
		a.ReadOperand(graph, arg)

		// Reserve 4 bytes if the argument is in out mode
		if isFunction {
			if symbol.(Function).Params[k+1].IsParamIn && symbol.(Function).Params[k+1].IsParamOut {
				a.Sub(SP, 4)
				a.CommentPreviousLine("Reserve space for the in out parameter")
				a.StoreAddress(arg, graph)
				a.AddComment("Stored the address of the in out parameter")

				removedOffset += 4
			}
			removedOffset += getTypeSize(symbol.(Function).Params[k+1].Type(), *graph.getScope(node))
		} else if isProcedure {
			if symbol.(Procedure).Params[k+1].IsParamIn && symbol.(Procedure).Params[k+1].IsParamOut {
				a.Sub(SP, 4)
				a.CommentPreviousLine("Reserve space for the in out parameter")
				a.StoreAddress(arg, graph)
				a.AddComment("Stored the address of the in out parameter")

				removedOffset += 4
			}
			removedOffset += getTypeSize(symbol.(Procedure).Params[k+1].Type(), *graph.getScope(node))
		}
	}

	a.AddComment("Arguments read, call the procedure")
	a.CallWithParameters(graph.symbols[name], graph.getScope(node), removedOffset) // TODO: record fix
}

func (a *AssemblyFile) StoreAddress(node int, graph Graph) {
	scope := graph.getScope(node)

	endScope, offset := goUpScope(graph, scope, node, graph.GetNode(node))

	if scope == endScope {
		a.MovRegister(R0, R11)
		a.Add(R0, offset)
	} else {
		// Loop through dynamic links until we reach the correct region
		random := strconv.Itoa(rand.Int()) + "_" + graph.GetNode(node)

		a.MovRegister(R9, R11)
		a.LdrFromFramePointer(R8, 4)
		a.Cmp(R8, endScope.Region)
		a.BranchToLabelWithCondition("notload_"+random, EQ)
		a.AddLabel("load_" + random)
		a.LdrFromFramePointer(R11, 8)
		a.LdrFromFramePointer(R8, 4)
		a.Cmp(R8, endScope.Region)
		a.BranchToLabelWithCondition("load_"+random, NE)

		a.AddLabel("notload_" + random)

		// Go 1 level up
		a.LdrFromFramePointer(R11, 8)

		a.MovRegister(R0, R11)
		a.Add(R0, offset)

		// Restore R9
		a.MovRegister(R11, R9)
	}
	// Move the stack pointer

	a.Str(R0)
	a.CommentPreviousLine("Store the value of " + graph.GetNode(node))
}

func (a *AssemblyFile) ReadWhile(graph Graph, node int) {
	a.AddComment("While statement")
	a.AddComment("Start of condition")

	randomLabel := strconv.Itoa(rand.Int())

	a.AddLabel("while" + randomLabel)

	// Read condition
	condition := graph.GetChildren(node)[0]
	a.ReadOperand(graph, condition)

	a.Ldr(R0, 0)
	a.CommentPreviousLine("Load result of while condition")
	a.Add(SP, 4)

	a.Cmp(R0, 0)
	a.AddComment("End of condition")
	a.BranchToLabelWithCondition("endwhile"+randomLabel, "EQ")

	// Read body
	a.ReadBody(graph, graph.GetChildren(node)[1])

	// Go to the beginning of the loop
	a.BranchToLabel("while" + randomLabel)

	// End of the while loop
	a.AddLabel("endwhile" + randomLabel)

	a.AddComment("End of while statement")
}

func (a *AssemblyFile) ReadFor(graph Graph, node int) {
	goodCounter := a.ForCounter
	a.ForCounter++

	a.AddComment("Loop #" + strconv.Itoa(goodCounter) + " start")

	a.StmfdMultiple([]Register{R10, R11})

	a.Mov(R10, getRegion(graph, node))
	a.MovRegister(R11, SP)
	a.Sub(R11, 4)

	// Reserve space for the index
	a.Sub(SP, 4)
	a.CommentPreviousLine("Reserve space for the index")

	children := graph.GetChildren(node)
	counterStart, err := strconv.Atoi(graph.GetNode(children[2]))
	if err != nil {
		a.ReadOperand(graph, children[2])
		a.Ldr(R0, 0)
		a.CommentPreviousLine("Load to R0 the value of the counter")
		a.Add(SP, 4)
	} else {
		a.Mov(R0, counterStart)
		a.CommentPreviousLine("Load to R0 the value of the counter: " + strconv.Itoa(counterStart))
	}

	a.Str(R0)
	a.CommentPreviousLine("Store the value of the counter")

	counterEnd, err := strconv.Atoi(graph.GetNode(children[3]))
	if err != nil {
		a.ReadOperand(graph, children[3])
		a.Ldr(R1, 0)
		a.CommentPreviousLine("Load to R1 the value of the max")
	} else {
		// Reserve space for the max
		a.Sub(SP, 4)
		a.CommentPreviousLine("Reserve space for the max")

		a.Mov(R1, counterEnd)
		a.CommentPreviousLine("Load to R1 the value of the max: " + strconv.Itoa(counterEnd))
	}
	a.Str(R1)
	a.CommentPreviousLine("Store the value of the max")

	a.AddLabel("for" + strconv.Itoa(goodCounter))

	a.Ldr(R0, 4)
	a.CommentPreviousLine("Load to R0 the value of the counter")

	a.Ldr(R1, 0)
	a.CommentPreviousLine("Load to R1 the value of the max")

	a.CmpRegisters(R0, R1)
	a.CommentPreviousLine("Compare the counter with the max")
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
	a.Add(SP, 8)

	a.LdmfdMultiple([]Register{R10, R11})
	a.CommentPreviousLine("End for loop")

	a.AddComment("Loop #" + strconv.Itoa(goodCounter) + " end")
}

type DeclMode int

const (
	OnlyFuncAndProc DeclMode = iota
	OnlyVar
	All
)

func (a *AssemblyFile) ReadDecl(graph Graph, node int, mode DeclMode) {
	// Read all the children
	children := graph.GetChildren(node)

	// Extend sort for var nodes
	slices.SortFunc(children, func(a, b int) int {
		nodeA := graph.GetNode(a)
		nodeB := graph.GetNode(b)

		if nodeA == "var" && nodeB == "var" {
			sortedA := maps.Keys(graph.gmap[a])
			slices.Sort(sortedA)
			sortedB := maps.Keys(graph.gmap[b])
			slices.Sort(sortedB)

			nameA := graph.GetNode(sortedA[0])
			nameB := graph.GetNode(sortedB[0])
			return strings.Compare(nameA, nameB)
		}
		return strings.Compare(nodeA, nodeB)
	})

	for _, child := range children {
		switch graph.GetNode(child) {
		case "var":
			if mode == All || mode == OnlyVar {
				a.ReadVar(graph, child)
			}
		case "procedure":
			if mode == All || mode == OnlyFuncAndProc {
				a.ReadProcedure(graph, child)
			}
		case "function":
			if mode == All || mode == OnlyFuncAndProc {
				a.ReadProcedure(graph, child)
			}
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
	a.AddLabel(graph.symbols[node])

	a.StmfdMultiple([]Register{R10, R11, LR})
	a.Mov(R10, getRegion(graph, node))
	a.CommentPreviousLine("Current region")
	a.MovRegister(R11, SP)
	a.Sub(R11, 4) // SP points to R10 so we need to subtract 4

	a.AddComment("Read variable declarations")
	a.ReadDecl(graph, declNode, OnlyVar)

	a.AddComment("Read the body of the procedure")
	// Read the body of the procedure
	a.ReadBody(graph, bodyNode)

	// Return the in out parameters
	symbol := graph.fullSymbols[node]
	_, isFunction := symbol.(Function)
	_, isProcedure := symbol.(Procedure)
	if isFunction {
		paramOffset := 0
		for _, param := range symbol.(Function).Params {
			if param.Offset > paramOffset {
				paramOffset = param.Offset
				if param.IsParamIn && param.IsParamOut {
					paramOffset += 4
				}
			}
		}
		for i := 1; i < len(symbol.(Function).Params)+1; i++ {
			if symbol.(Function).Params[i].IsParamIn && symbol.(Function).Params[i].IsParamOut {
				// Load the address of the in out parameter
				a.LdrFromFramePointer(R0, paramOffset-symbol.(Function).Params[i].Offset+16)
				a.CommentPreviousLine("Load the address of the in out parameter")
				a.LdrFromFramePointer(R1, paramOffset-symbol.(Function).Params[i].Offset+4+16)
				a.CommentPreviousLine("Load the value of the in out parameter")
				a.StrFrom(R1, R0, 0)
				a.CommentPreviousLine("Store the value of the in out parameter")
			}
		}
	} else if isProcedure {
		paramOffset := 0
		for _, param := range symbol.(Procedure).Params {
			if param.Offset > paramOffset {
				paramOffset = param.Offset
			}
		}
		for i := 1; i < len(symbol.(Procedure).Params)+1; i++ {
			if symbol.(Procedure).Params[i].IsParamIn && symbol.(Procedure).Params[i].IsParamOut {
				// Load the address of the in out parameter
				a.LdrFromFramePointer(R0, paramOffset-symbol.(Procedure).Params[i].Offset+16)
				a.CommentPreviousLine("Load the address of the in out parameter")
				a.LdrFromFramePointer(R1, paramOffset-symbol.(Procedure).Params[i].Offset+4+16)
				a.CommentPreviousLine("Load the value of the in out parameter")
				a.StrFrom(R1, R0, 0)
				a.CommentPreviousLine("Store the value of the in out parameter")
			}
		}
	}

	a.Add(SP, getDeclOffset(graph, node))
	a.CommentPreviousLine("Clear the stack of declarations: " + strconv.Itoa(getDeclOffset(graph, node)))

	a.LdmfdMultiple([]Register{R10, R11, PC})
	a.CommentPreviousLine("Base return of procedure " + procedureName)

	a.AddComment("End of procedure " + procedureName)

	a.WritingAtEnd = false

	a.ReadDecl(graph, declNode, OnlyFuncAndProc)
}

func (a *AssemblyFile) ReadVar(graph Graph, node int) {
	var nameList []string
	var hasValue bool

	name := graph.GetChildren(node)[0]

	if graph.GetNode(name) == "sametype" {
		for _, child := range graph.GetChildren(name) {
			nameList = append(nameList, graph.GetNode(child))
		}
	} else {
		nameList = append(nameList, graph.GetNode(name))
	}
	if len(graph.GetChildren(node)) > 2 {
		hasValue = true
	}

	// use values from the symbol table
	scope := graph.getScope(node)

	for _, name := range nameList {
		for _, symbol := range scope.Table[name] {
			if variable, ok := symbol.(Variable); ok {
				if hasValue {
					//offset := getTypeSize(variable.SType, *scope)

					a.AddComment("Read the value of " + name)

					value := graph.GetChildren(node)[2]
					a.ReadOperand(graph, value)
					a.Ldr(R0, 0)

					// Load the int value to r0
					a.Str(R0)
					a.CommentPreviousLine("Store the value of " + name)
				} else {
					offset := getTypeSize(variable.SType, *scope)

					// Move the stack pointer
					a.Sub(SP, offset)
					a.CommentPreviousLine("Reserve space for the value of " + name)
				}
			}
		}
	}
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
			a.Str(R0)
		} else {
			if graph.GetNode(node) == "true" {
				// Move the stack pointer
				a.Sub(SP, 4)

				// The operand is a bool
				// Load the bool value to r0
				a.Mov(R0, 1)
				a.Str(R0)
			} else if graph.GetNode(node) == "false" {
				// Move the stack pointer
				a.Sub(SP, 4)

				// The operand is a bool
				// Load the bool value to r0
				a.Mov(R0, 0)
				a.Str(R0)
			} else if graph.GetNode(node)[0] == '\'' {
				// Move the stack pointer
				a.Sub(SP, 4)

				// The operand is a char
				// Load the char value to r0
				a.Mov(R0, int(graph.GetRealNode(node)[1]))
				a.Str(R0)
			} else {
				// The operand is an ident
				// Load the ident value to r0

				// Get the address of the ident using the symbol table
				scope := graph.getScope(node)

				endScope, offset := goUpScope(graph, scope, node, graph.GetNode(node))

				a.MovRegister(R9, R11)
				if scope == endScope {
					a.AddComment(fmt.Sprintf("(S) Load the value of %v", graph.GetNode(node)))
				} else {
					// Loop through dynamic links until we reach the correct region
					random := strconv.Itoa(rand.Int()) + "_" + graph.GetNode(node)

					a.LdrFromFramePointer(R8, 4)
					a.Cmp(R8, endScope.Region)
					a.BranchToLabelWithCondition("notload_"+random, EQ)
					a.AddLabel("load_" + random)
					a.LdrFromFramePointer(R11, 8)
					a.LdrFromFramePointer(R8, 4)
					a.Cmp(R8, endScope.Region)
					a.BranchToLabelWithCondition("load_"+random, NE)

					a.AddLabel("notload_" + random)

					// Go 1 level up
					a.LdrFromFramePointer(R11, 8)

					a.AddComment(fmt.Sprintf("(NS) Load the value of %v", graph.GetNode(node)))
				}

				// What's the size of the type?
				typeSize := getTypeSize(endScope.Table[graph.GetNode(node)][0].Type(), *scope)

				for typeSize > 0 {
					a.LdrFromFramePointer(R0, offset)

					// Move the stack pointer
					a.Sub(SP, 4)
					a.CommentPreviousLine("Reserve space for the value of " + graph.GetNode(node))

					a.Str(R0)
					a.CommentPreviousLine("Store the value of " + graph.GetNode(node))

					offset += 4
					typeSize -= 4
				}

				// Restore R9
				a.MovRegister(R11, R9)
			}
		}
	}

	switch graph.GetNode(node) {
	case "+":
		// Read left operand
		a.ReadOperand(graph, children[0])

		// Read right operand
		a.ReadOperand(graph, children[1])
		a.Ldr(R0, 0)
		a.AddWithOffset(R0, R1, 4) // same as ldr from offset 8 then add

		a.Add(SP, 4)

		// Save the result in stack
		a.Str(R0)
	case "-":
		// Read left operand
		a.ReadOperand(graph, children[0])

		// Read right operand
		a.ReadOperand(graph, children[1])
		a.Ldr(R0, 0)
		a.SubWithOffset(R0, R1, 4)

		a.Add(SP, 4)

		// Save the result in stack
		a.Str(R0)
	case "*":
		// Read left operand
		a.ReadOperand(graph, children[0])

		// Read right operand
		a.ReadOperand(graph, children[1])

		// Left operand in R1, right operand in R2
		a.Ldr(R1, 0)
		a.Ldr(R2, 4)

		// Use the multiplication algorithm at the label mul
		a.CallProcedure("mul")

		a.Add(SP, 4)

		// Save the result in stack
		a.Str(R0)
	case "/", "rem":
		// Read left operand
		a.ReadOperand(graph, children[0])

		// Read right operand
		a.ReadOperand(graph, children[1])

		// Left operand in R0, right operand in R1
		a.Ldr(R2, 0)
		a.Ldr(R1, 4)

		// Use the division algorithm at the label div32
		a.CallProcedure("div32")

		if graph.GetNode(node) == "rem" {
			a.MovRegister(R0, R1)
		}

		a.Add(SP, 4)

		// Save the result in stack
		a.Str(R0)
	case "and":
		// Read left operand
		a.ReadOperand(graph, children[0])

		// Read right operand
		a.ReadOperand(graph, children[1])

		// Left operand in R1, right operand in R2
		a.Ldr(R1, 0)
		a.Ldr(R2, 4)

		// Use the AND operation
		a.And(R1, R2)

		a.Add(SP, 4)

		// Save the result in stack
		a.Str(R0)
	case ">", "=", "<", "<=", ">=", "/=", "!=":
		// Read left operand
		a.ReadOperand(graph, children[0])

		// Read right operand
		a.ReadOperand(graph, children[1])

		// Left operand in R0, right operand in R1
		a.Ldr(R1, 0)
		a.Ldr(R0, 4)

		// Compare the operands
		a.CmpRegisters(R0, R1)
		switch graph.GetNode(node) {
		case ">":
			a.MovCond(R0, 1, GT)
			a.MovCond(R0, 0, LE)
		case "=":
			a.MovCond(R0, 1, EQ)
			a.MovCond(R0, 0, NE)
		case "<":
			a.MovCond(R0, 1, LT)
			a.MovCond(R0, 0, GE)
		case "<=":
			a.MovCond(R0, 1, LE)
			a.MovCond(R0, 0, GT)
		case ">=":
			a.MovCond(R0, 1, GE)
			a.MovCond(R0, 0, LT)
		case "/=", "!=":
			a.MovCond(R0, 1, NE)
			a.MovCond(R0, 0, EQ)
		}

		a.Add(SP, 4)

		// Save the result in stack
		a.Str(R0)
	case "cast":
		// Read the operand
		a.ReadOperand(graph, children[1])

		// Move the result to R0
		a.Ldr(R0, 0)

		/*addr := a.Fill(12)
		a.LdrAddr(R3, addr)

		// Cast the result
		a.CallProcedure("to_ascii")

		a.LdrAddr(R0, addr)*/
	case "call":
		if graph.GetNode(children[0]) == "-" {
			// Read right operand
			a.ReadOperand(graph, children[1])

			a.Ldr(R0, 0)
			a.Negate(R0)

			a.Str(R0)
		} else if graph.GetNode(children[0]) == "not" {
			// Read right operand
			a.ReadOperand(graph, children[1])

			a.Ldr(R0, 0)
			a.Not(R0)

			a.Str(R0)
		} else {
			name := graph.GetChildren(node)[0]
			var args int
			if len(graph.GetChildren(node)) > 1 {
				args = graph.GetChildren(node)[1]
			} else {
				args = graph.GetChildren(node)[0]
			}

			a.Call(node, graph, name, args)

			// Move the stack pointer
			/*size := getTypeSize(graph.fullSymbols[graph.GetChildren(node)[0]].(Function).ReturnType, *graph.getScope(node))
			fmt.Println(size)
			a.Sub(SP, getTypeSize(graph.fullSymbols[graph.GetChildren(node)[0]].(Function).ReturnType, *graph.getScope(node)))

			o := 0
			for size > 0 {
				a.Ldr(R0, o)
				a.StrWithOffset(R0, o)
				size -= 4
				o += 4
			}*/
		}
	}

	if graph.GetNode(node) == "access" {
		// The operand is an ident
		// Load the ident value to r0

		// Get the address of the ident using the symbol table
		scope := graph.getScope(node)

		endScope, offset := goUpScope(graph, scope, node, graph.GetNode(node))

		if scope == endScope {
			a.LdrFrom(R0, R11, offset)
			a.CommentPreviousLine(fmt.Sprintf("(S) Load the value of %v", graph.GetNode(node)))
		} else {
			// Loop through dynamic links until we reach the correct region
			random := strconv.Itoa(rand.Int()) + "_" + graph.GetNode(node)

			a.MovRegister(R9, R11)
			a.LdrFromFramePointer(R8, 4)
			a.Cmp(R8, endScope.Region)
			a.BranchToLabelWithCondition("notload_"+random, EQ)
			a.AddLabel("load_" + random)
			a.LdrFromFramePointer(R11, 8)
			a.LdrFromFramePointer(R8, 4)
			a.Cmp(R8, endScope.Region)
			a.BranchToLabelWithCondition("load_"+random, NE)

			a.AddLabel("notload_" + random)

			// Go 1 level up
			a.LdrFromFramePointer(R11, 8)

			a.LdrFromFramePointer(R0, offset)

			// Restore R9
			a.MovRegister(R11, R9)

			a.CommentPreviousLine(fmt.Sprintf("(NS) Load the value of %v", graph.GetNode(node)))
		}
		// Move the stack pointer
		a.Sub(SP, 4)
		a.CommentPreviousLine("Reserve space for the value of " + graph.GetNode(node))

		a.Str(R0)
		a.CommentPreviousLine("Store the value of " + graph.GetNode(node))
	}
}
