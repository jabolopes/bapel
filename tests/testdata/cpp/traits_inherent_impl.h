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

struct Node;
struct IntList;
struct Point {
  int64_t x;
  int64_t y;
};
struct Node {
  int64_t val;
  ::Ptr<::Node> next;
};
struct Color
    : std::variant<std::monostate /* red */, std::monostate /* green */> {
  using std::variant<std::monostate /* red */,
                     std::monostate /* green */>::variant;
};
struct IntList
    : std::variant<std::monostate /* nil */,
                   std::tuple<int64_t, ::Ptr<::IntList>> /* cons */> {
  using std::variant<std::monostate /* nil */,
                     std::tuple<int64_t, ::Ptr<::IntList>> /* cons */>::variant;
};
int64_t test_vector(::Ptr<::Vector<int64_t>> v);
int64_t test_point();
int64_t test_node(::Ptr<::Node> n);
bool test_color();
int64_t test_int_list(::Ptr<::IntList> l);
namespace inherents {
template <typename a>
struct Vector {
  Vector() = delete;
  using Self = ::Vector<a>;
  static inline a get(::Ptr<Self> v, int64_t i) { return vector_get<a>(v, i); }
};
}  // namespace inherents
namespace inherents {
struct Point {
  Point() = delete;
  using Self = ::Point;
  static inline ::Point mk(int64_t x, int64_t y) { return {.x = x, .y = y}; }

  static inline int64_t sum(::Ptr<Self> p) {
    return (::inherents::Ptr<Self>::get(p).x) +
           (::inherents::Ptr<Self>::get(p).y);
  }

  static inline ::Point add(::Ptr<Self> a, ::Ptr<::Point> b) {
    return {.x = (::inherents::Ptr<Self>::get(a).x) +
                 (::inherents::Ptr<::Point>::get(b).x),
            .y = (::inherents::Ptr<Self>::get(a).y) +
                 (::inherents::Ptr<::Point>::get(b).y)};
  }
};
}  // namespace inherents
namespace inherents {
struct Node {
  Node() = delete;
  using Self = ::Node;
  static inline int64_t get_val(::Ptr<Self> n) {
    return ::inherents::Ptr<Self>::get(n).val;
  }

  static inline int64_t next_val(::Ptr<Self> n) {
    return ::inherents::Node::get_val(::inherents::Ptr<Self>::get(n).next);
  }
};
}  // namespace inherents
namespace inherents {
struct Color {
  Color() = delete;
  using Self = ::Color;
  static inline bool is_red(::Ptr<Self> c) {
    {
      auto __v_0 = ::inherents::Ptr<Self>::get(c);
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
};
}  // namespace inherents
namespace inherents {
struct IntList {
  IntList() = delete;
  using Self = ::IntList;
  static inline bool is_nil(::Ptr<Self> l) {
    {
      auto __v_1 = ::inherents::Ptr<Self>::get(l);
      switch (__v_1.index()) {
        case 0: {
          auto& u = std::get<0>(__v_1);
          return true;
        }
        default: {
          return false;
        }
      }
    };
  }

  static inline int64_t sum(::Ptr<Self> l) {
    {
      auto __v_2 = ::inherents::Ptr<Self>::get(l);
      switch (__v_2.index()) {
        case 0: {
          auto& u = std::get<0>(__v_2);
          return static_cast<int64_t>(0);
        }
        case 1: {
          auto& pair = std::get<1>(__v_2);
          return (std::get<0>(pair)) +
                 (::inherents::IntList::sum(std::get<1>(pair)));
        }
      }
    };
  }
};
}  // namespace inherents