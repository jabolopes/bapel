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

#include "tests/testdata/in/ptr.h"
#include "tests/testdata/in/vector.h"

struct IntList;
struct Node;
struct IntList
    : std::variant<std::monostate /* nil */,
                   std::tuple<int64_t, ::Ptr<::IntList>> /* cons */> {
  using std::variant<std::monostate /* nil */,
                     std::tuple<int64_t, ::Ptr<::IntList>> /* cons */>::variant;
};
struct Color
    : std::variant<std::monostate /* red */, std::monostate /* green */> {
  using std::variant<std::monostate /* red */,
                     std::monostate /* green */>::variant;
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
  using Self = ::Vector<a>;
  static inline a get(::Ptr<Self> v, int64_t index) {
    return vector_get<a>(v, index);
  }
};
}  // namespace traits
namespace traits {
template <>
struct Size<::Point> {
  using Self = ::Point;
  static inline int64_t size(::Ptr<Self> p) {
    return (::inherents::Ptr<Self>::get(p).x) +
           (::inherents::Ptr<Self>::get(p).y);
  }
};
}  // namespace traits
namespace traits {
template <>
struct Add<::Point> {
  using Self = ::Point;
  static inline int64_t add(::Ptr<Self> a, ::Ptr<Self> b) {
    return (((::inherents::Ptr<Self>::get(a).x) +
             (::inherents::Ptr<Self>::get(a).y)) +
            (::inherents::Ptr<Self>::get(b).x)) +
           (::inherents::Ptr<Self>::get(b).y);
  }
};
}  // namespace traits
namespace traits {
template <>
struct Size<::Node> {
  using Self = ::Node;
  static inline int64_t size(::Ptr<Self> n) {
    return (::inherents::Ptr<Self>::get(n).val) +
           (::traits::Size<::Node>::size(::inherents::Ptr<Self>::get(n).next));
  }
};
}  // namespace traits
namespace traits {
template <>
struct Size<::Color> {
  using Self = ::Color;
  static inline int64_t size(::Ptr<Self> c) {
    {
      auto __v_0 = ::inherents::Ptr<Self>::get(c);
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
  using Self = ::IntList;
  static inline int64_t size(::Ptr<Self> l) {
    {
      auto __v_1 = ::inherents::Ptr<Self>::get(l);
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