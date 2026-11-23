#pragma once

#include "lets.h"

BAPEL_VARIANT(V, (x, (int64_t)));
template <typename a>
a id(a x);
std::monostate lets();
template <typename a>
a id(a x) {
  return x;
}
