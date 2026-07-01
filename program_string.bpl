implements program

imports {
  bapel.core
  bapel.stl
}

fn mkString() -> () {
  let s: String = "...";
  ()
}

fn testStringView() -> () {
  let s: String = "hello";
  let sv: StringView = String::view s;

  core::print [bool] (StringView::empty sv);
  core::print [i64] (StringView::size sv);
  core::print [i8] (StringView::front sv);
  core::print [i8] (StringView::at (sv, 1));

  let sub: StringView = StringView::substr (sv, 1, 3);
  core::print [i64] (StringView::size sub);
  core::print [i8] (StringView::front sub);

  let empty_sv: StringView = StringView::substr (sv, 10, 3);
  core::print [bool] (StringView::empty empty_sv);

  ()
}

