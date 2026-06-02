
#include "returns_private.h"

int8_t return1() { return static_cast<int8_t>(0); }

int8_t return2() {
  if (true) {
    return static_cast<int8_t>(0);
    ;
  };
  return static_cast<int8_t>(1);
}

int8_t return3() {
  if (true) {
    return static_cast<int8_t>(0);
    ;
  } else {
    return static_cast<int8_t>(1);
    ;
  };
  return static_cast<int8_t>(2);
}
