implements program

imports {
  core
}

fn mkPoint() -> core.Point {
  let r: core.Point = struct { x = 0 [i32], y = 0 [i32] }
  let x: i32 = r->x
  let y: i32 = r->y

  Index.set r x 1
  Index.set r x 2

  core.noopPoint r

  r
}

fn mkAbsPoint() -> core.AbsPoint {
  let p: core.AbsPoint = core.mkAbsPoint ()
  let x: i32 = core.absPointX p
  core.noopAbsPoint p
  p
}
