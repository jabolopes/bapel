module /* tests/testdata/in/bad_recursive_type_non_regular.in:2 */bad_recursive_type_non_regular

in "tests/testdata/in/bad_recursive_type_non_regular.in" in lines 4-6
imports {
  /* tests/testdata/in/bad_recursive_type_non_regular.in:5 */bapel.stl
}

/* tests/testdata/in/bad_recursive_type_non_regular.in:8-11 */ type Nested :: ∗ -> ∗ = fun (a) (variant{nil (), cons ('a, SharedPtr (Nested ('a, 'a)))})
