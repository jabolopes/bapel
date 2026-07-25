implements bin.ir

pub type IrType = variant {
  name String,
  array (String, i64),
  fun (String, String),
  tuple (Vector String),
  variant_type (Vector IrField),
  struct_type (Vector IrField)
}
