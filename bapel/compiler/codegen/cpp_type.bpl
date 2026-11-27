implements bapel.compiler.codegen.cpp

imports {
  bapel.core
  bapel.stl
  bapel.compiler.ir
}

pub fn match_type(
  pat: &IrType,
  target: &IrType,
  vars: &Set String,
  subs: &Map String IrType
) -> bool {
  match (*pat).data {
    var_t v => {
      if vars.contains (&v) {
        let opt_existing: Optional IrType = subs.get (&v);
        if opt_existing.has_value () {
          let ex: IrType = opt_existing.get_value ();
          return ir_type_equals (&ex, target);
        }
        subs.insert (v, *target);
        return true;
      }
      ir_type_equals (pat, target)
    },
    app_t pa => {
      match (*target).data {
        app_t ta =>
          match_type (pa.fun.get (), ta.fun.get (), vars, subs) &&
          match_type (pa.arg.get (), ta.arg.get (), vars, subs),
        _ => false
      }
    },
    array_t pa => {
      match (*target).data {
        array_t ta =>
          pa.size == ta.size &&
          match_type (pa.elem.get (), ta.elem.get (), vars, subs),
        _ => false
      }
    },
    fun_t pf => {
      match (*target).data {
        fun_t tf =>
          match_type (pf.arg.get (), tf.arg.get (), vars, subs) &&
          match_type (pf.ret.get (), tf.ret.get (), vars, subs),
        _ => false
      }
    },
    struct_t pfields => {
      match (*target).data {
        struct_t tfields => {
          let np: i64 = pfields.size ();
          let nt: i64 = tfields.size ();
          if np != nt { return false; }
          let i: i64 = 0;
          let ok: bool = true;
          for i < np && ok {
            let pf: NamedType IrType = pfields.get i;
            let tf: NamedType IrType = tfields.get i;
            if pf.id != tf.id || !match_type (pf.typ.get (), tf.typ.get (), vars, subs) {
              ok <- false;
            }
            i <- i + 1;
          }
          ok
        },
        _ => false
      }
    },
    tuple_t pelems => {
      match (*target).data {
        tuple_t telems => {
          let np: i64 = pelems.size ();
          let nt: i64 = telems.size ();
          if np != nt { return false; }
          let i: i64 = 0;
          let ok: bool = true;
          for i < np && ok {
            let pe: IrType = pelems.get i;
            let te: IrType = telems.get i;
            if !match_type (&pe, &te, vars, subs) {
              ok <- false;
            }
            i <- i + 1;
          }
          ok
        },
        _ => false
      }
    },
    variant_t ptags => {
      match (*target).data {
        variant_t ttags => {
          let np: i64 = ptags.size ();
          let nt: i64 = ttags.size ();
          if np != nt { return false; }
          let i: i64 = 0;
          let ok: bool = true;
          for i < np && ok {
            let pt: NamedType IrType = ptags.get i;
            let tt: NamedType IrType = ttags.get i;
            if pt.id != tt.id || !match_type (pt.typ.get (), tt.typ.get (), vars, subs) {
              ok <- false;
            }
            i <- i + 1;
          }
          ok
        },
        _ => false
      }
    },
    _ => ir_type_equals (pat, target)
  }
}

fn try_print_generic_struct_in_list(
  p: &CppPrinter,
  list: &Vector IrDecl,
  typ: &IrType
) -> bool {
  let n: i64 = list.size ();
  let i: i64 = 0;
  let matched: bool = false;
  for i < n && !matched {
    let d: IrDecl = list.get i;
    match d.data {
      alias_decl ad => {
        match ad.typ.data {
          lambda_t _ => {
            let body: IrType = ad.typ.lambda_body ();
            match body.data {
              struct_t _ => {
                let lvars: Vector String = ad.typ.lambda_vars ();
                let var_set: Set String = Set::mk [String] ();
                let vn: i64 = lvars.size ();
                let vi: i64 = 0;
                for vi < vn {
                  let _: bool = var_set.insert (lvars.get vi);
                  vi <- vi + 1;
                }
                let subs: Map String IrType = Map::mk [String, IrType] ();
                if match_type (&body, typ, &var_set, &subs) && subs.size () == vn {
                  p.write "::";
                  let did: String = IrDecl::id (&d);
                  let cid: String = p.to_id (&did);
                  p.write_str (&cid);
                  p.write "<";
                  let j: i64 = 0;
                  for j < vn {
                    if j > 0 { p.write ", "; }
                    let v: String = lvars.get j;
                    let opt_sub: Optional IrType = subs.get (&v);
                    if opt_sub.has_value () {
                      let st: IrType = opt_sub.get_value ();
                      print_type (p, &st);
                    }
                    j <- j + 1;
                  }
                  p.write ">";
                  matched <- true;
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
    i <- i + 1;
  }
  matched
}

pub fn print_type(p: &CppPrinter, typ: &IrType) -> () {
  if (*p).auto_type {
    p.write "auto";
    return ();
  }

  match (*typ).data {
    app_t _ => {
      let fun: IrType = typ.app_fun ();
      print_type (p, &fun);
      let args: Vector IrType = typ.app_args ();
      p.write "<";
      let n: i64 = args.size ();
      let i: i64 = 0;
      for i < n {
        if i > 0 { p.write ", "; }
        let a: IrType = args.get i;
        print_type (p, &a);
        i <- i + 1;
      }
      p.write ">";
    },
    array_t arr => {
      p.write "std::array<";
      print_type (p, arr.elem.get ());
      p.write ", ";
      let sz: String = i64_to_string arr.size;
      p.write_str (&sz);
      p.write ">";
    },
    exist_var_t _ => {
      let s: String = ir_type_to_string typ;
      p.write_str (&s);
    },
    forall_t _ => {
      let tvars: Vector String = typ.forall_vars ();
      p.write "template <";
      let n: i64 = tvars.size ();
      let i: i64 = 0;
      for i < n {
        if i > 0 { p.write ", "; }
        p.write "typename ";
        let tv: String = tvars.get i;
        p.write_str (&tv);
        i <- i + 1;
      }
      p.write "> ";
      let body: IrType = typ.forall_body ();
      print_type (p, &body);
    },
    fun_t fn_t => {
      p.write "std::function<";
      let old_pos: Position = (*p).position;
      Ptr::set (p, set (*p) { position = variant { Position bind_position = () } });
      print_type (p, fn_t.ret.get ());
      Ptr::set (p, set (*p) { position = old_pos });
      p.write "(";
      print_type (p, fn_t.arg.get ());
      p.write ")>";
    },
    name_t n => {
      let nv: StringView = n.view ();
      let anonym_pfx: String = String::from_view "__anonym_";
      let sep: String = String::from_view "::";
      if nv == "bool" {
        p.write "bool";
      } else if nv == "i8" {
        p.write "int8_t";
      } else if nv == "i16" {
        p.write "int16_t";
      } else if nv == "i32" {
        p.write "int32_t";
      } else if nv == "i64" {
        p.write "int64_t";
      } else if nv == "f32" {
        p.write "float";
      } else if nv == "f64" {
        p.write "double";
      } else if n.starts_with (&anonym_pfx) {
        p.write_str (&n);
      } else {
        let opt_d: Optional IrDecl = p.find_decl (&n);
        if n.find (&sep, 0) < 0 && opt_d.has_value () {
          p.write "::";
        }
        let cid: String = p.to_id (&n);
        p.write_str (&cid);
      }
    },
    struct_t fields => {
      let opt_d: Optional IrDecl = p.find_decl_for_type typ;
      if opt_d.has_value () {
        let d: IrDecl = opt_d.get_value ();
        let did: String = IrDecl::id (&d);
        let anonym_pfx: String = String::from_view "__anonym_";
        if did.starts_with (&anonym_pfx) {
          p.write_str (&did);
        } else {
          p.write "::";
          let cid: String = p.to_id (&did);
          p.write_str (&cid);
        }
        return ();
      }
      let m: &IrTypedModule = (*p).mod;
      let decls: Vector IrDecl = (*m).decls;
      if try_print_generic_struct_in_list (p, &decls, typ) { return (); }
      let import_decls: Vector IrDecl = (*m).import_decls;
      if try_print_generic_struct_in_list (p, &import_decls, typ) { return (); }
      let ffi_decls: Vector IrDecl = (*m).ffi_decls;
      if try_print_generic_struct_in_list (p, &ffi_decls, typ) { return (); }

      p.write "struct {";
      let fn_len: i64 = fields.size ();
      let fi: i64 = 0;
      for fi < fn_len {
        if fi > 0 { p.write " "; }
        let f: NamedType IrType = fields.get fi;
        print_type (p, f.typ.get ());
        p.write " ";
        p.write_str (&(f.id));
        p.write ";";
        fi <- fi + 1;
      }
      p.write "}";
    },
    tuple_t elems => {
      let en: i64 = elems.size ();
      match (*p).position {
        type_position _ => {
          let i: i64 = 0;
          for i < en {
            if i > 0 { p.write ", "; }
            let e: IrType = elems.get i;
            print_type (p, &e);
            i <- i + 1;
          }
        },
        bind_position _ => {
          if en == 0 {
            p.write "std::monostate";
          } else if en == 1 {
            let e0: IrType = elems.get 0;
            print_type (p, &e0);
          } else {
            p.write "std::tuple<";
            let i: i64 = 0;
            for i < en {
              if i > 0 { p.write ", "; }
              let e: IrType = elems.get i;
              print_type (p, &e);
              i <- i + 1;
            }
            p.write ">";
          }
        }
      }
    },
    variant_t tags => {
      if !((*p).in_variant_base) {
        let opt_d: Optional IrDecl = p.find_decl_for_type typ;
        if opt_d.has_value () {
          let d: IrDecl = opt_d.get_value ();
          let did: String = IrDecl::id (&d);
          p.write "::";
          let cid: String = p.to_id (&did);
          p.write_str (&cid);
          return ();
        }
      }
      let saved: bool = (*p).in_variant_base;
      Ptr::set (p, set (*p) { in_variant_base = false });
      p.write "std::variant<";
      let tn: i64 = tags.size ();
      let ti: i64 = 0;
      for ti < tn {
        if ti > 0 { p.write ", "; }
        let tag: NamedType IrType = tags.get ti;
        let old_pos: Position = (*p).position;
        Ptr::set (p, set (*p) { position = variant { Position bind_position = () } });
        print_type (p, tag.typ.get ());
        Ptr::set (p, set (*p) { position = old_pos });
        p.write "/* ";
        let cid: String = p.to_id (&(tag.id));
        p.write_str (&cid);
        p.write " */";
        ti <- ti + 1;
      }
      p.write ">";
      Ptr::set (p, set (*p) { in_variant_base = saved });
    },
    mu_t m => {
      let opt_d: Optional IrDecl = p.find_decl_for_type typ;
      if opt_d.has_value () {
        let d: IrDecl = opt_d.get_value ();
        let did: String = IrDecl::id (&d);
        p.write "::";
        let cid: String = p.to_id (&did);
        p.write_str (&cid);
        return ();
      }
      print_type (p, m.body.get ());
    },
    var_t v => {
      let self_names: Map String String = (*p).mu_self_names;
      let opt_s: Optional String = self_names.get (&v);
      if opt_s.has_value () {
        let s: String = opt_s.get_value ();
        p.write_str (&s);
      } else {
        p.write_str (&v);
      }
    },
    _ => ()
  }
  ()
}

pub fn sfinae_constraint(p: &CppPrinter, vk_var: &String, bound: &IrType) -> String {
  let trait_name: String = bound.base_name ();
  let args: Vector IrType = bound.app_args ();
  let cpp_name: String = String::from_view "::";
  let tname: String = trait_cpp_name (&trait_name);
  cpp_name.append (tname.view ());

  let out: String = String::from_view "(sizeof(";
  out.append (cpp_name.view ());
  out.append "<";
  out.append (vk_var.view ());
  let an: i64 = args.size ();
  let ai: i64 = 0;
  for ai < an {
    out.append ", ";
    let temp: CppPrinter = CppPrinter::mk ((*p).mode, (*p).mod);
    let a: IrType = args.get ai;
    print_type (&temp, &a);
    let temp_str: &String = temp.out.get ();
    out.append (temp_str.view ());
    ai <- ai + 1;
  }
  out.append ">) > 0)";
  out
}

pub fn print_template_params(
  p: &CppPrinter,
  tvars: &Vector (TypeParam IrType),
  is_definition: bool
) -> () {
  let n: i64 = tvars.size ();
  if n == 0 { return (); }
  p.write "template <";
  let i: i64 = 0;
  for i < n {
    if i > 0 { p.write ", "; }
    let tp: TypeParam IrType = tvars.get i;
    p.write "typename ";
    p.write_str (&(tp.var));
    i <- i + 1;
  }

  let constraints: Vector String = Vector::mk [String] ();
  let j: i64 = 0;
  for j < n {
    let tp: TypeParam IrType = tvars.get j;
    let bounds: Vector IrType = tp.bounds;
    let bn: i64 = bounds.size ();
    let bi: i64 = 0;
    for bi < bn {
      let b: IrType = bounds.get bi;
      constraints.push_back (sfinae_constraint (p, &(tp.var), &b));
      bi <- bi + 1;
    }
    j <- j + 1;
  }

  let cn: i64 = constraints.size ();
  if cn > 0 {
    if is_definition {
      p.write ", typename";
    } else {
      p.write ", typename = std::enable_if_t<";
      let ci: i64 = 0;
      for ci < cn {
        if ci > 0 { p.write " && "; }
        let c: String = constraints.get ci;
        p.write_str (&c);
        ci <- ci + 1;
      }
      p.write ">";
    }
  }
  p.write "> ";
}
