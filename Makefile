all: bpl program
	./bpl query c.bpl
	./bpl query vector.bpl
	./bpl query program.bpl

.PHONY: bpl
bpl:
	go generate "./..."
	go build "./..."
	go test "./..."
	staticcheck "./..."
	go build -o $@ ./bin

a.bpl.cpp: bpl
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
	./bpl cpp program.bpl
	clang-format -i a.bpl.cpp

program: a.bpl.cpp
	g++ -std=c++20 -fmodules-ts -o c.o -c c.cpp
	g++ -std=c++20 -fmodules-ts -o vector.o -c vector.cpp
	g++ -o $@ -std=c++20 -fmodules-ts vector.o c.o a.bpl.cpp

debug:
	( cd bin; gdlv debug )
