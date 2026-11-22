# Self-Hosting Migration Analysis: Next Component & Blockers

This document complements [`designs/boot.md`](boot.md) with a concrete technical assessment of the current hybrid compiler codebase, recommending the immediate next component to migrate from C++ to Bapel and cataloging all language, compiler, and standard library blockers that stand in the way of full self-hosting.

---

## 1. Current State Summary

Currently, the Bapel compiler (`bpl`) is split across two layers:

* **Written in Bapel (`bin/`, `bapel/`)**:
  * CLI driver, Bazel workspace/BUILD generator, and module build graph ([`bin/main.bpl`](../bin/main.bpl)).
  * Module/file query engine and `workspace.bpl` package resolution ([`bin/query.bpl`](../bin/query.bpl)).
  * Core IR data model (`IrKind`, `IrType`, `IrTerm`, `IrDecl`, `IrFunction`, `IrTraitImpl`, `IrFile`, `IrModule`, `IrTypedModule`) in [`bapel/compiler/ir.bpl`](../bapel/compiler/ir.bpl) and its `ir_*.bpl` implementation files.
  * Standard library wrappers (`String`, `StringView`, `Vector`, `Deque`, `Optional`, `Map`, `Set`, `UnorderedMap`, `UnorderedSet`, `SharedPtr`, `UniquePtr`, `Ifstream`, `Ofstream`, `IStringStream`, `OStringStream`, `fs`, `os`, `args`).
* **Written in C++ (`ast/`, `comp/`, `ts/`, `cpp_parser/`, `bin/ir_*.h`, `bin/codegen_impl.h`)**:
  * ANTLR4 lexer/parser and AST data structures ([`cpp_parser/`](../cpp_parser), [`ast/`](../ast)).
  * Whole-module desugaring, recursive-type rewriting, import/impl resolution, `Session` caching, and `ModuleChecker` ([`comp/`](../comp)).
  * Bidirectional type inference, unification, trait solver, and read-only Core IR verifier ([`ts/`](../ts)).
  * Core IR definitions ([`bin/ir_*.h`](../bin)) and C++17 code generator ([`bin/codegen_impl.h`](../bin/codegen_impl.h)).

Communication between the Bapel frontend and C++ backend currently relies on `comp::Session` + FFI entry points (`typechecker::run` and `codegen::compile_module`).

---

## 2. Best Next Component to Migrate

### Recommendation: Core IR Data Model (`bapel.compiler.ir`) + C++ Code Generator (`bapel.compiler.codegen.cpp`)

The highest-leverage, lowest-risk next step is a two-part foundation replacing [`bin/ir_*.h`](../bin) and [`bin/codegen_impl.h`](../bin/codegen_impl.h):

#### Step 2.1: Define the Real Core IR in Bapel (`bapel.compiler.ir`) *(Completed)*
* **Why first**: Every middle-end and backend pass (`normalizer`, `resolver`, `inferencer`, `typecheck`, `codegen`) operates on `IrModule`, `IrTypedModule`, `IrFile`, `IrDecl`, `IrTraitImpl`, `IrFunction`, `IrType`, `IrTerm`, and `IrKind`.
* Implemented in [`bapel/compiler/ir.bpl`](../bapel/compiler/ir.bpl) and `ir_base.bpl`, `ir_kind.bpl`, `ir_type.bpl`, `ir_term.bpl`, `ir_decl.bpl`, using `variant`, `struct`, `Vector`, `Optional`, and `SharedPtr`, and verified to compile cleanly as a whole module with `clang++ -std=c++17 -fsyntax-only`.

#### Step 2.2: Port `CppPrinter` (`bin/codegen_impl.h`) to Bapel (`bapel.compiler.codegen.cpp`) *(Next)*
* **Leaf pipeline stage**: Code generation sits at the very end of the compilation pipeline (`IrTypedModule` $\rightarrow$ `.h`, `_private.h`, `.cc`). Parsing, normalization, and type inference can remain in C++ while a thin FFI bridge converts the typechecked C++ `ir::IrTypedModule` into Bapel's `IrTypedModule` and invokes the Bapel code generator.
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

#### 1. Recursive `type` Declarations *(Resolved for Self-Recursive Types)*
* Self-recursive types (including regular generic self-recursive types like `Tree ['a]`) are now implemented via `MuType` in Core IR (`bin/ir_base.h`, `bin/ir_type.h`, `bin/ir_substitute_type.h`), self-recursive desugaring (`comp/rewrite_recursive_types.h`), coinductive subtyping/unification (`ts/subtype.h`, `ts/unify.h`), and targeted forward declarations in `CppPrinter` (`bin/codegen_impl.h`). See [`designs/recursive-types.md`](recursive-types.md).
* Mutually recursive types (SCC & Bekić's theorem) are deferred to [`designs/features.md`](features.md) — not a hard blocker since helper structs can be inlined or parameterized over the recursive type.

#### 2. Anonymous Structs Inside `variant` Tags & Function Parameters *(Resolved)*
* `record_anonymous_types_from_module` in [`bin/codegen_impl.h`](../bin/codegen_impl.h) now lifts anonymous `StructType`s from `AliasDecl` bodies (including `variant` arms and nested `struct` fields), `TermDecl`s, `TraitDecl`s, `IrFunction` args/returns, and `IrTraitImpl` methods, handling `MuType` `'self` substitution, generic free `VarType` parameterization, `topo_sort_alias_decls` ordering, and `InjectionTerm` struct construction. Verified by `tests/testdata/in/anonymous_structs.in` and `tests/testdata/in/recursive_types.in`.

#### 3. String & Rune Literal Escaping and Comment-in-String Quirks
* **Unescaped literal storage**: In [`cpp_parser/ast_builder.h`](../cpp_parser/ast_builder.h) (`parse_string_literal` and `parse_rune_literal`), outer quotes are stripped without interpreting escape sequences (`\n`, `\"`, `\\`, `\t`), relying on [`bin/codegen_impl.h`](../bin/codegen_impl.h) (`ConstTerm`) emitting raw characters back between quotes (`"\"" << *lit.str_val << "\""`).
  * In Bapel code, this means `"a\nb".size` or character-by-character lexing sees literal backslash + `n` if constructed or round-tripped inconsistently, and escaping quotes inside Bapel string processing is fragile.
* **Parser comment bug workaround**: [`bin/main.bpl`](../bin/main.bpl) notes `// Safe construction of "//:" to avoid parser comment bugs` where `"//"` inside string literals interacted with comment handling or tooling.
* **Fix required**: Properly unescape string/rune literals at parse time into actual bytes in memory, and properly re-escape them when emitting C++ string/char literals in `CppPrinter`.

---

### 3.2 Language Ergonomics & Type System Blockers

While you can technically work around these, they will cause severe friction when writing 5,000+ lines of compiler passes (`IrTerm` has 16 cases; `IrType` has 11 cases):

#### 4. `match` Wildcard (`_`) Arms & Exhaustiveness Checks *(Resolved for Variant Matching)*
* Wildcard `_ =>` arms, duplicate-tag rejection, unreachable-arm rejection, exhaustiveness checking, and expected type propagation are now implemented across `cpp_parser/bapel.g4`, `cpp_parser/ast_builder.h`, `ast/ast_expr.h`, `bin/ir_term.h`, `ts/inferencer.h`, `ts/typecheck.h`, and `bin/codegen_impl.h`. See [`designs/pattern-matching.md`](pattern-matching.md).
* Non-variant literal matching and nested destructuring are deferred to [`designs/features.md`](features.md).

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
2. **Unblock Recursive ADTs in Seed Compiler (C++)** *(Completed)*:
   * ~~Implement self-recursive types (`MuType`) and targeted forward declarations in `CppPrinter`~~ *(Completed — see [`designs/recursive-types.md`](recursive-types.md))*.
   * ~~Lift anonymous structs from `AliasDecl` bodies (`variant` tags / nested structs) and function parameters~~ *(Completed)*.
   * ~~Add wildcard `_ =>` arm support to `match` expressions~~ *(Completed — see [`designs/pattern-matching.md`](pattern-matching.md))*.
   * ~~Support inherent `impl` blocks on `AliasDecl` types (non-recursive & recursive)~~ *(Completed)*.
   * Promote seed compiler (`make test && make bpl && make bootstrap`).
3. **Phase 1 + 2: Port `bapel.compiler.ir` & `bapel.compiler.codegen.cpp`** *(In Progress)*:
   * ~~Define `IrKind`, `IrType`, `IrTerm`, `IrDecl`, `IrFunction`, `IrTraitImpl`, `IrFile`, `IrModule`, and `IrTypedModule` in Bapel (`bapel/compiler/ir*.bpl`)~~ *(Completed — compiles as a whole module & passes `clang++ -std=c++17 -fsyntax-only` with 0 errors)*.
   * ~~Implement Whole-Module Compilation (`IrModule` $\rightarrow$ `ModuleChecker` $\rightarrow$ `IrTypedModule` $\rightarrow$ C++17, equal-weight base and `.bpl` `impls` files, per-file import scoping, signature-only `ImplDecl`, and `comp::Session` caching)~~ *(Completed — see [`designs/modules.md`](modules.md))*.
   * **Next:**
     1. Bump C++ codegen to `-std=c++20` with `bapel::Variant<bapel::Arm<"tag", T>...>` (`bapel/core_variant.h`), rewrite the legacy C++ IR (`bin/ir_*.h`) to match, and swap in the Bapel-generated `bapel/compiler/ir.h` directly (retiring `bin/ir_*.h` with no conversion bridge — see [`designs/variants.md`](variants.md)).
     2. Port `CppPrinter` to `bapel.compiler.codegen.cpp` operating directly on `bapel.compiler.ir`, and verify 100% pass on `CppPrinterTest.GoldenFiles` and `CppPrinterTest.IsValidCpp`.
4. **Phase 3: Port Read-Only Type Verifier (`bapel.compiler.types.check`)**:
   * Port [`ts/typecheck.h`](../ts/typecheck.h), [`ts/subtype.h`](../ts/subtype.h), [`ts/infer_kind.h`](../ts/infer_kind.h), and [`ts/wellformed.h`](../ts/wellformed.h) operating directly on the Bapel `IrModule` / `IrTypedModule`.
5. **Phase 4: Port Desugarer, Resolver & ModuleChecker (`bapel.compiler.normalizer`)**:
   * Port [`comp/desugar.h`](../comp/desugar.h), [`comp/rewrite_recursive_types.h`](../comp/rewrite_recursive_types.h), [`comp/resolver.h`](../comp/resolver.h), and [`comp/typecheck_module.h`](../comp/typecheck_module.h), unifying them with the existing query logic in [`bin/query.bpl`](../bin/query.bpl).
6. **Phase 5: Port Bidirectional Inferencer & Unification Solver (`bapel.compiler.types.infer`)**:
   * Port [`ts/inferencer.h`](../ts/inferencer.h), [`ts/unify.h`](../ts/unify.h), and [`ts/solver.h`](../ts/solver.h).
7. **Phase 6: Port Native Lexer & Recursive-Descent/Pratt Parser (`bapel.compiler.parser`)**:
   * Replace [`cpp_parser/`](../cpp_parser) and [`ast/`](../ast), drop `-lantlr4-runtime`, delete all legacy C++ compiler headers, and verify Stage 2 == Stage 3 fixed-point bootstrapping.
