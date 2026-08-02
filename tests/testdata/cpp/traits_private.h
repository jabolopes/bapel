#pragma once

#include "traits.h"

struct MyStruct {
  int64_t x;
};
std::tuple<int64_t, int8_t> run(::MyStruct s, ::Ptr<::Vector<int8_t>> v);
namespace traits {
template <>
struct Size<::MyStruct> {
  using Self = ::MyStruct;
  static inline int64_t size(Self s) { return s.x; }
};
}  // namespace traits
namespace traits {
template <typename a>
struct Indexable<::Vector<a>, a> {
  using Self = ::Vector<a>;
  static inline a get(::Ptr<Self> v, int64_t index) {
    return vector_get<a>(v, index);
  }
};
}  // namespace traits