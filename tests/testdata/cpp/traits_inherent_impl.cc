
#include "traits_inherent_impl_private.h"

int64_t test_vector(::Ptr<::Vector<int64_t>> v) {
  int64_t a = ::inherents::Vector<int64_t>::get(v, static_cast<int64_t>(0));
  int64_t b = ::inherents::Vector<int64_t>::get(v, static_cast<int64_t>(1));
  return (a) + (b);
}

int64_t test_point() {
  int64_t a = ::inherents::Point::sum(
      ::inherents::Ptr<::Point>::mk(::inherents::Point::mk(
          static_cast<int64_t>(1), static_cast<int64_t>(2))));
  ::Point p1 =
      ::inherents::Point::mk(static_cast<int64_t>(3), static_cast<int64_t>(4));
  int64_t b = ::inherents::Point::sum(::inherents::Ptr<::Point>::mk(p1));
  ::Ptr<::Point> ref_p1 = ::inherents::Ptr<::Point>::mk(p1);
  int64_t c = ::inherents::Point::sum(ref_p1);
  int64_t d = ::inherents::Point::sum(ref_p1);
  ::Point p2 =
      ::inherents::Point::mk(static_cast<int64_t>(5), static_cast<int64_t>(6));
  ::Point p3 = ::inherents::Point::add(::inherents::Ptr<::Point>::mk(p1),
                                       ::inherents::Ptr<::Point>::mk(p2));
  int64_t e = ::inherents::Point::sum(::inherents::Ptr<::Point>::mk(p3));
  return ((((a) + (b)) + (c)) + (d)) + (e);
}

int64_t test_node(::Ptr<::Node> n) {
  int64_t a = ::inherents::Node::get_val(n);
  int64_t b = ::inherents::Node::get_val(n);
  int64_t c = ::inherents::Node::next_val(n);
  int64_t d = ::inherents::Node::next_val(n);
  int64_t e = ::inherents::Node::get_val(n);
  return ((((a) + (b)) + (c)) + (d)) + (e);
}

bool test_color() {
  ::Color c1 = ::Color(std::in_place_index<0>, std::monostate());
  bool a = ::inherents::Color::is_red(::inherents::Ptr<::Color>::mk(c1));
  bool b = ::inherents::Color::is_red(::inherents::Ptr<::Color>::mk(c1));
  ::Ptr<::Color> ref_c1 = ::inherents::Ptr<::Color>::mk(c1);
  bool c = ::inherents::Color::is_red(ref_c1);
  bool d = ::inherents::Color::is_red(ref_c1);
  bool e = ::inherents::Color::is_red(::inherents::Ptr<::Color>::mk(c1));
  return ((((a) && (b)) && (c)) && (d)) && (e);
}

int64_t test_int_list(::Ptr<::IntList> l) {
  ::IntList empty = ::IntList(std::in_place_index<0>, std::monostate());
  bool a = ::inherents::IntList::is_nil(::inherents::Ptr<::IntList>::mk(empty));
  bool b = ::inherents::IntList::is_nil(::inherents::Ptr<::IntList>::mk(empty));
  bool c = ::inherents::IntList::is_nil(l);
  bool d = ::inherents::IntList::is_nil(l);
  int64_t s1 =
      ::inherents::IntList::sum(::inherents::Ptr<::IntList>::mk(empty));
  int64_t s2 =
      ::inherents::IntList::sum(::inherents::Ptr<::IntList>::mk(empty));
  int64_t s3 = ::inherents::IntList::sum(l);
  int64_t s4 = ::inherents::IntList::sum(l);
  int64_t s5 = ::inherents::IntList::sum(l);
  if ((((a) && (b)) && (c)) && (d)) {
    return (s1) + (s2);
  } else {
    return ((s3) + (s4)) + (s5);
  };
}
