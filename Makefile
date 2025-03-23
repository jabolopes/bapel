all: bpl program
	./bpl query arr.bpl
	./bpl query core.bpl
	./bpl query str.bpl
	./bpl query vec.bpl
	./bpl query program.bpl

.PHONY: bpl
bpl:
	go generate "./..."
	go build "./..."
	go test "./..."
	staticcheck "./..."
	go build -o $@ ./bin

program: bpl
	g++ -c -std=c++20 -fmodules-ts -xc++-system-header ctime \
		-xc++-system-header array \
		-xc++-system-header cassert \
		-xc++-system-header cerrno \
		-xc++-system-header cstdint \
		-xc++-system-header cstdlib \
		-xc++-system-header functional \
		-xc++-system-header iostream \
		-xc++-system-header string \
		-xc++-system-header string_view \
		-xc++-system-header tuple \
		-xc++-system-header variant \
		-xc++-system-header vector

	./bpl -vmodule="build=2" -alsologtostderr build program.bpl

program2: bpl
	clang++ -std=c++20 -x c++-module -fprebuilt-module-path=out arr_impl.cc --precompile -o out/arr-arr_impl.pcm
	clang++ -std=c++20 -x c++-module -fprebuilt-module-path=out arr.cc --precompile -o out/arr.pcm

	clang++ -std=c++20 -fprebuilt-module-path=out -c out/arr-arr_impl.pcm -o out/arr-arr_impl.o
	clang++ -std=c++20 -fprebuilt-module-path=out -c out/arr.pcm -o out/arr.o


	clang++ -std=c++20 -x c++-module -fprebuilt-module-path=out core_ecs.cc --precompile -o out/core-core_ecs.pcm
	clang++ -std=c++20 -x c++-module -fprebuilt-module-path=out core_impl.cc --precompile -o out/core-core_impl.pcm
	clang++ -std=c++20 -x c++-module -fprebuilt-module-path=out core_ref.cc --precompile -o out/core-core_ref.pcm
	clang++ -std=c++20 -x c++-module -fprebuilt-module-path=out core.cc --precompile -o out/core.pcm

	clang++ -std=c++20 -fprebuilt-module-path=out -c out/core-core_ecs.pcm -o out/core-core_ecs.o
	clang++ -std=c++20 -fprebuilt-module-path=out -c out/core-core_impl.pcm -o out/core-core_impl.o
	clang++ -std=c++20 -fprebuilt-module-path=out -c out/core-core_ref.pcm -o out/core-core_ref.o
	clang++ -std=c++20 -fprebuilt-module-path=out -c out/core.pcm -o out/core.o


	clang++ -std=c++20 -x c++-module -fprebuilt-module-path=out str_impl.cc --precompile -o out/str-str_impl.pcm
	clang++ -std=c++20 -x c++-module -fprebuilt-module-path=out str.cc --precompile -o out/str.pcm

	clang++ -std=c++20 -fprebuilt-module-path=out -c out/str-str_impl.pcm -o out/str-str_impl.o
	clang++ -std=c++20 -fprebuilt-module-path=out -c out/str.pcm -o out/str.o


	clang++ -std=c++20 -x c++-module -fprebuilt-module-path=out vec_impl.cc --precompile -o out/vec-vec_impl.pcm
	clang++ -std=c++20 -x c++-module -fprebuilt-module-path=out vec.cc --precompile -o out/vec.pcm

	clang++ -std=c++20 -fprebuilt-module-path=out -c out/vec-vec_impl.pcm -o out/vec-vec_impl.o
	clang++ -std=c++20 -fprebuilt-module-path=out -c out/vec.pcm -o out/vec.o


	clang++ -std=c++20 -x c++-module -fprebuilt-module-path=out -Ientt/single_include -ISDL/include game_impl.cc --precompile -o out/game-game_impl.pcm
	clang++ -std=c++20 -x c++-module -fprebuilt-module-path=out game.cc --precompile -o out/game.pcm

	clang++ -std=c++20 -fprebuilt-module-path=out -c out/game-game_impl.pcm -o out/game-game_impl.o
	clang++ -std=c++20 -fprebuilt-module-path=out -c out/game.pcm -o out/game.o


	clang++ -std=c++20 -x c++-module -fprebuilt-module-path=out program_array.cc --precompile -o out/program-program_array.pcm
	clang++ -std=c++20 -x c++-module -fprebuilt-module-path=out program_point.cc --precompile -o out/program-program_point.pcm
	clang++ -std=c++20 -x c++-module -fprebuilt-module-path=out program_string.cc --precompile -o out/program-program_string.pcm
	clang++ -std=c++20 -x c++-module -fprebuilt-module-path=out program_vector.cc --precompile -o out/program-program_vector.pcm
	clang++ -std=c++20 -x c++-module -fprebuilt-module-path=out program.cc --precompile -o out/program.pcm

	clang++ -std=c++20 -fprebuilt-module-path=out -c out/program-program_array.pcm -o out/program-program_array.o
	clang++ -std=c++20 -fprebuilt-module-path=out -c out/program-program_point.pcm -o out/program-program_point.o
	clang++ -std=c++20 -fprebuilt-module-path=out -c out/program-program_string.pcm -o out/program-program_string.o
	clang++ -std=c++20 -fprebuilt-module-path=out -c out/program-program_vector.pcm -o out/program-program_vector.o
	clang++ -std=c++20 -fprebuilt-module-path=out -c out/program.pcm -o out/program.o


	clang++ -std=c++20 -fprebuilt-module-path=out -o out/program \
		-Wl,-rpath,SDL/build \
		-LSDL/build -lSDL3 \
		out/arr-arr_impl.o \
		out/arr.o \
		out/core-core_ecs.o \
		out/core-core_impl.o \
		out/core-core_ref.o \
		out/core.o \
		out/game-game_impl.o \
		out/game.o \
		out/program-program_array.o \
		out/program-program_point.o \
		out/program-program_string.o \
		out/program-program_vector.o \
		out/program.o \
		out/str-str_impl.o \
		out/str.o \
		out/vec-vec_impl.o \
		out/vec.o

debug:
	( cd bin; gdlv debug )

test:
	go test "./..."

regen:
	go test ./bplparser2/... -regen
