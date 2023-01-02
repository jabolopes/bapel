// DO NOT EDIT - THIS CODE HAS BEEN AUTOMATICALLY GENERATED FROM vm/unary_ops.m4

include(`common_ops.m4')

ifelse(`PRINTU
mode: either immediate, variable, or stack.
typ: optype for op.
operand: value to print.')
define(PRINTU,
`opPrintImpl($3)')

ifelse(`PRINTS
mode: either immediate, variable, or stack.
typ: optype for op.
operand: value to print.')
define(PRINTS,
`opPrintImpl(GET_SIGNED(`$2')($3))')

ifelse(`PUSH
mode: either immediate, variable, or stack.
typ: optype for op.
operand: value to push.')
define(PUSH,
`ifelse(`$1', `immediate', `machine.Stack().Push$2($3)',
 ifelse(`$1', `variable', `machine.Stack().Push$2(varPc$2(machine))',
 ifelse(`$1', `stack', `return errors.New("Unimplemented")')))')

ifelse(`NEG
mode: either immediate, variable, or stack.
typ: optype for op.
operand: value to neg.')
define(NEG,
`machine.Stack().Push$2(- $3)')

package vm

import (
       "errors"

       "github.com/jabolopes/bapel/ir"
)

UNARY_OP_MODES(opPrintU, `PRINTU')
UNARY_OP_MODES(opPrintS, `PRINTS')
UNARY_OP_MODES(opPush, `PUSH')
UNARY_OP_MODES(opNeg, `NEG')

BINARY_OP_MODES(opAdd, `+')
