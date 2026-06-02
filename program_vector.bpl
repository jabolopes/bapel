implements program

imports {
  bapel.core
}

fn mkVector() -> () {
  let vec: core::Vector i8 = core::mk [i8] ();

  core::add [i8] (&vec, 10);
  let r: i8 = core::get [i8] (&vec, 0);
  core::set [i8] (&vec, 0, 10);

  let copy: core::Vector i8 = vec;
  ()
}
