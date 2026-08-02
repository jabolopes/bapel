#pragma once

#include "conditionals.h"

std::monostate conditionals();
bool conditionalsPolymorphic();
template <typename a, typename b>
a fconst(a, b);
bool ftrue();
template <typename a>
a id(a);
bool ifLastTerm();
std::monostate conditionals();
bool ifLastTerm();
bool ftrue();
template <typename a>
a id(a x) {
  return x;
}

template <typename a, typename b>
a fconst(a x, b y) {
  return x;
}

bool conditionalsPolymorphic();
