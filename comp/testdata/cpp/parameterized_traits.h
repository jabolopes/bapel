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

#include "testdata/in/ptr.h"
#include "testdata/in/vector.h"

namespace traits {
template <typename Self, typename elem>
struct Indexable;
}