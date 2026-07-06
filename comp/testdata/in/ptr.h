#pragma once
#include <type_traits>

template <typename A>
using Ptr = typename std::add_pointer<A>::type;

// @bpl: pub Ptr::mk: forall ['a] 'a -> Ptr 'a
// @bpl: pub Ptr::get: forall ['a] Ptr 'a -> 'a
namespace inherents {
template <typename A>
struct Ptr {
  Ptr() = delete;
  static inline ::Ptr<A> mk(A& a) { return &a; }
  static inline ::Ptr<A> mk(const A& a) { return const_cast<A*>(&a); }
  static inline ::Ptr<A> mk(A&& a) { return &a; }
  static inline A& get(::Ptr<A> ptr) { return *ptr; }
};
}
