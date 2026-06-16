all: bpl bootstrap/parser bootstrap/compiler bootstrap/querier program query

bootstrap/parser: $(wildcard cpp_parser/*.go) $(wildcard cpp_parser/parser/*.go)
	go build -o $@ ./cpp_parser

bootstrap/compiler: bootstrap/parser bin/cmd/compiler/compiler.go
	go build -o $@ ./bin/cmd/compiler/compiler.go

bootstrap/querier: bin/cmd/querier/querier.go
	go build -o $@ ./bin/cmd/querier/querier.go

.PHONY: bpl
bpl: bootstrap/parser bootstrap/compiler bootstrap/querier bootstrap/bpl
	go test "./..."
	staticcheck $$(go list ./... | grep -v /cpp_parser/parser)
	./bootstrap/bpl build bin.main
	rm -f $@
	cp out/bin.main $@

.PHONY: bootstrap
bootstrap: bpl bootstrap/parser bootstrap/compiler bootstrap/querier
	cp bpl bootstrap/bpl


program: bpl
	./bpl build program

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
