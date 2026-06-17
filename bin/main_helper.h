#pragma once
#include <string>
#include <vector>
#include <tuple>

namespace cli {


// 2. Process Execution
// @bpl: pub cli::exec: (String, Vector String) -> (i64, String)
std::tuple<int64_t, std::string> exec(std::string cmd, const std::vector<std::string>& args);



// 4. Parsing Structures and Functions



// @bpl: pub type cli::SourceFileInfo = struct { importModules Vector String, implFiles Vector String }
struct SourceFileInfo {
  std::vector<std::string> importModules;
  std::vector<std::string> implFiles;
};

// @bpl: pub cli::parseSourceFileFlat: String -> cli::SourceFileInfo
SourceFileInfo parseSourceFileFlat(std::string text);




















} // namespace cli
