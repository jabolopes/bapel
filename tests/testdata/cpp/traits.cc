
#include "traits_private.h"

::Point make_point(int64_t x, int64_t y) { return {.x = x, .y = y}; }

int64_t test_vector(::Ptr<::Vector<int64_t>> v) {
  int64_t a = ::traits::Indexable<::Vector<int64_t>, int64_t>::get(
      v, static_cast<int64_t>(0));
  int64_t b = ::traits::Indexable<::Vector<int64_t>, int64_t>::get(
      v, static_cast<int64_t>(1));
  return (a) + (b);
}

int64_t test_point() {
  int64_t tmp = ::traits::Size<::Point>::size(::inherents::Ptr<::Point>::mk(
      make_point(static_cast<int64_t>(1), static_cast<int64_t>(2))));
  ::Point p = make_point(static_cast<int64_t>(3), static_cast<int64_t>(4));
  ::Point p2 = make_point(static_cast<int64_t>(5), static_cast<int64_t>(6));
  int64_t a = ::traits::Size<::Point>::size(::inherents::Ptr<::Point>::mk(p));
  int64_t b = ::traits::Size<::Point>::size(::inherents::Ptr<::Point>::mk(p));
  ::Ptr<::Point> ref_p = ::inherents::Ptr<::Point>::mk(p);
  int64_t c = ::traits::Size<::Point>::size(ref_p);
  int64_t d = ::traits::Size<::Point>::size(ref_p);
  int64_t e = ::traits::Size<::Point>::size(::inherents::Ptr<::Point>::mk(p));
  int64_t f = ::traits::Add<::Point>::add(::inherents::Ptr<::Point>::mk(p),
                                          ::inherents::Ptr<::Point>::mk(p2));
  return ((((((tmp) + (a)) + (b)) + (c)) + (d)) + (e)) + (f);
}

int64_t test_node(::Ptr<::Node> n) {
  int64_t a = ::traits::Size<::Node>::size(n);
  int64_t b = ::traits::Size<::Node>::size(n);
  int64_t c = ::traits::Size<::Node>::size(n);
  return ((a) + (b)) + (c);
}

int64_t test_color() {
  ::Color col = ::Color(std::in_place_index<0>, std::monostate());
  int64_t a = ::traits::Size<::Color>::size(::inherents::Ptr<::Color>::mk(col));
  int64_t b = ::traits::Size<::Color>::size(::inherents::Ptr<::Color>::mk(col));
  ::Ptr<::Color> ref_col = ::inherents::Ptr<::Color>::mk(col);
  int64_t c = ::traits::Size<::Color>::size(ref_col);
  int64_t d = ::traits::Size<::Color>::size(ref_col);
  int64_t e = ::traits::Size<::Color>::size(::inherents::Ptr<::Color>::mk(col));
  return ((((a) + (b)) + (c)) + (d)) + (e);
}

int64_t test_int_list(::Ptr<::IntList> l) {
  ::IntList empty = ::IntList(std::in_place_index<0>, std::monostate());
  int64_t a =
      ::traits::Size<::IntList>::size(::inherents::Ptr<::IntList>::mk(empty));
  int64_t b =
      ::traits::Size<::IntList>::size(::inherents::Ptr<::IntList>::mk(empty));
  int64_t c = ::traits::Size<::IntList>::size(l);
  int64_t d = ::traits::Size<::IntList>::size(l);
  int64_t e = ::traits::Size<::IntList>::size(l);
  return ((((a) + (b)) + (c)) + (d)) + (e);
}
