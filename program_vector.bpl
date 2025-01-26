implements program

imports {
  vec
}

fn mkVector() -> () {
  let v: vec.Vector i8 = vec.mk [i8] ()
  vec.add [i8] (v, 10)
  let r: i8 = vec.get [i8] (v, 0)
  vec.set [i8] (v, 0, 10)
}
