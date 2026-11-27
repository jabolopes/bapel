
#include "bapel/core.h"
#include "bapel/stl.h"
#include "recursive_types_private.h"

::Node mkNode(int64_t v) {
  return {.val = v, .next = ::inherents::SharedPtr<::Node>::null()};
}

int64_t sumList(::Ptr<::IntList> l) {
  {
    auto __v_0 = ::inherents::Ptr<::IntList>::get(l);
    switch (__v_0.index()) {
      case 0: {
        auto& u = std::get<0>(__v_0);
        return static_cast<int64_t>(0);
      }
      case 1: {
        auto& pair = std::get<1>(__v_0);
        return (std::get<0>(pair)) +
               (sumList(::inherents::SharedPtr<::IntList>::get(
                   ::inherents::Ptr<::SharedPtr<::IntList>>::mk(
                       std::get<1>(pair)))));
      }
      default:
        __builtin_unreachable();
    }
  };
}

::Expr mkApp(::SharedPtr<::Expr> f, ::SharedPtr<::Expr> a) {
  return ::Expr(std::in_place_index<1>,
                __anonym_9642821253775059453{.fun = f, .arg = a});
}
