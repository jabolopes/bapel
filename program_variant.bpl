implements program

imports {
  core
  vec
}

type Choice ['a] = variant {left 'a, right i32}

fn mkLeft['a](value: 'a) -> (Choice 'a) {
  let v: Choice 'a = variant { (Choice 'a) left = value }

  let v1: 'a = v->left
  let v2: 'a = v->0

  let v3: 'a = case v {
    | left l = l
    | right r = v1
  }

  v <- variant {(Choice 'a) left = value}

  Index.set v left value
  Index.set v 0 value

  let r: Choice 'a = v
  r
}

fn mkRight['a](value: i32) -> (Choice 'a) {
  let v: Choice 'a = variant { (Choice 'a) right = value }

  let v1: i32 = v->right
  let v2: i32 = v->1

  let v3: i32 = case v {
    | left l = v2
    | right r = r
  }

  v <- variant {(Choice 'a) right = value}

  Index.set v right value
  Index.set v 1 value

  let r: Choice 'a = v
  r
}
