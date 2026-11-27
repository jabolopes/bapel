
#include "returns_private.h"

int8_t return1() { return static_cast<int8_t>(0); }

int8_t return2() {
  if (true) {
    return static_cast<int8_t>(0);
    ;
  };
  return static_cast<int8_t>(1);
}

int8_t return3() {
  if (true) {
    return static_cast<int8_t>(0);
    ;
  } else {
    return static_cast<int8_t>(1);
    ;
  };
  return static_cast<int8_t>(2);
}

std::monostate noop() { return std::monostate(); }

std::monostate unit_if_no_else(bool cond) {
  if (cond) {
    return noop();
  };
  return std::monostate();
}

std::monostate unit_early_return_if(bool cond) {
  if (cond) {
    return std::monostate();
    ;
  };
  if (cond) {
    return noop();
  };
  return std::monostate();
}

std::monostate unit_for_loop(int64_t n) {
  int64_t i = static_cast<int64_t>(0);
  while ((i) < (n)) {
    i = (i) + (static_cast<int64_t>(1));
    noop();
  };
  return std::monostate();
}

int64_t value_if_else_implicit(bool cond) {
  if (cond) {
    return static_cast<int64_t>(1);
  } else {
    return static_cast<int64_t>(2);
  };
}

int64_t value_for_early_return(int64_t n) {
  int64_t i = static_cast<int64_t>(0);
  while ((i) < (n)) {
    if ((i) == (static_cast<int64_t>(5))) {
      return i;
      ;
    };
    i = (i) + (static_cast<int64_t>(1));
  };
  return static_cast<int64_t>(0);
}

int64_t match_exhaustive_value(::Choice t) {
  {
    auto __v_0 = t;
    switch (__v_0.index()) {
      case 0: {
        auto& x = std::get<0>(__v_0);
        return x;
      }
      case 1: {
        auto& y = std::get<1>(__v_0);
        return (y) + (static_cast<int64_t>(1));
      }
      default:
        __builtin_unreachable();
    }
  };
}

int64_t match_wildcard_value(::Choice t) {
  {
    auto __v_1 = t;
    switch (__v_1.index()) {
      case 0: {
        auto& x = std::get<0>(__v_1);
        return x;
      }
      default: {
        return static_cast<int64_t>(0);
      }
    }
  };
}

std::monostate unit_match_exhaustive(::Choice t) {
  {
    auto __v_2 = t;
    switch (__v_2.index()) {
      case 0: {
        auto& v = std::get<0>(__v_2);
        {
          return noop();
        };
        break;
      }
      case 1: {
        auto& v = std::get<1>(__v_2);
        {
          return noop();
        };
        break;
      }
      default:
        __builtin_unreachable();
    }
  };
  return std::monostate();
}
