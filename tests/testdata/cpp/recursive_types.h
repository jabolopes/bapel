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

struct Node;
struct IntList;
template <typename a>
struct Tree;
struct Node {
  int64_t val;
  ::SharedPtr<::Node> next;
};
struct IntList
    : std::variant<std::monostate /* nil */,
                   std::tuple<int64_t, ::SharedPtr<::IntList>> /* cons */> {
  using std::variant<
      std::monostate /* nil */,
      std::tuple<int64_t, ::SharedPtr<::IntList>> /* cons */>::variant;
};
template <typename a>
struct Tree : std::variant<a /* leaf */,
                           std::tuple<::SharedPtr<::Tree<a>>,
                                      ::SharedPtr<::Tree<a>>> /* branch */> {
  using std::variant<a /* leaf */,
                     std::tuple<::SharedPtr<::Tree<a>>,
                                ::SharedPtr<::Tree<a>>> /* branch */>::variant;
};
::Node mkNode(int64_t v);
int64_t sumList(::Ptr<::IntList> l);
template <typename a>
::Tree<a> mkLeaf(a x);
template <typename a>
::Tree<a> mkLeaf(a x) {
  return ::Tree<a>(std::in_place_index<0>, x);
}
