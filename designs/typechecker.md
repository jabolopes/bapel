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
| `bin/cmd/typechecker/typechecker.go` | `bin/cmd/typechecker/main.cc` | CLI driver replacing `bootstrap/typechecker`. |

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

## 5. Implementation Roadmap (Phases 1 to 6)

### Phase 1: Context, Bindings, and Functional List Data Structures (`ts/list.h`, `ts/bind.h`, `ts/context.h`)
- [ ] Implement persistent singly linked list `ts::List<T>` with `cons`, `head`, `tail`, `reverse`, `filter`, and traversal utilities.
- [ ] Implement `ts::Binding` variants (`TVarBind`, `ExistVarBind`, `SolvedExistVarBind`, `TermVarBind`, `MarkerBind`, `DeclBind`, `TraitBind`, `ImplBind`).
- [ ] Implement `ts::Context` with ordered scope manipulation:
  - `add_tvar`, `add_exist_var`, `add_solved_exist_var`, `add_term_var`, `add_marker`.
  - `drop_marker`, `split_at_marker`, `contains_var`, `lookup_term_var`, `lookup_decl`.
  - Context substitution (`substitute_exist_vars`).
- [ ] Unit tests for context manipulation and scoping operations.

### Phase 2: Kind Inference, Type Well-Formedness, Reduction & Substitution (`ts/kind.h`, `ts/wellformed.h`, `ts/reduce.h`)
- [ ] Implement `ts::infer_kind(ctx, type)` computing kinds (`*`, `* -> *`) for types and type applications.
- [ ] Implement `ts::wellformed_type(ctx, type)` validating that all free type/existential variables are declared and well-scoped.
- [ ] Implement `ts::reduce_type(type)` performing type-level beta reduction: `(fun ('a) -> T) Arg` $\to$ `T['a := Arg]`.
- [ ] Implement `ts::rename_type_vars(type)` producing deterministic variable names (`'a`, `'b`, ...).
- [ ] Unit tests matching `rename_type_vars_test.go` and `wellformed_type_test.go`.

### Phase 3: Unification, Subtyping, and Algorithmic Bidirectional Engine (`ts/unify.h`, `ts/subtype.h`, `ts/inferencer.h`)
- [ ] Implement algorithmic subtyping `ts::subtype(ctx, t1, t2)`:
  - Higher-rank polymorphic type instantiation (`instantiate_l`, `instantiate_r`).
  - Structural subtyping for tuples, structs, and variants.
  - Function subtyping with contravariant arguments and covariant return types.
- [ ] Implement bidirectional typechecker (`ts::Typechecker`):
  - `infer_term(ctx, term)` $\to$ `(type, elaborated_term, updated_ctx)`.
  - `check_term(ctx, term, expected_type)` $\to$ `(elaborated_term, updated_ctx)`.
  - `infer_app(ctx, fun_type, arg_term)` $\to$ `(ret_type, elaborated_term, updated_ctx)`.
- [ ] Implement trait resolution:
  - Match required trait bounds against available trait impls in context.
  - Validate inherent impl method signatures and dispatch.
- [ ] Implement `ts::solve_term(ctx, term)` and `ts::solve_type(ctx, type)` to substitute all solved existential variables.
- [ ] Implement `ts::check_returns(term)` verifying exhaustiveness and return paths.
- [ ] Verify test suite against all STLC test cases (`ts/stlc/*.stlc.in`).

### Phase 4: AST Desugaring, Resolver, and Querier (`comp/desugar.h`, `comp/resolver.h`, `comp/querier.h`)
- [ ] Implement `comp::desugar_expr(ast_expr)` lowering high-level AST syntax (`let`, `if`, `for`, `match`, operators) to STLC IR terms.
- [ ] Implement `comp::ModuleFinder` to resolve module IDs (`bapel.core`, `bapel.stl`) to workspace paths using `ast::Workspace`.
- [ ] Implement `comp::Querier` scanning source files and extracting `@bpl:` header declarations.
- [ ] Implement `comp::Resolver`:
  - Symbol table aggregation from imported modules and header implementations.
  - Inherent implementation namespace resolution.
  - Declaration sorting and dependency ordering.

### Phase 5: C++ `bootstrap/typechecker` CLI Binary (`bin/cmd/typechecker/main.cc`)
- [ ] Implement `comp::typecheck_source_file(querier, options, filename)` orchestrating parsing, resolution, type inference, and IR unit synthesis.
- [ ] Implement `bin/cmd/typechecker/main.cc` with CLI arguments `-format=flat|json|ir`.
- [ ] Update `Makefile` target `bootstrap/typechecker` to compile with native C++ (`clang++ -O3 -std=c++17`).

### Phase 6: End-to-End Verification, Parity Testing & Cleanup
- [ ] Run full golden tests across `comp/testdata/` and `ts/stlc/`.
- [ ] Test `bpl build bin.main`, `bpl build program`, and `bpl query` with the C++ typechecker.
- [ ] Test `make bootstrap` self-bootstrapping cycle.
- [ ] Remove deprecated Go packages (`ts/`, `comp/`, `bin/cmd/typechecker/`).
