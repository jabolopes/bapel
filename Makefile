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

debug:
	( cd bin; gdlv debug )

test:
	go test "./..."

regen:
	go test ./bplparser2/... -regen
