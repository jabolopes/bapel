#pragma once
#include <string>
#include <vector>
#include <tuple>

namespace cli {


// 2. Process Execution
// @bpl: pub cli::exec: (String, Vector String) -> (i64, String)
std::tuple<int64_t, std::string> exec(std::string cmd, const std::vector<std::string>& args);























} // namespace cli
