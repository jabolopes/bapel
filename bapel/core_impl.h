#pragma once

#include <cassert>
#include <cerrno>
#include <cstdint>
#include <ctime>
#include <iostream>
#include <tuple>
#include <variant>

namespace core {

// @bpl: pub type core::Point = struct{x i32, y i32}
struct Point {
  int x;
  int y;
};

// @bpl: pub core::noopPoint: core::Point -> ()
inline std::monostate noopPoint(Point p) {
  return std::monostate();
}

// @bpl: pub type core::AbsPoint
struct AbsPoint {
  int x;
  int y;
};

// @bpl: pub core::mkAbsPoint: () -> core::AbsPoint
inline AbsPoint mkAbsPoint() { return AbsPoint{}; }

// @bpl: pub core::absPointX: core::AbsPoint -> i32
inline int absPointX(AbsPoint p) { return p.x; }

// @bpl: pub core::noopAbsPoint: core::AbsPoint -> ()
inline std::monostate noopAbsPoint(AbsPoint p) {
  return std::monostate();
}

// @bpl: pub core::time: () -> (i64, i64)
inline std::tuple<int64_t, int64_t> time() {
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

// @bpl: pub core::print: forall ['a] 'a -> ()
template <typename T>
std::monostate print(const T& value) {
  std::cout << value << std::endl;
  return std::monostate();
}

}  // namespace core
