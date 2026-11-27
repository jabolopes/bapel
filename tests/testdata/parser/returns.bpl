module /* tests/testdata/in/returns.in:1 */returns

in "tests/testdata/in/returns.in" in lines 3-5 fn return1() -> i8 {
  0
}

in "tests/testdata/in/returns.in" in lines 7-12 fn return2() -> i8 {
  ifthen (true, {
  return 0
})
  1
}

in "tests/testdata/in/returns.in" in lines 14-21 fn return3() -> i8 {
  ifelse (true, {
  return 0
}, {
  return 1
})
  2
}

in "tests/testdata/in/returns.in" in lines 23-25 fn noop() -> () {
  ()
}

in "tests/testdata/in/returns.in" in lines 27-31 fn unit_if_no_else(cond: bool) -> () {
  ifthen (cond, {
  noop ()
})
}

in "tests/testdata/in/returns.in" in lines 33-40 fn unit_early_return_if(cond: bool) -> () {
  ifthen (cond, {
  return ()
})
  ifthen (cond, {
  noop ()
})
}

in "tests/testdata/in/returns.in" in lines 42-48 fn unit_for_loop(n: i64) -> () {
  let i: i64 = 0
  for < (i, n) {
  i <- + (i, 1)
  noop ()
}
}

in "tests/testdata/in/returns.in" in lines 50-56 fn value_if_else_implicit(cond: bool) -> i64 {
  ifelse (cond, {
  1
}, {
  2
})
}

in "tests/testdata/in/returns.in" in lines 58-67 fn value_for_early_return(n: i64) -> i64 {
  let i: i64 = 0
  for < (i, n) {
  ifthen (== (i, 5), {
  return i
})
  i <- + (i, 1)
}
  0
}

/* tests/testdata/in/returns.in:69-72 */ type Choice = variant{a i64, b i64}

in "tests/testdata/in/returns.in" in lines 74-79 fn match_exhaustive_value(t: Choice) -> i64 {
  case t {
    a x -> x
    b y -> + (y, 1)
}
}

in "tests/testdata/in/returns.in" in lines 81-86 fn match_wildcard_value(t: Choice) -> i64 {
  case t {
    a x -> x
    _ -> 0
}
}

in "tests/testdata/in/returns.in" in lines 88-97 fn unit_match_exhaustive(t: Choice) -> () {
  case t {
    a v -> {
  noop ()
}
    b v -> {
  noop ()
}
}
}
