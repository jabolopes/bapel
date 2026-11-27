module /* tests/testdata/in/returns_bad7.in:2 */returns_bad7

/* tests/testdata/in/returns_bad7.in:4-7 */ type Choice = variant{a i64, b bool}

in "tests/testdata/in/returns_bad7.in" in lines 9-14 fn return_match_mismatch(c: Choice) -> i64 {
  case c {
    a x -> x
    b y -> y
}
}
