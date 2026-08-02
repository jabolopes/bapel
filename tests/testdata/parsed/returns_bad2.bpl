module /* tests/testdata/in/returns_bad2.in:2 */returns_bad2

in "tests/testdata/in/returns_bad2.in" in lines 4-16 fn return2() -> () {
  let f: i8 -> i8 = \(x: i8) -> {
  let y: i32 = 0
  ifthen (> (x, 0), {
  return y
})
  0
}
  ()
}
