module /* tests/testdata/in/bad_recursive_type_partial_app.in:2 */bad_recursive_type_partial_app

in "tests/testdata/in/bad_recursive_type_partial_app.in" in lines 4-6
imports {
  /* tests/testdata/in/bad_recursive_type_partial_app.in:5 */bapel.stl
}

/* tests/testdata/in/bad_recursive_type_partial_app.in:8-12 */ type Partial :: ∗ -> ∗ -> ∗ = fun (a) (fun (b) (struct{x: 'a, y: 'b, next: SharedPtr (Partial 'a)}))
