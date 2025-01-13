export module core:core_impl;

import <array>;
import <cassert>;
import <cerrno>;
import <cstdint>;
import <ctime>;
import <iostream>;
import <tuple>;
import <vector>;

// Needed because of import<vector> results in Bad file data:
// https://stackoverflow.com/questions/70456868/vector-in-c-module-causes-useless-bad-file-data-gcc-output
namespace std _GLIBCXX_VISIBILITY(default){}

export namespace core {

struct Point {
  int x;
  int y;
};

void noopPoint(Point p) {}

struct AbsPoint {
  int x;
  int y;
};

AbsPoint mkAbsPoint() { return AbsPoint{}; }
int absPointX(AbsPoint p) { return p.x; }
void noopAbsPoint(AbsPoint p) {}

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

template <typename T>
void print(T value) {
  std::cout << value << std::endl;
}

template <typename T>
std::array<T, 10> mkArray() {
  std::array<T, 10> a;
  a.fill(0);
  return a;
}

}  // namespace core
