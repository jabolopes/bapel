# Compiler Test Harness Architecture & Test Customization Directives

This document defines the architecture for test customization and stage skipping in Bapel's golden test suite.

---

## 1. Motivation and Background

Bapel tests its compiler pipeline through golden test suites located under [`tests/testdata/`](../tests/testdata):

| Test Suite | Runner | Pipeline Stage Tested |
| :--- | :--- | :--- |
| `parse` | [`tests/parser_test.cc`](../tests/parser_test.cc) | Lexer & ANTLR AST parser |
| `stlc` (Infer) | [`tests/stlc_test.cc`](../tests/stlc_test.cc) | Lambda calculus type inference & elaboration |
| `stlc` (Typecheck) | [`tests/stlc_test.cc`](../tests/stlc_test.cc) | Lambda calculus core read-only verification |
| `infer` | [`tests/infer_test.cc`](../tests/infer_test.cc) | Bapel type inference, constraint solving & elaboration |
| `typecheck` | [`tests/typecheck_test.cc`](../tests/typecheck_test.cc) | Bapel core read-only verification kernel |
| `cpp` | [`tests/cpp_printer_test.cc`](../tests/cpp_printer_test.cc) | C++ code emission (`.h`, `_private.h`, `.cc`) |
| `validity` | [`tests/cpp_printer_test.cc`](../tests/cpp_printer_test.cc) | Clang compilation of generated C++ code |

### Current Limitations

Currently, test-specific behaviors and stage skips are hardcoded inside the C++ test runner loops:
* **Hardcoded typechecker options**: [`tests/typecheck_test.cc`](../tests/typecheck_test.cc#L21-L23) explicitly checks `if (tests::path_base(inFile) == "order.in") options.skip_undefined_term_checks = true;`.
* **Hardcoded codegen skips**: [`tests/cpp_printer_test.cc`](../tests/cpp_printer_test.cc#L17) checks `if (tests::path_base(inFile) == "order.in") continue;`.
* **Hardcoded validity skips**: [`tests/cpp_printer_test.cc`](../tests/cpp_printer_test.cc#L72-L75) skips files importing `bapel.core` or testing invalid returns (`array.cc`, `context1.cc`, `loops.cc`, `returns_bad1.cc`, etc.).

This approach causes friction:
1. Adding a new test with custom requirements requires modifying C++ runner source files.
2. The reason a test skips a stage is disconnected from the test source file.
3. Test runners cannot be easily extended with new stages without updating lists of file exclusions.

---

## 2. In-File Header Directives (`@directive`)

To decouple test case requirements from the test runner implementations, test files (`.in`, `.stlc.in`) declare their stage requirements and configuration using top-of-file comment directives.

### Syntax Specification

Directives are declared at the beginning of the file using single-line comment syntax:

```bapel
// @<directive_name>: <arguments>
```

Directives must appear in the file header before any module declarations or code tokens.

### Canonical Ordered Stages

Compiler pipeline stages are ordered linearly:

$$\text{parse} \longrightarrow \text{infer} \longrightarrow \text{typecheck} \longrightarrow \text{cpp\_codegen} \longrightarrow \text{cpp\_compile}$$

| Stage Identifier | Order | Pipeline Action | Test Suite |
| :--- | :---: | :--- | :--- |
| **`parse`** | 1 | Lexer, parser & AST construction | [`ParserTest.GoldenFiles`](../tests/parser_test.cc) |
| **`infer`** | 2 | Elaboration, constraint solving & zonking | [`InferTest.GoldenFiles`](../tests/infer_test.cc), [`StlcTest.InferTerm`](../tests/stlc_test.cc) |
| **`typecheck`** | 3 | Core read-only verification kernel | [`TypecheckTest.GoldenFiles`](../tests/typecheck_test.cc), [`StlcTest.TypecheckTerm`](../tests/stlc_test.cc) |
| **`cpp_codegen`** | 4 | C++ emitter (`.h`, `_private.h`, `.cc`) | [`CppPrinterTest.GoldenFiles`](../tests/cpp_printer_test.cc) |
| **`cpp_compile`** | 5 | Clang compilation of generated `.cc` files to `.o` | [`CppPrinterTest.IsValidCpp`](../tests/cpp_printer_test.cc) |

### The `infer` vs. `typecheck` Invariant

As defined in [`designs/typechecking.md`](typechecking.md), Stage 2 (`typecheck`) is a trusted verification kernel that treats the AST as strictly **read-only**.

* **`infer`**: Runs `comp::typecheck_source_file` with `skip_term_typechecker = true`. Outputs elaborated IR.
* **`typecheck`**: Runs `comp::typecheck_source_file` with `skip_term_typechecker = false`. Validates the elaborated IR.
* **Zero Golden Duplication**: Because the verification kernel must not mutate the AST, both `infer` and `typecheck` output must be byte-for-byte identical, sharing the exact same `.out` golden file. Comparing `infer` output against `typecheck` output verifies that the typechecker kernel remains idempotent and read-only.

### Standard Directives

#### 1. `@expect-error: <stage>`
Specifies that this is a negative test case expected to fail with a diagnostic error at `<stage>` (`parse`, `infer`, or `typecheck`).

* **Preceding Stages**: All stages preceding `<stage>` must succeed (`ok == true`).
* **Target Stage**: The specified `<stage>` must fail (`ok == false`). The emitted error diagnostic is diffed directly against the golden file (`.out` / `.bpl`).
* **Subsequent Stages**: All stages downstream of `<stage>` are automatically skipped.

*(If `@expect-error` is omitted, the test is positive: all executed stages are expected to succeed).*

#### 2. `@skip-stage: <stage1>, <stage2>, ...`
Instructs the test harness to skip one or more specific stages for a positive test (e.g. skipping `cpp_compile` when testing code that depends on external runtime modules).

Valid values: `parse`, `infer`, `typecheck`, `cpp_codegen`, `cpp_compile`.

#### 3. `@typecheck-option: <option>=<value>`
Configures specific typechecker options for this test file:
* `skip_undefined_term_checks=true|false`: Disables validation that all declared symbols have definitions.
* `skip_default_context=true|false`: Omits the standard primitive types and operator declarations.

---

## 3. Test Case Directive Mapping

Below is the complete inventory of existing test files and their required directives:

### 1. Compiler Tests ([`tests/testdata/in/`](../tests/testdata/in/))

| Test File | Directives | Rationale |
| :--- | :--- | :--- |
| [`order.in`](../tests/testdata/in/order.in) | `// @typecheck-option: skip_undefined_term_checks=true`<br>`// @skip-stage: cpp_codegen, cpp_compile` | Forward declarations without definitions (`decl f: () -> ()`). |
| [`array.in`](../tests/testdata/in/array.in) | `// @skip-stage: cpp_compile` | Imports runtime module `bapel.core`; cannot compile C++ in isolation. |
| [`context1.in`](../tests/testdata/in/context1.in) | `// @skip-stage: cpp_compile` | Imports runtime module `bapel.core`; cannot compile C++ in isolation. |
| [`context2.in`](../tests/testdata/in/context2.in) | `// @expect-error: infer` | Negative test: Duplicate symbol declaration from implementation file. |
| [`loops.in`](../tests/testdata/in/loops.in) | `// @skip-stage: cpp_compile` | Imports `bapel.core` (uses `core::for`); cannot compile C++ in isolation. |
| [`polymorphism.in`](../tests/testdata/in/polymorphism.in) | `// @skip-stage: cpp_compile` | Generic functions without concrete instantiation in a standalone translation unit. |
| [`coherence_violation.in`](../tests/testdata/in/coherence_violation.in) | `// @expect-error: infer` | Negative test: Trait coherence violation caught during inference. |
| [`returns_bad1.in`](../tests/testdata/in/returns_bad1.in) | `// @expect-error: infer` | Negative test: Return type mismatch caught during inference unification. |
| [`returns_bad2.in`](../tests/testdata/in/returns_bad2.in) | `// @expect-error: infer` | Negative test: Return type mismatch caught during inference unification. |
| [`traits_bounds_error.in`](../tests/testdata/in/traits_bounds_error.in) | `// @expect-error: typecheck` | Negative test: Unsatisfied trait bound caught during typecheck verification. |
| *All other 18 `in/*.in` files* | *(None — default positive)* | Positive feature tests compiling cleanly through all five stages. |

### 2. Parser Tests ([`tests/testdata/parse/in/`](../tests/testdata/parse/in/))

| Test File | Directives | Rationale |
| :--- | :--- | :--- |
| [`bad_token.in`](../tests/testdata/parse/in/bad_token.in) | `// @expect-error: parse` | Negative parser test: Invalid token character `@`. |
| [`bad_unterminated_block_comment.in`](../tests/testdata/parse/in/bad_unterminated_block_comment.in) | `// @expect-error: parse` | Negative parser test: Unclosed block comment `/*`. |
| [`bad_unterminated_raw_string.in`](../tests/testdata/parse/in/bad_unterminated_raw_string.in) | `// @expect-error: parse` | Negative parser test: Unclosed raw string literal. |
| [`bad_unterminated_rune.in`](../tests/testdata/parse/in/bad_unterminated_rune.in) | `// @expect-error: parse` | Negative parser test: Unclosed rune literal `'a`. |
| [`bad_unterminated_string.in`](../tests/testdata/parse/in/bad_unterminated_string.in) | `// @expect-error: parse` | Negative parser test: Unclosed string literal `"foo`. |
| [`parser_test.in`](../tests/testdata/parse/in/parser_test.in) | *(None)* | Positive parser test. |
| [`traits.in`](../tests/testdata/parse/in/traits.in) | *(None)* | Positive parser test. |

### 3. STLC Tests ([`tests/testdata/stlc/`](../tests/testdata/stlc/))

| Test File | Directives | Rationale |
| :--- | :--- | :--- |
| *All 11 `*.stlc.in` files* | `// @typecheck-option: skip_default_context=true`<br>`// @skip-stage: cpp_codegen, cpp_compile` | Isolated lambda calculus fixtures that test `infer` and `typecheck` without the Bapel prelude context or C++ codegen. |

---

## 4. Architecture & Implementation

### 1. `TestDirectives` Data Structure

Defined in [`tests/test_util.h`](../tests/test_util.h):

```cpp
namespace tests {

struct TestDirectives {
  std::set<std::string> skipped_stages;
  comp::TypecheckOptions typecheck_options;
  std::string expect_error_stage; // Empty if expecting success

  static int stage_index(const std::string& stage) {
    if (stage == "parse") return 1;
    if (stage == "infer") return 2;
    if (stage == "typecheck") return 3;
    if (stage == "cpp_codegen") return 4;
    if (stage == "cpp_compile") return 5;
    return 999;
  }

  bool expects_error(const std::string& stage) const {
    return expect_error_stage == stage;
  }

  bool should_run_stage(const std::string& stage) const {
    if (skipped_stages.find(stage) != skipped_stages.end()) {
      return false;
    }
    if (!expect_error_stage.empty()) {
      int stage_idx = stage_index(stage);
      int err_idx = stage_index(expect_error_stage);
      if (stage_idx > err_idx) {
        return false;
      }
    }
    return true;
  }

  static TestDirectives parse_from_file(const std::string& filepath) {
    TestDirectives dir;
    std::ifstream ifs(filepath);
    if (!ifs) return dir;

    std::string line;
    while (std::getline(ifs, line)) {
      // Trim whitespace
      size_t start = line.find_first_not_of(" \t\r\n");
      if (start == std::string::npos) continue;
      line = line.substr(start);

      // Stop scanning when reaching non-comment lines
      if (line.rfind("//", 0) != 0) {
        break;
      }

      size_t at_pos = line.find('@');
      if (at_pos == std::string::npos) continue;

      size_t colon_pos = line.find(':', at_pos);
      if (colon_pos == std::string::npos) continue;

      std::string directive = line.substr(at_pos + 1, colon_pos - (at_pos + 1));
      std::string args = line.substr(colon_pos + 1);

      // Trim args
      size_t a_start = args.find_first_not_of(" \t\r\n");
      if (a_start != std::string::npos) {
        args = args.substr(a_start);
      }
      size_t a_end = args.find_last_not_of(" \t\r\n");
      if (a_end != std::string::npos) {
        args = args.substr(0, a_end + 1);
      }

      if (directive == "expect-error") {
        dir.expect_error_stage = args;
      } else if (directive == "skip-stage") {
        std::stringstream ss(args);
        std::string stage;
        while (std::getline(ss, stage, ',')) {
          size_t s = stage.find_first_not_of(" \t\r\n");
          size_t e = stage.find_last_not_of(" \t\r\n");
          if (s != std::string::npos && e != std::string::npos) {
            dir.skipped_stages.insert(stage.substr(s, e - s + 1));
          }
        }
      } else if (directive == "typecheck-option") {
        if (args.find("skip_undefined_term_checks=true") != std::string::npos) {
          dir.typecheck_options.skip_undefined_term_checks = true;
        } else if (args.find("skip_default_context=true") != std::string::npos) {
          dir.typecheck_options.skip_default_context = true;
        }
      }
    }
    return dir;
  }
};

} // namespace tests
```

### 2. Integration into Test Runners

Test suites query the parsed directives directly:

* **[`tests/parser_test.cc`](../tests/parser_test.cc)**:
  ```cpp
  auto directives = tests::TestDirectives::parse_from_file(inFile);
  if (!directives.should_run_stage("parse")) return;

  auto res = parser::parse_source_file_from_file(inFile);
  if (directives.expects_error("parse")) {
    if (res.ok) {
      sub_ctx.add_error("expected parse error for " + inFile + " but parser succeeded", __FILE__, __LINE__);
      return;
    }
  } else {
    if (!res.ok) {
      sub_ctx.add_error("parser failed on " + inFile + ":\n" + res.error(), __FILE__, __LINE__);
      return;
    }
  }

  std::string got = res.ok ? res.value.to_string(true) : res.error();
  if (got.empty() || got.back() != '\n') got += "\n";
  std::string diff_str;
  if (!tests::diff_out_regen(got, wantFile, diff_str)) {
    sub_ctx.add_error("in test " + inFile + ":\n" + diff_str, __FILE__, __LINE__);
  }
  ```

* **[`tests/infer_test.cc`](../tests/infer_test.cc)**:
  ```cpp
  auto directives = tests::TestDirectives::parse_from_file(inFile);
  if (!directives.should_run_stage("infer")) return;

  comp::TypecheckOptions options = directives.typecheck_options;
  options.skip_term_typechecker = true;

  ir::IrUnit unit;
  std::string err;
  bool ok = comp::typecheck_source_file(querier, options, inFile, unit, err);

  if (directives.expects_error("infer")) {
    if (ok) {
      sub_ctx.add_error("expected infer error for " + inFile + " but inference succeeded", __FILE__, __LINE__);
      return;
    }
  } else {
    if (!ok) {
      sub_ctx.add_error("infer failed on " + inFile + ":\n" + err, __FILE__, __LINE__);
      return;
    }
  }

  std::string got = ok ? unit.to_bpl_string(true) : (err + "\n");
  std::string diff_str;
  if (!tests::diff_out_regen(got, wantFile, diff_str)) {
    sub_ctx.add_error("in test " + inFile + ":\n" + diff_str, __FILE__, __LINE__);
  }
  ```

* **[`tests/typecheck_test.cc`](../tests/typecheck_test.cc)**:
  ```cpp
  auto directives = tests::TestDirectives::parse_from_file(inFile);
  if (!directives.should_run_stage("typecheck")) return;

  comp::TypecheckOptions options = directives.typecheck_options;
  options.skip_term_typechecker = false;

  ir::IrUnit unit;
  std::string err;
  bool ok = comp::typecheck_source_file(querier, options, inFile, unit, err);

  if (directives.expects_error("typecheck")) {
    if (ok) {
      sub_ctx.add_error("expected typecheck error for " + inFile + " but typechecking succeeded", __FILE__, __LINE__);
      return;
    }
  } else {
    if (!ok) {
      sub_ctx.add_error("typecheck failed on " + inFile + ":\n" + err, __FILE__, __LINE__);
      return;
    }
  }

  std::string got = ok ? unit.to_bpl_string(true) : (err + "\n");
  std::string diff_str;
  if (!tests::diff_out_regen(got, wantFile, diff_str)) {
    sub_ctx.add_error("in test " + inFile + ":\n" + diff_str, __FILE__, __LINE__);
  }
  ```

* **[`tests/cpp_printer_test.cc`](../tests/cpp_printer_test.cc)**:
  ```cpp
  auto directives = tests::TestDirectives::parse_from_file(inFile);
  if (!directives.should_run_stage("cpp_codegen")) return;
  ```

* **`CppPrinterTest.IsValidCpp`**:
  ```cpp
  std::string in_file = tests::replace_string(tests::replace_extension(inFile, ".in"), "/cpp/", "/in/");
  auto directives = tests::TestDirectives::parse_from_file(in_file);
  if (!directives.should_run_stage("cpp_compile")) return;
  ```

---

## 5. Complementary CLI Controls

In addition to in-file directives, [`tests/test_main.cc`](../tests/test_main.cc) supports CLI stage filtering:

```bash
# Run only parsing across all test cases
./tests/test_runner --stage=parse

# Run inference only
./tests/test_runner --stage=infer

# Run all tests, skipping slow clang++ compilation
./tests/test_runner --skip-stage=cpp_compile

# Filter by test name and compiler stage
./tests/test_runner --filter=TypecheckTest --stage=typecheck
```

---

## 6. Benefits Summary

1. **Zero-Code Test Addition**: New edge cases can skip specific stages or customize options without modifying C++ test files.
2. **Read-Only Kernel Verification**: Separate `infer` and `typecheck` stages assert that the core verifier does not mutate the elaborated AST.
3. **Unified & Non-Redundant**: `@expect-error: <stage>` declares both the expected failure and automatically prunes downstream stages.
4. **Hermetic & Robust**: The test harness remains simple, fast, and unified across all compiler test directories.
