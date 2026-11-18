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
