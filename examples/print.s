             ;       ########### Print a string ###########
STR_OUT      FILL    0x1000
STR_HELLO    DCD     0x6C6C6548, 0x77202C6F, 0x646C726F, 0x21
HEY          DCD     0x6C6C6548, 0x6C6C6548, 0x6C6C6548, 0x646C726F
start        B       main
             ;R0     contains the address of the null-terminated string to print
println      STMFD   SP!, {LR, R0-R3}
             MOV     R3, R0
             LDR     R1, =STR_OUT ; address of the output buffer
PRINTLN_LOOP LDRB    R2, [R0], #1
             STRB    R2, [R1], #1
             TST     R2, R2
             BNE     PRINTLN_LOOP
             MOV     R2, #10
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


             ;       print HEY
main         MOV     R0, #65
             STR     R0, [SP]
             MOV     R0, SP
             BL      println

             LDR     R0, =STR_HELLO
             BL      println

             B       exit
exit         END