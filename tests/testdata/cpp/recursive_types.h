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
    : std::variant<int64_t /* lit */, __anonym_14999960583225247649 /* app */> {
  using std::variant<int64_t /* lit */,
                     __anonym_14999960583225247649 /* app */>::variant;
};
template <typename a>
struct Tree : std::variant<a /* leaf */,
                           std::tuple<::SharedPtr<::Tree<a>>,
                                      ::SharedPtr<::Tree<a>>> /* branch */> {
  using std::variant<a /* leaf */,
                     std::tuple<::SharedPtr<::Tree<a>>,
                                ::SharedPtr<::Tree<a>>> /* branch */>::variant;
};
struct IntList
    : std::variant<std::monostate /* nil */,
                   std::tuple<int64_t, ::SharedPtr<::IntList>> /* cons */> {
  using std::variant<
      std::monostate /* nil */,
      std::tuple<int64_t, ::SharedPtr<::IntList>> /* cons */>::variant;
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
