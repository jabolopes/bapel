#pragma once

#include <cmath>

namespace core {

// @bpl: pub core::abs: forall ['a] 'a -> 'a
template <typename T>
T abs(T value) {
  return std::abs(value);
}

// @bpl: pub core::squareRoot: f32 -> f32
inline float squareRoot(float value) {
  return sqrtf(value);
}

}  // namespace core
