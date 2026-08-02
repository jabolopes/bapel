module /* tests/testdata/in/traits_bounds_error.in:2 */traits_bounds_error

in "tests/testdata/in/traits_bounds_error.in" in lines 4-6
impls {
  /* tests/testdata/in/traits_bounds_error.in:5 */ "ptr.h"
}

/* tests/testdata/in/traits_bounds_error.in:8 */ type Ptr :: ∗ -> ∗

in "tests/testdata/in/traits_bounds_error.in" in lines 10-12 trait Size {
  in "tests/testdata/in/traits_bounds_error.in" in line 11 fn size(s: Ptr Self) -> i64
}

/* tests/testdata/in/traits_bounds_error.in:14 */ type S = struct{x: i64}

in "tests/testdata/in/traits_bounds_error.in" in lines 16-20 impl Size for S {
  in "tests/testdata/in/traits_bounds_error.in" in lines 17-19 fn size(s: Ptr Self) -> i64 {
  Ptr::get s.x
}
}

/* tests/testdata/in/traits_bounds_error.in:22 */ type Unrelated = struct{y: i64}

in "tests/testdata/in/traits_bounds_error.in" in lines 24-26 fn printSize ['t: Size](x: Ptr 't) -> i64 {
  Size::size x
}

in "tests/testdata/in/traits_bounds_error.in" in lines 28-31 fn run() -> i64 {
  let u: Unrelated = struct{y = 42}
  printSize [Unrelated] (Ptr::mk u)
}
