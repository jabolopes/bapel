#pragma once
template <typename T>
struct Ptr {
  T* ptr;
};

// @bpl: pub Ptr::mk: forall ['a] 'a -> Ptr 'a
// @bpl: pub Ptr::get: forall ['a] Ptr 'a -> 'a
