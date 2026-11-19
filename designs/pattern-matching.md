# Pattern Matching in Bapel

This document designs pattern matching improvements for Bapel, starting with the Phase 1 MVP needed for self-hosting (wildcard `_ =>` arms and exhaustiveness/duplicate checks) and outlining future extensions.

---

## 1. Current State

Currently, `match` in Bapel only supports flat variant tag matching with a single bound variable:

```bapel
match e {
  lit x => x,
  add p => p.lhs + p.rhs
}
```

* **Grammar ([`cpp_parser/bapel.g4`](../cpp_parser/bapel.g4)):** `matchArm: id IDENTIFIER FAT_ARROW expression;`
* **Core IR ([`bin/ir_term.h`](../bin/ir_term.h)):** `MatchArm { std::string tag; std::string arg; std::shared_ptr<IrTerm> body; std::optional<int> index; }`
* **Typechecking ([`ts/inferencer.h`](../ts/inferencer.h), [`ts/typecheck.h`](../ts/typecheck.h)):** Checks that each `arm.tag` belongs to the matched `VariantType`, binds `arm.arg` to the tag's payload type, and unifies/subtypes all arm body types. Does **not** check for duplicate tags or exhaustiveness.
* **Codegen ([`bin/codegen_impl.h`](../bin/codegen_impl.h)):** Emits `switch (v.index()) { case idx: { auto &arg = std::get<idx>(v); body; } }`.

### Limitations for Self-Hosting
Writing compiler passes over large variants (`IrTerm` has 16 cases; `IrType` has 12 cases) without a wildcard arm forces every helper function (e.g. `is_cpp_statement`) to enumerate every tag of the variant.

---

## 2. Phase 1: Wildcard Arms (`_ =>`) & Exhaustiveness Checks (MVP)

### 2.1 Surface Syntax

```bapel
fn is_lit(e: Expr) -> bool {
  match e {
    lit x => true,
    _ => false
  }
}
```

### 2.2 Grammar & AST (`cpp_parser/bapel.g4`, `cpp_parser/ast_builder.h`)

* Extend `matchArm` in [`cpp_parser/bapel.g4`](../cpp_parser/bapel.g4):
  ```antlr
  matchArm: id IDENTIFIER FAT_ARROW expression
          | IDENTIFIER FAT_ARROW expression
          ;
  ```
  Keeping `_` as `IDENTIFIER` avoids introducing a `'_'` keyword token that would conflict with `_` used as a variable/parameter name.
* In `visitMatchArm` ([`cpp_parser/ast_builder.h`](../cpp_parser/ast_builder.h)):
  * If `ctx->id()` is present: normal tag arm (`arm.tag = tag.value`, `arm.arg = ctx->IDENTIFIER()->getText()`).
  * Otherwise: verify `ctx->IDENTIFIER()->getText() == "_"` (raising a syntax error if not), and construct a wildcard arm with `arm.tag = "_"` and `arm.arg = ""`.

### 2.3 AST & Core IR Representation (`ast/ast_expr.h`, `bin/ir_term.h`)

Both `ast::MatchArm` ([`ast/ast_expr.h`](../ast/ast_expr.h)) and `ir::MatchArm` ([`bin/ir_term.h`](../bin/ir_term.h)) represent a wildcard arm using `tag == "_"` and `arg == ""` (no extra fields needed, so `ast/ast_desugar.h` requires no changes):

* Add `bool is_wildcard() const { return tag == "_"; }` to both `ast::MatchArm` and `ir::MatchArm`.
* Update `ast::MatchArm::to_string(bool with_pos)` ([`ast/ast_expr.h`](../ast/ast_expr.h)) and `ir::IrTerm::to_string()` ([`bin/ir_term.h`](../bin/ir_term.h)) to format wildcard arms as `_ -> <body>` instead of `_  -> <body>`.

### 2.4 Type Inference & Verification (`ts/inferencer.h`, `ts/typecheck.h`)

For `match expr { arm_1, ..., arm_n }` on a `VariantType`:
1. **Wildcard placement:**
   * At most one `_` arm is allowed, and it must be the **last** arm (arms after `_` are unreachable and rejected).
   * A trailing `_` arm is allowed even if all variant tags are already covered (providing a forward-compatible fallback if new tags are added later).
2. **Duplicate tags:**
   * Track seen tag names; reject duplicate tag arms (e.g., two `lit` arms).
3. **Wildcard arm typing & expected type propagation:**
   * Do not look up `"_"` in `variant_type.tags()` and do not bind any variable in `context_`.
   * Leave `arm.index = std::nullopt`.
   * Propagate the enclosing expected type to the first arm (`expect_arm = match_type ? &*match_type : expect_type`), so polymorphic literals like `_ => 0` infer their type from the enclosing return type or `let` annotation.
   * Infer/typecheck `arm.body` against `expect_arm` / `match_type`.
4. **Exhaustiveness check:**
   * If no `_` arm is present, verify that every tag of `variant_type` is covered by an arm (or report the missing tag(s)).
   * *(Note: verify against existing golden tests in `tests/testdata/in/` to ensure no existing tests intentionally rely on non-exhaustive matches.)*

### 2.5 C++17 Code Generation (`bin/codegen_impl.h`)

In `CppPrinter::print_match_term`:
* For normal arms: emit `case <idx>: { auto &<arg> = std::get<<idx>>(<variant_id>); <body>; break; }` (or without `break;` when `last_term` emits `return`).
* For wildcard arms (`arm.is_wildcard()`):
  ```cpp
  default: {
    <body>;
  }
  ```

### 2.6 Verification Plan
* Regenerate parser (`make gen-parser`).
* Add positive golden test `tests/testdata/in/pattern_match.in` (covering `_ =>` on a multi-tag variant, redundant trailing `_ =>`, and expected type propagation with integer literals like `_ => 0`, including `CppPrinterTest.IsValidCpp`).
* Add negative tests (`tests/testdata/in/bad_pattern_match_*.in`) for:
  * Arm after `_ =>` (unreachable arm).
  * Duplicate tag arms.
  * Non-exhaustive `match` without `_ =>`.
* Verify `make test` passes.

---

## 3. Future Extensions (Deferred)

The following extensions are out of scope for Phase 1 and can be added later as needed:
* **Wildcard payload binding (`nil _ => 0`):** Allow `_` as the payload variable name without binding a C++ reference (avoiding unused-variable warnings).
* **Literal & primitive matching:** `match ch { 'a' => ..., '0' => ..., _ => ... }` on `i8`, `i64`, `bool`, and `StringView` (useful when porting the lexer in Phase 6 of [`designs/migration.md`](migration.md)).
* **Nested & destructuring patterns:** Tuple/struct destructuring inside arms (e.g., `cons (x, xs) => ...`), or-patterns (`a | b => ...`), and match guards (`lit x if x > 0 => ...`).
