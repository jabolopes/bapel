module /* tests/testdata/in/array.in:2 */array

in "tests/testdata/in/array.in" in lines 4-6
imports {
  /* tests/testdata/in/array.in:5 */bapel.core
}

in "tests/testdata/in/array.in" in lines 8-10 fn mkArray1() -> [i8, 10] {
  arr::mk [i8] ()
}

in "tests/testdata/in/array.in" in lines 12-24 fn mkArray2() -> [i8, 10] {
  let a: [i8, 10] = arr::mk [i8] ()
  let v1: i8 = arr::get [i8] (a, 0 [i64])
  let v2: i8 = arr::get [i8] (a, 0)
  let i: i64 = 0
  arr::set [i8] (a, i, 10 [i8])
  arr::set [i8] (a, i, 10)
  let r: [i8, 10] = a
  r
}
