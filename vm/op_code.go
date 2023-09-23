package vm

type OpCode = uint64

// opHalt halts the program.
//
// No operands.

// opCall calls a function at a given offset.
//
// call : (offset immediate : i64)
//
// offset: function to call identified by its absolute offset. The
// offset is an index in the program data.

// opReturn returns from a function.
//
// return : (leaveSize immediate : i16)
//
// leaveSize: size in bytes to deallocate from the stack. This size
// includes the size of locals and args.

// opIfThen tests and jumps.
//
// ifThen : (offset immediate : i64)
//
// opIfThen tests whether the value at the top of the stack is zero
// and if that is the case the pc is incremented by the value given by
// the operand. The value at the top of the stack is i8.

// opIfElse tests and jumps.
//
// ifElse : (offset immediate : i64)
//
// opIfElse tests whether the value at the top of the stack is *not*
// zero and if that is the case the pc is incremented by the value
// given by the operand. The value at the top of the stack is i8.

// opElse unconditionally jumps.
//
// else : (offset immediate : i64)
//
// opElse unconditionally increments the pc by the given offset.

// opPush pushes to the stack.
//
// Immediate mode:
//   push : Type => (value immediate : Type) -> (ret stack : Type)
//
//   Pushes an immediate value onto the stack. The type of the
//   immediate is determined by the type of the 'push' function.
//
// Var mode:
//   push : Type => (offset immediate : i16) -> (ret stack : Type)
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
//   pop : Type => (offset immediate : i16) -> ()
//
//   Pops a value from the stack and copies it to a variable. The
//   variable is identified by its offset relative to the fp.
//
// Stack mode:
//   pop : Type => () -> ()
//
//   Pops a value from the stack (and discards it).
//
//   No operands.
