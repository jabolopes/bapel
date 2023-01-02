// DO NOT EDIT - THIS CODE HAS BEEN AUTOMATICALLY GENERATED FROM binops.m4

ifelse(`GET_MODE:
mode: either immediate, variable, or stack.')
define(GET_MODE, `ifelse(`$1', `immediate', `ir.ImmediateMode',
                  ifelse(`$1', `variable', `ir.VarMode',
                  ifelse(`$1', `stack', `ir.StackMode')))')

ifelse(`GET_OPCODE:
mode: mode for op's 1st argument.
typ: optype for op.')
define(GET_OPCODE, `ir.UnaryOpCode(base, $1, ir.$2)')

ifelse(`GET_OPERAND:
mode: either immediate, variable, or stack.
typ: type of value to get.')
define(GET_OPERAND, `ifelse(`$1', `immediate', `machine.Tape().Get$2()',
                     ifelse(`$1', `variable', `varPc$2(machine)',
                     ifelse(`$1', `stack', `machine.Stack().Pop$2()')))')

ifelse(`GET_SIGNED:
typ: type of value.')
define(GET_SIGNED, `ifelse(`$1', `I8', `int8',
                    ifelse(`$1', `I16', `int16',
                    ifelse(`$1', `I32', `int32',
                    ifelse(`$1', `I64', `int64'))))')

ifelse(`UNARY_OP
mode: either immediate, variable, or stack.
typ: optype for op.
op: operation to perform on values, e.g., +.')
define(UNARY_OP,
`GET_OPCODE(GET_MODE($1), $2): func(machine *Machine)error {
  $3(`$2', GET_OPERAND(`$1', `$2'))
  return nil
},')

ifelse(`UNARY_OP_TYPES
mode: either immediate, variable, or stack.
op: operation to perform on values, e.g., +,')
define(UNARY_OP_TYPES,
`UNARY_OP(`$1', I8, `$2')
UNARY_OP(`$1', I16, `$2')
UNARY_OP(`$1', I32, `$2')
UNARY_OP(`$1', I64, `$2')')

ifelse(`UNARY_OP_MODES
symbol: name of the symbol to create
op: operation to perform on values, e.g., +.')
define(UNARY_OP_MODES,
`func $1(base ir.OpCode) opFamilyMap {
return opFamilyMap {
UNARY_OP_TYPES(`immediate', `$2')
UNARY_OP_TYPES(`variable', `$2')
UNARY_OP_TYPES(`stack', `$2')
}
}')

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
