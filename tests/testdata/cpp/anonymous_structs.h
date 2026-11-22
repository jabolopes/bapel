#pragma once

#include <array>
#include <cmath>
#include <cstdlib>
#include <functional>
#include <optional>
#include <string>
#include <tuple>
#include <variant>
#include <vector>

#include "bapel/core_variant.h"

struct __anonym_6241259967628382541 {
  int64_t lhs;
  int64_t rhs;
};
struct Expr
    : ::bapel::Variant<::bapel::Arm<"lit", int64_t>,
                       ::bapel::Arm<"add", __anonym_6241259967628382541>> {
  using ::bapel::Variant<
      ::bapel::Arm<"lit", int64_t>,
      ::bapel::Arm<"add", __anonym_6241259967628382541>>::Variant;
};
struct __anonym_14495726933898134737 {
  int64_t a;
  int64_t b;
};
struct Outer {
  int64_t label;
  __anonym_14495726933898134737 inner;
};
struct __anonym_13322383462023129514 {
  int64_t x;
  int64_t y;
};
int64_t eval(::Expr e);
int64_t sum_point(__anonym_13322383462023129514 p);
::Expr make_expr();
