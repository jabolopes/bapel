module /* tests/testdata/in/traits_inherent_impl.in:1 */traits_inherent_impl

in "tests/testdata/in/traits_inherent_impl.in" in lines 3-6
impls {
  /* tests/testdata/in/traits_inherent_impl.in:4 */ "ptr.h"
  /* tests/testdata/in/traits_inherent_impl.in:5 */ "vector.h"
}

/* tests/testdata/in/traits_inherent_impl.in:8 */ pub type Ptr :: ∗ -> ∗

/* tests/testdata/in/traits_inherent_impl.in:11 */ pub type Vector :: ∗ -> ∗

in "tests/testdata/in/traits_inherent_impl.in" in lines 13-17 impl ['a] Vector 'a {
  in "tests/testdata/in/traits_inherent_impl.in" in lines 14-16 fn get(v: Ptr Self, i: i64) -> 'a {
  vector_get (v, i)
}
}

/* tests/testdata/in/traits_inherent_impl.in:20-23 */ pub type Point = struct{x: i64, y: i64}

in "tests/testdata/in/traits_inherent_impl.in" in lines 25-37 impl Point {
  in "tests/testdata/in/traits_inherent_impl.in" in lines 26-28 fn mk(x: i64, y: i64) -> Point {
  struct{x = x, y = y}
}
  in "tests/testdata/in/traits_inherent_impl.in" in lines 30-32 fn sum(p: Ptr Self) -> i64 {
  + (Ptr::get p.x, Ptr::get p.y)
}
  in "tests/testdata/in/traits_inherent_impl.in" in lines 34-36 fn add(a: Ptr Self, b: Ptr Point) -> Point {
  struct{x = + (Ptr::get a.x, Ptr::get b.x), y = + (Ptr::get a.y, Ptr::get b.y)}
}
}

/* tests/testdata/in/traits_inherent_impl.in:40-43 */ pub type Node = struct{val: i64, next: Ptr Node}

in "tests/testdata/in/traits_inherent_impl.in" in lines 45-53 impl Node {
  in "tests/testdata/in/traits_inherent_impl.in" in lines 46-48 fn get_val(n: Ptr Self) -> i64 {
  Ptr::get n.val
}
  in "tests/testdata/in/traits_inherent_impl.in" in lines 50-52 fn next_val(n: Ptr Self) -> i64 {
  Ptr::get n.next.get_val ()
}
}

/* tests/testdata/in/traits_inherent_impl.in:56-59 */ pub type Color = variant{red (), green ()}

in "tests/testdata/in/traits_inherent_impl.in" in lines 61-68 impl Color {
  in "tests/testdata/in/traits_inherent_impl.in" in lines 62-67 fn is_red(c: Ptr Self) -> bool {
  case Ptr::get c {
    red u -> true
    _ -> false
}
}
}

/* tests/testdata/in/traits_inherent_impl.in:71-74 */ pub type IntList = variant{nil (), cons (i64, Ptr IntList)}

in "tests/testdata/in/traits_inherent_impl.in" in lines 76-90 impl IntList {
  in "tests/testdata/in/traits_inherent_impl.in" in lines 77-82 fn is_nil(l: Ptr Self) -> bool {
  case Ptr::get l {
    nil u -> true
    _ -> false
}
}
  in "tests/testdata/in/traits_inherent_impl.in" in lines 84-89 fn sum(l: Ptr Self) -> i64 {
  case Ptr::get l {
    nil u -> 0
    cons pair -> + (pair.0, pair.1.sum ())
}
}
}

in "tests/testdata/in/traits_inherent_impl.in" in lines 92-96 pub fn test_vector(v: Ptr (Vector i64)) -> i64 {
  let a: i64 = v.get 0
  let b: i64 = Vector::get [i64] (v, 1)
  + (a, b)
}

in "tests/testdata/in/traits_inherent_impl.in" in lines 98-109 pub fn test_point() -> i64 {
  let a: i64 = Point::mk (1, 2).sum ()
  let p1: Point = Point::mk (3, 4)
  let b: i64 = p1.sum
  let ref_p1: Ptr Point = Ptr::mk p1
  let c: i64 = ref_p1.sum ()
  let d: i64 = ref_p1.sum
  let p2: Point = Point::mk (5, 6)
  let p3: Point = p1.add (Ptr::mk p2)
  let e: i64 = Point::sum (Ptr::mk p3)
  + (+ (+ (+ (a, b), c), d), e)
}

in "tests/testdata/in/traits_inherent_impl.in" in lines 111-118 pub fn test_node(n: Ptr Node) -> i64 {
  let a: i64 = n.get_val ()
  let b: i64 = n.get_val
  let c: i64 = n.next_val ()
  let d: i64 = n.next_val
  let e: i64 = Node::get_val n
  + (+ (+ (+ (a, b), c), d), e)
}

in "tests/testdata/in/traits_inherent_impl.in" in lines 120-129 pub fn test_color() -> bool {
  let c1: Color = variant{Color red = ()}
  let a: bool = c1.is_red ()
  let b: bool = c1.is_red
  let ref_c1: Ptr Color = Ptr::mk c1
  let c: bool = ref_c1.is_red ()
  let d: bool = ref_c1.is_red
  let e: bool = Color::is_red (Ptr::mk c1)
  && (&& (&& (&& (a, b), c), d), e)
}

in "tests/testdata/in/traits_inherent_impl.in" in lines 131-147 pub fn test_int_list(l: Ptr IntList) -> i64 {
  let empty: IntList = variant{IntList nil = ()}
  let a: bool = empty.is_nil ()
  let b: bool = empty.is_nil
  let c: bool = l.is_nil ()
  let d: bool = l.is_nil
  let s1: i64 = empty.sum ()
  let s2: i64 = empty.sum
  let s3: i64 = l.sum ()
  let s4: i64 = l.sum
  let s5: i64 = IntList::sum l
  ifelse (&& (&& (&& (a, b), c), d), {
  + (s1, s2)
}, {
  + (+ (s3, s4), s5)
})
}
