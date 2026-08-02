module /* tests/testdata/in/context3.in:1 */context3

in "tests/testdata/in/context3.in" in lines 3-5 fn x() -> () {
  ()
}

in "tests/testdata/in/context3.in" in lines 8-10 fn y(x: i8) -> i8 {
  x
}

in "tests/testdata/in/context3.in" in lines 13-16 fn z() -> i16 {
  let x: i16 = 0
  x
}
