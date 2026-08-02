
#include "traits_private.h"

std::tuple<int64_t, int8_t> run(::MyStruct s, ::Ptr<::Vector<int8_t>> v) {
  return std::make_tuple(::traits::Size<::MyStruct>::size(s),
                         ::traits::Indexable<::Vector<int8_t>, int8_t>::get(
                             v, static_cast<int64_t>(0)));
}
