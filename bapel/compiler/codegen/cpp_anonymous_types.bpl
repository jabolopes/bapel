implements bapel.compiler.codegen.cpp

imports {
  bapel.core
  bapel.stl
  bapel.compiler.ir
}

fn vec_contains_str(v: &Vector String, target: &String) -> bool {
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

pub fn collect_free_vars(typ: &IrType, bound: &Vector String, out: &Set String) -> () {
  match (*typ).data {
    var_t v => {
      if !vec_contains_str (bound, &v) {
        let _: bool = out.insert v;
        ()
      }
      ()
    },
    app_t a => {
      collect_free_vars (a.fun.get (), bound, out);
      collect_free_vars (a.arg.get (), bound, out);
    },
    array_t arr => {
      collect_free_vars (arr.elem.get (), bound, out);
    },
    fun_t fn_t => {
      collect_free_vars (fn_t.arg.get (), bound, out);
      collect_free_vars (fn_t.ret.get (), bound, out);
    },
    forall_t fa => {
      let next_bound: Vector String = *bound;
      next_bound.push_back fa.param.var;
      let bn: i64 = fa.param.bounds.size ();
      let bi: i64 = 0;
      for bi < bn {
        let b: IrType = fa.param.bounds.get bi;
        collect_free_vars (&b, &next_bound, out);
        bi <- bi + 1;
      }
      collect_free_vars (fa.body.get (), &next_bound, out);
    },
    lambda_t l => {
      let next_bound: Vector String = *bound;
      next_bound.push_back l.var;
      collect_free_vars (l.body.get (), &next_bound, out);
    },
    mu_t m => {
      let next_bound: Vector String = *bound;
      next_bound.push_back m.var;
      collect_free_vars (m.body.get (), &next_bound, out);
    },
    struct_t fields => {
      let n: i64 = fields.size ();
      let i: i64 = 0;
      for i < n {
        let f: NamedType IrType = fields.get i;
        collect_free_vars (f.typ.get (), bound, out);
        i <- i + 1;
      }
    },
    tuple_t elems => {
      let n: i64 = elems.size ();
      let i: i64 = 0;
      for i < n {
        let e: IrType = elems.get i;
        collect_free_vars (&e, bound, out);
        i <- i + 1;
      }
    },
    variant_t tags => {
      let n: i64 = tags.size ();
      let i: i64 = 0;
      for i < n {
        let tag: NamedType IrType = tags.get i;
        collect_free_vars (tag.typ.get (), bound, out);
        i <- i + 1;
      }
    },
    _ => ()
  }
}

pub fn substitute_var_type(typ: &IrType, var_name: &String, target: &IrType) -> IrType {
  match (*typ).data {
    var_t v => {
      if v == *var_name {
        return *target;
      }
      *typ
    },
    app_t a => IrType::new_app (
      substitute_var_type (a.fun.get (), var_name, target),
      substitute_var_type (a.arg.get (), var_name, target)
    ),
    array_t arr => IrType::new_array (
      substitute_var_type (arr.elem.get (), var_name, target),
      arr.size
    ),
    fun_t fn_t => IrType::new_fun (
      substitute_var_type (fn_t.arg.get (), var_name, target),
      substitute_var_type (fn_t.ret.get (), var_name, target)
    ),
    forall_t fa => {
      if fa.param.var == *var_name {
        return *typ;
      }
      let new_bounds: Vector IrType = Vector::mk [IrType] ();
      let bn: i64 = fa.param.bounds.size ();
      let bi: i64 = 0;
      for bi < bn {
        let b: IrType = fa.param.bounds.get bi;
        new_bounds.push_back (substitute_var_type (&b, var_name, target));
        bi <- bi + 1;
      }
      let new_param: TypeParam IrType = struct {
        var = fa.param.var,
        kind = fa.param.kind,
        bounds = new_bounds
      };
      IrType::new_forall (new_param, substitute_var_type (fa.body.get (), var_name, target))
    },
    lambda_t l => {
      if l.var == *var_name {
        return *typ;
      }
      IrType::new_lambda (l.var, l.kind, substitute_var_type (l.body.get (), var_name, target))
    },
    mu_t m => {
      if m.var == *var_name {
        return *typ;
      }
      IrType::new_mu (m.var, m.kind, substitute_var_type (m.body.get (), var_name, target))
    },
    struct_t fields => {
      let new_fields: Vector (NamedType IrType) = Vector::mk [NamedType IrType] ();
      let n: i64 = fields.size ();
      let i: i64 = 0;
      for i < n {
        let f: NamedType IrType = fields.get i;
        new_fields.push_back (struct {
          id = f.id,
          typ = SharedPtr::make_shared [IrType] (substitute_var_type (f.typ.get (), var_name, target))
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
        new_elems.push_back (substitute_var_type (&e, var_name, target));
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
          typ = SharedPtr::make_shared [IrType] (substitute_var_type (tag.typ.get (), var_name, target))
        });
        i <- i + 1;
      }
      IrType::new_variant new_tags
    },
    _ => *typ
  }
}

pub fn gen_name_type(
  typ: IrType,
  export_flag: bool,
  anonymous_types: &Map String AnonymousType
) -> IrType {
  let bound: Vector String = Vector::mk [String] ();
  let free_set: Set String = Set::mk [String] ();
  collect_free_vars (&typ, &bound, &free_set);
  let free_vars: Vector String = free_set.values ();

  let decl_type: IrType = typ;
  let decl_kind: IrKind = IrKind::new_type ();
  let n: i64 = free_vars.size ();
  let i: i64 = n - 1;
  for i >= 0 {
    let fv: String = free_vars.get i;
    decl_type <- IrType::new_lambda (fv, IrKind::new_type (), decl_type);
    decl_kind <- IrKind::new_arrow (IrKind::new_type (), decl_kind);
    i <- i - 1;
  }

  let type_str: String = ir_type_to_string (&decl_type);
  let hash: String = fnv1a_hash_str (type_str.view ());
  let name: String = String::from_view "__anonym_";
  name.append (hash.view ());

  let existing: Optional AnonymousType = anonymous_types.get (&hash);
  if existing.has_value () {
    let cur: AnonymousType = existing.get_value ();
    if export_flag && !cur.alias_decl.export_flag {
      let updated_decl: IrDecl = struct {
        data = cur.alias_decl.data,
        export_flag = true,
        pos = cur.alias_decl.pos
      };
      anonymous_types.insert (hash, struct {
        name_type = cur.name_type,
        alias_decl = updated_decl
      });
    }
  } else {
    let new_decl: IrDecl = IrDecl::new_alias (name, decl_kind, decl_type, export_flag);
    anonymous_types.insert (hash, struct {
      name_type = IrType::new_name name,
      alias_decl = new_decl
    });
  }

  let ref_t: IrType = IrType::new_name name;
  let j: i64 = 0;
  for j < n {
    let fv: String = free_vars.get j;
    ref_t <- IrType::new_app (ref_t, IrType::new_var fv);
    j <- j + 1;
  }
  ref_t
}

pub fn record_anonymous_types(
  typ: &IrType,
  export_flag: bool,
  anonymous_types: &Map String AnonymousType,
  mu_subst: &Map String IrType
) -> IrType {
  let mapped: IrType = match (*typ).data {
    app_t a => IrType::new_app (
      record_anonymous_types (a.fun.get (), export_flag, anonymous_types, mu_subst),
      record_anonymous_types (a.arg.get (), export_flag, anonymous_types, mu_subst)
    ),
    array_t arr => IrType::new_array (
      record_anonymous_types (arr.elem.get (), export_flag, anonymous_types, mu_subst),
      arr.size
    ),
    fun_t fn_t => IrType::new_fun (
      record_anonymous_types (fn_t.arg.get (), export_flag, anonymous_types, mu_subst),
      record_anonymous_types (fn_t.ret.get (), export_flag, anonymous_types, mu_subst)
    ),
    forall_t fa => {
      let new_bounds: Vector IrType = Vector::mk [IrType] ();
      let bn: i64 = fa.param.bounds.size ();
      let bi: i64 = 0;
      for bi < bn {
        let b: IrType = fa.param.bounds.get bi;
        new_bounds.push_back (record_anonymous_types (&b, export_flag, anonymous_types, mu_subst));
        bi <- bi + 1;
      }
      let new_param: TypeParam IrType = struct {
        var = fa.param.var,
        kind = fa.param.kind,
        bounds = new_bounds
      };
      IrType::new_forall (new_param, record_anonymous_types (fa.body.get (), export_flag, anonymous_types, mu_subst))
    },
    lambda_t l => IrType::new_lambda (
      l.var,
      l.kind,
      record_anonymous_types (l.body.get (), export_flag, anonymous_types, mu_subst)
    ),
    mu_t m => IrType::new_mu (
      m.var,
      m.kind,
      record_anonymous_types (m.body.get (), export_flag, anonymous_types, mu_subst)
    ),
    struct_t fields => {
      let new_fields: Vector (NamedType IrType) = Vector::mk [NamedType IrType] ();
      let n: i64 = fields.size ();
      let i: i64 = 0;
      for i < n {
        let f: NamedType IrType = fields.get i;
        new_fields.push_back (struct {
          id = f.id,
          typ = SharedPtr::make_shared [IrType] (record_anonymous_types (f.typ.get (), export_flag, anonymous_types, mu_subst))
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
        new_elems.push_back (record_anonymous_types (&e, export_flag, anonymous_types, mu_subst));
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
          typ = SharedPtr::make_shared [IrType] (record_anonymous_types (tag.typ.get (), export_flag, anonymous_types, mu_subst))
        });
        i <- i + 1;
      }
      IrType::new_variant new_tags
    },
    _ => *typ
  };

  match mapped.data {
    struct_t _ => {
      let keys: Vector String = mu_subst.keys ();
      let kn: i64 = keys.size ();
      let ki: i64 = 0;
      for ki < kn {
        let k: String = keys.get ki;
        let opt_t: Optional IrType = mu_subst.get (&k);
        if opt_t.has_value () {
          let target: IrType = opt_t.get_value ();
          mapped <- substitute_var_type (&mapped, &k, &target);
        }
        ki <- ki + 1;
      }
      gen_name_type (mapped, export_flag, anonymous_types)
    },
    _ => mapped
  }
}

pub fn is_mu_alias_type(t: &IrType) -> bool {
  match (*t).data {
    lambda_t l => is_mu_alias_type (l.body.get ()),
    mu_t _ => true,
    _ => false
  }
}

pub fn get_free_names_from_type(typ: &IrType, out: &Set String) -> () {
  match (*typ).data {
    name_t n => {
      let _: bool = out.insert n;
      ()
    },
    app_t a => {
      get_free_names_from_type (a.fun.get (), out);
      get_free_names_from_type (a.arg.get (), out);
    },
    array_t arr => {
      get_free_names_from_type (arr.elem.get (), out);
    },
    fun_t fn_t => {
      get_free_names_from_type (fn_t.arg.get (), out);
      get_free_names_from_type (fn_t.ret.get (), out);
    },
    struct_t fields => {
      let n: i64 = fields.size ();
      let i: i64 = 0;
      for i < n {
        let f: NamedType IrType = fields.get i;
        get_free_names_from_type (f.typ.get (), out);
        i <- i + 1;
      }
    },
    tuple_t elems => {
      let n: i64 = elems.size ();
      let i: i64 = 0;
      for i < n {
        let e: IrType = elems.get i;
        get_free_names_from_type (&e, out);
        i <- i + 1;
      }
    },
    variant_t tags => {
      let n: i64 = tags.size ();
      let i: i64 = 0;
      for i < n {
        let tag: NamedType IrType = tags.get i;
        get_free_names_from_type (tag.typ.get (), out);
        i <- i + 1;
      }
    },
    lambda_t l => {
      get_free_names_from_type (l.body.get (), out);
    },
    mu_t m => {
      get_free_names_from_type (m.body.get (), out);
    },
    _ => ()
  }
}
