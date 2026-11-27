module /* tests/testdata/in/returns_bad4.in:2 */returns_bad4

in "tests/testdata/in/returns_bad4.in" in lines 4-9 fn return_for_only(n: i64) -> i64 {
  let i: i64 = 0
  for < (i, n) {
  return i
}
}
