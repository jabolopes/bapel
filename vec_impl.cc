export module vec:vec_impl;

import <cstdint>;
import <vector>;

export namespace vec {

template <typename T>
using Vector = std::vector<T>;

template <typename T>
Vector<T> mk() {
  return Vector<T>();
}

template <typename T>
void add(Vector<T>& v, const T& value) {
  v.push_back(value);
}

template <typename T>
T get(Vector<T>& v, int64_t index) {
  return v[index];
}

template <typename T>
void set(Vector<T>& v, int64_t index, T&& value) {
  v[index] = std::move(value);
}

}  // namespace vec
