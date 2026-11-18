module /* tests/testdata/in/bad_recursive_type_mutual.in:2 */bad_recursive_type_mutual

/* tests/testdata/in/bad_recursive_type_mutual.in:4-6 */ type A = struct{b: B}

/* tests/testdata/in/bad_recursive_type_mutual.in:8-10 */ type B = struct{a: A}
