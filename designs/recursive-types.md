# Equirecursive Types in Bapel (System $F_\omega^\mu$)

This document describes how equirecursive types work in Bapel, grounded in **System $F_\omega$ with Equirecursive Types ($F_\omega^\mu$)** (Cai, Giarrusso, Ostermann, POPL 2016).

---

## 1. Overview

Bapel supports **self-recursive types** (including regular generic self-recursive types) using natural surface syntax—no explicit `mu` keyword is required in `.bpl` files:

```bapel
pub type Node = struct {
  val: i64,
  next: SharedPtr Node
}

pub type Tree ['a] = variant {
  leaf 'a,
  branch (SharedPtr (Tree 'a), SharedPtr (Tree 'a))
}
```

In Core IR, recursive types are represented explicitly with a **fixed-point binder** ($\mu X :: \ast.\, T$), rather than cyclic `AliasDecl`s. This keeps `topo_sort_decls` acyclic and ensures that kind inference, reduction, subtyping, and unification terminate cleanly.

*(Mutually recursive types are deferred; see [`designs/features.md`](features.md).)*

---

## 2. Core IR Representation

`IrType` ([`bin/ir_base.h`](../bin/ir_base.h), [`bin/ir_type.h`](../bin/ir_type.h)) includes the `MuType` constructor:

$$T ::= \dots \mid \mu X :: \ast.\, T$$

```cpp
enum class IrTypeCase {
  // ...
  MuType = 6,
  // ...
};

struct MuTypeData {
  std::string var;
  IrKind kind;              // Restricted to TypeKind (*) for decidability
  std::shared_ptr<IrType> body;
};
```

Key Core IR operations:
* **Formatting (`to_string`):** Prints `mu 'self. <body>`.
* **Capture-avoiding substitution (`bin/ir_substitute_type.h`):** `ir::substitute_type` respects shadowing of `mu->var` and α-renames `mu->var` when it would capture a free variable of the replacement type.
* **Context opening (`ts/context.h`):** `Context::add_fresh_type` opens `MuType` (alongside `ForallType` and `LambdaType`) by introducing a fresh type parameter binding and substituting it into `body`.
* **Variable renaming & predication (`ts/rename_vars.h`, `ts/predicative.h`):** `MuType` binders are α-renamed canonically without being hoisted into outer `forall` quantifiers.

### Why $\mu$ Is Restricted to Kind $\ast$

In $F_\omega^\mu$, allowing $\mu$ at higher kinds ($\kappa_1 \to \kappa_2$) combined with full $F_\omega$ $\beta$-reduction makes type equivalence **undecidable**. Restricting $\mu$ to proper types (kind $\ast$) keeps type equivalence and subtyping **decidable**, while supporting all regular generic recursive ADTs by placing $\lambda$ outside $\mu$:

$$\text{Tree} \equiv \lambda \alpha :: \ast.\; \mu X :: \ast.\; \text{variant}\{\text{leaf } \alpha,\; \text{branch } (\text{SharedPtr } X,\; \text{SharedPtr } X)\}$$

Non-regular (nested) recursive types (e.g., `type Nested ['a] = variant { nil (), cons ('a, Nested ('a, 'a)) }`) are rejected during desugaring.

---

## 3. Desugaring Surface Syntax to `MuType`

Self-recursive desugaring runs in [`comp/desugar.h`](../comp/desugar.h) (`SelfRecursiveRewriter`, `desugar_self_recursive_decls`, `desugar_self_recursive_unit`), invoked from [`comp/normalizer.h`](../comp/normalizer.h) right after name resolution:

For each `AliasDecl` `type F ['a, 'b] = body` (represented in IR as `LambdaType('a, *, LambdaType('b, *, body))`):
1. **Self-reference check:** If `NameType("F")` does not appear in `body`, the `AliasDecl` is left untouched.
2. **Regularity check & rewriting:**
   * Choose a fresh binder name (`self`, avoiding collisions with `'a, 'b` and free type variables in `body`).
   * Replace every `AppType` chain headed by `NameType("F")` whose arguments are **exactly** `['a, 'b]` in order with `VarType(self)` (for 0-parameter aliases, bare `NameType("F")` becomes `VarType(self)`).
   * If `NameType("F")` appears partially applied, unapplied, or with different arguments, reject it with `"non-regular recursive type \"F\": recursive occurrences must match the declared type parameters"`.
3. **Wrap in `MuType`:**
   ```
   type Node = mu 'self. struct { val: i64, next: SharedPtr 'self }
   type Tree :: * -> * = fun (a) (mu 'self. variant { leaf 'a, branch (SharedPtr 'self, SharedPtr 'self) })
   ```

Because self-references are replaced by a bound `VarType`, the resulting `AliasDecl` has no self-edge in `topo_sort_decls`.

---

## 4. Type System

### 4.1 Well-Formedness & Contractiveness (`ts/wellformed.h`)

To prevent vacuous infinite types like `type Bad = Bad` ($\mu X.\, X$), `is_wellformed_type` checks that every `MuType`:
* Has kind `*` (`t.mu->kind.is_type_kind()`).
* Is **contractive** (its body is not bare `VarType(t.mu->var)`).
* Has a well-formed body under `c.add_fresh_type(t)`.

`UnitChecker::add_decl` ([`comp/typecheck_unit.h`](../comp/typecheck_unit.h)) validates local declarations via `ts::is_wellformed_decl` before adding them to the context.

### 4.2 Kinding (`ts/infer_kind.h`)

$$\frac{\Delta, X :: \ast \vdash T :: \ast}{\Delta \vdash (\mu X :: \ast.\, T) :: \ast} \quad (\text{K-Mu})$$

`KindInferencer` opens `MuType` with `add_fresh_type`, verifies `body :: *`, and returns `*`. Because `X` is a bound `VarType`, kind inference never loops.

### 4.3 Equirecursive Unfolding (`ts/reduce_type.h`)

Under equirecursive semantics, a $\mu$-type is equivalent to its one-step unfolding:

$$\mu X.\, T \;\equiv\; [X \mapsto \mu X.\, T]\, T$$

* `reduce_type` (weak-head normal form) expands non-recursive `AliasDecl`s and performs $\beta$-reduction, and **stops at `MuType`**.
* `unfold_mu_type(mu X. T)` performs one step of unfolding (`substitute_type(T, VarType(X), mu X. T)`), followed by head reduction.
* Elimination sites in [`ts/inferencer.h`](../ts/inferencer.h) and [`ts/typecheck.h`](../ts/typecheck.h) (`match`, field projection, `set`, struct construction, variant injection) call `unfold_mu_type` on demand.

### 4.4 Coinductive Subtyping & Unification (`ts/subtype.h`, `ts/unify.h`)

`Subtyper` ([`ts/subtype.h`](../ts/subtype.h)) and `Unifier` ([`ts/unify.h`](../ts/unify.h)) maintain a coinductive assumption list `assumed_` (`std::vector<std::pair<ir::IrType, ir::IrType>>`):

* **Reflexivity / Coinductive hypothesis:** If `ir::equals_type(left, right)` or `(left, right)` is already in `assumed_`, return `true`.
* **Three `MuType` cases** (after head reduction):
  1. **`MuType` vs `MuType` (α-equivalence):** Push `(left, right)` onto `assumed_`, open both with the *same* fresh variable, and compare bodies.
  2. **`MuType` vs non-`MuType`:** Push `(left, right)` onto `assumed_` and compare `unfold_mu_type(left)` against `right`.
  3. **non-`MuType` vs `MuType`:** Push `(left, right)` onto `assumed_` and compare `left` against `unfold_mu_type(right)`.
* **Unification (`solve_type`):** Unassigned `ExistVarType`s solve directly to `MuType` without unfolding.

---

## 5. C++17 Code Generation (`bin/codegen_impl.h`)

Recursive `struct` and `variant` types are emitted as nominal C++17 structs with forward declarations:

1. **Targeted forward declarations (`is_mu_alias`):** `CppPrinter::print_decls` emits forward declarations (`struct Node;`, `template <typename a> struct Tree;`) *only* for `AliasDecl`s whose body (after peeling `LambdaType`s) is a `MuType`. Non-recursive aliases are unaffected.
2. **Self-reference mapping (`mu_self_names_`):** When `print_alias_decl` peels `LambdaType`s and enters a `MuType`, it binds `mu_self_names_[var]` to the nominal C++ name (`::Node` or `::Tree<a>`) for the duration of printing the body. `print_type` resolves `VarType`s in `mu_self_names_` to that nominal name.
3. **Nominal lookup across α-renaming (`canonical_hash_type`):** `find_decl_for_type` hashes `MuType`s with α-normalized binder names so unfolded/renamed `MuType`s (e.g., `mu 'i. ...`) match their declared nominal C++ type (`::IntList`).

---

## 6. Test Coverage

* **Unit tests:** [`tests/ir_type_test.cc`](../tests/ir_type_test.cc) (`Formatting`, `SubstituteType` shadowing & capture avoidance, `MuTypeEquality`).
* **Positive golden test:** [`tests/testdata/in/recursive_types.in`](../tests/testdata/in/recursive_types.in) (self-recursive struct, self-recursive variant, generic self-recursive variant `Tree ['a]`, and C++17 syntax validation).
* **Negative tests:**
  * [`tests/testdata/in/bad_recursive_type_non_contractive.in`](../tests/testdata/in/bad_recursive_type_non_contractive.in) (`type Bad = Bad`)
  * [`tests/testdata/in/bad_recursive_type_non_regular.in`](../tests/testdata/in/bad_recursive_type_non_regular.in) (non-regular generic recursion)
  * [`tests/testdata/in/bad_recursive_type_partial_app.in`](../tests/testdata/in/bad_recursive_type_partial_app.in) (partially applied self-reference)
  * [`tests/testdata/in/bad_recursive_type_mutual.in`](../tests/testdata/in/bad_recursive_type_mutual.in) (mutually recursive aliases)
