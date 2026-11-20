# Modules

A module is the unit of abstraction. It is a collection of source
files that share the management of a set of symbols.
                                     
A module decides which symbols exist in that module, and which are visible or
hidden to the other modules.

## Related designs

* [Module discovery](module-discovery.md)

## Source files

A module is a collection of source files.

A source file can be a base file or an implementation file.

A module must have exactly 1 base file, and it can have 0 or more implementation
files.

### Roles of base and implementation files

The base file and all `.bpl` implementation files of a module have **equal
weight**. The base file's only special role is declaring the module's identity
(`module $MODULE_ID`), listing its implementation files (`impls { ... }`), and
specifying module-level build flags (`flags { ... }`).
All declarations and definitions—`type`, `fn`, `trait`, `impl Trait for T`, and
inherent `impl T` blocks—behave identically regardless of whether they appear in
the base file or in an implementation file. Each source file (base or
implementation) still has its own `imports` section, because imported symbols
are only visible within the file that imports them.

## Module identifiers

A module identifier (`$MODULE_ID`) is a collection of identifier components
separated by `.`, e.g., `mymodule`, `main`, `bapel.core`,
`bapel.internal.utils`, etc.

A source file filename (`$MODULE_FILENAME`) is, e.g., `"myfile.bpl"`,
`"myfile.h"`, `"main.bpl"`, `"main_impl.bpl"`, `"bapel/core.bpl"`,
`"bapel/core_impl.h"`, etc.

## `module` clause

A base file must begin with the `module` clause on the first line.

The `module` clause is `module $MODULE_ID`, e.g., `module bapel.core`.

The `module` clause is what identifies a file as being a base file, and the
resulting module as being an actual module.

## `implements` clause

An implementation file must begin with the `implements` clause on the first
line.

The `implements` clause is `implements $MODULE_ID`, e.g., `implements bapel.core`.

The `implements` clause identifies that a file is an implementation file, and it
must be consistent with the `impls` section (see below).

## `imports` section

A source file (base and implementation) can have at most one `imports` section
with 1 or more module identifiers (`$MODULE_ID`), e.g.:

    imports {
      core
      vec
    }

The module identifiers in the `imports` section must be distinct and must be
lexicographically sorted. Otherwise, it is an error.

The module names given in the `imports` section are "imported" into the current
module. What "import" means is that the symbols exported by the imported modules
are made visible to the current module, to be used. That's all it means.

Importing some other module does not mean to copy its module text into the
current module.

## `impls` section

A base file can have at most one `impls` section with 1 or more relative
filenames (`$RELATIVE_FILENAME`, relative to the base file's directory), e.g.:

    impls {
      "myfile1.bpl"
      "myfile2.h"
    }

An implementation file cannot have an `impls` section.

The `impls` section and the `implements` clauses must be consistent, i.e., the
base file must declare all its implementation files and all the implementation
files must declare the same `$MODULE_ID`. Otherwise, it is an error.

TODO: Does the `impls` section need to be sorted? See
https://github.com/jabolopes/bapel/issues/7.

## `flags` section

A base file can have at most one `flags` section specifying module-level build
or linker flags (e.g., C++ compiler/linker options passed to the build system).
An implementation file cannot have a `flags` section, because build flags apply
to the module as a whole.

## Symbol management

A module manages symbols. All files belonging to the same module (its base file
and all its implementation files) share the module's symbols directly.

Symbols visible within a source file:
* can be imported from other modules (visible only within the file that imports
  them).
* can be declared in the module (in the base file or any implementation file).
* can be defined in the module (in the base file or any implementation file).
* can be both declared & defined in the module.

Symbols can also be exported (`pub`).

## Symbol frequency

* A symbol can be imported from other modules at most once per source file.
* Across the entire module (base file + all implementation files), a symbol can
  be declared at most once and defined at most once.
* A trait implementation (`impl Trait for T`) can be defined at most once per
  `(Trait, T)` pair in `T`'s module.
* An inherent method (`T::method`) can be defined at most once per `(T, method)`
  pair in `T`'s module.

## Symbol shadowing

Shadowing follows the scope hierarchy: **inner block / function scope**
(parameters and `let` bindings) **> current module** (base and implementation
files) **> imported modules**.

In particular, symbols imported from other modules can be shadowed by symbols
declared or defined in the current module. The reason for this is to avoid the
following problem: let's say our module M imports some module C, and module C
happened to add a new function that module M already defined. Without shadowing,
module M would no longer compile. In other words, any module should be allowed
to add new symbols without breaking other modules that depend on it.

## All declared terms must be defined

All terms declared by a module must be defined within that module. In other
words, given a base file and 0 or more implementation files that belong to the
same module, all declared terms must also be defined in that module.

For example, if module `A` has:

    x: () -> ()

then there must be a definition of `x` in `A`'s base file or one of its
implementation files. A definition could be for example:

    fn x() -> () { ... }

Allowing a term's declaration and definition to live in separate files within
the same module aligns with the equal-weight module model, gives programmers
flexibility in how they organize code, and enables patterns such as
platform-specific implementation files (e.g., `*_linux.bpl`, `*_macos.bpl`,
once conditional file selection in `impls` is supported in the future).

## Type abstraction

A module can export a type (declaration), e.g.:

    pub type T

A module can also export a type (declaration) and have a hidden type definition,
e.g.:

    pub type T

    type T = (i8, i8)

In both cases, the type `T` is exported but its internal representation is
hidden from any modules that import it, but still visible to implementation
files in the same module where it is declared / defined.

A module can also export a type definition, e.g.:

    pub type T = (i8, i8)

In this case, type `T` is a tuple of `i8`s and it is exported. Type `T`'s
internal representation is made visible to any modules that import it.

The options above differ because the first 2 options hide the type's internal
representation, whereas the last option makes the internal representation
visible. This is how to employ abstraction.

A type declaration and its type definition must always be consistent (matching
kinds), and a type may have at most one declaration and at most one definition
across the module. They may appear in any order (including in separate files of
the same module).

## Term abstraction

A module exports a term with `pub fn`:

    pub fn f() -> () { () }

Only the term's type signature is exposed in the module's public interface
(`.h`); non-template function bodies live in `.cc` to preserve encapsulation and
incremental builds. (Generic/template functions emit their bodies in `.h` as
required by C++ template instantiation.)

Standalone `decl` declarations (e.g., `pub decl f: () -> ()`) are also supported
for C++ `@bpl` FFI declarations, self-contained typechecking tests without
imports, or separating a declaration from its definition across files. When both
a declaration and a definition exist, there may be at most one of each across
the module, in any order, with matching types.

## Traits and `impl` blocks

A module can declare and export traits (`pub trait Size { ... }`) and define
trait implementations (`impl Size for T { ... }`) and inherent implementations
(`impl T { ... }`) for any type defined in the same module (see
[`designs/traits.md`](traits.md)). When another module imports this module, all
exported traits, trait implementations for exported types, and inherent methods
for exported types are visible to the importing module.

# Imports

A source file can import from other modules by importing their module identifier.

For example, module `core` has:

    pub type T1

    pub type T2 = (i8, i8)

    pub fn f() -> () { () }

    pub trait Size {
      fn size(s: &Self) -> i64
    }

    impl Size for T2 {
      fn size(t: &Self) -> i64 { 2 }
    }

    impl T2 {
      fn first(t: &Self) -> i8 { (*t).0 }
    }

The type `T1` is exported but its internal representation is hidden.

The type `T2` is exported as well as its internal representation.

The term `f` is exported with its type signature.

The trait `Size`, the trait implementation `impl Size for T2`, and the inherent
method `T2::first` are also exported because `Size` and `T2` are exported.

For example, base file `main.bpl` has:

    imports {
      core
    }

This is normalized to:

    importDecls {
      /* core:... */ pub type T1
      /* core:... */ pub type T2 = (i8, i8)
      /* core:... */ pub trait Size { fn size(s: Ptr Self) -> i64 }
      /* core:... */ pub f: () -> ()
    }

    importedTraitImpls {
      /* core:... */ impl Size for T2 { fn size(t: Ptr Self) -> i64 }
      /* core:... */ impl T2 { fn first(t: Ptr Self) -> i8 }
    }

As a result, base file `main.bpl` can:
* use type `T1` in an abstracted way, without referring to its internal
  representation.
* use type `T2` including its internal representation, e.g., construct and
  destruct the tuple type.
* use term `f` via its exported type signature.
* use trait `Size` (e.g., as a generic bound `'t: Size` or via `Size::size`),
  call `t.size()` / `Size::size (&t)` on `T2`, and call inherent method
  `t.first()` / `T2::first (&t)`.

## Implicit declarations

All functions defined in a module are also declared, either explicitly or
implicitly.

An explicit declaration is when the source program contains a standalone
declaration (e.g., `decl f: () -> ()`).

An implicit declaration is a `TermDecl` that is automatically inserted when a
function (`fn f() -> () { () }`) is defined without an explicit `decl`. (Type
definitions like `type T2 = (i8, i8)` are already declarations in `decls` and do
not need a separate implicit `NameDecl`.)

For example, `main.bpl` has:

    type T1 = (i8, i8)

    decl f1: () -> ()
    fn f1() -> () { () }

    fn f2() -> () { () }

This is normalized so that `decls` contains:

    type T1 = (i8, i8)
    f1: () -> ()
    f2: () -> ()

The result is:
* term `f1` was already explicitly declared, so it has no implicit declaration.
* term `f2` was not explicitly declared, so an implicit declaration `f2: () -> ()`
  is inserted.

Implicit declarations obey the rule that symbols are declared / defined at most
once across the module.

The main purpose of implicit declarations is to allow functions to call each
other in any order (including mutually recursive functions and functions defined
in different files of the same module) without requiring explicit forward
declarations.

All declarations (types, traits, and implicit/explicit term declarations) are
automatically topologically sorted based on their dependencies.

Symbols imported from other modules are considered explicit declarations.

---

# Normalized `IrModule` Format

For `NormalizeTest`, a normalized module (output of parser + desugar +
resolver, **before** typechecking) preserves the **per-file structure** so the
typechecker knows each file's exact `imports` scope without relying on
`Pos::filename` lookups:

```text
module my_module

imports {
  bapel.core "bapel/core.bpl"
  bapel.stl "bapel/stl.bpl"
}

impls {
  "path/to/my_module_impl.bpl"
  "path/to/my_module_ffi.h"
}

flags {
  ...
}

importDecls {
  /* bapel.core:1 */ pub type Ptr :: * -> *
  ...
}

importedTraitImpls {
  /* bapel.stl:40 */ impl String { ... }
  /* bapel.stl:85 */ impl Size for String { ... }
  ...
}

file "path/to/my_module.bpl" (base) {
  imports {
    /* path/to/my_module.bpl:4 */ bapel.core
  }
  decls {
    /* path/to/my_module.bpl:10 */ pub type Foo = struct{x: i64}
  }
}

file "path/to/my_module_impl.bpl" (impl) {
  imports {
    /* path/to/my_module_impl.bpl:3 */ bapel.core
    /* path/to/my_module_impl.bpl:4 */ bapel.stl
  }
  decls {
    /* path/to/my_module_impl.bpl:6 */ pub foo_fn: Foo -> i64
  }
  traitImpls {
    /* path/to/my_module_impl.bpl:12 */ impl Foo { ... }
    /* path/to/my_module_impl.bpl:18 */ impl Size for Foo { ... }
  }
  functions {
    /* path/to/my_module_impl.bpl:6 */ pub fn foo_fn(f: Foo) -> i64 { ... }
  }
}

file "path/to/my_module_ffi.h" (ffi) {
  decls {
    /* path/to/my_module_ffi.h:1 */ pub ffi_helper: () -> ()
  }
}
```

Key properties:
* **Uniform per-file grouping:** Every file in the module (`base`, `impl`, or
  C++ `.h` `ffi`) has its own `file` block with its `imports`, `decls`,
  `traitImpls`, and `functions`.
* **Deterministic file ordering:** The `base` file is always first in
  `module.files`, followed by the implementation files (`impl` and `ffi`) in
  the exact order they appear in the base file's `impls` section.
* **Cached module-level imports:** `importDecls` and `importedTraitImpls` are
  loaded once per unique `ModuleID` across all files in the module.
* **All module resolution happens in the normalizer:** A core design goal is
  that module discovery (`workspace.bpl`, `ModuleFinder`, `Querier`) is
  completed entirely during normalization. By recording resolved paths in
  `imports` and `impls` and populating `importDecls` / `importedTraitImpls`,
  all downstream stages (typechecking, codegen, etc.) operate purely on
  `IrModule` / `IrUnit` and filenames without ever performing module
  resolution.
* **Workspace-relative filenames:** While the source `impls` section uses paths
  relative to the base file's directory, all filenames in `IrModule` (resolved
  `imports`, `impls`, and `file "..."`) are resolved paths (workspace-root-relative
  within the same workspace).

---

# Typechecked `IrModule` Format

`UnitChecker` validates and elaborates the `IrModule` **in-place**, preserving
the exact same `IrModule` structure and per-file grouping (no separate `IrUnit`
type). Every `file` is checked under its own `imports` scope plus the shared
module-local declarations across all `files`, and all inferred types, type
applications, method resolutions, and variant tag indices are filled in:

```text
module my_module

imports {
  bapel.core "bapel/core.bpl"
  bapel.stl "bapel/stl.bpl"
}

impls {
  "path/to/my_module_impl.bpl"
  "path/to/my_module_ffi.h"
}

flags {
  ...
}

importDecls {
  /* bapel.core:1 */ pub type Ptr :: * -> *
  ...
}

importedTraitImpls {
  /* bapel.stl:40 */ impl String { ... }
  /* bapel.stl:85 */ impl Size for String { ... }
  ...
}

file "path/to/my_module.bpl" (base) {
  imports {
    /* path/to/my_module.bpl:4 */ bapel.core
  }
  decls {
    /* path/to/my_module.bpl:10 */ pub type Foo = struct{x: i64}
  }
}

file "path/to/my_module_impl.bpl" (impl) {
  imports {
    /* path/to/my_module_impl.bpl:3 */ bapel.core
    /* path/to/my_module_impl.bpl:4 */ bapel.stl
  }
  decls {
    /* path/to/my_module_impl.bpl:6 */ pub foo_fn: Foo -> i64
  }
  traitImpls {
    /* path/to/my_module_impl.bpl:12 */ impl Foo { ... }
    /* path/to/my_module_impl.bpl:18 */ impl Size for Foo { ... }
  }
  functions {
    /* path/to/my_module_impl.bpl:6 */ pub fn foo_fn(f: Foo) -> i64 { (f:Foo.x):i64 }:i64
  }
}

file "path/to/my_module_ffi.h" (ffi) {
  decls {
    /* path/to/my_module_ffi.h:1 */ pub ffi_helper: () -> ()
  }
}
```

Key properties:
* **Single IR type across stages:** Both normalization and typechecking operate
  on `IrModule` (`NormalizeTest`, `InferTest`, and `TypecheckTest` share one
  format).
* **In-place elaboration:** All inferred types, type applications, method
  resolutions, and variant tag indices are populated inside each `file`.
* **Direct codegen input:** `CppPrinter` iterates over `module.files`—`ffi`
  files contribute `#include` directives (without re-declaring their `decls`),
  while `base` and `impl` files emit all `decls`, `traitImpls`, and `functions`
  into one `<module>.h`, `<module>_private.h`, and `<module>.cc`.
* **Empty sections omitted:** Sections with no entries are omitted.

---

# Compilation Model: Whole-Module Unit

To guarantee equal weight between a module's base file and its `.bpl`
implementation files without ad-hoc per-file special cases, Bapel compiles a
module as a **single `IrModule`** (analogous to Go packages and Rust crates),
while preserving **per-file import scoping**:

1. **Parse & assemble `IrModule`:** Given a module's base file (`module M`),
   parse the base file, all `.bpl` files, and all C++ `.h` FFI files listed in
   its `impls { ... }` section into `module.files` (each `file` keeping its own
   `filename`, kind `base | impl | ffi`, `imports`, `decls`, `trait_impls`, and
   `functions`).
2. **Load each imported module once (cached):** Collect the deduplicated union
   of imported `ModuleID`s across `module.files`, resolve each `ModuleID` to
   its path, and load its exported `importDecls` and `importedTraitImpls`
   **once** into `IrModule`, avoiding redundant parsing and desugaring of
   shared dependencies like `bapel.core` and `bapel.stl`.
3. **Typecheck in-place (`ModuleChecker`):** Validate and elaborate every
   `file` in `module.files` under its own `imports` scope plus the shared
   module-local declarations (see `# Typechecking an IrModule` below).
   Because all `.bpl` files are typechecked together in one `IrModule`, `type`,
   `fn`, `trait`, and `impl` blocks (including generic functions and generic
   `impl` blocks) behave identically regardless of which `.bpl` file they live
   in.
4. **Emit C++17:** Generate one `<module>.h`, `<module>_private.h`, and
   `<module>.cc` per module, `#include`ing the deduplicated union of imported
   module headers and `ffi` `.h` files.

---

# Typechecking an `IrModule`

`ModuleChecker` typechecks an `IrModule` in five steps, enforcing both
**module-wide symbol sharing** and **per-file import scoping**:

1. **Index cached imports by `ModuleID`:**
   Group `module.import_decls` and `module.imported_trait_impls` by their
   defining `ModuleID`. Because the normalizer already loaded each unique
   imported module once, building per-file import contexts requires no
   re-parsing.
2. **Collect & topologically sort module-local declarations:**
   * Gather all `decls` across all `files` (`base`, `impl`, and `ffi`).
   * Verify at most one declaration per symbol across the module.
   * Run `topo_sort_decls` over the combined module declarations.
   * Record all local type names (`NameDecl` and `AliasDecl`) for trait
     coherence.
3. **Build per-file checking contexts (`file_ctx[F]`):**
   For each `.bpl` file `F` (`base` or `impl`):
   * Start from `default_context`.
   * Add `import_decls` and `imported_trait_impls` for **only** the `ModuleID`s
     listed in `F.imports` (rejecting duplicate imported symbols within `F`).
   * Call `enter_scope()` so module-local symbols shadow imported symbols.
   * Add all topologically sorted module-local `decls` (`base`, `impl`, `ffi`)
     and register all module-local `trait_impls`.
   * Validate `is_wellformed_decl` on each declaration in `F.decls` using
     `file_ctx[F]` (ensuring a type/term declaration in `F` only references
     imports visible in `F`).
4. **Infer & verify bodies per file (in-place):**
   For each `.bpl` file `F`, using `file_ctx[F]`:
   * For each function `fn` in `F.functions`: verify at most one definition
     across the module, then run `infer_function` and `typecheck_function`,
     elaborating `fn` in-place.
   * For each `impl` block `ti` in `F.trait_impls`: run `check_trait_impl`
     (verifying coherence against module-local types and elaborating all method
     bodies in-place).
5. **Check undefined terms across the module:**
   Verify that every non-FFI `TermDecl` (declared in any `base` or `impl` file)
   has a matching definition in some `base` or `impl` file of the module.

---

# Execution Plan

1. **Core IR (`bin/ir_unit.h`, `bin/ir_base.h`):**
   * Define `IrFileCase` (`Base`, `Impl`, `Ffi`), `IrFile`, and `IrModule`
     (with `module_id`, `imports`, `impls`, `flags`, `import_decls`,
     `imported_trait_impls`, and `files`), plus `IrModule::to_bpl_string()`
     matching `# Normalized IrModule Format` / `# Typechecked IrModule Format`
     (printing `pub` instead of `export`).
2. **Normalizer & Resolver (`comp/resolver.h`, `comp/normalizer.h`, `comp/desugar.h`):**
   * Assemble `IrModule` by parsing the base file and all `.bpl` / `.h` files
     in `impls { ... }` into `module.files`.
   * Resolve all `ModuleID`s in `imports` once per unique `ModuleID`, populate
     `module.imports` (with resolved paths), `import_decls`, and
     `imported_trait_impls`.
   * Run `desugar_self_recursive` over each `.bpl` file in `module.files`.
3. **Module Typechecker (`comp/typecheck_unit.h`):**
   * Implement the 5-step `ModuleChecker` algorithm over `IrModule` (cached
     imports by `ModuleID`, combined `topo_sort_decls`, per-file `file_ctx[F]`,
     in-place elaboration of `functions` & `trait_impls`, module-wide
     undefined-term check).
4. **C++17 Codegen (`bin/codegen_impl.h`):**
   * Update `CppPrinter` to consume `IrModule`: collect `decls` / `trait_impls`
     / `functions` from `base` + `impl` files, lift anonymous structs,
     topo-sort once, emit `#include`s for `imports` and `ffi` files, and output
     one `<module>.h`, `<module>_private.h`, and `<module>.cc`.
   * Remove temporary `print_trait_impl_out_of_line` workaround.
5. **CLI Driver (`bin/main.bpl`) & Bapel IR (`bapel/compiler/ir_decl.bpl`):**
   * Update `bin/main.bpl` so each module only invokes `codegen::compile_unit`
     on its base file and `BUILD` `srcs` lists `<module>.cc`.
   * Update `bapel/compiler/ir_decl.bpl` to mirror `IrFile` / `IrModule`.

---

# Verification Plan

1. **Golden Tests (`make test`):**
   * Regenerate `NormalizeTest`, `InferTest`, `TypecheckTest`, and
     `CppPrinterTest` golden files for the new `IrModule` format and verify all
     22+ test suites pass (including `CppPrinterTest.IsValidCpp`).
2. **Multi-File Module Tests:**
   * Add a multi-file test verifying equal-weight `type`, `fn`, `trait`, and
     `impl` blocks across base and `.bpl` `impls` files, plus per-file import
     scoping (a symbol imported only by file B is not visible in file A).
   * Compile `bapel/compiler/ir.bpl` (base + 5 `.bpl` `impls` files with
     inherent `impl` blocks) as a single whole-module `IrModule` and verify
     `clang++ -std=c++17 -fsyntax-only` passes with 0 errors.
3. **Bootstrap:**
   * Promote seed compiler (`make bpl && make bootstrap`) and verify
     `make program && make query`.
