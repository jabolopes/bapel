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

}  // namespace c

export namespace ecs {

int64_t addEntity() {
  static int64_t idgen = 0;
  return idgen++;
}

template<typename V>
struct Iterator {
  std::tuple<int64_t, V, bool> Next();
};

template<typename V>
struct Component {
  std::pair<V, bool> Get(int64_t key) const;
  void Set(int64_t key, V value) const;

  template <template <typename> typename I, typename std::enable_if<std::is_default_constructible<Iterator<V>>::value>::type>
  I<V> Iterate(const Iterator<V>& iterator) const;
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
    std::tuple<int64_t, V, bool> Next(IteratorImpl<V>& it) const {
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

    std::pair<V, bool> Get(int64_t key) const {
      return pool_.Get(key);
    }

    void Set(int64_t key, V value) const {
      pool_.Set(key, value);
    }

    IteratorImpl<V> Iterate() const {
      return IteratorImpl<V>(pool_);
    }
  };
};

template <int64_t Size>
template <typename V>
StaticPool<V, Size> StaticComponent<Size>::Component<V>::pool_;

// Component a => i64 -> a
template<typename V>
std::pair<V, bool> get(const Component<V>& component, int64_t entityId) {
  return component.Get(entityId);
}

// Component a => i64 -> a -> ()
template<typename V>
void set(const Component<V>& component, int64_t entityId, V value) {
  component.Set(entityId, std::move(value));
}

// (Component a, Iterator b a) => b a
template <typename V, typename I>
I iterate(const Component<V>& component) {
  static_assert(std::is_default_constructible<Iterator<V>>::value, "Missing");
  return component.Iterate();
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
