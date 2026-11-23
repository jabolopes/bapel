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

#ifndef __anonym_4591088945429812398_DEFINED
#define __anonym_4591088945429812398_DEFINED
struct __anonym_4591088945429812398 {
  int64_t lhs;
  int64_t rhs;
};
#endif
BAPEL_VARIANT(Expr, (lit, (int64_t)), (add, (__anonym_4591088945429812398)));
#ifndef __anonym_3433834150179409729_DEFINED
#define __anonym_3433834150179409729_DEFINED
struct __anonym_3433834150179409729 {
  int64_t a;
  int64_t b;
};
#endif
struct Outer {
  int64_t label;
  __anonym_3433834150179409729 inner;
};
#ifndef __anonym_11548725740611542807_DEFINED
#define __anonym_11548725740611542807_DEFINED
struct __anonym_11548725740611542807 {
  int64_t x;
  int64_t y;
};
#endif
int64_t eval(::Expr e);
int64_t sum_point(__anonym_11548725740611542807 p);
::Expr make_expr();
