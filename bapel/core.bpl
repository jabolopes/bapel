module bapel.core

impls {
  "core_array.ccm"
  "core_impl.ccm"
  "core_optional.ccm"
  "core_math.ccm"
  "core_pointer.ccm"
  "core_vector.ccm"
}

pub fn forCount(i: i64, f: () -> ()) -> () {
  if i == 0 {
    return ()
  }
  f ();
  forCount (i - 1, f)
}
