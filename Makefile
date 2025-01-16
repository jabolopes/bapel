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
	./bpl cpp vec.bpl
	./bpl cpp -m program program_point.bpl
	./bpl cpp program.bpl

	g++ -c -std=c++20 -fmodules-ts -xc++-system-header ctime \
		-xc++-system-header array \
		-xc++-system-header cassert \
		-xc++-system-header cerrno \
		-xc++-system-header cstdint \
		-xc++-system-header cstdlib \
		-xc++-system-header functional \
		-xc++-system-header iostream \
		-xc++-system-header tuple \
		-xc++-system-header variant \
		-xc++-system-header vector

	g++ -std=c++20 -fmodules-ts -o main \
		core_impl.cpp \
		core_ecs.cpp \
		core.cpp \
		vec_impl.cpp \
		vec.cpp \
		program_point.cpp \
		program.cpp

debug:
	( cd bin; gdlv debug )
