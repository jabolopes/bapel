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

template <typename a>
struct Maybe : std::variant<std::monostate /* none */, a /* some */> {};
struct One : std::variant<int64_t /* one */> {};
template <typename a>
::Maybe<a> mkNone();
::One mkOne();
template <typename a>
::Maybe<a> mkSome(a);
template <typename a>
::Maybe<a> mkNone() {
  ::Maybe<a> v = ::Maybe<a>{std::in_place_index<0>, std::monostate()};
  std::monostate v1 = std::get<0>(v);
  std::monostate v2 = std::get<0>(v);
  std::monostate v3;
  {
    auto __v_0 = v;
    switch (__v_0.index()) {
      case 0: {
        auto& l = std::get<0>(__v_0);
        v3 = l;
      }
      case 1: {
        auto& r = std::get<1>(__v_0);
        v3 = v1;
      }
    }
  };
  {
    auto __v_1 = v;
    switch (__v_1.index()) {
      case 0: {
        auto& l = std::get<0>(__v_1);
        l;
      }
      case 1: {
        auto& r = std::get<1>(__v_1);
        v1;
      }
    }
  };
  v = ::Maybe<a>{std::in_place_index<0>, std::monostate()};
  ::Maybe<a> r = v;
  return r;
}

template <typename a>
::Maybe<a> mkSome(a value) {
  ::Maybe<a> v = ::Maybe<a>{std::in_place_index<1>, value};
  a v1 = std::get<1>(v);
  a v2 = std::get<1>(v);
  a v3;
  {
    auto __v_2 = v;
    switch (__v_2.index()) {
      case 0: {
        auto& l = std::get<0>(__v_2);
        v3 = v2;
      }
      case 1: {
        auto& r = std::get<1>(__v_2);
        v3 = r;
      }
    }
  };
  {
    auto __v_3 = v;
    switch (__v_3.index()) {
      case 0: {
        auto& l = std::get<0>(__v_3);
        v2;
      }
      case 1: {
        auto& r = std::get<1>(__v_3);
        r;
      }
    }
  };
  v = ::Maybe<a>{std::in_place_index<1>, value};
  ::Maybe<a> r = v;
  return r;
}

::One mkOne();
