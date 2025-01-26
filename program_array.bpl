implements program

imports {
  arr
}

fn mkArray() -> [i32, 10] {
  let a: [i32, 10] = arr.mk [i32] ()

  let v: i32 = a->0

  let r1: i32 = arr.get [i32] (a, 0)
  let r2: i32 = Index.get a 0

  let i: i64 = 0 [i64]
  arr.set [i32] (a, i, 10 [i32])
  Index.set a i (10 [i32])

  a
}
