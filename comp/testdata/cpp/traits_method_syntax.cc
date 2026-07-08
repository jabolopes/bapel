
#include "traits_method_syntax_private.h"

::S make_s(int64_t val) { return {.x = val}; }

int64_t run() {
  int64_t a = ::traits::Size<::S>::size(
      ::inherents::Ptr<::S>::mk(make_s(static_cast<int64_t>(10))));
  ::S s1 = make_s(static_cast<int64_t>(20));
  int64_t b = ::traits::Size<::S>::size(::inherents::Ptr<::S>::mk(s1));
  ::Ptr<::S> ref_s1 = ::inherents::Ptr<::S>::mk(s1);
  int64_t c = ::traits::Size<::S>::size(ref_s1);
  int64_t d = printSize<::S>(::inherents::Ptr<::S>::mk(s1));
  ::S s2 = make_s(static_cast<int64_t>(30));
  int64_t e = ::traits::Add<::S>::add(::inherents::Ptr<::S>::mk(s1),
                                      ::inherents::Ptr<::S>::mk(s2));
  int64_t f = ::traits::Size<::S>::size(ref_s1);
  return (((((a) + (b)) + (c)) + (d)) + (e)) + (f);
}
