all: bpl bootstrap/parser program query

bootstrap/parser: $(wildcard cpp_parser/*.go) $(wildcard cpp_parser/parser/*.go)
	go build -o $@ ./cpp_parser

.PHONY: bpl
bpl: bootstrap/parser
	go build "./..."
	go test "./..."
	staticcheck $$(go list ./... | grep -v /cpp_parser/parser)
	go build -o $@ ./bin

program: bpl
	./bpl -vmodule="build=2" -alsologtostderr build program.bpl

query: bpl
	./bpl query bapel/core
	./bpl query ./bapel/core.bpl
	./bpl query ./bapel/core_impl.h

debug:
	( cd bin; gdlv debug )

test:
	go test -p 8 "./..."

regen:
	go test ./parse/... -regen
	go test ./comp/... -regen
	go test ./ts/stlc/... -regen



.PHONY: gen-parser
gen-parser:
	antlr4 -Dlanguage=Go -visitor -Xexact-output-dir -o cpp_parser/parser cpp_parser/bapel.g4
