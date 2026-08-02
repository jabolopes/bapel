module /* tests/testdata/in/returns.in:1 */returns

in "tests/testdata/in/returns.in" in lines 3-5 fn return1() -> i8 {
  0
}

in "tests/testdata/in/returns.in" in lines 7-12 fn return2() -> i8 {
  ifthen (true, {
  return 0
})
  1
}

in "tests/testdata/in/returns.in" in lines 14-21 fn return3() -> i8 {
  ifelse (true, {
  return 0
}, {
  return 1
})
  2
}
