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

template <typename V>
class Iterator {
 public:
  virtual std::tuple<int64_t, V, bool> Next() = 0;
};

template<typename V>
class Component {
 public:
  virtual std::pair<V, bool> Get(int64_t key) = 0;
  virtual void Set(int64_t key, V value) = 0;
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

  void Iterate() {
  }

 private:
  std::array<bool, Size> has_;
  std::array<V, Size> dense_;
};

template <typename V, int64_t Size>
class StaticIterator final : public Iterator<V> {
 private:
  using Pool = StaticPool<V, Size>;

 public:
  explicit StaticIterator(Pool& pool) : pool_(pool), key_(0) {}

  std::tuple<int64_t, V, bool> Next() override {
    while (true) {
      if (key_ >= Size) {
        return std::make_tuple(0, V{}, false);
      }

      auto [value, ok] = pool_.Get(key_);
      if (!ok) {
        ++key_;
        continue;
      }

      return std::make_tuple(key_++, value, true);
    }
  }

 private:
  Pool &pool_;
  int64_t key_;
};

template <typename V, int size>
class StaticComponent final : public Component<V> {
 public:
  std::pair<V, bool> Get(int64_t key) override { return pool_.Get(key); }
  void Set(int64_t key, V value) override { return pool_.Set(key, value); }

 private:
  StaticPool<V, size> pool_;
};

// Component a => i64 -> a
template<typename V>
std::pair<V, bool> get(Component<V>* component, int64_t entityId) {
  return component->Get(entityId);
}

// Component a => i64 -> a -> ()
template<typename V>
void set(Component<V>* component, int64_t entityId, V value) {
  component->Set(entityId, value);
}

}  // namespace entity
