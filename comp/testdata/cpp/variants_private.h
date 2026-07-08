#pragma once

#include "variants.h"

template <typename a>
struct Choice : std::variant<a /* left */, int8_t /* right */> {
  using std::variant<a /* left */, int8_t /* right */>::variant;
};
template <typename a>
::Choice<a> mkLeft(a);
template <typename a>
::Choice<a> mkRight(int8_t);
template <typename a>
::Choice<a> mkLeft(a value) {
  ::Choice<a> v = ::Choice<a>(std::in_place_index<0>, value);
  a v1 = std::get<0>(v);
  a v2 = std::get<0>(v);
  a v3;
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
  v = ::Choice<a>(std::in_place_index<0>, value);
  ::Choice<a> r = v;
  return r;
}

template <typename a>
::Choice<a> mkRight(int8_t value) {
  ::Choice<a> v = ::Choice<a>(std::in_place_index<1>, value);
  int8_t v1 = std::get<1>(v);
  int8_t v2 = std::get<1>(v);
  int8_t v3;
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
  v = ::Choice<a>(std::in_place_index<1>, value);
  ::Choice<a> r = v;
  return r;
}
