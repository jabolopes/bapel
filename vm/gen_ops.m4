// DO NOT EDIT - THIS CODE HAS BEEN AUTOMATICALLY GENERATED.
package vm

import (
       "errors"
       "fmt"

       "github.com/jabolopes/bapel/ir"
)

include(`common_ops.m4')

ifelse(`GET_OPERAND:
mode: either immediate, variable, or stack.
typ: type of value to get.')
define(GET_OPERAND, `ifelse(`$1', `immediate', `machine.Tape().Get$2()',
                     ifelse(`$1', `variable', `varPc$2(machine)',
                     ifelse(`$1', `stack', `machine.Stack().Pop$2()')))')

ifelse(`PRINTU
mode: either immediate, variable, or stack.
typ: optype for op.')
define(PRINTU,
`func(machine *Machine) error {
   fmt.Printf("%d\n", GET_OPERAND(`$1', `$2'))
   return nil
}')

ifelse(`PRINTS
mode: either immediate, variable, or stack.
typ: optype for op.')
define(PRINTS,
`func(machine *Machine) error {
   fmt.Printf("%d\n", GET_SIGNED(`$2')(GET_OPERAND(`$1', `$2')))
   return nil
}')

ifelse(`PUSH
mode: either immediate, variable, or stack.
typ: optype for op.')
define(PUSH,
`func(machine *Machine) error {
 ifelse(`$1', `immediate', `machine.Stack().Push$2(GET_OPERAND(`$1', `$2'))',
 ifelse(`$1', `variable', `machine.Stack().Push$2(varPc$2(machine))',
 ifelse(`$1', `stack', `return errors.New("Unimplemented")')))
 return nil
}')

ifelse(`POP
mode: either immediate, variable, or stack.
typ: optype for op.
operand: value to pop.')
define(POP,
`func(machine *Machine) error {
 ifelse(`$1', `immediate', `return errors.New("Unimplemented")',
 ifelse(`$1', `variable', `setVarPc$2(machine, machine.Stack().Pop$2())',
 ifelse(`$1', `stack', `_ = machine.Stack().Pop$2()')))
 return nil
}')

ifelse(`NEG
mode: either immediate, variable, or stack.
typ: optype for op.
operand: value to neg.')
define(NEG,
`func(machine *Machine) error {
 machine.Stack().Push$2(- GET_OPERAND(`$1', `$2'))
 return nil
}')

ifelse(`ADD
mode1: either immediate, variable, or stack.
mode2: either immediate, variable, or stack.
typ: optype for op.')
define(ADD,
`func(machine *Machine) error {
 machine.Stack().Push$3(GET_OPERAND($1, $3) + GET_OPERAND($2, $3))
 return nil
}')

UNARY_OP_MODES(opPrintU, `PRINTU')
UNARY_OP_MODES(opPrintS, `PRINTS')
UNARY_OP_MODES(opPush, `PUSH')
UNARY_OP_MODES(opPop, `POP')
UNARY_OP_MODES(opNeg, `NEG')

BINARY_OP_MODES(opAdd, `ADD')
