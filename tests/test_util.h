#pragma once

#include <algorithm>
#include <chrono>
#include <cstdlib>
#include <filesystem>
#include <fstream>
#include <functional>
#include <iomanip>
#include <iostream>
#include <memory>
#include <regex>
#include <set>
#include <sstream>
#include <string>
#include <utility>
#include <vector>

#include "comp/typecheck_unit.h"

template <typename T>
inline std::ostream& operator<<(std::ostream& os, const std::vector<T>& vec) {
  os << "[";
  for (size_t i = 0; i < vec.size(); ++i) {
    if (i > 0) os << ", ";
    os << vec[i];
  }
  os << "]";
  return os;
}

namespace tests {

namespace fs = std::filesystem;

// Global test runner configuration
struct TestConfig {
  bool regen = false;
  std::string filter;
  std::string glob_override;
  bool verbose = false;
  bool color = true;

  static TestConfig& instance() {
    static TestConfig cfg;
    return cfg;
  }
};

// Exception for fatal test assertions (ASSERT_*)
class TestFatalException : public std::exception {
 public:
  const char* what() const noexcept override { return "Test assertion failed (fatal)"; }
};

// Test Context for recording failures and running nested subtests
class TestContext {
 public:
  explicit TestContext(std::string name, TestContext* parent = nullptr)
      : name_(std::move(name)), parent_(parent) {}

  const std::string& name() const { return name_; }
  bool failed() const { return failed_; }
  const std::vector<std::string>& errors() const { return errors_; }

  void add_error(const std::string& err, const char* file, int line) {
    failed_ = true;
    std::stringstream ss;
    ss << file << ":" << line << ": " << err;
    errors_.push_back(ss.str());
    if (parent_) {
      parent_->failed_ = true;
    }
  }

  void log(const std::string& msg) {
    if (TestConfig::instance().verbose) {
      std::cout << "        [INFO] " << msg << "\n";
    }
  }

  void run(const std::string& sub_name, std::function<void(TestContext&)> sub_fn) {
    std::string full_name = name_ + "/" + sub_name;
    TestContext sub_ctx(full_name, this);
    try {
      sub_fn(sub_ctx);
    } catch (const TestFatalException&) {
      // Subtest aborted on fatal assert
    } catch (const std::exception& e) {
      sub_ctx.add_error(std::string("Unhandled exception: ") + e.what(), __FILE__, __LINE__);
    } catch (...) {
      sub_ctx.add_error("Unhandled unknown exception", __FILE__, __LINE__);
    }

    if (sub_ctx.failed()) {
      failed_ = true;
      for (const auto& err : sub_ctx.errors()) {
        errors_.push_back(err);
      }
    }
  }

 private:
  std::string name_;
  TestContext* parent_ = nullptr;
  bool failed_ = false;
  std::vector<std::string> errors_;
};

// Test case registration
using TestFunc = std::function<void(TestContext&)>;

struct TestCase {
  std::string suite;
  std::string name;
  TestFunc fn;
  std::string file;
  int line;

  std::string full_name() const {
    return suite + "." + name;
  }
};

class TestRegistry {
 public:
  static TestRegistry& instance() {
    static TestRegistry reg;
    return reg;
  }

  void register_test(const std::string& suite, const std::string& name, TestFunc fn, const char* file, int line) {
    tests_.push_back({suite, name, std::move(fn), file, line});
  }

  const std::vector<TestCase>& tests() const { return tests_; }

 private:
  std::vector<TestCase> tests_;
};

struct TestRegistrar {
  TestRegistrar(const std::string& suite, const std::string& name, TestFunc fn, const char* file, int line) {
    TestRegistry::instance().register_test(suite, name, std::move(fn), file, line);
  }
};

#define TEST_CONCAT_INNER(a, b) a##b
#define TEST_CONCAT(a, b) TEST_CONCAT_INNER(a, b)

#define TEST(suite, name) \
  void test_fn_##suite##_##name(::tests::TestContext& ctx); \
  static ::tests::TestRegistrar registrar_##suite##_##name( \
      #suite, #name, test_fn_##suite##_##name, __FILE__, __LINE__); \
  void test_fn_##suite##_##name(::tests::TestContext& ctx)

// Assertion macros
#define EXPECT_TRUE(cond) \
  do { \
    if (!(cond)) { \
      ctx.add_error(std::string("EXPECT_TRUE failed: ") + #cond, __FILE__, __LINE__); \
    } \
  } while (0)

#define EXPECT_FALSE(cond) \
  do { \
    if (cond) { \
      ctx.add_error(std::string("EXPECT_FALSE failed: ") + #cond, __FILE__, __LINE__); \
    } \
  } while (0)

#define EXPECT_EQ(val1, val2) \
  do { \
    auto _v1 = (val1); \
    auto _v2 = (val2); \
    if (!(_v1 == _v2)) { \
      std::stringstream _ss; \
      _ss << "EXPECT_EQ failed (" << #val1 << " == " << #val2 << "): \n" \
          << "  Expected: " << _v2 << "\n" \
          << "  Actual:   " << _v1; \
      ctx.add_error(_ss.str(), __FILE__, __LINE__); \
    } \
  } while (0)

#define EXPECT_NE(val1, val2) \
  do { \
    auto _v1 = (val1); \
    auto _v2 = (val2); \
    if (_v1 == _v2) { \
      std::stringstream _ss; \
      _ss << "EXPECT_NE failed (" << #val1 << " != " << #val2 << "): both values are " << _v1; \
      ctx.add_error(_ss.str(), __FILE__, __LINE__); \
    } \
  } while (0)

#define ASSERT_TRUE(cond) \
  do { \
    if (!(cond)) { \
      ctx.add_error(std::string("ASSERT_TRUE failed: ") + #cond, __FILE__, __LINE__); \
      throw ::tests::TestFatalException(); \
    } \
  } while (0)

#define ASSERT_FALSE(cond) \
  do { \
    if (cond) { \
      ctx.add_error(std::string("ASSERT_FALSE failed: ") + #cond, __FILE__, __LINE__); \
      throw ::tests::TestFatalException(); \
    } \
  } while (0)

#define ASSERT_EQ(val1, val2) \
  do { \
    auto _v1 = (val1); \
    auto _v2 = (val2); \
    if (!(_v1 == _v2)) { \
      std::stringstream _ss; \
      _ss << "ASSERT_EQ failed (" << #val1 << " == " << #val2 << "): \n" \
          << "  Expected: " << _v2 << "\n" \
          << "  Actual:   " << _v1; \
      ctx.add_error(_ss.str(), __FILE__, __LINE__); \
      throw ::tests::TestFatalException(); \
    } \
  } while (0)

// File & Path utilities
inline std::string read_file(const std::string& path, std::string* err = nullptr) {
  std::ifstream ifs(path, std::ios::binary);
  if (!ifs) {
    if (err) *err = "failed to open file for reading: " + path;
    return "";
  }
  return std::string((std::istreambuf_iterator<char>(ifs)), std::istreambuf_iterator<char>());
}

inline bool write_file(const std::string& path, const std::string& content, std::string* err = nullptr) {
  fs::path p(path);
  if (p.has_parent_path()) {
    std::error_code ec;
    fs::create_directories(p.parent_path(), ec);
    if (ec) {
      if (err) *err = "failed to create directory " + p.parent_path().string() + ": " + ec.message();
      return false;
    }
  }

  std::ofstream ofs(path, std::ios::binary | std::ios::trunc);
  if (!ofs) {
    if (err) *err = "failed to open file for writing: " + path;
    return false;
  }
  ofs.write(content.data(), content.size());
  return ofs.good();
}

inline std::string replace_extension(const std::string& filepath, const std::string& new_ext) {
  size_t last_dot = filepath.rfind('.');
  size_t last_slash = filepath.rfind('/');
  if (last_dot == std::string::npos || (last_slash != std::string::npos && last_dot < last_slash)) {
    return filepath + new_ext;
  }
  return filepath.substr(0, last_dot) + new_ext;
}

inline std::string trim_extension(const std::string& filepath) {
  size_t last_dot = filepath.rfind('.');
  size_t last_slash = filepath.rfind('/');
  if (last_dot == std::string::npos || (last_slash != std::string::npos && last_dot < last_slash)) {
    return filepath;
  }
  return filepath.substr(0, last_dot);
}

inline std::string path_base(const std::string& filepath) {
  size_t last_slash = filepath.rfind('/');
  if (last_slash == std::string::npos) return filepath;
  return filepath.substr(last_slash + 1);
}

inline std::string path_dir(const std::string& filepath) {
  size_t last_slash = filepath.rfind('/');
  if (last_slash == std::string::npos) return ".";
  if (last_slash == 0) return "/";
  return filepath.substr(0, last_slash);
}

inline std::string replace_string(std::string str, const std::string& from, const std::string& to) {
  size_t pos = str.find(from);
  if (pos != std::string::npos) {
    str.replace(pos, from.length(), to);
  }
  return str;
}

// Convert wildcard pattern like "tests/testdata/parse/in/*.in" to regex
inline std::regex pattern_to_regex(const std::string& pattern) {
  std::string regex_str;
  for (size_t i = 0; i < pattern.size(); ++i) {
    char c = pattern[i];
    if (c == '*') {
      regex_str += ".*";
    } else if (c == '?') {
      regex_str += ".";
    } else if (c == '.' || c == '+' || c == '(' || c == ')' || c == '[' || c == ']' ||
               c == '{' || c == '}' || c == '^' || c == '$' || c == '|' || c == '\\') {
      regex_str += '\\';
      regex_str += c;
    } else {
      regex_str += c;
    }
  }
  return std::regex("^" + regex_str + "$");
}

// Filesystem glob utility
inline std::vector<std::string> glob(const std::string& pattern_in) {
  std::string pattern = pattern_in;
  if (!TestConfig::instance().glob_override.empty()) {
    pattern = TestConfig::instance().glob_override;
  }

  // Find fixed directory prefix before wildcard characters
  size_t first_wildcard = pattern.find_first_of("*?");
  std::string search_dir;
  if (first_wildcard == std::string::npos) {
    search_dir = path_dir(pattern);
  } else {
    size_t last_slash_before = pattern.rfind('/', first_wildcard);
    if (last_slash_before == std::string::npos) {
      search_dir = ".";
    } else {
      search_dir = pattern.substr(0, last_slash_before);
    }
  }

  std::regex re = pattern_to_regex(pattern);
  std::vector<std::string> matches;

  std::error_code ec;
  if (fs::exists(search_dir, ec)) {
    for (const auto& entry : fs::recursive_directory_iterator(search_dir, ec)) {
      if (entry.is_regular_file(ec)) {
        std::string p = entry.path().string();
        // Normalize leading "./"
        if (p.rfind("./", 0) == 0 && pattern.rfind("./", 0) != 0) {
          p = p.substr(2);
        }
        if (std::regex_match(p, re)) {
          matches.push_back(p);
        }
      }
    }
  }

  std::sort(matches.begin(), matches.end());
  return matches;
}

// Line-by-line diff algorithm (LCS based)
inline std::vector<std::string> split_lines(const std::string& str) {
  std::vector<std::string> lines;
  std::stringstream ss(str);
  std::string line;
  while (std::getline(ss, line)) {
    if (!line.empty() && line.back() == '\r') {
      line.pop_back();
    }
    lines.push_back(line);
  }
  return lines;
}

inline std::string diff(const std::string& got, const std::string& want) {
  if (got == want) return "";

  std::vector<std::string> got_lines = split_lines(got);
  std::vector<std::string> want_lines = split_lines(want);

  size_t n = want_lines.size();
  size_t m = got_lines.size();

  // Compute Longest Common Subsequence matrix
  std::vector<std::vector<int>> dp(n + 1, std::vector<int>(m + 1, 0));
  for (size_t i = 1; i <= n; ++i) {
    for (size_t j = 1; j <= m; ++j) {
      if (want_lines[i - 1] == got_lines[j - 1]) {
        dp[i][j] = dp[i - 1][j - 1] + 1;
      } else {
        dp[i][j] = std::max(dp[i - 1][j], dp[i][j - 1]);
      }
    }
  }

  // Backtrack to build diff edits
  enum class EditType { Same, Insert, Delete };
  struct Edit {
    EditType type;
    std::string text;
  };
  std::vector<Edit> edits;

  int i = static_cast<int>(n);
  int j = static_cast<int>(m);
  while (i > 0 || j > 0) {
    if (i > 0 && j > 0 && want_lines[i - 1] == got_lines[j - 1]) {
      edits.push_back({EditType::Same, want_lines[i - 1]});
      --i;
      --j;
    } else if (j > 0 && (i == 0 || dp[i][j - 1] >= dp[i - 1][j])) {
      edits.push_back({EditType::Insert, got_lines[j - 1]});
      --j;
    } else if (i > 0 && (j == 0 || dp[i][j - 1] < dp[i - 1][j])) {
      edits.push_back({EditType::Delete, want_lines[i - 1]});
      --i;
    }
  }
  std::reverse(edits.begin(), edits.end());

  // Format unified diff output
  std::stringstream out;
  out << "--- Expected (want)\n+++ Actual   (got)\n";

  for (const auto& edit : edits) {
    if (edit.type == EditType::Same) {
      out << "  " << edit.text << "\n";
    } else if (edit.type == EditType::Delete) {
      out << "- " << edit.text << "\n";
    } else if (edit.type == EditType::Insert) {
      out << "+ " << edit.text << "\n";
    }
  }

  return out.str();
}

// Diff with golden file and automatic -regen support
inline bool diff_out_regen(const std::string& got, const std::string& want_file, std::string& out_diff) {
  if (TestConfig::instance().regen) {
    std::string write_err;
    if (!write_file(want_file, got, &write_err)) {
      out_diff = "Failed to regenerate golden file: " + write_err;
      return false;
    }
  }

  std::string read_err;
  std::string want = read_file(want_file, &read_err);
  if (!read_err.empty()) {
    out_diff = "Could not read golden file \"" + want_file + "\": " + read_err;
    return false;
  }

  out_diff = diff(got, want);
  return out_diff.empty();
}

inline bool diff_out_regen_file(const std::string& got_file, const std::string& want_file, std::string& out_diff) {
  std::string read_err;
  std::string got = read_file(got_file, &read_err);
  if (!read_err.empty()) {
    out_diff = "Could not read generated file \"" + got_file + "\": " + read_err;
    return false;
  }
  return diff_out_regen(got, want_file, out_diff);
}

struct TestDirectives {
  std::set<std::string> skipped_stages;
  comp::TypecheckOptions typecheck_options;
  std::string expect_error_stage;

  static int stage_index(const std::string& stage) {
    if (stage == "parse") return 1;
    if (stage == "normalize") return 2;
    if (stage == "infer") return 3;
    if (stage == "typecheck") return 4;
    if (stage == "cpp_codegen") return 5;
    if (stage == "cpp_compile") return 6;
    return 999;
  }

  bool expects_error(const std::string& stage) const {
    return expect_error_stage == stage;
  }

  bool should_run_stage(const std::string& stage) const {
    if (skipped_stages.find(stage) != skipped_stages.end()) {
      return false;
    }
    if (!expect_error_stage.empty()) {
      int stage_idx = stage_index(stage);
      int err_idx = stage_index(expect_error_stage);
      if (stage_idx > err_idx) {
        return false;
      }
    }
    return true;
  }

  static TestDirectives parse_from_file(const std::string& filepath) {
    TestDirectives dir;
    std::ifstream ifs(filepath);
    if (!ifs) return dir;

    std::string line;
    while (std::getline(ifs, line)) {
      size_t start = line.find_first_not_of(" \t\r\n");
      if (start == std::string::npos) continue;
      line = line.substr(start);

      // Stop scanning when reaching non-comment lines
      if (line.rfind("//", 0) != 0) {
        break;
      }

      size_t at_pos = line.find('@');
      if (at_pos == std::string::npos) continue;

      size_t colon_pos = line.find(':', at_pos);
      if (colon_pos == std::string::npos) continue;

      std::string directive = line.substr(at_pos + 1, colon_pos - (at_pos + 1));
      std::string args = line.substr(colon_pos + 1);

      size_t a_start = args.find_first_not_of(" \t\r\n");
      if (a_start != std::string::npos) {
        args = args.substr(a_start);
      }
      size_t a_end = args.find_last_not_of(" \t\r\n");
      if (a_end != std::string::npos) {
        args = args.substr(0, a_end + 1);
      }

      if (directive == "expect-error") {
        dir.expect_error_stage = args;
      } else if (directive == "skip-stage") {
        std::stringstream ss(args);
        std::string stage;
        while (std::getline(ss, stage, ',')) {
          size_t s = stage.find_first_not_of(" \t\r\n");
          size_t e = stage.find_last_not_of(" \t\r\n");
          if (s != std::string::npos && e != std::string::npos) {
            dir.skipped_stages.insert(stage.substr(s, e - s + 1));
          }
        }
      } else if (directive == "typecheck-option") {
        if (args.find("skip_undefined_term_checks=true") != std::string::npos) {
          dir.typecheck_options.skip_undefined_term_checks = true;
        } else if (args.find("skip_default_context=true") != std::string::npos) {
          dir.typecheck_options.skip_default_context = true;
        }
      }
    }
    return dir;
  }
};

} // namespace tests
