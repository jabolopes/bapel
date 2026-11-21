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
Currently, `UnitChecker::add_trait_impl` (`comp/typecheck_module.h`) registers:
* **`InherentImpl`** under the nominal alias name (`impl->type_name`), while
* **`TraitImpl`** under `reduce_type(impl->type_name)` (the expanded structural `StructType` / `VariantType` / `MuType`).

This has two drawbacks:
1. **Structural/nominal mismatch with C++17:** If `type A = struct { x: i64 }` and `type B = struct { x: i64 }` both exist, `impl Size for A` matches `B` in Bapel's typechecker (same reduced type), but C++17 only generates `template <> struct traits::Size<::A>`, causing a C++ compile error if `Size::size` is called on `B`.
2. **2-step method lookup:** `Inferencer::lookup_method_for_receiver` (`ts/inferencer.h`) has to check both the unreduced type (for inherent methods) and the reduced type (for trait methods).

### Planned Design
Key `TraitImplBind` on the nominal `impl->type_name` (matching `InherentImpl` and C++'s nominal struct/variant translation), and update `ts/solver.h` and `Context::lookup_method` to match trait impls nominally.

### Follow-Up Test Coverage
Add cross-module import tests verifying that when module A imports module B, A sees all exported traits, trait implementations (via `Trait::method`, dot syntax, and generic bounds), and inherent methods on `AliasDecl`s (including those defined in B's `impls` files).

---

## 4. `impl` Blocks in `impls` Files (Equal-Weight Module Files)

**Status:** Deferred

### Motivation
In Bapel, base module files (`module foo`) and implementation files (`implements foo`) are meant to have equal weight for all declarations. Currently, however, `CppPrinter` only emits `.h` from the base module file (`BaseUnit`), and `unit->imported_trait_impls` collected from `impls` files are not typechecked during base-module compilation—so `impl` blocks defined in an `impls` file (e.g., `impl IrType` in `ir_type.bpl`) are not emitted into the module's `.h` header.

### Workaround Today
Keep thin `impl` blocks in the base module file (matching [`bapel/stl.bpl`](../bapel/stl.bpl)) delegating to `pub fn` helpers in `impls` files, or declare `impl` methods in the base module file.

### Planned Design
Either emit method declarations (signatures) for `impl` blocks from `impls` files in the module's `.h` header (with out-of-line definitions in each `impls` file's `.cc` for non-template methods, and typechecked inline definitions for template methods), or generate per-file headers included by the module's umbrella `.h`.

---

## 5. Abstract Type Definitions (`pub type T` + Private `type T = ...`)

**Status:** Deferred (documented in [`designs/modules.md`](modules.md))

### Motivation
Allow a module to export an opaque type (`pub type T`) while providing a private definition (`type T = ...`) visible only within the same module (base and implementation files), hiding the internal representation from importing modules. Currently `Context::is_wellformed` rejects defining an `AliasDecl` with the same name as an existing `NameDecl`.

---

## 6. Platform-Specific File Selection in `impls`

**Status:** Deferred (mentioned in [`designs/modules.md`](modules.md))

### Motivation
Allow conditional selection of implementation files (e.g., `*_linux.bpl`, `*_macos.bpl`) in a module's `impls` section so platform-specific definitions can satisfy shared module-level declarations.

---

## 7. Module-Qualified Imports and C++ Namespaces

**Status:** Deferred (`export` $\to$ `pub` in normalized IR is done; module-qualified access & C++ namespaces deferred)

### Motivation
* Currently, normalized IR prints `export` (legacy term) instead of `pub`, and imported symbols from different modules share a flat namespace with no way to disambiguate colliding names.
* C++ uses `namespace`s to avoid polluting the global namespace, whereas Bapel uses module imports.

### Planned Design Direction
1. **Normalized IR:** Replace `export` with `pub` and preserve each imported symbol's defining `ModuleID`.
2. **Suffix-qualified module access:** Importing `bapel.compiler.ir` makes `IrType` accessible as `IrType` (when unambiguous), `ir.IrType`, `compiler.ir.IrType`, or `bapel.compiler.ir.IrType`—no explicit import-aliasing syntax required by default.
3. **Coexistence with C++ namespaces (`module.ns::sym`):** Use `.` for Bapel module paths and `::` for C++ / inherent / trait namespaces (e.g., `stl.fs::read_file`). Pure Bapel modules can emit into `namespace a::b::c { ... }`, while `@bpl` C++ FFI declarations keep their existing C++ namespaces (e.g. `::fs::read_file`).

---

## 8. Accurate Source Positions for `@bpl:` FFI Declarations (`<stdin>:1`)

**Status:** Deferred

### Motivation
When `comp::query_source_file` parses `@bpl:` FFI declarations from C++ `.h` headers (`comp/querier.h`), it currently passes `"<stdin>"` as the filename and parses each line as line 1 (`parser::parse_decl(decl_text, "<stdin>")`). As a result, normalized/typechecked golden test outputs and error messages show `/* <stdin>:1 */` instead of the actual header filename and line number.

### Planned Fix
Track the 1-based line number while scanning `.h` files in `comp::query_source_file`, pass the real header `filename` to `parser::parse_decl`, and set `decl.pos.begin_line` / `decl.pos.end_line` to the actual line number in the `.h` file.

---

## 9. `ImplDecl` Well-Formedness Checking in `ts/wellformed.h`

**Status:** Deferred

### Motivation
Currently `IrDeclCase::ImplDecl` only appears in `import_decls` (already typechecked in its defining module), so `ts::is_wellformed_decl` in `ts/wellformed.h` returns `true` for `ImplDecl`. If `ImplDecl` is ever used for local declarations or validated standalone, `is_wellformed_decl` should check all types inside it.

### Planned Fix
In `ts::is_wellformed_decl` (`case ir::IrDeclCase::ImplDecl`):
- Verify `decl.impl_decl` is non-null, with no duplicate `type_params` or duplicate method names.
- Enter scope with `Self` and `impl_decl->type_params`, and validate `is_wellformed_type` on `type_name` (and `trait_type` for `TraitImpl`).
- For each method signature, bind its `type_params` and validate `is_wellformed_type` on all `args` and `ret_type`.

---

## 10. Handle `ImplDecl` Directly in `Context::add_symbol`

**Status:** Deferred

### Motivation
Currently `TraitImplBind` stores `reduced_type` and `reduced_trait_type`, which requires calling `ts::Typechecker::reduce_type` when registering a `TraitImpl` or `ImplDecl`. Because `ts/reduce_type.h` includes `ts/context.h`, `Context::add_symbol` cannot call `reduce_type`, so `ModuleChecker` has to register `ImplDecl`s separately in `comp/typecheck_module.h`.

### Planned Fix
Store the unreduced `type_name` and `trait_type` in `TraitImplBind` (matching `InherentImpl`), then handle `IrDeclCase::ImplDecl` (both `InherentImpl` and `TraitImpl`) directly inside `Context::add_symbol` in `ts/context.h`, removing `build_impl_decl_bindings` and `add_impl_decl_to_ctx` from `comp/typecheck_module.h`.

---

## 11. Eliminate `codegen::topo_sort_decls` via In-Order `__anonym_*` Insertion

**Status:** Deferred

### Motivation
`IrTypedModule.decls` is already topologically sorted by `ModuleChecker`. The only reason `bin/codegen_impl.h` still runs `codegen::topo_sort_decls` / `codegen::topo_sort_alias_decls` is that `record_anonymous_types_from_unit` appends newly-lifted `__anonym_<hash>` structs to the end of `unit->decls`, requiring a second sort (which also reverses independent zero-in-degree declarations due to its LIFO stack).

### Planned Fix
Insert each newly-lifted `__anonym_<hash>` `AliasDecl` immediately before the declaration that uses it during the single pass over `unit->decls` (and before `functions` / `trait_impls` for anonymous structs lifted from function/method signatures). This preserves `IrTypedModule`'s exact declaration order and removes `codegen::topo_sort_decls` and `codegen::topo_sort_alias_decls` completely.

---

## 12. `DefID`-Based Name Resolution & Strict Direct-Import Scoping

**Status:** Deferred

### Motivation
Currently, `ModuleChecker` computes `reachable_imports(file.imports)` (the transitive closure of imported modules) and loads all transitive `import_decls` into `file_ctx[F]`. Because Bapel uses flat string names in `ts::Context`, a file `F` that imports `bapel.stl` (which depends on `bapel.core`) can currently reference `bapel.core` symbols by name even without directly importing `bapel.core`.

### Planned Design Direction
Separate **Name Resolution** from **Typechecking** (as in Go and Rust):
1. **Per-file Lexical Scope (`string` $\to$ `DefID`):** Contains only local variables, module-local declarations, and **direct** `file.imports`. Source identifiers are resolved exclusively here and rewritten to canonical `DefID`s `(ModuleID, Name)`.
2. **Global Definition Store (`DefID` $\to$ `Definition`):** Holds all local, direct, and transitive definitions in `IrTypedModule`. Typechecking looks up resolved `DefID`s here, so transitive types referenced by exported APIs typecheck cleanly without leaking into a file's lexical name scope.
