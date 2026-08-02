.NOTPARALLEL:

HDRS = $(wildcard ast/*.h) $(wildcard bin/*.h) $(wildcard comp/*.h) $(wildcard ts/*.h) $(wildcard cpp_parser/*.h)
TEST_SRCS = $(wildcard tests/*_test.cc)
PARSER_GEN_SRCS = cpp_parser/generated/bapelLexer.cpp cpp_parser/generated/bapelParser.cpp cpp_parser/generated/bapelBaseVisitor.cpp
PARSER_INCLUDES = -I. -Icpp_parser -Icpp_parser/generated -I/usr/include/antlr4-runtime

.PHONY: all
all: bpl program query

tests/test_runner: tests/test_main.cc tests/test_util.h $(TEST_SRCS) $(HDRS) $(PARSER_GEN_SRCS)
	clang++ -O3 -std=c++17 -Wno-deprecated-declarations $(PARSER_INCLUDES) tests/test_main.cc $(TEST_SRCS) $(PARSER_GEN_SRCS) -lantlr4-runtime -o $@

.PHONY: test
test: tests/test_runner
	./tests/test_runner

.PHONY: regen
regen: tests/test_runner
	rm -f tests/testdata/parser/* tests/testdata/normalize/* tests/testdata/typecheck/* tests/testdata/cpp/* tests/testdata/stlc/*.out
	./tests/test_runner -regen

.PHONY: bpl
bpl: bootstrap/bpl
	./bootstrap/bpl build bin.main
	rm -f $@
	cp out/bin.main $@

.PHONY: bootstrap
bootstrap: bpl
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


