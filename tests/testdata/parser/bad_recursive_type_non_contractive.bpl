module /* tests/testdata/in/bad_recursive_type_non_contractive.in:2 */bad_recursive_type_non_contractive

/* tests/testdata/in/bad_recursive_type_non_contractive.in:4 */ type Loop = Loop

in "tests/testdata/in/bad_recursive_type_non_contractive.in" in lines 6-8 fn f(x: Loop) -> Loop {
  x
}
