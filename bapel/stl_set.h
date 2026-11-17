#pragma once
#include <cstdint>
#include <set>
#include <utility>
#include <variant>
#include <vector>
#include "bapel/stl_vector.h"

template <typename T>
using Set = std::set<T>;

// @bpl: pub SetImpl::mk: forall ['a] () -> Set 'a
// @bpl: pub SetImpl::insert: forall ['a] (&Set 'a, 'a) -> bool
// @bpl: pub SetImpl::remove: forall ['a] (&Set 'a, &'a) -> bool
// @bpl: pub SetImpl::clear: forall ['a] &Set 'a -> ()
// @bpl: pub SetImpl::size: forall ['a] &Set 'a -> i64
// @bpl: pub SetImpl::empty: forall ['a] &Set 'a -> bool
// @bpl: pub SetImpl::contains: forall ['a] (&Set 'a, &'a) -> bool
// @bpl: pub SetImpl::values: forall ['a] &Set 'a -> Vector 'a
struct SetImpl {
  SetImpl() = delete;

  template <typename T>
  static inline Set<T> mk() {
    return Set<T>();
  }

  template <typename T>
  static inline bool insert(Set<T>* s, T val) {
    return s->insert(std::move(val)).second;
  }

  template <typename T>
  static inline bool remove(Set<T>* s, const T* val) {
    return s->erase(*val) > 0;
  }

  template <typename T>
  static inline std::monostate clear(Set<T>* s) {
    s->clear();
    return std::monostate();
  }

  template <typename T>
  static inline int64_t size(const Set<T>* s) {
    return static_cast<int64_t>(s->size());
  }

  template <typename T>
  static inline bool empty(const Set<T>* s) {
    return s->empty();
  }

  template <typename T>
  static inline bool contains(const Set<T>* s, const T* val) {
    return s->find(*val) != s->end();
  }

  template <typename T>
  static inline Vector<T> values(const Set<T>* s) {
    Vector<T> res;
    res.reserve(s->size());
    for (const auto& v : *s) {
      res.push_back(v);
    }
    return res;
  }
};
