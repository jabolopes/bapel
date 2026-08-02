module /* tests/testdata/in/traits_method_syntax.in:1 */traits_method_syntax

in "tests/testdata/in/traits_method_syntax.in" in lines 3-5
impls {
  /* tests/testdata/in/traits_method_syntax.in:4 */ "ptr.h"
}

/* tests/testdata/in/traits_method_syntax.in:7 */ type Ptr :: ∗ -> ∗

in "tests/testdata/in/traits_method_syntax.in" in lines 9-11 trait Size {
  in "tests/testdata/in/traits_method_syntax.in" in line 10 fn size(s: Ptr Self) -> i64
}

in "tests/testdata/in/traits_method_syntax.in" in lines 13-15 trait Add {
  in "tests/testdata/in/traits_method_syntax.in" in line 14 fn add(a: Ptr Self, b: Ptr Self) -> i64
}

/* tests/testdata/in/traits_method_syntax.in:17 */ type S = struct{x: i64}

in "tests/testdata/in/traits_method_syntax.in" in lines 19-23 impl Size for S {
  in "tests/testdata/in/traits_method_syntax.in" in lines 20-22 fn size(s: Ptr Self) -> i64 {
  Ptr::get s.x
}
}

in "tests/testdata/in/traits_method_syntax.in" in lines 25-29 impl Add for S {
  in "tests/testdata/in/traits_method_syntax.in" in lines 26-28 fn add(a: Ptr Self, b: Ptr Self) -> i64 {
  + (Ptr::get a.x, Ptr::get b.x)
}
}

in "tests/testdata/in/traits_method_syntax.in" in lines 31-33 fn make_s(val: i64) -> S {
  struct{x = val}
}

in "tests/testdata/in/traits_method_syntax.in" in lines 35-37 fn printSize ['t: Size](x: Ptr 't) -> i64 {
  x.size ()
}

in "tests/testdata/in/traits_method_syntax.in" in lines 39-62 fn run() -> i64 {
  let a: i64 = make_s 10.size ()
  let s1: S = make_s 20
  let b: i64 = s1.size
  let ref_s1: Ptr S = Ptr::mk s1
  let c: i64 = ref_s1.size ()
  let d: i64 = printSize [S] (Ptr::mk s1)
  let s2: S = make_s 30
  let e: i64 = s1.add (Ptr::mk s2)
  let f: i64 = ref_s1.size
  + (+ (+ (+ (+ (a, b), c), d), e), f)
}
