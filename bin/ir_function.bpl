implements bin.ir

pub type IrFunction = struct {
  name: String,
  ret_type: String,
  params: Vector IrField,
  type_params: Vector String,
  body: String,
  is_export: bool
}
