#pragma once
#include <string>
#include <vector>
#include <tuple>

namespace cli {

// 1. CLI Args Access
// @bpl: pub cli::getArgCount: () -> i64
int64_t getArgCount();

// @bpl: pub cli::getArg: i64 -> String
std::string getArg(int64_t i);

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

// @bpl: pub cli::joinPath: (String, String) -> String
std::string joinPath(std::string a, std::string b);

// @bpl: pub cli::replaceSeparator: (String, String, String) -> String
std::string replaceSeparator(std::string s, std::string from, std::string to);

// @bpl: pub cli::isPrefixOf: (String, String) -> bool
bool isPrefixOf(std::string prefix, std::string s);

// @bpl: pub cli::concat: (String, String) -> String
std::string concat(std::string a, std::string b);




// @bpl: pub cli::addTarget: (String, String, core::Vector String, core::Vector String, core::Vector String) -> ()
void addTarget(std::string type, std::string name, const std::vector<std::string>& srcs, const std::vector<std::string>& hdrs, const std::vector<std::string>& deps);

// @bpl: pub cli::writeBuildFile: () -> bool
bool writeBuildFile();

// @bpl: pub cli::ensureWorkspaceSetup: () -> bool
bool ensureWorkspaceSetup();



// @bpl: pub cli::getSubArgs: i64 -> core::Vector String
std::vector<std::string> getSubArgs(int64_t start);



} // namespace cli
