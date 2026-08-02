module /* tests/testdata/in/parser_test.in:2 */parser

in "tests/testdata/in/parser_test.in" in lines 4-7
imports {
  /* tests/testdata/in/parser_test.in:5 */core
  /* tests/testdata/in/parser_test.in:6 */vec
}

in "tests/testdata/in/parser_test.in" in lines 9-12
impls {
  /* tests/testdata/in/parser_test.in:10 */ "f1.bpl"
  /* tests/testdata/in/parser_test.in:11 */ "f2.cc"
}

/* tests/testdata/in/parser_test.in:14 */ export f: i8

in "tests/testdata/in/parser_test.in" in lines 16-27 fn assign() -> () {
  r <- 1
  r <- x
  (r1, r2) <- (a1, a2)
  r <- f0 ()
  r <- f1 x
  r <- f2 (x, y)
  r <- x.1
  r <- x.y
  r <- - (0, a)
  r <- + (a, b)
}

in "tests/testdata/in/parser_test.in" in lines 29-40 fn expression() -> () {
  x
  f0 ()
  f1 x
  f2 (x, y)
  a.1
  a.x
  - (0, a)
  + (a, b)
  ! a
  1 [i8]
}

in "tests/testdata/in/parser_test.in" in lines 42-54 fn operators() -> () {
  != (a, b)
  == (a, b)
  > (a, b)
  >= (a, b)
  < (a, b)
  <= (a, b)
  + (a, b)
  - (a, b)
  * (a, b)
  / (a, b)
  ! a
}

in "tests/testdata/in/parser_test.in" in lines 56-59 fn projection() -> () {
  x.a
  x.a ()
}

in "tests/testdata/in/parser_test.in" in lines 61-87 fn term() -> () {
  x <- 1
  ifelse (x, {
  0
}, {
  1
})
  ifelse (! x, {
  0
}, {
  1
})
  ifelse (x [i8], {
  0
}, {
  1
})
  ()
  (x, x)
  x
}

in "tests/testdata/in/parser_test.in" in lines 89-91 fn fn1() -> () {
  ()
}

in "tests/testdata/in/parser_test.in" in lines 93-95 fn fn2(a: i32) -> () {
  ()
}

in "tests/testdata/in/parser_test.in" in lines 97-99 fn fn3() -> i64 {
  ()
}

in "tests/testdata/in/parser_test.in" in lines 101-103 fn fn4(a: [i32, 9223372036854775807], b: i64) -> () {
  ()
}

in "tests/testdata/in/parser_test.in" in lines 105-107 fn fn5(a: [i32, 9223372036854775807], b: i64) -> (i32, [i64, 9223372036854775807]) {
  ()
}

in "tests/testdata/in/parser_test.in" in lines 109-111 fn fn6 ['a](x: 'a) -> 'a {
  ()
}

in "tests/testdata/in/parser_test.in" in lines 113-115 fn fn7 ['a, 'b](x: 'a, y: 'b) -> ('a, 'b) {
  ()
}

in "tests/testdata/in/parser_test.in" in lines 117-182 fn matchTerms() -> () {
  case v { none l -> l }
  case v {
    none l -> l
    some r -> v1
}
  case v { none l -> l }
  case v {
    none l -> l
    some r -> v1
}
  case v {
    none l -> {
  l
}
    some r -> v1
}
  case v {
    none l -> {
  l
}
    some r -> v1
}
  case v {
    none l -> {
  l
}
    some r -> {
  v1
}
}
  let v2: V = case v { none l -> l }
  let v2: V = case v {
    none l -> l
    some r -> v1
}
  let v2: V = case v { none l -> l }
  let v2: V = case v {
    none l -> l
    some r -> v1
}
  let v2: V = case v {
    none l -> {
  l
}
    some r -> v1
}
  let v2: V = case v {
    none l -> {
  l
}
    some r -> v1
}
  let v2: V = case v {
    none l -> {
  l
}
    some r -> {
  v1
}
}
  let v2: V = case variant{V left = l} {
    left l -> {
  l
}
    right r -> 0
}
  ()
}

in "tests/testdata/in/parser_test.in" in lines 184-221 fn setTerms() -> () {
  set p {x = 0}
  set p {x = 0, y = 1}
  set p {x = 0}
  set p {x = 0, y = 1}
  let p2: P = set p {x = 0}
  let p2: P = set p {x = 0, y = 1}
  let p2: P = set p {x = 0}
  let p2: P = set p {x = 0, y = 1}
  p2 <- set p {x = 0}
  p2 <- set p {x = 0, y = 1}
  p2 <- set p {x = 0}
  p2 <- set p {x = 0, y = 1}
  let p2: P = set struct{a = 0, b = 1} {x = 0, y = 1}
  ()
}

in "tests/testdata/in/parser_test.in" in lines 223-258 fn structTerms() -> () {
  struct{}
  struct{a = 0}
  struct{a = 0, b = 1}
  struct{a = 0}
  struct{a = 0, b = 1}
  let s: S = struct{}
  let s: S = struct{a = 0}
  let s: S = struct{a = 0, b = 1}
  let s: S = struct{a = 0}
  let s: S = struct{a = 0, b = 1}
  s <- struct{}
  s <- struct{a = 0}
  s <- struct{a = 0, b = 1}
  s <- struct{a = 0}
  s <- struct{a = 0, b = 1}
  ()
}

in "tests/testdata/in/parser_test.in" in lines 260-277 fn variantTerms() -> () {
  variant{V left = l}
  variant{V left = l}
  let v: V = variant{V left = l}
  let v: V = variant{V left = l}
  v <- variant{V left = l}
  v <- variant{P left = l}
  ()
}

in "tests/testdata/in/parser_test.in" in lines 280-283 fn comments() -> () {
  ()
}

in "tests/testdata/in/parser_test.in" in lines 285-291 fn conditionals() -> () {
  ifthen (a, {
  b
})
  ifelse (a, {
  b
}, {
  c
})
  ifelse (a, {
  b
}, ifthen (c, {
  d
}))
  ifelse (a, {
  b
}, ifelse (c, {
  d
}, {
  e
}))
  ()
}

in "tests/testdata/in/parser_test.in" in lines 293-300 fn blocks() -> () {
  {
  ()
}
  {
  {
  ()
}
  {
  ()
}
}
  ()
}

in "tests/testdata/in/parser_test.in" in lines 302-307 fn loops() -> () {
  for < (x, 10) {
  ()
}
  ()
}

/* tests/testdata/in/parser_test.in:309 */ x: i8

/* tests/testdata/in/parser_test.in:310 */ x: i16

/* tests/testdata/in/parser_test.in:311 */ x: struct{}

/* tests/testdata/in/parser_test.in:312 */ x: struct{a: i8}

/* tests/testdata/in/parser_test.in:313 */ x: struct{a: i8, b: i16}

/* tests/testdata/in/parser_test.in:314-316 */ x: struct{a: i8}

/* tests/testdata/in/parser_test.in:317-320 */ x: struct{a: i8, b: i16}

/* tests/testdata/in/parser_test.in:321 */ x: variant{left i8}

/* tests/testdata/in/parser_test.in:322 */ x: variant{left i8, right i16}

/* tests/testdata/in/parser_test.in:323-325 */ x: variant{left i8}

/* tests/testdata/in/parser_test.in:326-329 */ x: variant{left i8, right i16}

/* tests/testdata/in/parser_test.in:330 */ x: ()

/* tests/testdata/in/parser_test.in:331 */ x: (i8, i16)

/* tests/testdata/in/parser_test.in:332 */ x: [i8, 10]

/* tests/testdata/in/parser_test.in:333 */ x: () -> ()

/* tests/testdata/in/parser_test.in:334 */ x: i8 -> i16

/* tests/testdata/in/parser_test.in:335 */ x: i8 -> (i8, i16)

/* tests/testdata/in/parser_test.in:336 */ x: (i8, i16) -> i8

/* tests/testdata/in/parser_test.in:337 */ x: (i8, i16) -> (i8, i16)

/* tests/testdata/in/parser_test.in:338 */ x: forall ['a] 'a -> 'a

/* tests/testdata/in/parser_test.in:340 */ type T = struct{}

/* tests/testdata/in/parser_test.in:341 */ type T = struct{a: i8}

/* tests/testdata/in/parser_test.in:342 */ type T = struct{a: i8, b: i16}

/* tests/testdata/in/parser_test.in:343-345 */ type T = struct{a: i8}

/* tests/testdata/in/parser_test.in:346-349 */ type T = struct{a: i8, b: i16}

/* tests/testdata/in/parser_test.in:350 */ type T = variant{left i8}

/* tests/testdata/in/parser_test.in:351 */ type T = variant{left i8, right i16}

/* tests/testdata/in/parser_test.in:352-354 */ type T = variant{left i8}

/* tests/testdata/in/parser_test.in:355-358 */ type T = variant{left i8, right i16}
