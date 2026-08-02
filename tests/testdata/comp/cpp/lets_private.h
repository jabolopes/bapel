#pragma once

#include "lets.h"

struct V : std::variant<int64_t /* x */> {
  using std::variant<int64_t /* x */>::variant;
};
template <typename a>
a id(a);
std::monostate lets();
template <typename a>
a id(a x) {
  return x;
}

std::monostate lets();
