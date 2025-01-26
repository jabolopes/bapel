all: bpl main
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

main: bpl
	./bpl cc core.bpl
	./bpl cc arr.bpl
	./bpl cc str.bpl
	./bpl cc vec.bpl
	./bpl cc program_array.bpl
	./bpl cc program_point.bpl
	./bpl cc program_string.bpl
	./bpl cc program_vector.bpl
	./bpl cc program.bpl

	clang-format -i core.cc
	clang-format -i arr.cc
	clang-format -i str.cc
	clang-format -i vec.cc
	clang-format -i program_array.cc
	clang-format -i program_point.cc
	clang-format -i program_string.cc
	clang-format -i program_vector.cc
	clang-format -i program.cc

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

	g++ -std=c++20 -fmodules-ts -o main \
		core_impl.cc \
		core_ecs.cc \
		core.cc \
		arr_impl.cc \
		arr.cc \
		str_impl.cc \
		str.cc \
		vec_impl.cc \
		vec.cc \
		program_array.cc \
		program_point.cc \
		program_string.cc \
		program_vector.cc \
		program.cc

debug:
	( cd bin; gdlv debug )

test:
	go test "./..."

regen:
	go test ./bplparser2/... -regen
