module /* tests/testdata/in/traits_bounds_ok.in:1 */traits_bounds_ok

in "tests/testdata/in/traits_bounds_ok.in" in lines 3-5
impls {
  /* tests/testdata/in/traits_bounds_ok.in:4 */ "ptr.h"
}

/* tests/testdata/in/traits_bounds_ok.in:7 */ type Ptr :: ∗ -> ∗

in "tests/testdata/in/traits_bounds_ok.in" in lines 9-11 trait Size {
  in "tests/testdata/in/traits_bounds_ok.in" in line 10 fn size(s: Ptr Self) -> i64
}

/* tests/testdata/in/traits_bounds_ok.in:13 */ type S = struct{x: i64}

in "tests/testdata/in/traits_bounds_ok.in" in lines 15-19 impl Size for S {
  in "tests/testdata/in/traits_bounds_ok.in" in lines 16-18 fn size(s: Ptr Self) -> i64 {
  Ptr::get s.x
}
}

in "tests/testdata/in/traits_bounds_ok.in" in lines 21-23 fn printSize ['t: Size](x: Ptr 't) -> i64 {
  Size::size x
}

in "tests/testdata/in/traits_bounds_ok.in" in lines 25-28 fn run() -> i64 {
  let s: S = struct{x = 42}
  printSize [S] (Ptr::mk s)
}
