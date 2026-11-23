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

#ifndef __anonym_6241259967628382541_DEFINED
#define __anonym_6241259967628382541_DEFINED
struct __anonym_6241259967628382541 {
  int64_t lhs;
  int64_t rhs;
};
#endif
BAPEL_VARIANT(Expr, (lit, (int64_t)), (add, (__anonym_6241259967628382541)));
#ifndef __anonym_14495726933898134737_DEFINED
#define __anonym_14495726933898134737_DEFINED
struct __anonym_14495726933898134737 {
  int64_t a;
  int64_t b;
};
#endif
struct Outer {
  int64_t label;
  __anonym_14495726933898134737 inner;
};
#ifndef __anonym_13322383462023129514_DEFINED
#define __anonym_13322383462023129514_DEFINED
struct __anonym_13322383462023129514 {
  int64_t x;
  int64_t y;
};
#endif
int64_t eval(::Expr e);
int64_t sum_point(__anonym_13322383462023129514 p);
::Expr make_expr();
