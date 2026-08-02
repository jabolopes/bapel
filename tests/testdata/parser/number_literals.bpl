module /* tests/testdata/in/number_literals.in:1 */number_literals

in "tests/testdata/in/number_literals.in" in lines 3-8 fn callWithIDs() -> () {
  let i: i8 = 0
  let j: i8 = 0
  + (i, j)
  ()
}

in "tests/testdata/in/number_literals.in" in lines 10-14 fn callWithIDAndLiterals() -> () {
  let i: i8 = 0
  + (i, 1)
  ()
}

in "tests/testdata/in/number_literals.in" in lines 16-20 fn callWithLiterals() -> () {
  let i: i8 = 0
  + (i, 1)
  ()
}

in "tests/testdata/in/number_literals.in" in lines 22-27 fn letWithIDs() -> () {
  let i: i8 = 0
  let j: i8 = 0
  let x: i8 = + (i, j)
  ()
}

in "tests/testdata/in/number_literals.in" in lines 29-33 fn letWithIDAndLiterals() -> () {
  let i: i8 = 0
  let x: i8 = + (i, 1)
  ()
}

in "tests/testdata/in/number_literals.in" in lines 35-38 fn letWithLiterals() -> () {
  let x: i8 = + (1, 2)
  ()
}

in "tests/testdata/in/number_literals.in" in lines 40-46 fn assignWithIDs() -> () {
  let x: i8 = 0
  let i: i8 = 0
  let j: i8 = 0
  x <- + (i, j)
  ()
}

in "tests/testdata/in/number_literals.in" in lines 48-53 fn assignWithIDAndLiterals() -> () {
  let x: i8 = 0
  let i: i8 = 0
  x <- + (i, 1)
  ()
}

in "tests/testdata/in/number_literals.in" in lines 55-59 fn assignWithLiterals() -> () {
  let x: i8 = 0
  x <- + (1, 2)
  ()
}

in "tests/testdata/in/number_literals.in" in lines 61-65 fn returnWithIDs() -> i8 {
  let i: i8 = 0
  let j: i8 = 0
  + (i, j)
}

in "tests/testdata/in/number_literals.in" in lines 67-70 fn returnWithIDAndLiterals() -> i8 {
  let i: i8 = 0
  + (i, 1)
}

in "tests/testdata/in/number_literals.in" in lines 72-74 fn returnWithLiterals() -> i8 {
  + (1, 2)
}
