# Future Features & Long-Term Work

This document collects potential language and compiler features deferred for the long term.

---

## 1. Mutually Recursive Types (SCC & Bekić's Theorem)

**Status:** Deferred (Phase 2 of [`designs/recursive-types.md`](recursive-types.md))

### Motivation
Phase 1 of [`designs/recursive-types.md`](recursive-types.md) supports self-recursive types (including regular generic self-recursive types like `Tree ['a]`) via a $\mu$-binder (`MuType`) in Core IR. Mutually recursive type aliases—such as `Expr` $\leftrightarrow$ `Stmt` or a variant referencing named payload structs that refer back to the variant—are currently rejected during topological sorting:

```bapel
pub type AppData = struct {
  fun: SharedPtr Expr,
  arg: SharedPtr Expr
}

pub type Expr = variant {
  lit i64,
  app AppData
}
```

### Workarounds Today
- Inline anonymous `struct { ... }` directly inside `variant` arms, or
- Parameterize helper structs over the recursive type (e.g. `type AppData ['e] = struct { fun: SharedPtr 'e, arg: SharedPtr 'e }`, then use `AppData Expr` inside `Expr`).

### Planned Design
1. **SCC Detection:** Replace cycle rejection in `topo_sort_alias_decls` with Tarjan's / Kosaraju's strongly connected component (SCC) decomposition.
2. **Bekić Elimination:** Desugar each mutually recursive SCC $\{A_1, \dots, A_n\}$ into nested `MuType` binders by substituting one alias into the next (standard in System $F_\omega^\mu$).
3. **C++ Codegen:** Emit forward declarations for all aliases participating in an SCC before their definitions in `CppPrinter::print_decls`.

---

## 2. Extended Pattern Matching

**Status:** Deferred (beyond Phase 1 of [`designs/pattern-matching.md`](pattern-matching.md))

### Potential Extensions
* **Wildcard payload binding (`nil _ => 0`):** Allow `_` as the payload variable name in a tag arm without emitting a C++ `auto &_ = std::get<idx>(...);` binding (avoiding unused-variable warnings).
* **Literal & primitive matching:** Support `match` on `i8`, `i64`, `bool`, and `StringView` (e.g., `match ch { 'a' => ..., '0' => ..., _ => ... }`), useful when porting the lexer in Phase 6 of [`designs/migration.md`](migration.md).
* **Nested & destructuring patterns:** Tuple/struct destructuring inside arms (e.g., `cons (x, xs) => ...`, `app { fun, arg } => ...`), or-patterns (`a | b => ...`), and match guards (`lit x if x > 0 => ...`).

---

## 3. Nominal `TraitImplBind` Keying

**Status:** Deferred (see [`designs/traits.md`](traits.md))

### Motivation
Currently, `UnitChecker::add_trait_impl` (`comp/typecheck_unit.h`) registers:
* **`InherentImpl`** under the nominal alias name (`impl->type_name`), while
* **`TraitImpl`** under `reduce_type(impl->type_name)` (the expanded structural `StructType` / `VariantType` / `MuType`).

This has two drawbacks:
1. **Structural/nominal mismatch with C++17:** If `type A = struct { x: i64 }` and `type B = struct { x: i64 }` both exist, `impl Size for A` matches `B` in Bapel's typechecker (same reduced type), but C++17 only generates `template <> struct traits::Size<::A>`, causing a C++ compile error if `Size::size` is called on `B`.
2. **2-step method lookup:** `Inferencer::lookup_method_for_receiver` (`ts/inferencer.h`) has to check both the unreduced type (for inherent methods) and the reduced type (for trait methods).

### Planned Design
Key `TraitImplBind` on the nominal `impl->type_name` (matching `InherentImpl` and C++'s nominal struct/variant translation), and update `ts/solver.h` and `Context::lookup_method` to match trait impls nominally.
