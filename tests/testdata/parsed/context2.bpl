module /* tests/testdata/in/context2.in:2 */context2

in "tests/testdata/in/context2.in" in lines 4-6
imports {
  /* tests/testdata/in/context2.in:5 */bapel.core
}

in "tests/testdata/in/context2.in" in lines 8-10
impls {
  /* tests/testdata/in/context2.in:9 */ "context2_impl.bpl"
}

in "tests/testdata/in/context2.in" in lines 13-15 fn context2_impl(x: i8) -> i8 {
  x
}
