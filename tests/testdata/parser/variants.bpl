module /* tests/testdata/in/variants.in:1 */variants

/* tests/testdata/in/variants.in:3 */ type Choice :: ∗ -> ∗ = fun (a) (variant{left 'a, right i8})

in "tests/testdata/in/variants.in" in lines 5-25 fn mkLeft ['a](value: 'a) -> Choice 'a {
  let v: Choice 'a = variant{Choice 'a left = value}
  let v1: 'a = v.left
  let v2: 'a = v.0
  let v3: 'a = case v {
    left l -> l
    right r -> v1
}
  case v {
    left l -> l
    right r -> v1
}
  v <- variant{Choice 'a left = value}
  let r: Choice 'a = v
  r
}

in "tests/testdata/in/variants.in" in lines 27-47 fn mkRight ['a](value: i8) -> Choice 'a {
  let v: Choice 'a = variant{Choice 'a right = value}
  let v1: i8 = v.right
  let v2: i8 = v.1
  let v3: i8 = case v {
    left l -> v2
    right r -> r
}
  case v {
    left l -> v2
    right r -> r
}
  v <- variant{Choice 'a right = value}
  let r: Choice 'a = v
  r
}

/* tests/testdata/in/variants.in:51 */ export type Maybe :: ∗ -> ∗ = fun (a) (variant{none (), some 'a})

in "tests/testdata/in/variants.in" in lines 53-73 pub fn mkNone ['a]() -> Maybe 'a {
  let v: Maybe 'a = variant{Maybe 'a none = ()}
  let v1: () = v.none
  let v2: () = v.0
  let v3: () = case v {
    none l -> l
    some r -> v1
}
  case v {
    none l -> l
    some r -> v1
}
  v <- variant{Maybe 'a none = ()}
  let r: Maybe 'a = v
  r
}

in "tests/testdata/in/variants.in" in lines 75-95 pub fn mkSome ['a](value: 'a) -> Maybe 'a {
  let v: Maybe 'a = variant{Maybe 'a some = value}
  let v1: 'a = v.some
  let v2: 'a = v.1
  let v3: 'a = case v {
    none l -> v2
    some r -> r
}
  case v {
    none l -> v2
    some r -> r
}
  v <- variant{Maybe 'a some = value}
  let r: Maybe 'a = v
  r
}

/* tests/testdata/in/variants.in:99 */ export type One = variant{one i64}

in "tests/testdata/in/variants.in" in lines 101-119 pub fn mkOne() -> One {
  let v: One = variant{One one = 1}
  let v1: i64 = v.one
  let v2: i64 = v.0
  let v3: i64 = case v { one l -> l }
  case v { one l -> l }
  v <- variant{One one = 2}
  let r: One = v
  r
}
