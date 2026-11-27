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

  fn new_tuple(elems: Vector IrTerm) -> IrTerm {
    IrTerm::mk (
      variant { (IrTermData IrTerm) tuple_t = elems },
      Optional::none [IrType] (),
      Pos::empty ()
    )
  }

  fn app_types(self: &Self) -> (IrTerm, Vector IrType) {
    match (*self).data {
      app_type at => {
        let res: (IrTerm, Vector IrType) = IrTerm::app_types (at.fun.get ());
        let types: Vector IrType = res.1;
        types.push_back at.arg;
        (res.0, types)
      },
      _ => (*self, Vector::mk [IrType] ())
    }
  }

  fn app_args(self: &Self) -> (IrTerm, Vector IrType, IrTerm) {
    match (*self).data {
      app_term at => {
        let app_t: (IrTerm, Vector IrType) = IrTerm::app_types (at.fun.get ());
        (app_t.0, app_t.1, *(at.arg.get ()))
      },
      _ => (*self, Vector::mk [IrType] (), IrTerm::new_tuple (Vector::mk [IrTerm] ()))
    }
  }

  fn to_function(self: &Self) -> (Vector (TypeParam IrType), Vector FunctionArg, IrTerm) {
    let tvars: Vector (TypeParam IrType) = Vector::mk [TypeParam IrType] ();
    let args: Vector FunctionArg = Vector::mk [FunctionArg] ();
    let cur: IrTerm = *self;
    let peeling_t: bool = true;
    for peeling_t {
      match cur.data {
        type_abs ta => {
          tvars.push_back ta.param;
          cur <- *(ta.body.get ());
          ()
        },
        _ => {
          peeling_t <- false;
          ()
        }
      }
    }
    let peeling_l: bool = true;
    for peeling_l {
      match cur.data {
        lambda l => {
          args.push_back l.arg;
          cur <- *(l.body.get ());
          ()
        },
        _ => {
          peeling_l <- false;
          ()
        }
      }
    }
    (tvars, args, cur)
  }
}
