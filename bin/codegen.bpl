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
