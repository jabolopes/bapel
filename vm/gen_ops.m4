// DO NOT EDIT - THIS CODE HAS BEEN AUTOMATICALLY GENERATED FROM vm/unary_ops.m4

include(`common_ops.m4')

ifelse(`PRINTU
typ: optype for op.
value: value to print.')
define(PRINTU,
`opPrintImpl($2)')

ifelse(`PRINTS
typ: optype for op.
value: value to print.')
define(PRINTS,
`opPrintImpl(GET_SIGNED(`$1')($2))')

ifelse(`NEG
typ: optype for op.
value: value to neg.')
define(NEG,
`machine.Stack().Push$1(- $2)')

package vm

import (
       "github.com/jabolopes/bapel/ir"
)

UNARY_OP_MODES(opPrintU, `PRINTU')
UNARY_OP_MODES(opPrintS, `PRINTS')
UNARY_OP_MODES(opNeg, `NEG')

BINARY_OP_MODES(opAdd, `+')
