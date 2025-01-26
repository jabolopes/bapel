implements program

imports {
  arr
}

fn mkArray() -> [i32, 10] {
  let a: [i32, 10] = arr.mk [i32] ()

  let v1: i32 = a->0
  let v2: i32 = arr.get [i32] (a, 0)

  let i: i64 = 0 [i64]
  arr.set [i32] (a, i, 10 [i32])
  Index.set a i (10 [i32])

  a
}
