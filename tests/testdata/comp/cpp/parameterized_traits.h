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

#include "../tests/testdata/comp/in/ptr.h"
#include "../tests/testdata/comp/in/vector.h"

namespace traits {
template <typename Self, typename elem>
struct Indexable;
}