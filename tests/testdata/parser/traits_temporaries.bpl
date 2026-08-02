module /* tests/testdata/in/traits_temporaries.in:1 */traits_temporaries

in "tests/testdata/in/traits_temporaries.in" in lines 3-5
impls {
  /* tests/testdata/in/traits_temporaries.in:4 */ "ptr.h"
}

/* tests/testdata/in/traits_temporaries.in:7 */ type Ptr :: ∗ -> ∗

in "tests/testdata/in/traits_temporaries.in" in lines 9-11 trait Size {
  in "tests/testdata/in/traits_temporaries.in" in line 10 fn size(s: Ptr Self) -> i64
}

/* tests/testdata/in/traits_temporaries.in:13 */ type S = struct{x: i64}

in "tests/testdata/in/traits_temporaries.in" in lines 15-19 impl Size for S {
  in "tests/testdata/in/traits_temporaries.in" in lines 16-18 fn size(s: Ptr Self) -> i64 {
  Ptr::get s.x
}
}

in "tests/testdata/in/traits_temporaries.in" in lines 21-23 fn make_s(val: i64) -> S {
  struct{x = val}
}

in "tests/testdata/in/traits_temporaries.in" in lines 25-27 fn printSize ['t: Size](x: Ptr 't) -> i64 {
  Size::size x
}

in "tests/testdata/in/traits_temporaries.in" in lines 29-33 fn run() -> i64 {
  let a: i64 = printSize [S] (Ptr::mk (make_s 10))
  let b: i64 = Size::size (Ptr::mk (make_s 20))
  + (a, b)
}
