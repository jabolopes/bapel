#pragma once

#include <cstdint>
#include <memory>
#include <utility>
#include <variant>

template <typename T>
using SharedPtr = std::shared_ptr<T>;

// @bpl: pub SharedPtrImpl::make_shared: forall ['a] 'a -> SharedPtr 'a
// @bpl: pub SharedPtrImpl::null: forall ['a] () -> SharedPtr 'a
// @bpl: pub SharedPtrImpl::is_null: forall ['a] &SharedPtr 'a -> bool
// @bpl: pub SharedPtrImpl::get: forall ['a] &SharedPtr 'a -> 'a
// @bpl: pub SharedPtrImpl::use_count: forall ['a] &SharedPtr 'a -> i64
struct SharedPtrImpl {
  SharedPtrImpl() = delete;

  template <typename T>
  static inline SharedPtr<T> make_shared(T val) {
    return std::make_shared<T>(std::move(val));
  }

  template <typename T>
  static inline SharedPtr<T> null() {
    return nullptr;
  }

  template <typename T>
  static inline bool is_null(const SharedPtr<T>* ptr) {
    return *ptr == nullptr;
  }

  template <typename T>
  static inline T& get(const SharedPtr<T>* ptr) {
    return **ptr;
  }

  template <typename T>
  static inline int64_t use_count(const SharedPtr<T>* ptr) {
    return ptr->use_count();
  }
};
