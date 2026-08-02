#pragma once

#include "polymorphism.h"

std::monostate callPolymorphic();
std::monostate functionSubtyping();
template <typename a>
a id(a);
template <typename a>
a id(a x) {
  return x;
}

std::monostate callPolymorphic();
std::monostate functionSubtyping();
