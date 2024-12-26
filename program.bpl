import c.bpl

decls {
  print10 : i32 -> ()
  addints : (i32, i32) -> i32
  tuple12 : () -> (i8, i8)
  tuple10 : i16 -> (i16, i16)
  fconst : forall ['a, 'b] ('a, 'b) -> 'a
}

func id['a](i 'a) -> (r 'a) {
  r <- i
}

func fconst['a, 'b](i 'a, j 'b) -> (r 'a) {
  r <- i
}

func ftrue() -> (r i32) {
  r <- 1
}

export func ns.myfunc() -> () {
}

func ns.myotherfunc() -> () {
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
  }

  a1 <- 0
  if !a1 {
    c.print [i8] 2
  }

  a1 <- 1
  if a1 {
    c.print [i8] 3
  } else {
    c.print [i8] 0
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

func main() -> (r i32) {
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
}

func print10(a1 i32) -> () {
  let l1 i32 = a1 + 10
  c.print [i32] l1
}

func addints(a1 i32, a2 i32) -> (r1 i32) {
  r1 <- a1 + a2
}

func tuple12() -> (r1 i8, r2 i8) {
  r1 <- 1
  r2 <- 2
}

func tuple10(a1 i16) -> (r1 i16, r2 i16) {
  r1 <- a1
  r2 <- 10
}

func mkArray() -> (r1 [i32, 10]) {
}

func getArray(a [i32, 10], i i64) -> (r1 i32) {
  r1 <- Index.get a i
}

func setArray(a [i32, 10], i i64, v i32) -> () {
  Index.set a i v
}

export struct ExportedStruct{a i8}

struct Hello{a i32, b i64}

func mkStruct() -> (r Hello) {
  let h Hello
  Index.set h a 1
  Index.set h b 2
  r <- h
}

func getStructByIndex(a Hello) -> (r i32) {
  r <- Index.get a 0
}

func setStructByIndex(a Hello, v i32) -> () {
  Index.set a 0 v
}

func getStructByID(a Hello) -> (r i64) {
  r <- Index.get a b
}

func setStructByID(a Hello) -> () {
  Index.set a b 0
}

func mkTuple() -> (r (i32, i32)) {
  Index.set r 0 1
  Index.set r 0 2
}

func mkTuple2() -> (r (i32, i32)) {
  let a (i32, i32) = (1, 2)
  r <- a
}

func getTupleByIndex(a (i32, i32)) -> (r i32) {
  r <- Index.get a 0
}

func setTupleByIndex(a (i32, i32), b i32) -> () {
  Index.set a 0 b
}

type Choice ['a] {{left 'a, right i32}}

func getLeftByLabel['a](c (Choice 'a)) -> (r 'a) {
  r <- Index.get c left
}

func getRightByLabel['a](c (Choice 'a)) -> (r i32) {
  r <- Index.get c right
}

func getLeftByIndex['a](c (Choice 'a)) -> (r 'a) {
  r <- Index.get c 0
}

func getRightByIndex['a](c (Choice 'a)) -> (r i32) {
  r <- Index.get c 1
}

func mkLeft['a](value 'a) -> (r (Choice 'a)) {
  Index.set r left value
}

func mkRight['a](value i32) -> (r (Choice 'a)) {
  Index.set r right value
}

func mkLeft2['a](value 'a) -> (r (Choice 'a)) {
  Index.set r 0 value
}

func mkRight2['a](value i32) -> (r (Choice 'a)) {
  Index.set r 1 value
}

func f['a](x 'a) -> () {
  f ['a] x
}

func foo() -> () {
  let var1 (i8, i8) = tuple12 ()
  c.print [(i8, i8)] var1
}

func mkPoint() -> (p c.Point) {
  Index.set p x 1
  Index.set p x 2
  c.noopPoint p
}

func pointX() -> (x i32) {
  let p c.Point
  x <- Index.get p x
  c.noopPoint p
}

func mkAbsPoint() -> (p c.AbsPoint) {
  let x i32 = c.absPointX p
  c.noopAbsPoint p
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
}
