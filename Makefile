CXX = clang++
CXXFLAGS = -O1 -std=c++20 -Wno-deprecated-declarations $(PARSER_INCLUDES)
TEST_OBJ_DIR = tests/.obj
TEST_SRCS = tests/test_main.cc $(wildcard tests/*_test.cc)
PARSER_GEN_SRCS = cpp_parser/generated/bapelLexer.cpp cpp_parser/generated/bapelParser.cpp cpp_parser/generated/bapelBaseVisitor.cpp
PARSER_INCLUDES = -I. -Icpp_parser -Icpp_parser/generated -I/usr/include/antlr4-runtime

TEST_OBJS = $(patsubst tests/%.cc,$(TEST_OBJ_DIR)/%.o,$(TEST_SRCS))
PARSER_OBJS = $(patsubst cpp_parser/generated/%.cpp,$(TEST_OBJ_DIR)/%.o,$(PARSER_GEN_SRCS))
ALL_TEST_OBJS = $(TEST_OBJS) $(PARSER_OBJS)
DEPS = $(ALL_TEST_OBJS:.o=.d)

-include $(DEPS)

.PHONY: all
all: bpl program query

$(TEST_OBJ_DIR):
	mkdir -p $(TEST_OBJ_DIR)

$(TEST_OBJ_DIR)/%.o: tests/%.cc | $(TEST_OBJ_DIR)
	$(CXX) $(CXXFLAGS) -MMD -MP -c $< -o $@

$(TEST_OBJ_DIR)/%.o: cpp_parser/generated/%.cpp | $(TEST_OBJ_DIR)
	$(CXX) $(CXXFLAGS) -MMD -MP -c $< -o $@

tests/test_runner: $(ALL_TEST_OBJS)
	$(CXX) $(ALL_TEST_OBJS) -lantlr4-runtime -o $@

.PHONY: test
test:
	$(MAKE) -j$$(nproc) tests/test_runner
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
	rm -rf out/ bpl tests/test_runner $(TEST_OBJ_DIR)

.PHONY: gen-parser
gen-parser:
	mkdir -p cpp_parser/generated
	antlr4 -Dlanguage=Cpp -visitor -no-listener -Xexact-output-dir -o cpp_parser/generated cpp_parser/bapel.g4

.PHONY: gen-parser-cpp
gen-parser-cpp: gen-parser


