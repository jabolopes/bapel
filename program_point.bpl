implements program

imports {
  bapel.core
}

fn mkPoint() -> core::Point {
  let r: core::Point = struct {x = 0, y = 0}
  let x: i32 = r->x
  let y: i32 = r->y

  set r {x = 3, y = 4}
  r <- set r {x = 3, y = 4}

  set r {0 = 3, 1 = 4}
  r <- set r {0 = 3, 1 = 4}

  core::noopPoint r

  r
}

fn mkAbsPoint() -> core::AbsPoint {
  let p: core::AbsPoint = core::mkAbsPoint ()
  let x: i32 = core::absPointX p
  core::noopAbsPoint p
  p
}
