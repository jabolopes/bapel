
#include "loops_private.h"

std::monostate testLoop() {
  forCount(static_cast<int64_t>(10),
           [&]() { return ::core::print<int8_t>(static_cast<int8_t>(1)); });
  return std::monostate();
}
