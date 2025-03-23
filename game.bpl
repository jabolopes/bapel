exports {
  type Material
  newRect: (i64, i64, i64, i64) -> Material
  addMaterial: Material -> ()

  gameInit: () -> i64
}

impls {
  game_impl.cc
}

flags {
  "-ISDL/include"
  "-Wl,-rpath,SDL/build"
  "-LSDL/build"
  "-lSDL3"

  "-Ientt/single_include"
}
