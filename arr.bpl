exports {
  arr.mk : forall ['a] () -> ['a, 10]
  arr.get : forall ['a] (['a, 10], i64) -> 'a
  arr.set : forall ['a] (['a, 10], i64, 'a) -> ()
}

impls {
  arr_impl.cc
}
