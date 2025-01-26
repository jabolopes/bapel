all: bpl main
	./bpl query core.bpl
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
	./bpl cpp core.bpl
	./bpl cpp str.bpl
	./bpl cpp vec.bpl
	./bpl cpp program_point.bpl
	./bpl cpp program_string.bpl
	./bpl cpp program.bpl

	clang-format -i core.cpp
	clang-format -i str.cpp
	clang-format -i vec.cpp
	clang-format -i program_point.cpp
	clang-format -i program_string.cpp
	clang-format -i program.cpp

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
		core_impl.cpp \
		core_ecs.cpp \
		core.cpp \
		str_impl.cpp \
		str.cpp \
		vec_impl.cpp \
		vec.cpp \
		program_point.cpp \
		program_string.cpp \
		program.cpp

debug:
	( cd bin; gdlv debug )

test:
	go test "./..."

regen:
	go test ./bplparser2/... -regen
