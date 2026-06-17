#pragma once

#include <deque>
#include <variant>

#include "bapel/core.h"

// @bpl: pub type Deque ['a]

template <typename A>
using Deque = std::deque<A>;

// @bpl: pub Deque_::mk: forall ['a] () -> Deque 'a
// @bpl: pub Deque_::push_back: forall ['a] (&Deque 'a, 'a) -> ()
// @bpl: pub Deque_::pop_front: forall ['a] &Deque 'a -> ()
// @bpl: pub Deque_::front: forall ['a] &Deque 'a -> 'a
// @bpl: pub Deque_::empty: forall ['a] &Deque 'a -> bool
namespace Deque_ {

template <typename A>
Deque<A> mk() { return Deque<A>(); }

template <typename A>
std::monostate push_back(Deque<A>* deque, A a) {
	deque->push_back(std::move(a));
	return std::monostate();
}

template <typename A>
std::monostate pop_front(Deque<A>* deque) {
	deque->pop_front();
	return std::monostate();
}

template <typename A>
A front(Deque<A>* deque) {
	return deque->front();
}

template <typename A>
bool empty(Deque<A>* deque) {
	return deque->empty();
}

}  // namespace Deque_
