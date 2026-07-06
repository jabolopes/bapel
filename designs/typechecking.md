# Bapel Type System Architecture: Elaboration to an Explicitly-Typed Core

This document defines the architectural separation of concerns between Bapel's **Type Inferencer** and **Typechecker**.

In Bapel, type checking follows the industry-standard **Elaboration to an Explicitly-Typed Core** architecture (analogous to GHC/Haskell's Core, Rust's MIR/rustc, OCaml, and Scala). This architecture cleanly separates complex, heuristic-heavy inference and desugaring from simple, deterministic type verification.

---

## 1. The Two-Stage Pipeline

### Stage 1: The Type Inferencer (The Elaborator)
*   **Location:** [ts/stlc/inferencer.go](file:///usr/local/google/home/jabolopes/.gemini/jetski/scratch/bapel/ts/stlc/inferencer.go) (invoked via [InferFunction](file:///usr/local/google/home/jabolopes/.gemini/jetski/scratch/bapel/ts/stlc/typechecker.go#L43)).
*   **Role:** The inferencer is responsible for all complex, heuristic-heavy type analysis and program transformations:
    *   Unification and constraint solving.
    *   Managing existential variables (`evars`).
    *   Method resolution and trait lookup ([LookupMethod](file:///usr/local/google/home/jabolopes/.gemini/jetski/scratch/bapel/ts/stlc/context.go#L244)).
    *   Automatic receiver adjustment (auto-borrowing via `Ptr::mk` and auto-dereferencing via `Ptr::get`).
    *   Inferring omitted generic type arguments.
*   **Output:** The inferencer outputs a **Fully Annotated Core IR**. By the end of inference, every single expression, sub-expression, variable, and let-binding has an explicit, concrete `.Type` assigned, and all implicit behaviors or syntactic sugar are explicitly desugared.

### Stage 2: The Typechecker (The Core Verifier / Trusted Kernel)
*   **Location:** [ts/stlc/typechecker_typecheck.go](file:///usr/local/google/home/jabolopes/.gemini/jetski/scratch/bapel/ts/stlc/typechecker_typecheck.go) (invoked via [TypecheckFunction](file:///usr/local/google/home/jabolopes/.gemini/jetski/scratch/bapel/ts/stlc/typechecker.go)).
*   **Role:** The typechecker acts as a lightweight, fast, deterministic **Trusted Kernel**. Because it expects the incoming IR from the inferencer to be fully annotated and desugared, the typechecker requires **zero unification, zero constraint solving, zero existential variables, and zero heuristics**.
*   **Verification:** It performs simple, mechanical bottom-up or top-down validation:
    *   In a function application `f x`, it verifies that `f.Type` is exactly `A -> B` and `x.Type` is exactly `A`, confirming the application has type `B`.
    *   In a let-binding `let x: T = val`, it verifies that `val.Type == T`.
*   **Safety Guarantee:** By keeping the verifier simple and small, it acts as an independent safety check. Even if the complex inferencer has a bug or makes an error during constraint solving, this simple verifier will immediately catch ill-typed IR before C++ code generation in [comp/cpp_printer.go](file:///usr/local/google/home/jabolopes/.gemini/jetski/scratch/bapel/comp/cpp_printer.go).

---

## 2. Desugaring Architecture: Syntactic vs. Type-Directed

A common question in compiler design is whether all desugaring should be separated into its own distinct pass before or after type inference. In Bapel, desugaring is bifurcated based on whether it requires type information:

### Syntactic Desugaring $\rightarrow$ Pre-Inference Pass
*   **Examples:** Converting `while` loops into `loop/break`, desugaring `@derive` attributes, expanding `impl Trait` parameter bounds into generic parameters, or expanding string interpolation.
*   **Location:** Handled in a separate pass before inference in [ast/desugar.go](file:///usr/local/google/home/jabolopes/.gemini/jetski/scratch/bapel/ast/desugar.go).
*   **Rationale:** These transformations are purely structural and depend only on syntax. Performing them before inference keeps the IR small and clean, meaning the inferencer never has to deal with syntactic sugar constructs.

### Type-Directed Desugaring (Elaboration) $\rightarrow$ Interleaved in Inferencer
*   **Examples:**
    *   **Method Resolution:** Rewriting dot-method calls (`s.size()`) to fully qualified calls (`String::size(s)` vs `Vector::size(s)`).
    *   **Receiver Adjustment:** Wrapping receivers in address-of (`Ptr::mk s`) or dereferencing (`Ptr::get s`).
    *   **Operator Overloading:** Rewriting `a + b` to `Add::add(&a, &b)` for custom structs while keeping primitive integer addition.
    *   **Shadowing Resolution:** Determining whether `s.foo` accesses a struct function field `foo` or invokes a method `foo()`.
*   **Location:** Interleaved directly inside the inferencer ([ts/stlc/inferencer.go](file:///usr/local/google/home/jabolopes/.gemini/jetski/scratch/bapel/ts/stlc/inferencer.go)).
*   **Why it cannot be separated into a distinct pass:**
    *   **Cannot be done BEFORE inference:** The receiver's type is unknown before inference, making it impossible to know whether `s.size()` calls `String::size`, `Vector::size`, or accesses a struct field.
    *   **Cannot be done AFTER inference:** In an expression like `let x = s.size() + 5`, the inferencer **must** resolve `s.size()` immediately to know that its return type is `i64` in order to typecheck the addition and infer the type of `x`. Postponing method resolution to a post-inference pass would cause type inference to stall or fail.

---

## 3. Requirements for Clean Separation

To maintain the simplicity of the typechecker and prevent logic duplication, the inferencer enforces three core invariants before handing the IR over to the typechecker:

1.  **100% AST/IR Annotation:** Every single `IrTerm` node (including intermediate sub-expressions in applications, projections, and tuples) must have its `.Type` field explicitly populated by the end of inference.
2.  **Existential Resolution ("Zonking"):** During inference, types contain temporary unification variables (`evars`). At the very end of [InferFunction](file:///usr/local/google/home/jabolopes/.gemini/jetski/scratch/bapel/ts/stlc/typechecker.go#L43), a final substitution pass (zonking) must replace all unresolved `evars` in `.Type` annotations with their solved concrete types.
3.  **Explicit Desugaring & Instantiation:** All implicit behaviors must be explicitly written into the IR during inference:
    *   **Method Calls:** Rewritten to explicit fully qualified function applications (`Type::method(adjusted_s, args)`).
    *   **Generic Instantiations:** When calling a polymorphic function `f` without explicit type arguments, the inferencer explicitly wraps it in an [AppTypeTerm](file:///usr/local/google/home/jabolopes/.gemini/jetski/scratch/bapel/ir/ir_term.go#L540) (`f [T]`), ensuring the simple typechecker never has to guess generic type parameters.
4.  **Elimination of Duplication:** With this architecture established, duplicated reduction and subtyping logic in [ts/stlc/typechecker.go](file:///usr/local/google/home/jabolopes/.gemini/jetski/scratch/bapel/ts/stlc/typechecker.go) (marked with `// TODO: Deduplicate with Inferencer`) can be removed, stripping the typechecker down to a minimal, trusted verifier.
