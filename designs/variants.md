# Variant C++ Codegen & Interop Design (`BapelVariant`)

## 1. Motivation

Bapel's `variant { tag0 T0, tag1 T1, ... }` currently compiles to `std::variant<T0, T1, ...>`. While `std::variant` has an optimal runtime layout (`max(sizeof(Ti)) + tag`, zero heap allocations, RAII), raw `std::variant` relies on positional numeric indices (`std::in_place_index<3>`, `std::get<3>(v)`, `v.index() == 3`), making C++ interop awkward.

Meanwhile, the legacy hand-written C++ IR (`bin/ir_*.h`) uses `enum class case_val` + a bag of `std::shared_ptr` fields, creating a mismatch with Bapel's `bapel.compiler.ir` types.

By moving generated C++ to **C++20** and using a small `BapelVariant<Arm<"tag", T>...>` template, we get:
- **Zero-overhead `std::variant` layout** with compile-time tag-name resolution.
- **100% uniformity across named and anonymous variants** (works inline in function signatures without lifting or macros).
- **Ergonomic C++ interop** (`mk<"tag">(...)`, `is<"tag">()`, `as<"tag">()`, `index_of<"tag">()`), allowing us to align `bin/ir_*.h` with Bapel's generated C++ and swap in `bapel/compiler/ir.h` directly with no conversion bridge.

---

## 2. Design

### 2.1 Shared Runtime Template (`bapel/core_variant.h`)

A ~35-line header included by generated `.h` files (and hand-written C++):

```cpp
#pragma once

#include <array>
#include <cstddef>
#include <string_view>
#include <utility>
#include <variant>

namespace bapel {

template <size_t N>
struct FixedString {
  char data[N];
  constexpr FixedString(const char (&s)[N]) {
    for (size_t i = 0; i < N; ++i) data[i] = s[i];
  }
  constexpr std::string_view sv() const { return {data, N - 1}; }
};

template <FixedString Tag, typename T>
struct Arm {
  static constexpr auto tag = Tag;
  using type = T;
};

template <typename... Arms>
struct Variant : std::variant<typename Arms::type...> {
  using Base = std::variant<typename Arms::type...>;
  using Base::Base;

  template <FixedString Tag>
  static constexpr size_t index_of() {
    constexpr std::array<std::string_view, sizeof...(Arms)> tags = { Arms::tag.sv()... };
    for (size_t i = 0; i < tags.size(); ++i) {
      if (tags[i] == Tag.sv()) return i;
    }
    return sizeof...(Arms);
  }

  template <FixedString Tag, typename... Args>
  static Variant mk(Args&&... args) {
    constexpr size_t idx = index_of<Tag>();
    static_assert(idx < sizeof...(Arms), "unknown variant tag");
    return Variant(std::in_place_index<idx>, std::forward<Args>(args)...);
  }

  template <FixedString Tag>
  bool is() const {
    constexpr size_t idx = index_of<Tag>();
    static_assert(idx < sizeof...(Arms), "unknown variant tag");
    return this->index() == idx;
  }

  template <FixedString Tag>
  auto& as() {
    constexpr size_t idx = index_of<Tag>();
    static_assert(idx < sizeof...(Arms), "unknown variant tag");
    return std::get<idx>(*this);
  }

  template <FixedString Tag>
  const auto& as() const {
    constexpr size_t idx = index_of<Tag>();
    static_assert(idx < sizeof...(Arms), "unknown variant tag");
    return std::get<idx>(*this);
  }
};

} // namespace bapel
```

### 2.2 Generated C++ Shape

For a Bapel variant:
```bapel
pub type IntList = variant {
  nil (),
  cons (i64, Ptr IntList)
}
```

`CppPrinter` emits:
```cpp
struct IntList
    : ::bapel::Variant<::bapel::Arm<"nil", std::monostate>,
                       ::bapel::Arm<"cons", std::tuple<int64_t, ::Ptr<::IntList>>>> {
  using ::bapel::Variant<::bapel::Arm<"nil", std::monostate>,
                         ::bapel::Arm<"cons", std::tuple<int64_t, ::Ptr<::IntList>>>>::Variant;
};
```

And inline anonymous variants emit `::bapel::Variant<::bapel::Arm<"a", int64_t>, ::bapel::Arm<"b", bool>>` directly—with the exact same `mk<"a">()`, `is<"a">()`, `as<"a">()`, and `index_of<"a">()` API!

---

## 3. Execution Plan

1. **Bump C++ standard to `-std=c++20` and add `bapel/core_variant.h`:**
   - Update `-std=c++17` $\to$ `-std=c++20` in `Makefile`, `tests/cpp_printer_test.cc`, `tests/typecheck_test.cc`, and `bin/main.bpl` (Bazel `copts`).
   - Add `bapel/core_variant.h` defining `bapel::FixedString`, `bapel::Arm`, and `bapel::Variant`.
2. **Emit `::bapel::Variant<::bapel::Arm<...>>` in `CppPrinter` (`bin/codegen_impl.h`):**
   - Include `"bapel/core_variant.h"` in generated public headers.
   - Update `VariantType` printing in `bin/codegen_impl.h` to emit `::bapel::Variant<::bapel::Arm<"tag", T>...>`.
   - Regenerate `tests/testdata/cpp/*.h` golden files, run `make test` and `make`, and promote `cp bpl bootstrap/bpl`.
3. **Rewrite existing C++ IR (`bin/ir_*.h`) to use `bapel::Variant` incrementally:**
   - Convert `IrLiteral`, `IrKind`, `IrType`, `IrDecl`, and `IrTerm` one at a time, keeping `make test` (23/23) and `make` green at every step.
4. **Swap in Bapel-generated `bapel/compiler/ir.h`:**
   - Retire the duplicate hand-written C++ IR definitions in `bin/ir_*.h` and use `bapel/compiler/ir.h` directly with zero conversion bridge.
5. **Proceed to Phase 2 (`bapel.compiler.codegen.cpp`)** emitting the same `::bapel::Variant<::bapel::Arm<...>>` codegen.

---

## 4. Verification Plan

1. **Step 1–2:** `make test` (all 23/23 suites, including `CppPrinterTest.GoldenFiles` and `CppPrinterTest.IsValidCpp` under `-std=c++20`), `make`, `cp bpl bootstrap/bpl`, and a second `make` to verify clean self-hosting.
2. **Step 3:** `make test` and `make` green after each IR type rewrite (`IrLiteral` $\to$ `IrKind` $\to$ `IrType` $\to$ `IrDecl` $\to$ `IrTerm`).
3. **Step 4:** `make test`, `make`, and bootstrap fixed-point check with the single Bapel-defined `bapel/compiler/ir.h`.
