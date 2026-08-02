module /* tests/testdata/in/context1.in:2 */context1

in "tests/testdata/in/context1.in" in lines 4-6
imports {
  /* tests/testdata/in/context1.in:5 */bapel.core
}

in "tests/testdata/in/context1.in" in lines 9-11 fn core::print() -> () {
  ()
}

in "tests/testdata/in/context1.in" in lines 13-15 fn callPrint() -> () {
  core::print ()
}
