implements bapel.compiler.codegen.cpp

imports {
  bapel.core
  bapel.stl
  bapel.compiler.ir
}

fn open_namespace(p: &CppPrinter, id: &String) -> (bool, String) {
  let sep: String = String::from_view "::";
  let pos: i64 = id.rfind (&sep);
  if pos < 0 {
    return (false, *id);
  }
  let ns: String = id.substr (0, pos);
  let base_id: String = id.substr (pos + 2, id.size () - (pos + 2));
  p.write "namespace ";
  p.write_str (&ns);
  p.write " { ";
  (true, base_id)
}

fn close_namespace(p: &CppPrinter, opened: bool) -> () {
  if opened {
    p.write " }";
  }
  ()
}

fn is_mode_public(mode: &PrinterMode) -> bool {
  match *mode {
    mode_public_header _ => true,
    _ => false
  }
}

fn is_mode_private(mode: &PrinterMode) -> bool {
  match *mode {
    mode_private_header _ => true,
    _ => false
  }
}

fn is_mode_source(mode: &PrinterMode) -> bool {
  match *mode {
    mode_source _ => true,
    _ => false
  }
}

pub fn print_alias_decl(p: &CppPrinter, id: &String, typ: &IrType, tvars: &Vector String) -> () {
  match (*typ).data {
    lambda_t _ => {
      let lvars: Vector String = typ.lambda_vars ();
      p.write "template <";
      let n: i64 = lvars.size ();
      let i: i64 = 0;
      for i < n {
        if i > 0 { p.write ", "; }
        p.write "typename ";
        let tv: String = lvars.get i;
        p.write_str (&tv);
        i <- i + 1;
      }
      p.write "> ";
      let body: IrType = typ.lambda_body ();
      print_alias_decl (p, id, &body, &lvars);
    },
    mu_t m => {
      let self_name: String = String::from_view "::";
      let cid: String = p.to_id id;
      self_name.append (cid.view ());
      let tn: i64 = tvars.size ();
      if tn > 0 {
        self_name.append "<";
        let ti: i64 = 0;
        for ti < tn {
          if ti > 0 { self_name.append ", "; }
          let tv: String = tvars.get ti;
          self_name.append (tv.view ());
          ti <- ti + 1;
        }
        self_name.append ">";
      }
      let cur_map: Map String String = (*p).mu_self_names;
      let prev: Optional String = cur_map.get (&(m.var));
      cur_map.insert (m.var, self_name);
      Ptr::set (p, set (*p) { mu_self_names = cur_map });
      print_alias_decl (p, id, m.body.get (), tvars);
      let restore_map: Map String String = (*p).mu_self_names;
      if prev.has_value () {
        restore_map.insert (m.var, prev.get_value ());
      } else {
        let _: bool = restore_map.remove (&(m.var));
      }
      Ptr::set (p, set (*p) { mu_self_names = restore_map });
    },
    name_t _ => {
      p.write "using ";
      p.write_str id;
      p.write " = ";
      print_type (p, typ);
    },
    struct_t fields => {
      p.write "struct ";
      p.write_str id;
      p.write " {\n";
      let n: i64 = fields.size ();
      let i: i64 = 0;
      for i < n {
        let f: NamedType IrType = fields.get i;
        print_type (p, f.typ.get ());
        p.write " ";
        p.write_str (&(f.id));
        p.write ";\n";
        i <- i + 1;
      }
      p.write "}";
    },
    tuple_t _ => {
      p.write "struct ";
      p.write_str id;
      p.write " : ";
      let old_pos: Position = (*p).position;
      Ptr::set (p, set (*p) { position = variant { Position bind_position = () } });
      print_type (p, typ);
      p.write " { ";
      p.write_str id;
      p.write "(const ";
      print_type (p, typ);
      p.write "& arg) : ";
      print_type (p, typ);
      p.write "(arg) {}}";
      Ptr::set (p, set (*p) { position = old_pos });
    },
    variant_t tags => {
      p.write "BAPEL_VARIANT(";
      p.write_str id;
      let n: i64 = tags.size ();
      let i: i64 = 0;
      for i < n {
        let tag: NamedType IrType = tags.get i;
        p.write ", (";
        let tid: String = p.to_id (&(tag.id));
        p.write_str (&tid);
        p.write ", (";
        let old_pos: Position = (*p).position;
        Ptr::set (p, set (*p) { position = variant { Position bind_position = () } });
        print_type (p, tag.typ.get ());
        Ptr::set (p, set (*p) { position = old_pos });
        p.write "))";
        i <- i + 1;
      }
      p.write ")";
    },
    _ => ()
  }
}

pub fn print_trait_decl(p: &CppPrinter, d: &IrDecl) -> () {
  match (*d).data {
    trait_decl td => {
      let tname: String = trait_cpp_name (&(td.id));
      let ns_res: (bool, String) = open_namespace (p, &tname);
      let base_id: String = ns_res.1;
      p.write "template <typename Self";
      let tps: Vector (TypeParam IrType) = td.type_params;
      let n: i64 = tps.size ();
      let i: i64 = 0;
      for i < n {
        let tp: TypeParam IrType = tps.get i;
        p.write ", typename ";
        p.write_str (&(tp.var));
        i <- i + 1;
      }
      p.write ">\nstruct ";
      p.write_str (&base_id);
      p.write ";\n";
      close_namespace (p, ns_res.0);
    },
    _ => ()
  }
}

pub fn print_decl(p: &CppPrinter, d: &IrDecl) -> () {
  let m: PrinterMode = (*p).mode;
  if is_mode_source (&m) { return (); }
  if is_mode_public (&m) && !((*d).export_flag) { return (); }
  if is_mode_private (&m) && (*d).export_flag { return (); }

  match (*d).data {
    trait_decl _ => {
      print_trait_decl (p, d);
    },
    name_decl _ => (),
    alias_decl ad => {
      let ns_res: (bool, String) = open_namespace (p, &(ad.id));
      let base_id: String = ns_res.1;
      match ad.typ.data {
        lambda_t _ => {
          let tvars: Vector String = ad.typ.lambda_vars ();
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
          p.write "> struct ";
          p.write_str (&base_id);
        },
        _ => {
          p.write "struct ";
          p.write_str (&base_id);
        }
      }
      p.write ";\n";
      close_namespace (p, ns_res.0);
    },
    term_decl td => {
      if p.has_function (&(td.id)) { return (); }
      let typ: IrType = td.typ;
      match typ.data {
        forall_t _ => {
          let ns_res: (bool, String) = open_namespace (p, &(td.id));
          let base_id: String = ns_res.1;
          let tvars: Vector (TypeParam IrType) = typ.forall_type_params ();
          print_template_params (p, &tvars, false);
          let inner_decl: IrDecl = IrDecl::new_term (base_id, typ.forall_body (), (*d).export_flag, (*d).pos);
          print_decl (p, &inner_decl);
          close_namespace (p, ns_res.0);
        },
        fun_t fn_t => {
          let ns_res: (bool, String) = open_namespace (p, &(td.id));
          let base_id: String = ns_res.1;
          let old_pos: Position = (*p).position;
          Ptr::set (p, set (*p) { position = variant { Position bind_position = () } });
          print_type (p, fn_t.ret.get ());
          Ptr::set (p, set (*p) { position = old_pos });
          p.write " ";
          p.write_str (&base_id);
          p.write "(";
          print_type (p, fn_t.arg.get ());
          p.write ");";
          close_namespace (p, ns_res.0);
        },
        app_t _ => {
          let ns_res: (bool, String) = open_namespace (p, &(td.id));
          let base_id: String = ns_res.1;
          let old_pos: Position = (*p).position;
          Ptr::set (p, set (*p) { position = variant { Position bind_position = () } });
          print_type (p, &typ);
          p.write " ";
          p.write_str (&base_id);
          Ptr::set (p, set (*p) { position = old_pos });
          close_namespace (p, ns_res.0);
        },
        array_t _ => {
          let ns_res: (bool, String) = open_namespace (p, &(td.id));
          let base_id: String = ns_res.1;
          let old_pos: Position = (*p).position;
          Ptr::set (p, set (*p) { position = variant { Position bind_position = () } });
          print_type (p, &typ);
          p.write " ";
          p.write_str (&base_id);
          Ptr::set (p, set (*p) { position = old_pos });
          close_namespace (p, ns_res.0);
        },
        exist_var_t _ => {
          let ns_res: (bool, String) = open_namespace (p, &(td.id));
          let base_id: String = ns_res.1;
          let old_pos: Position = (*p).position;
          Ptr::set (p, set (*p) { position = variant { Position bind_position = () } });
          print_type (p, &typ);
          p.write " ";
          p.write_str (&base_id);
          Ptr::set (p, set (*p) { position = old_pos });
          close_namespace (p, ns_res.0);
        },
        name_t _ => {
          let ns_res: (bool, String) = open_namespace (p, &(td.id));
          let base_id: String = ns_res.1;
          let old_pos: Position = (*p).position;
          Ptr::set (p, set (*p) { position = variant { Position bind_position = () } });
          print_type (p, &typ);
          p.write " ";
          p.write_str (&base_id);
          Ptr::set (p, set (*p) { position = old_pos });
          close_namespace (p, ns_res.0);
        },
        tuple_t _ => {
          let ns_res: (bool, String) = open_namespace (p, &(td.id));
          let base_id: String = ns_res.1;
          let old_pos: Position = (*p).position;
          Ptr::set (p, set (*p) { position = variant { Position bind_position = () } });
          print_type (p, &typ);
          p.write " ";
          p.write_str (&base_id);
          Ptr::set (p, set (*p) { position = old_pos });
          close_namespace (p, ns_res.0);
        },
        var_t _ => {
          let ns_res: (bool, String) = open_namespace (p, &(td.id));
          let base_id: String = ns_res.1;
          let old_pos: Position = (*p).position;
          Ptr::set (p, set (*p) { position = variant { Position bind_position = () } });
          print_type (p, &typ);
          p.write " ";
          p.write_str (&base_id);
          Ptr::set (p, set (*p) { position = old_pos });
          close_namespace (p, ns_res.0);
        },
        _ => ()
      }
    },
    impl_decl _ => ()
  }
  ()
}

pub fn print_type_def(p: &CppPrinter, d: &IrDecl) -> () {
  let m: PrinterMode = (*p).mode;
  if is_mode_source (&m) { return (); }
  if is_mode_public (&m) && !((*d).export_flag) { return (); }
  if is_mode_private (&m) && (*d).export_flag { return (); }

  match (*d).data {
    name_decl nd => {
      let nt: IrType = IrType::new_name nd.id;
      print_type (p, &nt);
      p.write ";\n";
    },
    alias_decl ad => {
      let anonym_pfx: String = String::from_view "__anonym_";
      let is_anonym: bool = ad.id.starts_with (&anonym_pfx);
      if is_anonym {
        p.write "#ifndef ";
        p.write_str (&(ad.id));
        p.write "_DEFINED\n#define ";
        p.write_str (&(ad.id));
        p.write "_DEFINED\n";
      }
      let ns_res: (bool, String) = open_namespace (p, &(ad.id));
      let base_id: String = ns_res.1;
      let empty_tvars: Vector String = Vector::mk [String] ();
      print_alias_decl (p, &base_id, &(ad.typ), &empty_tvars);
      p.write ";\n";
      close_namespace (p, ns_res.0);
      if is_anonym {
        p.write "#endif\n";
      }
      ()
    },
    _ => ()
  }
}

pub fn print_function_signature(p: &CppPrinter, f: &IrFunction) -> () {
  let ns_res: (bool, String) = open_namespace (p, &((*f).id));
  let base_id: String = ns_res.1;
  let tps: Vector (TypeParam IrType) = (*f).type_params;
  print_template_params (p, &tps, false);
  let old_pos: Position = (*p).position;
  Ptr::set (p, set (*p) { position = variant { Position bind_position = () } });
  print_type (p, &((*f).ret_type));
  Ptr::set (p, set (*p) { position = old_pos });
  p.write " ";
  p.write_str (&base_id);
  p.write "(";
  let args: Vector FunctionArg = (*f).args;
  let n: i64 = args.size ();
  let i: i64 = 0;
  for i < n {
    if i > 0 { p.write ", "; }
    let a: FunctionArg = args.get i;
    print_type (p, &(a.typ));
    p.write " ";
    p.write_str (&(a.id));
    i <- i + 1;
  }
  p.write ");\n";
  close_namespace (p, ns_res.0);
}

pub fn print_function_full(p: &CppPrinter, f: &IrFunction) -> () {
  let ns_res: (bool, String) = open_namespace (p, &((*f).id));
  let base_id: String = ns_res.1;
  let tps: Vector (TypeParam IrType) = (*f).type_params;
  print_template_params (p, &tps, true);
  let old_pos: Position = (*p).position;
  Ptr::set (p, set (*p) { position = variant { Position bind_position = () } });
  print_type (p, &((*f).ret_type));
  Ptr::set (p, set (*p) { position = old_pos });
  p.write " ";
  p.write_str (&base_id);
  p.write "(";
  let args: Vector FunctionArg = (*f).args;
  let n: i64 = args.size ();
  let i: i64 = 0;
  for i < n {
    if i > 0 { p.write ", "; }
    let a: FunctionArg = args.get i;
    print_type (p, &(a.typ));
    p.write " ";
    p.write_str (&(a.id));
    i <- i + 1;
  }
  p.write ")\n";
  let old_last: bool = (*p).last_term;
  Ptr::set (p, set (*p) { last_term = true });
  print_term (p, &((*f).body));
  Ptr::set (p, set (*p) { last_term = old_last });
  p.write "\n";
  close_namespace (p, ns_res.0);
}

fn print_impl_methods(p: &CppPrinter, methods: &Vector IrFunction) -> () {
  let mn: i64 = methods.size ();
  let mi: i64 = 0;
  for mi < mn {
    let m: IrFunction = methods.get mi;
    let mtps: Vector (TypeParam IrType) = m.type_params;
    print_template_params (p, &mtps, false);
    p.write "  static inline ";
    let old_pos: Position = (*p).position;
    Ptr::set (p, set (*p) { position = variant { Position bind_position = () } });
    print_type (p, &(m.ret_type));
    Ptr::set (p, set (*p) { position = old_pos });
    p.write " ";
    p.write_str (&(m.id));
    p.write "(";
    let margs: Vector FunctionArg = m.args;
    let an: i64 = margs.size ();
    let ai: i64 = 0;
    for ai < an {
      if ai > 0 { p.write ", "; }
      let a: FunctionArg = margs.get ai;
      print_type (p, &(a.typ));
      p.write " ";
      p.write_str (&(a.id));
      ai <- ai + 1;
    }
    p.write ") ";
    let old_last: bool = (*p).last_term;
    Ptr::set (p, set (*p) { last_term = true });
    print_term (p, &(m.body));
    Ptr::set (p, set (*p) { last_term = old_last });
    p.write "\n";
    mi <- mi + 1;
  }
  ()
}

pub fn print_trait_impl(p: &CppPrinter, ti: &IrTraitImpl) -> () {
  let exported: bool = false;
  let base_name: String = (*ti).type_name.base_name ();
  if !(base_name.empty ()) {
    let opt_d: Optional IrDecl = p.find_decl (&base_name);
    if opt_d.has_value () {
      let d: IrDecl = opt_d.get_value ();
      exported <- d.export_flag;
    }
  }

  let m: PrinterMode = (*p).mode;
  if is_mode_source (&m) { return (); }
  if is_mode_public (&m) && !exported { return (); }
  if is_mode_private (&m) && exported { return (); }

  let tps: Vector (TypeParam IrType) = (*ti).type_params;
  let tn: i64 = tps.size ();
  let methods: Vector IrFunction = (*ti).methods;

  match (*ti).case_val {
    inherent_impl _ => {
      let inh_name: String = inherent_cpp_name (&base_name);
      let ns_res: (bool, String) = open_namespace (p, &inh_name);
      let base_id: String = ns_res.1;
      if tn > 0 {
        p.write "template <";
        let i: i64 = 0;
        for i < tn {
          if i > 0 { p.write ", "; }
          let tp: TypeParam IrType = tps.get i;
          p.write "typename ";
          p.write_str (&(tp.var));
          i <- i + 1;
        }
        p.write ">\n";
      }
      p.write "struct ";
      p.write_str (&base_id);
      p.write " {\n  ";
      p.write_str (&base_id);
      p.write "() = delete;\n";
      print_impl_methods (p, &methods);
      p.write "};\n";
      close_namespace (p, ns_res.0);
    },
    trait_impl _ => {
      let trait_name: String = (*ti).trait_type.base_name ();
      let tcpp: String = trait_cpp_name (&trait_name);
      let ns_res: (bool, String) = open_namespace (p, &tcpp);
      let base_id: String = ns_res.1;
      if tn > 0 {
        p.write "template <";
        let i: i64 = 0;
        for i < tn {
          if i > 0 { p.write ", "; }
          let tp: TypeParam IrType = tps.get i;
          p.write "typename ";
          p.write_str (&(tp.var));
          i <- i + 1;
        }
        p.write ">\n";
      } else {
        p.write "template <>\n";
      }
      p.write "struct ";
      p.write_str (&base_id);
      p.write "<";
      let old_pos: Position = (*p).position;
      Ptr::set (p, set (*p) { position = variant { Position bind_position = () } });
      print_type (p, &((*ti).type_name));
      let targs: Vector IrType = (*ti).trait_type.app_args ();
      let an: i64 = targs.size ();
      let ai: i64 = 0;
      for ai < an {
        p.write ", ";
        let a: IrType = targs.get ai;
        print_type (p, &a);
        ai <- ai + 1;
      }
      Ptr::set (p, set (*p) { position = old_pos });
      p.write "> {\n";
      print_impl_methods (p, &methods);
      p.write "};\n";
      close_namespace (p, ns_res.0);
    }
  }
  ()
}

pub fn print_module_top(p: &CppPrinter) -> () {
  let mod_ref: &IrTypedModule = (*p).mod;
  match (*p).mode {
    mode_public_header _ => {
      p.write "#pragma once\n\n#include <array>\n#include <cstdlib>\n#include <cmath>\n#include <functional>\n#include <optional>\n#include <string>\n#include <tuple>\n#include <variant>\n#include <vector>\n\n#include \"bapel/core_variant.h\"\n\n";
    },
    mode_private_header _ => {
      p.write "#pragma once\n\n#include \"";
      let h_path: String = to_header_path (&((*mod_ref).module_id));
      p.write_str (&h_path);
      p.write "\"\n\n";
    },
    mode_source _ => {
      p.write "\n#include \"";
      let h_path: String = to_header_path (&((*mod_ref).module_id));
      let dot_h: String = String::from_view ".h";
      let pos: i64 = h_path.rfind (&dot_h);
      let priv_path: String = if pos >= 0 { h_path.substr (0, pos) } else { h_path };
      priv_path.append "_private.h";
      p.write_str (&priv_path);
      p.write "\"\n\n";
    }
  }
  ()
}

pub fn print_imports(p: &CppPrinter) -> () {
  let m: PrinterMode = (*p).mode;
  if is_mode_public (&m) || is_mode_source (&m) {
    let mod_ref: &IrTypedModule = (*p).mod;
    let imps: Vector IrResolvedImport = (*mod_ref).mod_imports;
    p.write "\n";
    let n: i64 = imps.size ();
    let i: i64 = 0;
    for i < n {
      let imp: IrResolvedImport = imps.get i;
      p.write "#include \"";
      let hp: String = to_header_path (&(imp.module_id));
      p.write_str (&hp);
      p.write "\"\n";
      i <- i + 1;
    }
    p.write "\n";
  }
  ()
}

pub fn print_impls(p: &CppPrinter) -> () {
  let m: PrinterMode = (*p).mode;
  if !is_mode_public (&m) { return (); }
  let mod_ref: &IrTypedModule = (*p).mod;
  let mod_impls: Vector IrImpl = (*mod_ref).mod_impls;
  let dot_h: String = String::from_view ".h";
  p.write "\n";
  let n: i64 = mod_impls.size ();
  let i: i64 = 0;
  for i < n {
    let imp: IrImpl = mod_impls.get i;
    let fn_str: String = imp.relative_filename.value;
    if fn_str.ends_with (&dot_h) {
      p.write "#include \"";
      p.write_str (&fn_str);
      p.write "\"\n";
    }
    i <- i + 1;
  }
  p.write "\n";
}

pub fn print_decls(p: &CppPrinter, decls: &Vector IrDecl) -> () {
  let n: i64 = decls.size ();
  let i: i64 = 0;
  for i < n {
    let d: IrDecl = decls.get i;
    match d.data {
      alias_decl ad => {
        if ad.typ.is_mu_alias () {
          print_decl (p, &d);
        }
      },
      _ => ()
    }
    i <- i + 1;
  }
  let j: i64 = 0;
  for j < n {
    let d: IrDecl = decls.get j;
    match d.data {
      alias_decl _ => print_type_def (p, &d),
      _ => print_decl (p, &d)
    }
    j <- j + 1;
  }
  ()
}

pub fn print_module_content(p: &CppPrinter) -> () {
  let mod_ref: &IrTypedModule = (*p).mod;
  print_module_top p;
  print_impls p;
  print_imports p;

  let decls: Vector IrDecl = (*mod_ref).decls;
  print_decls (p, &decls);

  let m: PrinterMode = (*p).mode;
  let fns: Vector IrFunction = (*mod_ref).functions;
  let fn_len: i64 = fns.size ();

  if is_mode_public (&m) || is_mode_private (&m) {
    let i: i64 = 0;
    for i < fn_len {
      let f: IrFunction = fns.get i;
      let is_pub: bool = f.export_flag;
      if (is_mode_public (&m) && is_pub) || (is_mode_private (&m) && !is_pub) {
        print_function_signature (p, &f);
      }
      i <- i + 1;
    }
    let j: i64 = 0;
    for j < fn_len {
      let f: IrFunction = fns.get j;
      let is_pub: bool = f.export_flag;
      let f_tps: Vector (TypeParam IrType) = f.type_params;
      let is_template: bool = f_tps.size () > 0;
      if is_template && ((is_mode_public (&m) && is_pub) || (is_mode_private (&m) && !is_pub)) {
        print_function_full (p, &f);
      }
      j <- j + 1;
    }
  } else if is_mode_source (&m) {
    let i: i64 = 0;
    for i < fn_len {
      let f: IrFunction = fns.get i;
      let f_tps: Vector (TypeParam IrType) = f.type_params;
      if f_tps.size () == 0 {
        print_function_full (p, &f);
      }
      i <- i + 1;
    }
    return ();
  }

  let tis: Vector IrTraitImpl = (*mod_ref).trait_impls;
  let tin: i64 = tis.size ();
  let ti: i64 = 0;
  for ti < tin {
    let impl_item: IrTraitImpl = tis.get ti;
    print_trait_impl (p, &impl_item);
    ti <- ti + 1;
  }
  ()
}

pub fn print_module_to_string(mod_in: &IrTypedModule, mode: PrinterMode) -> String {
  let p: CppPrinter = CppPrinter::mk (mode, mod_in);
  print_module_content (&p);
  *(p.out.get ())
}
