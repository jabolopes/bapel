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
	./bpl -vmodule="build=2" -alsologtostderr build program.bpl

debug:
	( cd bin; gdlv debug )

test:
	go test "./..."

regen:
	go test ./bplparser2/... -regen
