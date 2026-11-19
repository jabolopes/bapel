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

struct Color
    : std::variant<std::monostate /* red */, std::monostate /* green */,
                   std::monostate /* blue */> {
  using std::variant<std::monostate /* red */, std::monostate /* green */,
                     std::monostate /* blue */>::variant;
};
bool is_red(::Color c);
int64_t color_code(::Color c);
int64_t always_zero(::Color c);
