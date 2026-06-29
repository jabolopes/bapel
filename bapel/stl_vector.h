#pragma once
#include <vector>
#include <cstdint>
#include <utility>
#include <variant>

// @bpl: pub type Vector ['a]
template <typename T>
using Vector = std::vector<T>;

// @bpl: pub Vector::mk: forall ['a] () -> Vector 'a
// @bpl: pub Vector::push_back: forall ['a] (&Vector 'a, 'a) -> ()
// @bpl: pub Vector::size: forall ['a] &Vector 'a -> i64
// @bpl: pub Vector::get: forall ['a] (&Vector 'a, i64) -> 'a
// @bpl: pub Vector::set: forall ['a] (&Vector 'a, i64, 'a) -> ()
template <typename T>
struct Vector_ {
  Vector_() = delete;

  static inline Vector<T> mk() {
    return Vector<T>();
  }

  static inline std::monostate push_back(Vector<T>* v, T value) {
    v->push_back(std::move(value));
    return std::monostate();
  }

  static inline int64_t size(Vector<T>* v) {
    return v->size();
  }

  static inline T get(Vector<T>* v, int64_t index) {
    return (*v)[index];
  }

  static inline void set(Vector<T>* v, int64_t index, T value) {
    (*v)[index] = std::move(value);
  }
};
