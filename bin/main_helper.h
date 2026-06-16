#pragma once
#include <string>
#include <vector>
#include <tuple>

namespace cli {


// 2. Process Execution
// @bpl: pub cli::exec: (String, core::Vector String) -> (i64, String)
std::tuple<int64_t, std::string> exec(std::string cmd, const std::vector<std::string>& args);



// 4. Parsing Structures and Functions

// @bpl: pub type cli::PackageMapping = struct { is_prefix bool, name String, path String }
struct PackageMapping {
  bool is_prefix;
  std::string name;
  std::string path;
};

// @bpl: pub cli::parseWorkspaceFlat: String -> core::Vector cli::PackageMapping
std::vector<PackageMapping> parseWorkspaceFlat(std::string text);

// @bpl: pub type cli::SourceFileInfo = struct { importModules core::Vector String, implFiles core::Vector String }
struct SourceFileInfo {
  std::vector<std::string> importModules;
  std::vector<std::string> implFiles;
};

// @bpl: pub cli::parseSourceFileFlat: String -> cli::SourceFileInfo
SourceFileInfo parseSourceFileFlat(std::string text);




// @bpl: pub cli::replaceSeparator: (String, String, String) -> String
std::string replaceSeparator(std::string s, std::string from, std::string to);

















} // namespace cli
