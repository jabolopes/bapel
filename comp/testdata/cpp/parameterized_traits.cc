
#include "parameterized_traits_private.h"

int8_t run(Ptr<Vector<int8_t>> v) {
  return ::traits::Indexable<Vector<int8_t>, int8_t>::get(
      v, static_cast<int64_t>(0));
}
