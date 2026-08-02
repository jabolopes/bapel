module /* tests/testdata/in/structs.in:1 */structs

in "tests/testdata/in/structs.in" in lines 3-5 fn mkStruct1() -> struct{x: i8, y: i16} {
  struct{x = 1, y = 2}
}

in "tests/testdata/in/structs.in" in lines 7-9 fn mkStruct2() -> struct{x: i8, y: i16} {
  struct{x = 1 [i8], y = 2 [i16]}
}

in "tests/testdata/in/structs.in" in lines 11-26 fn mkStruct3() -> struct{x: i8, y: i16} {
  let r1: struct{x: i8, y: i16} = struct{x = 0, y = 0}
  let r2: struct{x: i8, y: i16} = struct{x = 0 [i8], y = 0 [i16]}
  let x: i8 = r1.x
  let y: i16 = r1.y
  set r1 {x = 3, y = 4}
  r1 <- set r1 {x = 3, y = 4}
  set r1 {0 = 3, 1 = 4}
  r1 <- set r1 {0 = 3, 1 = 4}
  let r: struct{x: i8, y: i16} = r1
  r
}

in "tests/testdata/in/structs.in" in lines 28-30 fn getStruct1(r: struct{x: i8}) -> i8 {
  r.x
}

in "tests/testdata/in/structs.in" in lines 32-34 fn getStruct2(r: struct{x: i8, y: i16}) -> i16 {
  r.y
}

/* tests/testdata/in/structs.in:36 */ type Point = struct{x: i8, y: i16}

in "tests/testdata/in/structs.in" in lines 38-40 fn mkPoint1() -> Point {
  struct{x = 1, y = 2}
}

in "tests/testdata/in/structs.in" in lines 42-44 fn mkPoint2() -> Point {
  struct{x = 1 [i8], y = 2 [i16]}
}

in "tests/testdata/in/structs.in" in lines 46-61 fn mkPoint3() -> Point {
  let r1: Point = struct{x = 0, y = 0}
  let r2: Point = struct{x = 0 [i8], y = 0 [i16]}
  let x: i8 = r1.x
  let y: i16 = r1.y
  set r1 {x = 3, y = 4}
  r1 <- set r1 {x = 3, y = 4}
  set r1 {0 = 3, 1 = 4}
  r1 <- set r1 {0 = 3, 1 = 4}
  let r: Point = r1
  r
}

in "tests/testdata/in/structs.in" in lines 63-65 fn getPoint1(p: Point) -> i8 {
  p.x
}

in "tests/testdata/in/structs.in" in lines 67-69 fn getPoint2(p: Point) -> i16 {
  p.y
}
