export module c;

import <array>;
import <cerrno>;
import <cstdint>;
import <ctime>;
import <iostream>;
import <tuple>;

export namespace c {

std::tuple<int64_t, int64_t> time() {
  auto res = ::time(nullptr);
  if (res == -1) {
    return {0, errno};
  }
  return {res, 0};
}

template <typename T>
void print(T value) {
  std::cout << value << std::endl;
}

}  // namespace c

export namespace ecs {

template <typename V>
class Pool {
 public:
  virtual V Get(int key) = 0;
  virtual void Set(int key, V value) = 0;
};

template <typename V, int64_t Capacity>
class StaticPool final : public Pool<V>{
 public:
  V Get(int key) override { return dense_.at(key); }

  void Set(int key, V value) override { dense_[key] = value; }

 private:
  std::array<V, Capacity> dense_;
};

template<typename V>
struct Component {
  virtual V Get(int key) { return GetPool()->Get(key); }
  virtual void Set(int key, V value) { return GetPool()->Set(key, value); }
  virtual Pool<V>* GetPool() { return nullptr; }
};

template <typename V, int size>
struct StaticComponent final : public Component<V> {
  Pool<V>* GetPool() { return &pool_; }

  StaticPool<V, size> pool_;
};

// Component a => int -> a
template<typename V>
V get(Component<V>* component, int entityId) {
  return component->Get(entityId);
}

// Component a => int -> a -> ()
template<typename V>
void set(Component<V>* component, int entityId, V value) {
  component->Set(entityId, value);
}

struct Hello {};

StaticComponent<Hello, 1024> Component_Hello{};

void f() {
  Hello value = get(&Component_Hello, 0);
  set(&Component_Hello, 0, value);
}

}  // namespace entity
