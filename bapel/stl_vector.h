#pragma once
#include <vector>
#include <cstdint>
#include <utility>

// @bpl: pub type Vector ['a]
template <typename T>
using Vector = std::vector<T>;

// @bpl: pub Vector_::mk: forall ['a] () -> Vector 'a
// @bpl: pub Vector_::push_back: forall ['a] (& (Vector 'a), 'a) -> ()
// @bpl: pub Vector_::size: forall ['a] & (Vector 'a) -> i64
// @bpl: pub Vector_::get: forall ['a] (& (Vector 'a), i64) -> 'a
// @bpl: pub Vector_::set: forall ['a] (& (Vector 'a), i64, 'a) -> ()
namespace Vector_ {

template <typename T>
inline Vector<T> mk() {
  return Vector<T>();
}

template <typename T>
inline void push_back(Vector<T>* v, const T& value) {
  v->push_back(value);
}

template <typename T>
inline int64_t size(Vector<T>* v) {
  return v->size();
}

template <typename T>
inline T get(Vector<T>* v, int64_t index) {
  return (*v)[index];
}

template <typename T>
inline void set(Vector<T>* v, int64_t index, T value) {
  (*v)[index] = std::move(value);
}

}  // namespace Vector_
