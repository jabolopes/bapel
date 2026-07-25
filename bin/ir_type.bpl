implements bin.ir

pub type IrType = variant {
  name String,
  array (String, i64),
  fun (String, String),
  tuple (Vector String),
  variant_type (Vector (String, String)),
  struct_type (Vector (String, String))
}
