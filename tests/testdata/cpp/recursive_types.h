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

#include "bapel/core.h"
#include "bapel/core_variant.h"
#include "bapel/stl.h"

struct Expr;
template <typename a>
struct Tree;
struct IntList;
struct Node;
struct __anonym_14999960583225247649 {
  ::SharedPtr<::Expr> fun;
  ::SharedPtr<::Expr> arg;
};
struct Expr
    : ::bapel::Variant<::bapel::Arm<"lit", int64_t>,
                       ::bapel::Arm<"app", __anonym_14999960583225247649>> {
  using ::bapel::Variant<
      ::bapel::Arm<"lit", int64_t>,
      ::bapel::Arm<"app", __anonym_14999960583225247649>>::Variant;
};
template <typename a>
struct Tree : ::bapel::Variant<
                  ::bapel::Arm<"leaf", a>,
                  ::bapel::Arm<"branch", std::tuple<::SharedPtr<::Tree<a>>,
                                                    ::SharedPtr<::Tree<a>>>>> {
  using ::bapel::Variant<
      ::bapel::Arm<"leaf", a>,
      ::bapel::Arm<"branch", std::tuple<::SharedPtr<::Tree<a>>,
                                        ::SharedPtr<::Tree<a>>>>>::Variant;
};
struct IntList
    : ::bapel::Variant<
          ::bapel::Arm<"nil", std::monostate>,
          ::bapel::Arm<"cons", std::tuple<int64_t, ::SharedPtr<::IntList>>>> {
  using ::bapel::Variant<
      ::bapel::Arm<"nil", std::monostate>,
      ::bapel::Arm<"cons",
                   std::tuple<int64_t, ::SharedPtr<::IntList>>>>::Variant;
};
struct Node {
  int64_t val;
  ::SharedPtr<::Node> next;
};
::Node mkNode(int64_t v);
int64_t sumList(::Ptr<::IntList> l);
template <typename a>
::Tree<a> mkLeaf(a x);
::Expr mkApp(::SharedPtr<::Expr> f, ::SharedPtr<::Expr> a);
template <typename a>
::Tree<a> mkLeaf(a x) {
  return ::Tree<a>(std::in_place_index<0>, x);
}
