implements bapel.compiler.ir

imports {
  bapel.core
  bapel.stl
}

pub type TypeParam ['t] = struct {
  var: String,
  kind: IrKind,
  bounds: Vector 't
}

pub type NamedType ['t] = struct {
  id: String,
  typ: SharedPtr 't
}

pub type IrTypeData ['t] = variant {
  app_t struct {
    fun: SharedPtr 't,
    arg: SharedPtr 't
  },
  array_t struct {
    elem: SharedPtr 't,
    size: i64
  },
  exist_var_t i64,
  forall_t struct {
    param: TypeParam 't,
    body: SharedPtr 't
  },
  fun_t struct {
    arg: SharedPtr 't,
    ret: SharedPtr 't
  },
  lambda_t struct {
    var: String,
    kind: IrKind,
    body: SharedPtr 't
  },
  mu_t struct {
    var: String,
    kind: IrKind,
    body: SharedPtr 't
  },
  name_t String,
  struct_t Vector (NamedType 't),
  tuple_t Vector 't,
  variant_t Vector (NamedType 't),
  var_t String
}

pub type IrType = struct {
  data: IrTypeData IrType,
  pos: Pos
}

pub type FunctionArg = struct {
  id: String,
  typ: IrType
}

impl IrType {
  fn mk(data: IrTypeData IrType, pos: Pos) -> IrType {
    struct { data = data, pos = pos }
  }

  fn new_name(name: String) -> IrType {
    IrType::mk (variant { (IrTypeData IrType) name_t = name }, Pos::empty ())
  }

  fn new_var(var: String) -> IrType {
    IrType::mk (variant { (IrTypeData IrType) var_t = var }, Pos::empty ())
  }

  fn new_exist_var(id: i64) -> IrType {
    IrType::mk (variant { (IrTypeData IrType) exist_var_t = id }, Pos::empty ())
  }

  fn new_app(fun: IrType, arg: IrType) -> IrType {
    IrType::mk (
      variant {
        (IrTypeData IrType) app_t = struct {
          fun = SharedPtr::make_shared [IrType] fun,
          arg = SharedPtr::make_shared [IrType] arg
        }
      },
      Pos::empty ()
    )
  }

  fn new_fun(arg: IrType, ret: IrType) -> IrType {
    IrType::mk (
      variant {
        (IrTypeData IrType) fun_t = struct {
          arg = SharedPtr::make_shared [IrType] arg,
          ret = SharedPtr::make_shared [IrType] ret
        }
      },
      Pos::empty ()
    )
  }

  fn new_tuple(elems: Vector IrType) -> IrType {
    IrType::mk (variant { (IrTypeData IrType) tuple_t = elems }, Pos::empty ())
  }

  fn new_struct(fields: Vector (NamedType IrType)) -> IrType {
    IrType::mk (variant { (IrTypeData IrType) struct_t = fields }, Pos::empty ())
  }

  fn new_variant(tags: Vector (NamedType IrType)) -> IrType {
    IrType::mk (variant { (IrTypeData IrType) variant_t = tags }, Pos::empty ())
  }

  fn new_lambda(var: String, kind: IrKind, body: IrType) -> IrType {
    IrType::mk (
      variant {
        (IrTypeData IrType) lambda_t = struct {
          var = var,
          kind = kind,
          body = SharedPtr::make_shared [IrType] body
        }
      },
      Pos::empty ()
    )
  }

  fn new_mu(var: String, kind: IrKind, body: IrType) -> IrType {
    IrType::mk (
      variant {
        (IrTypeData IrType) mu_t = struct {
          var = var,
          kind = kind,
          body = SharedPtr::make_shared [IrType] body
        }
      },
      Pos::empty ()
    )
  }

  fn new_forall(param: TypeParam IrType, body: IrType) -> IrType {
    IrType::mk (
      variant {
        (IrTypeData IrType) forall_t = struct {
          param = param,
          body = SharedPtr::make_shared [IrType] body
        }
      },
      Pos::empty ()
    )
  }

  fn new_array(elem: IrType, size: i64) -> IrType {
    IrType::mk (
      variant {
        (IrTypeData IrType) array_t = struct {
          elem = SharedPtr::make_shared [IrType] elem,
          size = size
        }
      },
      Pos::empty ()
    )
  }
}

pub fn ir_type_equals(a: &IrType, b: &IrType) -> bool {
  match (*a).data {
    app_t aa => match (*b).data {
      app_t bb =>
        ir_type_equals (aa.fun.get(), bb.fun.get()) &&
        ir_type_equals (aa.arg.get(), bb.arg.get()),
      _ => false
    },
    array_t aa => match (*b).data {
      array_t bb =>
        aa.size == bb.size &&
        ir_type_equals (aa.elem.get(), bb.elem.get()),
      _ => false
    },
    exist_var_t ea => match (*b).data {
      exist_var_t eb => ea == eb,
      _ => false
    },
    forall_t fa => match (*b).data {
      forall_t fb =>
        fa.param.var == fb.param.var &&
        ir_kind_equals (&fa.param.kind, &fb.param.kind) &&
        equals_type_vec (&fa.param.bounds, &fb.param.bounds) &&
        ir_type_equals (fa.body.get(), fb.body.get()),
      _ => false
    },
    fun_t fa => match (*b).data {
      fun_t fb =>
        ir_type_equals (fa.arg.get(), fb.arg.get()) &&
        ir_type_equals (fa.ret.get(), fb.ret.get()),
      _ => false
    },
    lambda_t la => match (*b).data {
      lambda_t lb =>
        la.var == lb.var &&
        ir_kind_equals (&la.kind, &lb.kind) &&
        ir_type_equals (la.body.get(), lb.body.get()),
      _ => false
    },
    mu_t ma => match (*b).data {
      mu_t mb =>
        ma.var == mb.var &&
        ir_kind_equals (&ma.kind, &mb.kind) &&
        ir_type_equals (ma.body.get(), mb.body.get()),
      _ => false
    },
    name_t na => match (*b).data {
      name_t nb => na == nb,
      _ => false
    },
    struct_t sa => match (*b).data {
      struct_t sb => equals_named_type_vec (&sa, &sb),
      _ => false
    },
    tuple_t ta => match (*b).data {
      tuple_t tb => equals_type_vec (&ta, &tb),
      _ => false
    },
    variant_t va => match (*b).data {
      variant_t vb => equals_named_type_vec (&va, &vb),
      _ => false
    },
    var_t va => match (*b).data {
      var_t vb => va == vb,
      _ => false
    }
  }
}

pub fn ir_type_to_string(t: &IrType) -> String {
  match (*t).data {
    app_t a => {
      let s: String = String::from_view "(";
      let fs: String = ir_type_to_string (a.fun.get());
      s.append (fs.view ());
      s.append " ";
      let as_: String = ir_type_to_string (a.arg.get());
      s.append (as_.view ());
      s.append ")";
      s
    },
    array_t a => {
      let s: String = String::from_view "[";
      let es: String = ir_type_to_string (a.elem.get());
      s.append (es.view ());
      s.append ", ";
      let ns: String = i64_to_string a.size;
      s.append (ns.view ());
      s.append "]";
      s
    },
    exist_var_t id => {
      let s: String = String::from_view "^";
      let ns: String = i64_to_string id;
      s.append (ns.view ());
      s
    },
    forall_t f => {
      let s: String = String::from_view "forall ['";
      s.append (f.param.var.view ());
      s.append "] ";
      let bs: String = ir_type_to_string (f.body.get());
      s.append (bs.view ());
      s
    },
    fun_t f => {
      let s: String = String::from_view "";
      let as_: String = ir_type_to_string (f.arg.get());
      s.append (as_.view ());
      s.append " -> ";
      let rs: String = ir_type_to_string (f.ret.get());
      s.append (rs.view ());
      s
    },
    lambda_t l => {
      let s: String = String::from_view "fun (";
      s.append (l.var.view ());
      s.append ") (";
      let bs: String = ir_type_to_string (l.body.get());
      s.append (bs.view ());
      s.append ")";
      s
    },
    mu_t m => {
      let s: String = String::from_view "mu '";
      s.append (m.var.view ());
      s.append ". ";
      let bs: String = ir_type_to_string (m.body.get());
      s.append (bs.view ());
      s
    },
    name_t n => n,
    struct_t fields => {
      let s: String = String::from_view "struct{";
      let i: i64 = 0;
      let n: i64 = fields.size ();
      for i < n {
        if i > 0 {
          s.append ", ";
        }
        let f: NamedType IrType = fields.get i;
        s.append (f.id.view ());
        s.append ": ";
        let fs: String = ir_type_to_string (f.typ.get());
        s.append (fs.view ());
        i <- i + 1;
      }
      s.append "}";
      s
    },
    tuple_t elems => {
      let s: String = String::from_view "(";
      let i: i64 = 0;
      let n: i64 = elems.size ();
      for i < n {
        if i > 0 {
          s.append ", ";
        }
        let e: IrType = elems.get i;
        let es: String = ir_type_to_string (&e);
        s.append (es.view ());
        i <- i + 1;
      }
      s.append ")";
      s
    },
    variant_t tags => {
      let s: String = String::from_view "variant{";
      let i: i64 = 0;
      let n: i64 = tags.size ();
      for i < n {
        if i > 0 {
          s.append ", ";
        }
        let tag: NamedType IrType = tags.get i;
        s.append (tag.id.view ());
        s.append " ";
        let ts: String = ir_type_to_string (tag.typ.get());
        s.append (ts.view ());
        i <- i + 1;
      }
      s.append "}";
      s
    },
    var_t v => {
      let s: String = String::from_view "'";
      s.append (v.view ());
      s
    }
  }
}

pub fn equals_type_vec(a: &Vector IrType, b: &Vector IrType) -> bool {
  let n: i64 = a.size ();
  if n != b.size () {
    return false
  }
  let ok: bool = true;
  let i: i64 = 0;
  for i < n && ok {
    let ea: IrType = a.get i;
    let eb: IrType = b.get i;
    if !(ir_type_equals (&ea, &eb)) {
      ok <- false;
    }
    i <- i + 1;
  }
  ok
}

pub fn equals_named_type_vec(a: &Vector (NamedType IrType), b: &Vector (NamedType IrType)) -> bool {
  let n: i64 = a.size ();
  if n != b.size () {
    return false
  }
  let ok: bool = true;
  let i: i64 = 0;
  for i < n && ok {
    let fa: NamedType IrType = a.get i;
    let fb: NamedType IrType = b.get i;
    if fa.id != fb.id || !(ir_type_equals (fa.typ.get(), fb.typ.get())) {
      ok <- false;
    }
    i <- i + 1;
  }
  ok
}
