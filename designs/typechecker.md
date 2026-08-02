# Bapel Native C++ Typechecker Design

## 1. Executive Summary & Goals

The goal of this project is to port the Bapel type system, unification engine, and typechecking pipeline from Go (`ts/` and `comp/`) to native C++17. 

By replacing the Go `bootstrap/typechecker` with a native C++ binary, the Bapel compiler achieves a major milestone in self-bootstrapping:
1. **Self-Bootstrapping**: Replaces all Go type inference, subtyping, constraint solving, and symbol resolution code with C++.
2. **Performance**: Eliminates Go runtime overhead, garbage collection pauses, and process-spawn latency during compilation.
3. **Drop-in Compatibility**: Retains 100% behavioral and format parity with the existing `bootstrap/typechecker` CLI (`-format=flat|json|ir`).
4. **Seamless Integration**: Consumes the C++ AST (`ast/*.h`) and emits C++ IR (`bin/ir_*.h`), seamlessly driving `bpl build` and `bpl query`.

---

## 2. Architecture & Component Decomposition

```
                    +--------------------------------+
                    |       C++ AST / Parser         |
                    | (ast/*.h, bootstrap/parser)    |
                    +---------------+----------------+
                                    |
                                    v
+-----------------------------------------------------------------------+
|                       comp/ (Compiler Pipeline)                       |
|                                                                       |
|  +------------------------+             +--------------------------+  |
|  |     ModuleFinder       |             |         Querier          |  |
|  |  (Workspace traversal) |             | (Header & @bpl: parsing) |  |
|  +-----------+------------+             +------------+-------------+  |
|              |                                       |                |
|              +-------------------+-------------------+                |
|                                  |                                    |
|                                  v                                    |
|                   +-------------------------------+                   |
|                   |     Resolver & Desugarer      |                   |
|                   |  - Symbol Table Construction  |                   |
|                   |  - AST Lowering to IR Terms   |                   |
|                   +---------------+---------------+                   |
+-----------------------------------|-----------------------------------+
                                    |
                                    v
+-----------------------------------------------------------------------+
|                         ts/ (Type System)                             |
|                                                                       |
|  +-------------------+  +--------------------+  +------------------+  |
|  |  Typing Context   |  |   Kind Inference   |  |  Type Reduction  |  |
|  | (Scoped Bindings) |  |  & Well-Formedness |  | (Beta-reduction) |  |
|  +---------+---------+  +---------+----------+  +--------+---------+  |
|            |                      |                      |            |
|            +----------------------+----------------------+            |
|                                   |                                   |
|                                   v                                   |
|               +---------------------------------------+               |
|               |  Bidirectional Inferencer & Subtyping |               |
|               |  - Bidirectional Check (⇐) / Synthesize (⇒)           |
|               |  - Higher-Rank Polymorphic Subtyping  |               |
|               |  - Trait & Inherent Impl Solving      |               |
|               |  - Let-Polymorphism & Generalization  |               |
|               +-------------------+-------------------+               |
|                                   |                                   |
|                                   v                                   |
|               +---------------------------------------+               |
|               |        Solver & Canonicalizer         |               |
|               |  - Existential Substitution           |               |
|               |  - Return Flow & Exhaustiveness Check |               |
|               |  - Type Variable Renaming             |               |
|               +-------------------+-------------------+               |
+-----------------------------------|-----------------------------------+
                                    |
                                    v
                    +--------------------------------+
                    |    Elaborated Typed IrUnit     |
                    | (bin/ir_unit.h, IrFunction)    |
                    +--------------------------------+
```

### 2.1 Component Mapping (Go to C++)

| Go Package / File | Target C++ Header / Source | Description |
| :--- | :--- | :--- |
| `ts/list/list.go` | `ts/list.h` | Persistent/immutable singly linked list with structural sharing for typing contexts. |
| `ts/stlc/bind.go` | `ts/bind.h` | Context bindings (`TVarBind`, `ExistVarBind`, `SolvedExistVarBind`, `TermVarBind`, `MarkerBind`, `DeclBind`, `TraitBind`, `ImplBind`). |
| `ts/stlc/context.go` | `ts/context.h` | Ordered typing context with scope manipulation, marker drops, lookup, and substitution. |
| `ts/stlc/wellformed_*.go` | `ts/wellformed.h` | Well-formedness checks for bindings, types, and contexts. |
| `ts/stlc/infer_kind.go` | `ts/infer_kind.h` | Kind inference (`*`, `* -> *`) for types and type constructors. |
| `ts/stlc/reduce_type.go` | `ts/reduce_type.h` | Type-level beta reduction for type applications over type lambdas. |
| `ts/stlc/unify.go` | `ts/unify.h` | Bidirectional unification and existential variable instantiations (`instantiateL`, `instantiateR`). |
| `ts/stlc/typechecker_subtype.go` | `ts/subtype.h` | Algorithmic subtyping for higher-rank polymorphism, tuples, structs, and variants. |
| `ts/stlc/typechecker_typecheck.go` | `ts/typecheck.h` | Bidirectional typing rules: checking (`CheckTerm`) and synthesizing (`InferTerm`, `InferApp`). |
| `ts/stlc/inferencer.go` | `ts/inferencer.h` | High-level inference engine coordinating let-generalization, trait solving, and term elaboration. |
| `ts/stlc/solve_*.go` | `ts/solver.h` | Resolving solved existential variables into terms and types. |
| `ts/stlc/returns.go` | `ts/returns.h` | Control flow return validation and terminal branch checking. |
| `ts/stlc/rename_type_vars.go` | `ts/rename_vars.h` | Deterministic canonicalization of type variable names (`'a`, `'b`, ...). |
| `comp/desugar.go` | `comp/desugar.h` | Lowering AST expressions (`ast::Expr`) into STLC IR terms (`ir::IrTerm`). |
| `comp/resolver.go` | `comp/resolver.h` | Module symbol resolution, import resolution, and inherent/trait method dispatch. |
| `comp/module_finder.go` | `comp/module_finder.h` | Workspace discovery, prefix mapping, and module ID to filesystem path mapping. |
| `comp/querier.go` | `comp/querier.h` | Module header and `@bpl:` declaration scanner for querying. |
| `comp/typecheck_source_file.go` | `comp/typecheck_file.h` | Top-level typechecking driver producing `ir::IrUnit`. |
| `bin/cmd/typechecker/typechecker.go` | `cpp_typechecker/main.cc` | CLI driver replacing `bootstrap/typechecker`. |
| `parse/parser_test.go` | `tests/parser_test.cc` | Parser golden tests (in `tests/`). |
| `comp/typecheck_source_file_test.go` | `tests/typecheck_test.cc` | Typechecker golden tests (in `tests/`). |
| `comp/cpp_printer_test.go` | `tests/cpp_printer_test.cc` | Codegen/printer golden tests (in `tests/`). |
| `ts/stlc/stlc_test.go` | `tests/stlc_test.cc` | STLC inference & typechecking golden tests (in `tests/`). |
| `ts/list/list_test.go` | `ts/list_test.cc` | Persistent list unit tests (in `ts/`). |
| `ts/stlc/rename_type_vars_test.go` | `ts/rename_vars_test.cc` | Type variable renaming unit tests (in `ts/`). |
| `ir/ir_type_test.go` | `bin/ir_type_test.cc` | IR type operations unit tests (in `bin/`). |

---

## 3. Formal Type System & Algorithms

The Bapel type system is based on **Complete and Easy Bidirectional Typechecking for Higher-Rank Polymorphism** (Dunfield & Krishnaswami), extended with:
1. **Higher-Kinded Types & Type Lambdas**: `fun ('a :: *) -> variant{none (), some 'a}` with type-level beta reduction.
2. **Trait Bounds & Trait Implementation Resolution**: `forall ['a: Size + Printable] ...` solved via context search.
3. **Inherent Implementation Blocks**: `impl String { fn size(...) -> i64 { ... } }` resolving methods under namespaces.
4. **Structural & Nominal Record Types**: Structs (`struct{x: i32, y: i32}`), Tuples (`(i32, String)`), Variants (`variant{left 'a, right i8}`).
5. **Mutable References & Pointers**: `Ptr T` / `&T` with auto-dereferencing and projection syntax.

### 3.1 Ordered Typing Context ($\Gamma$)

The typing context is an **ordered sequence of bindings**:
$$
\Gamma ::= \emptyset \mid \Gamma, \alpha \mid \Gamma, \widehat{\alpha} \mid \Gamma, \widehat{\alpha} = \tau \mid \Gamma, x : \sigma \mid \Gamma, \blacktriangleright_{\widehat{\alpha}} \mid \Gamma, \text{decl}(d) \mid \Gamma, \text{impl}(i)
$$
- $\alpha$: Universal type variable (bound by forall/lambda).
- $\widehat{\alpha}$: Unsolved existential type variable (unification variable).
- $\widehat{\alpha} = \tau$: Solved existential type variable.
- $x : \sigma$: Term variable binding.
- $\blacktriangleright_{\widehat{\alpha}}$: Marker tracking existential variable scoping during instantiation.

### 3.2 Bidirectional Typing Judgments

1. **Synthesis (Inference)**: $\Gamma \vdash e \Rightarrow \tau \dashv \Delta$
   - Given context $\Gamma$ and term $e$, synthesizes type $\tau$ and updated context $\Delta$.
2. **Checking**: $\Gamma \vdash e \Leftarrow \tau \dashv \Delta$
   - Given context $\Gamma$, term $e$, and expected type $\tau$, verifies $e$ and produces updated context $\Delta$.
3. **Subtyping**: $\Gamma \vdash \tau_1 \le \tau_2 \dashv \Delta$
   - Verifies that $\tau_1$ is a subtype of $\tau_2$, solving existential variables as necessary.
4. **Instantiation**:
   - $\Gamma \vdash \widehat{\alpha} \le \tau \dashv \Delta$ (Instantiate Left)
   - $\Gamma \vdash \tau \le \widehat{\alpha} \dashv \Delta$ (Instantiate Right)

---

## 4. CLI Protocol & Output Compatibility

The C++ binary `bootstrap/typechecker` must accept the exact CLI arguments as the Go implementation:

```bash
bootstrap/typechecker [-format=flat|json|ir] <input_file>
```

### 4.1 Output Formats

1. **`flat` (Default for compiler & `bpl build`)**:
   ```
   MODULE <module.id>
   CASE <base|impl>
   IMPORT <imported.module.id>
   IMPL <relative/impl/file.h>
   DECL <formatted_decl>
   DECL_DEF <0|1> <id> <escaped_decl>
   TRAIT_IMPL <formatted_trait_impl>
   TRAIT_DEF <trait_type> <type_name> <escaped_trait_impl>
   FUNC <formatted_fn_signature>
   FUNC_DEF <0|1> <id> <ret_type> <escaped_body_term>
   ```
2. **`json`**:
   Emits formatted JSON matching the `ir::IrUnit` serialization schema.
3. **`ir`**:
   Emits textual representation of the fully typechecked `ir::IrUnit`.

---

## 5. Implementation Roadmap & Status Summary

| Phase | Description | Key Deliverables | Status |
| :--- | :--- | :--- | :--- |
| **Phase 1** | Context, Bindings & List Structures | `ts/list.h`, `ts/bind.h`, `ts/context.h` | **COMPLETED** |
| **Phase 2** | Kind Inference, Well-Formedness & Reduction | `ts/infer_kind.h`, `ts/wellformed.h`, `ts/reduce_type.h`, `ts/rename_vars.h` | **COMPLETED** |
| **Phase 3** | Unification, Subtyping & Bidirectional Engine | `ts/unify.h`, `ts/subtype.h`, `ts/inferencer.h`, `ts/typecheck.h`, `ts/solver.h` | **COMPLETED** |
| **Phase 4** | AST Desugaring, Resolver & Querier | `ast/ast_desugar.h`, `comp/desugar.h`, `comp/resolver.h`, `comp/querier.h` | **COMPLETED** |
| **Phase 5** | C++ `bootstrap/typechecker` CLI Binary | `cpp_typechecker/main.cc`, `comp/typecheck_unit.h`, `Makefile` | **COMPLETED** |
| **Phase 6** | End-to-End Verification & Bootstrapping | `make all`, `make bootstrap`, `make program`, `./out/program`, `make query` | **COMPLETED** |
| **Phase 7** | C++ Golden Tests & Native Test Runner | `tests/test_main.cc`, `tests/typecheck_test.cc`, `tests/stlc_test.cc`, `Makefile` | **PLANNED** |

---

### Phase 1: Context, Bindings, and Functional List Data Structures (`ts/list.h`, `ts/bind.h`, `ts/context.h`) — [COMPLETED]
- [x] Implement persistent singly linked list `ts::List<T>` with `cons`, `head`, `tail`, `reverse`, `filter`, and traversal utilities.
- [x] Implement `ts::Binding` variants (`TVarBind`, `ExistVarBind`, `SolvedExistVarBind`, `TermVarBind`, `MarkerBind`, `DeclBind`, `TraitBind`, `ImplBind`).
- [x] Implement `ts::Context` with ordered scope manipulation:
  - `add_tvar`, `add_exist_var`, `add_solved_exist_var`, `add_term_var`, `add_marker`.
  - `drop_marker`, `split_at_marker`, `contains_var`, `lookup_term_var`, `lookup_decl`.
  - Context substitution (`substitute_exist_vars`).
- [x] Verified context manipulation, scoping operations, and structural sharing.

### Phase 2: Kind Inference, Type Well-Formedness, Reduction & Substitution (`ts/infer_kind.h`, `ts/wellformed.h`, `ts/reduce_type.h`, `ts/rename_vars.h`) — [COMPLETED]
- [x] Implement `ts::infer_kind(ctx, type)` computing kinds (`*`, `* -> *`) for types and type applications.
- [x] Implement `ts::wellformed_type(ctx, type)` validating that all free type/existential variables are declared and well-scoped.
- [x] Implement `ts::reduce_type(type)` performing type-level beta reduction: `(fun ('a) -> T) Arg` $\to$ `T['a := Arg]`.
- [x] Implement `ts::rename_type_vars(type)` producing deterministic canonical variable names (`'a`, `'b`, ...).
- [x] Verified against all kind inference and type reduction test suites.

### Phase 3: Unification, Subtyping, and Algorithmic Bidirectional Engine (`ts/unify.h`, `ts/subtype.h`, `ts/inferencer.h`) — [COMPLETED]
- [x] Implement algorithmic subtyping `ts::subtype(ctx, t1, t2)`:
  - Higher-rank polymorphic type instantiation (`instantiate_l`, `instantiate_r`).
  - Structural subtyping for tuples, structs, and variants.
  - Higher-kinded `LambdaType` alpha-equivalence and subtyping.
  - Function subtyping with contravariant arguments and covariant return types.
- [x] Implement bidirectional typechecker (`ts::Typechecker` and `ts::Inferencer`):
  - `infer_term(ctx, term)` $\to$ `(type, elaborated_term, updated_ctx)`.
  - `check_term(ctx, term, expected_type)` $\to$ `(elaborated_term, updated_ctx)`.
  - `infer_app(ctx, fun_type, arg_term)` $\to$ `(ret_type, elaborated_term, updated_ctx)`.
- [x] Implement trait resolution:
  - Match required trait bounds against available trait impls in context.
  - Validate inherent impl method signatures and dispatch.
- [x] Implement `ts::solve_term(ctx, term)` and `ts::solve_type(ctx, type)` to substitute all solved existential variables.
- [x] Implement `ts::check_returns(term)` verifying exhaustiveness and return paths.
- [x] Verified against all STLC test cases (`ts/stlc/*.stlc.in`).

### Phase 4: AST Desugaring, Resolver, and Querier (`ast/ast_desugar.h`, `comp/resolver.h`, `comp/querier.h`) — [COMPLETED]
- [x] Implement `ast::desugar_expr(ast_expr)` lowering high-level AST syntax (`let`, `if`, `for`, `match`, operators, variants, structs, projections) to STLC IR terms.
- [x] Implement AST JSON deserialization (`ast::deserialize_ast_source_file`, `ast::deserialize_ast_function`, `ast::deserialize_ast_impl`).
- [x] Implement `comp::Querier` scanning source files and extracting `@bpl:` header declarations with subprocess escaping.
- [x] Implement `comp::Resolver`:
  - Symbol table aggregation from imported modules and header implementations.
  - Inherent implementation namespace resolution.
  - Declaration topological sorting and dependency ordering.
  - Cycle detection for recursive module imports.

### Phase 5: C++ `bootstrap/typechecker` CLI Binary (`cpp_typechecker/main.cc`, `comp/typecheck_unit.h`) — [COMPLETED]
- [x] Implement `comp::typecheck_source_file(querier, options, filename)` orchestrating parsing, resolution, type inference, and IR unit synthesis.
- [x] Implement `comp::new_default_context()` with standard primitives (`bool`, `i8`, `i16`, `i32`, `i64`, `f32`, `f64`, `void`), operators, and built-in intrinsics.
- [x] Implement `cpp_typechecker/main.cc` with CLI arguments `-format=flat|json|ir`, `-v`, and `-quiet`.
- [x] Update `Makefile` target `bootstrap/typechecker` to compile with native C++ (`clang++ -O3 -std=c++17`).

### Phase 6: End-to-End Verification, Parity Testing & Cleanup — [COMPLETED]
- [x] Verified full golden test parity across `comp/testdata/` and `ts/stlc/`.
- [x] Tested `bpl build bin.main`, `bpl build program`, and `bpl query` with the C++ typechecker.
- [x] Verified `make bootstrap` self-bootstrapping cycle.
- [x] Verified full execution of `program.bpl` and all implementation modules (`program_array.bpl`, `program_conditionals.bpl`, `program_point.bpl`, `program_string.bpl`, `program_unordered_map.bpl`, `program_variant.bpl`, `program_vector.bpl`).

---

### Phase 7: Comprehensive C++ Golden Tests & Native Test Runner Suite — [PLANNED]

All golden tests, golden testdata, and test runner utilities reside in the top-level `tests/` directory (`tests/testdata/`), while component unit tests remain within their respective component packages (`ts/`, `bin/`):

```
tests/
├── test_main.cc           # CLI test runner entry point
├── test_util.h            # Test registration macros, assertions, diffing, -regen, and globbing
├── parser_test.cc         # Golden tests for C++ parser
├── typecheck_test.cc      # Golden tests for C++ typechecker
├── cpp_printer_test.cc    # Golden tests for C++ codegen printer
├── stlc_test.cc           # Golden tests for STLC inference & typechecking
└── testdata/
    ├── parse/
    │   ├── in/            # Parser test inputs (*.in)
    │   └── parsed/        # Parser AST golden outputs (*.bpl)
    ├── comp/
    │   ├── in/            # Compiler/typechecker test inputs (*.in)
    │   ├── typecheck/     # Typecheck IR golden outputs (*.out)
    │   └── cpp/           # C++ codegen golden files (*.h, *_private.h, *.cc)
    └── stlc/
        ├── *.stlc.in      # STLC input terms
        └── *.stlc.out     # STLC inferred golden outputs

ts/
├── list_test.cc           # Unit tests for functional persistent singly-linked list
└── rename_vars_test.cc    # Unit tests for type variable alpha-renaming

bin/
└── ir_type_test.cc        # Unit tests for IR type operations and formatting
```

- [x] **Phase 7.1: Testdata Migration to `tests/testdata/`** — [COMPLETED]:
  - Move parser testdata from `parse/testdata/` to `tests/testdata/parse/`.
  - Move compiler/typechecker and printer testdata from `comp/testdata/` to `tests/testdata/comp/`.
  - Move STLC test files from `ts/stlc/*.stlc.*` to `tests/testdata/stlc/`.
  - *Migration Note*: Ensure Go tests or `Makefile` target dependencies are updated concurrently so that `make bpl` (which runs `go test ./...`) continues to pass until the native C++ test suite replaces it.

- [x] **Phase 7.2: Native C++ Test Harness (`tests/test_main.cc`, `tests/test_util.h`)** — [COMPLETED]:
  - Implement lightweight, self-contained C++ test framework with test registration macros (`TEST(Suite, Case)`), test filters (`--filter=...`), and failure assertions.
  - Implement filesystem traversal and globbing utilities (`tests::glob("tests/testdata/.../*.in")`).
  - Implement golden comparison and line-by-line diff formatting (`tests::diff(got, want)`).
  - Support `-regen` CLI flag to automatically regenerate and update golden `.out`, `.bpl`, and `.cc` files in `tests/testdata/` when intentional AST/IR/codegen changes occur.

- [x] **Phase 7.3: Parser Golden Tests (`tests/parser_test.cc`)** — [COMPLETED]:
  - Port `parse/parser_test.go` to native C++.
  - Parse test inputs in `tests/testdata/parse/in/*.in` by executing `bootstrap/parser --symbol=SourceFile --format=bpl --with-pos` (or formatting AST via `ast::SourceFile::to_string()`).
  - Validate formatted AST output against golden files in `tests/testdata/parse/parsed/*.bpl`.
  - Verify syntax and parse error diagnostics for negative test cases (`bad_*.in`).

- [ ] **Phase 7.4: Typechecker Golden Tests (`tests/typecheck_test.cc`)**:
  - Port `comp/typecheck_source_file_test.go` to native C++.
  - Run `comp::typecheck_source_file` across all test files in `tests/testdata/comp/in/*.in`.
  - Configure `comp::TypecheckOptions`: enable `options.skip_undefined_term_checks = true` for `order.in`.
  - Validate formatted `IrUnit` / typecheck output against golden files in `tests/testdata/comp/typecheck/*.out`.
  - Verify error messages for typechecking failure cases.

- [ ] **Phase 7.5: C++ Codegen & Printer Golden Tests (`tests/cpp_printer_test.cc`)**:
  - Port `comp/cpp_printer_test.go` to native C++.
  - Run `codegen::compile_unit` (from `bin/codegen_impl.h`) across all inputs in `tests/testdata/comp/in/*.in` (skipping `order.in`).
  - Compare generated C++ outputs (header `.h`, private header `_private.h`, and translation unit `.cc`) against golden references in `tests/testdata/comp/cpp/`.
  - `TestCppPrinterIsValidCpp`: Compile generated `.cc` files with `clang++ -std=c++17 -c`, skipping files that import `bapel.core` (`array.cc`, `context1.cc`, `loops.cc`, `polymorphism.cc`).

- [ ] **Phase 7.6: STLC Inference & Typechecking Golden Tests (`tests/stlc_test.cc`)**:
  - Port `ts/stlc/stlc_test.go` to native C++.
  - `TestInferTerm`: Run `comp::typecheck_source_file` across all test inputs in `tests/testdata/stlc/*.stlc.in` with `TypecheckOptions{skip_default_context = true, skip_term_typechecker = true, skip_undefined_term_checks = true}` and compare formatted `IrUnit` against golden `tests/testdata/stlc/*.stlc.out`.
  - `TestTypecheckTerm`: Run `comp::typecheck_source_file` with `TypecheckOptions{skip_default_context = true, skip_term_typechecker = false, skip_undefined_term_checks = true}` verifying exhaustive type checking across all STLC test cases.

- [ ] **Phase 7.7: Component Unit Tests (`ts/list_test.cc`, `ts/rename_vars_test.cc`, `bin/ir_type_test.cc`)**:
  - Refactor existing phase tests (`ts/phase[1-4]_test.cc`) and port Go unit tests:
    - **`ts/list_test.cc`**: Port `ts/list/list_test.go` testing persistent singly-linked list immutability, structural sharing, iteration, filtering, and reversal.
    - **`ts/rename_vars_test.cc`**: Port `ts/stlc/rename_type_vars_test.go` testing canonical type variable naming and alpha-renaming.
    - **`bin/ir_type_test.cc`**: Port `ir/ir_type_test.go` testing type representations, subtyping, and string formatting.

- [ ] **Phase 7.8: Makefile & CI Integration**:
  - Add `bootstrap/test_runner` build target in `Makefile` compiling all test suites with `clang++ -O3 -std=c++17`.
  - Add `make test-cpp` target to execute all native C++ tests.
  - Update `make test` and `make regen` to execute native C++ test suites directly without Go toolchain dependencies.
  - Update `bpl:` target in `Makefile` to run the native C++ test suite.
