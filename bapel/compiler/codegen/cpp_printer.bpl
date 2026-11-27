implements bapel.compiler.codegen.cpp

imports {
  bapel.core
  bapel.stl
  bapel.compiler.ir
}

pub type CppPrinter = struct {
  out: String,
  mode: PrinterMode,
  idgen: i64,
  position: Position,
  auto_type: bool,
  last_term: bool,
  var_destination: String,
  mu_self_names: Map String String,
  in_variant_base: bool,
  mod: &IrTypedModule
}

pub fn is_type_decl(d: &IrDecl) -> bool {
  match (*d).data {
    name_decl _ => true,
    alias_decl _ => true,
    _ => false
  }
}

pub fn hash_type(typ: &IrType) -> String {
  let s: String = ir_type_to_string typ;
  fnv1a_hash_str (s.view ())
}

pub fn canonical_mu_hash(typ: &IrType) -> String {
  match (*typ).data {
    mu_t m => {
      let mu_name: String = String::from_view "__mu";
      let target: IrType = IrType::new_var mu_name;
      let norm_body: IrType = substitute_var_type (m.body.get (), &(m.var), &target);
      let norm: IrType = IrType::new_mu (mu_name, m.kind, norm_body);
      hash_type (&norm)
    },
    _ => hash_type typ
  }
}

pub fn find_decl_in_list(list: &Vector IrDecl, id: &String) -> Optional IrDecl {
  let n: i64 = list.size ();
  let i: i64 = 0;
  let found: Optional IrDecl = Optional::none [IrDecl] ();
  for i < n && !(found.has_value ()) {
    let d: IrDecl = list.get i;
    let did: String = IrDecl::id (&d);
    if did == *id {
      found <- Optional::make_optional [IrDecl] d;
    }
    i <- i + 1;
  }
  found
}

pub fn find_mu_name_in_list(list: &Vector IrDecl, h: &String) -> String {
  let n: i64 = list.size ();
  let i: i64 = 0;
  let res: String = String::from_view "";
  for i < n && res.empty () {
    let d: IrDecl = list.get i;
    match d.data {
      alias_decl ad => {
        match ad.typ.data {
          mu_t _ => {
            if canonical_mu_hash (&(ad.typ)) == *h {
              res <- ad.id;
            }
            ()
          },
          _ => ()
        }
      },
      _ => ()
    }
    i <- i + 1;
  }
  res
}

fn named_ids_equal(a: &Vector (NamedType IrType), b: &Vector (NamedType IrType)) -> bool {
  let na: i64 = a.size ();
  let nb: i64 = b.size ();
  if na != nb {
    return false;
  }
  let i: i64 = 0;
  let ok: bool = true;
  for i < na && ok {
    let fa: NamedType IrType = a.get i;
    let fb: NamedType IrType = b.get i;
    if fa.id != fb.id {
      ok <- false;
    }
    i <- i + 1;
  }
  ok
}

pub fn match_mu_shape_in_list(list: &Vector IrDecl, tbody: &IrType) -> Optional IrDecl {
  let n: i64 = list.size ();
  let i: i64 = 0;
  let found: Optional IrDecl = Optional::none [IrDecl] ();
  for i < n && !(found.has_value ()) {
    let d: IrDecl = list.get i;
    match d.data {
      alias_decl ad => {
        match ad.typ.data {
          mu_t dm => {
            let dbody: &IrType = dm.body.get ();
            match (*tbody).data {
              struct_t tfields => {
                match (*dbody).data {
                  struct_t dfields => {
                    if named_ids_equal (&tfields, &dfields) {
                      found <- Optional::make_optional [IrDecl] d;
                    }
                    ()
                  },
                  _ => ()
                }
              },
              variant_t ttags => {
                match (*dbody).data {
                  variant_t dtags => {
                    if named_ids_equal (&ttags, &dtags) {
                      found <- Optional::make_optional [IrDecl] d;
                    }
                    ()
                  },
                  _ => ()
                }
              },
              _ => ()
            }
          },
          _ => ()
        }
      },
      _ => ()
    }
    i <- i + 1;
  }
  found
}

pub fn find_trait_decl_in_list(list: &Vector IrDecl, id: &String) -> Optional IrDecl {
  let n: i64 = list.size ();
  let i: i64 = 0;
  let found: Optional IrDecl = Optional::none [IrDecl] ();
  for i < n && !(found.has_value ()) {
    let d: IrDecl = list.get i;
    match d.data {
      trait_decl td => {
        if td.id == *id {
          found <- Optional::make_optional [IrDecl] d;
        }
        ()
      },
      _ => ()
    }
    i <- i + 1;
  }
  found
}

impl CppPrinter {
  fn mk(mode: PrinterMode, mod: &IrTypedModule) -> Self {
    struct {
      out = String::from_view "",
      mode = mode,
      idgen = 0,
      position = variant { Position type_position = () },
      auto_type = false,
      last_term = false,
      var_destination = String::from_view "",
      mu_self_names = Map::mk [String, String] (),
      in_variant_base = false,
      mod = mod
    }
  }

  fn write(self: &Self, s: StringView) -> () {
    (*self).out.append s;
    ()
  }

  fn write_str(self: &Self, s: &String) -> () {
    (*self).out.append (s.view ());
    ()
  }

  fn gen_id(self: &Self) -> String {
    let id: i64 = (*self).idgen;
    Ptr::set (self, set (*self) { idgen = id + 1 });
    let res: String = String::from_view "__v_";
    let num: String = i64_to_string id;
    res.append (num.view ());
    res
  }

  fn find_decl(self: &Self, id: &String) -> Optional IrDecl {
    let m: &IrTypedModule = (*self).mod;
    let decls: Vector IrDecl = (*m).decls;
    let r1: Optional IrDecl = find_decl_in_list (&decls, id);
    if r1.has_value () {
      return r1;
    }
    let import_decls: Vector IrDecl = (*m).import_decls;
    let r2: Optional IrDecl = find_decl_in_list (&import_decls, id);
    if r2.has_value () {
      return r2;
    }
    let ffi_decls: Vector IrDecl = (*m).ffi_decls;
    find_decl_in_list (&ffi_decls, id)
  }

  fn normalize_mu_refs(self: &Self, typ: &IrType) -> IrType {
    match (*typ).data {
      mu_t _ => {
        let h: String = canonical_mu_hash typ;
        let m: &IrTypedModule = (*self).mod;
        let decls: Vector IrDecl = (*m).decls;
        let name: String = find_mu_name_in_list (&decls, &h);
        if name.empty () {
          let import_decls: Vector IrDecl = (*m).import_decls;
          name <- find_mu_name_in_list (&import_decls, &h);
        }
        if name.empty () {
          let ffi_decls: Vector IrDecl = (*m).ffi_decls;
          name <- find_mu_name_in_list (&ffi_decls, &h);
        }
        if !(name.empty ()) {
          return IrType::new_name name;
        }
        *typ
      },
      app_t a => IrType::new_app (
        CppPrinter::normalize_mu_refs (self, a.fun.get ()),
        CppPrinter::normalize_mu_refs (self, a.arg.get ())
      ),
      array_t arr => IrType::new_array (
        CppPrinter::normalize_mu_refs (self, arr.elem.get ()),
        arr.size
      ),
      fun_t fn_t => IrType::new_fun (
        CppPrinter::normalize_mu_refs (self, fn_t.arg.get ()),
        CppPrinter::normalize_mu_refs (self, fn_t.ret.get ())
      ),
      struct_t fields => {
        let new_fields: Vector (NamedType IrType) = Vector::mk [NamedType IrType] ();
        let n: i64 = fields.size ();
        let i: i64 = 0;
        for i < n {
          let f: NamedType IrType = fields.get i;
          new_fields.push_back (struct {
            id = f.id,
            typ = SharedPtr::make_shared [IrType] (CppPrinter::normalize_mu_refs (self, f.typ.get ()))
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
          new_elems.push_back (CppPrinter::normalize_mu_refs (self, &e));
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
            typ = SharedPtr::make_shared [IrType] (CppPrinter::normalize_mu_refs (self, tag.typ.get ()))
          });
          i <- i + 1;
        }
        IrType::new_variant new_tags
      },
      _ => *typ
    }
  }

  fn canonical_hash_type(self: &Self, typ: &IrType) -> String {
    match (*typ).data {
      mu_t _ => canonical_mu_hash typ,
      _ => {
        let norm: IrType = CppPrinter::normalize_mu_refs (self, typ);
        hash_type (&norm)
      }
    }
  }

  fn find_decl_for_hash_in_list(self: &Self, list: &Vector IrDecl, target_hash: &String) -> Optional IrDecl {
    let n: i64 = list.size ();
    let i: i64 = 0;
    let found: Optional IrDecl = Optional::none [IrDecl] ();
    for i < n && !(found.has_value ()) {
      let d: IrDecl = list.get i;
      match d.data {
        alias_decl ad => {
          if CppPrinter::canonical_hash_type (self, &(ad.typ)) == *target_hash {
            found <- Optional::make_optional [IrDecl] d;
          }
          ()
        },
        _ => ()
      }
      i <- i + 1;
    }
    found
  }

  fn find_decl_for_type(self: &Self, typ: &IrType) -> Optional IrDecl {
    let m: &IrTypedModule = (*self).mod;
    let target_hash: String = CppPrinter::canonical_hash_type (self, typ);
    let decls: Vector IrDecl = (*m).decls;
    let r1: Optional IrDecl = CppPrinter::find_decl_for_hash_in_list (self, &decls, &target_hash);
    if r1.has_value () { return r1; }
    let import_decls: Vector IrDecl = (*m).import_decls;
    let r2: Optional IrDecl = CppPrinter::find_decl_for_hash_in_list (self, &import_decls, &target_hash);
    if r2.has_value () { return r2; }
    let ffi_decls: Vector IrDecl = (*m).ffi_decls;
    let r3: Optional IrDecl = CppPrinter::find_decl_for_hash_in_list (self, &ffi_decls, &target_hash);
    if r3.has_value () { return r3; }

    match (*typ).data {
      mu_t tm => {
        let tbody: &IrType = tm.body.get ();
        let m1: Optional IrDecl = match_mu_shape_in_list (&decls, tbody);
        if m1.has_value () { return m1; }
        let m2: Optional IrDecl = match_mu_shape_in_list (&import_decls, tbody);
        if m2.has_value () { return m2; }
        match_mu_shape_in_list (&ffi_decls, tbody)
      },
      _ => Optional::none [IrDecl] ()
    }
  }

  fn find_trait_decl(self: &Self, id: &String) -> Optional IrDecl {
    let m: &IrTypedModule = (*self).mod;
    let decls: Vector IrDecl = (*m).decls;
    let r1: Optional IrDecl = find_trait_decl_in_list (&decls, id);
    if r1.has_value () { return r1; }
    let import_decls: Vector IrDecl = (*m).import_decls;
    let r2: Optional IrDecl = find_trait_decl_in_list (&import_decls, id);
    if r2.has_value () { return r2; }
    let ffi_decls: Vector IrDecl = (*m).ffi_decls;
    find_trait_decl_in_list (&ffi_decls, id)
  }

  fn has_function(self: &Self, id: &String) -> bool {
    let m: &IrTypedModule = (*self).mod;
    let fns: Vector IrFunction = (*m).functions;
    let n: i64 = fns.size ();
    let i: i64 = 0;
    let found: bool = false;
    for i < n && !found {
      let f: IrFunction = fns.get i;
      if f.id == *id {
        found <- true;
      }
      i <- i + 1;
    }
    found
  }

  fn to_id(self: &Self, id: &String) -> String {
    let sep: String = String::from_view "::";
    let last_pos: i64 = id.rfind (&sep);
    if last_pos >= 0 {
      let type_pfx: String = id.substr (0, last_pos);
      let method_name: String = id.substr (last_pos + 2, id.size () - (last_pos + 2));
      let opt_d: Optional IrDecl = CppPrinter::find_decl (self, &type_pfx);
      if opt_d.has_value () {
        let d: IrDecl = opt_d.get_value ();
        if is_type_decl (&d) {
          let res: String = String::from_view "::";
          let inh: String = inherent_cpp_name (&type_pfx);
          res.append (inh.view ());
          res.append "::";
          res.append (method_name.view ());
          return res;
        }
      }
      let res: String = String::from_view "::";
      res.append (id.view ());
      return res;
    }
    *id
  }
}
