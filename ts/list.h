#pragma once

#include <cstddef>
#include <memory>
#include <stdexcept>
#include <utility>
#include <vector>

namespace ts {

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

  List<T> remove() const {
    if (!node_) return *this;
    return node_->parent;
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

  List<T> reverse() const {
    List<T> rev;
    auto curr = *this;
    while (curr.node_) {
      rev = rev.add(curr.node_->value);
      curr = curr.node_->parent;
    }
    return rev;
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

 private:
  std::shared_ptr<Node> node_;
  size_t size_ = 0;
};

} // namespace ts
