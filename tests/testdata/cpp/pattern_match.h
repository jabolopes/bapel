#pragma once

#include <array>
#include <cmath>
#include <cstdlib>
#include <functional>
#include <optional>
#include <string>
#include <tuple>
#include <variant>
#include <vector>

#include "bapel/core_variant.h"

struct Color : ::bapel::Variant<::bapel::Arm<"red", std::monostate>,
                                ::bapel::Arm<"green", std::monostate>,
                                ::bapel::Arm<"blue", std::monostate>> {
  using ::bapel::Variant<::bapel::Arm<"red", std::monostate>,
                         ::bapel::Arm<"green", std::monostate>,
                         ::bapel::Arm<"blue", std::monostate>>::Variant;
};
bool is_red(::Color c);
int64_t color_code(::Color c);
int64_t always_zero(::Color c);
