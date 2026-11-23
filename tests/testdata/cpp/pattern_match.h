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

BAPEL_VARIANT(Color, (red, (std::monostate)), (green, (std::monostate)),
              (blue, (std::monostate)));
bool is_red(::Color c);
int64_t color_code(::Color c);
int64_t always_zero(::Color c);
