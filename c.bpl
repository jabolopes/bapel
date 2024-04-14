exports {
  type Point = {x i32, y i32}
  c.noopPoint : Point -> ()

  type AbsPoint
  c.mkAbsPoint : () -> AbsPoint
  c.noopAbsPoint : AbsPoint -> ()

  c.time : () -> (i64, i64)
  c.print : forall ['a] 'a -> ()
}
