ifelse(`GET_MODE:
mode: either immediate, variable, or stack.')
define(GET_MODE, `ifelse(`$1', `immediate', `ir.ImmediateMode',
                  ifelse(`$1', `variable', `ir.VarMode',
                  ifelse(`$1', `stack', `ir.StackMode')))')

ifelse(`GET_OPCODE1:
mode: mode for op's 1st argument.
typ: optype for op.')
define(GET_OPCODE1, `ir.UnaryOpCode(base, $1, ir.$2)')

ifelse(`GET_OPCODE2:
mode1: mode for op's 1st argument.
mode2: mode for op's 2nd argument.
typ: optype for op.')
define(GET_OPCODE2, `ir.BinaryOpCode(base, $1, $2, ir.$3)')

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
`GET_OPCODE1(GET_MODE($1), $2): $3(`$1', `$2'),')

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


ifelse(`BINARY_OP:
mode1: either immediate, variable, or stack.
mode2: either immediate, variable, or stack.
typ: optype for op.
op: operation to perform on values, e.g., +.')
define(BINARY_OP,
`GET_OPCODE2(GET_MODE($1), GET_MODE($2), $3): $4(`$1', `$2', `$3'),')

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
