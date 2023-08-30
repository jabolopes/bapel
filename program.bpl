imports {
  c.time : () -> (i64, i64)
}

exports {
  assignPrint : () -> ()
  ns.myfunc : () -> ()
}

decls {
  print10 : (i32) -> ()
  addints : (i32, i32) -> (i32)
  tuple12 : () -> (i8, i8)
  tuple10 : (i16) -> (i16, i16)
  type Hello : {a i32, b i64}
}

func id(i i32) -> (r i32) {
  r <- i
}

func ftrue() -> (r i32) {
  r <- 1
}

func ns.myfunc() -> () {
}

func ns.myotherfunc() -> () {
}

func assignPrint() -> () {
  let a1 i8
  a1 <- 123
  printU a1
}

func addVars() -> () {
  let a1 i16
  a1 <- 1024

  let a2 i16
  a2 <- 10

  a1 <- a1 + a2
  printU a1
}

func addVarConstant() -> () {
  let a1 i32
  a1 <- 2048
  a1 <- a1 + 10
  printU a1
}

func addConstants() -> () {
  let a1 i32
  a1 <- 4096 + 10
  printU a1
}

func ifs() -> () {
  let a1 i8
  a1 <- 1
  if a1 {
    printU i8 1
  }

  a1 <- 0
  if a1 else {
    printU i8 2
  }

  a1 <- 1
  if a1 {
    printU i8 3
  } else {
    printU i8 0
  }

  a1 <- 0
  if a1 {
    printU i8 0
  } else {
    printU i8 4
  }

  a1 <- 1
  if a1 {
    let a2 i32
    let a3 i32
    a2 <- 2
    a3 <- 3
    a2 <- addints a2 a3
    printU a2
  } else {
    printU i8 0
  }

  if id 1 {
    printU i8 1
  } else {
    printU i8 0
  }

  if ftrue {
    printU i8 1
  } else {
    printU i8 0
  }
}

func main() -> (r i32) {
  let var1 i8
  let var2 i8
  let var3 i16
  let var4 i16
  let var5 i32
  let var6 i32
  let var7 i64

  ns.myfunc
  ns.myotherfunc
  assignPrint
  addVars
  addVarConstant
  addConstants
  ifs

  var5 <- 1024
  print10 var5

  var5 <- 10
  var6 <- 22
  var5 <- addints var5 var6
  printU var5

  var1 var2 <- tuple12
  printU var1
  printU var2

  var3 <- 5
  var3 var4 <- tuple10 var3
  printU var3
  printU var4

  var1 <- 1
  var1 <- - var1
  printS var1

  let time i64
  let err i64
  err time <- c.time

  printU i8 99
  printU err
  printU time
}

func print10(a1 i32) -> () {
  let l1 i32
  l1 <- a1 + 10
  printU l1
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

func mkArray() -> (r1 [i32 10]) {
}

func getArray(a [i32 10], i i64) -> (r1 i32) {
  r1 <- Index.get a i
}

func setArray(a [i32 10], i i64, v i32) -> () {
  Index.set a i v
}

struct Hello{a i32, b i64}

func mkStruct() -> (r Hello) {
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

entity Hello {}
