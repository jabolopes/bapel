module /* testdata/in/traits.in:1 */test.traits

/* testdata/in/traits.in:3-5 */pub trait Size {
  fn size(s: Ptr Self) -> i64
}
/* testdata/in/traits.in:7-11 */impl Size for String {
  fn size(s: Ptr Self) -> i64 {
  String::size s
}
}
