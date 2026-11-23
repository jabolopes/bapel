#pragma once

#include "structs.h"

#ifndef __anonym_6676189813872021938_DEFINED
#define __anonym_6676189813872021938_DEFINED
struct __anonym_6676189813872021938 {
  int8_t x;
  int16_t y;
};
#endif
#ifndef __anonym_17929389288086054735_DEFINED
#define __anonym_17929389288086054735_DEFINED
struct __anonym_17929389288086054735 {
  int8_t x;
};
#endif
struct Point {
  int8_t x;
  int16_t y;
};
__anonym_6676189813872021938 mkStruct1();
__anonym_6676189813872021938 mkStruct2();
__anonym_6676189813872021938 mkStruct3();
int8_t getStruct1(__anonym_17929389288086054735 r);
int16_t getStruct2(__anonym_6676189813872021938 r);
::Point mkPoint1();
::Point mkPoint2();
::Point mkPoint3();
int8_t getPoint1(::Point p);
int16_t getPoint2(::Point p);
