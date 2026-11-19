module /* tests/testdata/in/bad_pattern_match_duplicate.in:2 */bad_pattern_match_duplicate

/* tests/testdata/in/bad_pattern_match_duplicate.in:4-7 */ type Color = variant{red (), green ()}

in "tests/testdata/in/bad_pattern_match_duplicate.in" in lines 9-15 fn bad(c: Color) -> i64 {
  case c {
    red u -> 1
    red v -> 2
    green w -> 3
}
}
