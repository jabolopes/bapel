
#include "context1_private.h"

namespace core {
std::monostate print() { return std::monostate(); }

}  // namespace core
std::monostate callPrint() { return ::core::print(); }
