implements bin.main

// C++ Code Generator for Bapel (bin/codegen.bpl)

fn cpp_to_header_path(module_id: &String) -> String {
  let dot: String = ".".to_string;
  let slash: String = "/".to_string;
  let path_str: String = replaceSeparator (*module_id, &dot, &slash);
  path_str.concat &".h".to_string
}

fn cpp_sanitize_id(id: &String) -> String {
  let colon_colon: String = "::".to_string;
  let dot: String = ".".to_string;
  let under: String = "_".to_string;
  let s1: String = replaceSeparator (*id, &colon_colon, &under);
  replaceSeparator (s1, &dot, &under)
}

fn cpp_indent_step(level: i64, acc: String) -> String {
  if level <= 0 {
    acc
  } else {
    let spaces: String = "  ".to_string;
    cpp_indent_step (level - 1, acc.concat &spaces)
  }
}

fn cpp_indent(level: i64) -> String {
  cpp_indent_step (level, "".to_string)
}

fn cpp_format_type(t: &String) -> String {
  if *t == "i8".to_string {
    "int8_t".to_string
  } else if *t == "i16".to_string {
    "int16_t".to_string
  } else if *t == "i32".to_string {
    "int32_t".to_string
  } else if *t == "i64".to_string {
    "int64_t".to_string
  } else if *t == "bool".to_string {
    "bool".to_string
  } else if *t == "f32".to_string {
    "float".to_string
  } else if *t == "f64".to_string {
    "double".to_string
  } else if *t == "()".to_string {
    "void".to_string
  } else {
    *t
  }
}

fn cpp_format_ptr_type(elem_type: &String) -> String {
  let formatted: String = cpp_format_type elem_type;
  formatted.concat &"*".to_string
}

fn cpp_format_ref_type(elem_type: &String) -> String {
  let formatted: String = cpp_format_type elem_type;
  let pref: String = "const ".to_string;
  (pref.concat &formatted).concat &"&".to_string
}

fn cpp_format_array_type(elem_type: &String, size: i64) -> String {
  let formatted: String = cpp_format_type elem_type;
  let pref: String = "std::array<".to_string;
  let mid: String = (pref.concat &formatted).concat &", ".to_string;
  (mid.concat &(to_string size)).concat &">".to_string
}

// Phase 9.2: Composite Types & Anonymous Types

fn cpp_format_tuple_elems_step(elems: &Vector String, index: i64, acc: String) -> String {
  if index >= elems.size {
    acc
  } else {
    let elem: String = elems.get index;
    let formatted: String = cpp_format_type &elem;
    let next_acc: String = if index == 0 {
      acc.concat &formatted
    } else {
      (acc.concat &", ".to_string).concat &formatted
    };
    cpp_format_tuple_elems_step (elems, index + 1, next_acc)
  }
}

fn cpp_format_tuple_type(elems: &Vector String) -> String {
  if elems.size == 0 {
    "std::monostate".to_string
  } else if elems.size == 1 {
    let e: String = elems.get 0;
    cpp_format_type &e
  } else {
    let list: String = cpp_format_tuple_elems_step (elems, 0, "".to_string);
    ("std::tuple<".to_string.concat &list).concat &">".to_string
  }
}

fn cpp_format_fun_type(arg_type: &String, ret_type: &String) -> String {
  let r: String = cpp_format_type ret_type;
  let a: String = cpp_format_type arg_type;
  let s1: String = ("std::function<".to_string.concat &r).concat &"(".to_string;
  (s1.concat &a).concat &")>".to_string
}

fn cpp_format_struct_fields_step(fields: &Vector IrField, index: i64, acc: String) -> String {
  if index >= fields.size {
    acc
  } else {
    let f: IrField = fields.get index;
    let f_type: String = cpp_format_type (&f.type_name);
    let f_decl: String = ((f_type.concat &" ".to_string).concat &f.name).concat &"; ".to_string;
    cpp_format_struct_fields_step (fields, index + 1, acc.concat &f_decl)
  }
}

fn cpp_format_struct_type(fields: &Vector IrField) -> String {
  let fields_str: String = cpp_format_struct_fields_step (fields, 0, "".to_string);
  let s1: String = "struct { ".to_string.concat &fields_str;
  s1.concat &"}".to_string
}

fn cpp_format_anonym_struct_name(hash_str: &String) -> String {
  "__anonym_".to_string.concat hash_str
}

fn cpp_format_params_step(params: &Vector String, index: i64, acc: String) -> String {
  if index >= params.size {
    acc
  } else {
    let p: String = params.get index;
    let next_acc: String = if index == 0 {
      acc.concat &p
    } else {
      (acc.concat &", ".to_string).concat &p
    };
    cpp_format_params_step (params, index + 1, next_acc)
  }
}

fn cpp_format_template_params_step(tvars: &Vector String, index: i64, acc: String) -> String {
  if index >= tvars.size {
    acc
  } else {
    let tvar: String = tvars.get index;
    let t_decl: String = "typename ".to_string.concat &tvar;
    let next_acc: String = if index == 0 {
      acc.concat &t_decl
    } else {
      (acc.concat &", ".to_string).concat &t_decl
    };
    cpp_format_template_params_step (tvars, index + 1, next_acc)
  }
}

fn cpp_format_app_type(fun_name: &String, type_args: &Vector String) -> String {
  let f: String = cpp_format_type fun_name;
  if type_args.size == 0 {
    f
  } else {
    let args_str: String = cpp_format_params_step (type_args, 0, "".to_string);
    let s1: String = (f.concat &"<".to_string).concat &args_str;
    s1.concat &">".to_string
  }
}

fn cpp_format_variant_type_step(tags: &Vector IrField, index: i64, acc: String) -> String {
  if index >= tags.size {
    acc
  } else {
    let pair: IrField = tags.get index;
    let formatted_type: String = cpp_format_type &pair.type_name;
    let next_acc: String = if index == 0 {
      acc.concat &formatted_type
    } else {
      (acc.concat &", ".to_string).concat &formatted_type
    };
    cpp_format_variant_type_step (tags, index + 1, next_acc)
  }
}

fn cpp_format_variant_type(tags: &Vector IrField) -> String {
  let list: String = cpp_format_variant_type_step (tags, 0, "".to_string);
  ("std::variant<".to_string.concat &list).concat &">".to_string
}

fn cpp_format_ir_type(t: &IrType) -> String {
  match *t {
    name n => cpp_format_type (&n),
    app a => cpp_format_app_type (&a.0, &a.1),
    array a => cpp_format_array_type (&a.0, a.1),
    fun f => cpp_format_fun_type (&f.0, &f.1),
    tuple elems => cpp_format_tuple_type (&elems),
    variant_type tags => cpp_format_variant_type (&tags),
    struct_type fields => cpp_format_struct_type (&fields),
    ptr elem => cpp_format_ptr_type (&elem),
    ref elem => cpp_format_ref_type (&elem)
  }
}

fn cpp_format_function_signature(
    fn_name: &String,
    return_type: &String,
    params: &Vector String,
    type_params: &Vector String) -> String {
  let ret: String = cpp_format_type return_type;
  let param_list: String = cpp_format_params_step (params, 0, "".to_string);
  let base_sig: String = ((ret.concat &" ".to_string).concat fn_name).concat &"(".to_string;
  let full_sig: String = (base_sig.concat &param_list).concat &")".to_string;

  if type_params.size == 0 {
    full_sig
  } else {
    let t_list: String = cpp_format_template_params_step (type_params, 0, "".to_string);
    let t_head: String = ("template <".to_string.concat &t_list).concat &">\n".to_string;
    t_head.concat &full_sig
  }
}

// Phase 3: Header Generation Routines

fn cpp_emit_std_includes() -> String {
  let res: String = "#pragma once\n\n".to_string;
  let r1: String = res.concat &"#include <array>\n#include <cstdlib>\n#include <cmath>\n".to_string;
  let r2: String = r1.concat &"#include <functional>\n#include <optional>\n#include <string>\n".to_string;
  r2.concat &"#include <tuple>\n#include <variant>\n#include <vector>\n\n".to_string
}

fn cpp_emit_import_includes_step(imp_modules: &Vector String, index: i64, acc: String) -> String {
  if index >= imp_modules.size {
    acc
  } else {
    let imp: String = imp_modules.get index;
    let h_path: String = cpp_to_header_path (&imp);
    let inc_line: String = ("#include \"".to_string.concat &h_path).concat &"\"\n".to_string;
    cpp_emit_import_includes_step (imp_modules, index + 1, acc.concat &inc_line)
  }
}

fn cpp_emit_import_includes(imp_modules: &Vector String) -> String {
  cpp_emit_import_includes_step (imp_modules, 0, "".to_string)
}

// Emits trait base template DECLARATION ONLY (rule: never define base template)
fn cpp_emit_trait_base_decl(trait_name: &String, type_params: &Vector String) -> String {
  let self_param: String = "typename Self".to_string;
  let t_list: String = if type_params.size == 0 {
    self_param
  } else {
    let rest: String = cpp_format_template_params_step (type_params, 0, "".to_string);
    let comma: String = ", ".to_string;
    (self_param.concat &comma).concat &rest
  };
  let t_open: String = "template <".to_string;
  let t_close: String = ">\nstruct ".to_string;
  let t_head: String = (t_open.concat &t_list).concat &t_close;
  let s_end: String = ";\n\n".to_string;
  (t_head.concat trait_name).concat &s_end
}

// Phase 4: Source File Expression & Statement Code Generation

fn cpp_emit_let(var_name: &String, var_type: &String, expr_val: &String, indent_level: i64) -> String {
  let ind: String = cpp_indent indent_level;
  let formatted_type: String = cpp_format_type var_type;
  let line: String = (((ind.concat &formatted_type).concat &" ".to_string).concat var_name).concat &" = ".to_string;
  (line.concat expr_val).concat &";\n".to_string
}

fn cpp_emit_assign(var_name: &String, expr_val: &String, indent_level: i64) -> String {
  let ind: String = cpp_indent indent_level;
  let line: String = (ind.concat var_name).concat &" = ".to_string;
  (line.concat expr_val).concat &";\n".to_string
}

fn cpp_emit_return(expr_val: &String, indent_level: i64) -> String {
  let ind: String = cpp_indent indent_level;
  let line: String = ind.concat &"return ".to_string;
  (line.concat expr_val).concat &";\n".to_string
}

fn cpp_emit_if_head(cond: &String, indent_level: i64) -> String {
  let ind: String = cpp_indent indent_level;
  let line: String = ind.concat &"if (".to_string;
  (line.concat cond).concat &") {\n".to_string
}

fn cpp_emit_if_term(cond: &String, then_term: &String, else_term: &String, indent_level: i64) -> String {
  let head: String = cpp_emit_if_head (cond, indent_level);
  let ind: String = cpp_indent (indent_level + 1);
  let then_line: String = (ind.concat then_term).concat &";\n".to_string;
  let close_if: String = (cpp_indent indent_level).concat &"}\n".to_string;
  if (*else_term).size == 0 {
    (head.concat &then_line).concat &close_if
  } else {
    let else_head: String = (cpp_indent indent_level).concat &"} else {\n".to_string;
    let else_line: String = (ind.concat else_term).concat &";\n".to_string;
    let s1: String = (head.concat &then_line).concat &else_head;
    (s1.concat &else_line).concat &close_if
  }
}

fn cpp_emit_for_loop(var_name: &String, start_val: &String, end_val: &String, body: &String, indent_level: i64) -> String {
  let ind: String = cpp_indent indent_level;
  let p1: String = (ind.concat &"for (int64_t ".to_string).concat var_name;
  let p2: String = ((p1.concat &" = ".to_string).concat start_val).concat &"; ".to_string;
  let p3: String = (((p2.concat var_name).concat &" < ".to_string).concat end_val).concat &"; ++".to_string;
  let head: String = (p3.concat var_name).concat &") {\n".to_string;
  let ind_inner: String = cpp_indent (indent_level + 1);
  let b_line: String = (ind_inner.concat body).concat &";\n".to_string;
  let close: String = ind.concat &"}\n".to_string;
  (head.concat &b_line).concat &close
}

fn cpp_emit_match_arm_step(
    var_name: &String,
    arms: &Vector MatchArm,
    index: i64,
    indent_level: i64,
    acc: String) -> String {
  if index >= arms.size {
    acc
  } else {
    let arm: MatchArm = arms.get index;
    let arm_idx: i64 = arm.index;
    let arg_id: String = arm.arg_id;
    let body_str: String = arm.body;

    let ind: String = cpp_indent indent_level;
    let ind_inner: String = cpp_indent (indent_level + 1);

    let idx_str: String = to_string arm_idx;
    let c1: String = ind.concat &"case ".to_string;
    let c2: String = c1.concat &idx_str;
    let case_head: String = c2.concat &": {\n".to_string;

    let g1: String = ind_inner.concat &"auto &".to_string;
    let g2: String = (g1.concat &arg_id).concat &" = std::get<".to_string;
    let g3: String = (g2.concat &idx_str).concat &">(".to_string;
    let get_line: String = (g3.concat var_name).concat &");\n".to_string;

    let body_line: String = (ind_inner.concat &body_str).concat &";\n".to_string;
    let case_tail: String = ind.concat &"}\n".to_string;

    let a1: String = case_head.concat &get_line;
    let a2: String = a1.concat &body_line;
    let arm_code: String = a2.concat &case_tail;
    let next_acc: String = acc.concat &arm_code;
    cpp_emit_match_arm_step (var_name, arms, index + 1, indent_level, next_acc)
  }
}

fn cpp_emit_match(var_name: &String, arms: &Vector MatchArm, indent_level: i64) -> String {
  let ind: String = cpp_indent indent_level;
  let h1: String = ind.concat &"switch (".to_string;
  let head: String = h1.concat var_name;
  let switch_open: String = head.concat &".index()) {\n".to_string;
  let arms_code: String = cpp_emit_match_arm_step (var_name, arms, 0, indent_level + 1, "".to_string);
  let switch_close: String = ind.concat &"}\n".to_string;
  let res1: String = switch_open.concat &arms_code;
  res1.concat &switch_close
}

fn cpp_emit_sfinae_constraint(trait_name: &String, type_param: &String) -> String {
  let s1: String = "std::enable_if_t<(sizeof(".to_string;
  let s2: String = s1.concat trait_name;
  let s3: String = s2.concat &"<".to_string;
  let s4: String = s3.concat type_param;
  s4.concat &">) > 0), int> = 0".to_string
}

fn cpp_emit_namespace_start(ns: &String) -> String {
  let n: String = cpp_sanitize_id ns;
  let ns_pref: String = "namespace ".to_string.concat &n;
  ns_pref.concat &" {\n\n".to_string
}

fn cpp_emit_namespace_end(ns: &String) -> String {
  let n: String = cpp_sanitize_id ns;
  let ns_pref: String = "}\n// namespace ".to_string.concat &n;
  ns_pref.concat &"\n\n".to_string
}

// Phase 9.2: Complete Term Lowering & Expression Generation

fn cpp_format_args_step(args: &Vector String, index: i64, acc: String) -> String {
  if index >= args.size {
    acc
  } else {
    let arg: String = args.get index;
    let next_acc: String = if index == 0 {
      acc.concat &arg
    } else {
      (acc.concat &", ".to_string).concat &arg
    };
    cpp_format_args_step (args, index + 1, next_acc)
  }
}

fn cpp_format_args(args: &Vector String) -> String {
  cpp_format_args_step (args, 0, "".to_string)
}

fn cpp_emit_app_type_term(term_str: &String, type_args: &Vector String) -> String {
  if type_args.size == 0 {
    *term_str
  } else {
    let t_list: String = cpp_format_params_step (type_args, 0, "".to_string);
    let s1: String = (term_str.concat &"<".to_string).concat &t_list;
    s1.concat &">".to_string
  }
}

fn cpp_emit_app_term(fn_name: &String, args: &Vector String) -> String {
  let args_str: String = cpp_format_args args;
  let s1: String = (fn_name.concat &"(".to_string).concat &args_str;
  s1.concat &")".to_string
}

fn cpp_emit_projection(term_str: &String, field_or_index: &String, is_struct: bool) -> String {
  if is_struct {
    (term_str.concat &".".to_string).concat field_or_index
  } else {
    let s1: String = "std::get<".to_string.concat field_or_index;
    let s2: String = (s1.concat &">(".to_string).concat term_str;
    s2.concat &")".to_string
  }
}

fn cpp_emit_injection(variant_type: &String, index: i64, value_str: &String) -> String {
  let vtype: String = cpp_format_type variant_type;
  let idx_str: String = to_string index;
  let s1: String = (vtype.concat &"(std::in_place_index<".to_string).concat &idx_str;
  let s2: String = (s1.concat &">, ".to_string).concat value_str;
  s2.concat &")".to_string
}

fn cpp_emit_set(term_str: &String, field: &String, val_str: &String, is_struct: bool) -> String {
  if is_struct {
    let s1: String = "([&, __v_0 = ".to_string.concat term_str;
    let s2: String = s1.concat &"]() mutable {\n  __v_0.".to_string;
    let s3: String = ((s2.concat field).concat &" = ".to_string).concat val_str;
    s3.concat &";\n  return __v_0;\n})()".to_string
  } else {
    let s1: String = "([__v_0 = ".to_string.concat term_str;
    let s2: String = s1.concat &"]() mutable {\n  std::get<".to_string;
    let s3: String = (s2.concat field).concat &">(__v_0) = ".to_string;
    let s4: String = (s3.concat val_str).concat &";\n  return __v_0;\n})()".to_string;
    s4
  }
}

fn cpp_emit_lambda(type_params: &Vector String, params: &Vector String, body: &String) -> String {
  let t_head: String = if type_params.size == 0 {
    "".to_string
  } else {
    let t_list: String = cpp_format_template_params_step (type_params, 0, "".to_string);
    ("<".to_string.concat &t_list).concat &">".to_string
  };
  let p_list: String = cpp_format_params_step (params, 0, "".to_string);
  let l_head: String = (("[&]".to_string.concat &t_head).concat &"(".to_string).concat &p_list;
  let l_mid: String = (l_head.concat &") {\n  return ".to_string).concat body;
  l_mid.concat &";\n}".to_string
}

fn cpp_emit_tuple(elems: &Vector String) -> String {
  if elems.size == 0 {
    "std::monostate()".to_string
  } else {
    let list: String = cpp_format_args elems;
    ("std::make_tuple(".to_string.concat &list).concat &")".to_string
  }
}

fn cpp_emit_struct_fields_init_step(fields: &Vector IrField, index: i64, acc: String) -> String {
  if index >= fields.size {
    acc
  } else {
    let f: IrField = fields.get index;
    let init_str: String = ((".".to_string.concat &f.name).concat &" = ".to_string).concat &f.type_name;
    let next_acc: String = if index == 0 {
      acc.concat &init_str
    } else {
      (acc.concat &", ".to_string).concat &init_str
    };
    cpp_emit_struct_fields_init_step (fields, index + 1, next_acc)
  }
}

fn cpp_emit_struct(fields: &Vector IrField) -> String {
  let fields_init: String = cpp_emit_struct_fields_init_step (fields, 0, "".to_string);
  ("{ ".to_string.concat &fields_init).concat &" }".to_string
}

fn cpp_emit_block_lines_step(lines: &Vector String, index: i64, indent_level: i64, acc: String) -> String {
  if index >= lines.size {
    acc
  } else {
    let line: String = lines.get index;
    let ind: String = cpp_indent indent_level;
    let line_str: String = (ind.concat &line).concat &";\n".to_string;
    cpp_emit_block_lines_step (lines, index + 1, indent_level, acc.concat &line_str)
  }
}

fn cpp_emit_block(lines: &Vector String, indent_level: i64) -> String {
  let ind: String = cpp_indent indent_level;
  let b_open: String = ind.concat &"{\n".to_string;
  let body: String = cpp_emit_block_lines_step (lines, 0, indent_level + 1, "".to_string);
  let b_close: String = ind.concat &"}\n".to_string;
  (b_open.concat &body).concat &b_close
}

fn cpp_emit_term(term: &IrTerm, indent_level: i64) -> String {
  match *term {
    var id => id,
    const_int i => to_string i,
    const_float f => to_string f,
    const_str s => ("\"".to_string.concat &s).concat &"\"".to_string,
    const_bool b => if b { "true".to_string } else { "false".to_string },
    let_term t => cpp_emit_let (&t.0, &t.1, &t.2, indent_level),
    assign_term t => cpp_emit_assign (&t.0, &t.1, indent_level),
    return_term s => cpp_emit_return (&s, indent_level),
    block lines => cpp_emit_block (&lines, indent_level),
    if_term t => cpp_emit_if_term (&t.0, &t.1, &t.2, indent_level),
    for_loop t => cpp_emit_for_loop (&t.0, &t.1, &t.2, &t.3, indent_level),
    match_term t => cpp_emit_match (&t.0, &t.1, indent_level),
    app_type_term t => cpp_emit_app_type_term (&t.0, &t.1),
    app_term t => cpp_emit_app_term (&t.0, &t.1),
    projection_term t => cpp_emit_projection (&t.0, &t.1, t.2),
    injection_term t => cpp_emit_injection (&t.0, t.1, &t.2),
    set_term t => cpp_emit_set (&t.0, &t.1, &t.2, t.3),
    lambda_term t => cpp_emit_lambda (&t.0, &t.1, &t.3),
    tuple_term elems => cpp_emit_tuple (&elems),
    struct_term fields => cpp_emit_struct (&fields)
  }
}

// Phase 6.4: Direct File I/O Integration Routines

fn cpp_write_file(path: &String, content: &String) -> bool {
  let parent_dir: String = fs::parent_path (*path);
  if parent_dir.size > 0 {
    if !fs::create_directories parent_dir {
      return false
    }
  };
  let f: Ofstream = Ofstream::open path;
  if !f.is_open {
    return false
  }
  f.write (*content);
  f.close;
  true
}

fn cpp_emit_public_header_content(module_id: &String, imp_modules: &Vector String) -> String {
  let std_inc: String = cpp_emit_std_includes ();
  let imp_inc: String = cpp_emit_import_includes imp_modules;
  let ns_start: String = cpp_emit_namespace_start module_id;
  let ns_end: String = cpp_emit_namespace_end module_id;
  let h_decl: String = (std_inc.concat &imp_inc).concat &ns_start;
  h_decl.concat &ns_end
}

fn cpp_emit_private_header_content(module_id: &String) -> String {
  let pragma: String = "#pragma once\n\n".to_string;
  let h_path: String = cpp_to_header_path module_id;
  let inc_pub: String = ("#include \"".to_string.concat &h_path).concat &"\"\n\n".to_string;
  let ns_start: String = cpp_emit_namespace_start module_id;
  let ns_end: String = cpp_emit_namespace_end module_id;
  let content1: String = (pragma.concat &inc_pub).concat &ns_start;
  content1.concat &ns_end
}

fn cpp_emit_source_content(module_id: &String, imp_modules: &Vector String) -> String {
  let h_path: String = cpp_to_header_path module_id;
  let inc_pub: String = ("#include \"".to_string.concat &h_path).concat &"\"\n".to_string;
  let dot: String = ".".to_string;
  let slash: String = "/".to_string;
  let base_path: String = replaceSeparator (*module_id, &dot, &slash);
  let priv_h_path: String = base_path.concat &"_private.h".to_string;
  let inc_priv: String = ("#include \"".to_string.concat &priv_h_path).concat &"\"\n\n".to_string;
  inc_pub.concat &inc_priv
}

fn cpp_write_module_files(out_dir: &String, module_id: &String, imp_modules: &Vector String) -> bool {
  let dot: String = ".".to_string;
  let slash: String = "/".to_string;
  let base_name: String = replaceSeparator (*module_id, &dot, &slash);
  let out_path: String = fs::join (*out_dir, base_name);

  let header_path: String = out_path.concat &".h".to_string;
  let priv_header_path: String = out_path.concat &"_private.h".to_string;
  let source_path: String = out_path.concat &".cc".to_string;

  let pub_content: String = cpp_emit_public_header_content (module_id, imp_modules);
  let priv_content: String = cpp_emit_private_header_content module_id;
  let src_content: String = cpp_emit_source_content (module_id, imp_modules);

  if !cpp_write_file (&header_path, &pub_content) {
    return false
  };
  if !cpp_write_file (&priv_header_path, &priv_content) {
    return false
  };
  if !cpp_write_file (&source_path, &src_content) {
    return false
  };
  true
}

fn cpp_emit_decl(d: &IrDecl, mode: i64) -> String {
  let is_pub: bool = (*d).is_export;
  if mode == 2 {
    "".to_string
  } else if mode == 0 && !is_pub {
    "".to_string
  } else if mode == 1 && is_pub {
    "".to_string
  } else {
    let k: String = (*d).decl_kind;
    if k.size == 0 {
      "".to_string
    } else {
      let semi: String = ";\n".to_string;
      let line: String = if String::ends_with (&k, &semi) {
        k
      } else {
        let semi2: String = ";".to_string;
        if String::ends_with (&k, &semi2) {
          k.concat &"\n".to_string
        } else {
          k.concat &";\n".to_string
        }
      };
      line
    }
  }
}

fn cpp_emit_decls_step(decls: &Vector IrDecl, index: i64, mode: i64, acc: String) -> String {
  if index >= (*decls).size {
    acc
  } else {
    let d: IrDecl = Vector::get [IrDecl] (decls, index);
    let s: String = cpp_emit_decl (&d, mode);
    let next_acc: String = acc.concat &s;
    cpp_emit_decls_step (decls, index + 1, mode, next_acc)
  }
}

fn cpp_emit_decls(decls: &Vector IrDecl, mode: i64) -> String {
  cpp_emit_decls_step (decls, 0, mode, "".to_string)
}

fn cpp_emit_trait_methods_step(methods: &Vector String, index: i64, acc: String) -> String {
  if index >= (*methods).size {
    acc
  } else {
    let m: String = Vector::get [String] (methods, index);
    let ind: String = "  ".to_string;
    let semi: String = ";\n".to_string;
    let m_line: String = (ind.concat &m).concat &semi;
    let next_acc: String = acc.concat &m_line;
    cpp_emit_trait_methods_step (methods, index + 1, next_acc)
  }
}

fn cpp_emit_trait_impl(trait_impl: &IrTraitImpl, mode: i64) -> String {
  if mode == 2 {
    "".to_string
  } else {
    let t_name: String = (*trait_impl).trait_name;
    let type_nm: String = (*trait_impl).type_name;
    let formatted_type: String = cpp_format_type &type_nm;
    let methods_code: String = cpp_emit_trait_methods_step (&(*trait_impl).methods, 0, "".to_string);
    
    if t_name.size == 0 || t_name == "inherent".to_string {
      // Inherent impl
      let t_head: String = if (*trait_impl).type_params.size == 0 {
        "".to_string
      } else {
        let t_list: String = cpp_format_template_params_step (&(*trait_impl).type_params, 0, "".to_string);
        ("template <".to_string.concat &t_list).concat &">\n".to_string
      };
      let s1: String = (t_head.concat &"struct inherents::".to_string).concat &type_nm;
      let s2: String = (s1.concat &" {\n  using Self = ".to_string).concat &formatted_type;
      let s3: String = (s2.concat &";\n".to_string).concat &methods_code;
      s3.concat &"};\n\n".to_string
    } else {
      // Trait specialization
      let t_head: String = if (*trait_impl).type_params.size == 0 {
        "template <>\n".to_string
      } else {
        let t_list: String = cpp_format_template_params_step (&(*trait_impl).type_params, 0, "".to_string);
        ("template <".to_string.concat &t_list).concat &">\n".to_string
      };
      let s1: String = (t_head.concat &"struct traits::".to_string).concat &t_name;
      let s2: String = ((s1.concat &"<".to_string).concat &formatted_type).concat &"> {\n  using Self = ".to_string;
      let s3: String = ((s2.concat &formatted_type).concat &";\n".to_string).concat &methods_code;
      s3.concat &"};\n\n".to_string
    }
  }
}

fn cpp_emit_trait_impls_step(trait_impls_list: &Vector IrTraitImpl, index: i64, mode: i64, acc: String) -> String {
  if index >= (*trait_impls_list).size {
    acc
  } else {
    let impl_item: IrTraitImpl = Vector::get [IrTraitImpl] (trait_impls_list, index);
    let s: String = cpp_emit_trait_impl (&impl_item, mode);
    let next_acc: String = acc.concat &s;
    cpp_emit_trait_impls_step (trait_impls_list, index + 1, mode, next_acc)
  }
}

fn cpp_emit_trait_impls(trait_impls_list: &Vector IrTraitImpl, mode: i64) -> String {
  cpp_emit_trait_impls_step (trait_impls_list, 0, mode, "".to_string)
}

// Phase 6.6 & 9.1: Native Function Transpilation & Self-Bootstrapping Routines

fn cpp_emit_function(f: &IrFunction, mode: i64) -> String {
  let is_pub: bool = (*f).is_export;
  if mode == 0 && !is_pub {
    "".to_string
  } else if mode == 1 && is_pub {
    "".to_string
  } else {
    let sig: String = cpp_format_function_signature (&(*f).name, &(*f).ret_type, &(*f).params, &(*f).type_params);
    if mode == 2 {
      ((sig.concat &" {\n".to_string).concat &(cpp_indent_step (1, (*f).body))).concat &"\n}\n\n".to_string
    } else {
      sig.concat &";\n".to_string
    }
  }
}

fn cpp_emit_functions_step(funcs: &Vector IrFunction, index: i64, mode: i64, acc: String) -> String {
  if index >= (*funcs).size {
    acc
  } else {
    let f: IrFunction = Vector::get [IrFunction] (funcs, index);
    let s: String = cpp_emit_function (&f, mode);
    let next_acc: String = acc.concat &s;
    cpp_emit_functions_step (funcs, index + 1, mode, next_acc)
  }
}

fn cpp_emit_functions(funcs: &Vector IrFunction, mode: i64) -> String {
  cpp_emit_functions_step (funcs, 0, mode, "".to_string)
}

// CodegenMode flags: 0 = public_header, 1 = private_header, 2 = source_file

fn cpp_emit_unit(unit: &IrUnit, mode: i64) -> String {
  let ns_start: String = cpp_emit_namespace_start (&(*unit).module_id);
  let ns_end: String = cpp_emit_namespace_end (&(*unit).module_id);
  let decl_content: String = cpp_emit_decls (&(*unit).decls, mode);
  let fn_content: String = cpp_emit_functions (&(*unit).functions, mode);
  let trait_content: String = cpp_emit_trait_impls (&(*unit).trait_impls, mode);

  let body_content: String = (decl_content.concat &fn_content).concat &trait_content;

  if mode == 0 {
    let std_inc: String = cpp_emit_std_includes ();
    let imp_inc: String = cpp_emit_import_includes (&(*unit).import_modules);
    let head: String = (std_inc.concat &imp_inc).concat &ns_start;
    (head.concat &body_content).concat &ns_end
  } else if mode == 1 {
    let pragma: String = "#pragma once\n\n".to_string;
    let h_path: String = cpp_to_header_path (&(*unit).module_id);
    let inc_pub: String = ("#include \"".to_string.concat &h_path).concat &"\"\n\n".to_string;
    let head: String = (pragma.concat &inc_pub).concat &ns_start;
    (head.concat &body_content).concat &ns_end
  } else {
    let h_path: String = cpp_to_header_path (&(*unit).module_id);
    let inc_pub: String = ("#include \"".to_string.concat &h_path).concat &"\"\n".to_string;
    let dot: String = ".".to_string;
    let slash: String = "/".to_string;
    let base_path: String = replaceSeparator ((*unit).module_id, &dot, &slash);
    let priv_h_path: String = base_path.concat &"_private.h".to_string;
    let inc_priv: String = ("#include \"".to_string.concat &priv_h_path).concat &"\"\n\n".to_string;
    let head: String = (inc_pub.concat &inc_priv).concat &ns_start;
    (head.concat &body_content).concat &ns_end
  }
}

fn cpp_write_module_unit(out_dir: &String, unit: &IrUnit) -> bool {
  let dot: String = ".".to_string;
  let slash: String = "/".to_string;
  let base_name: String = replaceSeparator ((*unit).module_id, &dot, &slash);
  let out_path: String = fs::join (*out_dir, base_name);

  let header_path: String = out_path.concat &".h".to_string;
  let priv_header_path: String = out_path.concat &"_private.h".to_string;
  let source_path: String = out_path.concat &".cc".to_string;

  let pub_content: String = cpp_emit_unit (unit, 0);
  let priv_content: String = cpp_emit_unit (unit, 1);
  let src_content: String = cpp_emit_unit (unit, 2);

  if !cpp_write_file (&header_path, &pub_content) {
    return false
  }
  if !cpp_write_file (&priv_header_path, &priv_content) {
    return false
  }
  if !cpp_write_file (&source_path, &src_content) {
    return false
  }
  true
}
