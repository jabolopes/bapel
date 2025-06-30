module game

imports {
  "bapel.core"
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

  "-DENTT_NO_ETO"
  "-Ientt/single_include"
}
