export module c;

import <array>;
import <cerrno>;
import <cstdint>;
import <ctime>;
import <iostream>;
import <tuple>;

export namespace c {

struct Point {
  int x;
  int y;
};

void noopPoint(Point p) {}

struct AbsPoint {
  int x;
  int y;
};

AbsPoint mkAbsPoint() { return AbsPoint{}; }
int absPointX(AbsPoint p) { return p.x; }
void noopAbsPoint(AbsPoint p) {}

std::tuple<int64_t, int64_t> time() {
  auto res = ::time(nullptr);
  if (res == -1) {
    return {0, errno};
  }
  return {res, 0};
}

template <typename T1, typename T2>
std::ostream& operator<<(std::ostream& os, std::tuple<T1, T2> const& v) {
  return os << "("
            << std::get<0>(v)
            << ", "
            << std::get<1>(v)
            << ")"
            << std::endl;
}

template <typename T>
void print(T value) {
  std::cout << value << std::endl;
}

}  // namespace c

export namespace ecs {

int64_t addEntity() {
  static int64_t idgen = 0;
  return idgen++;
}

template<typename T>
struct Iterator {
  using Value = T::Value;
  static std::tuple<int64_t, Value, bool> Next(T& iterator);
};

template<typename A>
struct Component {
  using Value = A::Value;
  using Iterator = A::Iterator;

  std::pair<Value, bool> Get(int64_t key) const;
  void Set(int64_t key, Value value) const;
  Iterator Iterate() const;
};

template <typename V, int64_t Size>
class StaticPool final {
 public:
  StaticPool() {
    for (size_t i = 0; i < has_.size(); ++i) {
      has_[i] = false;
    }
  }

  std::pair<V, bool> Get(int64_t key) {
    if (!has_.at(key)) {
      return std::make_pair(V{}, false);
    }
    return std::make_pair(dense_.at(key), true);
  }

  void Set(int64_t key, V value) {
    has_.at(key) = true;
    dense_[key] = value;
  }

  void Unset(int64_t key) {
    has_.at(key) = false;
  }

 private:
  std::array<bool, Size> has_;
  std::array<V, Size> dense_;
};

template <typename V, int64_t Size>
class StaticIterator final {
 public:
  explicit StaticIterator(StaticPool<V, Size>& pool) : pool_(pool), key_(0) {}

  StaticPool<V, Size> &pool_;
  int64_t key_;
};

template <typename V, int64_t Size>
struct Iterator<StaticIterator<V, Size>> {
  using Value = V;

  static std::tuple<int64_t, V, bool> Next(StaticIterator<V, Size>& iterator) {
    while (true) {
      if (iterator.key_ >= Size) {
        return std::make_tuple(0, V{}, false);
      }

      auto [value, ok] = iterator.pool_.Get(iterator.key_);
      if (!ok) {
        ++iterator.key_;
        continue;
      }

      return std::make_tuple(iterator.key_++, value, true);
    }
  }
};

template <typename V, int64_t Size>
struct StaticComponent {};

template <typename V, int64_t Size>
struct Component<StaticComponent<V, Size>> {
 public:
  using Value = V;
  using Iterator = StaticIterator<V, Size>;

  static StaticPool<V, Size> pool_;

  std::pair<V, bool> Get(int64_t key) const {
    return pool_.Get(key);
  }

  void Set(int64_t key, V value) const {
    pool_.Set(key, value);
  }

  Iterator Iterate() const {
    return Iterator(pool_);
  }
};

template <typename V, int64_t Size>
StaticPool<V, Size> Component<StaticComponent<V, Size>>::pool_;

// Component Value a => a -> i64 -> Value
template<typename C>
std::pair<typename Component<C>::Value, bool> get(const Component<C>& component, int64_t entityId) {
  return component.Get(entityId);
}

// Component Value a => a -> i64 -> Value -> ()
template<typename C>
void set(const Component<C>& component, int64_t entityId, typename Component<C>::Value value) {
  component.Set(entityId, value);
}

template <typename C>
Component<C>::Iterator iterate(const Component<C>& component) {
  return component.Iterate();
}

void f() {
  get<StaticComponent<int, 10>>(Component<StaticComponent<int, 10>>{}, 1);
  set<StaticComponent<int, 10>>(Component<StaticComponent<int, 10>>{}, 1, 1);

  {
    auto it = iterate<StaticComponent<int, 10>>({});
    std::tuple<int64_t, int, bool> v = Iterator<StaticIterator<int, 10>>::Next(it);
  }
}

}  // namespace entity
