define(IMMEDIATE, machine.Tape().Get$1())dnl
define(VARIABLE, machine.Frame().Var$1(uint64(machine.Tape().GetI16())))dnl
define(STACK, machine.Stack().Pop$1())dnl
define(BINARY_OP,
`func(machine *Machine)error {
  machine.Stack().Push$2($3($2) $1 $4($2))
  return nil
},')dnl
define(BINARY_OP_TYPES,
`BINARY_OP($1, I8, `$2', `$3')
BINARY_OP($1, I16, `$2', `$3')
BINARY_OP($1, I32, `$2', `$3')
BINARY_OP($1, I64, `$2', `$3')')dnl
define(BINARY_OP_MODES,
`var $1 = []func(*Machine)error {
BINARY_OP_TYPES($2, `IMMEDIATE', `IMMEDIATE')
BINARY_OP_TYPES($2, `IMMEDIATE', `VARIABLE')
BINARY_OP_TYPES($2, `IMMEDIATE', `STACK')
BINARY_OP_TYPES($2, `VARIABLE', `IMMEDIATE')
BINARY_OP_TYPES($2, `VARIABLE', `VARIABLE')
BINARY_OP_TYPES($2, `VARIABLE', `STACK')
BINARY_OP_TYPES($2, `STACK', `IMMEDIATE')
BINARY_OP_TYPES($2, `STACK', `VARIABLE')
BINARY_OP_TYPES($2, `STACK', `STACK')
}')dnl
// DO NOT EDIT - THIS CODE HAS BEEN AUTOMATICALLY GENERATED FROM binops.m4
package vm

BINARY_OP_MODES(opAdd, `+')
