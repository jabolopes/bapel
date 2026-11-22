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
int64_t test_vector(::Ptr<::Vector<int64_t>> v);
int64_t test_point();
int64_t test_node(::Ptr<::Node> n);
bool test_color();
int64_t test_int_list(::Ptr<::IntList> l);
namespace inherents {
template <typename a>
struct Vector {
  Vector() = delete;
  static inline a get(::Ptr<::Vector<a>> v, int64_t i) {
    return vector_get<a>(v, i);
  }
};
}  // namespace inherents
namespace inherents {
struct Point {
  Point() = delete;
  static inline ::Point mk(int64_t x, int64_t y) { return {.x = x, .y = y}; }

  static inline int64_t sum(::Ptr<::Point> p) {
    return (::inherents::Ptr<::Point>::get(p).x) +
           (::inherents::Ptr<::Point>::get(p).y);
  }

  static inline ::Point add(::Ptr<::Point> a, ::Ptr<::Point> b) {
    return {.x = (::inherents::Ptr<::Point>::get(a).x) +
                 (::inherents::Ptr<::Point>::get(b).x),
            .y = (::inherents::Ptr<::Point>::get(a).y) +
                 (::inherents::Ptr<::Point>::get(b).y)};
  }
};
}  // namespace inherents
namespace inherents {
struct Node {
  Node() = delete;
  static inline int64_t get_val(::Ptr<::Node> n) {
    return ::inherents::Ptr<::Node>::get(n).val;
  }

  static inline int64_t next_val(::Ptr<::Node> n) {
    return ::inherents::Node::get_val(::inherents::Ptr<::Node>::get(n).next);
  }
};
}  // namespace inherents
namespace inherents {
struct Color {
  Color() = delete;
  static inline bool is_red(::Ptr<::Color> c) {
    {
      auto __v_0 = ::inherents::Ptr<::Color>::get(c);
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
  static inline bool is_nil(::Ptr<::IntList> l) {
    {
      auto __v_1 = ::inherents::Ptr<::IntList>::get(l);
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

  static inline int64_t sum(::Ptr<::IntList> l) {
    {
      auto __v_2 = ::inherents::Ptr<::IntList>::get(l);
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