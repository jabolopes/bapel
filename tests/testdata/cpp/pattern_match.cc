
#include "pattern_match_private.h"

bool is_red(::Color c) {
  {
    auto __v_0 = c;
    switch (__v_0.index()) {
      case 0: {
        auto& u = std::get<0>(__v_0);
        return true;
      }
      default: {
        return false;
      }
    }
  };
}

int64_t color_code(::Color c) {
  {
    auto __v_1 = c;
    switch (__v_1.index()) {
      case 0: {
        auto& u = std::get<0>(__v_1);
        return static_cast<int64_t>(1);
      }
      case 1: {
        auto& u = std::get<1>(__v_1);
        return static_cast<int64_t>(2);
      }
      case 2: {
        auto& u = std::get<2>(__v_1);
        return static_cast<int64_t>(3);
      }
      default: {
        return static_cast<int64_t>(0);
      }
    }
  };
}

int64_t always_zero(::Color c) {
  {
    auto __v_2 = c;
    switch (__v_2.index()) {
      default: {
        return static_cast<int64_t>(0);
      }
    }
  };
}
