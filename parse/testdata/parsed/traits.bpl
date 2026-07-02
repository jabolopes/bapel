module /* testdata/in/traits.in:1 */test.traits

/* testdata/in/traits.in:3-5 */pub trait Size {
  fn size(s: Ptr Self) -> i64
}
/* testdata/in/traits.in:7-11 */impl Size for String {
  fn size(s: Ptr Self) -> i64 {
  String::size s
}
}
/* testdata/in/traits.in:13-15 */pub trait Indexable ['elem] {
  fn get(v: Ptr Self, index: i64) -> 'elem
}
/* testdata/in/traits.in:17-21 */impl ['a] Indexable 'a for Vector 'a {
  fn get(v: Ptr Self, index: i64) -> 'a {
  Vector::get (v, index)
}
}
