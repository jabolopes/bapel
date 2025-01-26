exports {
  type vec.Vector ['a]

  vec.mk: () -> vec.Vector 'a
  vec.add: forall ['a] (vec.Vector 'a, 'a) -> ()
  vec.get: forall ['a] (vec.Vector 'a, i64) -> 'a
  vec.set: forall ['a] (vec.Vector 'a, i64, 'a) -> ()
}

impls {
  vec_impl.cc
}
