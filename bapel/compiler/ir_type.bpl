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

  fn app_fun(self: &Self) -> IrType {
    match (*self).data {
      app_t a => IrType::app_fun (a.fun.get ()),
      _ => *self
    }
  }

  fn app_args(self: &Self) -> Vector IrType {
    match (*self).data {
      app_t a => {
        let args: Vector IrType = IrType::app_args (a.fun.get ());
        args.push_back (*(a.arg.get ()));
        args
      },
      _ => Vector::mk [IrType] ()
    }
  }

  fn base_name(self: &Self) -> String {
    match (*self).data {
      name_t n => n,
      app_t _ => {
        let fun: IrType = IrType::app_fun self;
        IrType::base_name (&fun)
      },
      _ => String::from_view ""
    }
  }

  fn forall_vars(self: &Self) -> Vector String {
    let out: Vector String = Vector::mk [String] ();
    let cur: IrType = *self;
    let running: bool = true;
    for running {
      match cur.data {
        forall_t fa => {
          out.push_back fa.param.var;
          cur <- *(fa.body.get ());
          ()
        },
        _ => {
          running <- false;
          ()
        }
      }
    }
    out
  }

  fn forall_type_params(self: &Self) -> Vector (TypeParam IrType) {
    let out: Vector (TypeParam IrType) = Vector::mk [TypeParam IrType] ();
    let cur: IrType = *self;
    let running: bool = true;
    for running {
      match cur.data {
        forall_t fa => {
          out.push_back fa.param;
          cur <- *(fa.body.get ());
          ()
        },
        _ => {
          running <- false;
          ()
        }
      }
    }
    out
  }

  fn forall_body(self: &Self) -> IrType {
    match (*self).data {
      forall_t fa => IrType::forall_body (fa.body.get ()),
      _ => *self
    }
  }

  fn lambda_vars(self: &Self) -> Vector String {
    let out: Vector String = Vector::mk [String] ();
    let cur: IrType = *self;
    let running: bool = true;
    for running {
      match cur.data {
        lambda_t l => {
          out.push_back l.var;
          cur <- *(l.body.get ());
          ()
        },
        _ => {
          running <- false;
          ()
        }
      }
    }
    out
  }

  fn lambda_body(self: &Self) -> IrType {
    match (*self).data {
      lambda_t l => IrType::lambda_body (l.body.get ()),
      _ => *self
    }
  }

  fn count_alias_type_vars(self: &Self) -> i64 {
    match (*self).data {
      lambda_t l => 1 + IrType::count_alias_type_vars (l.body.get ()),
      _ => 0
    }
  }

  fn is_mu_alias(self: &Self) -> bool {
    match (*self).data {
      lambda_t l => IrType::is_mu_alias (l.body.get ()),
      mu_t _ => true,
      _ => false
    }
  }

  fn collect_free_vars(self: &Self, bound: &Vector String, out: &Set String) -> () {
    match (*self).data {
      var_t v => {
        if !ir_vec_contains_str (bound, &v) {
          let _: bool = out.insert v;
          ()
        }
        ()
      },
      app_t a => {
        IrType::collect_free_vars (a.fun.get (), bound, out);
        IrType::collect_free_vars (a.arg.get (), bound, out);
      },
      array_t arr => {
        IrType::collect_free_vars (arr.elem.get (), bound, out);
      },
      fun_t fn_t => {
        IrType::collect_free_vars (fn_t.arg.get (), bound, out);
        IrType::collect_free_vars (fn_t.ret.get (), bound, out);
      },
      forall_t fa => {
        let next_bound: Vector String = *bound;
        next_bound.push_back fa.param.var;
        let bn: i64 = fa.param.bounds.size ();
        let bi: i64 = 0;
        for bi < bn {
          let b: IrType = fa.param.bounds.get bi;
          IrType::collect_free_vars (&b, &next_bound, out);
          bi <- bi + 1;
        }
        IrType::collect_free_vars (fa.body.get (), &next_bound, out);
      },
      lambda_t l => {
        let next_bound: Vector String = *bound;
        next_bound.push_back l.var;
        IrType::collect_free_vars (l.body.get (), &next_bound, out);
      },
      mu_t m => {
        let next_bound: Vector String = *bound;
        next_bound.push_back m.var;
        IrType::collect_free_vars (m.body.get (), &next_bound, out);
      },
      struct_t fields => {
        let n: i64 = fields.size ();
        let i: i64 = 0;
        for i < n {
          let f: NamedType IrType = fields.get i;
          IrType::collect_free_vars (f.typ.get (), bound, out);
          i <- i + 1;
        }
      },
      tuple_t elems => {
        let n: i64 = elems.size ();
        let i: i64 = 0;
        for i < n {
          let e: IrType = elems.get i;
          IrType::collect_free_vars (&e, bound, out);
          i <- i + 1;
        }
      },
      variant_t tags => {
        let n: i64 = tags.size ();
        let i: i64 = 0;
        for i < n {
          let tag: NamedType IrType = tags.get i;
          IrType::collect_free_vars (tag.typ.get (), bound, out);
          i <- i + 1;
        }
      },
      _ => ()
    }
    ()
  }

  fn substitute_var(self: &Self, var_name: &String, target: &IrType) -> IrType {
    match (*self).data {
      var_t v => {
        if v == *var_name {
          return *target;
        }
        *self
      },
      app_t a => IrType::new_app (
        IrType::substitute_var (a.fun.get (), var_name, target),
        IrType::substitute_var (a.arg.get (), var_name, target)
      ),
      array_t arr => IrType::new_array (
        IrType::substitute_var (arr.elem.get (), var_name, target),
        arr.size
      ),
      fun_t fn_t => IrType::new_fun (
        IrType::substitute_var (fn_t.arg.get (), var_name, target),
        IrType::substitute_var (fn_t.ret.get (), var_name, target)
      ),
      forall_t fa => {
        if fa.param.var == *var_name {
          return *self;
        }
        let new_bounds: Vector IrType = Vector::mk [IrType] ();
        let bn: i64 = fa.param.bounds.size ();
        let bi: i64 = 0;
        for bi < bn {
          let b: IrType = fa.param.bounds.get bi;
          new_bounds.push_back (IrType::substitute_var (&b, var_name, target));
          bi <- bi + 1;
        }
        let new_param: TypeParam IrType = struct {
          var = fa.param.var,
          kind = fa.param.kind,
          bounds = new_bounds
        };
        IrType::new_forall (new_param, IrType::substitute_var (fa.body.get (), var_name, target))
      },
      lambda_t l => {
        if l.var == *var_name {
          return *self;
        }
        IrType::new_lambda (l.var, l.kind, IrType::substitute_var (l.body.get (), var_name, target))
      },
      mu_t m => {
        if m.var == *var_name {
          return *self;
        }
        IrType::new_mu (m.var, m.kind, IrType::substitute_var (m.body.get (), var_name, target))
      },
      struct_t fields => {
        let new_fields: Vector (NamedType IrType) = Vector::mk [NamedType IrType] ();
        let n: i64 = fields.size ();
        let i: i64 = 0;
        for i < n {
          let f: NamedType IrType = fields.get i;
          new_fields.push_back (struct {
            id = f.id,
            typ = SharedPtr::make_shared [IrType] (IrType::substitute_var (f.typ.get (), var_name, target))
          });
          i <- i + 1;
        }
        IrType::new_struct new_fields
      },
      tuple_t elems => {
        let new_elems: Vector IrType = Vector::mk [IrType] ();
        let n: i64 = elems.size ();
        let i: i64 = 0;
        for i < n {
          let e: IrType = elems.get i;
          new_elems.push_back (IrType::substitute_var (&e, var_name, target));
          i <- i + 1;
        }
        IrType::new_tuple new_elems
      },
      variant_t tags => {
        let new_tags: Vector (NamedType IrType) = Vector::mk [NamedType IrType] ();
        let n: i64 = tags.size ();
        let i: i64 = 0;
        for i < n {
          let tag: NamedType IrType = tags.get i;
          new_tags.push_back (struct {
            id = tag.id,
            typ = SharedPtr::make_shared [IrType] (IrType::substitute_var (tag.typ.get (), var_name, target))
          });
          i <- i + 1;
        }
        IrType::new_variant new_tags
      },
      _ => *self
    }
  }

  fn get_free_names(self: &Self, out: &Set String) -> () {
    match (*self).data {
      name_t n => {
        let _: bool = out.insert n;
        ()
      },
      app_t a => {
        IrType::get_free_names (a.fun.get (), out);
        IrType::get_free_names (a.arg.get (), out);
      },
      array_t arr => {
        IrType::get_free_names (arr.elem.get (), out);
      },
      fun_t fn_t => {
        IrType::get_free_names (fn_t.arg.get (), out);
        IrType::get_free_names (fn_t.ret.get (), out);
      },
      struct_t fields => {
        let n: i64 = fields.size ();
        let i: i64 = 0;
        for i < n {
          let f: NamedType IrType = fields.get i;
          IrType::get_free_names (f.typ.get (), out);
          i <- i + 1;
        }
      },
      tuple_t elems => {
        let n: i64 = elems.size ();
        let i: i64 = 0;
        for i < n {
          let e: IrType = elems.get i;
          IrType::get_free_names (&e, out);
          i <- i + 1;
        }
      },
      variant_t tags => {
        let n: i64 = tags.size ();
        let i: i64 = 0;
        for i < n {
          let tag: NamedType IrType = tags.get i;
          IrType::get_free_names (tag.typ.get (), out);
          i <- i + 1;
        }
      },
      lambda_t l => {
        IrType::get_free_names (l.body.get (), out);
      },
      mu_t m => {
        IrType::get_free_names (m.body.get (), out);
      },
      _ => ()
    }
    ()
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

fn arg_needs_parens(t: &IrType) -> bool {
  match (*t).data {
    app_t _ => true,
    forall_t _ => true,
    fun_t _ => true,
    lambda_t _ => true,
    _ => false
  }
}

pub fn ir_type_to_string(t: &IrType) -> String {
  match (*t).data {
    app_t a => {
      let s: String = ir_type_to_string (a.fun.get ());
      s.append " ";
      let arg_ptr: &IrType = a.arg.get ();
      let needs_parens: bool = arg_needs_parens arg_ptr;
      if needs_parens {
        s.append "(";
      }
      let as_: String = ir_type_to_string arg_ptr;
      s.append (as_.view ());
      if needs_parens {
        s.append ")";
      }
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

pub fn count_type_vars(kind: &IrKind) -> i64 {
  match *kind {
    arrow_kind ak => 1 + count_type_vars (ak.right.get ()),
    type_kind _ => 0
  }
}

pub fn ir_vec_contains_str(v: &Vector String, target: &String) -> bool {
  let n: i64 = v.size ();
  let i: i64 = 0;
  let found: bool = false;
  for i < n && !found {
    let s: String = v.get i;
    if s == *target {
      found <- true;
    }
    i <- i + 1;
  }
  found
}
