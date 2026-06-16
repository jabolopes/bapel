module /* testdata/in/parser_test.in:1 */parser

/* testdata/in/parser_test.in:3-6 */imports {
  /* testdata/in/parser_test.in:4 */core
  /* testdata/in/parser_test.in:5 */vec
}

/* testdata/in/parser_test.in:8-11 */impls {
  /* testdata/in/parser_test.in:9 */"f1.bpl"
  /* testdata/in/parser_test.in:10 */"f2.cc"
}

/* testdata/in/parser_test.in:13 */export f: i8
/* testdata/in/parser_test.in:15-26 */fn assign() -> () {
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
/* testdata/in/parser_test.in:28-39 */fn expression() -> () {
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
/* testdata/in/parser_test.in:41-53 */fn operators() -> () {
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
/* testdata/in/parser_test.in:55-58 */fn projection() -> () {
  x.a
  x.a ()
}
/* testdata/in/parser_test.in:60-86 */fn term() -> () {
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
/* testdata/in/parser_test.in:88-90 */fn fn1() -> () {
  ()
}
/* testdata/in/parser_test.in:92-94 */fn fn2(a: i32) -> () {
  ()
}
/* testdata/in/parser_test.in:96-98 */fn fn3() -> i64 {
  ()
}
/* testdata/in/parser_test.in:100-102 */fn fn4(a: [i32], b: i64) -> () {
  ()
}
/* testdata/in/parser_test.in:104-106 */fn fn5(a: [i32], b: i64) -> (i32, [i64]) {
  ()
}
/* testdata/in/parser_test.in:108-110 */fn fn6['a ∗](x: 'a) -> 'a {
  ()
}
/* testdata/in/parser_test.in:112-114 */fn fn7['a ∗, 'b ∗](x: 'a, y: 'b) -> ('a, 'b) {
  ()
}
/* testdata/in/parser_test.in:116-181 */fn matchTerms() -> () {
  case v { none l -> l }
  case v {
    none l -> l
    some r -> v1}
  case v { none l -> l }
  case v {
    none l -> l
    some r -> v1}
  case v {
    none l -> {
  l
}
    some r -> v1}
  case v {
    none l -> {
  l
}
    some r -> v1}
  case v {
    none l -> {
  l
}
    some r -> {
  v1
}}
  let v2: V = case v { none l -> l }
  let v2: V = case v {
    none l -> l
    some r -> v1}
  let v2: V = case v { none l -> l }
  let v2: V = case v {
    none l -> l
    some r -> v1}
  let v2: V = case v {
    none l -> {
  l
}
    some r -> v1}
  let v2: V = case v {
    none l -> {
  l
}
    some r -> v1}
  let v2: V = case v {
    none l -> {
  l
}
    some r -> {
  v1
}}
  let v2: V = case variant{V left = l} {
    left l -> {
  l
}
    right r -> 0}
  ()
}
/* testdata/in/parser_test.in:183-220 */fn setTerms() -> () {
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
/* testdata/in/parser_test.in:222-257 */fn structTerms() -> () {
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
/* testdata/in/parser_test.in:259-276 */fn variantTerms() -> () {
  variant{V left = l}
  variant{V left = l}
  let v: V = variant{V left = l}
  let v: V = variant{V left = l}
  v <- variant{V left = l}
  v <- variant{P left = l}
  ()
}
/* testdata/in/parser_test.in:279-282 */fn comments() -> () {
  ()
}
/* testdata/in/parser_test.in:284-290 */fn conditionals() -> () {
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
/* testdata/in/parser_test.in:292-299 */fn blocks() -> () {
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
/* testdata/in/parser_test.in:301-306 */fn loops() -> () {
  for < (x, 10) {
  ()
}
  ()
}
/* testdata/in/parser_test.in:308 */x: i8
/* testdata/in/parser_test.in:309 */x: i16
/* testdata/in/parser_test.in:310 */x: struct{}
/* testdata/in/parser_test.in:311 */x: struct{a i8}
/* testdata/in/parser_test.in:312 */x: struct{a i8, b i16}
/* testdata/in/parser_test.in:313-315 */x: struct{a i8}
/* testdata/in/parser_test.in:316-319 */x: struct{a i8, b i16}
/* testdata/in/parser_test.in:320 */x: variant{left i8}
/* testdata/in/parser_test.in:321 */x: variant{left i8, right i16}
/* testdata/in/parser_test.in:322-324 */x: variant{left i8}
/* testdata/in/parser_test.in:325-328 */x: variant{left i8, right i16}
/* testdata/in/parser_test.in:329 */x: ()
/* testdata/in/parser_test.in:330 */x: (i8, i16)
/* testdata/in/parser_test.in:331 */x: [i8, 10]
/* testdata/in/parser_test.in:332 */x: () -> ()
/* testdata/in/parser_test.in:333 */x: i8 -> i16
/* testdata/in/parser_test.in:334 */x: i8 -> (i8, i16)
/* testdata/in/parser_test.in:335 */x: (i8, i16) -> i8
/* testdata/in/parser_test.in:336 */x: (i8, i16) -> (i8, i16)
/* testdata/in/parser_test.in:337 */x: forall ['a] 'a -> 'a
/* testdata/in/parser_test.in:339 */type T = struct{}
/* testdata/in/parser_test.in:340 */type T = struct{a i8}
/* testdata/in/parser_test.in:341 */type T = struct{a i8, b i16}
/* testdata/in/parser_test.in:342-344 */type T = struct{a i8}
/* testdata/in/parser_test.in:345-348 */type T = struct{a i8, b i16}
/* testdata/in/parser_test.in:349 */type T = variant{left i8}
/* testdata/in/parser_test.in:350 */type T = variant{left i8, right i16}
/* testdata/in/parser_test.in:351-353 */type T = variant{left i8}
/* testdata/in/parser_test.in:354-357 */type T = variant{left i8, right i16}
