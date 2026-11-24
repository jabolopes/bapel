implements bapel.compiler.ir

imports {
  bapel.core
  bapel.stl
}

pub type IrFunction = struct {
  id: String,
  type_params: Vector (TypeParam IrType),
  args: Vector FunctionArg,
  ret_type: IrType,
  body: IrTerm,
  export_flag: bool,
  pos: Pos
}

pub type IrSignature = struct {
  id: String,
  type_params: Vector (TypeParam IrType),
  args: Vector FunctionArg,
  ret_type: IrType
}

pub type IrDeclData = variant {
  term_decl struct {
    id: String,
    typ: IrType
  },
  alias_decl struct {
    id: String,
    kind: IrKind,
    typ: IrType
  },
  name_decl struct {
    id: String,
    kind: IrKind
  },
  trait_decl struct {
    id: String,
    type_params: Vector (TypeParam IrType),
    methods: Vector IrSignature
  }
}

pub type IrDecl = struct {
  data: IrDeclData,
  export_flag: bool,
  pos: Pos
}

pub type ImplCase = variant {
  trait_impl (),
  inherent_impl ()
}

pub type IrTraitImpl = struct {
  case_val: ImplCase,
  type_params: Vector (TypeParam IrType),
  trait_type: IrType,
  type_name: IrType,
  methods: Vector IrFunction,
  pos: Pos
}

pub type IrFileCase = variant {
  file_base (),
  file_impl (),
  file_ffi ()
}

pub type IrResolvedImport = struct {
  module_id: ModuleID,
  resolved_filename: Filename,
  deps: Vector ModuleID
}

pub type IrFile = struct {
  case_val: IrFileCase,
  filename: Filename,
  file_imports: Vector IrImport,
  decls: Vector IrDecl,
  trait_impls: Vector IrTraitImpl,
  functions: Vector IrFunction
}

pub type IrModule = struct {
  module_id: ModuleID,
  mod_imports: Vector IrResolvedImport,
  mod_impls: Vector IrImpl,
  mod_flags: Vector String,
  import_decls: Vector IrDecl,
  files: Vector IrFile
}

pub type IrTypedModule = struct {
  module_id: ModuleID,
  mod_imports: Vector IrResolvedImport,
  mod_impls: Vector IrImpl,
  mod_flags: Vector String,
  import_decls: Vector IrDecl,
  ffi_decls: Vector IrDecl,
  decls: Vector IrDecl,
  trait_impls: Vector IrTraitImpl,
  functions: Vector IrFunction
}

impl IrDecl {
  fn id(d: &Self) -> String {
    match (*d).data {
      term_decl t => t.id,
      alias_decl a => a.id,
      name_decl n => n.id,
      trait_decl tr => tr.id
    }
  }
}
