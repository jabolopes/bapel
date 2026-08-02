all: bpl bootstrap/parser bootstrap/typechecker program query

bootstrap/parser: cpp_parser/main.cc cpp_parser/ast_builder.h cpp_parser/error_listener.h $(wildcard cpp_parser/generated/*.cpp) $(wildcard ast/*.h) $(wildcard bin/*.h)
	clang++ -O3 -std=c++17 -Wno-deprecated-declarations -I. -Icpp_parser -Icpp_parser/generated -I/usr/include/antlr4-runtime cpp_parser/main.cc cpp_parser/generated/bapelLexer.cpp cpp_parser/generated/bapelParser.cpp cpp_parser/generated/bapelBaseVisitor.cpp -lantlr4-runtime -o $@



bootstrap/typechecker: bootstrap/parser cpp_typechecker/main.cc $(wildcard comp/*.h) $(wildcard ts/*.h) $(wildcard bin/*.h)
	clang++ -O3 -std=c++17 -I. cpp_typechecker/main.cc -o $@

.PHONY: bpl
bpl: bootstrap/parser bootstrap/typechecker bootstrap/bpl
	go test "./..."
	staticcheck $$(go list ./...)
	./bootstrap/bpl build bin.main
	rm -f $@
	cp out/bin.main $@

.PHONY: bootstrap
bootstrap: bpl bootstrap/parser bootstrap/typechecker
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
	mkdir -p cpp_parser/generated
	antlr4 -Dlanguage=Cpp -visitor -no-listener -Xexact-output-dir -o cpp_parser/generated cpp_parser/bapel.g4

.PHONY: gen-parser-cpp
gen-parser-cpp: gen-parser


