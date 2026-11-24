implements bapel.compiler.ir

imports {
  bapel.core
  bapel.stl
}

pub type IrKind = variant {
  type_kind (),
  arrow_kind struct {
    left: SharedPtr IrKind,
    right: SharedPtr IrKind
  }
}

impl IrKind {
  fn new_type() -> IrKind {
    variant { IrKind type_kind = () }
  }

  fn new_arrow(l: IrKind, r: IrKind) -> IrKind {
    variant {
      IrKind arrow_kind = struct {
        left = SharedPtr::make_shared [IrKind] l,
        right = SharedPtr::make_shared [IrKind] r
      }
    }
  }

  fn is_type_kind(k: &Self) -> bool {
    match *k {
      type_kind u => true,
      _ => false
    }
  }

  fn is_arrow_kind(k: &Self) -> bool {
    match *k {
      arrow_kind a => true,
      _ => false
    }
  }

  fn equals(a: &Self, b: &IrKind) -> bool {
    ir_kind_equals (a, b)
  }

  fn to_string(k: &Self) -> String {
    match *k {
      type_kind u => String::from_view "*",
      arrow_kind a => {
        let s: String = String::from_view "";
        if (a.left.get()).is_arrow_kind() {
          s.append "(";
          let ls: String = (a.left.get()).to_string();
          s.append (ls.view());
          s.append ")";
        } else {
          let ls: String = (a.left.get()).to_string();
          s.append (ls.view());
        }
        s.append " -> ";
        let rs: String = (a.right.get()).to_string();
        s.append (rs.view());
        s
      }
    }
  }
}

pub fn ir_kind_equals(a: &IrKind, b: &IrKind) -> bool {
  match *a {
    type_kind u => match *b {
      type_kind v => true,
      _ => false
    },
    arrow_kind ka => match *b {
      arrow_kind kb =>
        ir_kind_equals (ka.left.get(), kb.left.get()) &&
        ir_kind_equals (ka.right.get(), kb.right.get()),
      _ => false
    }
  }
}
