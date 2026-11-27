module /* tests/testdata/in/returns_bad3.in:2 */returns_bad3

in "tests/testdata/in/returns_bad3.in" in lines 4-8 fn return_if_no_else(cond: bool) -> i64 {
  ifthen (cond, {
  return 1
})
}
