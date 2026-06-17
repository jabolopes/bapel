implements program

imports {
  bapel.core
  bapel.stl
}

fn mkVector() -> () {
  let vec: Vector i8 = Vector_::mk [i8] ();

  Vector_::push_back [i8] (&vec, 10);
  let r: i8 = Vector_::get [i8] (&vec, 0);
  Vector_::set [i8] (&vec, 0, 10);

  let copy: Vector i8 = vec;
  ()
}
