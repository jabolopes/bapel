#pragma once
#include <unordered_map>
#include <cstdint>
#include <utility>
#include <variant>
#include "bapel/stl_optional.h"

template <typename K, typename V>
using UnorderedMap = std::unordered_map<K, V>;

// @bpl: pub UnorderedMapImpl::mk: forall ['k, 'v] () -> UnorderedMap 'k 'v
// @bpl: pub UnorderedMapImpl::insert: forall ['k, 'v] (&UnorderedMap 'k 'v, 'k, 'v) -> ()
// @bpl: pub UnorderedMapImpl::size: forall ['k, 'v] &UnorderedMap 'k 'v -> i64
// @bpl: pub UnorderedMapImpl::empty: forall ['k, 'v] &UnorderedMap 'k 'v -> bool
// @bpl: pub UnorderedMapImpl::contains: forall ['k, 'v] (&UnorderedMap 'k 'v, &'k) -> bool
// @bpl: pub UnorderedMapImpl::get: forall ['k, 'v] (&UnorderedMap 'k 'v, &'k) -> Optional 'v
struct UnorderedMapImpl {
  UnorderedMapImpl() = delete;

  template <typename K, typename V>
  static inline UnorderedMap<K, V> mk() {
    return UnorderedMap<K, V>();
  }

  template <typename K, typename V>
  static inline std::monostate insert(UnorderedMap<K, V>* m, K key, V val) {
    (*m)[std::move(key)] = std::move(val);
    return std::monostate();
  }

  template <typename K, typename V>
  static inline int64_t size(UnorderedMap<K, V>* m) {
    return m->size();
  }

  template <typename K, typename V>
  static inline bool empty(UnorderedMap<K, V>* m) {
    return m->empty();
  }

  template <typename K, typename V>
  static inline bool contains(UnorderedMap<K, V>* m, const K* key) {
    return m->find(*key) != m->end();
  }

  template <typename K, typename V>
  static inline Optional<V> get(UnorderedMap<K, V>* m, const K* key) {
    auto it = m->find(*key);
    if (it == m->end()) {
      return std::nullopt;
    }
    return std::make_optional(it->second);
  }
};

