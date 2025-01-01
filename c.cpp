export module c;

import <array>;
import <cassert>;
import <cerrno>;
import <cstdint>;
import <ctime>;
import <iostream>;
import <tuple>;
import <vector>;

// Needed because of import<vector> results in Bad file data:
// https://stackoverflow.com/questions/70456868/vector-in-c-module-causes-useless-bad-file-data-gcc-output
namespace std _GLIBCXX_VISIBILITY(default){}

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

template <typename T>
std::array<T, 10> mkArray() {
  std::array<T, 10> a;
  a.fill(0);
  return a;
}

}  // namespace c

export namespace ecs {

int64_t addEntity() {
  static int64_t idgen = 0;
  return idgen++;
}

template<typename V>
struct Iterator {
  template <template <typename> typename I>
  static std::tuple<int64_t, V, bool> Next(I<V>& iterator);
};

template<typename V>
struct Component {
  static std::pair<V, bool> Get(int64_t key);
  static void Set(int64_t key, V value);

  template <template <typename> typename I>
  static I<V> Iterate();
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

template <int64_t Size>
struct StaticComponent {
  template <typename V>
  struct IteratorImpl {
    IteratorImpl() : pool_(nullptr), key_(0) {}

    explicit IteratorImpl(StaticPool<V, Size>& pool) : pool_(&pool), key_(0) {}

    StaticPool<V, Size> *pool_;
    int64_t key_;
  };

  template <typename V>
  struct Iterator {
    static std::tuple<int64_t, V, bool> Next(IteratorImpl<V>& it) {
      assert(it.pool_ != nullptr);

      while (true) {
        if (it.key_ >= Size) {
          return std::make_tuple(0, V{}, false);
        }

        auto [value, ok] = it.pool_->Get(it.key_);
        if (!ok) {
          ++it.key_;
          continue;
        }

        return std::make_tuple(it.key_++, value, true);
      }
    }
  };

  template <typename V>
  struct Component {
    static StaticPool<V, Size> pool_;

    static std::pair<V, bool> Get(int64_t key) {
      return pool_.Get(key);
    }

    static void Set(int64_t key, V value) {
      pool_.Set(key, value);
    }

    static IteratorImpl<V> Iterate() {
      return IteratorImpl<V>(pool_);
    }
  };
};

template <int64_t Size>
template <typename V>
StaticPool<V, Size> StaticComponent<Size>::Component<V>::pool_;

// Component a => i64 -> a
template<typename V>
std::pair<V, bool> get(int64_t entityId) {
  return Component<V>::Get(entityId);
}

// Component a => i64 -> a -> ()
template<typename V>
void set(int64_t entityId, V value) {
  Component<V>::Set(entityId, std::move(value));
}

// (Component a, Iterator b a) => b a
template <typename V, typename I>
I iterate() {
  return Component<V>::Iterate();
}

// (Component a, Iterator b a) => b a
template <typename V, typename I>
std::tuple<int64_t, V, bool> next(I& iterator) {
  return Iterator<V>::Next(iterator);
}

template <typename V>
std::vector<V> collect(const Component<V>& component) {
  std::vector<V> values;
  auto it = component.Iterate();
  while (true) {
    auto [id, value, ok] = (Iterator<V>{}).Next(it);
    if (!ok) {
      break;
    }
    values.push_back(std::move(value));
  }
  return values;
}

}  // namespace entity
