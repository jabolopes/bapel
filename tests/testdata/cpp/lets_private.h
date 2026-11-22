#pragma once

#include "lets.h"

struct V : ::bapel::Variant<::bapel::Arm<"x", int64_t>> {
  using ::bapel::Variant<::bapel::Arm<"x", int64_t>>::Variant;
};
template <typename a>
a id(a x);
std::monostate lets();
template <typename a>
a id(a x) {
  return x;
}
