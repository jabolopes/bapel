
#include "traits_bounds_ok_private.h"

int64_t run() {
  S s = {.x = static_cast<int64_t>(42)};
  return printSize<S>(::Ptr_<S>::mk(s));
}
