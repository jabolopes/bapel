module /* tests/testdata/in/lambda.in:1 */lambda

in "tests/testdata/in/lambda.in" in lines 3-6 fn lambda() -> i8 {
  let add: i8 -> i8 = \(x: i8) -> {
  + (x, 1 [i8])
}
  add 2
}

in "tests/testdata/in/lambda.in" in lines 8-11 pub fn lambda2() -> i8 {
  let add: i8 -> i8 = \(x: i8) -> {
  + (x, 1 [i8])
}
  add 2
}

in "tests/testdata/in/lambda.in" in lines 13-29 pub fn lambda3() -> i8 {
  let add: i32 -> i32 = \(x: i32) -> {
  ifthen (> (x, 0), {
  return 0 [i32]
})
  1 [i32]
}
  let add2 = \(x: i32) -> {
  ifthen (> (x, 0), {
  return 0 [i32]
})
  1 [i32]
}
  0
}
