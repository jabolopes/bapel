# Plan: Porting C++ Code Generation via C++ Implementation Header (`bin/codegen_impl.h`)

## 1. Overview & Scope

The Goal is to port the C++17 code generator ([comp/cpp_printer.go](comp/cpp_printer.go), approx. 1,300 LOC) to a native C++ implementation header (`bin/codegen_impl.h`) exposed to the `bin.main` module using Bapel annotations (`// @bpl:`).

Currently, the Bapel compiler driver ([bin/main.bpl](bin/main.bpl)) queries module dependencies natively in Bapel, but delegates C++ code generation to `bootstrap/typechecker -o`. Porting code generation to a native C++ implementation header enables the compiler driver to directly invoke code generation without external compiler wrappers or Go toolchain dependencies, completing the backend half of the self-bootstrapping pipeline.

---

## 2. Proposed Architecture & Data Structures

We create [bin/codegen.bpl](bin/codegen.bpl) as an implementation file (`implements bin.main`) of [bin/main.bpl](bin/main.bpl).

### Core Data Structures in `bin/codegen.bpl`

```bapel
implements bin.main

// Output emission targets
type CodegenMode = variant {
  public_header (),
  private_header (),
  source_file ()
}

// Code generator printer state
type CppPrinter = struct {
  mode: CodegenMode,
  indent_level: i64,
  buffer: String,
  current_module: String
}
```

---

## 3. C++17 Translation Rules & Mapping

The Bapel C++ code generator translates Bapel AST / IR representations into idiomatic C++17 conforming to the project rules:

| Bapel Construct | C++17 Target Code Pattern | Notes & User Rules |
| :--- | :--- | :--- |
| `i8`, `i16`, `i32`, `i64` | `int8_t`, `int16_t`, `int32_t`, `int64_t` | Uses `<cstdint>` standard types. |
| `bool`, `f32`, `f64` | `bool`, `float`, `double` | Standard C++ primitives. |
| `&T` / Reference | `const T&` or `T&` | Reference parameters. |
| `Ptr 'a` | `T*` | Raw pointer types. |
| `struct { x: i32, y: i32 }` | `struct TypeName { int32_t x; int32_t y; };` | Emitted with member initializers. |
| `variant { left 'a, right i32 }` | `std::variant<...>` / struct tagged union | Emitted with tag discriminators. |
| Trait Base Template | `template <typename Self> struct TraitName;` | **Declaration only** (no definition) for C++17 SFINAE correctness. |
| `impl Trait for Type` | `template <> struct TraitName<Type> { ... };` | Full template specialization in the same file. |
| `fn f['a](x: 'a) -> 'a` | `template <typename A> A f(A x) { ... }` | Generic function templates. |

---

## 4. Phased Implementation Plan

### Phase 1: Foundation & Buffer Utilities (COMPLETED)
- Implemented string formatting, indentation helpers (`cpp_indent`), and C++ keyword/namespace sanitization (`cpp_sanitize_id`) in [bin/codegen.bpl](bin/codegen.bpl).

### Phase 2: Type & Signature Formatting (COMPLETED)
- Implemented type translation routines in [bin/codegen.bpl](bin/codegen.bpl):
  - `cpp_format_type`: Translates primitives (`i8`, `i16`, `i32`, `i64`, `bool`, `f32`, `f64`, `()`) into C++ standard types (`int8_t`, `int64_t`, `void`, etc.).
  - `cpp_format_ptr_type`, `cpp_format_ref_type`, `cpp_format_array_type`: Formats pointer (`T*`), reference (`const T&`), and array (`std::array<T, N>`) types.
  - `cpp_format_function_signature`: Formats complete C++ function prototypes and template headers (`template <typename ...>`).

### Phase 3: Header Generation (`.h` and `_impl.h`) (COMPLETED)
- Implemented header emission routines in [bin/codegen.bpl](bin/codegen.bpl):
  - `cpp_emit_std_includes`: Emits `#pragma once` and standard library `#include` directives (`<array>`, `<cstdlib>`, `<functional>`, `<variant>`, `<vector>`, etc.).
  - `cpp_emit_import_includes`: Emits `#include "mod.h"` for all imported module dependencies.
  - `cpp_emit_trait_base_decl`: Emits trait base template **declarations only** (`template <typename Self, ...> struct TraitName;`) ensuring C++17 SFINAE trait resolution correctness per `.agents/AGENTS.md`.

### Phase 4: Source File Code Generation (`.cc`) (COMPLETED)
- Implemented statement and expression emission routines in [bin/codegen.bpl](bin/codegen.bpl):
  - `cpp_emit_let`: Emits let variable declarations and initializers (`formatted_type var = expr;`).
  - `cpp_emit_assign`: Emits assignment statements (`var = expr;`).
  - `cpp_emit_return`: Emits return statements (`return expr;`).
  - `cpp_emit_if_head`: Emits conditional branching (`if (cond) {`).
  - `cpp_emit_for_loop`: Emits iterative loop bounds (`for (int64_t i = start; i < end; ++i) {`).

### Phase 5: Driver Integration & Verification (COMPLETED)
- Integrated [bin/codegen.bpl](bin/codegen.bpl) into [bin/main.bpl](bin/main.bpl).
- Verified full compilation with `go test ./...`, `staticcheck`, `make all program query`, and `./bootstrap/bpl build bin.main`.

### Phase 6: Full AST/IR Emission & Traversal (COMPLETED)

Port complete AST/IR node traversal, type formatting, and C++ printer logic from [comp/cpp_printer.go](comp/cpp_printer.go) into [bin/codegen.bpl](bin/codegen.bpl) in 5 concrete sub-steps:

1. **Bapel AST / IR Data Structures (`bin.ir` Module) (COMPLETED):**
   - **Target Files:** `bin/ir.bpl`, `bin/ir_type.bpl`, `bin/ir_term.bpl`, `bin/ir_decl.bpl`, `bin/ir_function.bpl`, `bin/ir_unit.bpl`
   - Define a native Bapel IR module (`module bin.ir`) in `bin/` with types matching the Go IR definitions ([ir/ir_unit.go](ir/ir_unit.go), [ir/ir_term.go](ir/ir_term.go), [ir/ir_type.go](ir/ir_type.go), [ir/ir_decl.go](ir/ir_decl.go)):
     - `bin/ir.bpl`: Module header definition (`module bin.ir`).
     - `bin/ir_type.bpl`: `IrType` variants (`NameType`, `AppType`, `FunType`, `StructType`, `TupleType`, `VariantType`).
     - `bin/ir_term.bpl`: `IrTerm` variants (`VarTerm`, `ConstTerm`, `AppTerm`, `LetTerm`, `AssignTerm`, `BlockTerm`, `MatchTerm`, `ProjectionTerm`, `SetTerm`, `LambdaTerm`, `TupleTerm`, `StructTerm`, `InjectionTerm`).
     - `bin/ir_decl.bpl`: `IrDecl` and `IrTraitImpl`.
     - `bin/ir_function.bpl`: `IrFunction`.
     - `bin/ir_unit.bpl`: `IrUnit`.

2. **Parser & Typechecker IR Export (COMPLETED):**
   - **Target Files:** [bin/query.bpl](bin/query.bpl), `bootstrap/parser`
   - Extend `bootstrap/parser` or the compiler query interface in [bin/query.bpl](bin/query.bpl) to serialize and expose full typechecked AST/IR structures to [bin/codegen.bpl](bin/codegen.bpl).

3. **Port Core Traversal & Printing Routines (COMPLETED):**
   - **Target Files:** [bin/codegen.bpl](bin/codegen.bpl)
   - **Term Traversal (`cpp_emit_term`):** Port `PrintTerm` from [comp/cpp_printer.go](comp/cpp_printer.go) to handle recursive term printing, variable bindings, assignment destinations (`varDestination`), block expressions, and function calls.
   - **Pattern Matching (`cpp_emit_match`):** Port `printMatchTerm` to translate Bapel variant `match` expressions into C++ `std::variant::index()` switch statements and `std::get` bindings.
   - **Complex Types (`cpp_format_type` expansion):** Support `FunType` (`std::function`), `VariantType` (`std::variant`), `TupleType` (`std::tuple`/`std::monostate`), `StructType`, and anonymous struct mangling/hashing.
   - **SFINAE Trait Constraints (`cpp_emit_trait_constraints`):** Port `sfinaeConstraint` and `printTemplateParams` to generate SFINAE trait constraints (`std::enable_if_t<sizeof(...) > 0>`).
   - **Emission Modes & Namespaces:** Support `ModePublicHeader`, `ModePrivateHeader`, and `ModeSource` splitting with namespace wrapping (`inherents::`, `traits::`).

4. **Direct File I/O Integration (COMPLETED):**
   - **Target Files:** [bin/codegen.bpl](bin/codegen.bpl), [bin/main.bpl](bin/main.bpl)
   - Wire [bin/codegen.bpl](bin/codegen.bpl) routines to write generated header (`.h`, `_private.h`) and source (`.cc`) output files directly using `bapel.stl` `Ofstream`.

5. **Verification & Parity Test Harness (Phase 6.5) (COMPLETED):**
   - **Target Files:** [comp/cpp_printer_test.go](comp/cpp_printer_test.go), Bapel compiler e2e tests
   - **Objectives & Key Requirements:**
     - **Output Equivalence Verification:** For every test input in `comp/testdata/in/*.in`, generate code using both the legacy Go C++ printer (`bootstrap/compiler` / `comp.CompileBPL`) and native Bapel C++ printer ([bin/codegen.bpl](bin/codegen.bpl)), comparing emitted `.h`, `_private.h`, and `.cc` files against golden references in `comp/testdata/cpp/`.
     - **Build Output Parity Verification:** Compare the generated C++ files written directly to the `out/` directory when building complete modules (`out/bin.main build <module>`) using the native Bapel code generator against those produced by `bootstrap/compiler`.
     - **C++17 Validity Testing:** Run `clang++ -std=c++17` on all output files to guarantee that C++ produced by [bin/codegen.bpl](bin/codegen.bpl) compiles cleanly.
     - **Rule Compliance Verification:** Verify strictly per [.agents/AGENTS.md](.agents/AGENTS.md):
       - Trait Base Templates: Declared only (`template <typename Self> struct TraitName;`), never defined (ensuring SFINAE correctness).
       - SFINAE Trait Constraints: Emitted as `std::enable_if_t<(sizeof(...) > 0), int> = 0`.
       - Fully Qualified Trait Invocations: Correct namespace and scope resolution (`Size::size`).
   - **Implementation Steps:**
     1. Implement `cpp_emit_unit(unit: &IrUnit, mode: CodegenMode)` in [bin/codegen.bpl](bin/codegen.bpl) for top-level file/module code generation (`.h`, `_private.h`, `.cc`).
     2. Add `TestBapelCodegenParity` in [comp/cpp_printer_test.go](comp/cpp_printer_test.go) to compile inputs via `out/bin.main` and verify output diffs.
     3. Compare generated C++ header and source files in the `out/` directory between builds using `bootstrap/compiler` and native Bapel `bin/codegen.bpl`.
     4. Run `go test ./comp/...` to confirm parity and compilation pass.

6. **Full Function Transpilation & Self-Bootstrapping (Phase 6.6) (COMPLETED):**
   - **Target Files:** [bin/codegen.bpl](bin/codegen.bpl), `Makefile`
   - **Objectives & Key Requirements:**
     - **Full Function Codegen (`cpp_emit_function`):** Port `printFunction` from Go into [bin/codegen.bpl](bin/codegen.bpl) to transpile full `IrFunction` definitions (signatures, template parameters, return types, and body terms via `cpp_emit_term`) into `.cc` files natively in Bapel.
     - **Pre-generated Bootstrap C++ Files:** Pre-generate C++ sources (`.h`, `_private.h`, `.cc`) for `bin.main` and store them in a dedicated bootstrap directory (`bootstrap/gen/`).
     - **Native Bootstrapping Driver:** Update `Makefile` to compile pre-generated C++ files directly using `gcc`/`clang++` to produce the `bpl` executable without running `bootstrap/compiler` (Go tool).

---

## 5. Phase 7: Decouple Typechecker & Expose Fully-Annotated IR (`bootstrap/typechecker`) (COMPLETED)

Isolate the Go typechecker and inferencer ([ts/stlc](ts/stlc/)) into a standalone CLI binary (`bootstrap/typechecker`) that exports fully elaborated `IrUnit` IR objects directly to the native Bapel driver:

1. **Standalone Go Typechecker CLI (`bootstrap/typechecker`) (COMPLETED):**
   - **Target Files:** `bin/cmd/typechecker/typechecker.go`, `Makefile`
   - Created dedicated Go binary `bootstrap/typechecker` that parses source code, runs type inference & elaboration ([ts/stlc/](ts/stlc/)), and outputs the fully-annotated `IrUnit` (with concrete types, resolved method calls, auto-borrowing, and variant tags) as flat (`-format=flat`), JSON (`-format=json`), or IR (`-format=ir`) output, or compiles directly via `-o`.

2. **Driver & Query Integration (`bin/main.bpl`) (COMPLETED):**
   - **Target Files:** [bin/query.bpl](bin/query.bpl), [bin/main.bpl](bin/main.bpl), [bin/codegen.bpl](bin/codegen.bpl)
   - Updated [bin/query.bpl](bin/query.bpl) with `query_typechecked_unit` to query and retrieve annotated `IrUnit` objects from `bootstrap/typechecker`.
   - Updated [bin/main.bpl](bin/main.bpl) and [bin/codegen.bpl](bin/codegen.bpl) to use `bootstrap/typechecker` for module and implementation file compilation, completely replacing calls to `bootstrap/compiler`.

---

## 6. Phase 8: Port `comp/cpp_printer.go` to C++ via Bapel Header Annotations (PLANNED)

Port the C++17 code generator ([comp/cpp_printer.go](comp/cpp_printer.go) and [comp/cpp_printer_atypes.go](comp/cpp_printer_atypes.go)) to native C++ implementation headers structured with `ir_` prefixed files, and expose the entrypoint to [bin/main.bpl](bin/main.bpl) via Bapel annotations (`// @bpl:`):

1. **C++ IR Data Structures & JSON/IR Parser (`ir_*.h` files):**
   - **Target Files:**
     - `bin/ir.h`: Core typedefs, enums (`PrinterMode`, `Position`), namespace helpers, and forward declarations.
     - `bin/ir_type.h`: C++ `IrType` definition and type variants (`NameType`, `AppType`, `ArrayType`, `FunType`, `StructType`, `TupleType`, `VariantType`, `ForallType`, `LambdaType`).
     - `bin/ir_term.h`: C++ `IrTerm` definition and term variants (`VarTerm`, `ConstTerm`, `AppTypeTerm`, `AppTermTerm`, `LetTerm`, `AssignTerm`, `BlockTerm`, `MatchTerm`, `ProjectionTerm`, `SetTerm`, `LambdaTerm`, `TupleTerm`, `StructTerm`, `InjectionTerm`, `ReturnTerm`).
     - `bin/ir_decl.h`: `IrDecl` (`NameDecl`, `AliasDecl`, `TermDecl`, `TraitDecl`) and `IrTraitImpl`.
     - `bin/ir_function.h`: `IrFunction`, `FunctionArg`, and `TypeParam`.
     - `bin/ir_unit.h`: `IrUnit`, `IrImport`, and `IrImpl`.
     - `bin/ir_parser.h`: In-memory parser/deserializer to load the fully-annotated `IrUnit` from JSON (`bootstrap/typechecker -format=json`) or structured IR.

2. **Port Core Printer & Expression Lowering to C++ (`bin/codegen_impl.h`):**
   - **Target Files:** `bin/codegen_impl.h` (includes `bin/ir_unit.h`, etc.)
   - Port all `CppPrinter` methods from [comp/cpp_printer.go](comp/cpp_printer.go) to C++17:
     - Recursive term printing: `printType`, `printAppTypeTerm`, `printAppTermTerm`, `printLetTerm`, `printMatchTerm`, `printProjectionTerm`, `printReturnTerm`, `printTupleTerm`, `printSetTerm`, `printStructTerm`, `PrintTerm`.
     - Context state management: `varDestination`, `lastTerm`, `isCppStatement`, `withBindPosition`, `withAutoType`.
     - Symbol table lookup: `findDecl`, `findTraitDecl`, `findDeclForType`.
     - Trait & namespace formatting: `inherentCppName`, `traitCppName`, `printInNamespace`, and SFINAE constraint generation (`sfinaeConstraint`).

3. **Port Anonymous Struct Extraction & Topological Sorting (`bin/codegen_impl.h`):**
   - **Target Files:** `bin/codegen_impl.h`
   - Port the AST transformation pass from [comp/cpp_printer_atypes.go](comp/cpp_printer_atypes.go) to C++:
     - `recordAnonymousTypesFromUnit`, SHA-1 type hashing (`hashType`), `genNameType` (`__anonym_<hash>`).
     - Topological sorting of declarations (`TopoSortDecls`) in C++.

4. **Export Interface via Bapel Annotations & Integrate into Driver (`bin/main.bpl`):**
   - **Target Files:** `bin/codegen_impl.h`, [bin/main.bpl](bin/main.bpl)
   - Add Bapel annotation to export the C++ compilation entrypoint:
     ```cpp
     // @bpl: pub codegen::compile_unit: (String, String) -> i64
     inline int64_t compile_unit(const std::string& input_file, const std::string& output_base) { ... }
     ```
   - Update `bin/main.bpl` to include `codegen_impl.h` in `impls { ... }` and call `codegen::compile_unit` directly in `buildModule` and `buildImpls`, removing the subprocess call to `bootstrap/typechecker -o`.

5. **Parity Verification Against Golden Tests:**
   - **Target Files:** [comp/cpp_printer_test.go](comp/cpp_printer_test.go), `Makefile`
   - Verify that C++ emitted by `bin/codegen_impl.h` matches all golden references in `comp/testdata/cpp/` with 0 diffs.
   - Run `clang++ -std=c++17` validation across all outputs.

6. **Deprecate & Remove Go Code Generator:**
   - **Target Files:** [comp/cpp_printer.go](comp/cpp_printer.go), `comp/cpp_printer_atypes.go`, `bin/cmd/typechecker/typechecker.go`, [comp/compile.go](comp/compile.go)
   - Remove `-o` flag from `bootstrap/typechecker`.
   - Delete legacy Go C++ printer files ([comp/cpp_printer.go](comp/cpp_printer.go) and `comp/cpp_printer_atypes.go`).
