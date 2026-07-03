#pragma once

#include <deque>
#include <variant>

#include "bapel/core.h"

// @bpl: pub type Deque ['a]

template <typename A>
using Deque = std::deque<A>;

// @bpl: pub Deque::mk: forall ['a] () -> Deque 'a
// @bpl: pub Deque::push_back: forall ['a] (&Deque 'a, 'a) -> ()
// @bpl: pub Deque::pop_front: forall ['a] &Deque 'a -> ()
// @bpl: pub Deque::front: forall ['a] &Deque 'a -> 'a
// @bpl: pub Deque::empty: forall ['a] &Deque 'a -> bool
namespace inherents {
template <typename A>
struct Deque {
  Deque() = delete;

  static inline ::Deque<A> mk() { return ::Deque<A>(); }

  static inline std::monostate push_back(::Deque<A>* deque, A a) {
    deque->push_back(std::move(a));
    return std::monostate();
  }

  static inline std::monostate pop_front(::Deque<A>* deque) {
    deque->pop_front();
    return std::monostate();
  }

  static inline A front(::Deque<A>* deque) {
    return deque->front();
  }

  static inline bool empty(::Deque<A>* deque) {
    return deque->empty();
  }
};
}
