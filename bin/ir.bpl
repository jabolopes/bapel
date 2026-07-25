module bin.ir

imports {
  bapel.core
  bapel.stl
}

impls {
  "ir_type.bpl"
  "ir_term.bpl"
  "ir_decl.bpl"
  "ir_function.bpl"
  "ir_unit.bpl"
}

pub type IrField = struct {
  name: String,
  type_name: String
}

pub type MatchArm = struct {
  index: i64,
  arg_id: String,
  body: String
}



