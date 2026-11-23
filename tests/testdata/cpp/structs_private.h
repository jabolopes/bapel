#pragma once

#include "structs.h"

#ifndef __anonym_1564688777576866135_DEFINED
#define __anonym_1564688777576866135_DEFINED
struct __anonym_1564688777576866135 {
  int8_t x;
  int16_t y;
};
#endif
#ifndef __anonym_14667179875479736764_DEFINED
#define __anonym_14667179875479736764_DEFINED
struct __anonym_14667179875479736764 {
  int8_t x;
};
#endif
struct Point {
  int8_t x;
  int16_t y;
};
__anonym_1564688777576866135 mkStruct1();
__anonym_1564688777576866135 mkStruct2();
__anonym_1564688777576866135 mkStruct3();
int8_t getStruct1(__anonym_14667179875479736764 r);
int16_t getStruct2(__anonym_1564688777576866135 r);
::Point mkPoint1();
::Point mkPoint2();
::Point mkPoint3();
int8_t getPoint1(::Point p);
int16_t getPoint2(::Point p);
