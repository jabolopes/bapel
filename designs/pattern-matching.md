# Pattern Matching in Bapel

This document describes how `match` expressions work in Bapel.

*(Deferred pattern-matching extensions—wildcard payload bindings, literal/primitive matching, nested destructuring, or-patterns, and guards—are tracked in [`designs/features.md`](features.md).)*

---

## 1. Surface Syntax

`match` matches a `variant` value against tag arms and an optional trailing wildcard (`_`) arm:

```bapel
pub type Color = variant {
  red (),
  green (),
  blue ()
}

pub fn is_red(c: Color) -> bool {
  match c {
    red u => true,
    _ => false
  }
}
```

* **Tag arm (`<tag> <var> => <expr>`):** Matches a specific variant tag and binds its payload to `<var>` in `<expr>`.
* **Wildcard arm (`_ => <expr>`):** Matches any remaining variant tags without binding a variable. At most one `_` arm is allowed and it must be the last arm. A trailing `_` arm is allowed even when all tags are already covered (providing a forward-compatible fallback if new tags are added later).

---

## 2. Grammar, AST & Core IR

* **Grammar ([`cpp_parser/bapel.g4`](../cpp_parser/bapel.g4)):**
  ```antlr
  matchTerm: MATCH expression LBRACE matchArms (COMMA?) RBRACE;
  matchArms: matchArm (COMMA matchArm)*;
  matchArm: id IDENTIFIER FAT_ARROW expression
          | IDENTIFIER FAT_ARROW expression
          ;
  ```
  Keeping `_` as `IDENTIFIER` avoids lexer conflicts with `_` used as a variable or parameter name. `visitMatchArm` ([`cpp_parser/ast_builder.h`](../cpp_parser/ast_builder.h)) validates that a single-identifier arm is `"_"` and sets `tag = "_"`, `arg = ""`.
* **AST & Core IR ([`ast/ast_expr.h`](../ast/ast_expr.h), [`bin/ir_term.h`](../bin/ir_term.h)):**
  * Both `ast::MatchArm` and `ir::MatchArm` represent a wildcard arm with `tag == "_"` and `arg == ""` (`ir::MatchArm::is_wildcard()`).
  * `ast::MatchArm::to_string()` and `ir::IrTerm::to_string()` format wildcard arms as `_ -> <body>` and tag arms as `<tag> <arg> -> <body>`.

---

## 3. Type Inference & Verification

### 3.1 Inference (`ts/inferencer.h` — `infer_match`)
1. Infers the scrutinee and reduces/unfolds it (including `MuType` unfolding) to a `VariantType`.
2. Iterates through arms tracking `std::vector<bool> seen(tags.size(), false)`, `remaining = tags.size()`, and `has_wildcard`:
   * Rejects any arm appearing after `_`.
   * For `_`: sets `has_wildcard = true` and `arm.index = std::nullopt` (no variable bound).
   * For tag arms: resolves `arm.tag` in `tags`, rejects invalid or duplicate tags, records `arm.index`, marks `seen[idx] = true`, decrements `remaining`, and binds `arm.arg` to the tag's payload type.
   * Propagates the expected result type (`expect_arm = match_type ? &*match_type : expect_type`) so polymorphic literals like `_ => 0` infer their type from the enclosing return type or `let` annotation.
3. If `!has_wildcard && remaining > 0`, reports the missing tag as a non-exhaustive match error.

### 3.2 Trusted Kernel Verification (`ts/typecheck.h` — `typecheck_match`)
Independently verifies the elaborated `MatchTerm`:
* Rejects arms after `_`, duplicate tags, and non-exhaustive matches without `_`.
* Verifies wildcard arms have `!arm.index.has_value()`.
* Verifies non-wildcard arms in `O(1)` using `arm.index` (`0 <= *arm.index < tags.size() && tags[*arm.index].id == arm.tag`).
* Verifies every arm body subtypes `match_type` and `match_type <: *term->type`.

---

## 4. C++17 Code Generation (`bin/codegen_impl.h`)

`CppPrinter::print_match_term` lowers `MatchTerm` to a `switch` on `std::variant::index()`:
* **Tag arm:** `case <idx>: { auto &<arg> = std::get<<idx>>(<v>); <body>; }`
* **Wildcard arm:** `default: { <body>; }`

---

## 5. Test Coverage

* **Positive golden tests:**
  * [`tests/testdata/in/pattern_match.in`](../tests/testdata/in/pattern_match.in) (wildcard fallback, redundant trailing wildcard, wildcard-only match with expected type propagation for `_ => 0`, plus `CppPrinterTest.IsValidCpp`).
  * [`tests/testdata/in/variants.in`](../tests/testdata/in/variants.in) (exhaustive variant matches).
* **Negative tests:**
  * [`tests/testdata/in/bad_pattern_match_unreachable.in`](../tests/testdata/in/bad_pattern_match_unreachable.in) (arm after `_ =>`)
  * [`tests/testdata/in/bad_pattern_match_duplicate.in`](../tests/testdata/in/bad_pattern_match_duplicate.in) (duplicate tag arm)
  * [`tests/testdata/in/bad_pattern_match_non_exhaustive.in`](../tests/testdata/in/bad_pattern_match_non_exhaustive.in) (missing tag without `_ =>`)
