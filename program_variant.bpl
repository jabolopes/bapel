implements program

imports {
  bapel.core
}

type Choice ['a] = variant {left 'a, right i32}

fn mkLeft['a](value: 'a) -> Choice 'a {
  let v: Choice 'a = variant { (Choice 'a) left = value }

  let v1: 'a = v.left
  let v2: 'a = v.0

  let v3: 'a = case v {
    left l = l
    right r = v1
  }

  v <- variant {(Choice 'a) left = value}

  let r: Choice 'a = v
  r
}

fn mkRight['a](value: i32) -> Choice 'a {
  let v: Choice 'a = variant { (Choice 'a) right = value }

  let v1: i32 = v.right
  let v2: i32 = v.1

  let v3: i32 = case v {
    left l = v2
    right r = r
  }

  v <- variant {(Choice 'a) right = value}

  let r: Choice 'a = v
  r
}

export type Maybe ['a] = variant {none (), some 'a}

export fn mkNone['a]() -> Maybe 'a {
  let v: Maybe 'a = variant { (Maybe 'a) none = () }

  let v1: () = v.none
  let v2: () = v.0

  let v3: () = case v {
    none l = l
    some r = v1
  }

  v <- variant {(Maybe 'a) none = ()}

  let r: Maybe 'a = v
  r
}

export fn mkSome['a](value: 'a) -> Maybe 'a {
  let v: Maybe 'a = variant { (Maybe 'a) some = value }

  let v1: 'a = v.some
  let v2: 'a = v.1

  let v3: 'a = case v {
    none l = v2
    some r = r
  }

  v <- variant {(Maybe 'a) some = value}

  let r: Maybe 'a = v
  r
}
