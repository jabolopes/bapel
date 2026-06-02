#pragma once

#include "lets.h"

struct V : std::variant<int64_t /* x */> {};
template <typename a>
a id(a);
std::monostate lets();
template <typename a>
a id(a x) {
  return x;
}

std::monostate lets();
