implements program

imports {
  bapel.core
  bapel.stl
}

fn mkVector() -> () {
  let vec: Vector i8 = Vector::mk [i8] ();

  Vector::push_back [i8] (&vec, 10);
  let r: i8 = Vector::get [i8] (&vec, 0);
  Vector::set [i8] (&vec, 0, 10);

  let copy: Vector i8 = vec;
  ()
}
