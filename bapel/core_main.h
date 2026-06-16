#pragma once
#include <vector>
#include <string>
#include <variant>

// @bpl: pub type core::Argc
// @bpl: pub type core::Argv
namespace core {

using Argc = int;
using Argv = char**;

namespace internal {
inline std::vector<std::string> main_args;
} // namespace internal

// @bpl: pub core::init: (core::Argc, core::Argv) -> ()
inline std::monostate init(Argc argc, Argv argv) {
    internal::main_args.clear();
    for (int i = 0; i < argc; ++i) {
        internal::main_args.push_back(argv[i]);
    }
    return std::monostate();
}

// @bpl: pub core::get_args: () -> core::Vector String
inline const std::vector<std::string>& get_args() {
    return internal::main_args;
}

} // namespace core
