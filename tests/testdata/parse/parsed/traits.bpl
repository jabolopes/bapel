module /* tests/testdata/parse/in/traits.in:1 */test.traits

in "tests/testdata/parse/in/traits.in" in lines 3-5 pub trait Size {
  in "tests/testdata/parse/in/traits.in" in line 4 fn size(s: Ptr Self) -> i64
}

in "tests/testdata/parse/in/traits.in" in lines 7-11 impl Size for String {
  in "tests/testdata/parse/in/traits.in" in lines 8-10 fn size(s: Ptr Self) -> i64 {
  String::size s
}
}

in "tests/testdata/parse/in/traits.in" in lines 13-15 pub trait Indexable ['elem] {
  in "tests/testdata/parse/in/traits.in" in line 14 fn get(v: Ptr Self, index: i64) -> elem
}

in "tests/testdata/parse/in/traits.in" in lines 17-21 impl ['a] Indexable a for Vector a {
  in "tests/testdata/parse/in/traits.in" in lines 18-20 fn get(v: Ptr Self, index: i64) -> a {
  Vector::get (v, index)
}
}

