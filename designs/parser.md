# Bapel C++ Parser & Lexer Design

## Overview & Goal

The goal of this project is to port the Bapel Lexer and Parser from Go ([cpp_parser/](cpp_parser/), [parse/](parse/)) to native C++17, continuing the progress toward making the Bapel compiler fully self-bootstrapping.

We will continue using **ANTLR4** by targeting its official C++ backend (`-Dlanguage=Cpp`) and linking against the standard ANTLR4 C++ runtime library (`libantlr4-runtime`).

---

## Technical Feasibility & Prerequisites

1. **ANTLR Tool:** `antlr4` (Version 4.9.2) is available and supports `-Dlanguage=Cpp`.
2. **ANTLR C++ Runtime:**
   - Headers: `/usr/include/antlr4-runtime/antlr4-runtime.h`
   - Library: `/usr/lib/x86_64-linux-gnu/libantlr4-runtime.a` / `libantlr4-runtime.so`
3. **Grammar:** [cpp_parser/bapel.g4](cpp_parser/bapel.g4) defines complete grammar rules for:
   - `sourceFile` (base and impl files)
   - `workspace` (package rules)
   - `decl` (term and type declarations, including `@bpl:` header annotations)
   - `type_` and `expression`

---

## Architecture & Data Flow

```
                      ┌────────────────────────┐
                      │  cpp_parser/bapel.g4   │
                      └───────────┬────────────┘
                                  │ antlr4 -Dlanguage=Cpp -visitor
                                  ▼
           ┌──────────────────────────────────────────────┐
           │ Generated C++ Lexer, Parser & Visitor        │
           │ (cpp_parser/generated/bapelLexer.cpp, etc.)  │
           └──────────────────────┬───────────────────────┘
                                  │
                                  ▼
           ┌──────────────────────────────────────────────┐
           │ C++ AST Builder & Visitor (ast_builder.h)    │
           │ Implements bapelBaseVisitor & validation     │
           └──────────────────────┬───────────────────────┘
                                  │
                                  ▼
           ┌──────────────────────────────────────────────┐
           │ C++ AST Data Structures (ast/ast_*.h)        │
           │ (SourceFile, Workspace, Decl, Expr, Type)    │
           └──────────────────────┬───────────────────────┘
                                  │
                                  ▼
           ┌──────────────────────────────────────────────┐
           │ bootstrap/parser CLI (JSON & Flat Emitter)   │
           └──────────────────────────────────────────────┘
```

---

## Detailed Implementation Phases

### Phase 1: Native C++ AST Data Structures (`ast/ast_*.h`) - [COMPLETED]
- [x] Defined value-oriented, serializable C++ AST data structures matching Go AST in [ast/](ast/).
- [x] Reused core IR types from [bin/ir_base.h](bin/ir_base.h) and [bin/ir_type.h](bin/ir_type.h).
- [x] Created [ast/ast_pos.h](ast/ast_pos.h), [ast/ast_expr.h](ast/ast_expr.h), [ast/ast_decl.h](ast/ast_decl.h), [ast/ast_source_file.h](ast/ast_source_file.h), [ast/ast_workspace.h](ast/ast_workspace.h), and [ast/ast.h](ast/ast.h).
- [x] Implemented `to_string(bool with_pos)` and `to_json()` matching Go `encoding/json` output schema.
- [x] Verified compilation and serialization with standalone test suite.

### Phase 2: ANTLR4 C++ Generation & Build Pipeline - [COMPLETED]
- [x] Generated C++ parser, lexer, and visitor sources into `cpp_parser/generated/` using ANTLR 4.9.2.
- [x] Verified clean compilation of generated sources against system `<antlr4-runtime/antlr4-runtime.h>`.
- [x] Added `gen-parser-cpp` build target to [Makefile](Makefile).

### Phase 3: C++ AST Visitor & Error Listener (`cpp_parser/ast_builder.h`) - [COMPLETED]
- [x] Implemented `BapelErrorListener` in [cpp_parser/error_listener.h](cpp_parser/error_listener.h) capturing syntax error messages with line/column and unterminated literal diagnostics.
- [x] Implemented `AstBuilder` in [cpp_parser/ast_builder.h](cpp_parser/ast_builder.h) subclassing `bapelBaseVisitor` and constructing typed C++ AST objects for all rules.
- [x] Tested parser with complete traits and impl source file parsing and JSON serialization.

### Phase 4: Native C++ Parser Driver (`cpp_parser/main.cc`)
- Implement `cpp_parser/main.cc`:
  - CLI Flags:
    - `--symbol=<SourceFile|Workspace|Decl|Type|Expr>` (default: `SourceFile`)
    - `--format=<json|flat>` (default: `json`)
    - `--filename=<name>` (for position info when reading stdin)
    - Positional argument: input file (or stdin if omitted)
  - Output Modes:
    - `json`: Prints formatted JSON matching Go AST schema.
    - `flat`: Prints line-oriented format for `bpl query` (`IMPORT ...`, `IMPL ...`, `DECL ...`, `FUNC ...`, `TRAIT_IMPL ...`).
- Replace the Go `bootstrap/parser` target in `Makefile` with the native C++ binary.

### Phase 5: Verification & Parity Testing
- Verify all test inputs in [parse/testdata/in/*.in](parse/testdata/in/) and [parse/testdata/parsed/](parse/testdata/parsed/) (both valid files and `bad_*.in` error cases).
- Run `go test ./...` to confirm `bootstrap/typechecker` consumes C++ `bootstrap/parser` output without errors.
- Confirm 100% diff-free parity against existing golden outputs.

### Phase 6: Deprecate & Remove Go Parser
- Remove Go ANTLR generated code: `cpp_parser/parser/*.go`.
- Remove Go parser driver: `cpp_parser/main.go`.
- Update `Makefile` and documentation.

