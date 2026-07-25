implements bin.ir

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
