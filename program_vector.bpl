implements program

imports {
  bapel.core
}

fn mkVector() -> () {
  let v: ref.Ref (core.Vector i8) = ref.mk [core.Vector i8] (core.mk [i8] ())

  core.add [i8] (ref.get [core.Vector i8] v, 10)
  let r: i8 = core.get [i8] (ref.get [core.Vector i8] v, 0)
  core.set [i8] (ref.get [core.Vector i8] v, 0, 10)

  let alias: ref.Ref (core.Vector i8) = v

  let copy: core.Vector i8 = ref.get [core.Vector i8] v
  ()
}
