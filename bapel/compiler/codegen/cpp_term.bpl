implements bapel.compiler.codegen.cpp

imports {
  bapel.core
  bapel.stl
  bapel.compiler.ir
}

fn emit_last_term_prefix(p: &CppPrinter) -> bool {
  let was_last: bool = (*p).last_term;
  if was_last {
    let dest: String = (*p).var_destination;
    if dest.empty () {
      p.write "return ";
    } else {
      p.write_str (&dest);
      p.write " = ";
    }
    Ptr::set (p, set (*p) { last_term = false });
  }
  was_last
}

fn restore_last_term(p: &CppPrinter, was_last: bool) -> () {
  if was_last {
    Ptr::set (p, set (*p) { last_term = true });
  }
}

fn print_types_bind(p: &CppPrinter, types: &Vector IrType) -> () {
  let old_pos: Position = (*p).position;
  Ptr::set (p, set (*p) { position = variant { Position bind_position = () } });
  let n: i64 = types.size ();
  let i: i64 = 0;
  for i < n {
    if i > 0 { p.write ", "; }
    let t: IrType = types.get i;
    print_type (p, &t);
    i <- i + 1;
  }
  Ptr::set (p, set (*p) { position = old_pos });
}

fn print_arg_terms(p: &CppPrinter, arg: &IrTerm) -> () {
  match (*arg).data {
    tuple_t elems => {
      let n: i64 = elems.size ();
      let i: i64 = 0;
      for i < n {
        if i > 0 { p.write ", "; }
        let e: IrTerm = elems.get i;
        print_term (p, &e);
        i <- i + 1;
      }
    },
    _ => print_term (p, arg)
  }
}

pub fn print_app_type_term(p: &CppPrinter, term: &IrTerm) -> () {
  let app_res: (IrTerm, Vector IrType) = term.app_types ();
  let arg: IrTerm = app_res.0;
  let types: Vector IrType = app_res.1;
  match arg.data {
    const_t _ => {
      p.write "static_cast<";
      print_types_bind (p, &types);
      p.write ">(";
      print_term (p, &arg);
      p.write ")";
      return ();
    },
    _ => ()
  }
  print_term (p, &arg);
  p.write "<";
  print_types_bind (p, &types);
  p.write ">";
}

pub fn print_app_term_term(p: &CppPrinter, term: &IrTerm) -> () {
  let app_res: (IrTerm, Vector IrType, IrTerm) = term.app_args ();
  let id: IrTerm = app_res.0;
  let types: Vector IrType = app_res.1;
  let arg: IrTerm = app_res.2;

  match id.data {
    var_t vname => {
      let vv: StringView = vname.view ();
      if vv == "ifthen" {
        match arg.data {
          tuple_t elems => {
            if elems.size () >= 2 {
              let cond: IrTerm = elems.get 0;
              let then_term: IrTerm = elems.get 1;
              p.write "if (";
              let old_last: bool = (*p).last_term;
              Ptr::set (p, set (*p) { last_term = false });
              print_term (p, &cond);
              Ptr::set (p, set (*p) { last_term = old_last });
              p.write ") ";
              print_term (p, &then_term);
              return ();
            }
          },
          _ => ()
        }
      }
      if vv == "ifelse" {
        match arg.data {
          tuple_t elems => {
            if elems.size () >= 3 {
              let cond: IrTerm = elems.get 0;
              let then_term: IrTerm = elems.get 1;
              let else_term: IrTerm = elems.get 2;
              p.write "if (";
              let old_last: bool = (*p).last_term;
              Ptr::set (p, set (*p) { last_term = false });
              print_term (p, &cond);
              Ptr::set (p, set (*p) { last_term = old_last });
              p.write ") ";
              print_term (p, &then_term);
              p.write " else ";
              print_term (p, &else_term);
              return ();
            }
          },
          _ => ()
        }
      }
      if vv == "core::for" {
        match arg.data {
          tuple_t elems => {
            if elems.size () >= 2 {
              let cond: IrTerm = elems.get 0;
              let lam: IrTerm = elems.get 1;
              let fn_res: (Vector (TypeParam IrType), Vector FunctionArg, IrTerm) = lam.to_function ();
              let body: IrTerm = fn_res.2;
              let old_last: bool = (*p).last_term;
              Ptr::set (p, set (*p) { last_term = false });
              p.write "while (";
              print_term (p, &cond);
              p.write ") ";
              print_term (p, &body);
              Ptr::set (p, set (*p) { last_term = old_last });
              return ();
            }
          },
          _ => ()
        }
      }
      if is_operator (&vname) {
        match arg.data {
          tuple_t elems => {
            let en: i64 = elems.size ();
            if en > 0 {
              p.write "(";
              let e0: IrTerm = elems.get 0;
              print_term (p, &e0);
              p.write ")";
              print_term (p, &id);
              p.write "(";
              let rest: Vector IrTerm = Vector::mk [IrTerm] ();
              let i: i64 = 1;
              for i < en {
                rest.push_back (elems.get i);
                i <- i + 1;
              }
              let rest_tuple: IrTerm = IrTerm::new_tuple rest;
              print_term (p, &rest_tuple);
              p.write ")";
              return ();
            }
          },
          _ => ()
        }
        print_term (p, &id);
        p.write " ";
        print_term (p, &arg);
        return ();
      }

      let sep: String = String::from_view "::";
      let sep_pos: i64 = vname.rfind (&sep);
      if sep_pos >= 0 {
        let trait_name: String = vname.substr (0, sep_pos);
        let method_name: String = vname.substr (sep_pos + 2, vname.size () - (sep_pos + 2));
        let opt_td: Optional IrDecl = p.find_trait_decl (&trait_name);
        if opt_td.has_value () {
          p.write "::";
          let tc: String = trait_cpp_name (&trait_name);
          p.write_str (&tc);
          p.write "<";
          print_types_bind (p, &types);
          p.write ">::";
          p.write_str (&method_name);
          p.write "(";
          print_arg_terms (p, &arg);
          p.write ")";
          return ();
        }

        let opt_d: Optional IrDecl = p.find_decl (&trait_name);
        if opt_d.has_value () {
          let d: IrDecl = opt_d.get_value ();
          if is_type_decl (&d) {
            let arity: i64 = match d.data {
              name_decl nd => count_type_vars (&(nd.kind)),
              alias_decl ad => count_alias_type_vars (&(ad.typ)),
              _ => 0
            };
            let tn: i64 = types.size ();
            if tn >= arity {
              let type_args: Vector IrType = Vector::mk [IrType] ();
              let method_args: Vector IrType = Vector::mk [IrType] ();
              let ti: i64 = 0;
              for ti < arity {
                type_args.push_back (types.get ti);
                ti <- ti + 1;
              }
              for ti < tn {
                method_args.push_back (types.get ti);
                ti <- ti + 1;
              }
              let mapped_type: String = inherent_cpp_name (&trait_name);
              p.write "::";
              p.write_str (&mapped_type);
              if arity > 0 {
                p.write "<";
                print_types_bind (p, &type_args);
                p.write ">";
              }
              p.write "::";
              p.write_str (&method_name);
              if method_args.size () > 0 {
                p.write "<";
                print_types_bind (p, &method_args);
                p.write ">";
              }
              p.write "(";
              print_arg_terms (p, &arg);
              p.write ")";
              return ();
            }
          }
        }
      }
    },
    _ => ()
  }

  print_term (p, &id);
  match id.data {
    var_t vname => {
      if !is_operator (&vname) && types.size () > 0 {
        p.write "<";
        print_types_bind (p, &types);
        p.write ">";
      }
    },
    _ => ()
  }
  p.write "(";
  print_arg_terms (p, &arg);
  p.write ")";
}

pub fn print_let_term(p: &CppPrinter, term: &IrTerm) -> () {
  match (*term).data {
    let_t c => {
      let val_term: &IrTerm = c.value.get ();
      let auto_flag: bool = match (*val_term).data {
        type_abs _ => true,
        _ => false
      };
      let old_auto: bool = (*p).auto_type;
      let old_pos: Position = (*p).position;
      Ptr::set (p, set (*p) { auto_type = auto_flag, position = variant { Position bind_position = () } });
      if c.var_type.has_value () {
        let vt: IrType = c.var_type.get_value ();
        print_type (p, &vt);
        p.write " ";
      }
      p.write_str (&(c.var));
      Ptr::set (p, set (*p) { auto_type = old_auto, position = old_pos });

      if is_cpp_statement val_term {
        let old_dest: String = (*p).var_destination;
        let old_last: bool = (*p).last_term;
        Ptr::set (p, set (*p) { var_destination = c.var, last_term = true });
        p.write ";\n";
        print_term (p, val_term);
        Ptr::set (p, set (*p) { var_destination = old_dest, last_term = old_last });
      } else {
        p.write " = ";
        print_term (p, val_term);
      }
    },
    _ => ()
  }
}

pub fn print_match_term(p: &CppPrinter, term: &IrTerm) -> () {
  match (*term).data {
    match_t c => {
      let variant_id: String = p.gen_id ();
      p.write "{\nauto ";
      p.write_str (&variant_id);
      p.write " = ";
      let old_last: bool = (*p).last_term;
      Ptr::set (p, set (*p) { last_term = false });
      print_term (p, c.term.get ());
      Ptr::set (p, set (*p) { last_term = old_last });
      p.write ";\nswitch (";
      p.write_str (&variant_id);
      p.write ".index()) {\n";

      let arms: Vector (MatchArm IrTerm) = c.arms;
      let n: i64 = arms.size ();
      let i: i64 = 0;
      for i < n {
        let arm: MatchArm IrTerm = arms.get i;
        if arm.tag.view () == "_" {
          p.write "default: {\n";
          print_term (p, arm.body.get ());
          p.write ";\n}\n";
        } else {
          let idx: i64 = if arm.index.has_value () { arm.index.get_value () } else { 0 };
          let idx_s: String = i64_to_string idx;
          p.write "case ";
          p.write_str (&idx_s);
          p.write ": {\nauto &";
          p.write_str (&(arm.arg));
          p.write " = std::get<";
          p.write_str (&idx_s);
          p.write ">(";
          p.write_str (&variant_id);
          p.write ");\n";
          print_term (p, arm.body.get ());
          p.write ";\n}\n";
        }
        i <- i + 1;
      }
      p.write "} }";
    },
    _ => ()
  }
}

fn find_field_name(fields: &Vector (NamedType IrType), label: &String) -> String {
  let n: i64 = fields.size ();
  let i: i64 = 0;
  let found: String = String::from_view "";
  for i < n && found.empty () {
    let f: NamedType IrType = fields.get i;
    if f.id == *label {
      found <- f.id;
    }
    i <- i + 1;
  }
  if found.empty () {
    let opt_idx: Optional i64 = parse_i64 (label.view ());
    if opt_idx.has_value () {
      let idx: i64 = opt_idx.get_value ();
      if idx >= 0 && idx < n {
        let f: NamedType IrType = fields.get idx;
        found <- f.id;
      }
    }
  }
  found
}

fn find_tag_index(tags: &Vector (NamedType IrType), label: &String) -> i64 {
  let n: i64 = tags.size ();
  let i: i64 = 0;
  let found: i64 = 0 - 1;
  for i < n && found < 0 {
    let tag: NamedType IrType = tags.get i;
    if tag.id == *label {
      found <- i;
    }
    i <- i + 1;
  }
  found
}

pub fn print_projection_term(p: &CppPrinter, term: &IrTerm) -> () {
  match (*term).data {
    projection c => {
      let obj_type: IrType = if c.reduced_type.has_value () {
        c.reduced_type.get_value ()
      } else if (*term).typ.has_value () {
        (*term).typ.get_value ()
      } else {
        IrType::new_name (String::from_view "")
      };
      match obj_type.data {
        mu_t m => {
          obj_type <- *(m.body.get ());
          ()
        },
        _ => ()
      }

      match obj_type.data {
        struct_t fields => {
          let fname: String = find_field_name (&fields, &(c.label));
          if !(fname.empty ()) {
            print_term (p, c.term.get ());
            p.write ".";
            p.write_str (&fname);
            return ();
          }
        },
        tuple_t elems => {
          let opt_idx: Optional i64 = parse_i64 (c.label.view ());
          if opt_idx.has_value () {
            let idx: i64 = opt_idx.get_value ();
            if idx >= 0 && idx < elems.size () {
              p.write "std::get<";
              let idx_s: String = i64_to_string idx;
              p.write_str (&idx_s);
              p.write ">(";
              print_term (p, c.term.get ());
              p.write ")";
              return ();
            }
          }
        },
        variant_t tags => {
          let idx: i64 = find_tag_index (&tags, &(c.label));
          if idx >= 0 {
            p.write "std::get<";
            let idx_s: String = i64_to_string idx;
            p.write_str (&idx_s);
            p.write ">(";
            print_term (p, c.term.get ());
            p.write ")";
            return ();
          }
        },
        _ => ()
      }

      p.write "std::get<";
      p.write_str (&(c.label));
      p.write ">(";
      print_term (p, c.term.get ());
      p.write ")";
    },
    _ => ()
  }
}

pub fn print_set_term(p: &CppPrinter, term: &IrTerm) -> () {
  match (*term).data {
    set_t c => {
      let obj_type: IrType = if c.reduced_type.has_value () {
        c.reduced_type.get_value ()
      } else if (*term).typ.has_value () {
        (*term).typ.get_value ()
      } else {
        IrType::new_name (String::from_view "")
      };
      match obj_type.data {
        mu_t m => {
          obj_type <- *(m.body.get ());
          ()
        },
        _ => ()
      }

      let values: Vector (LabelValue IrTerm) = c.values;
      let vn: i64 = values.size ();
      match obj_type.data {
        struct_t fields => {
          let struct_id: String = p.gen_id ();
          p.write "([&, ";
          p.write_str (&struct_id);
          p.write " = ";
          print_term (p, c.term.get ());
          p.write "]() mutable {\n";
          let vi: i64 = 0;
          for vi < vn {
            let lv: LabelValue IrTerm = values.get vi;
            let fname: String = find_field_name (&fields, &(lv.label));
            if !(fname.empty ()) {
              p.write_str (&struct_id);
              p.write ".";
              p.write_str (&fname);
              p.write " = ";
              print_term (p, lv.value.get ());
              p.write ";\n";
            }
            vi <- vi + 1;
          }
          p.write "return ";
          p.write_str (&struct_id);
          p.write ";\n})()";
        },
        tuple_t _ => {
          let tuple_id: String = p.gen_id ();
          p.write "([";
          p.write_str (&tuple_id);
          p.write " = ";
          print_term (p, c.term.get ());
          p.write "]() mutable {\n";
          let vi: i64 = 0;
          for vi < vn {
            let lv: LabelValue IrTerm = values.get vi;
            p.write "std::get<";
            p.write_str (&(lv.label));
            p.write ">(";
            p.write_str (&tuple_id);
            p.write ") = ";
            print_term (p, lv.value.get ());
            p.write ";\n";
            vi <- vi + 1;
          }
          p.write "return ";
          p.write_str (&tuple_id);
          p.write ";\n})()";
        },
        _ => ()
      }
    },
    _ => ()
  }
}

pub fn print_lambda_term(p: &CppPrinter, term: &IrTerm) -> () {
  let fn_res: (Vector (TypeParam IrType), Vector FunctionArg, IrTerm) = term.to_function ();
  let tvars: Vector (TypeParam IrType) = fn_res.0;
  let args: Vector FunctionArg = fn_res.1;
  let body: IrTerm = fn_res.2;

  p.write "[&]";
  let tn: i64 = tvars.size ();
  if tn > 0 {
    p.write "<";
    let ti: i64 = 0;
    for ti < tn {
      if ti > 0 { p.write ", "; }
      let tp: TypeParam IrType = tvars.get ti;
      p.write "typename ";
      p.write_str (&(tp.var));
      ti <- ti + 1;
    }
    p.write ">";
  }

  let an: i64 = args.size ();
  let argless: bool = false;
  if an == 1 {
    let a0: FunctionArg = args.get 0;
    match a0.typ.data {
      tuple_t elems => {
        if elems.size () == 0 {
          argless <- true;
        }
        ()
      },
      _ => ()
    }
  }

  p.write "(";
  if !argless {
    let ai: i64 = 0;
    for ai < an {
      if ai > 0 { p.write ", "; }
      let a: FunctionArg = args.get ai;
      print_type (p, &(a.typ));
      p.write " ";
      let cid: String = p.to_id (&(a.id));
      p.write_str (&cid);
      ai <- ai + 1;
    }
  }
  p.write ")";

  let old_last: bool = (*p).last_term;
  Ptr::set (p, set (*p) { last_term = true });
  print_term (p, &body);
  Ptr::set (p, set (*p) { last_term = old_last });
}

pub fn print_term(p: &CppPrinter, term: &IrTerm) -> () {
  match (*term).data {
    app_type _ => {
      let was_last: bool = emit_last_term_prefix p;
      print_app_type_term (p, term);
      restore_last_term (p, was_last);
    },
    app_term _ => {
      if is_cpp_statement term {
        print_app_term_term (p, term);
      } else {
        let was_last: bool = emit_last_term_prefix p;
        print_app_term_term (p, term);
        restore_last_term (p, was_last);
      }
    },
    assign a => {
      let was_last: bool = emit_last_term_prefix p;
      let old_pos: Position = (*p).position;
      Ptr::set (p, set (*p) { position = variant { Position bind_position = () } });
      print_term (p, a.ret.get ());
      Ptr::set (p, set (*p) { position = old_pos });
      p.write " = ";
      print_term (p, a.arg.get ());
      restore_last_term (p, was_last);
    },
    block terms => {
      p.write "{\n";
      let n: i64 = terms.size ();
      let i: i64 = 0;
      for i < n {
        let t: IrTerm = terms.get i;
        if i + 1 == n {
          print_term (p, &t);
          p.write ";}\n";
        } else {
          let old_last: bool = (*p).last_term;
          Ptr::set (p, set (*p) { last_term = false });
          print_term (p, &t);
          p.write ";\n";
          Ptr::set (p, set (*p) { last_term = old_last });
        }
        i <- i + 1;
      }
    },
    const_t lit => {
      let was_last: bool = emit_last_term_prefix p;
      match lit {
        int_lit v => {
          let s: String = i64_to_string v;
          p.write_str (&s);
        },
        float_lit fl => {
          let is: String = i64_to_string fl.integer;
          let ds: String = i64_to_string fl.decimal;
          p.write_str (&is);
          p.write ".";
          p.write_str (&ds);
        },
        rune_lit r => {
          p.write "'";
          p.write_str (&r);
          p.write "'";
        },
        str_lit s => {
          p.write "\"";
          p.write_str (&s);
          p.write "\"";
        }
      }
      restore_last_term (p, was_last);
    },
    injection inj => {
      let was_last: bool = emit_last_term_prefix p;
      print_type (p, &(inj.variant_type));
      p.write "(";
      if inj.tag_index.has_value () {
        p.write "std::in_place_index<";
        let idx_s: String = i64_to_string (inj.tag_index.get_value ());
        p.write_str (&idx_s);
        p.write ">, ";
      }
      let val_term: &IrTerm = inj.value.get ();
      match (*val_term).data {
        struct_t _ => {
          let tag_type: IrType = IrType::new_name (String::from_view "");
          let found_tag_type: bool = false;
          if inj.tag_index.has_value () {
            let vname: String = base_type_name (&(inj.variant_type));
            if !(vname.empty ()) {
              let opt_vd: Optional IrDecl = p.find_decl (&vname);
              if opt_vd.has_value () {
                let vd: IrDecl = opt_vd.get_value ();
                match vd.data {
                  alias_decl ad => {
                    let vparams: Vector String = type_lambda_vars (&(ad.typ));
                    let vbody: IrType = type_lambda_body (&(ad.typ));
                    match vbody.data {
                      mu_t m => {
                        vbody <- *(m.body.get ());
                        ()
                      },
                      _ => ()
                    }
                    match vbody.data {
                      variant_t tags => {
                        let idx: i64 = inj.tag_index.get_value ();
                        if idx >= 0 && idx < tags.size () {
                          let tag: NamedType IrType = tags.get idx;
                          tag_type <- *(tag.typ.get ());
                          let vargs: Vector IrType = type_app_args (&(inj.variant_type));
                          let pn: i64 = vparams.size ();
                          let an: i64 = vargs.size ();
                          let pi: i64 = 0;
                          for pi < pn && pi < an {
                            let vp: String = vparams.get pi;
                            let va: IrType = vargs.get pi;
                            tag_type <- substitute_var_type (&tag_type, &vp, &va);
                            pi <- pi + 1;
                          }
                          found_tag_type <- true;
                        }
                        ()
                      },
                      _ => ()
                    }
                  },
                  _ => ()
                }
              }
            }
          }
          let old_pos: Position = (*p).position;
          Ptr::set (p, set (*p) { position = variant { Position bind_position = () } });
          if found_tag_type {
            print_type (p, &tag_type);
          } else if (*val_term).typ.has_value () {
            let vt: IrType = (*val_term).typ.get_value ();
            print_type (p, &vt);
          }
          Ptr::set (p, set (*p) { position = old_pos });
        },
        _ => ()
      }
      print_term (p, val_term);
      p.write ")";
      restore_last_term (p, was_last);
    },
    lambda _ => {
      let was_last: bool = emit_last_term_prefix p;
      let old_last: bool = (*p).last_term;
      Ptr::set (p, set (*p) { last_term = false });
      print_lambda_term (p, term);
      Ptr::set (p, set (*p) { last_term = old_last });
      restore_last_term (p, was_last);
    },
    type_abs _ => {
      let was_last: bool = emit_last_term_prefix p;
      let old_last: bool = (*p).last_term;
      Ptr::set (p, set (*p) { last_term = false });
      print_lambda_term (p, term);
      Ptr::set (p, set (*p) { last_term = old_last });
      restore_last_term (p, was_last);
    },
    let_t _ => {
      let old_last: bool = (*p).last_term;
      Ptr::set (p, set (*p) { last_term = false });
      print_let_term (p, term);
      Ptr::set (p, set (*p) { last_term = old_last });
    },
    match_t _ => {
      print_match_term (p, term);
    },
    projection _ => {
      let was_last: bool = emit_last_term_prefix p;
      print_projection_term (p, term);
      restore_last_term (p, was_last);
    },
    return_t r => {
      let old_last: bool = (*p).last_term;
      Ptr::set (p, set (*p) { last_term = false });
      p.write "return ";
      print_term (p, r.get ());
      p.write ";";
      Ptr::set (p, set (*p) { last_term = old_last });
    },
    tuple_t elems => {
      let was_last: bool = emit_last_term_prefix p;
      let n: i64 = elems.size ();
      match (*p).position {
        bind_position _ => p.write "std::tie(",
        type_position _ => {
          if n == 0 {
            p.write "std::monostate(";
          } else {
            p.write "std::make_tuple(";
          }
        }
      }
      let i: i64 = 0;
      for i < n {
        if i > 0 { p.write ", "; }
        let e: IrTerm = elems.get i;
        print_term (p, &e);
        i <- i + 1;
      }
      p.write ")";
      restore_last_term (p, was_last);
    },
    set_t _ => {
      let was_last: bool = emit_last_term_prefix p;
      print_set_term (p, term);
      restore_last_term (p, was_last);
    },
    struct_t values => {
      let was_last: bool = emit_last_term_prefix p;
      p.write "{";
      let n: i64 = values.size ();
      let i: i64 = 0;
      for i < n {
        if i > 0 { p.write ", "; }
        let lv: LabelValue IrTerm = values.get i;
        p.write ".";
        p.write_str (&(lv.label));
        p.write " = ";
        print_term (p, lv.value.get ());
        i <- i + 1;
      }
      p.write "}";
      restore_last_term (p, was_last);
    },
    var_t v => {
      let was_last: bool = emit_last_term_prefix p;
      let cid: String = p.to_id (&v);
      p.write_str (&cid);
      restore_last_term (p, was_last);
    }
  }
}
