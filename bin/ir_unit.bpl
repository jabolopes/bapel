implements bin.ir

pub type IrUnit = struct {
  module_id: String,
  import_modules: Vector String,
  impl_files: Vector String,
  decls: Vector IrDecl,
  functions: Vector IrFunction,
  trait_impls: Vector IrTraitImpl
}
