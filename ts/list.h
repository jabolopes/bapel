#pragma once

#include <cstddef>
#include <functional>
#include <memory>
#include <optional>
#include <stdexcept>
#include <utility>
#include <vector>

namespace ts {

template <typename T>
class ListIterator;

template <typename T>
class List {
 public:
  struct Node {
    List<T> parent;
    T value;
  };

  List() : node_(nullptr), size_(0) {}
  List(std::shared_ptr<Node> node, size_t size) : node_(std::move(node)), size_(size) {}

  bool empty() const { return node_ == nullptr; }
  size_t size() const { return size_; }

  List<T> add(T value) const {
    auto n = std::make_shared<Node>(Node{*this, std::move(value)});
    return List<T>(std::move(n), size_ + 1);
  }

  List<T> cons(T value) const {
    return add(std::move(value));
  }

  List<T> remove() const {
    if (!node_) return *this;
    return node_->parent;
  }

  List<T> tail() const {
    return remove();
  }

  List<T> pop() const {
    return remove();
  }

  bool value(T& out) const {
    if (!node_) return false;
    out = node_->value;
    return true;
  }

  std::optional<T> head() const {
    if (!node_) return std::nullopt;
    return node_->value;
  }

  const T* head_ptr() const {
    if (!node_) return nullptr;
    return &node_->value;
  }

  const T& front() const {
    if (!node_) throw std::runtime_error("List is empty");
    return node_->value;
  }

  // Collect from oldest to newest: [1, 2, 3]
  std::vector<T> collect() const {
    std::vector<T> result(size_);
    auto curr = *this;
    for (size_t i = size_; i > 0; --i) {
      result[i - 1] = curr.node_->value;
      curr = curr.node_->parent;
    }
    return result;
  }

  std::vector<T> to_vector() const {
    return collect();
  }

  List<T> reverse() const {
    List<T> rev;
    auto curr = *this;
    while (curr.node_) {
      rev = rev.add(curr.node_->value);
      curr = curr.node_->parent;
    }
    return rev;
  }

  template <typename Pred>
  List<T> filter(Pred&& pred) const {
    auto items = collect();
    List<T> result;
    for (const auto& item : items) {
      if (pred(item)) {
        result = result.add(item);
      }
    }
    return result;
  }

  template <typename Func>
  auto map(Func&& f) const -> List<decltype(f(std::declval<T>()))> {
    using U = decltype(f(std::declval<T>()));
    auto items = collect();
    List<U> result;
    for (const auto& item : items) {
      result = result.add(f(item));
    }
    return result;
  }

  template <typename Func>
  void for_each(Func&& f) const {
    const auto* curr = node_.get();
    while (curr) {
      f(curr->value);
      curr = curr->parent.node_.get();
    }
  }

  template <typename Pred>
  const T* find_if(Pred&& pred) const {
    const auto* curr = node_.get();
    while (curr) {
      if (pred(curr->value)) {
        return &curr->value;
      }
      curr = curr->parent.node_.get();
    }
    return nullptr;
  }

  const Node* head_node_ptr() const {
    return node_.get();
  }

  static List<T> from_vector(const std::vector<T>& vec) {
    List<T> l;
    for (const auto& v : vec) {
      l = l.add(v);
    }
    return l;
  }

  ListIterator<T> iterate() const {
    return ListIterator<T>(*this);
  }

  bool operator==(const List<T>& other) const {
    if (size_ != other.size_) return false;
    if (node_ == other.node_) return true;
    auto a = *this;
    auto b = other;
    while (a.node_ && b.node_) {
      if (!(a.node_->value == b.node_->value)) return false;
      a = a.node_->parent;
      b = b.node_->parent;
    }
    return a.node_ == nullptr && b.node_ == nullptr;
  }

  bool operator!=(const List<T>& other) const {
    return !(*this == other);
  }

 private:
  std::shared_ptr<Node> node_;
  size_t size_ = 0;
};

template <typename T>
class ListIterator {
 public:
  explicit ListIterator(const List<T>& list) : curr_(list.head_node_ptr()), size_(list.size()) {}

  bool next(size_t& index, T& value) {
    if (!curr_) return false;
    index = --size_;
    value = curr_->value;
    curr_ = curr_->parent.head_node_ptr();
    return true;
  }

  std::optional<std::pair<size_t, T>> next() {
    if (!curr_) return std::nullopt;
    size_t index = --size_;
    T value = curr_->value;
    curr_ = curr_->parent.head_node_ptr();
    return std::make_pair(index, std::move(value));
  }

  std::vector<T> collect() const {
    std::vector<T> result;
    result.reserve(size_);
    const auto* c = curr_;
    while (c) {
      result.push_back(c->value);
      c = c->parent.head_node_ptr();
    }
    std::reverse(result.begin(), result.end());
    return result;
  }

 private:
  const typename List<T>::Node* curr_ = nullptr;
  size_t size_ = 0;
};

} // namespace ts
