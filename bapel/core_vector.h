#pragma once

#include <cstdint>
#include <vector>

namespace core {

// @bpl: pub type core::Vector ['a]
template <typename T>
using Vector = std::vector<T>;

// @bpl: pub core::mk: () -> core::Vector 'a
template <typename T>
Vector<T> mk() {
  return Vector<T>();
}

// @bpl: pub core::add: forall ['a] (& (core::Vector 'a), 'a) -> ()
template <typename T>
void add(Vector<T>* v, const T& value) {
  v->push_back(value);
}

// @bpl: pub core::get: forall ['a] (& (core::Vector 'a), i64) -> 'a
template <typename T>
T get(Vector<T>* v, int64_t index) {
  return (*v)[index];
}

// @bpl: pub core::set: forall ['a] (& (core::Vector 'a), i64, 'a) -> ()
template <typename T>
void set(Vector<T>* v, int64_t index, T&& value) {
  (*v)[index] = std::move(value);
}

// @bpl: pub core::vec_copy: forall ['a] (& (core::Vector 'a)) -> core::Vector 'a
template <typename T>
Vector<T> vec_copy(Vector<T>* v) {
  return *v;
}

// @bpl: pub core::vec_size: forall ['a] & (core::Vector 'a) -> i64
template <typename T>
int64_t vec_size(Vector<T>* v) {
  return v->size();
}

}  // namespace core
