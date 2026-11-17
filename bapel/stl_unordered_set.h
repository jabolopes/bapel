#pragma once
#include <cstdint>
#include <unordered_set>
#include <utility>
#include <variant>
#include <vector>
#include "bapel/stl_vector.h"

template <typename T>
using UnorderedSet = std::unordered_set<T>;

// @bpl: pub UnorderedSetImpl::mk: forall ['a] () -> UnorderedSet 'a
// @bpl: pub UnorderedSetImpl::insert: forall ['a] (&UnorderedSet 'a, 'a) -> bool
// @bpl: pub UnorderedSetImpl::remove: forall ['a] (&UnorderedSet 'a, &'a) -> bool
// @bpl: pub UnorderedSetImpl::clear: forall ['a] &UnorderedSet 'a -> ()
// @bpl: pub UnorderedSetImpl::size: forall ['a] &UnorderedSet 'a -> i64
// @bpl: pub UnorderedSetImpl::empty: forall ['a] &UnorderedSet 'a -> bool
// @bpl: pub UnorderedSetImpl::contains: forall ['a] (&UnorderedSet 'a, &'a) -> bool
// @bpl: pub UnorderedSetImpl::values: forall ['a] &UnorderedSet 'a -> Vector 'a
struct UnorderedSetImpl {
  UnorderedSetImpl() = delete;

  template <typename T>
  static inline UnorderedSet<T> mk() {
    return UnorderedSet<T>();
  }

  template <typename T>
  static inline bool insert(UnorderedSet<T>* s, T val) {
    return s->insert(std::move(val)).second;
  }

  template <typename T>
  static inline bool remove(UnorderedSet<T>* s, const T* val) {
    return s->erase(*val) > 0;
  }

  template <typename T>
  static inline std::monostate clear(UnorderedSet<T>* s) {
    s->clear();
    return std::monostate();
  }

  template <typename T>
  static inline int64_t size(const UnorderedSet<T>* s) {
    return static_cast<int64_t>(s->size());
  }

  template <typename T>
  static inline bool empty(const UnorderedSet<T>* s) {
    return s->empty();
  }

  template <typename T>
  static inline bool contains(const UnorderedSet<T>* s, const T* val) {
    return s->find(*val) != s->end();
  }

  template <typename T>
  static inline Vector<T> values(const UnorderedSet<T>* s) {
    Vector<T> res;
    res.reserve(s->size());
    for (const auto& v : *s) {
      res.push_back(v);
    }
    return res;
  }
};
