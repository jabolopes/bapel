module /* tests/testdata/in/bad_traits_anonymous_struct.in:2 */bad_traits_anonymous_struct

in "tests/testdata/in/bad_traits_anonymous_struct.in" in lines 4-6
impls {
  /* tests/testdata/in/bad_traits_anonymous_struct.in:5 */ "ptr.h"
}

/* tests/testdata/in/bad_traits_anonymous_struct.in:8 */ type Ptr :: ∗ -> ∗

in "tests/testdata/in/bad_traits_anonymous_struct.in" in lines 10-12 trait Size {
  in "tests/testdata/in/bad_traits_anonymous_struct.in" in line 11 fn size(s: Ptr Self) -> i64
}

/* tests/testdata/in/bad_traits_anonymous_struct.in:14-17 */ type Point = struct{x: i64, y: i64}

in "tests/testdata/in/bad_traits_anonymous_struct.in" in lines 19-23 impl Size for Point {
  in "tests/testdata/in/bad_traits_anonymous_struct.in" in lines 20-22 fn size(p: Ptr Self) -> i64 {
  + (Ptr::get p.x, Ptr::get p.y)
}
}

in "tests/testdata/in/bad_traits_anonymous_struct.in" in lines 25-28 fn run() -> i64 {
  let p = struct{x = 10, y = 20}
  Ptr::mk p.size ()
}
