implements bapel.compiler.codegen.cpp

imports {
  bapel.core
  bapel.stl
  bapel.compiler.ir
}

pub type AnonymousType = struct {
  name_type: IrType,
  alias_decl: IrDecl
}

pub fn to_header_path(module_id: &ModuleID) -> String {
  let out: String = String::from_view "";
  let s: &String = &((*module_id).name);
  let n: i64 = s.size ();
  let i: i64 = 0;
  for i < n {
    let c: i8 = s.at i;
    if c == '.' {
      out.push_back '/';
    } else {
      out.push_back c;
    }
    i <- i + 1;
  }
  out.append ".h";
  out
}

pub fn inherent_cpp_name(bapel_name: &String) -> String {
  let sep: String = String::from_view "::";
  let pos: i64 = bapel_name.rfind (&sep);
  if pos < 0 {
    let res: String = String::from_view "inherents::";
    res.append (bapel_name.view ());
    return res;
  }
  let pfx: String = bapel_name.substr (0, pos);
  let name: String = bapel_name.substr (pos + 2, bapel_name.size () - (pos + 2));
  let res: String = pfx;
  res.append "::inherents::";
  res.append (name.view ());
  res
}

pub fn trait_cpp_name(bapel_name: &String) -> String {
  let sep: String = String::from_view "::";
  let pos: i64 = bapel_name.rfind (&sep);
  if pos < 0 {
    let res: String = String::from_view "traits::";
    res.append (bapel_name.view ());
    return res;
  }
  let pfx: String = bapel_name.substr (0, pos);
  let name: String = bapel_name.substr (pos + 2, bapel_name.size () - (pos + 2));
  let res: String = pfx;
  res.append "::traits::";
  res.append (name.view ());
  res
}

pub fn type_app_fun(typ: &IrType) -> IrType {
  match (*typ).data {
    app_t a => type_app_fun (a.fun.get ()),
    _ => *typ
  }
}

pub fn base_type_name(typ: &IrType) -> String {
  match (*typ).data {
    name_t n => n,
    app_t _ => {
      let fun: IrType = type_app_fun typ;
      base_type_name (&fun)
    },
    _ => String::from_view ""
  }
}

pub fn count_type_vars(kind: &IrKind) -> i64 {
  match *kind {
    arrow_kind ak => 1 + count_type_vars (ak.right.get ()),
    type_kind _ => 0
  }
}

pub fn count_alias_type_vars(typ: &IrType) -> i64 {
  match (*typ).data {
    lambda_t l => 1 + count_alias_type_vars (l.body.get ()),
    _ => 0
  }
}

pub fn is_operator(id: &String) -> bool {
  let v: StringView = id.view ();
  if v == "+" { return true; }
  if v == "-" { return true; }
  if v == "*" { return true; }
  if v == "/" { return true; }
  if v == "%" { return true; }
  if v == "==" { return true; }
  if v == "!=" { return true; }
  if v == "<" { return true; }
  if v == "<=" { return true; }
  if v == ">" { return true; }
  if v == ">=" { return true; }
  if v == "&&" { return true; }
  if v == "||" { return true; }
  if v == "!" { return true; }
  if v == "&" { return true; }
  if v == "|" { return true; }
  if v == "^" { return true; }
  if v == "~" { return true; }
  if v == "<<" { return true; }
  if v == ">>" { return true; }
  false
}

pub fn is_cpp_statement(term: &IrTerm) -> bool {
  match (*term).data {
    assign _ => true,
    block _ => true,
    let_t _ => true,
    match_t _ => true,
    return_t _ => true,
    app_term at => {
      let fun_term: &IrTerm = at.fun.get ();
      match (*fun_term).data {
        var_t v => {
          let sv: StringView = v.view ();
          if sv == "ifthen" { return true; }
          if sv == "ifelse" { return true; }
          if sv == "core::for" { return true; }
          false
        },
        app_type ap => {
          let inner: &IrTerm = ap.fun.get ();
          match (*inner).data {
            var_t v => {
              let sv: StringView = v.view ();
              if sv == "ifthen" { return true; }
              if sv == "ifelse" { return true; }
              if sv == "core::for" { return true; }
              false
            },
            _ => false
          }
        },
        _ => false
      }
    },
    _ => false
  }
}
