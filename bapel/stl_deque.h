#pragma once

#include <deque>
#include <utility>
#include <variant>

#include "bapel/core.h"

template <typename A>
using Deque = std::deque<A>;

// @bpl: pub DequeImpl::mk: forall ['a] () -> Deque 'a
// @bpl: pub DequeImpl::push_back: forall ['a] (&Deque 'a, 'a) -> ()
// @bpl: pub DequeImpl::pop_front: forall ['a] &Deque 'a -> ()
// @bpl: pub DequeImpl::front: forall ['a] &Deque 'a -> 'a
// @bpl: pub DequeImpl::empty: forall ['a] &Deque 'a -> bool
struct DequeImpl {
  DequeImpl() = delete;

  template <typename A>
  static inline Deque<A> mk() { return Deque<A>(); }

  template <typename A>
  static inline std::monostate push_back(Deque<A>* deque, A a) {
    deque->push_back(std::move(a));
    return std::monostate();
  }

  template <typename A>
  static inline std::monostate pop_front(Deque<A>* deque) {
    deque->pop_front();
    return std::monostate();
  }

  template <typename A>
  static inline A front(Deque<A>* deque) {
    return deque->front();
  }

  template <typename A>
  static inline bool empty(Deque<A>* deque) {
    return deque->empty();
  }
};
