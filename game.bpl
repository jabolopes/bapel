exports {
  type Entity
  add: () -> Entity
  init: forall ['a] (Entity, 'a) -> Entity

  type Material
  newRect: (i64, i64, i64, i64) -> Material
  addMaterial: Material -> ()

  setUpdate: (() -> ()) -> ()
  gameInit: () -> i64
}

impls {
  game_game.cc
  game_impl.cc
  game_material.cc
}

flags {
  "-ISDL/include"
  "-Wl,-rpath,SDL/build"
  "-LSDL/build"
  "-lSDL3"

  "-Ientt/single_include"
}
