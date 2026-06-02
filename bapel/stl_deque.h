#pragma once

#include <deque>
#include <variant>

#include "bapel/core.h"

// @bpl: pub type Deque ['a]

template <typename A>
using Deque = std::deque<A>;

// @bpl: pub Deque_::mk: forall ['a] () -> Deque 'a
// @bpl: pub Deque_::push_back: forall ['a] (&Deque 'a, 'a) -> ()
namespace Deque_ {

template <typename A>
Deque<A> mk() { return Deque<A>(); }

template <typename A>
std::monostate push_back(Deque<A>* deque, A a) {
	deque->push_back(std::move(a));
	return std::monostate();
}

}  // namespace Deque_
