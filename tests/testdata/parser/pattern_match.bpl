module /* tests/testdata/in/pattern_match.in:1 */pattern_match

/* tests/testdata/in/pattern_match.in:3-7 */ export type Color = variant{red (), green (), blue ()}

in "tests/testdata/in/pattern_match.in" in lines 10-15 pub fn is_red(c: Color) -> bool {
  case c {
    red u -> true
    _ -> false
}
}

in "tests/testdata/in/pattern_match.in" in lines 18-25 pub fn color_code(c: Color) -> i64 {
  case c {
    red u -> 1
    green u -> 2
    blue u -> 3
    _ -> 0
}
}

in "tests/testdata/in/pattern_match.in" in lines 28-32 pub fn always_zero(c: Color) -> i64 {
  case c { _ -> 0 }
}
