#pragma once

#include <array>
#include <cmath>
#include <cstdlib>
#include <functional>
#include <optional>
#include <string>
#include <tuple>
#include <variant>
#include <vector>

#include "bapel/core_variant.h"
#include "tests/testdata/in/ptr.h"
#include "tests/testdata/in/vector.h"

struct IntList;
struct Node;
BAPEL_VARIANT(IntList, (nil, (std::monostate)),
              (cons, (std::tuple<int64_t, ::Ptr<::IntList>>)));
BAPEL_VARIANT(Color, (red, (std::monostate)), (green, (std::monostate)));
struct Node {
  int64_t val;
  ::Ptr<::Node> next;
};
struct Point {
  int64_t x;
  int64_t y;
};
namespace traits {
template <typename Self>
struct Size;
}
template <typename t,
          typename = std::enable_if_t<(sizeof(::traits::Size<t>) > 0)>>
int64_t print_size(::Ptr<t> x);
template <typename t,
          typename = std::enable_if_t<(sizeof(::traits::Size<t>) > 0)>>
int64_t print_size_dot(::Ptr<t> x);
int64_t test_vector(::Ptr<::Vector<int64_t>> v);
int64_t test_point();
int64_t test_node(::Ptr<::Node> n);
int64_t test_color();
int64_t test_int_list(::Ptr<::IntList> l);
template <typename t, typename>
int64_t print_size(::Ptr<t> x) {
  return ::traits::Size<t>::size(x);
}

template <typename t, typename>
int64_t print_size_dot(::Ptr<t> x) {
  return ::traits::Size<t>::size(x);
}

namespace traits {
template <typename a>
struct Size<::Vector<a>> {
  static inline int64_t size(::Ptr<::Vector<a>> v) {
    return static_cast<int64_t>(0);
  }
};
}  // namespace traits
namespace traits {
template <>
struct Size<::Point> {
  static inline int64_t size(::Ptr<::Point> p) {
    return (::inherents::Ptr<::Point>::get(p).x) +
           (::inherents::Ptr<::Point>::get(p).y);
  }
};
}  // namespace traits
namespace traits {
template <>
struct Size<::Node> {
  static inline int64_t size(::Ptr<::Node> n) {
    return ::inherents::Ptr<::Node>::get(n).val;
  }
};
}  // namespace traits
namespace traits {
template <>
struct Size<::Color> {
  static inline int64_t size(::Ptr<::Color> c) {
    {
      auto __v_0 = ::inherents::Ptr<::Color>::get(c);
      switch (__v_0.index()) {
        case 0: {
          auto& u = std::get<0>(__v_0);
          return static_cast<int64_t>(1);
        }
        default: {
          return static_cast<int64_t>(2);
        }
      }
    };
  }
};
}  // namespace traits
namespace traits {
template <>
struct Size<::IntList> {
  static inline int64_t size(::Ptr<::IntList> l) {
    {
      auto __v_1 = ::inherents::Ptr<::IntList>::get(l);
      switch (__v_1.index()) {
        case 0: {
          auto& u = std::get<0>(__v_1);
          return static_cast<int64_t>(0);
        }
        case 1: {
          auto& pair = std::get<1>(__v_1);
          return std::get<0>(pair);
        }
        default:
          __builtin_unreachable();
      }
    };
  }
};
}  // namespace traits