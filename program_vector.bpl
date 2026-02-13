implements program

imports {
  bapel.core
}

fn mkVector() -> () {
  let v: ref.Ref (core.Vector i8) = ref.mk [core.Vector i8] (core.mk [i8] ())

  let b: ref.Borrow (core.Vector i8) = ref.get v

  core.add [i8] (b, 10)
  let r: i8 = core.get [i8] (b, 0)
  core.set [i8] (b, 0, 10)

  let alias: ref.Ref (core.Vector i8) = v

  let copy: core.Vector i8 = core.vec_copy b
  ()
}

fn mkVectorSynchronized() -> () {
  let v: ref.Ref (core.Vector i8) = ref.mk [core.Vector i8] (core.mk [i8] ())

  let l: ref.Guard = ref.lock v
  let b: ref.Borrow (core.Vector i8) = ref.get v

  core.add [i8] (b, 10)
  let r: i8 = core.get [i8] (b, 0)
  core.set [i8] (b, 0, 10)

  let alias: ref.Ref (core.Vector i8) = v

  let copy: core.Vector i8 = core.vec_copy b
  ()
}
