
#include "anonymous_structs_private.h"

int64_t eval(::Expr e) {
  {
    auto __v_0 = e;
    switch (__v_0.index()) {
      case 0: {
        auto& x = std::get<0>(__v_0);
        return x;
      }
      case 1: {
        auto& p = std::get<1>(__v_0);
        return (p.lhs) + (p.rhs);
      }
      default:
        __builtin_unreachable();
    }
  };
}

int64_t sum_point(__anonym_11548725740611542807 p) { return (p.x) + (p.y); }

::Expr make_expr() {
  return ::Expr(std::in_place_index<1>,
                __anonym_4591088945429812398{.lhs = static_cast<int64_t>(1),
                                             .rhs = static_cast<int64_t>(2)});
}
