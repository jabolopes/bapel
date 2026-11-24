implements bapel.compiler.ir

imports {
  bapel.core
  bapel.stl
}

pub type MatchArm ['e] = struct {
  tag: String,
  arg: String,
  body: SharedPtr 'e,
  index: Optional i64
}

pub type LabelValue ['e] = struct {
  label: String,
  value: SharedPtr 'e
}

pub type IrTermData ['e] = variant {
  app_term struct {
    fun: SharedPtr 'e,
    arg: SharedPtr 'e
  },
  app_type struct {
    fun: SharedPtr 'e,
    arg: IrType
  },
  assign struct {
    arg: SharedPtr 'e,
    ret: SharedPtr 'e
  },
  block Vector 'e,
  const_t IrLiteral,
  injection struct {
    variant_type: IrType,
    tag: String,
    value: SharedPtr 'e,
    tag_index: Optional i64
  },
  lambda struct {
    arg: FunctionArg,
    body: SharedPtr 'e
  },
  let_t struct {
    var: String,
    var_type: Optional IrType,
    value: SharedPtr 'e
  },
  match_t struct {
    term: SharedPtr 'e,
    arms: Vector (MatchArm 'e)
  },
  projection struct {
    term: SharedPtr 'e,
    label: String,
    reduced_type: Optional IrType
  },
  return_t SharedPtr 'e,
  set_t struct {
    term: SharedPtr 'e,
    values: Vector (LabelValue 'e),
    reduced_type: Optional IrType
  },
  struct_t Vector (LabelValue 'e),
  tuple_t Vector 'e,
  type_abs struct {
    param: TypeParam IrType,
    body: SharedPtr 'e
  },
  var_t String
}

pub type IrTerm = struct {
  data: IrTermData IrTerm,
  typ: Optional IrType,
  pos: Pos
}

impl IrTerm {
  fn mk(data: IrTermData IrTerm, typ: Optional IrType, pos: Pos) -> IrTerm {
    struct { data = data, typ = typ, pos = pos }
  }

  fn new_var(id: String) -> IrTerm {
    IrTerm::mk (
      variant { (IrTermData IrTerm) var_t = id },
      Optional::none [IrType] (),
      Pos::empty ()
    )
  }

  fn new_const(lit: IrLiteral) -> IrTerm {
    IrTerm::mk (
      variant { (IrTermData IrTerm) const_t = lit },
      Optional::none [IrType] (),
      Pos::empty ()
    )
  }
}
