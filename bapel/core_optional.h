#pragma once

#include <optional>

// @bpl: pub type std::optional ['a]
// @bpl: pub std::make_optional: forall ['a] 'a -> (std::optional 'a)

// @bpl: pub none: forall ['a] () -> (std::optional 'a)
template <typename A>
std::optional<A> none() {
  return std::nullopt;
}

// @bpl: pub has_value: forall ['a] std::optional 'a -> bool
template <typename A>
bool has_value(const std::optional<A>& opt) {
  return opt.has_value();
}

// @bpl: pub get_value: forall ['a] std::optional 'a -> 'a
template <typename A>
A get_value(const std::optional<A>& opt) {
  return *opt;
}
