exports {
  type core.Point = struct {x i32, y i32}
  core.noopPoint : core.Point -> ()

  type core.AbsPoint
  core.mkAbsPoint : () -> core.AbsPoint
  core.absPointX : core.AbsPoint -> i32
  core.noopAbsPoint : core.AbsPoint -> ()

  core.time : () -> (i64, i64)
  core.print : forall ['a] 'a -> ()

  core.mkArray : forall ['a] () -> ['a, 10]

  ecs.addEntity : () -> i64
  ecs.get : forall ['a] i64 -> ('a, i8)
  ecs.set : forall ['a] (i64, 'a) -> ()
  ecs.iterate : forall ['a, 'b] () -> 'b
  ecs.next : forall ['a, 'b] 'b -> (i64, 'a, i8)
}

impls {
  core_impl.cc
  core_ecs.cc
}
