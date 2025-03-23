module;

#include <array>
#include <cstdint>

export module arr:arr_impl;

export namespace arr {

template <typename T>
std::array<T, 10> mk() {
  std::array<T, 10> a;
  a.fill(0);
  return a;
}

template <typename T>
T get(const std::array<T, 10>& a, int64_t index) {
  return a[index];
}

template <typename T>
void set(std::array<T, 10>& a, int64_t index, T&& value) {
  a[index] = std::move(value);
}

}  // namespace arr
