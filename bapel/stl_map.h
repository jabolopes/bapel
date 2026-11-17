#pragma once
#include <cstdint>
#include <map>
#include <tuple>
#include <utility>
#include <variant>
#include <vector>
#include "bapel/stl_optional.h"
#include "bapel/stl_vector.h"

template <typename K, typename V>
using Map = std::map<K, V>;

// @bpl: pub MapImpl::mk: forall ['k, 'v] () -> Map 'k 'v
// @bpl: pub MapImpl::insert: forall ['k, 'v] (&Map 'k 'v, 'k, 'v) -> ()
// @bpl: pub MapImpl::remove: forall ['k, 'v] (&Map 'k 'v, &'k) -> bool
// @bpl: pub MapImpl::clear: forall ['k, 'v] &Map 'k 'v -> ()
// @bpl: pub MapImpl::size: forall ['k, 'v] &Map 'k 'v -> i64
// @bpl: pub MapImpl::empty: forall ['k, 'v] &Map 'k 'v -> bool
// @bpl: pub MapImpl::contains: forall ['k, 'v] (&Map 'k 'v, &'k) -> bool
// @bpl: pub MapImpl::get: forall ['k, 'v] (&Map 'k 'v, &'k) -> Optional 'v
// @bpl: pub MapImpl::keys: forall ['k, 'v] &Map 'k 'v -> Vector 'k
// @bpl: pub MapImpl::entries: forall ['k, 'v] &Map 'k 'v -> Vector ('k, 'v)
struct MapImpl {
  MapImpl() = delete;

  template <typename K, typename V>
  static inline Map<K, V> mk() {
    return Map<K, V>();
  }

  template <typename K, typename V>
  static inline std::monostate insert(Map<K, V>* m, K key, V val) {
    (*m)[std::move(key)] = std::move(val);
    return std::monostate();
  }

  template <typename K, typename V>
  static inline bool remove(Map<K, V>* m, const K* key) {
    return m->erase(*key) > 0;
  }

  template <typename K, typename V>
  static inline std::monostate clear(Map<K, V>* m) {
    m->clear();
    return std::monostate();
  }

  template <typename K, typename V>
  static inline int64_t size(const Map<K, V>* m) {
    return static_cast<int64_t>(m->size());
  }

  template <typename K, typename V>
  static inline bool empty(const Map<K, V>* m) {
    return m->empty();
  }

  template <typename K, typename V>
  static inline bool contains(const Map<K, V>* m, const K* key) {
    return m->find(*key) != m->end();
  }

  template <typename K, typename V>
  static inline Optional<V> get(const Map<K, V>* m, const K* key) {
    auto it = m->find(*key);
    if (it == m->end()) {
      return std::nullopt;
    }
    return std::make_optional(it->second);
  }

  template <typename K, typename V>
  static inline Vector<K> keys(const Map<K, V>* m) {
    Vector<K> res;
    res.reserve(m->size());
    for (const auto& kv : *m) {
      res.push_back(kv.first);
    }
    return res;
  }

  template <typename K, typename V>
  static inline Vector<std::tuple<K, V>> entries(const Map<K, V>* m) {
    Vector<std::tuple<K, V>> res;
    res.reserve(m->size());
    for (const auto& kv : *m) {
      res.push_back(std::make_tuple(kv.first, kv.second));
    }
    return res;
  }
};
