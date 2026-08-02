module /* tests/testdata/in/returns_bad1.in:2 */returns_bad1

in "tests/testdata/in/returns_bad1.in" in lines 4-12 fn return1() -> i8 {
  let x: i32 = 0
  ifthen (> (x, 0), {
  return x
})
  0
}
