#pragma once
#include <vector>
#include <string>
#include <variant>

// @bpl: pub type args::Argc
// @bpl: pub type args::Argv
namespace args {

using Argc = int;
using Argv = char**;

namespace internal {
inline std::vector<std::string> main_args;
} // namespace internal

// @bpl: pub args::init: (args::Argc, args::Argv) -> ()
inline std::monostate init(Argc argc, Argv argv) {
    internal::main_args.clear();
    for (int i = 0; i < argc; ++i) {
        internal::main_args.push_back(argv[i]);
    }
    return std::monostate();
}

// @bpl: pub args::get_args: () -> core::Vector String
inline const std::vector<std::string>& get_args() {
    return internal::main_args;
}

} // namespace args
