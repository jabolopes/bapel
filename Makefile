.NOTPARALLEL:

HDRS = $(wildcard ast/*.h) $(wildcard bin/*.h) $(wildcard comp/*.h) $(wildcard ts/*.h)
TEST_SRCS = $(wildcard tests/*_test.cc)

.PHONY: all
all: bpl bootstrap/parser program query

bootstrap/parser: cpp_parser/main.cc cpp_parser/ast_builder.h cpp_parser/error_listener.h $(wildcard cpp_parser/generated/*.cpp) $(HDRS)
	clang++ -O3 -std=c++17 -Wno-deprecated-declarations -I. -Icpp_parser -Icpp_parser/generated -I/usr/include/antlr4-runtime cpp_parser/main.cc cpp_parser/generated/bapelLexer.cpp cpp_parser/generated/bapelParser.cpp cpp_parser/generated/bapelBaseVisitor.cpp -lantlr4-runtime -o $@

tests/test_runner: tests/test_main.cc tests/test_util.h $(TEST_SRCS) $(HDRS) bootstrap/parser
	clang++ -O3 -std=c++17 -I. tests/test_main.cc $(TEST_SRCS) -o $@

.PHONY: test
test: tests/test_runner
	./tests/test_runner

.PHONY: regen
regen: tests/test_runner
	./tests/test_runner -regen

.PHONY: bpl
bpl: bootstrap/parser bootstrap/bpl
	./bootstrap/bpl build bin.main
	rm -f $@
	cp out/bin.main $@

.PHONY: bootstrap
bootstrap: bpl bootstrap/parser
	cp bpl bootstrap/bpl

.PHONY: program
program: bpl
	./bpl build program

.PHONY: query
query: bpl
	./bpl query bapel/core
	./bpl query ./bapel/core.bpl
	./bpl query ./bapel/core_impl.h

.PHONY: clean
clean:
	rm -rf out/ bpl tests/test_runner

.PHONY: gen-parser
gen-parser:
	mkdir -p cpp_parser/generated
	antlr4 -Dlanguage=Cpp -visitor -no-listener -Xexact-output-dir -o cpp_parser/generated cpp_parser/bapel.g4

.PHONY: gen-parser-cpp
gen-parser-cpp: gen-parser


