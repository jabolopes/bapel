.NOTPARALLEL:

all: bpl bootstrap/parser bootstrap/typechecker bootstrap/test_runner program query

bootstrap/parser: cpp_parser/main.cc cpp_parser/ast_builder.h cpp_parser/error_listener.h $(wildcard cpp_parser/generated/*.cpp) $(wildcard ast/*.h) $(wildcard bin/*.h)
	clang++ -O3 -std=c++17 -Wno-deprecated-declarations -I. -Icpp_parser -Icpp_parser/generated -I/usr/include/antlr4-runtime cpp_parser/main.cc cpp_parser/generated/bapelLexer.cpp cpp_parser/generated/bapelParser.cpp cpp_parser/generated/bapelBaseVisitor.cpp -lantlr4-runtime -o $@

bootstrap/typechecker: cpp_typechecker/main.cc $(wildcard comp/*.h) $(wildcard ts/*.h) $(wildcard bin/*.h)
	clang++ -O3 -std=c++17 -I. cpp_typechecker/main.cc -o $@

TEST_SRCS = $(wildcard tests/*_test.cc)

tests/test_runner: tests/test_main.cc tests/test_util.h $(TEST_SRCS) bootstrap/parser bootstrap/typechecker
	clang++ -O3 -std=c++17 -I. tests/test_main.cc $(TEST_SRCS) -o $@

bootstrap/test_runner: tests/test_runner
	cp $< $@

.PHONY: test test-cpp
test: tests/test_runner
	./tests/test_runner

test-cpp: tests/test_runner
	./tests/test_runner

.PHONY: regen
regen: tests/test_runner
	./tests/test_runner -regen

.PHONY: bpl
bpl: bootstrap/parser bootstrap/typechecker bootstrap/bpl tests/test_runner
	./tests/test_runner
	./bootstrap/bpl build bin.main
	rm -f $@
	cp out/bin.main $@

.PHONY: bootstrap
bootstrap: bpl bootstrap/parser bootstrap/typechecker bootstrap/test_runner
	cp bpl bootstrap/bpl

program: bpl
	./bpl build program

query: bpl
	./bpl query bapel/core
	./bpl query ./bapel/core.bpl
	./bpl query ./bapel/core_impl.h

debug:
	( cd bin; gdlv debug )

.PHONY: gen-parser
gen-parser:
	mkdir -p cpp_parser/generated
	antlr4 -Dlanguage=Cpp -visitor -no-listener -Xexact-output-dir -o cpp_parser/generated cpp_parser/bapel.g4

.PHONY: gen-parser-cpp
gen-parser-cpp: gen-parser


