implements program

imports {
  bapel.core
  vec
}

fn mkVector() -> () {
  let v: ref.Ref (vec.Vector i8) = ref.mk [vec.Vector i8] (vec.mk [i8] ())

  vec.add [i8] (ref.get [vec.Vector i8] v, 10)
  let r: i8 = vec.get [i8] (ref.get [vec.Vector i8] v, 0)
  vec.set [i8] (ref.get [vec.Vector i8] v, 0, 10)

  let alias: ref.Ref (vec.Vector i8) = v

  let copy: vec.Vector i8 = ref.get [vec.Vector i8] v
  ()
}
