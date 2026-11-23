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
BAPEL_VARIANT(Expr, (lit, (int64_t)), (app, (__anonym_14999960583225247649)));
template <typename a>
BAPEL_VARIANT(Tree, (leaf, (a)),
              (branch,
               (std::tuple<::SharedPtr<::Tree<a>>, ::SharedPtr<::Tree<a>>>)));
BAPEL_VARIANT(IntList, (nil, (std::monostate)),
              (cons, (std::tuple<int64_t, ::SharedPtr<::IntList>>)));
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
