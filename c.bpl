exports {
  type c.Point = struct {x i32, y i32}
  c.noopPoint : c.Point -> ()

  type c.AbsPoint
  c.mkAbsPoint : () -> c.AbsPoint
  c.absPointX : c.AbsPoint -> i32
  c.noopAbsPoint : c.AbsPoint -> ()

  c.time : () -> (i64, i64)
  c.print : forall ['a] 'a -> ()

  c.mkArray : forall ['a] () -> ['a, 10]

  ecs.addEntity : () -> i64
  ecs.get : forall ['a] i64 -> ('a, i8)
  ecs.set : forall ['a] (i64, 'a) -> ()
  ecs.iterate : forall ['a, 'b] () -> 'b
  ecs.next : forall ['a, 'b] 'b -> (i64, 'a, i8)
}
