# Plan: Porting `query` to Bapel (`bapel.query`)

## 1. Overview & Scope

The Goal is to port the Go `query` package (approx. 400 LOC across 5 files in [query/](file:///usr/local/google/home/jabolopes/.gemini/jetski/scratch/bapel/query)) into Bapel, located at [bin/query.bpl](file:///usr/local/google/home/jabolopes/.gemini/jetski/scratch/bapel/bin/query.bpl) as part of the `bin.main` module.

Currently, the CLI driver [bin/main.bpl](file:///usr/local/google/home/jabolopes/.gemini/jetski/scratch/bapel/bin/main.bpl) implements an ad-hoc version of module resolution (`resolveModule`, `PackageMapping`, `findBestMatch`) and relies on subprocess calls to `bootstrap/parser -format=flat` for dependency discovery. Porting `query` to Bapel will formalize these data structures into a reusable standard module, unify workspace path resolution, and lay the foundation for in-process compilation driving.

---

## 2. Proposed Module Architecture & Data Types

We create [bin/query.bpl](file:///usr/local/google/home/jabolopes/.gemini/jetski/scratch/bapel/bin/query.bpl) as an implementation file (`implements bin.main`) of [bin/main.bpl](file:///usr/local/google/home/jabolopes/.gemini/jetski/scratch/bapel/bin/main.bpl). This allows sharing types and functions directly without namespace fragmentation or duplicate ad-hoc symbol definitions.

### Data Structures in `bin/query.bpl`

```bapel
implements bin.main

// Represents a workspace package mapping (from module_finder.go and main.bpl)
type PackageMapping = struct {
  is_prefix: bool,
  name: String,
  path: String
}

// Encapsulates module lookup tables
type ModuleFinder = struct {
  modules_by_name: UnorderedMap String String,
  modules_by_prefix: UnorderedMap String String
}

// Results of querying a single source file (from source_file_query.go)
type SourceFileQuery = struct {
  import_modules: Vector String,
  impl_files: Vector String,
  flag_files: Vector String,
  declarations: Vector String,
  trait_implementations: Vector String
}

// Results of querying a full module and its implementation files (from module_query.go)
type ModuleQuery = struct {
  import_modules: Vector String,
  impl_files: Vector String,
  flag_files: Vector String,
  declarations: Vector String,
  trait_implementations: Vector String
}
```

---

## 3. Function Mapping (Go to Bapel)

The following table maps existing Go functions to their Bapel implementations in `bin/query.bpl`:

| Go Source File | Go Function / Method | Proposed Bapel API (`bin/query.bpl`) | Notes & Implementation Strategy |
| :--- | :--- | :--- | :--- |
| [module_finder.go](file:///usr/local/google/home/jabolopes/.gemini/jetski/scratch/bapel/query/module_finder.go) | `newModuleFinder` | `fn mk_module_finder() -> ModuleFinder` | Reads `workspace.bpl` (or default paths) using `fs::exists` and parses mappings. |
| [module_finder.go](file:///usr/local/google/home/jabolopes/.gemini/jetski/scratch/bapel/query/module_finder.go) | `lookupModuleByName` / `ByPrefix` | `fn lookup_module(finder: &ModuleFinder, mod_id: &String) -> (bool, String)` | Unifies and refines `findBestMatch` in [bin/main.bpl](file:///usr/local/google/home/jabolopes/.gemini/jetski/scratch/bapel/bin/main.bpl). |
| [module_finder.go](file:///usr/local/google/home/jabolopes/.gemini/jetski/scratch/bapel/query/module_finder.go) | `baseSourceFilename` | `fn base_filename(finder: &ModuleFinder, mod_id: &String) -> String` | Replaced `resolveModule` in `main.bpl`. Uses `fs::join` and string separator replacement. |
| [module_finder.go](file:///usr/local/google/home/jabolopes/.gemini/jetski/scratch/bapel/query/module_finder.go) | `implSourceFilename` | `fn impl_filename(base_file: &String, rel_impl: &String) -> String` | Computes `fs::join (fs::parent_path (*base_file), *rel_impl)`. |
| [query_source_file.go](file:///usr/local/google/home/jabolopes/.gemini/jetski/scratch/bapel/query/query_source_file.go) | `queryAnnotationNonBplFile` | `fn query_annotation_file(path: &String) -> SourceFileQuery` | Reads file line-by-line using `IStringStream` / `getline`. Scans for `import ` prefixes and `// @bpl: ` annotation strings. |
| [query_source_file.go](file:///usr/local/google/home/jabolopes/.gemini/jetski/scratch/bapel/query/query_source_file.go) | `queryDeclsBplFile` | `fn query_bpl_file(path: &String) -> SourceFileQuery` | For MVP, invokes `bootstrap/parser -format=flat <path>` and extracts imports, impls, and decls from flat text (extending `parseSourceFileFlat`). |
| [query_source_file.go](file:///usr/local/google/home/jabolopes/.gemini/jetski/scratch/bapel/query/query_source_file.go) | `QuerySourceFile` | `fn query_source_file(path: &String) -> SourceFileQuery` | Dispatches to `query_bpl_file` if extension is `.bpl`, otherwise `query_annotation_file`. |
| [querier.go](file:///usr/local/google/home/jabolopes/.gemini/jetski/scratch/bapel/query/querier.go) | `QueryModule` | `fn query_module(finder: &ModuleFinder, mod_id: &String) -> ModuleQuery` | Queries base filename, iterates over `impls` to query implementation files, merges vectors, deduplicates, and sorts. |
| [querier.go](file:///usr/local/google/home/jabolopes/.gemini/jetski/scratch/bapel/query/querier.go) | `QueryModuleExports` | `fn query_module_exports(finder: &ModuleFinder, mod_id: &String) -> ModuleQuery` | Calls `query_module` and filters `decls` for exported symbols (or flag prefix in flat format). |

---

## 4. Implementation Phases

### Phase 0: Standard Library & Compiler Prerequisites (COMPLETED)
Before implementing `query`, five prerequisites were addressed:
1. **File Input Reading (`Ifstream`) (COMPLETED):** Added `pub type Ifstream` and an `impl` block in [bapel/stl.bpl](file:///usr/local/google/home/jabolopes/.gemini/jetski/scratch/bapel/bapel/stl.bpl#L204-L217), backed by `IfstreamImpl` in [bapel/stl_fstream.h](file:///usr/local/google/home/jabolopes/.gemini/jetski/scratch/bapel/bapel/stl_fstream.h#L36-L57). Both `Ofstream::open` and `Ifstream::open` take reference arguments (`&String`). [getline](file:///usr/local/google/home/jabolopes/.gemini/jetski/scratch/bapel/bapel/stl_string.h#L14) automatically works with `Ifstream` since it is polymorphic over stream types.
2. **Generic Methods in `impl` Blocks (COMPLETED):** Updated [comp/cpp_printer.go](file:///usr/local/google/home/jabolopes/.gemini/jetski/scratch/bapel/comp/cpp_printer.go#L1233-L1277) to emit C++ template parameters for generic methods inside non-generic `impl` blocks (enabling polymorphic methods like `Ifstream::read`).
3. **String Utilities (COMPLETED):** Added `ends_with`, `remove_prefix`, `remove_suffix`, `trim_prefix`, `trim_suffix`, and `rfind` to `StringView` and `String` in [bapel/stl.bpl](file:///usr/local/google/home/jabolopes/.gemini/jetski/scratch/bapel/bapel/stl.bpl#L27-L130) and [bapel/stl_string.h](file:///usr/local/google/home/jabolopes/.gemini/jetski/scratch/bapel/bapel/stl_string.h#L68-L129).
4. **Hash Map (`UnorderedMap`) (COMPLETED):** Added `pub type UnorderedMap ['k, 'v]` in [bapel/stl.bpl](file:///usr/local/google/home/jabolopes/.gemini/jetski/scratch/bapel/bapel/stl.bpl#L219-L238), backed by `std::unordered_map` in [bapel/stl_unordered_map.h](file:///usr/local/google/home/jabolopes/.gemini/jetski/scratch/bapel/bapel/stl_unordered_map.h), including `mk`, `insert`, `size`, `empty`, `contains`, and `get`.
5. **Vector Sorting & Deduplication (COMPLETED):** Added `sort` and `dedup` methods to `Vector` in [bapel/stl.bpl](file:///usr/local/google/home/jabolopes/.gemini/jetski/scratch/bapel/bapel/stl.bpl#L148-L153) and [bapel/stl_vector.h](file:///usr/local/google/home/jabolopes/.gemini/jetski/scratch/bapel/bapel/stl_vector.h), enabling in-place sorting and deduplication of vector elements.

### Phase 1: Module Finder & Workspace Resolution (COMPLETED)
1. **Created `bin/query.bpl`:** Implemented as `implements bin.main` with `ModuleFinder` and package mapping logic (`mk_module_finder`, `lookup_module`, `base_filename`, `impl_filename`).
2. **Updated `bin/main.bpl`:** Changed header to `module bin.main`, added `impls { "query.bpl" }`, removed old ad-hoc mapping functions (`MatchResult`, `PackageMapping`, `resolveMappedPath`, etc.), and replaced `resolveModule` with `mk_module_finder` and `base_filename`.

### Phase 2: Source File Querying (`query_source_file`)
1. Implement `query_annotation_file(&String)`:
   - Read header/source files line-by-line using `IStringStream` and `getline`.
   - Check for `import <mod>;` statements and extract module names.
   - Check for `// @bpl: ` prefixes and extract embedded Bapel declarations.
2. Implement `query_bpl_file(&String)`:
   - Adapt `parseSourceFileFlat` from `main.bpl` to populate `SourceFileQuery` (including decls and flags).
3. Implement the unified `query_source_file(&String)` entrypoint.

### Phase 3: Module Querier & Deduping (`query_module`)
1. Implement vector helper functions in `bin/query.bpl`:
   - `merge_unique_strings(dst: &Vector String, src: &Vector String)`
   - Leverage `Vector::sort` and `Vector::dedup` from `bapel.stl` for sorting and compaction of module IDs and filenames.
2. Implement `query_module` and `query_module_exports`.

### Phase 4: Driver Integration ([bin/main.bpl](file:///usr/local/google/home/jabolopes/.gemini/jetski/scratch/bapel/bin/main.bpl))
1. Replace remaining ad-hoc file parsing functions (`parseSourceFileFlat`, `collectImplImports`) with calls to `query_source_file()` and `query_module()`.
2. Verify clean compilation and execution of `./bpl build`.

### Phase 5: Verification & Testing
1. Add a test target in [Makefile](file:///usr/local/google/home/jabolopes/.gemini/jetski/scratch/bapel/Makefile) or create a new test program `tests/query_test.bpl`.
2. Execute queries against `bapel/core` and `./bapel/core_impl.h` using both `bootstrap/querier` (Go) and the new Bapel implementation, asserting identical output.

---

## 6. Verification Strategy

To ensure parity between the Go querier and `bapel.query`:
1. Run `./bpl query bapel/core` (which tests the existing Go querier).
2. Run a new Bapel test binary that invokes `query_module(&finder, &"bapel.core".to_string)` and prints formatted results.
3. Diff the output of both tools across multiple target modules in the repo (`bapel.core`, `bapel.stl`, `bapel.os`).
