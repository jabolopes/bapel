all: bpl program query

.PHONY: lexer
lexer:
	./bootstrap/bpl build cpp_lexer/lex.bpl

.PHONY: bpl
bpl: lexer
	go build "./..."
	go test "./..."
	staticcheck "./..."
	go build -o $@ ./bin

program: bpl
	./bpl -vmodule="module_actions=2" -alsologtostderr build program.bpl

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

bootstrap: lexer
	cp out/cpp_lexer.lex bootstrap/lexer
