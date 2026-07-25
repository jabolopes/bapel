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
  let under: String = "_".to_string;
  replaceSeparator (*id, &colon_colon, &under)
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

fn cpp_emit_for_loop(var_name: &String, start_val: &String, end_val: &String, indent_level: i64) -> String {
  let ind: String = cpp_indent indent_level;
  let p1: String = (ind.concat &"for (int64_t ".to_string).concat var_name;
  let p2: String = ((p1.concat &" = ".to_string).concat start_val).concat &"; ".to_string;
  let p3: String = (((p2.concat var_name).concat &" < ".to_string).concat end_val).concat &"; ++".to_string;
  ((p3.concat var_name).concat &") {\n".to_string)
}

// Phase 6.3: Core Traversal & Complex Type Printing Routines

fn cpp_format_fun_type(arg_type: &String, ret_type: &String) -> String {
  let formatted_arg: String = cpp_format_type arg_type;
  let formatted_ret: String = cpp_format_type ret_type;
  let p1: String = ("std::function<".to_string.concat &formatted_ret).concat &"(".to_string;
  (p1.concat &formatted_arg).concat &")>".to_string
}

fn cpp_format_tuple_type_step(elems: &Vector String, index: i64, acc: String) -> String {
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
    cpp_format_tuple_type_step (elems, index + 1, next_acc)
  }
}

fn cpp_format_tuple_type(elems: &Vector String) -> String {
  if elems.size == 0 {
    "std::monostate".to_string
  } else if elems.size == 1 {
    cpp_format_type (&(elems.get 0))
  } else {
    let list: String = cpp_format_tuple_type_step (elems, 0, "".to_string);
    ("std::tuple<".to_string.concat &list).concat &">".to_string
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
  let ns_pref: String = "// namespace ".to_string.concat &n;
  ns_pref.concat &"\n\n".to_string
}


