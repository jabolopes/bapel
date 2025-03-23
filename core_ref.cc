module;

#include <functional>

export module core:core_ref;

export namespace ref {

template <typename T>
class Ref final {
 private:
  struct Impl {
    explicit Impl(T&& data) : count(1), data(std::move(data)) {}

    int count;
    T data;
  };

 public:
  explicit Ref(T&& data) : impl_(new Impl(std::move(data))) {}

  Ref(Ref& ref) {
    if (impl_ != ref.impl_) {
      acquire(ref);
    }
  }

  ~Ref() {
    release();
  }

  Ref& operator=(const Ref& ref) {
    if (impl_ != ref.impl_) {
      release();
      acquire(ref);
    }
  }

  T* operator->() {
    return &impl_->data;
  }

  const T* operator->() const {
    return &impl_->data;
  }

  T& operator*() {
    return impl_->data;
  }

  const T& operator*() const {
    return impl_->data;
  }

 private:
  void acquire(const Ref& ref) {
    impl_ = ref.impl_;
    impl_->count++;
  }

  void release() {
    impl_->count--;
    if (impl_->count <= 0) {
      delete impl_;
    }

    impl_ = nullptr;
  }

  Impl *impl_;
};

template <typename T>
Ref<T> mk(T&& data) {
  return Ref<T>(std::move(data));
}

template <typename T>
T& get(Ref<T>& ref) {
  return *ref;
}

template <typename T>
void set(Ref<T>& ref, T&& data) {
  ref->data = std::move(data);
}

}  // namespace ref
