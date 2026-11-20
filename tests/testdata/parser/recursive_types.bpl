module /* tests/testdata/in/recursive_types.in:2 */recursive_types

in "tests/testdata/in/recursive_types.in" in lines 4-7
imports {
  /* tests/testdata/in/recursive_types.in:5 */bapel.core
  /* tests/testdata/in/recursive_types.in:6 */bapel.stl
}

/* tests/testdata/in/recursive_types.in:10-13 */ pub type Node = struct{val: i64, next: SharedPtr Node}

in "tests/testdata/in/recursive_types.in" in lines 15-20 pub fn mkNode(v: i64) -> Node {
  struct{val = v, next = SharedPtr::null [Node] ()}
}

/* tests/testdata/in/recursive_types.in:23-26 */ pub type IntList = variant{nil (), cons (i64, SharedPtr IntList)}

in "tests/testdata/in/recursive_types.in" in lines 28-33 pub fn sumList(l: Ptr IntList) -> i64 {
  case Ptr::get l {
    nil u -> 0 [i64]
    cons pair -> + (pair.0, sumList (pair.1.get ()))
}
}

/* tests/testdata/in/recursive_types.in:36-39 */ pub type Tree :: ∗ -> ∗ = fun (a) (variant{leaf 'a, branch (SharedPtr (Tree 'a), SharedPtr (Tree 'a))})

in "tests/testdata/in/recursive_types.in" in lines 41-43 pub fn mkLeaf ['a](x: 'a) -> Tree 'a {
  variant{Tree 'a leaf = x}
}

/* tests/testdata/in/recursive_types.in:46-52 */ pub type Expr = variant{lit i64, app struct{fun: SharedPtr Expr, arg: SharedPtr Expr}}

in "tests/testdata/in/recursive_types.in" in lines 54-56 pub fn mkApp(f: SharedPtr Expr, a: SharedPtr Expr) -> Expr {
  variant{Expr app = struct{fun = f, arg = a}}
}
