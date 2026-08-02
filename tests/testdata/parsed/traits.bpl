module /* tests/testdata/in/traits.in:1 */traits

in "tests/testdata/in/traits.in" in lines 3-6
impls {
  /* tests/testdata/in/traits.in:4 */ "ptr.h"
  /* tests/testdata/in/traits.in:5 */ "vector.h"
}

/* tests/testdata/in/traits.in:8 */ type Ptr :: ∗ -> ∗

/* tests/testdata/in/traits.in:10 */ type Vector :: ∗ -> ∗

/* tests/testdata/in/traits.in:12 */ type MyStruct = struct{x: i64}

in "tests/testdata/in/traits.in" in lines 14-16 pub trait Size {
  in "tests/testdata/in/traits.in" in line 15 fn size(s: Self) -> i64
}

in "tests/testdata/in/traits.in" in lines 18-22 impl Size for MyStruct {
  in "tests/testdata/in/traits.in" in lines 19-21 fn size(s: Self) -> i64 {
  s.x
}
}

in "tests/testdata/in/traits.in" in lines 24-26 pub trait Indexable ['elem] {
  in "tests/testdata/in/traits.in" in line 25 fn get(v: Ptr Self, index: i64) -> 'elem
}

in "tests/testdata/in/traits.in" in lines 28-32 impl ['a] Indexable 'a for Vector 'a {
  in "tests/testdata/in/traits.in" in lines 29-31 fn get(v: Ptr Self, index: i64) -> 'a {
  vector_get (v, index)
}
}

in "tests/testdata/in/traits.in" in lines 34-36 fn run(s: MyStruct, v: Ptr (Vector i8)) -> (i64, i8) {
  (Size::size s, Indexable::get (v, 0))
}
