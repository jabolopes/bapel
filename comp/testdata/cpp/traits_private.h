#pragma once

#include "traits.h"

struct MyStruct {
  int64_t x;
};
int64_t run(MyStruct);
int64_t run(MyStruct s);
namespace traits {
template <>
struct Size<::MyStruct> {
  using Self = ::MyStruct;
  static int64_t size(Self s) { return s.x; }
};
}  // namespace traits