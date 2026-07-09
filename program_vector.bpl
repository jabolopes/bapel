implements program

imports {
  bapel.core
  bapel.stl
}

fn mkVector() -> () {
  let vec: Vector i8 = Vector::mk [i8] ();

  Vector::push_back [i8] (&vec, 10);
  let r: i8 = Vector::get [i8] (&vec, 0);
  Vector::set_at [i8] (&vec, 0, 10);

  let copy: Vector i8 = vec;
  ()
}

fn testVectorSort() -> () {
  let vec: Vector i64 = Vector::mk [i64] ();
  Vector::push_back [i64] (&vec, 30);
  Vector::push_back [i64] (&vec, 10);
  Vector::push_back [i64] (&vec, 20);
  Vector::push_back [i64] (&vec, 10);
  Vector::push_back [i64] (&vec, 30);

  Vector::sort [i64] &vec;
  Vector::dedup [i64] &vec;
  core::print [i64] (Vector::size [i64] &vec);
  core::print [i64] (Vector::get [i64] (&vec, 0));
  core::print [i64] (Vector::get [i64] (&vec, 1));
  core::print [i64] (Vector::get [i64] (&vec, 2));
  ()
}

