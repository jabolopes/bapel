
#include "traits_bounds_ok_private.h"

int64_t test_vector(::Ptr<::Vector<int64_t>> v) {
  return (print_size<::Vector<int64_t>>(v)) +
         (print_size_dot<::Vector<int64_t>>(v));
}

int64_t test_point() {
  ::Point p = {.x = static_cast<int64_t>(10), .y = static_cast<int64_t>(32)};
  return (print_size<::Point>(::inherents::Ptr<::Point>::mk(p))) +
         (print_size_dot<::Point>(::inherents::Ptr<::Point>::mk(p)));
}

int64_t test_node(::Ptr<::Node> n) { return print_size<::Node>(n); }

int64_t test_color() {
  ::Color c = ::Color(std::in_place_index<0>, std::monostate());
  return print_size<::Color>(::inherents::Ptr<::Color>::mk(c));
}

int64_t test_int_list(::Ptr<::IntList> l) { return print_size<::IntList>(l); }
