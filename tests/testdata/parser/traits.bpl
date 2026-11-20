module /* tests/testdata/in/traits.in:1 */traits

in "tests/testdata/in/traits.in" in lines 3-6
impls {
  /* tests/testdata/in/traits.in:4 */ "ptr.h"
  /* tests/testdata/in/traits.in:5 */ "vector.h"
}

/* tests/testdata/in/traits.in:8 */ pub type Ptr :: ∗ -> ∗

in "tests/testdata/in/traits.in" in lines 10-12 pub trait Size {
  in "tests/testdata/in/traits.in" in line 11 fn size(s: Ptr Self) -> i64
}

in "tests/testdata/in/traits.in" in lines 14-16 pub trait Add {
  in "tests/testdata/in/traits.in" in line 15 fn add(a: Ptr Self, b: Ptr Self) -> i64
}

in "tests/testdata/in/traits.in" in lines 18-20 pub trait Indexable ['elem] {
  in "tests/testdata/in/traits.in" in line 19 fn get(v: Ptr Self, index: i64) -> 'elem
}

/* tests/testdata/in/traits.in:23 */ pub type Vector :: ∗ -> ∗

in "tests/testdata/in/traits.in" in lines 25-29 impl ['a] Indexable 'a for Vector 'a {
  in "tests/testdata/in/traits.in" in lines 26-28 fn get(v: Ptr Self, index: i64) -> 'a {
  vector_get (v, index)
}
}

/* tests/testdata/in/traits.in:32-35 */ pub type Point = struct{x: i64, y: i64}

in "tests/testdata/in/traits.in" in lines 37-41 impl Size for Point {
  in "tests/testdata/in/traits.in" in lines 38-40 fn size(p: Ptr Self) -> i64 {
  + (Ptr::get p.x, Ptr::get p.y)
}
}

in "tests/testdata/in/traits.in" in lines 43-47 impl Add for Point {
  in "tests/testdata/in/traits.in" in lines 44-46 fn add(a: Ptr Self, b: Ptr Self) -> i64 {
  + (+ (+ (Ptr::get a.x, Ptr::get a.y), Ptr::get b.x), Ptr::get b.y)
}
}

in "tests/testdata/in/traits.in" in lines 49-51 pub fn make_point(x: i64, y: i64) -> Point {
  struct{x = x, y = y}
}

/* tests/testdata/in/traits.in:54-57 */ pub type Node = struct{val: i64, next: Ptr Node}

in "tests/testdata/in/traits.in" in lines 59-63 impl Size for Node {
  in "tests/testdata/in/traits.in" in lines 60-62 fn size(n: Ptr Self) -> i64 {
  + (Ptr::get n.val, Size::size Ptr::get n.next)
}
}

/* tests/testdata/in/traits.in:66-69 */ pub type Color = variant{red (), green ()}

in "tests/testdata/in/traits.in" in lines 71-78 impl Size for Color {
  in "tests/testdata/in/traits.in" in lines 72-77 fn size(c: Ptr Self) -> i64 {
  case Ptr::get c {
    red u -> 1
    _ -> 2
}
}
}

/* tests/testdata/in/traits.in:81-84 */ pub type IntList = variant{nil (), cons (i64, Ptr IntList)}

in "tests/testdata/in/traits.in" in lines 86-93 impl Size for IntList {
  in "tests/testdata/in/traits.in" in lines 87-92 fn size(l: Ptr Self) -> i64 {
  case Ptr::get l {
    nil u -> 0
    cons pair -> + (pair.0, Size::size pair.1)
}
}
}

in "tests/testdata/in/traits.in" in lines 95-99 pub fn test_vector(v: Ptr (Vector i64)) -> i64 {
  let a: i64 = v.get 0
  let b: i64 = Indexable::get (v, 1)
  + (a, b)
}

in "tests/testdata/in/traits.in" in lines 101-113 pub fn test_point() -> i64 {
  let tmp: i64 = make_point (1, 2).size ()
  let p: Point = make_point (3, 4)
  let p2: Point = make_point (5, 6)
  let a: i64 = p.size ()
  let b: i64 = p.size
  let ref_p: Ptr Point = Ptr::mk p
  let c: i64 = ref_p.size ()
  let d: i64 = ref_p.size
  let e: i64 = Size::size (Ptr::mk p)
  let f: i64 = p.add (Ptr::mk p2)
  + (+ (+ (+ (+ (+ (tmp, a), b), c), d), e), f)
}

in "tests/testdata/in/traits.in" in lines 115-120 pub fn test_node(n: Ptr Node) -> i64 {
  let a: i64 = n.size ()
  let b: i64 = n.size
  let c: i64 = Size::size n
  + (+ (a, b), c)
}

in "tests/testdata/in/traits.in" in lines 122-131 pub fn test_color() -> i64 {
  let col: Color = variant{Color red = ()}
  let a: i64 = col.size ()
  let b: i64 = col.size
  let ref_col: Ptr Color = Ptr::mk col
  let c: i64 = ref_col.size ()
  let d: i64 = ref_col.size
  let e: i64 = Size::size (Ptr::mk col)
  + (+ (+ (+ (a, b), c), d), e)
}

in "tests/testdata/in/traits.in" in lines 133-141 pub fn test_int_list(l: Ptr IntList) -> i64 {
  let empty: IntList = variant{IntList nil = ()}
  let a: i64 = empty.size ()
  let b: i64 = empty.size
  let c: i64 = l.size ()
  let d: i64 = l.size
  let e: i64 = Size::size l
  + (+ (+ (+ (a, b), c), d), e)
}
