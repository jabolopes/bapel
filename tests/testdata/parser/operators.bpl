module /* tests/testdata/in/operators.in:1 */operators

in "tests/testdata/in/operators.in" in lines 3-42 fn operators() -> bool {
  let x: i8 = 0
  let y: bool = true
  || (y, y)
  && (y, y)
  != (x, x)
  == (x, x)
  > (x, x)
  >= (x, x)
  < (x, x)
  <= (x, x)
  + (x, 1)
  - (x, 1)
  * (x, 1)
  / (x, 1)
  ! y
  - (0, x)
  != [i8] (1, 2)
  == [i8] (1, 2)
  > [i8] (1, 2)
  >= [i8] (1, 2)
  < [i8] (1, 2)
  <= [i8] (1, 2)
  + [i8] (1, 2)
  - [i8] (1, 2)
  * [i8] (1, 2)
  / [i8] (1, 2)
  - [i8] (0, 1)
  || (true, false)
  && (true, false)
  ! (&& (y, y))
  ! (|| (y, y))
  ! (== (x, 1))
  - (0, + (x, 1))
  ! (! y)
  ! true
}
