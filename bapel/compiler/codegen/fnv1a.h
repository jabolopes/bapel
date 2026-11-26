#pragma once

#include <cstdint>
#include <string>
#include <string_view>
#include "bapel/stl_string.h"

// @bpl: pub fnv1a_hash_str: StringView -> String
inline String fnv1a_hash_str(StringView s) {
  uint64_t h = 14695981039346656037ULL;
  for (unsigned char c : s) {
    h ^= c;
    h *= 1099511628211ULL;
  }
  return std::to_string(h);
}
