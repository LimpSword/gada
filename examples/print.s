             ;       ########### Print a string ###########
STR_OUT      FILL    0x1000
STR_HELLO    DCD     0x6C6C6548, 0x77202C6F, 0x646C726F, 0x21
start        B       main

             ;R0     contains the address of the null-terminated string to print
print        STMFA   SP!, {R0-R2}
             LDR     R1, =STR_OUT ; address of the output buffer
PRINT_LOOP   LDRB    R2, [R0], #1
             STRB    R2, [R1], #1
             TST     R2, R2
             BNE     PRINT_LOOP
             LDMFA   SP!, {R0-R2}
             LDR     PC, [R13, #-4]!

             ;R0     contains the address of the null-terminated string to print
println      STMFA   SP!, {R0-R2}
             LDR     R1, =STR_OUT ; address of the output buffer
PRINTLN_LOOP LDRB    R2, [R0], #1
             STRB    R2, [R1], #1
             TST     R2, R2
             BNE     PRINTLN_LOOP
             MOV     R2, #10
             STRB    R2, [R1, #-1]
             MOV     R2, #0
             STRB    R2, [R1]
             LDMFA   SP!, {R0-R2}
             LDR     PC, [R13, #-4]!

main         MOV     R0, #STR_HELLO
             STR     PC, [R13], #4
             BL      println

             B       exit
exit         END