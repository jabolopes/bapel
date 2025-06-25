module;

#include <cassert>
#include <cerrno>
#include <cstdint>
#include <ctime>
#include <iostream>
#include <tuple>
#include <variant>
#include <vector>

export module core:core_impl;

// Needed because of import<vector> results in Bad file data:
// https://stackoverflow.com/questions/70456868/vector-in-c-module-causes-useless-bad-file-data-gcc-output
namespace std _GLIBCXX_VISIBILITY(default){}

export namespace core {

// @bpl: export type core.Point = struct{x i32, y i32}
struct Point {
  int x;
  int y;
};

// @bpl: export core.noopPoint: core.Point -> ()
std::monostate noopPoint(Point p) {
  return std::monostate();
}

// @bpl: export type core.AbsPoint
struct AbsPoint {
  int x;
  int y;
};

// @bpl: export core.mkAbsPoint: () -> core.AbsPoint
AbsPoint mkAbsPoint() { return AbsPoint{}; }

// @bpl: export core.absPointX: core.AbsPoint -> i32
int absPointX(AbsPoint p) { return p.x; }

// @bpl: export core.noopAbsPoint: core.AbsPoint -> ()
std::monostate noopAbsPoint(AbsPoint p) {
  return std::monostate();
}

// @bpl: export core.time: () -> (i64, i64)
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

// @bpl: export core.print: forall ['a] 'a -> ()
template <typename T>
std::monostate print(T value) {
  std::cout << value << std::endl;
  return std::monostate();
}

// @bpl: export type std.optional ['a]
//
// No definition because it re-exports from C++.

// @bpl: export type std.string
//
// No definition because it re-exports from C++.

}  // namespace core
