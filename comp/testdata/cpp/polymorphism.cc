
#include "polymorphism_private.h"

std::monostate callPolymorphic() {
  return ::core::print<int8_t>(static_cast<int8_t>(1));
}

std::monostate functionSubtyping() {
  std::function<int8_t(int8_t)> id2 = id<int8_t>;
  return std::monostate();
}
