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

pub type IrType = variant {
  name String,
  app (String, Vector String),
  array (String, i64),
  fun (String, String),
  tuple (Vector String),
  variant_type (Vector IrField),
  struct_type (Vector IrField),
  ptr String,
  ref String
}

pub type IrTerm = variant {
  var String,
  const_int i64,
  const_float f64,
  const_str String,
  const_bool bool,
  let_term (String, String, String),
  assign_term (String, String),
  return_term String,
  block (Vector String),
  if_term (String, String, String),
  for_loop (String, String, String, String),
  match_term (String, Vector MatchArm),
  app_type_term (String, Vector String),
  app_term (String, Vector String),
  projection_term (String, String, bool),
  injection_term (String, i64, String),
  set_term (String, String, String, bool),
  lambda_term (Vector String, Vector String, String, String),
  tuple_term (Vector String),
  struct_term (Vector IrField)
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
  params: Vector String,
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
