implements bin.ir

pub type MatchArm = struct {
  index: i64,
  arg_id: String,
  body: String
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
  match_term (String, Vector (i64, String, String))
}
