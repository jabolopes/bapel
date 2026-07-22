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
