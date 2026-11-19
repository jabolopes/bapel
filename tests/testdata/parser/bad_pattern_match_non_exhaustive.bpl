module /* tests/testdata/in/bad_pattern_match_non_exhaustive.in:2 */bad_pattern_match_non_exhaustive

/* tests/testdata/in/bad_pattern_match_non_exhaustive.in:4-8 */ type Color = variant{red (), green (), blue ()}

in "tests/testdata/in/bad_pattern_match_non_exhaustive.in" in lines 10-15 fn bad(c: Color) -> i64 {
  case c {
    red u -> 1
    green v -> 2
}
}
