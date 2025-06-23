module;

#include <cstdint>
#include <vector>

export module vec:vec_impl;

export namespace vec {

// @bpl: export type vec.Vector ['a]
template <typename T>
using Vector = std::vector<T>;

// @bpl: export vec.mk: () -> vec.Vector 'a
template <typename T>
Vector<T> mk() {
  return Vector<T>();
}

// @bpl: export vec.add: forall ['a] (vec.Vector 'a, 'a) -> ()
template <typename T>
void add(Vector<T>& v, const T& value) {
  v.push_back(value);
}

// @bpl: export vec.get: forall ['a] (vec.Vector 'a, i64) -> 'a
template <typename T>
T get(Vector<T>& v, int64_t index) {
  return v[index];
}

// @bpl: export vec.set: forall ['a] (vec.Vector 'a, i64, 'a) -> ()
template <typename T>
void set(Vector<T>& v, int64_t index, T&& value) {
  v[index] = std::move(value);
}

}  // namespace vec
