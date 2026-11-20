module /* tests/testdata/in/anonymous_structs.in:1 */anonymous_structs

/* tests/testdata/in/anonymous_structs.in:3-9 */ pub type Expr = variant{lit i64, add struct{lhs: i64, rhs: i64}}

/* tests/testdata/in/anonymous_structs.in:11-17 */ pub type Outer = struct{label: i64, inner: struct{a: i64, b: i64}}

in "tests/testdata/in/anonymous_structs.in" in lines 19-24 pub fn eval(e: Expr) -> i64 {
  case e {
    lit x -> x
    add p -> + (p.lhs, p.rhs)
}
}

in "tests/testdata/in/anonymous_structs.in" in lines 26-28 pub fn sum_point(p: struct{x: i64, y: i64}) -> i64 {
  + (p.x, p.y)
}

in "tests/testdata/in/anonymous_structs.in" in lines 30-32 pub fn make_expr() -> Expr {
  variant{Expr add = struct{lhs = 1 [i64], rhs = 2 [i64]}}
}
