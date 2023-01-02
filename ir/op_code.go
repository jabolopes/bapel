package ir

type OpCode = uint64

// opPush pushes to the stack.
//
// Immediate mode:
//   push : Type => (value : Type)
//
//   Pushes an immediate value onto the stack. The type of the
//   immediate is determined by the type of the 'push' function.
//
// Var mode:
//   push : Type => (offset : i16)
//
//   Pushes the value of a varible onto the stack. The variable is
//   identified by its offset relative to the fp.
//
// Stack mode: unimplemented.

// opPop pops from the stack.
//
// Immediate mode: unimplemented
//
// Var mode:
//   pop : Type => (offset : i16)
//
//   Pops a value from the stack and copies it to a variable. The
//   variable is identified by its offset relative to the fp.
//
// Stack mode:
//   pop : Type
//
//   Pops a value from the stack (and discards it).
//
//   No operands.
