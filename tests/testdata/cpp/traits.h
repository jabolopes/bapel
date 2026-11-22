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
struct IntList
    : ::bapel::Variant<
          ::bapel::Arm<"nil", std::monostate>,
          ::bapel::Arm<"cons", std::tuple<int64_t, ::Ptr<::IntList>>>> {
  using ::bapel::Variant<
      ::bapel::Arm<"nil", std::monostate>,
      ::bapel::Arm<"cons", std::tuple<int64_t, ::Ptr<::IntList>>>>::Variant;
};
struct Color : ::bapel::Variant<::bapel::Arm<"red", std::monostate>,
                                ::bapel::Arm<"green", std::monostate>> {
  using ::bapel::Variant<::bapel::Arm<"red", std::monostate>,
                         ::bapel::Arm<"green", std::monostate>>::Variant;
};
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
struct Add;
}
namespace traits {
template <typename Self, typename elem>
struct Indexable;
}
namespace traits {
template <typename Self>
struct Size;
}
::Point make_point(int64_t x, int64_t y);
int64_t test_vector(::Ptr<::Vector<int64_t>> v);
int64_t test_point();
int64_t test_node(::Ptr<::Node> n);
int64_t test_color();
int64_t test_int_list(::Ptr<::IntList> l);
namespace traits {
template <typename a>
struct Indexable<::Vector<a>, a> {
  static inline a get(::Ptr<::Vector<a>> v, int64_t index) {
    return vector_get<a>(v, index);
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
struct Add<::Point> {
  static inline int64_t add(::Ptr<::Point> a, ::Ptr<::Point> b) {
    return (((::inherents::Ptr<::Point>::get(a).x) +
             (::inherents::Ptr<::Point>::get(a).y)) +
            (::inherents::Ptr<::Point>::get(b).x)) +
           (::inherents::Ptr<::Point>::get(b).y);
  }
};
}  // namespace traits
namespace traits {
template <>
struct Size<::Node> {
  static inline int64_t size(::Ptr<::Node> n) {
    return (::inherents::Ptr<::Node>::get(n).val) +
           (::traits::Size<::Node>::size(
               ::inherents::Ptr<::Node>::get(n).next));
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
          return (std::get<0>(pair)) +
                 (::traits::Size<::IntList>::size(std::get<1>(pair)));
        }
      }
    };
  }
};
}  // namespace traits