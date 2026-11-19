module /* tests/testdata/in/bad_pattern_match_unreachable.in:2 */bad_pattern_match_unreachable

/* tests/testdata/in/bad_pattern_match_unreachable.in:4-7 */ type Color = variant{red (), green ()}

in "tests/testdata/in/bad_pattern_match_unreachable.in" in lines 9-14 fn bad(c: Color) -> i64 {
  case c {
    _ -> 0
    red u -> 1
}
}
