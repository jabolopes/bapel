module /* tests/testdata/in/coherence_violation.in:2 */coherence_violation

in "tests/testdata/in/coherence_violation.in" in lines 4-6
imports {
  /* tests/testdata/in/coherence_violation.in:5 */bapel.stl
}

in "tests/testdata/in/coherence_violation.in" in lines 8-10 trait MyTrait {
  in "tests/testdata/in/coherence_violation.in" in line 9 fn foo(x: Ptr Self) -> ()
}

in "tests/testdata/in/coherence_violation.in" in lines 12-16 impl MyTrait for String {
  in "tests/testdata/in/coherence_violation.in" in lines 13-15 fn foo(x: Ptr Self) -> () {
  ()
}
}
