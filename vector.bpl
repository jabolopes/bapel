exports {
  type vector.Vector ['a]

  vector.mk: () -> vector.Vector 'a
  vector.add: (vector.Vector 'a, 'a) -> ()
}
