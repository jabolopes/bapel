
#include "traits_private.h"

int64_t run(MyStruct s) { return ::traits::Size<MyStruct>::size(s); }
