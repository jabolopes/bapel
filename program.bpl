import c.bpl

decls {
  print10 : i32 -> ()
  addints : (i32, i32) -> i32
  tuple12 : () -> (i8, i8)
  tuple10 : i16 -> (i16, i16)
  fconst : forall ['a, 'b] ('a, 'b) -> 'a
}

func id['a](i 'a) -> 'a {
  i
}

func fconst['a, 'b](i 'a, j 'b) -> 'a {
  i
}

func ftrue() -> i32 {
  1 [i32]
}

export func ns.myfunc() -> () {
  ()
}

func ns.myotherfunc() -> () {
  ()
}

export func assignPrint() -> () {
  let a1 i8 = 123
  c.print [i8] a1
}

func addVars() -> () {
  let a1 i16 = 1024
  let a2 i16 = 10
  a1 <- a1 + a2
  c.print [i16] a1
}

func addVarConstant() -> () {
  let a1 i32 = 2048
  a1 <- a1 + 10
  c.print [i32] a1
}

func addConstants() -> () {
  let a1 i32
  a1 <- 4096 + 10
  c.print [i32] a1
}

func ifs() -> () {
  let a1 i8 = 1
  if a1 {
    c.print [i8] 1
    return ()
  }

  a1 <- 0
  if !a1 {
    c.print [i8] 2
  }

  a1 <- 1
  if a1 {
    c.print [i8] 3
    return ()
  } else {
    c.print [i8] 0
    return ()
  }

  a1 <- 0
  if a1 {
    c.print [i8] 0
  } else {
    c.print [i8] 4
  }

  a1 <- 1
  if a1 {
    let a2 i32 = 2
    let a3 i32 = 3
    a2 <- addints (a2, a3)
    c.print [i32] a2
  } else {
    c.print [i8] 0
  }

  if id [i8] 1 {
    c.print [i8] 1
  } else {
    c.print [i8] 0
  }

  if fconst [i8, i8] (1, 2) {
    c.print [i8] 1
  } else {
    c.print [i8] 0
  }

  if ftrue () {
    c.print [i8] 1
  } else {
    c.print [i8] 0
  }

  if 1 [i8] {
    c.print [i8] 1
  }
}

func main() -> i32 {
  ns.myfunc
  ns.myotherfunc
  assignPrint
  addVars
  addVarConstant
  addConstants
  ifs

  let var5 i32 = 1024
  print10 var5

  var5 <- 10
  let var6 i32 = 22
  var5 <- addints (var5, var6)
  c.print [i32] var5

  let var1 i8 = 0
  let var2 i8 = 0
  (var1, var2) <- tuple12 ()
  c.print [i8] var1
  c.print [i8] var2

  let var3 i16 = 5
  let var4 i16 = 0
  var3 <- 5
  (var3, var4) <- tuple10 var3
  c.print [i16] var3
  c.print [i16] var4

  var1 <- 1
  var1 <- 0 - var1
  c.print [i8] var1

  let time i64 = 0
  let err i64 = 0
  (err, time) <- c.time ()

  c.print [i8] 99
  c.print [i64] err
  c.print [i64] time

  0 [i32]
}

func print10(a1 i32) -> () {
  let l1 i32 = a1 + 10
  c.print [i32] l1
}

func addints(a1 i32, a2 i32) -> i32 {
  a1 + a2
}

func tuple12() -> (i8, i8) {
  (1 [i8], 2 [i8])
}

func tuple10(a1 i16) -> (i16, i16) {
  (a1, 10 [i16])
}

func mkArray() -> [i32, 10] {
  let r [i32, 10]
  r
}

func getArray(a [i32, 10], i i64) -> i32 {
  let r i32 = Index.get a i
  r
}

func setArray(a [i32, 10], i i64, v i32) -> () {
  Index.set a i v
}

export type ExportedStruct = struct {a i8}

type Hello = struct {a i32, b i64}

func mkHello() -> Hello {
  let h Hello
  Index.set h a 1
  Index.set h b 2
  h
}

func mkHello2() -> Hello {
  let h Hello = struct { a = 1 [i32], b = 2 [i64] }
  h
}

func getStructByIndex(a Hello) -> i32 {
  let r i32 = Index.get a 0
  r
}

func setStructByIndex(a Hello, v i32) -> () {
  Index.set a 0 v
  ()
}

func getStructByID(a Hello) -> i64 {
  let r i64 = Index.get a b
  r
}

func setStructByID(a Hello) -> () {
  Index.set a b 0
  ()
}

func mkTuple() -> (i32, i32) {
  let r (i32, i32)
  Index.set r 0 1
  Index.set r 0 2
  r
}

func mkTuple2() -> (i32, i32) {
  let a (i32, i32) = (1, 2)
  let r (i32, i32)
  r <- a
}

func getTupleByIndex(a (i32, i32)) -> i32 {
  let r i32 = Index.get a 0
  r
}

func setTupleByIndex(a (i32, i32), b i32) -> () {
  Index.set a 0 b
  ()
}

type Choice ['a] = variant {left 'a, right i32}

func getLeftByLabel['a](c (Choice 'a)) -> 'a {
  let r 'a = Index.get c left
  r
}

func getRightByLabel['a](c (Choice 'a)) -> i32 {
  let r i32 = Index.get c right
  r
}

func getLeftByIndex['a](c (Choice 'a)) -> 'a {
  let r 'a = Index.get c 0
  r
}

func getRightByIndex['a](c (Choice 'a)) -> i32 {
  let r i32 = Index.get c 1
  r
}

func mkLeft['a](value 'a) -> (Choice 'a) {
  let r Choice 'a
  Index.set r left value
  r
}

func mkRight['a](value i32) -> (Choice 'a) {
  let r Choice 'a
  Index.set r right value
  r
}

func mkLeft2['a](value 'a) -> (Choice 'a) {
  let r Choice 'a
  Index.set r 0 value
  r
}

func mkRight2['a](value i32) -> (Choice 'a) {
  let r Choice 'a
  Index.set r 1 value
  r
}

func mkLeft3() -> () {
  let a Choice i8 = variant { (Choice i8) left = 10 }
  ()
}

func f['a](x 'a) -> () {
  f ['a] x
}

func foo() -> () {
  let var1 (i8, i8) = tuple12 ()
  c.print [(i8, i8)] var1
}

func mkPoint() -> c.Point {
  let r c.Point
  Index.set r x 1
  Index.set r x 2
  c.noopPoint r
  r
}

func pointX() -> i32 {
  let p c.Point
  let x i32 = Index.get p x
  c.noopPoint p
  x
}

func mkAbsPoint() -> c.AbsPoint {
  let p c.AbsPoint
  let x i32 = c.absPointX p
  c.noopAbsPoint p
  p
}

component [Hello, 100]

func addEntity() -> () {
  let e i64 = ecs.addEntity ()
  c.print [i64] e

  let v Hello
  let ok i8 = 0
  (v, ok) <- ecs.get [Hello] e
  ecs.set [Hello] (e, v)

  let it Hello_iterator = ecs.iterate [Hello, Hello_iterator] ()
  (e, v, ok) <- ecs.next [Hello, Hello_iterator] it

  ()
}
