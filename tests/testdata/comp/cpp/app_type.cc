
#include "app_type_private.h"

std::monostate typeApplicativeConst() {
  static_cast<int8_t>(1);
  return std::monostate();
}
