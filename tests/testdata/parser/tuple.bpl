module /* tests/testdata/in/tuple.in:1 */tuple

in "tests/testdata/in/tuple.in" in lines 3-5 fn mkTuple1() -> (i8, i16) {
  (1, 2)
}

in "tests/testdata/in/tuple.in" in lines 7-9 fn mkTuple2() -> (i8, i16) {
  (1 [i8], 2 [i16])
}

in "tests/testdata/in/tuple.in" in lines 11-23 fn mkTuple3() -> (i8, i16) {
  let r1: (i8, i16) = (1, 2)
  let r2: (i8, i16) = (1 [i8], 2 [i16])
  let x: i8 = r1.0
  let y: i16 = r1.1
  set r1 {0 = 3, 1 = 4}
  r1 <- set r1 {0 = 3, 1 = 4}
  let r: (i8, i16) = r1
  r1
}

/* tests/testdata/in/tuple.in:25 */ type Point = (i8, i16)

in "tests/testdata/in/tuple.in" in lines 27-29 fn mkPoint1() -> Point {
  (1, 2)
}

in "tests/testdata/in/tuple.in" in lines 31-33 fn mkPoint2() -> Point {
  (1 [i8], 2 [i16])
}

in "tests/testdata/in/tuple.in" in lines 35-47 fn mkPoint3() -> Point {
  let r1: Point = (1, 2)
  let r2: Point = (1 [i8], 2 [i16])
  let x: i8 = r1.0
  let y: i16 = r1.1
  set r1 {0 = 3, 1 = 4}
  r1 <- set r1 {0 = 3, 1 = 4}
  let r: Point = r1
  r1
}
