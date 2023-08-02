export module c;

import <cerrno>;
import <cstdint>;
import <ctime>;
import <tuple>;

export namespace c {

std::tuple<int64_t, int64_t> time() {
  auto res = ::time(nullptr);
  if (res == -1) {
    return {0, errno};
  }
  return {res, 0};
}

}
