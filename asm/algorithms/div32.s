       MOV     R0, #-12
       MOV     R1, #3
       BL      div32
       END

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