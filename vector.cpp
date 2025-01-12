export module vector;

import <vector>;

export namespace vector {

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

}  // namespace
