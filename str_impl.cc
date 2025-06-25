module;

#include <string>
#include <string_view>

export module str:str_impl;

export namespace str {

// @bpl: export type str.String
using String = std::string;

}  // namespace str
