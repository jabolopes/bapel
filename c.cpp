export module c;

import <array>;
import <cerrno>;
import <cstdint>;
import <ctime>;
import <tuple>;

export namespace c {

std::tuple<int64_t, int64_t> time() {
  auto res = ::time(nullptr);
  if (res == -1) {
    return {0, errno};
  }
  return {res, 0};
}

}  // namespace c

export namespace ecs {

template <typename K, typename V, int64_t Capacity>
class StaticPool final {
 public:
  StaticPool() {}

  V Get(K key) { return dense_.at(key); }

  void Set(K key, V value) { dense_[key] = value; }

 private:
  std::array<V, Capacity> dense_;
};

}  // namespace entity
