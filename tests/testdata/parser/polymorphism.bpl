module /* tests/testdata/in/polymorphism.in:2 */polymorphism

in "tests/testdata/in/polymorphism.in" in lines 4-6
imports {
  /* tests/testdata/in/polymorphism.in:5 */bapel.core
}

in "tests/testdata/in/polymorphism.in" in lines 8-10 fn id ['a](x: 'a) -> 'a {
  x
}

in "tests/testdata/in/polymorphism.in" in lines 12-14 fn callPolymorphic() -> () {
  core::print [i8] 1
}

in "tests/testdata/in/polymorphism.in" in lines 16-19 fn functionSubtyping() -> () {
  let id2: i8 -> i8 = id [i8]
  ()
}
