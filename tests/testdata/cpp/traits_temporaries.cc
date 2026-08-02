
#include "traits_temporaries_private.h"

::S make_s(int64_t val) { return {.x = val}; }

int64_t run() {
  int64_t a = printSize<::S>(
      ::inherents::Ptr<::S>::mk(make_s(static_cast<int64_t>(10))));
  int64_t b = ::traits::Size<::S>::size(
      ::inherents::Ptr<::S>::mk(make_s(static_cast<int64_t>(20))));
  return (a) + (b);
}
