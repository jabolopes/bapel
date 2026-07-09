# Plan: Porting `query` to Bapel (`bapel.query`)

## 1. Overview & Scope

The Goal is to port the Go `query` package (approx. 400 LOC across 5 files in [query/](file:///usr/local/google/home/jabolopes/.gemini/jetski/scratch/bapel/query)) into a native Bapel library module named `bapel.query`, located at `bapel/query.bpl`.

Currently, the CLI driver [bin/main.bpl](file:///usr/local/google/home/jabolopes/.gemini/jetski/scratch/bapel/bin/main.bpl) implements an ad-hoc version of module resolution (`resolveModule`, `PackageMapping`, `findBestMatch`) and relies on subprocess calls to `bootstrap/parser -format=flat` for dependency discovery. Porting `query` to Bapel will formalize these data structures into a reusable standard module, unify workspace path resolution, and lay the foundation for in-process compilation driving.

---

## 2. Proposed Module Architecture & Data Types

We will create `bapel/query.bpl`, importing `bapel.core`, `bapel.stl` (for `Vector`, `String`, `fs`, `IStringStream`), and `bapel.os`.

### Data Structures in `bapel.query`

```bapel
module query

imports {
  bapel.core
  bapel.os
  bapel.stl
}

// Represents a workspace package mapping (from module_finder.go and main.bpl)
pub type PackageMapping = struct {
  is_prefix: bool,
  name: String,
  path: String
}

// Encapsulates module lookup tables
pub type ModuleFinder = struct {
  modules_by_name: UnorderedMap String String,
  modules_by_prefix: UnorderedMap String String
}

// Results of querying a single source file (from source_file_query.go)
pub type SourceFileQuery = struct {
  imports: Vector String,
  impls: Vector String,
  flags: Vector String,
  decls: Vector String,
  trait_impls: Vector String
}

// Results of querying a full module and its implementation files (from module_query.go)
pub type ModuleQuery = struct {
  imports: Vector String,
  impls: Vector String,
  flags: Vector String,
  decls: Vector String,
  trait_impls: Vector String
}
```

---

## 3. Function Mapping (Go to Bapel)

The following table maps existing Go functions to their proposed Bapel implementations:

| Go Source File | Go Function / Method | Proposed Bapel API (`bapel.query`) | Notes & Implementation Strategy |
| :--- | :--- | :--- | :--- |
| [module_finder.go](file:///usr/local/google/home/jabolopes/.gemini/jetski/scratch/bapel/query/module_finder.go) | `newModuleFinder` | `pub fn ModuleFinder::mk() -> ModuleFinder` | Reads `workspace.bpl` (or default paths) using `fs::exists` and parses mappings. |
| [module_finder.go](file:///usr/local/google/home/jabolopes/.gemini/jetski/scratch/bapel/query/module_finder.go) | `lookupModuleByName` / `ByPrefix` | `fn ModuleFinder::lookup(finder: &ModuleFinder, mod_id: &String) -> (bool, String)` | Unifies and refines `findBestMatch` currently in [bin/main.bpl](file:///usr/local/google/home/jabolopes/.gemini/jetski/scratch/bapel/bin/main.bpl#L91). |
| [module_finder.go](file:///usr/local/google/home/jabolopes/.gemini/jetski/scratch/bapel/query/module_finder.go) | `baseSourceFilename` | `pub fn ModuleFinder::base_filename(finder: &ModuleFinder, mod_id: &String) -> String` | Replaces `resolveModule` in `main.bpl`. Uses `fs::join` and string separator replacement. |
| [module_finder.go](file:///usr/local/google/home/jabolopes/.gemini/jetski/scratch/bapel/query/module_finder.go) | `implSourceFilename` | `pub fn ModuleFinder::impl_filename(base_file: &String, rel_impl: &String) -> String` | Computes `fs::join (fs::parent_path (*base_file), *rel_impl)`. |
| [query_source_file.go](file:///usr/local/google/home/jabolopes/.gemini/jetski/scratch/bapel/query/query_source_file.go) | `queryAnnotationNonBplFile` | `fn query_annotation_file(path: &String) -> SourceFileQuery` | Reads file line-by-line using `IStringStream` / `getline`. Scans for `import ` prefixes and `// @bpl: ` annotation strings. |
| [query_source_file.go](file:///usr/local/google/home/jabolopes/.gemini/jetski/scratch/bapel/query/query_source_file.go) | `queryDeclsBplFile` | `fn query_bpl_file(path: &String) -> SourceFileQuery` | For MVP, invokes `bootstrap/parser -format=flat <path>` and extracts imports, impls, and decls from flat text (extending `parseSourceFileFlat`). |
| [query_source_file.go](file:///usr/local/google/home/jabolopes/.gemini/jetski/scratch/bapel/query/query_source_file.go) | `QuerySourceFile` | `pub fn query_source_file(path: &String) -> SourceFileQuery` | Dispatches to `query_bpl_file` if extension is `.bpl`, otherwise `query_annotation_file`. |
| [querier.go](file:///usr/local/google/home/jabolopes/.gemini/jetski/scratch/bapel/query/querier.go) | `QueryModule` | `pub fn query_module(finder: &ModuleFinder, mod_id: &String) -> ModuleQuery` | Queries base filename, iterates over `impls` to query implementation files, merges vectors, deduplicates, and sorts. |
| [querier.go](file:///usr/local/google/home/jabolopes/.gemini/jetski/scratch/bapel/query/querier.go) | `QueryModuleExports` | `pub fn query_module_exports(finder: &ModuleFinder, mod_id: &String) -> ModuleQuery` | Calls `query_module` and filters `decls` for exported symbols (or flag prefix in flat format). |

---

## 4. Implementation Phases

### Phase 0: Standard Library & Compiler Prerequisites (COMPLETED)
Before implementing `query`, two minor prerequisites were addressed:
1. **File Input Reading (`Ifstream`) (COMPLETED):** Added `pub type Ifstream` and an `impl` block in [bapel/stl.bpl](file:///usr/local/google/home/jabolopes/.gemini/jetski/scratch/bapel/bapel/stl.bpl#L148-L158), backed by `IfstreamImpl` in [bapel/stl_fstream.h](file:///usr/local/google/home/jabolopes/.gemini/jetski/scratch/bapel/bapel/stl_fstream.h#L36-L57). Both `Ofstream::open` and `Ifstream::open` take reference arguments (`&String`). [getline](file:///usr/local/google/home/jabolopes/.gemini/jetski/scratch/bapel/bapel/stl_string.h#L14) automatically works with `Ifstream` since it is polymorphic over stream types.
2. **Generic Methods in `impl` Blocks (COMPLETED):** Updated [comp/cpp_printer.go](file:///usr/local/google/home/jabolopes/.gemini/jetski/scratch/bapel/comp/cpp_printer.go#L1233-L1277) to emit C++ template parameters for generic methods inside non-generic `impl` blocks (enabling polymorphic methods like `Ifstream::read`).
3. **String Utilities (COMPLETED):** Added `ends_with`, `remove_prefix`, `remove_suffix`, `trim_prefix`, and `trim_suffix` to `StringView` and `String` in [bapel/stl.bpl](file:///usr/local/google/home/jabolopes/.gemini/jetski/scratch/bapel/bapel/stl.bpl#L25-L125) and [bapel/stl_string.h](file:///usr/local/google/home/jabolopes/.gemini/jetski/scratch/bapel/bapel/stl_string.h#L68-L129).
4. **Hash Map (`UnorderedMap`) (COMPLETED):** Added `pub type UnorderedMap ['k, 'v]` in [bapel/stl.bpl](file:///usr/local/google/home/jabolopes/.gemini/jetski/scratch/bapel/bapel/stl.bpl), backed by `std::unordered_map` in [bapel/stl_unordered_map.h](file:///usr/local/google/home/jabolopes/.gemini/jetski/scratch/bapel/bapel/stl_unordered_map.h), including `mk`, `insert`, `size`, `empty`, `contains`, and `get`.

### Phase 1: Module Finder & Workspace Resolution
1. Create `bapel/query.bpl`.
2. Port package mapping logic by migrating `parseWorkspaceFlat` and populating `modules_by_name` and `modules_by_prefix` (`UnorderedMap String String`) in `ModuleFinder`.
3. Implement `ModuleFinder::lookup`, `ModuleFinder::base_filename` and `ModuleFinder::impl_filename`.

### Phase 2: Source File Querying (`query_source_file`)
1. Implement `query_annotation_file(&String)`:
   - Read header/source files line-by-line using `IStringStream` and `getline`.
   - Check for `import <mod>;` statements and extract module names.
   - Check for `// @bpl: ` prefixes and extract embedded Bapel declarations.
2. Implement `query_bpl_file(&String)`:
   - Adapt `parseSourceFileFlat` from `main.bpl` to populate `SourceFileQuery` (including decls and flags).
3. Implement the unified `query_source_file(&String)` entrypoint.

### Phase 3: Module Querier & Deduping (`query_module`)
1. Implement vector helper functions in `bapel/query.bpl` (or leverage `bapel.stl`):
   - `merge_unique_strings(dst: &Vector String, src: &Vector String)`
   - Sorting and compaction for module IDs and filenames.
2. Implement `query_module` and `query_module_exports`.

### Phase 4: Driver Integration ([bin/main.bpl](file:///usr/local/google/home/jabolopes/.gemini/jetski/scratch/bapel/bin/main.bpl))
1. Update `bin/main.bpl` to import `bapel.query`.
2. Replace ad-hoc functions (`PackageMapping`, `resolveModule`, `parseSourceFileFlat`, `collectImplImports`) with calls to `bapel.query` functions (`ModuleFinder::mk()`, `ModuleFinder::base_filename()`, `query_source_file()`).
3. Verify clean compilation of `main.bpl`.

### Phase 5: Verification & Testing
1. Add a test target in [Makefile](file:///usr/local/google/home/jabolopes/.gemini/jetski/scratch/bapel/Makefile) or create a new test program `tests/query_test.bpl`.
2. Execute queries against `bapel/core` and `./bapel/core_impl.h` using both `bootstrap/querier` (Go) and the new Bapel implementation, asserting identical output.

---

## 6. Verification Strategy

To ensure parity between the Go querier and `bapel.query`:
1. Run `./bpl query bapel/core` (which tests the existing Go querier).
2. Run a new Bapel test binary that invokes `query_module(&finder, &"bapel.core".to_string)` and prints formatted results.
3. Diff the output of both tools across multiple target modules in the repo (`bapel.core`, `bapel.stl`, `bapel.os`).
