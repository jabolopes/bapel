#pragma once

#include "tuple.h"

struct Point : std::tuple<int8_t, int16_t> {
  Point(const std::tuple<int8_t, int16_t>& arg)
      : std::tuple<int8_t, int16_t>(arg) {}
};
Point mkPoint1();
Point mkPoint2();
Point mkPoint3();
std::tuple<int8_t, int16_t> mkTuple1();
std::tuple<int8_t, int16_t> mkTuple2();
std::tuple<int8_t, int16_t> mkTuple3();
std::tuple<int8_t, int16_t> mkTuple1();
std::tuple<int8_t, int16_t> mkTuple2();
std::tuple<int8_t, int16_t> mkTuple3();
Point mkPoint1();
Point mkPoint2();
Point mkPoint3();
