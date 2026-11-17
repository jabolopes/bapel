# Self-Hosting Migration Analysis: Next Component & Blockers

This document complements [`designs/boot.md`](boot.md) with a concrete technical assessment of the current hybrid compiler codebase, recommending the immediate next component to migrate from C++ to Bapel and cataloging all language, compiler, and standard library blockers that stand in the way of full self-hosting.

---

## 1. Current State Summary

Currently, the Bapel compiler (`bpl`) is split across two layers:

* **Written in Bapel (`bin/`, `bapel/`)**:
  * CLI driver, Bazel workspace/BUILD generator, and module build graph ([`bin/main.bpl`](../bin/main.bpl)).
  * Module/file query engine and `workspace.bpl` package resolution ([`bin/query.bpl`](../bin/query.bpl)).
  * Standard library wrappers (`String`, `StringView`, `Vector`, `Deque`, `Optional`, `UnorderedMap`, `SharedPtr`, `UniquePtr`, `Ifstream`, `Ofstream`, `IStringStream`, `fs`, `os`, `args`).
* **Written in C++ (`ast/`, `comp/`, `ts/`, `cpp_parser/`, `bin/ir_*.h`, `bin/codegen_impl.h`)**:
  * ANTLR4 lexer/parser and AST builder ([`cpp_parser/`](../cpp_parser), [`ast/`](../ast)).
  * Syntactic desugaring, import/impl resolution, and trait coherence ([`comp/`](../comp)).
  * Bidirectional type inference, unification, trait solver, and read-only Core IR verifier ([`ts/`](../ts)).
  * Core IR definitions ([`bin/ir_*.h`](../bin)) and C++17 code generator ([`bin/codegen_impl.h`](../bin/codegen_impl.h)).

Communication between the Bapel frontend and C++ backend currently relies on flat-text line-oriented serialization via FFI entry points (`bapel_parser::run`, `typechecker::run`, and `codegen::compile_unit`).

---

## 2. Best Next Component to Migrate

### Recommendation: Core IR Data Model (`bapel.compiler.ir`) + C++ Code Generator (`bapel.compiler.codegen.cpp`)

The highest-leverage, lowest-risk next step is a two-part foundation replacing [`bin/ir_*.h`](../bin) and [`bin/codegen_impl.h`](../bin/codegen_impl.h):

#### Step 2.1: Define the Real Core IR in Bapel (`bapel.compiler.ir`)
* **Why first**: Every middle-end and backend pass (`normalizer`, `resolver`, `inferencer`, `typecheck`, `codegen`) operates on `IrUnit`, `IrDecl`, `IrTraitImpl`, `IrFunction`, `IrType`, `IrTerm`, and `IrKind`.
* Currently, [`bin/main.bpl`](../bin/main.bpl) only defines a shallow string-based stub of `IrUnit` for dependency discovery.
* With `SharedPtr` and `UniquePtr` now available in [`bapel/stl.bpl`](../bapel/stl.bpl), Bapel can represent recursive algebraic data types (`IrType`, `IrTerm`, `IrKind`) directly using `variant`, `struct`, `Vector`, `Optional`, and `SharedPtr`.

#### Step 2.2: Port `CppPrinter` (`bin/codegen_impl.h`) to Bapel (`bapel.compiler.codegen.cpp`)
* **Leaf pipeline stage**: Code generation sits at the very end of the compilation pipeline (`Elaborated IrUnit` $\rightarrow$ `.h`, `_private.h`, `.cc`). Parsing, normalization, and type inference can remain in C++ while a thin FFI bridge converts the typechecked C++ `ir::IrUnit` into Bapel's `IrUnit` and invokes the Bapel code generator.
* **Pure tree-walk & string emission**: Unlike `ts/inferencer.h` (which requires mutable unification state, existential variables, and interleaved method/receiver elaboration) or `cpp_parser/` (lexing and precedence-climbing), `CppPrinter` is a straightforward recursive traversal plus declaration topological sorting.
* **Authoritative test oracle**: [`tests/cpp_printer_test.cc`](../tests/cpp_printer_test.cc) already validates every golden test in `tests/testdata/cpp/*` byte-for-byte (`CppPrinterTest.GoldenFiles`) and runs `clang++ -std=c++17 -fsyntax-only` on all generated outputs (`CppPrinterTest.IsValidCpp`).

#### Alternative Option: Native Lexer & Parser (`bapel.compiler.parser`)
If you prefer to start at the front of the pipeline with zero FFI data-structure bridging:
* Porting [`cpp_parser/`](../cpp_parser) and [`ast/`](../ast) to a hand-written Recursive Descent + Pratt parser in Bapel immediately eliminates the external `antlr4-runtime` C++ library dependency from [`Makefile`](../Makefile) and generated `out/BUILD` files.
* Validated directly against `ParserTest.GoldenFiles` (`tests/testdata/parser/*.bpl`).
* However, it still requires solving the recursive ADT blockers below for `ast::Expr` and `ir::IrType`.

---

## 3. Blockers Preventing Self-Hosting in Bapel

The blockers fall into three categories: **hard compiler/codegen blockers** (cannot compile recursive compiler AST/IR types without fixing these), **language ergonomics gaps**, and **standard library gaps**.

### 3.1 Hard Compiler & Codegen Blockers

#### 1. Recursive and Mutually Recursive `type` Declarations Fail Topological Sort & C++ Emission
Compiler IR and AST types are inherently self-recursive and mutually recursive:
* `IrType` references `IrType` (via `AppType`, `FunType`, `StructField`, `VariantTag`, `TypeParam` bounds).
* `IrTerm` references both `IrTerm` (sub-expressions) and `IrType` ( `.type` annotations, `AppTypeTerm`, `InjectionTerm`).

Currently, two mechanisms block this:
* **Strict cycle rejection in `topo_sort_alias_decls`**: In [`bin/ir_unit.h`](../bin/ir_unit.h) (`collect_name_types` / `topo_sort_alias_decls`), every `NameType` referenced inside an `AliasDecl` creates a hard dependency edge—even when wrapped inside an indirection type like `SharedPtr`, `UniquePtr`, `Vector`, or `Ptr`. Any self-recursive or mutually recursive alias fails normalization/typechecking with `"cyclic dependency between type declarations"`.
* **C++ `std::variant` completeness rules**: In [`bin/codegen_impl.h`](../bin/codegen_impl.h), `type X = variant { ... }` emits `struct X : std::variant<...>`. In C++17, `std::variant` requires all alternative types to be complete types at the point of template instantiation.
* **Fix required**:
  1. Treat `NameType`s inside pointer/container indirection (`SharedPtr`, `UniquePtr`, `Ptr`, `Vector`) as weak edges during `topo_sort_alias_decls` (breaking cycles while preserving value-embedding order).
  2. Emit C++ forward declarations (`struct X;` / `template <typename a> struct X;`) in headers *before* emitting `AliasDecl` full struct/variant definitions (`print_decls` currently calls `print_type_def` directly without a prior forward-declaration pass).

#### 2. Anonymous Structs Inside `variant` Tags Are Invalid in Generated C++
* Idiomatic AST/IR definitions use inline structs for variant payloads (e.g., `type Expr = variant { app struct { fun: SharedPtr Expr, arg: SharedPtr Expr }, ... }`).
* In [`bin/codegen_impl.h`](../bin/codegen_impl.h), `record_anonymous_types_from_unit` only extracts `__anonym_<hash>` structs from `TermDecl`s and function return types—it **does not** run on `AliasDecl` bodies or function parameter types.
* Defining an unnamed `struct { ... }` inline inside `std::variant<struct { ... }>` is a C++ compile error (`types may not be defined in a type-id`).
* **Fix required**: Extend `record_anonymous_types_from_unit` to rewrite `AliasDecl` types and `IrFunction` argument types (or require naming every variant payload struct explicitly until fixed).

#### 3. String & Rune Literal Escaping and Comment-in-String Quirks
* **Unescaped literal storage**: In [`cpp_parser/ast_builder.h`](../cpp_parser/ast_builder.h) (`parse_string_literal` and `parse_rune_literal`), outer quotes are stripped without interpreting escape sequences (`\n`, `\"`, `\\`, `\t`), relying on [`bin/codegen_impl.h`](../bin/codegen_impl.h) (`ConstTerm`) emitting raw characters back between quotes (`"\"" << *lit.str_val << "\""`).
  * In Bapel code, this means `"a\nb".size` or character-by-character lexing sees literal backslash + `n` if constructed or round-tripped inconsistently, and escaping quotes inside Bapel string processing is fragile.
* **Parser comment bug workaround**: [`bin/main.bpl`](../bin/main.bpl) notes `// Safe construction of "//:" to avoid parser comment bugs` where `"//"` inside string literals interacted with comment handling or tooling.
* **Fix required**: Properly unescape string/rune literals at parse time into actual bytes in memory, and properly re-escape them when emitting C++ string/char literals in `CppPrinter`.

---

### 3.2 Language Ergonomics & Type System Blockers

While you can technically work around these, they will cause severe friction when writing 5,000+ lines of compiler passes (`IrTerm` has 16 cases; `IrType` has 11 cases):

#### 4. `match` Lacks Wildcard (`_`) Arms, Non-Variant Matching, and Exhaustiveness Checks
* In [`cpp_parser/bapel.g4`](../cpp_parser/bapel.g4), `matchArm` is strictly `id IDENTIFIER FAT_ARROW expression`.
* **No wildcard/default arm**: Writing a helper like `is_statement(term: &IrTerm) -> bool` currently requires listing all 16 variant tags of `IrTerm` instead of matching 5 tags and `_ => false`.
* **No literal matching**: Cannot `match` on `i8` bytes (in a lexer) or integers/enums without chaining `if / else if`.
* **Fix required**: Add `_ => expr` default arm support to `bapel.g4`, `IrTerm::MatchTerm`, `ts/inferencer.h`, `ts/typecheck.h`, and `CppPrinter::print_match_term` (`default: { ... }`), plus exhaustiveness checking when `_` is absent.

#### 5. Operators (`==`, `!=`, `<`) Only Support Primitives and `String`
* Trait-based operator overloading (`Eq`, `Ord`) is not yet implemented ([`designs/traits.md`](traits.md)).
* Comparing `IrKind`, `IrType`, `ModuleID`, or lexer `TokenKind` tags requires manual equality functions (`equals_type`, `equals_kind`), and custom structs/variants cannot be used directly with `Vector::sort`, `Vector::dedup`, or as keys in `UnorderedMap` without C++ `operator==` / `operator<` / `std::hash` support.

#### 6. Loops Lack `break` and `continue` (and Range/Collection Iteration)
* [`cpp_parser/bapel.g4`](../cpp_parser/bapel.g4) only provides `for <condition> { ... }` (desugared to `core::for` $\rightarrow$ C++ `while`).
* Because there is no `break`, `continue`, or `for elem in vec`, [`bin/main.bpl`](../bin/main.bpl) and [`bin/query.bpl`](../bin/query.bpl) have to write almost every traversal as an explicit index-passing tail-recursive function (`fn step(..., index: i64)`). Adding `break`/`continue` (or iterator/slice helpers) will substantially simplify lexers, parsers, and graph passes.

#### 7. Polymorphic Lambdas Not Yet Finished
* Noted as a `TODO` in [`program.bpl`](../program.bpl) and [`ideas.txt`](../ideas.txt). While top-level generic functions work well, local generic closures inside compiler passes are not yet fully supported.

---

### 3.3 Standard Library (`bapel.stl` / `bapel.core`) Gaps *(Resolved)*

The following stdlib essentials have been implemented in [`bapel/stl.bpl`](../bapel/stl.bpl) and [`bapel/core_impl.h`](../bapel/core_impl.h):

| Capability | Status | Implementation |
| :--- | :--- | :--- |
| **`String` in-place building, indexing & slicing (`append`, `push_back`, `reserve`, `at`, `substr`)** | Done | [`bapel/stl_string.h`](../bapel/stl_string.h), [`bapel/stl.bpl`](../bapel/stl.bpl) |
| **`OStringStream` (`mk`, `write`, `write_char`, `str`)** | Done | [`bapel/stl_sstream.h`](../bapel/stl_sstream.h), [`bapel/stl.bpl`](../bapel/stl.bpl) |
| **Number parsing (`parse_i64` for decimal & `0x` hex)** | Done | [`bapel/stl_string.h`](../bapel/stl_string.h) |
| **Deterministic `Map ['k, 'v]`, `Set ['a]`, and `UnorderedSet ['a]`** | Done | [`bapel/stl_map.h`](../bapel/stl_map.h), [`bapel/stl_set.h`](../bapel/stl_set.h), [`bapel/stl_unordered_set.h`](../bapel/stl_unordered_set.h), [`bapel/stl.bpl`](../bapel/stl.bpl) |
| **Zero-copy `SharedPtr::get` & `UniquePtr::get` returning `&'a`** | Done | [`bapel/stl_memory.h`](../bapel/stl_memory.h), [`bapel/stl.bpl`](../bapel/stl.bpl) |
| **Filesystem helpers (`fs::read_file`, `fs::write_file`, `fs::lexically_normal`)** | Done *(with TODO to eventually remove)* | [`bapel/stl_filesystem.h`](../bapel/stl_filesystem.h) |
| **`stderr` output (`core::eprint`)** | Done | [`bapel/core_impl.h`](../bapel/core_impl.h) |

*(Remaining optional pure-Bapel utility: deterministic 64-bit string hashing, e.g., FNV-1a / rolling hash, to be written in Bapel when porting `CppPrinter`.)*

---

## 4. Recommended Execution Sequence

To unblock self-hosting with minimal risk and keep `make test` and `make bpl` green at every commit:

1. **Expand `bapel.stl` & `bapel.core`** *(Completed)*:
   * Added `String` methods (`append`, `push_back`, `reserve`, `at`, `substr`), `parse_i64`, `OStringStream`, `Map ['k, 'v]`, `Set ['a]`, `UnorderedSet ['a]`, zero-copy `SharedPtr::get` / `UniquePtr::get`, `fs` helpers, and `core::eprint`.
2. **Unblock Recursive ADTs in Seed Compiler (C++)** *(Next)*:
   * Update `topo_sort_alias_decls` in [`bin/ir_unit.h`](../bin/ir_unit.h) and [`bin/codegen_impl.h`](../bin/codegen_impl.h) so `SharedPtr`/`UniquePtr`/`Vector`/`Ptr` break alias dependency cycles.
   * Emit struct/variant forward declarations before alias definitions in `CppPrinter`.
   * Run `record_anonymous_types` over `AliasDecl` and function parameter types (if inline structs in variants are desired).
   * Add wildcard `_ =>` arm support to `match` expressions.
   * Promote seed compiler (`make test && make bpl && make bootstrap`).
3. **Phase 1 + 2: Port `bapel.compiler.ir` & `bapel.compiler.codegen.cpp`**:
   * Define `IrKind`, `IrType`, `IrTerm`, `IrDecl`, `IrFunction`, `IrTraitImpl`, and `IrUnit` in Bapel.
   * Port `CppPrinter` to Bapel, bridged from C++ `typecheck_source_file`, and verify 100% pass on `CppPrinterTest.GoldenFiles` and `CppPrinterTest.IsValidCpp`.
4. **Phase 3: Port Read-Only Type Verifier (`bapel.compiler.types.check`)**:
   * Port [`ts/typecheck.h`](../ts/typecheck.h), [`ts/subtype.h`](../ts/subtype.h), [`ts/infer_kind.h`](../ts/infer_kind.h), and [`ts/wellformed.h`](../ts/wellformed.h) operating directly on the Bapel `IrUnit`.
5. **Phase 4: Port Desugarer, Resolver & UnitChecker (`bapel.compiler.normalizer`)**:
   * Port [`ast/ast_desugar.h`](../ast/ast_desugar.h), [`comp/desugar.h`](../comp/desugar.h), [`comp/resolver.h`](../comp/resolver.h), and [`comp/typecheck_unit.h`](../comp/typecheck_unit.h), unifying them with the existing query logic in [`bin/query.bpl`](../bin/query.bpl).
6. **Phase 5: Port Bidirectional Inferencer & Unification Solver (`bapel.compiler.types.infer`)**:
   * Port [`ts/inferencer.h`](../ts/inferencer.h), [`ts/unify.h`](../ts/unify.h), and [`ts/solver.h`](../ts/solver.h).
7. **Phase 6: Port Native Lexer & Recursive-Descent/Pratt Parser (`bapel.compiler.parser`)**:
   * Replace [`cpp_parser/`](../cpp_parser) and [`ast/`](../ast), drop `-lantlr4-runtime`, delete all legacy C++ compiler headers, and verify Stage 2 == Stage 3 fixed-point bootstrapping.
