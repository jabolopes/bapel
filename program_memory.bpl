implements program

imports {
  bapel.core
  bapel.stl
}

pub fn testSharedPtr() -> i64 {
  let p: SharedPtr i64 = SharedPtr::make_shared 42;
  let val: i64 = *(p.get());
  let count: i64 = p.use_count();
  let null_ptr: SharedPtr i64 = SharedPtr::null();
  let is_n: bool = null_ptr.is_null();
  
  if is_n {
    return val + count
  }
  0
}

pub fn testUniquePtr() -> i64 {
  let u: UniquePtr i64 = UniquePtr::make_unique 100;
  let val: i64 = *(u.get());
  let null_u: UniquePtr i64 = UniquePtr::null();
  let is_n: bool = null_u.is_null();
  
  u.reset();
  let is_n_after_reset: bool = u.is_null();
  
  if is_n && is_n_after_reset {
    return val
  }
  0
}

