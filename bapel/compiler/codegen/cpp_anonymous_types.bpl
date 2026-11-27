implements bapel.compiler.codegen.cpp

imports {
  bapel.core
  bapel.stl
  bapel.compiler.ir
}

pub fn gen_name_type(
  typ: IrType,
  export_flag: bool,
  anonymous_types: &Map String AnonymousType
) -> IrType {
  let bound: Vector String = Vector::mk [String] ();
  let free_set: Set String = Set::mk [String] ();
  typ.collect_free_vars (&bound, &free_set);
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
          mapped <- mapped.substitute_var (&k, &target);
        }
        ki <- ki + 1;
      }
      gen_name_type (mapped, export_flag, anonymous_types)
    },
    _ => mapped
  }
}

pub fn record_anonymous_types_in_alias(
  d: &IrDecl,
  anonymous_types: &Map String AnonymousType
) -> IrDecl {
  match (*d).data {
    alias_decl ad => {
      let lambda_vars: Vector String = Vector::mk [String] ();
      let lambda_kinds: Vector IrKind = Vector::mk [IrKind] ();
      let cur: IrType = ad.typ;
      let peeling: bool = true;
      for peeling {
        match cur.data {
          lambda_t l => {
            lambda_vars.push_back l.var;
            lambda_kinds.push_back l.kind;
            cur <- *(l.body.get ());
            ()
          },
          _ => {
            peeling <- false;
            ()
          }
        }
      }

      let self_type: IrType = IrType::new_name ad.id;
      let ln: i64 = lambda_vars.size ();
      let li: i64 = 0;
      for li < ln {
        let lv: String = lambda_vars.get li;
        self_type <- IrType::new_app (self_type, IrType::new_var lv);
        li <- li + 1;
      }

      let mu_subst: Map String IrType = Map::mk [String, IrType] ();
      let has_mu: bool = false;
      let mu_var: String = String::from_view "";
      let mu_kind: IrKind = IrKind::new_type ();
      match cur.data {
        mu_t m => {
          has_mu <- true;
          mu_var <- m.var;
          mu_kind <- m.kind;
          mu_subst.insert (m.var, self_type);
          cur <- *(m.body.get ());
          ()
        },
        _ => ()
      }

      match cur.data {
        struct_t fields => {
          let new_fields: Vector (NamedType IrType) = Vector::mk [NamedType IrType] ();
          let fn_len: i64 = fields.size ();
          let fi: i64 = 0;
          for fi < fn_len {
            let f: NamedType IrType = fields.get fi;
            let mapped_ft: IrType = record_anonymous_types (f.typ.get (), (*d).export_flag, anonymous_types, &mu_subst);
            new_fields.push_back (struct {
              id = f.id,
              typ = SharedPtr::make_shared [IrType] mapped_ft
            });
            fi <- fi + 1;
          }
          cur <- IrType::new_struct new_fields;
        },
        _ => {
          cur <- record_anonymous_types (&cur, (*d).export_flag, anonymous_types, &mu_subst);
        }
      }

      if has_mu {
        cur <- IrType::new_mu (mu_var, mu_kind, cur);
      }
      let ri: i64 = ln - 1;
      for ri >= 0 {
        let lv: String = lambda_vars.get ri;
        let lk: IrKind = lambda_kinds.get ri;
        cur <- IrType::new_lambda (lv, lk, cur);
        ri <- ri - 1;
      }
      let updated: IrDecl = IrDecl::new_alias (ad.id, ad.kind, cur, (*d).export_flag);
      struct {
        data = updated.data,
        export_flag = (*d).export_flag,
        pos = (*d).pos
      }
    },
    _ => *d
  }
}

pub fn topo_sort_alias_decls(decls: &Vector IrDecl) -> Vector IrDecl {
  let n: i64 = decls.size ();
  let nodes_by_name: Map String i64 = Map::mk [String, i64] ();
  let decl_refs: Vector (Set String) = Vector::mk [Set String] ();
  let adj: Vector (Vector i64) = Vector::mk [Vector i64] ();
  let in_degree: Vector i64 = Vector::mk [i64] ();

  let i: i64 = 0;
  for i < n {
    let d: IrDecl = decls.get i;
    let refs: Set String = Set::mk [String] ();
    match d.data {
      alias_decl ad => {
        nodes_by_name.insert (ad.id, i);
        ad.typ.get_free_names (&refs);
      },
      _ => {
        nodes_by_name.insert (String::from_view "", i);
      }
    }
    decl_refs.push_back refs;
    adj.push_back (Vector::mk [i64] ());
    in_degree.push_back 0;
    i <- i + 1;
  }

  let anonym_pfx: String = String::from_view "__anonym_";
  let j: i64 = 0;
  for j < n {
    let d: IrDecl = decls.get j;
    match d.data {
      alias_decl ad => {
        let is_anonym: bool = ad.id.starts_with (&anonym_pfx);
        let refs_set: Set String = decl_refs.get j;
        let ref_names: Vector String = refs_set.values ();
        let rn: i64 = ref_names.size ();
        let ri: i64 = 0;
        for ri < rn {
          let ref_name: String = ref_names.get ri;
          let opt_idx: Optional i64 = nodes_by_name.get (&ref_name);
          if opt_idx.has_value () {
            let target_idx: i64 = opt_idx.get_value ();
            if target_idx != j {
              let skip: bool = false;
              if is_anonym {
                let target_decl: IrDecl = decls.get target_idx;
                match target_decl.data {
                  alias_decl tad => {
                    if tad.typ.is_mu_alias () {
                      let target_refs: Set String = decl_refs.get target_idx;
                      if target_refs.contains (&(ad.id)) {
                        skip <- true;
                      }
                    }
                    ()
                  },
                  _ => ()
                }
              }
              if !skip {
                let target_adj: Vector i64 = adj.get target_idx;
                target_adj.push_back j;
                VectorImpl::set (&adj, target_idx, target_adj);
                VectorImpl::set (&in_degree, j, in_degree.get j + 1);
              }
            }
          }
          ri <- ri + 1;
        }
      },
      _ => ()
    }
    j <- j + 1;
  }

  let zero_queue: Vector i64 = Vector::mk [i64] ();
  let k: i64 = 0;
  for k < n {
    if in_degree.get k == 0 {
      zero_queue.push_back k;
    }
    k <- k + 1;
  }

  let sorted: Vector IrDecl = Vector::mk [IrDecl] ();
  let q_len: i64 = zero_queue.size ();
  for q_len > 0 {
    q_len <- q_len - 1;
    let u: i64 = zero_queue.get q_len;
    sorted.push_back (decls.get u);
    let neighbors: Vector i64 = adj.get u;
    let nn: i64 = neighbors.size ();
    let ni: i64 = 0;
    for ni < nn {
      let v: i64 = neighbors.get ni;
      let deg: i64 = in_degree.get v - 1;
      VectorImpl::set (&in_degree, v, deg);
      if deg == 0 {
        if q_len < zero_queue.size () {
          VectorImpl::set (&zero_queue, q_len, v);
        } else {
          zero_queue.push_back v;
        }
        q_len <- q_len + 1;
      }
      ni <- ni + 1;
    }
  }

  if sorted.size () != n {
    return *decls;
  }
  sorted
}

pub fn topo_sort_decls(decls: &Vector IrDecl) -> Vector IrDecl {
  let name_decls: Vector IrDecl = Vector::mk [IrDecl] ();
  let alias_decls: Vector IrDecl = Vector::mk [IrDecl] ();
  let trait_decls: Vector IrDecl = Vector::mk [IrDecl] ();
  let term_decls: Vector IrDecl = Vector::mk [IrDecl] ();

  let n: i64 = decls.size ();
  let i: i64 = 0;
  for i < n {
    let d: IrDecl = decls.get i;
    match d.data {
      name_decl _ => name_decls.push_back d,
      alias_decl _ => alias_decls.push_back d,
      trait_decl _ => trait_decls.push_back d,
      term_decl _ => term_decls.push_back d,
      impl_decl _ => ()
    }
    i <- i + 1;
  }

  let sorted_aliases: Vector IrDecl = topo_sort_alias_decls (&alias_decls);
  let out: Vector IrDecl = Vector::mk [IrDecl] ();
  let idx: i64 = 0;
  for idx < name_decls.size () {
    out.push_back (name_decls.get idx);
    idx <- idx + 1;
  }
  idx <- 0;
  for idx < sorted_aliases.size () {
    out.push_back (sorted_aliases.get idx);
    idx <- idx + 1;
  }
  idx <- 0;
  for idx < trait_decls.size () {
    out.push_back (trait_decls.get idx);
    idx <- idx + 1;
  }
  idx <- 0;
  for idx < term_decls.size () {
    out.push_back (term_decls.get idx);
    idx <- idx + 1;
  }
  out
}

pub fn record_anonymous_types_from_module(mod_in: &IrTypedModule) -> IrTypedModule {
  let anonymous_types: Map String AnonymousType = Map::mk [String, AnonymousType] ();
  let empty_mu: Map String IrType = Map::mk [String, IrType] ();

  let new_decls: Vector IrDecl = Vector::mk [IrDecl] ();
  let mod_decls: Vector IrDecl = (*mod_in).decls;
  let dn: i64 = mod_decls.size ();
  let di: i64 = 0;
  for di < dn {
    let d: IrDecl = mod_decls.get di;
    let rewritten: IrDecl = match d.data {
      term_decl td => IrDecl::new_term (
        td.id,
        record_anonymous_types (&(td.typ), d.export_flag, &anonymous_types, &empty_mu),
        d.export_flag,
        d.pos
      ),
      alias_decl _ => record_anonymous_types_in_alias (&d, &anonymous_types),
      trait_decl td => {
        let new_methods: Vector IrSignature = Vector::mk [IrSignature] ();
        let td_methods: Vector IrSignature = td.methods;
        let mn: i64 = td_methods.size ();
        let mi: i64 = 0;
        for mi < mn {
          let m: IrSignature = td_methods.get mi;
          let new_args: Vector FunctionArg = Vector::mk [FunctionArg] ();
          let m_args: Vector FunctionArg = m.args;
          let an: i64 = m_args.size ();
          let ai: i64 = 0;
          for ai < an {
            let a: FunctionArg = m_args.get ai;
            new_args.push_back (struct {
              id = a.id,
              typ = record_anonymous_types (&(a.typ), d.export_flag, &anonymous_types, &empty_mu)
            });
            ai <- ai + 1;
          }
          let new_ret: IrType = record_anonymous_types (&(m.ret_type), d.export_flag, &anonymous_types, &empty_mu);
          new_methods.push_back (struct {
            id = m.id,
            type_params = m.type_params,
            args = new_args,
            ret_type = new_ret
          });
          mi <- mi + 1;
        }
        IrDecl::new_trait (td.id, td.type_params, new_methods, d.export_flag, d.pos)
      },
      _ => d
    };
    new_decls.push_back rewritten;
    di <- di + 1;
  }

  let new_functions: Vector IrFunction = Vector::mk [IrFunction] ();
  let mod_functions: Vector IrFunction = (*mod_in).functions;
  let fn_len: i64 = mod_functions.size ();
  let fi: i64 = 0;
  for fi < fn_len {
    let f: IrFunction = mod_functions.get fi;
    let new_args: Vector FunctionArg = Vector::mk [FunctionArg] ();
    let f_args: Vector FunctionArg = f.args;
    let an: i64 = f_args.size ();
    let ai: i64 = 0;
    for ai < an {
      let a: FunctionArg = f_args.get ai;
      new_args.push_back (struct {
        id = a.id,
        typ = record_anonymous_types (&(a.typ), f.export_flag, &anonymous_types, &empty_mu)
      });
      ai <- ai + 1;
    }
    let new_ret: IrType = record_anonymous_types (&(f.ret_type), f.export_flag, &anonymous_types, &empty_mu);
    new_functions.push_back (struct {
      id = f.id,
      type_params = f.type_params,
      args = new_args,
      ret_type = new_ret,
      body = f.body,
      export_flag = f.export_flag,
      pos = f.pos
    });
    fi <- fi + 1;
  }

  let new_trait_impls: Vector IrTraitImpl = Vector::mk [IrTraitImpl] ();
  let mod_trait_impls: Vector IrTraitImpl = (*mod_in).trait_impls;
  let tin: i64 = mod_trait_impls.size ();
  let tii: i64 = 0;
  for tii < tin {
    let ti: IrTraitImpl = mod_trait_impls.get tii;
    let new_methods: Vector IrFunction = Vector::mk [IrFunction] ();
    let ti_methods: Vector IrFunction = ti.methods;
    let mn: i64 = ti_methods.size ();
    let mi: i64 = 0;
    for mi < mn {
      let m: IrFunction = ti_methods.get mi;
      let new_args: Vector FunctionArg = Vector::mk [FunctionArg] ();
      let m_args: Vector FunctionArg = m.args;
      let an: i64 = m_args.size ();
      let ai: i64 = 0;
      for ai < an {
        let a: FunctionArg = m_args.get ai;
        new_args.push_back (struct {
          id = a.id,
          typ = record_anonymous_types (&(a.typ), m.export_flag, &anonymous_types, &empty_mu)
        });
        ai <- ai + 1;
      }
      let new_ret: IrType = record_anonymous_types (&(m.ret_type), m.export_flag, &anonymous_types, &empty_mu);
      new_methods.push_back (struct {
        id = m.id,
        type_params = m.type_params,
        args = new_args,
        ret_type = new_ret,
        body = m.body,
        export_flag = m.export_flag,
        pos = m.pos
      });
      mi <- mi + 1;
    }
    new_trait_impls.push_back (struct {
      case_val = ti.case_val,
      type_params = ti.type_params,
      trait_type = ti.trait_type,
      type_name = ti.type_name,
      methods = new_methods,
      pos = ti.pos
    });
    tii <- tii + 1;
  }

  let anon_keys: Vector String = anonymous_types.keys ();
  let akn: i64 = anon_keys.size ();
  let aki: i64 = 0;
  for aki < akn {
    let k: String = anon_keys.get aki;
    let opt_a: Optional AnonymousType = anonymous_types.get (&k);
    if opt_a.has_value () {
      let a: AnonymousType = opt_a.get_value ();
      new_decls.push_back a.alias_decl;
    }
    aki <- aki + 1;
  }

  let sorted_decls: Vector IrDecl = topo_sort_decls (&new_decls);
  struct {
    module_id = (*mod_in).module_id,
    mod_imports = (*mod_in).mod_imports,
    mod_impls = (*mod_in).mod_impls,
    mod_flags = (*mod_in).mod_flags,
    import_decls = (*mod_in).import_decls,
    ffi_decls = (*mod_in).ffi_decls,
    decls = sorted_decls,
    trait_impls = new_trait_impls,
    functions = new_functions
  }
}
