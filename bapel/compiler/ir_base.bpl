implements bapel.compiler.ir

imports {
  bapel.core
  bapel.stl
}

pub type Pos = struct {
  filename: String,
  begin_line: i64,
  end_line: i64
}

impl Pos {
  fn empty() -> Pos {
    struct {
      filename = String::from_view "",
      begin_line = 0,
      end_line = 0
    }
  }

  fn mk(filename: String, begin_line: i64, end_line: i64) -> Pos {
    struct {
      filename = filename,
      begin_line = begin_line,
      end_line = end_line
    }
  }

  fn is_empty(p: &Self) -> bool {
    (*p).filename.empty() && (*p).begin_line == 0 && (*p).end_line == 0
  }
}

pub type ModuleID = struct {
  name: String,
  pos: Pos
}

impl ModuleID {
  fn mk(name: String) -> ModuleID {
    struct { name = name, pos = Pos::empty () }
  }
}

pub type Filename = struct {
  value: String,
  pos: Pos
}

pub type IrImport = struct {
  module_id: ModuleID,
  pos: Pos
}

pub type IrImpl = struct {
  relative_filename: Filename,
  pos: Pos
}

pub type FloatLit = struct {
  integer: i64,
  decimal: i64
}

pub type IrLiteral = variant {
  int_lit i64,
  float_lit FloatLit,
  rune_lit String,
  str_lit String
}

impl IrLiteral {
  fn is_int(l: &Self) -> bool {
    match *l {
      int_lit v => true,
      _ => false
    }
  }

  fn is_float(l: &Self) -> bool {
    match *l {
      float_lit v => true,
      _ => false
    }
  }

  fn is_rune(l: &Self) -> bool {
    match *l {
      rune_lit v => true,
      _ => false
    }
  }

  fn is_str(l: &Self) -> bool {
    match *l {
      str_lit v => true,
      _ => false
    }
  }
}

pub type PrinterMode = variant {
  mode_public_header (),
  mode_private_header (),
  mode_source ()
}

pub type Position = variant {
  type_position (),
  bind_position ()
}

pub fn i64_to_string(n: i64) -> String {
  if n == 0 {
    return String::from_view "0"
  }
  let neg: bool = n < 0;
  let v: i64 = if neg { 0 - n } else { n };
  let digits: String = String::from_view "0123456789";
  let rev: String = String::from_view "";
  for v > 0 {
    let d: i64 = v - (v / 10) * 10;
    rev.push_back (digits.at d);
    v <- v / 10;
  }
  let out: String = String::from_view "";
  if neg {
    out.push_back '-';
  }
  let i: i64 = rev.size () - 1;
  for i >= 0 {
    out.push_back (rev.at i);
    i <- i - 1;
  }
  out
}
