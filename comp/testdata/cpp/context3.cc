
#include "context3_private.h"

std::monostate x() { return std::monostate(); }

int8_t y(int8_t x) { return x; }

int16_t z() {
  int16_t x = static_cast<int16_t>(0);
  return x;
}
