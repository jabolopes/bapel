all: bpl program query

.PHONY: bpl
bpl:
	go generate "./..."
	go build "./..."
	go test "./..."
	staticcheck "./..."
	go build -o $@ ./bin

program: bpl
	./bpl -vmodule="build=2" -alsologtostderr build program.bpl

query: bpl
	./bpl query bapel/core
	./bpl query bapel/core.bpl
	./bpl query bapel/core_impl.cc

debug:
	( cd bin; gdlv debug )

test:
	go test "./..."

regen:
	go test ./bplparser2/... -regen
	go test ./comp/... -regen
	go test ./ts/stlc/... -regen
