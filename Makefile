all: bpl program query

.PHONY: bpl
bpl:
	go generate "./..."
	go build "./..."
	go test "./..."
	staticcheck "./..."
	go build -o $@ ./bin

program: bpl
	./bpl -vmodule="module_actions=2" -alsologtostderr build program.bpl

query: bpl
	./bpl query bapel/core
	./bpl query ./bapel/core.bpl
	./bpl query ./bapel/core_impl.ccm

debug:
	( cd bin; gdlv debug )

test:
	go test "./..."

regen:
	go test ./parse/... -regen
	go test ./comp/... -regen
	go test ./ts/stlc/... -regen
