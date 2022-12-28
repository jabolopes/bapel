// DO NOT EDIT - THIS CODE HAS BEEN AUTOMATICALLY GENERATED FROM binops.m4

ifelse(`GET_MODE:
mode: either immediate, variable, or stack.')
define(GET_MODE, `ifelse(`$1', `immediate', `ImmediateMode',
                  ifelse(`$1', `variable', `VarMode',
                  ifelse(`$1', `stack', `StackMode')))')

ifelse(`GET_OPCODE:
mode: mode for op's 1st argument.
typ: optype for op.')
define(GET_OPCODE, `unaryOpCode(base, $1, $2)')

ifelse(`GET_VALUE:
mode: either immediate, variable, or stack.
typ: type of value to get.')
define(GET_VALUE, `ifelse(`$1', `immediate', `machine.Tape().Get$2()',
                   ifelse(`$1', `variable', `machine.Frame().Var$2(uint64(machine.Tape().GetI16()))',
                   ifelse(`$1', `stack', `machine.Stack().Pop$2()')))')

ifelse(`UNARY_OP
mode: either immediate, variable, or stack.
typ: optype for op.
op: operation to perform on values, e.g., +.')
define(UNARY_OP,
`GET_OPCODE(GET_MODE($1), $2): func(machine *Machine)error {
  $3(GET_VALUE($1, $2))
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
`func $1(base OpCode) map[OpCode]func(*Machine)error {
return map[OpCode]func(*Machine)error {
UNARY_OP_TYPES(`immediate', `$2')
UNARY_OP_TYPES(`variable', `$2')
UNARY_OP_TYPES(`stack', `$2')
}
}')

ifelse(`PRINT
value: value to print.')
define(PRINT,
`fmt.Printf("%d\n", $1)')

package vm

import (
  "fmt"
)

UNARY_OP_MODES(opPrint, `PRINT')
