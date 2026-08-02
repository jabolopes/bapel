#pragma once

#include "conditionals.h"

std::monostate conditionals();
bool ifLastTerm();
bool ftrue();
template <typename a>
a id(a x);
template <typename a, typename b>
a fconst(a x, b y);
bool conditionalsPolymorphic();
template <typename a>
a id(a x) {
  return x;
}

template <typename a, typename b>
a fconst(a x, b y) {
  return x;
}
