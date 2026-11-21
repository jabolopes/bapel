# Bapel Type System Architecture: Elaboration to an Explicitly-Typed Core

This document defines the architectural separation of concerns between Bapel's **Type Inferencer** and **Typechecker**, as well as the core invariants of Bapel's nominal/structural type system.

In Bapel, type checking follows the industry-standard **Elaboration to an Explicitly-Typed Core** architecture (analogous to GHC/Haskell's Core, Rust's MIR/rustc, OCaml, and Scala). This architecture cleanly separates complex, heuristic-heavy inference and desugaring from simple, deterministic type verification.

---

## 1. The Two-Stage Pipeline

### Stage 1: The Type Inferencer (The Elaborator)
*   **Location:** [ts/inferencer.h](../ts/inferencer.h) (invoked via `Typechecker::infer_function` in [ts/typecheck.h](../ts/typecheck.h)).
*   **Role:** The inferencer is responsible for all complex, heuristic-heavy type analysis and program transformations:
    *   Unification and constraint solving ([ts/unify.h](../ts/unify.h), [ts/solver.h](../ts/solver.h)).
    *   Managing existential variables (`evars`).
    *   Method resolution and trait lookup (`Context::lookup_method` in [ts/context.h](../ts/context.h)).
    *   Automatic receiver adjustment (auto-borrowing via `Ptr::mk` and auto-dereferencing via `Ptr::get`).
    *   Inferring omitted generic type arguments.
*   **Output:** The inferencer outputs a **Fully Annotated Core IR**. By the end of inference, every single expression, sub-expression, variable, and let-binding has an explicit, concrete `.type` assigned, and all implicit behaviors or syntactic sugar are explicitly desugared.

### Stage 2: The Typechecker (The Core Verifier / Trusted Kernel)
*   **Location:** [ts/typecheck.h](../ts/typecheck.h) (`Typechecker::typecheck_function`).
*   **Role:** The typechecker acts as a lightweight, fast, deterministic **Trusted Kernel**. Because it expects the incoming IR from the inferencer to be fully annotated and desugared, the typechecker requires **zero unification, zero constraint solving, zero existential variables, and zero heuristics**.
*   **Verification:** It performs simple, mechanical bottom-up or top-down validation:
    *   In a function application `f x`, it verifies that `f.type` is `A -> B` and `x.type` is a subtype of `A`, confirming the application has type `B`.
    *   In a let-binding `let x: T = val`, it verifies that `val.type` is a subtype of `T`.
*   **Safety Guarantee:** By keeping the verifier simple and small, it acts as an independent safety check. Even if the complex inferencer has a bug or makes an error during constraint solving, this simple verifier will immediately catch ill-typed IR before C++ code generation in [bin/codegen_impl.h](../bin/codegen_impl.h).

---

## 2. Desugaring Architecture: Syntactic vs. Type-Directed

A common question in compiler design is whether all desugaring should be separated into its own distinct pass before or after type inference. In Bapel, desugaring is bifurcated based on whether it requires type information:

### Syntactic Desugaring $\rightarrow$ Pre-Inference Pass
*   **Examples:**
    *   **`for` loops:** Rewriting `for cond { body }` into `core::for(cond, |_: ()| body)`.
    *   **`Self` in `impl` blocks (`desugar_impl`):** Replacing `NameType("Self")` with the concrete `impl.type_name` across all method type parameter bounds, argument types, return types, and body terms (via `ir::rewrite_types_in_term` in [bin/ir_visit_term.h](../bin/ir_visit_term.h)).
    *   **`Self` in `trait` declarations (`desugar_trait`):** Converting `NameType("Self")` to `VarType("Self")` (`'Self`) so `forall ['Self]` binds it like any ordinary type variable.
*   **Location:** Handled in a separate pass before normalization and inference in [comp/desugar.h](../comp/desugar.h).
*   **Rationale:** These transformations are purely structural and depend only on syntax. Performing them before inference keeps the IR small and clean, meaning downstream stages ([ts/context.h](../ts/context.h), [ts/wellformed.h](../ts/wellformed.h), [comp/typecheck_module.h](../comp/typecheck_module.h), [bin/codegen_impl.h](../bin/codegen_impl.h)) never have to handle `NameType("Self")` or maintain a `Self` alias binding.

### Type-Directed Desugaring (Elaboration) $\rightarrow$ Interleaved in Inferencer
*   **Examples:**
    *   **Method Resolution:** Rewriting dot-method calls (`s.size()`) to fully qualified calls (`Size::size(s)` or `Point::size(s)`).
    *   **Receiver Adjustment:** Wrapping receivers in address-of (`Ptr::mk s`) or dereferencing (`Ptr::get s`).
    *   **Shadowing Resolution:** Determining whether `s.foo` accesses a struct field `foo` or invokes a method `foo()`.
*   **Location:** Interleaved directly inside the inferencer ([ts/inferencer.h](../ts/inferencer.h)).
*   **Why it cannot be separated into a distinct pass:**
    *   **Cannot be done BEFORE inference:** The receiver's type is unknown before inference, making it impossible to know whether `s.size()` calls `Point::size`, `Size::size`, or accesses a struct field.
    *   **Cannot be done AFTER inference:** In an expression like `let x = s.size() + 5`, the inferencer **must** resolve `s.size()` immediately to know that its return type is `i64` in order to typecheck the addition and infer the type of `x`. Postponing method resolution to a post-inference pass would cause type inference to stall or fail.

---

## 3. Type Alias Preservation & Nominal Traits

1.  **Preserve Type Aliases (Upcast, Never Lower):**
    *   Bapel's type system preserves nominal type aliases (`Point`, `Node`, etc.) wherever possible rather than lowering them to their underlying structural `struct{...}` / `variant{...}` types.
    *   Structural literals can be **upcast** into a named alias (e.g., `let p: Point = struct{ x = 1, y = 2 }`), while named aliases are kept intact through unification, inference, and verification.
2.  **Nominal Traits and Methods:**
    *   `TraitImplBind` in [ts/context.h](../ts/context.h) stores unreduced `type_name` and `trait_type`, and both inherent methods (`Type::method`) and trait bounds (`satisfies_bound` in [ts/subtype.h](../ts/subtype.h)) match **nominally**.
    *   An anonymous `struct{...}` or `variant{...}` without a nominal type does not inherit methods or trait implementations from a structurally identical alias (see [tests/testdata/in/bad_traits_anonymous_struct.in](../tests/testdata/in/bad_traits_anonymous_struct.in)).
3.  **Equirecursive `MuType` Alias Recovery:**
    *   Recursive type aliases (e.g., `type Node = struct{ val: i64, next: &Node }`) are rewritten to `mu 'self. struct{...}` in [comp/rewrite_recursive_types.h](../comp/rewrite_recursive_types.h).
    *   Unfolding a recursive alias when projecting a self-referential field (e.g., `(*n).next`) yields a `MuType` rather than `NameType("Node")`.
    *   `Context::lookup_mu_alias_name` ([ts/context.h](../ts/context.h)) resolves a `MuType` back to its defining alias name via alpha-equivalent `match_type`, so nominal method lookup and trait-bound checking work seamlessly on recursive types.

---

## 4. Whole-Module Orchestration (`ModuleChecker`)

At the module level, `ModuleChecker::check_module` ([comp/typecheck_module.h](../comp/typecheck_module.h)) coordinates checking across all files in an `IrModule` and produces a flat `IrTypedModule`:

1.  **Topological Sort:** Merges and topologically sorts module-local declarations (`ir::topo_sort_decls`).
2.  **Per-File Contexts & Well-Formedness:** Builds an isolated `ts::Context` per `.bpl` file containing its reachable imports, module-local declarations, and `ImplDecl` bindings, then validates every declaration with `ts::is_wellformed_decl` ([ts/wellformed.h](../ts/wellformed.h)) and kind inference ([ts/infer_kind.h](../ts/infer_kind.h)).
3.  **Function & Impl Checking:** Runs `infer_function` $\to$ `typecheck_function` on every function and `impl` method. For `impl` blocks, enforces trait coherence (orphan rule: only module-local types may implement traits) and verifies each trait impl method's signature against its `TraitDecl`.
4.  **Undefined-Symbol Check:** Verifies that every non-FFI `TermDecl` in the module has a matching definition.

---

## 5. Requirements for Clean Separation

To maintain the simplicity of the typechecker and prevent logic duplication, the inferencer enforces four core invariants before handing the IR over to the typechecker:

1.  **100% AST/IR Annotation & Read-Only Verification:** Every single `IrTerm` node (including intermediate sub-expressions in applications, projections, and tuples) must have its `.type` field explicitly populated by the end of inference. When Stage 2 (`Typechecker`) receives the IR, it treats the AST as strictly **read-only**, verifying existing annotations without mutating or re-assigning types.
2.  **Existential Resolution ("Zonking"):** During inference, types contain temporary unification variables (`evars`). At the very end of `infer_function`, a final substitution pass ([ts/solver.h](../ts/solver.h)) replaces all unresolved `evars` in `.type` annotations with their solved concrete types.
3.  **Explicit Desugaring & Elaboration:** All implicit behaviors and code-generation metadata must be explicitly written into the IR during Stage 1 inference and zonking:
    *   **Method Calls:** Rewritten to explicit fully qualified function applications (`Type::method(adjusted_s, args)`).
    *   **Generic Instantiations:** When calling a polymorphic function `f` without explicit type arguments, the inferencer explicitly wraps it in an `AppTypeTerm` (`f [T]`), ensuring the simple typechecker never has to guess generic type parameters.
    *   **Variant and Match Metadata:** Unconditionally populates `tag_index` on variant injections, `index` on match arms, and `reduced_type` on projections and set terms.
4.  **Minimal Kernel:** Duplicated type synthesis and reduction logic has been stripped out of [ts/typecheck.h](../ts/typecheck.h), leaving a minimal, trusted verifier.
