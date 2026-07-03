#pragma once

#include "parameterized_traits.h"

int8_t run(Vector<int8_t>*);
int8_t run(Vector<int8_t>* v);
namespace traits {
template <typename a>
struct Indexable<::Vector<a>, a> {
  using Self = ::Vector<a>;
  static a get(Self* v, int64_t index) { return vector_get<a>(v, index); }
};
}  // namespace traits