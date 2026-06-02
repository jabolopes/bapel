
#include "variants_private.h"

One mkOne() {
  One v = One{static_cast<int64_t>(1)};
  int64_t v1 = std::get<0>(v);
  int64_t v2 = std::get<0>(v);
  int64_t v3;
  {
    auto __v_0 = v;
    switch (__v_0.index()) {
      case 0: {
        auto& l = std::get<0>(__v_0);
        v3 = l;
      }
    }
  };
  {
    auto __v_1 = v;
    switch (__v_1.index()) {
      case 0: {
        auto& l = std::get<0>(__v_1);
        l;
      }
    }
  };
  v = One{static_cast<int64_t>(2)};
  One r = v;
  return r;
}
