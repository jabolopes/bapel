#pragma once
#include <optional>
#include <variant>

template <typename T>
using Optional = std::optional<T>;

// @bpl: pub OptionalImpl::none: forall ['a] () -> Optional 'a
// @bpl: pub OptionalImpl::make_optional: forall ['a] 'a -> Optional 'a
// @bpl: pub OptionalImpl::has_value: forall ['a] &Optional 'a -> bool
// @bpl: pub OptionalImpl::get_value: forall ['a] &Optional 'a -> 'a
struct OptionalImpl {
  OptionalImpl() = delete;

  template <typename T>
  static inline Optional<T> none() {
    return std::nullopt;
  }

  template <typename T>
  static inline Optional<T> make_optional(T val) {
    return std::make_optional(std::move(val));
  }

  template <typename T>
  static inline bool has_value(const Optional<T>* opt) {
    return opt->has_value();
  }

  template <typename T>
  static inline T get_value(const Optional<T>* opt) {
    return **opt;
  }
};
