
#include "structs_private.h"

__anonym_e283ac1df1bc153369a5d1a9c526e758514edfc3 mkStruct1() {
  return {.x = static_cast<int8_t>(1), .y = static_cast<int16_t>(2)};
}

__anonym_e283ac1df1bc153369a5d1a9c526e758514edfc3 mkStruct2() {
  return {.x = static_cast<int8_t>(1), .y = static_cast<int16_t>(2)};
}

__anonym_e283ac1df1bc153369a5d1a9c526e758514edfc3 mkStruct3() {
  __anonym_e283ac1df1bc153369a5d1a9c526e758514edfc3 r1 = {
      .x = static_cast<int8_t>(0), .y = static_cast<int16_t>(0)};
  __anonym_e283ac1df1bc153369a5d1a9c526e758514edfc3 r2 = {
      .x = static_cast<int8_t>(0), .y = static_cast<int16_t>(0)};
  int8_t x = r1.x;
  int16_t y = r1.y;
  ([&, __v_0 = r1]() mutable {
    __v_0.x = static_cast<int8_t>(3);
    __v_0.y = static_cast<int16_t>(4);
    return __v_0;
  })();
  r1 = ([&, __v_1 = r1]() mutable {
    __v_1.x = static_cast<int8_t>(3);
    __v_1.y = static_cast<int16_t>(4);
    return __v_1;
  })();
  ([&, __v_2 = r1]() mutable {
    __v_2.x = static_cast<int8_t>(3);
    __v_2.y = static_cast<int16_t>(4);
    return __v_2;
  })();
  r1 = ([&, __v_3 = r1]() mutable {
    __v_3.x = static_cast<int8_t>(3);
    __v_3.y = static_cast<int16_t>(4);
    return __v_3;
  })();
  __anonym_e283ac1df1bc153369a5d1a9c526e758514edfc3 r = r1;
  return r;
}

int8_t getStruct1(__anonym_a48d9280c7723679aaf5b528e9fe9a7a9f455500 r) {
  return r.x;
}

int16_t getStruct2(__anonym_e283ac1df1bc153369a5d1a9c526e758514edfc3 r) {
  return r.y;
}

Point mkPoint1() {
  return {.x = static_cast<int8_t>(1), .y = static_cast<int16_t>(2)};
}

Point mkPoint2() {
  return {.x = static_cast<int8_t>(1), .y = static_cast<int16_t>(2)};
}

Point mkPoint3() {
  Point r1 = {.x = static_cast<int8_t>(0), .y = static_cast<int16_t>(0)};
  Point r2 = {.x = static_cast<int8_t>(0), .y = static_cast<int16_t>(0)};
  int8_t x = r1.x;
  int16_t y = r1.y;
  ([&, __v_4 = r1]() mutable {
    __v_4.x = static_cast<int8_t>(3);
    __v_4.y = static_cast<int16_t>(4);
    return __v_4;
  })();
  r1 = ([&, __v_5 = r1]() mutable {
    __v_5.x = static_cast<int8_t>(3);
    __v_5.y = static_cast<int16_t>(4);
    return __v_5;
  })();
  ([&, __v_6 = r1]() mutable {
    __v_6.x = static_cast<int8_t>(3);
    __v_6.y = static_cast<int16_t>(4);
    return __v_6;
  })();
  r1 = ([&, __v_7 = r1]() mutable {
    __v_7.x = static_cast<int8_t>(3);
    __v_7.y = static_cast<int16_t>(4);
    return __v_7;
  })();
  Point r = r1;
  return r;
}

int8_t getPoint1(Point p) { return p.x; }

int16_t getPoint2(Point p) { return p.y; }
