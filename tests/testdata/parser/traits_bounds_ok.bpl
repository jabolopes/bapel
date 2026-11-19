module /* tests/testdata/in/traits_bounds_ok.in:1 */traits_bounds_ok

in "tests/testdata/in/traits_bounds_ok.in" in lines 3-6
impls {
  /* tests/testdata/in/traits_bounds_ok.in:4 */ "ptr.h"
  /* tests/testdata/in/traits_bounds_ok.in:5 */ "vector.h"
}

/* tests/testdata/in/traits_bounds_ok.in:8 */ export type Ptr :: ∗ -> ∗

in "tests/testdata/in/traits_bounds_ok.in" in lines 10-12 pub trait Size {
  in "tests/testdata/in/traits_bounds_ok.in" in line 11 fn size(s: Ptr Self) -> i64
}

/* tests/testdata/in/traits_bounds_ok.in:15 */ export type Vector :: ∗ -> ∗

in "tests/testdata/in/traits_bounds_ok.in" in lines 17-21 impl ['a] Size for Vector 'a {
  in "tests/testdata/in/traits_bounds_ok.in" in lines 18-20 fn size(v: Ptr Self) -> i64 {
  0
}
}

/* tests/testdata/in/traits_bounds_ok.in:24-27 */ export type Point = struct{x: i64, y: i64}

in "tests/testdata/in/traits_bounds_ok.in" in lines 29-33 impl Size for Point {
  in "tests/testdata/in/traits_bounds_ok.in" in lines 30-32 fn size(p: Ptr Self) -> i64 {
  + (Ptr::get p.x, Ptr::get p.y)
}
}

/* tests/testdata/in/traits_bounds_ok.in:36-39 */ export type Node = struct{val: i64, next: Ptr Node}

in "tests/testdata/in/traits_bounds_ok.in" in lines 41-45 impl Size for Node {
  in "tests/testdata/in/traits_bounds_ok.in" in lines 42-44 fn size(n: Ptr Self) -> i64 {
  Ptr::get n.val
}
}

/* tests/testdata/in/traits_bounds_ok.in:48-51 */ export type Color = variant{red (), green ()}

in "tests/testdata/in/traits_bounds_ok.in" in lines 53-60 impl Size for Color {
  in "tests/testdata/in/traits_bounds_ok.in" in lines 54-59 fn size(c: Ptr Self) -> i64 {
  case Ptr::get c {
    red u -> 1
    _ -> 2
}
}
}

/* tests/testdata/in/traits_bounds_ok.in:63-66 */ export type IntList = variant{nil (), cons (i64, Ptr IntList)}

in "tests/testdata/in/traits_bounds_ok.in" in lines 68-75 impl Size for IntList {
  in "tests/testdata/in/traits_bounds_ok.in" in lines 69-74 fn size(l: Ptr Self) -> i64 {
  case Ptr::get l {
    nil u -> 0
    cons pair -> pair.0
}
}
}

in "tests/testdata/in/traits_bounds_ok.in" in lines 77-79 pub fn print_size ['t: Size](x: Ptr 't) -> i64 {
  Size::size x
}

in "tests/testdata/in/traits_bounds_ok.in" in lines 81-83 pub fn print_size_dot ['t: Size](x: Ptr 't) -> i64 {
  x.size ()
}

in "tests/testdata/in/traits_bounds_ok.in" in lines 85-87 pub fn test_vector(v: Ptr (Vector i64)) -> i64 {
  + (print_size [Vector i64] v, print_size_dot [Vector i64] v)
}

in "tests/testdata/in/traits_bounds_ok.in" in lines 89-92 pub fn test_point() -> i64 {
  let p: Point = struct{x = 10, y = 32}
  + (print_size [Point] (Ptr::mk p), print_size_dot [Point] (Ptr::mk p))
}

in "tests/testdata/in/traits_bounds_ok.in" in lines 94-96 pub fn test_node(n: Ptr Node) -> i64 {
  print_size [Node] n
}

in "tests/testdata/in/traits_bounds_ok.in" in lines 98-101 pub fn test_color() -> i64 {
  let c: Color = variant{Color red = ()}
  print_size [Color] (Ptr::mk c)
}

in "tests/testdata/in/traits_bounds_ok.in" in lines 103-105 pub fn test_int_list(l: Ptr IntList) -> i64 {
  print_size [IntList] l
}
