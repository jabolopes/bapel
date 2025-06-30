module program

imports {
  core
  game
  vec
}

impls {
  program_array.bpl
  program_point.bpl
  program_string.bpl
  program_variant.bpl
  program_vector.bpl
}

fn id['a](i: 'a) -> 'a {
  i
}

fn fconst['a, 'b](i: 'a, j: 'b) -> 'a {
  i
}

fn ftrue() -> bool {
  true
}

export fn ns.myfunc() -> () {
  ()
}

fn ns.myotherfunc() -> () {
  ()
}

export fn assignPrint() -> () {
  let a1: i8 = 123
  core.print [i8] a1
}

fn addVars() -> () {
  let a1: i16 = 1024
  let a2: i16 = 10
  a1 <- a1 + a2
  core.print [i16] a1
}

fn addVarConstant() -> () {
  let a1: i32 = 2048
  a1 <- a1 + 10
  core.print [i32] a1
}

fn addConstants() -> () {
  let a1: i32 = 0
  a1 <- 4096 + 10
  core.print [i32] a1
}

fn ifs() -> () {
  let a1: i8 = 1
  if a1 != 0 {
    core.print [i8] 1
    return ()
  }

  a1 <- 0
  if a1 == 0 {
    core.print [i8] 2
  }

  a1 <- 1
  if a1 != 0 {
    core.print [i8] 3
    return ()
  } else {
    core.print [i8] 0
    return ()
  }

  a1 <- 0
  if a1 != 0 {
    core.print [i8] 0
  } else {
    core.print [i8] 4
  }

  a1 <- 1
  if a1 != 0 {
    let a2: i32 = 2
    let a3: i32 = 3
    a2 <- addints (a2, a3)
    core.print [i32] a2
  } else {
    core.print [i8] 0
  }

  if id [i8] 1 != 0 {
    core.print [i8] 1
  } else {
    core.print [i8] 0
  }

  if fconst [i8, i8] (1, 2) != 0{
    core.print [i8] 1
  } else {
    core.print [i8] 0
  }

  if ftrue () {
    core.print [i8] 1
  } else {
    core.print [i8] 0
  }

  if 1 [i8] != 0 {
    core.print [i8] 1
  }
}

type Player = struct {}

fn update () -> () {
  let delay: i64 = 16
  let e : std.optional (Entity, Player) = iterateAny [Player] ()
  ()
}

fn main() -> i32 {
  let material: Material = newRect (100, 100, 50, 50)
  let e1: Entity = init2 [Player, Material] (add (), struct {}, material)

  setUpdate update
  gameInit ()

  mkVector ()
  ns.myfunc ()
  ns.myotherfunc ()
  assignPrint ()
  addVars ()
  addVarConstant ()
  addConstants ()
  ifs ()

  let var5: i32 = 1024
  print10 var5

  var5 <- 10
  let var6: i32 = 22
  var5 <- addints (var5, var6)
  core.print [i32] var5

  let var1: i8 = 0
  let var2: i8 = 0
  (var1, var2) <- tuple12 ()
  core.print [i8] var1
  core.print [i8] var2

  let var3: i16 = 5
  let var4: i16 = 0
  var3 <- 5
  (var3, var4) <- tuple10 var3
  core.print [i16] var3
  core.print [i16] var4

  var1 <- 1
  var1 <- 0 - var1
  core.print [i8] var1

  let time: i64 = 0
  let err: i64 = 0
  (err, time) <- core.time ()

  core.print [i8] 99
  core.print [i64] err
  core.print [i64] time

  0
}

fn print10(a1: i32) -> () {
  let l1: i32 = a1 + 10
  core.print [i32] l1
}

fn addints(a1: i32, a2: i32) -> i32 {
  a1 + a2
}

fn tuple12() -> (i8, i8) {
  (1 [i8], 2 [i8])
}

fn tuple10(a1: i16) -> (i16, i16) {
  (a1, 10 [i16])
}

export type ExportedStruct = struct {a i8}

type Hello = struct {a i32, b i64}

fn mkHello() -> Hello {
  let h: Hello = struct {a = 1, b = 2}

  h <- set h {a = 3, b = 4}

  let v1: i32 = h->a
  let v2: i64 = h->b
  let v3: i32 = h->0
  let v4: i64 = h->1

  h
}

fn mkTuple() -> (i32, i64) {
  let t: (i32, i64) = (1, 2)

  let v1: i32 = t->0
  let v2: i64 = t->1

  set t {0 = 3, 1 = 4}
  t <- set t {0 = 3, 1 = 4}

  let r: (i32, i64) = t
  r
}

fn f['a](x: 'a) -> () {
  f ['a] x
}

fn foo() -> () {
  let var1: (i8, i8) = tuple12 ()
  core.print [(i8, i8)] var1
}

component [Hello, 100]

fn addEntity() -> () {
  let e: i64 = ecs.addEntity ()
  core.print [i64] e

  let v: Hello = struct {a = 0, b = 0}
  let ok: i8 = 0
  (v, ok) <- ecs.get [Hello] e
  ecs.set [Hello] (e, v)

  let it: Hello_iterator = ecs.iterate [Hello, Hello_iterator] ()
  (e, v, ok) <- ecs.next [Hello, Hello_iterator] it

  ()
}

fn lambda() -> i32 {
  let add: i32 -> i32 = \ x: i32 = x + 1 [i32]
  add 2
}

/*
 * TODO: Finish.
fn polymorphicLambda['a](value: 'a) -> 'a {
  let id2: forall ['b] 'b -> 'b = \ ['b] x: 'b = x
  id2 ['a] value
}
 */

fn functionSubtyping1() -> () {
  let id2: i8 -> i8 = id [i8]
  ()
}

fn functionSubtyping2['a](x: 'a -> 'a) -> ('a -> 'a) {
  id ['a -> 'a] (id ['a])
}
