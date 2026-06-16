
#include "loops_private.h"

std::monostate testLoop() {
  int64_t i = static_cast<int64_t>(0);
  while ((i) < (static_cast<int64_t>(10))) {
    ::core::print<int8_t>(static_cast<int8_t>(1));
    i = (i) + (static_cast<int64_t>(1));
  };
  return std::monostate();
}
