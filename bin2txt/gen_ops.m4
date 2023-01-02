// DO NOT EDIT - THIS CODE HAS BEEN AUTOMATICALLY GENERATED.
package bin2txt

import (
       "fmt"

       "github.com/jabolopes/bapel/ir"
)

include(`../vm/common_ops.m4')

ifelse(`GET_OPERAND:
mode: either immediate, variable, or stack.
typ: type of value to get.')
define(GET_OPERAND, `ifelse(`$1', `immediate', `disassembler.dec().Get$2()',
                     ifelse(`$1', `variable', `disassembler.dec().Get$2()',
                     ifelse(`$1', `stack', `"sp"')))')

ifelse(`PRINTU
mode: either immediate, variable, or stack.
typ: optype for op.')
define(PRINTU,
`func(disassembler *disassembler) error {
   fmt.Printf("printu %s %v\n", "$2", GET_OPERAND(`$1', `$2'))
   return nil
}')

ifelse(`PRINTS
mode: either immediate, variable, or stack.
typ: optype for op.')
define(PRINTS,
`func(disassembler *disassembler) error {
   fmt.Printf("prints %s %v\n", "$2", GET_OPERAND(`$1', `$2'))
   return nil
}')

ifelse(`PUSH
mode: either immediate, variable, or stack.
typ: optype for op.')
define(PUSH,
`func(disassembler *disassembler) error {
   fmt.Printf("push %s %v\n", "$2", GET_OPERAND(`$1', `$2'))
   return nil
}')

ifelse(`POP
mode: either immediate, variable, or stack.
typ: optype for op.')
define(POP,
`func(disassembler *disassembler) error {
   fmt.Printf("pop %s %v\n", "$2", GET_OPERAND(`$1', `$2'))
   return nil
}')

ifelse(`NEG
mode: either immediate, variable, or stack.
typ: optype for op.')
define(NEG,
`func(disassembler *disassembler) error {
   fmt.Printf("neg %s %v\n", "$2", GET_OPERAND(`$1', `$2'))
   return nil
}')


ifelse(`ADD
mode1: either immediate, variable, or stack.
mode2: either immediate, variable, or stack.
typ: optype for op.')
define(ADD,
`func(disassembler *disassembler) error {
   fmt.Printf("add %s %v\n", "$3", GET_OPERAND(`$1', `$3'), GET_OPERAND(`$2', `$3'))
   return nil
}')

UNARY_OP_MODES(opPrintU, `PRINTU')
UNARY_OP_MODES(opPrintS, `PRINTS')
UNARY_OP_MODES(opPush, `PUSH')
UNARY_OP_MODES(opPop, `POP')
UNARY_OP_MODES(opNeg, `NEG')

BINARY_OP_MODES(opAdd, `ADD')
