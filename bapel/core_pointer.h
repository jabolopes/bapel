#pragma once

#include <variant>

// @bpl: pub type Ptr ['a]
template <typename A>
using Ptr = typename std::add_pointer<A>::type;

// @bpl: pub Ptr_::mk: forall ['a] 'a -> Ptr 'a
// @bpl: pub Ptr_::get: forall ['a] Ptr 'a -> 'a
// @bpl: pub Ptr_::set: forall ['a] (Ptr 'a, 'a) -> ()
namespace Ptr_ {

template <typename A>
Ptr<A> mk(A& a) { return &a; }

template <typename A>
A& get(Ptr<A> ptr) { return *ptr; }

template <typename A>
std::monostate set(Ptr<A> p, A a) {
  *p = std::move(a);
  return std::monostate();
}

}  // namespace Ptr_
