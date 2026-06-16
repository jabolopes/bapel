module bapel.core

impls {
  "core_array.h"
  "core_impl.h"
  "core_main.h"
  "core_optional.h"
  "core_math.h"
  "core_pointer.h"
  "core_string.h"
  "core_vector.h"
}

pub fn forCount(i: i64, f: () -> ()) -> () {
  if i == 0 {
    return ()
  }
  f ();
  forCount (i - 1, f)
}
