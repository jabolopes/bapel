#pragma once

#include "parameterized_traits.h"

int8_t run(::Ptr<::Vector<int8_t>> v);
namespace traits {
template <typename a>
struct Indexable<::Vector<a>, a> {
  static inline a get(::Ptr<::Vector<a>> v, int64_t index) {
    return vector_get<a>(v, index);
  }
};
}  // namespace traits