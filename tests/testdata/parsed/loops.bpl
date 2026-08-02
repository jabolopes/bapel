module /* tests/testdata/in/loops.in:2 */loops

in "tests/testdata/in/loops.in" in lines 4-6
imports {
  /* tests/testdata/in/loops.in:5 */bapel.core
}

in "tests/testdata/in/loops.in" in lines 8-15 fn testLoop() -> () {
  let i: i64 = 0
  for < (i, 10) {
  core::print [i8] 1
  i <- + (i, 1)
}
  ()
}
