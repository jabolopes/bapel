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

#include "comp/testdata/in/ptr.h"
#include "comp/testdata/in/vector.h"

namespace traits {
template <typename Self, typename elem>
struct Indexable;
}