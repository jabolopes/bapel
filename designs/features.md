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

## 3. Abstract Type Definitions (`pub type T` + Private `type T = ...`)

**Status:** Deferred (documented in [`designs/modules.md`](modules.md))

### Motivation
Allow a module to export an opaque type (`pub type T`) while providing a private definition (`type T = ...`) visible only within the same module (base and implementation files), hiding the internal representation from importing modules. Currently `Context::is_wellformed` rejects defining an `AliasDecl` with the same name as an existing `NameDecl`.

---

## 4. Platform-Specific File Selection in `impls`

**Status:** Deferred (mentioned in [`designs/modules.md`](modules.md))

### Motivation
Allow conditional selection of implementation files (e.g., `*_linux.bpl`, `*_macos.bpl`) in a module's `impls` section so platform-specific definitions can satisfy shared module-level declarations.

---

## 5. Module-Qualified Imports and C++ Namespaces

**Status:** Deferred (`export` $\to$ `pub` in normalized IR is done; module-qualified access & C++ namespaces deferred)

### Motivation
* Currently, imported symbols from different modules share a flat namespace with no way to disambiguate colliding names.
* C++ uses `namespace`s to avoid polluting the global namespace, whereas Bapel uses module imports.

### Planned Design Direction
1. **Preserve defining `ModuleID`** on imported symbols.
2. **Suffix-qualified module access:** Importing `bapel.compiler.ir` makes `IrType` accessible as `IrType` (when unambiguous), `ir.IrType`, `compiler.ir.IrType`, or `bapel.compiler.ir.IrType`—no explicit import-aliasing syntax required by default.
3. **Coexistence with C++ namespaces (`module.ns::sym`):** Use `.` for Bapel module paths and `::` for C++ / inherent / trait namespaces (e.g., `stl.fs::read_file`). Pure Bapel modules can emit into `namespace a::b::c { ... }`, while `@bpl` C++ FFI declarations keep their existing C++ namespaces (e.g. `::fs::read_file`).

---

## 6. Accurate Source Positions for `@bpl:` FFI Declarations (`<stdin>:1`)

**Status:** Deferred

### Motivation
When `comp::query_source_file` parses `@bpl:` FFI declarations from C++ `.h` headers (`comp/querier.h`), it currently passes `"<stdin>"` as the filename and parses each line as line 1 (`parser::parse_decl(decl_text, "<stdin>")`). As a result, normalized/typechecked golden test outputs and error messages show `/* <stdin>:1 */` instead of the actual header filename and line number.

### Planned Fix
Track the 1-based line number while scanning `.h` files in `comp::query_source_file`, pass the real header `filename` to `parser::parse_decl`, and set `decl.pos.begin_line` / `decl.pos.end_line` to the actual line number in the `.h` file.

---

## 7. `ImplDecl` Well-Formedness Checking in `ts/wellformed.h`

**Status:** Deferred

### Motivation
Currently `IrDeclCase::ImplDecl` only appears in `import_decls` (already typechecked in its defining module), so `ts::is_wellformed_decl` in `ts/wellformed.h` returns `true` for `ImplDecl`. If `ImplDecl` is ever used for local declarations or validated standalone, `is_wellformed_decl` should check all types inside it.

### Planned Fix
In `ts::is_wellformed_decl` (`case ir::IrDeclCase::ImplDecl`):
- Verify `decl.impl_decl` is non-null, with no duplicate `type_params` or duplicate method names.
- Enter scope with `impl_decl->type_params`, and validate `is_wellformed_type` on `type_name` (and `trait_type` for `TraitImpl`).
- For each method signature, bind its `type_params` and validate `is_wellformed_type` on all `args` and `ret_type`.

---

## 8. Eliminate `codegen::topo_sort_decls` via In-Order `__anonym_*` Insertion

**Status:** Deferred

### Motivation
`IrTypedModule.decls` is already topologically sorted by `ModuleChecker`. The only reason `bin/codegen_impl.h` still runs `codegen::topo_sort_decls` / `codegen::topo_sort_alias_decls` is that `record_anonymous_types_from_unit` appends newly-lifted `__anonym_<hash>` structs to the end of `unit->decls`, requiring a second sort (which also reverses independent zero-in-degree declarations due to its LIFO stack).

### Planned Fix
Insert each newly-lifted `__anonym_<hash>` `AliasDecl` immediately before the declaration that uses it during the single pass over `unit->decls` (and before `functions` / `trait_impls` for anonymous structs lifted from function/method signatures). This preserves `IrTypedModule`'s exact declaration order and removes `codegen::topo_sort_decls` and `codegen::topo_sort_alias_decls` completely.

---

## 9. `DefID`-Based Name Resolution & Strict Direct-Import Scoping

**Status:** Deferred

### Motivation
Currently, `ModuleChecker` computes `reachable_imports(file.imports)` (the transitive closure of imported modules) and loads all transitive `import_decls` into `file_ctx[F]`. Because Bapel uses flat string names in `ts::Context`, a file `F` that imports `bapel.stl` (which depends on `bapel.core`) can currently reference `bapel.core` symbols by name even without directly importing `bapel.core`.

### Planned Design Direction
Separate **Name Resolution** from **Typechecking** (as in Go and Rust):
1. **Per-file Lexical Scope (`string` $\to$ `DefID`):** Contains only local variables, module-local declarations, and **direct** `file.imports`. Source identifiers are resolved exclusively here and rewritten to canonical `DefID`s `(ModuleID, Name)`.
2. **Global Definition Store (`DefID` $\to$ `Definition`):** Holds all local, direct, and transitive definitions in `IrTypedModule`. Typechecking looks up resolved `DefID`s here, so transitive types referenced by exported APIs typecheck cleanly without leaking into a file's lexical name scope.

---

## 10. Revisit Trait Method Qualification Syntax (`Size::size` vs `Size.size`)

**Status:** Deferred

### Motivation
Currently, trait methods are qualified with `::` (e.g., `Size::size`), and `Context::add_symbol` flattens each `TraitDecl` method into an individual `TermDeclBind` named `Trait::method`.

### Open Design Question
Evaluate whether `Size::size` is the right syntax, or whether a projection-based syntax like `Size.size` (where `Size` is a trait and `.size` is a projection on that trait) is cleaner and more consistent with the rest of the language — and how that interacts with `TraitBind` vs synthesized `TermDeclBind`s in `ts/context.h`.

---

## 11. Operator Overloading via Traits

**Status:** Deferred

### Motivation
Currently, operators (`+`, `-`, `*`, `/`, `==`, etc.) are built-in polymorphic terms registered in `new_default_context()` (`comp/typecheck_module.h`).

### Planned Direction
Support type-directed operator overloading in `ts/inferencer.h` (e.g., rewriting `a + b` to `Add::add(&a, &b)` for custom types implementing `Add`, while keeping primitive arithmetic direct).

---

## 12. Implicit `bapel.core` Import & Runtime Headers

**Status:** Deferred

### Motivation
Certain core runtime types and helpers (e.g., `bapel/core_variant.h` for `bapel::Variant`, `Ptr`, `core::for`) are needed by generated C++ even in modules that do not explicitly `import bapel.core`. Currently `core_variant.h` is listed in `bapel/core.bpl`'s `impls`, which only exposes it to Bazel targets that import `bapel.core`.

### Planned Direction
Have the compiler implicitly import `bapel.core` (or attach a core runtime `cc_library` dependency automatically) for any module that does not explicitly import it.

---

## 13. Support Building Library Modules with `bpl build`

**Status:** Deferred

### Motivation
Currently, `bpl build <module>` (`bin/main.bpl`) always generates a Bazel `cc_binary` rule for the requested top-level module. Building a library module that does not define `main` (e.g., `./bpl build bapel.compiler.codegen.cpp`) compiles all C++ sources cleanly but fails at the final link step with `undefined reference to 'main'`.

### Planned Direction
Detect whether the top-level module defines a `main` function (e.g. via the queried/typechecked declarations):
- If `main` is present, emit a `cc_binary` rule and copy the resulting executable to `out/<module>`.
- If `main` is absent, emit a `cc_library` rule and build it with Bazel so library modules can be typechecked and compiled standalone without link errors.

---

## 14. Multiple Inherent `impl` Blocks for the Same Type Across Module Files

**Status:** Deferred

### Motivation
In large multi-file modules (such as `bapel.compiler.codegen.cpp` with `CppPrinter` split across `cpp_printer.bpl`, `cpp_type.bpl`, `cpp_term.bpl`, and `cpp_decl.bpl`), it is natural to define inherent `impl T { ... }` blocks for the same type `T` across multiple implementation files of the same module. Currently, each `impl T` emits its own `namespace inherents { struct T { ... }; }` in C++, causing a C++ redefinition error if `impl T` appears more than once.

### Planned Direction & Challenges
1. **Merging inherent `impl` blocks:** Group all inherent `IrTraitImpl` blocks targeting the same `(type_params, type_name)` within a module into a single C++ `inherents::T` struct (checking for duplicate method names across blocks).
2. **Declaration ordering / mutual recursion:** Because inherent `impl` methods are currently emitted inline in the header (`.h` / `_private.h`), methods split across files may call each other or private module functions. Either emit forward declarations for all methods inside `struct inherents::T` with out-of-line definitions (in `.cc` for non-generic methods or after all declarations in the header for generic methods), or topologically sort / split method declarations from definitions.


---

## 15. [Bug] Missing `return` in Unit-Returning (`-> ()`) Functions and Exhaustive `match`

**Status:** Open (currently worked around by adding explicit trailing `()` in affected functions).

### Problem

In Bapel, unit `()` compiles to `std::monostate`, which is a **non-void** C++ return type. In C++, falling off the end of a non-void function without executing a `return` statement is **Undefined Behavior** (`-Wreturn-type`).

Under `-O1`/`-O2`, Clang and GCC assume any fall-off path is unreachable, causing silent miscompilations and crashes:
1. **Tail `if` without `else` in a `-> ()` function:**
   ```bapel
   fn foo(cond: bool) -> () {
     if cond {
       do_something ();
     }
   }
   ```
   The generated C++ returns inside the `if` branch (or falls through both branches) with no `return std::monostate();` after the `if`. The optimizer assumes `cond` is always `true` and deletes the `if` check completely.
2. **Tail `for` loop in a `-> ()` function:**
   ```bapel
   fn bar(n: i64) -> () {
     let i: i64 = 0;
     for i < n {
       i <- i + 1;
     }
   }
   ```
   With no `return std::monostate();` after the `while` loop, the optimizer assumes the loop never terminates and omits `ret` (causing a segfault).
3. **Exhaustive `match` without a wildcard `_` arm:**
   Generated `switch (__v.index())` statements over `std::variant` have no `default:` case when all tags are covered, so C++ compilers emit `-Wreturn-type` warnings (they do not know `index()` is bounded to `0..N-1`).

### Proposed Fix (in `bapel/compiler/codegen/cpp_decl.bpl` & `cpp_term.bpl`)

1. **Unit-returning functions (`-> ()`):** Always emit `return std::monostate();` at the end of `-> ()` function, method, and lambda bodies (or whenever the tail expression is a statement such as `if`, `for`, or `match` that can fall through).
2. **Exhaustive `match` without `_`:** Emit `default: __builtin_unreachable();` (or `std::abort();`) in `switch (__v.index())` when no wildcard arm is present.


---

## 16. [Bug] Pointer Dereference + Field Projection (`(*p).field`) Copies the Struct Instead of Projecting In-Place

**Status:** Open (currently worked around using `SharedPtr` fields and `Ptr::set (p, set (*p) { ... })`).

### Problem

In Bapel, dereferencing a pointer (`*p`) lowers to `Ptr::get(p)`, which returns `T` **by value** (a full copy of `*p`). As a result:
1. **Field mutation via method call silently mutates a temporary:** `(*self).out.append s` copies `*self`, appends to the temporary copy's `out`, and discards it, leaving the original `self->out` unchanged.
2. **Taking a pointer to a field (`&((*p).field)`) points to a temporary** rather than to the field inside `*p`.
3. **Updating a single field requires copying the entire struct:** `Ptr::set (p, set (*p) { field = val })` instead of in-place `(*p).field <- val` (or `p.field <- val` with auto-deref).

### Motivation / Context (`CppPrinter.out`)

In `bapel/compiler/codegen/cpp_printer.bpl`, `CppPrinter.out` had to be typed as `SharedPtr String` instead of `String` solely because `(*self).out.append s` would otherwise mutate a temporary copy, and `Ptr::set (self, set (*self) { out = ... })` on every write would copy the entire output buffer ($O(N^2)$).

Once this bug is fixed, `CppPrinter` should be simplified to:
- Change `CppPrinter.out` from `SharedPtr String` back to `String`.
- Replace `Ptr::set (p, set (*p) { ... })` updates throughout `bapel/compiler/codegen/cpp_*.bpl` with direct in-place field assignments.

### Proposed Fix

- Treat dereference + projection (`(*p).field`, `&((*p).field)`, `(*p).field.method(...)`, `(*p).field <- val`) as an lvalue projection that lowers to `p->field` (or `(&p->field)`) in C++ without calling `Ptr::get(p)` by value.


---

## 17. [Bug] Cross-Module Construction of `variant` Tags with Anonymous-Struct Payloads

**Status:** Open (currently worked around by defining constructor methods inside the defining module, e.g. `IrType::new_*` and `IrDecl::new_*` in `bapel.compiler.ir`).

### Problem

When a `pub type` in module `A` defines a `variant` whose tag payload is an anonymous struct (e.g. `app_t struct { fun: SharedPtr 't, arg: SharedPtr 't }` in `IrTypeData`), constructing that variant tag from an importing module `B` (e.g. `variant { (IrTypeData IrType) app_t = struct { fun = ..., arg = ... } }`) fails or produces mismatched C++ anonymous struct types (`__anonym_<hash>`).

### Context / Workaround

In `bapel.compiler.ir`, we added constructor methods (`IrType::new_app`, `IrType::new_fun`, `IrDecl::new_alias`, etc.) so `bapel.compiler.codegen.cpp` can construct `IrType` and `IrDecl` values without instantiating anonymous-struct variant payloads across module boundaries.

### Proposed Fix

Ensure anonymous-struct payloads inside exported `variant` definitions are either:
1. Resolved to the exact same canonical `__anonym_<hash>` type across modules (including identical `#ifndef __anonym_<hash>_DEFINED` guards and identical type-var parameterizations), or
2. Constructed via the variant tag's named payload type / generated helper without requiring the caller module to re-synthesize a matching anonymous struct.
