export module vec:vec_impl;

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

}  // namespace vec
