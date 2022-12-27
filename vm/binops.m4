define(GET_VALUE, `ifelse(`$2', `immediate', `machine.Tape().Get$1()',
                   ifelse(`$2', `variable', `machine.Frame().Var$1(uint64(machine.Tape().GetI16()))',
                   ifelse(`$2', `stack', `machine.Stack().Pop$1()')))')dnl
define(BINARY_OP,
`func(machine *Machine)error {
  machine.Stack().Push$2(GET_VALUE($2, $3) $1 GET_VALUE($2, $4))
  return nil
},')dnl
define(BINARY_OP_TYPES,
`BINARY_OP($1, I8, `$2', `$3')
BINARY_OP($1, I16, `$2', `$3')
BINARY_OP($1, I32, `$2', `$3')
BINARY_OP($1, I64, `$2', `$3')')dnl
define(BINARY_OP_MODES,
`var $1 = []func(*Machine)error {
BINARY_OP_TYPES($2, `immediate', `immediate')
BINARY_OP_TYPES($2, `immediate', `variable')
BINARY_OP_TYPES($2, `immediate', `stack')
BINARY_OP_TYPES($2, `variable', `immediate')
BINARY_OP_TYPES($2, `variable', `variable')
BINARY_OP_TYPES($2, `variable', `stack')
BINARY_OP_TYPES($2, `stack', `immediate')
BINARY_OP_TYPES($2, `stack', `variable')
BINARY_OP_TYPES($2, `stack', `stack')
}')dnl
// DO NOT EDIT - THIS CODE HAS BEEN AUTOMATICALLY GENERATED FROM binops.m4
package vm

BINARY_OP_MODES(opAdd, `+')
