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

  static std::pair<Value, bool> Get(A &a, int64_t key);
  static void Set(A &a, int64_t key, Value value);
  static Iterator Iterate(A &a);
};

// template <typename T>
// concept Componentable = requires Component<T>;

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
struct StaticComponent {
  StaticPool<V, Size> pool_;
};

template <typename V, int64_t Size>
struct Component<StaticComponent<V, Size>> {
 public:
  using Value = V;
  using Iterator = StaticIterator<V, Size>;

  static std::pair<V, bool> Get(StaticComponent<V, Size> &component, int64_t key) {
    return component.pool_.Get(key);
  }

  static void Set(StaticComponent<V, Size> &component, int64_t key, V value) {
    component.pool_.Set(key, value);
  }

  static Iterator Iterate(StaticComponent<V, Size> &component) {
    return Iterator(component.pool_);
  }
};

// Component a => i64 -> a::Value
template<typename C>
std::pair<typename Component<C>::Value, bool> get(C& component, int64_t entityId) {
  return Component<C>::Get(component, entityId);
}

// Component a => i64 -> a -> ()
template<typename C>
void set(C& component, int64_t entityId, typename Component<C>::Value value) {
  Component<C>::Set(component, entityId, value);
}

template <typename C>
Component<C>::Iterator iterate(C& component) {
  return Component<C>::Iterate(component);
}

void f() {
  StaticComponent<int, 10> component;
  Component<StaticComponent<int, 10>>::Get(component, 1);
  StaticIterator<int, 10> it = Component<StaticComponent<int, 10>>::Iterate(component);

  get<StaticComponent<int, 10>>(component, 1);
  set<StaticComponent<int, 10>>(component, 1, 1);

  {
    auto it = iterate<StaticComponent<int, 10>>(component);
    std::tuple<int64_t, int, bool> v = Iterator<StaticIterator<int, 10>>::Next(it);
  }
}

}  // namespace entity
