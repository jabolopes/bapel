#pragma once
#include <algorithm>
#include <vector>
#include <cstdint>
#include <utility>
#include <variant>

template <typename T>
using Vector = std::vector<T>;

// @bpl: pub VectorImpl::mk: forall ['a] () -> Vector 'a
// @bpl: pub VectorImpl::push_back: forall ['a] (&Vector 'a, 'a) -> ()
// @bpl: pub VectorImpl::size: forall ['a] &Vector 'a -> i64
// @bpl: pub VectorImpl::get: forall ['a] (&Vector 'a, i64) -> 'a
// @bpl: pub VectorImpl::set: forall ['a] (&Vector 'a, i64, 'a) -> ()
// @bpl: pub VectorImpl::sort: forall ['a] &Vector 'a -> ()
// @bpl: pub VectorImpl::dedup: forall ['a] &Vector 'a -> ()
struct VectorImpl {
  VectorImpl() = delete;

  template <typename T>
  static inline Vector<T> mk() {
    return Vector<T>();
  }

  template <typename T>
  static inline std::monostate push_back(Vector<T>* v, T value) {
    v->push_back(std::move(value));
    return std::monostate();
  }

  template <typename T>
  static inline int64_t size(Vector<T>* v) {
    return v->size();
  }

  template <typename T>
  static inline T get(Vector<T>* v, int64_t index) {
    return (*v)[index];
  }

  template <typename T>
  static inline std::monostate set(Vector<T>* v, int64_t index, T value) {
    (*v)[index] = std::move(value);
    return std::monostate();
  }

  template <typename T>
  static inline std::monostate sort(Vector<T>* v) {
    std::sort(v->begin(), v->end());
    return std::monostate();
  }

  template <typename T>
  static inline std::monostate dedup(Vector<T>* v) {
    v->erase(std::unique(v->begin(), v->end()), v->end());
    return std::monostate();
  }
};


