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
