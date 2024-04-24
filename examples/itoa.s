STR_OUT       FILL    0x1000
main          LDR     R0, =0x75BCD15 ; =-123456789 à transformer en String
addr          FILL    12
              LDR     R3, =addr
              bl      to_ascii
              ldr     r0, =addr
              bl      println
              end

println       STMFD   SP!, {LR, R0-R3}
              MOV     R3, R0
              LDR     R1, =STR_OUT ; address of the output buffer
PRINTLN_LOOP  LDRB    R2, [R0], #1
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
CLEAN         LDRB    R2, [R0], #1
              MOV     R3, #0
              STRB    R3, [R1], #1
              TST     R2, R2
              BNE     CLEAN
              ;       clear 3 more because why not
              STRB    R3, [R1], #1
              STRB    R3, [R1], #1
              STRB    R3, [R1], #1

              LDMFD   SP!, {PC, R0-R3}

to_ascii      STMFD   SP!, {LR, R4-R6}
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

              LDMFD   SP!, {PC, R4-R6}


              ;       Multiplication algorithm
              ;       R0 = result, R1 = multiplicand, R2 = multiplier
mul           STMFD   SP!, {LR}
              MOV     R0, #0
mul_loop      LSRS    R2, R2, #1
              ADDCS   R0, R0, R1
              LSL     R1, R1, #1
              TST     R2, R2
              BNE     mul_loop
              LDMFD   SP!, {PC}

              ;       Integer division routine
              ;       Arguments:
              ;       R0 = Dividend
              ;       R1 = Divisor
              ;       Returns:
              ;       R0 = Quotient
              ;       R1 = Remainder
div32         STMFD   SP!, {LR, R2-R5}
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
div_max       LSL     R4, R4, #1
              LSL     R5, R5, #1
              CMP     R4, R1
              BLE     div_max
div_loop      LSR     R4, R4, #1
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
div_exit      LDMFD   SP!, {PC, R2-R5}
