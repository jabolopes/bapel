#pragma once

#include <variant>

// @bpl: pub type ptr::Ptr ['a]
// @bpl: pub ptr::mk: forall ['a] 'a -> ptr::Ptr 'a
// @bpl: pub ptr::get: forall ['a] ptr::Ptr 'a -> 'a
// @bpl: pub ptr::set: forall ['a] (ptr::Ptr 'a, 'a) -> ()
namespace ptr {

template <typename A>
using Ptr = typename std::add_pointer<A>::type;

template <typename A>
Ptr<A> mk(A& a) { return &a; }

template <typename A>
A& get(Ptr<A> ptr) { return *ptr; }

template <typename A>
std::monostate set(Ptr<A> p, A a) {
  *p = std::move(a);
  return std::monostate();
}

}  // namespace ptr
