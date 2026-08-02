#pragma once

#include "polymorphism.h"

template <typename a>
a id(a x);
std::monostate callPolymorphic();
std::monostate functionSubtyping();
template <typename a>
a id(a x) {
  return x;
}
