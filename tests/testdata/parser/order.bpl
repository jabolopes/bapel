module /* tests/testdata/in/order.in:3 */order

in "tests/testdata/in/order.in" in lines 5-7 fn mkC() -> C {
  struct{b = struct{a = mkA ()}}
}

/* tests/testdata/in/order.in:9 */ type C = struct{b: B}

/* tests/testdata/in/order.in:11 */ type B = struct{a: A}

/* tests/testdata/in/order.in:13 */ mkA: () -> A

/* tests/testdata/in/order.in:15 */ type A

/* tests/testdata/in/order.in:17 */ type a

/* tests/testdata/in/order.in:19 */ type b = struct{a: a}

/* tests/testdata/in/order.in:21 */ type c = struct{b: b}

/* tests/testdata/in/order.in:23 */ mka: () -> a
