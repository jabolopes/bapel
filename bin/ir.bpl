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

pub type IrDecl = struct {
  id: String,
  is_export: bool,
  decl_kind: String
}

pub type IrTraitImpl = struct {
  trait_name: String,
  type_name: String,
  type_params: Vector String,
  methods: Vector String
}

pub type IrFunction = struct {
  name: String,
  ret_type: String,
  params: Vector IrField,
  type_params: Vector String,
  body: String,
  is_export: bool
}

pub type IrUnit = struct {
  module_id: String,
  import_modules: Vector String,
  impl_files: Vector String,
  decls: Vector IrDecl,
  functions: Vector IrFunction,
  trait_impls: Vector IrTraitImpl
}




