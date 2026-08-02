module /* tests/testdata/in/parameterized_traits.in:1 */parameterized_traits

in "tests/testdata/in/parameterized_traits.in" in lines 3-6
impls {
  /* tests/testdata/in/parameterized_traits.in:4 */ "ptr.h"
  /* tests/testdata/in/parameterized_traits.in:5 */ "vector.h"
}

/* tests/testdata/in/parameterized_traits.in:8 */ type Ptr :: ∗ -> ∗

/* tests/testdata/in/parameterized_traits.in:10 */ type Vector :: ∗ -> ∗

in "tests/testdata/in/parameterized_traits.in" in lines 13-15 pub trait Indexable ['elem] {
  in "tests/testdata/in/parameterized_traits.in" in line 14 fn get(v: Ptr Self, index: i64) -> 'elem
}

in "tests/testdata/in/parameterized_traits.in" in lines 17-21 impl ['a] Indexable 'a for Vector 'a {
  in "tests/testdata/in/parameterized_traits.in" in lines 18-20 fn get(v: Ptr Self, index: i64) -> 'a {
  vector_get (v, index)
}
}

in "tests/testdata/in/parameterized_traits.in" in lines 23-25 fn run(v: Ptr (Vector i8)) -> i8 {
  Indexable::get (v, 0)
}
