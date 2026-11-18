# Equirecursive Types for Bapel (System $F_\omega^\mu$)

This document proposes a principled design for recursive and mutually recursive types in Bapel, grounded in **System $F_\omega$ with Equirecursive Types ($F_\omega^\mu$)** (Cai, Giarrusso, Ostermann, POPL 2016).

---

## 1. Motivation

Writing the Bapel compiler in Bapel itself (see [`designs/migration.md`](migration.md)) requires self-recursive and mutually recursive algebraic data types—for example:

```bapel
pub type AppData = struct {
  fun: SharedPtr Expr,
  arg: SharedPtr Expr
}

pub type Expr = variant {
  lit i64,
  app AppData
}

pub type Tree ['a] = variant {
  leaf 'a,
  branch (SharedPtr (Tree 'a), SharedPtr (Tree 'a))
}
```

Currently, Bapel's Core IR is **System $F_\omega$ without a fixed-point ($\mu$) binder**:
* All `type X = T` declarations are transparent, non-recursive type aliases (`AliasDecl`).
* `topo_sort_alias_decls` rejects any cycle between `AliasDecl`s.
* `reduce_type`, `infer_kind`, `subtype`, and `unify` assume that expanding aliases always terminates in a finite tree.

Patching individual passes with ad-hoc `visiting` sets is fragile and unsound (e.g., it breaks kinding for parameterized recursive types and causes lockstep unfolding divergence in `subtype`). Instead, Bapel should represent recursive types explicitly with a **fixed-point binder** ($\mu X.\, T$) in the Core IR.

---

## 2. Type Theory: System $F_\omega^\mu$

### 2.1 Syntax of Types in Core IR

We extend `IrType` ([`bin/ir_type.h`](../bin/ir_type.h)) with a single new constructor:

$$T ::= \dots \mid \mu X :: \ast.\, T$$

In Core IR:
```cpp
enum class IrTypeCase {
  // ... existing cases ...
  MuType = 11,
};

struct MuTypeData {
  std::string var;
  IrKind kind;              // Restricted to TypeKind (*) for decidability
  std::shared_ptr<IrType> body;
};
```

### 2.2 Why Restrict $\mu$ to Kind $\ast$?

In $F_\omega^\mu$, allowing $\mu$ at arbitrary higher kinds ($\kappa_1 \to \kappa_2$) combined with full $F_\omega$ $\beta$-reduction makes type-level normalization able to encode the untyped $\lambda$-calculus, rendering type equivalence **undecidable**.

Restricting $\mu$ to proper types (kind $\ast$) keeps type equivalence and subtyping **decidable**, while still supporting **all regular generic recursive ADTs** by placing $\lambda$ outside $\mu$:

```bapel
type Tree ['a] = variant {
  leaf 'a,
  branch (SharedPtr (Tree 'a), SharedPtr (Tree 'a))
}
```

Desugars to:
$$\text{Tree} \equiv \lambda \alpha :: \ast.\; \mu X :: \ast.\; \text{variant}\{\text{leaf } \alpha,\; \text{branch } (\text{SharedPtr } X,\; \text{SharedPtr } X)\}$$

*(Non-regular / nested recursive types like `type Nested ['a] = variant { nil (), cons ('a, Nested ('a, 'a)) }` are rejected at desugaring time, which is standard and not needed for AST/IR trees.)*

### 2.3 Kinding Rule

$$\frac{\Delta, X :: \ast \vdash T :: \ast}{\Delta \vdash (\mu X :: \ast.\, T) :: \ast} \quad (\text{K-Mu})$$

In `ts/infer_kind.h`:
* `MuType(X, *, body)` extends the context with `X :: *` (just like `ForallType` / `LambdaType` via `add_fresh_type`), infers the kind of `body`, verifies `body :: *`, and returns `*`.
* Because `X` is a bound type variable (`VarType`), `infer_kind` **never loops**—no `visiting` set is needed!

### 2.4 Contractiveness (Guardedness)

To prevent vacuous infinite types like $\mu X.\, X$ or $\mu X.\, (\lambda \alpha.\, \alpha)\, X$, every $\mu X.\, T$ must be **contractive**:
* $X$ must not appear unguarded in head position of $T$ (after $\beta$-reduction).
* Any occurrence of $X$ inside a type constructor (`StructType`, `VariantType`, `TupleType`, `FunType`, `ArrayType`, or `AppType` headed by a `NameDecl` like `SharedPtr` / `Ptr` / `Vector`) is **guarded** (contractive).
* Non-contractive definitions (e.g., `type Bad = Bad`) are rejected with a clear compile error during desugaring/well-formedness checking.

### 2.5 Equirecursive Unfolding (Weak-Head Normal Form)

Under equirecursive semantics, a $\mu$-type is equivalent to its one-step unfolding:

$$\mu X.\, T \;\equiv\; [X \mapsto \mu X.\, T]\, T$$

Crucially, **we do NOT eagerly unfold $\mu$-types deeply inside `StructType` / `VariantType` / `AppType` arguments**:
* **`reduce_type` (WHNF / head reduction):**
  * Expands non-recursive `AliasDecl`s and performs $\beta$-reduction $(\lambda \alpha.\, T_1)\, T_2 \to [\alpha \mapsto T_2]\, T_1$.
  * Stops at `MuType` unless a caller specifically asks to **unfold one step** (`unfold_mu(mu X. T) = substitute(T, X, mu X. T)`), e.g., when projecting a field from a struct, matching a variant, or comparing a `MuType` against a structural type in `subtype` / `unify`.

### 2.6 Coinductive Subtyping & Unification

Because `MuType` is only unfolded on demand, we implement subtyping via a `Subtyper` helper class (matching `TypeReducer` and `KindInferencer`) holding the **coinductive hypothesis set** `std::vector<std::pair<ir::IrType, ir::IrType>> assumed_` (written $\Sigma$):

* **Reflexivity:** If `ir::equals_type(left, right)`, return `true` immediately.
* **Coinductive check:** If `(left, right)` is already in `assumed_` (checked via `ir::equals_type`), return `true`.
* **Three `MuType` cases** (after `reduce_and_predicate_type`, which stops at `MuType`):
  1. **`MuType` vs `MuType` (α-equivalence):** Push `(left, right)` onto `assumed_`, open both with the *same* fresh variable via `Context::add_fresh_type` (handling `mu 'x. ...` vs `mu 'y. ...`), and check `subtype(left_body, right_body)`.
  2. **`MuType` vs non-`MuType`:** Push `(left, right)` onto `assumed_` and check `subtype(unfold_mu_type(left), right)`.
  3. **non-`MuType` vs `MuType`:** Push `(left, right)` onto `assumed_` and check `subtype(left, unfold_mu_type(right))`.
* **Unification (`ts/unify.h`):**
  * If one side is an unassigned `ExistVarType` (`^0`) and the other is `MuType`, solve `^0 = MuType` **directly** without unfolding.
  * Otherwise, apply the exact same three `MuType` cases with an `assumed_` vector.

Because every `MuType` is contractive and restricted to kind $\ast$, the set of reachable sub-terms is finite, so coinductive subtyping and unification are **guaranteed to terminate**.

---

## 3. Desugaring Surface Syntax to `MuType`

Users continue to write natural Bapel syntax—no explicit `mu` keyword is required in surface `.bpl` files (though `mu` is printed in Core IR golden tests):

### 3.1 Self-Recursive Types (Phase 1)
For an `AliasDecl` `type F ['a, 'b] = body` (represented in IR as `LambdaType('a, *, LambdaType('b, *, body))`):
1. **Detect self-reference:** Check whether `NameType("F")` appears anywhere in `body`. If not, leave the `AliasDecl` untouched.
2. **Regularity check & rewriting:**
   * Choose a fresh binder name (e.g. `self`, avoiding collisions with `'a, 'b` or free type variables in `body`).
   * Replace every `AppType` chain headed by `NameType("F")` whose arguments are **exactly** `['a, 'b]` in order with `VarType(self)` (for 0-parameter aliases like `Node`, bare `NameType("Node")` becomes `VarType(self)`).
   * If `NameType("F")` appears partially applied, unapplied, or with different arguments (e.g., `Tree ('a, 'a)`), reject it immediately with `"non-regular recursive type \"F\": recursive occurrences must match the declared type parameters"`.
3. **Wrap in `MuType`:**
   ```
   type Node = mu 'self. struct { val: i64, next: SharedPtr 'self }
   type Tree :: * -> * = fun (a) (mu 'self. variant { leaf 'a, branch (SharedPtr 'self, SharedPtr 'self) })
   ```
4. **Trusted Kernel verification (`ts/wellformed.h`, `ts/infer_kind.h`):**
   Independently verifies `(K-Mu)` kinding, contractiveness (no unguarded `self` at head position), and well-formedness on the resulting `MuType` IR.

### 3.2 Mutually Recursive Types — SCC & Bekić's Theorem *(Deferred to Phase 2)*
For a strongly connected component (SCC) of mutually recursive aliases, e.g.:
```bapel
type AppData = struct { fun: SharedPtr Expr, arg: SharedPtr Expr }
type Expr = variant { lit i64, app AppData }
```
Using **Bekić's rule** (standard in $F_\omega^\mu$):
1. Substitute `AppData` into `Expr` to form a single fixed point:
   $$\text{Expr} = \mu E.\; \text{variant}\{\text{lit } \text{i64},\; \text{app } \text{struct}\{\text{fun}: \text{SharedPtr } E,\; \text{arg}: \text{SharedPtr } E\}\}$$
2. Define `AppData` non-recursively in terms of `Expr`:
   $$\text{AppData} = \text{struct}\{\text{fun}: \text{SharedPtr Expr},\; \text{arg}: \text{SharedPtr Expr}\}$$

In general, any SCC $\{A_1, \dots, A_n\}$ of mutually recursive types can be resolved into nested $\mu$-binders by eliminating one variable at a time.

> **Scope note:** SCC detection and Bekić elimination for mutually recursive types are **deferred to Phase 2**. Phase 1 focuses exclusively on **self-recursive types** (including generic self-recursive types like `Tree ['a]`), keeping the initial implementation small and easy to verify.

---

## 4. C++17 Code Generation (`bin/codegen_impl.h`)

In C++17, recursive `struct` and `variant` types are naturally represented by **nominal C++ structs with forward declarations**:

1. **Forward declarations for recursive aliases only:**
   * In Phase 1, an `AliasDecl` needs a forward declaration if and only if its body (after peeling any `LambdaType`s) is a `MuType` (no cycle-finder needed in `CppPrinter`).
   * Before emitting definitions in `CppPrinter::print_decls`, emit `print_decl(d)` for those `MuType` aliases only (`struct Node;`, `template <typename a> struct Tree;`).
   * Existing non-recursive golden tests remain 100% unchanged.
2. **Emitting `MuType` bodies in `print_alias_decl`:**
   * Maintain `std::map<std::string, std::string> mu_self_names_;` on `CppPrinter`.
   * When `print_alias_decl(id, typ)` peels `LambdaType`s (`tvars = {"a"}`) and hits `MuType(var, *, body)`, bind `mu_self_names_[var]` to the nominal C++ name (`"::" + id` or `"::" + id + "<a, ...>"`), recurse into `print_alias_decl(id, body)`, and erase `var` on exit.
   * In `print_type` for `VarType`: if `typ.var` is in `mu_self_names_`, emit the nominal C++ name (`::Node` / `::Tree<a>`).
3. **Emitting `MuType` / `VariantType` references in expressions:**
   * When `print_type` encounters a `MuType` (or `VariantType` outside its own `print_alias_decl` definition), use `find_decl_for_type(typ, decl)` to emit the nominal C++ name (`::IntList`) so pointers like `::Ptr<::IntList>` and `::SharedPtr<::IntList>` typecheck in C++.

---

## 5. Implementation & Verification Plan

### Phase 1: Self-Recursive Types (MVP)
1. **Core IR (`bin/ir_base.h`, `bin/ir_type.h`):**
   * Add `IrTypeCase::MuType`, `MuTypeData`, `new_mu_type`, `to_string` (`mu 'x. ...`), `equals_type`, `collect_free_vars`, and capture-avoiding `substitute_type`.
2. **Self-Recursive Desugaring (`bin/ir_unit.h` / `comp/resolver.h`):**
   * Detect self-recursive `AliasDecl`s (`type X = ... X ...` and `type X ['a] = ... (X 'a) ...`), enforce regularity, and wrap their body in `MuType` (preserving acyclic topological sort for `order.in`).
3. **Type System (`ts/`):**
   * `ts/context.h`: Extend `Context::add_fresh_type` to open `MuType`.
   * `ts/rename_vars.h`: Handle `MuType` via `with_type_substitution`.
   * `ts/predicative.h`: Predicate `MuType` body without hoisting its binder into `tvars_`.
   * `ts/wellformed.h` & `ts/infer_kind.h`: Implement `(K-Mu)` and contractiveness check.
   * `ts/reduce_type.h`: Keep `MuType` unexpanded in `reduce_type`, add `unfold_mu_type` for one-step head unfolding.
   * `ts/subtype.h` & `ts/unify.h`: Implement coinductive `Subtyper` and `unify` with `assumed_` vector, α-equivalence, `solve_type` for `MuType`, and on-demand `unfold_mu_type`.
   * `ts/inferencer.h` & `ts/typecheck.h`: Unfold `MuType` at elimination sites (`match`, field projection, `set`, `struct`, variant injection).
4. **C++ Codegen (`bin/codegen_impl.h`):**
   * Handle `MuType` in `record_anonymous_types` and `get_free_vars_from_type`.
   * Emit forward declarations for `MuType` aliases only, map `MuType` self-variable to nominal C++ type name (`mu_self_names_`), and preserve nominal `VariantType` / `MuType` names in references.
5. **Testing:**
   * Unit tests in `tests/ir_type_test.cc` for `MuType` formatting, substitution, and equality.
   * Golden test `tests/testdata/in/recursive_types.in` covering self-recursive struct, self-recursive variant, generic self-recursive variant, and non-contractive error rejection.
   * Verify `make test` and `clang++ -std=c++17 -fsyntax-only` pass with 0 errors.

### Phase 2: Mutually Recursive Types *(Deferred)*
* Add SCC detection in `topo_sort_alias_decls` and Bekić elimination to desugar mutually recursive `AliasDecl` groups (e.g., `AppData` $\leftrightarrow$ `Expr`) into nested `MuType`s, plus multi-alias forward declarations in `CppPrinter`.
