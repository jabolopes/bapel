#include "tests/test_util.h"

#include <chrono>
#include <iomanip>
#include <iostream>
#include <string>
#include <vector>

namespace tests {

void print_help(const char* prog_name) {
  std::cout << "Usage: " << prog_name << " [options]\n\n"
            << "Options:\n"
            << "  --filter=PATTERN   Run only tests matching PATTERN (substring or glob)\n"
            << "  -regen, --regen    Regenerate golden test output files\n"
            << "  --glob=PATTERN     Override glob pattern for test inputs\n"
            << "  -v, --verbose      Enable verbose output logging\n"
            << "  --no-color         Disable ANSI color output\n"
            << "  -h, --help         Show this help message\n";
}

bool matches_filter(const std::string& test_name, const std::string& filter) {
  if (filter.empty()) return true;
  if (test_name.find(filter) != std::string::npos) return true;
  try {
    std::regex re = pattern_to_regex(filter);
    if (std::regex_search(test_name, re)) return true;
  } catch (...) {
  }
  return false;
}

int run_all_tests(int argc, char* argv[]) {
  auto& cfg = TestConfig::instance();

  for (int i = 1; i < argc; ++i) {
    std::string arg = argv[i];
    if (arg == "-regen" || arg == "--regen") {
      cfg.regen = true;
    } else if (arg.rfind("--filter=", 0) == 0) {
      cfg.filter = arg.substr(9);
    } else if (arg.rfind("-filter=", 0) == 0) {
      cfg.filter = arg.substr(8);
    } else if (arg.rfind("--glob=", 0) == 0) {
      cfg.glob_override = arg.substr(7);
    } else if (arg.rfind("-glob=", 0) == 0) {
      cfg.glob_override = arg.substr(6);
    } else if (arg == "-v" || arg == "--verbose") {
      cfg.verbose = true;
    } else if (arg == "--no-color") {
      cfg.color = false;
    } else if (arg == "-h" || arg == "--help") {
      print_help(argv[0]);
      return 0;
    } else {
      std::cerr << "Unknown argument: " << arg << "\n";
      print_help(argv[0]);
      return 1;
    }
  }

  // Check if standard output is a terminal for colors
  const char* term = std::getenv("TERM");
  if (!term || std::string(term) == "dumb") {
    cfg.color = false;
  }

  const std::string red = cfg.color ? "\033[31m" : "";
  const std::string green = cfg.color ? "\033[32m" : "";
  const std::string yellow = cfg.color ? "\033[33m" : "";
  const std::string reset = cfg.color ? "\033[0m" : "";

  const auto& all_tests = TestRegistry::instance().tests();
  std::vector<const TestCase*> matching_tests;
  for (const auto& t : all_tests) {
    if (matches_filter(t.full_name(), cfg.filter)) {
      matching_tests.push_back(&t);
    }
  }

  std::cout << green << "[==========]" << reset << " Running " << matching_tests.size()
            << " tests from " << all_tests.size() << " registered tests.\n";

  auto overall_start = std::chrono::steady_clock::now();
  int passed_count = 0;
  int failed_count = 0;
  std::vector<std::pair<std::string, std::vector<std::string>>> failed_details;

  for (const auto* test_ptr : matching_tests) {
    const auto& test = *test_ptr;
    std::cout << green << "[ RUN      ]" << reset << " " << test.full_name() << "\n";

    TestContext ctx(test.full_name());
    auto test_start = std::chrono::steady_clock::now();

    try {
      test.fn(ctx);
    } catch (const TestFatalException&) {
      // Handled via ctx
    } catch (const std::exception& e) {
      ctx.add_error(std::string("Unhandled exception: ") + e.what(), test.file.c_str(), test.line);
    } catch (...) {
      ctx.add_error("Unhandled unknown exception", test.file.c_str(), test.line);
    }

    auto test_end = std::chrono::steady_clock::now();
    auto elapsed_ms = std::chrono::duration_cast<std::chrono::milliseconds>(test_end - test_start).count();

    if (ctx.failed()) {
      ++failed_count;
      std::cout << red << "[  FAILED  ]" << reset << " " << test.full_name()
                << " (" << elapsed_ms << " ms)\n";
      for (const auto& err : ctx.errors()) {
        std::cout << "  " << err << "\n";
      }
      failed_details.push_back({test.full_name(), ctx.errors()});
    } else {
      ++passed_count;
      std::cout << green << "[       OK ]" << reset << " " << test.full_name()
                << " (" << elapsed_ms << " ms)\n";
    }
  }

  auto overall_end = std::chrono::steady_clock::now();
  auto total_ms = std::chrono::duration_cast<std::chrono::milliseconds>(overall_end - overall_start).count();

  std::cout << green << "[==========]" << reset << " " << matching_tests.size()
            << " tests ran. (" << total_ms << " ms total)\n";
  std::cout << green << "[  PASSED  ]" << reset << " " << passed_count << " tests.\n";

  if (failed_count > 0) {
    std::cout << red << "[  FAILED  ]" << reset << " " << failed_count << " tests, listed below:\n";
    for (const auto& f : failed_details) {
      std::cout << red << "[  FAILED  ]" << reset << " " << f.first << "\n";
    }
    return 1;
  }

  return 0;
}

} // namespace tests

int main(int argc, char* argv[]) {
  return tests::run_all_tests(argc, argv);
}
