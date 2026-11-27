module /* tests/testdata/in/returns_bad5.in:2 */returns_bad5

in "tests/testdata/in/returns_bad5.in" in lines 4-10 fn return_if_else_mismatch(cond: bool) -> i64 {
  ifelse (cond, {
  1
}, {
  true
})
}
