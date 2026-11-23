# Variant C++ Codegen & Interop Design (`BAPEL_VARIANT`)

## 1. Motivation

Bapel's `variant { tag0 T0, tag1 T1, ... }` compiles to `std::variant<T0, T1, ...>`, which gives an optimal runtime layout (`max(sizeof(Ti)) + tag`, zero heap allocations, RAII). However, raw `std::variant` relies on positional numeric indices (`std::in_place_index<3>`, `std::get<3>(v)`, `v.index() == 3`), making C++ interop awkward.

Meanwhile, the legacy hand-written C++ IR (`bin/ir_*.h`) uses `enum class case_val` + a bag of `std::shared_ptr` fields, creating a mismatch with Bapel's `bapel.compiler.ir` types.

By emitting `enum class Tag`, `tag()`, `is(Tag)`, and `BAPEL_VARIANT_ARM` / `BAPEL_VARIANT_UNIT_ARM` helpers on **named** `type Foo = variant { ... }` structs (while keeping **anonymous** inline variants as plain `std::variant<...>`), we get:
- **Zero-overhead `std::variant` layout** with real C++ identifiers (`Tag::nil`, `mk_cons(...)`, `as_cons()`).
- **Zero `.template` noise** inside templates/lambdas and full `-Wswitch` exhaustiveness warnings in hand-written C++.
- **No lifting or ODR risk** for anonymous variants; Bapel `InjectionTerm` and `MatchTerm` codegen stay uniformly on `std::in_place_index<k>` / `.index()` / `std::get<k>`.
- **Direct drop-in alignment** between `bin/ir_*.h` and Bapel-generated `bapel/compiler/ir.h`.

---

## 2. Design

### 2.1 Shared Macro Header (`bapel/core_variant.h`)

```cpp
#pragma once

#include <cstddef>
#include <utility>
#include <variant>

#define BAPEL_VARIANT_UNIT_ARM(Self, Index, TagName)                \
  static Self mk_##TagName() {                                      \
    return Self(std::in_place_index<Index>, std::monostate{});      \
  }                                                                 \
  std::monostate& as_##TagName() { return std::get<Index>(*this); } \
  const std::monostate& as_##TagName() const {                      \
    return std::get<Index>(*this);                                  \
  }

#define BAPEL_VARIANT_ARM(Self, Index, TagName, ...)                \
  static Self mk_##TagName(__VA_ARGS__ v) {                         \
    return Self(std::in_place_index<Index>, std::move(v));          \
  }                                                                 \
  __VA_ARGS__& as_##TagName() { return std::get<Index>(*this); }    \
  const __VA_ARGS__& as_##TagName() const {                         \
    return std::get<Index>(*this);                                  \
  }
```

### 2.2 Generated C++ Shape

For a named Bapel variant:
```bapel
pub type IntList = variant {
  nil (),
  cons (i64, Ptr IntList)
}
```

`CppPrinter` emits:
```cpp
BAPEL_VARIANT(IntList,
  (nil, (std::monostate)),
  (cons, (std::tuple<int64_t, ::Ptr<::IntList>>))
);
```

Anonymous inline `variant { ... }` types continue to emit plain `std::variant<...>`.

---

## 3. Execution Plan

1. **Add `BAPEL_VARIANT_ARM` / `BAPEL_VARIANT_UNIT_ARM` in `bapel/core_variant.h` and emit them in `CppPrinter` (`bin/codegen_impl.h`):**
   - Revert `print_type` (`VariantType`) to `std::variant<...>` for anonymous variants.
   - Update `print_alias_decl` (`VariantType`) to emit `enum class Tag`, `tag()`, `is(Tag)`, and `BAPEL_VARIANT_UNIT_ARM` / `BAPEL_VARIANT_ARM` for each tag.
   - Regenerate `tests/testdata/cpp/*.h` golden files, run `make test` and `make`, and promote `cp bpl bootstrap/bpl`.
2. **Rewrite existing C++ IR (`bin/ir_*.h`) to use `std::variant` + `BAPEL_VARIANT_ARM` incrementally:**
   - Convert `IrLiteral`, `IrKind`, `IrType`, `IrDecl`, and `IrTerm` one at a time, keeping `make test` (23/23) and `make` green at every step.
3. **Swap in Bapel-generated `bapel/compiler/ir.h`:**
   - Retire the duplicate hand-written C++ IR definitions in `bin/ir_*.h` and use `bapel/compiler/ir.h` directly with zero conversion bridge.
4. **Proceed to Phase 2 (`bapel.compiler.codegen.cpp`)** emitting the same named variant helpers.

---

## 4. Verification Plan

1. **Step 1:** `make test` (all 23/23 suites, including `CppPrinterTest.GoldenFiles` and `CppPrinterTest.IsValidCpp`), `make`, `cp bpl bootstrap/bpl`, and a second `make` to verify clean self-hosting.
2. **Step 2:** `make test` and `make` green after each IR type rewrite (`IrLiteral` $\to$ `IrKind` $\to$ `IrType` $\to$ `IrDecl` $\to$ `IrTerm`).
3. **Step 3:** `make test`, `make`, and bootstrap fixed-point check with the single Bapel-defined `bapel/compiler/ir.h`.
