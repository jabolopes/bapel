#pragma once

#include "traits_temporaries.h"

struct S {
  int64_t x;
};
namespace traits {
template <typename Self>
struct Size;
}
::S make_s(int64_t val);
template <typename t,
          typename = std::enable_if_t<(sizeof(::traits::Size<t>) > 0)>>
int64_t printSize(::Ptr<t> x);
int64_t run();
template <typename t, typename>
int64_t printSize(::Ptr<t> x) {
  return ::traits::Size<t>::size(x);
}

namespace traits {
template <>
struct Size<::S> {
  static inline int64_t size(::Ptr<::S> s) {
    return ::inherents::Ptr<::S>::get(s).x;
  }
};
}  // namespace traits