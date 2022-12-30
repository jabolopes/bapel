// DO NOT EDIT - THIS CODE HAS BEEN AUTOMATICALLY GENERATED FROM binops.m4

ifelse(`GET_MODE:
mode: either immediate, variable, or stack.')
define(GET_MODE, `ifelse(`$1', `immediate', `ir.ImmediateMode',
                  ifelse(`$1', `variable', `ir.VarMode',
                  ifelse(`$1', `stack', `ir.StackMode')))')

ifelse(`GET_OPCODE:
mode1: mode for op's 1st argument.
mode2: mode for op's 2nd argument.
typ: optype for op.')
define(GET_OPCODE, `ir.BinaryOpCode(base, $1, $2, ir.$3)')

ifelse(`GET_VALUE:
mode: either immediate, variable, or stack.
typ: type of value to get.')
define(GET_VALUE, `ifelse(`$1', `immediate', `machine.Tape().Get$2()',
                   ifelse(`$1', `variable', `varPc$2(machine)',
                   ifelse(`$1', `stack', `machine.Stack().Pop$2()')))')

ifelse(`BINARY_OP
mode1: either immediate, variable, or stack.
mode2: either immediate, variable, or stack.
typ: optype for op.
op: operation to perform on values, e.g., +.')
define(BINARY_OP,
`GET_OPCODE(GET_MODE($1), GET_MODE($2), $3): func(machine *Machine)error {
  machine.Stack().Push$3(GET_VALUE($1, $3) $4 GET_VALUE($2, $3))
  return nil
},')

ifelse(`BINARY_OP_TYPES
mode1: either immediate, variable, or stack.
mode2: either immediate, variable, or stack.
op: operation to perform on values, e.g., +,')
define(BINARY_OP_TYPES,
`BINARY_OP(`$1', `$2', I8, `$3')
BINARY_OP(`$1', `$2', I16, `$3')
BINARY_OP(`$1', `$2', I32, `$3')
BINARY_OP(`$1', `$2', I64, `$3')')

ifelse(`BINARY_OP_MODES
symbol: name of the symbol to create
op: operation to perform on values, e.g., +.')
define(BINARY_OP_MODES,
`func $1(base ir.OpCode) opFamilyMap {
return opFamilyMap {
BINARY_OP_TYPES(`immediate', `immediate', `$2')
BINARY_OP_TYPES(`immediate', `variable', `$2')
BINARY_OP_TYPES(`immediate', `stack', `$2')
BINARY_OP_TYPES(`variable', `immediate', `$2')
BINARY_OP_TYPES(`variable', `variable', `$2')
BINARY_OP_TYPES(`variable', `stack', `$2')
BINARY_OP_TYPES(`stack', `immediate', `$2')
BINARY_OP_TYPES(`stack', `variable', `$2')
BINARY_OP_TYPES(`stack', `stack', `$2')
}
}')

package vm

import (
       "github.com/jabolopes/bapel/ir"
)

BINARY_OP_MODES(opAdd, `+')
