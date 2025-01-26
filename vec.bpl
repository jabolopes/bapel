exports {
  type vec.Vector ['a]

  vec.mk: () -> vec.Vector 'a
  vec.add: (vec.Vector 'a, 'a) -> ()
}

impls {
  vec_impl.cc
}
