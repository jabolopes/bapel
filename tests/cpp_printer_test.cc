#include "tests/test_util.h"
#include "bin/codegen_impl.h"

#include <chrono>
#include <cstdlib>
#include <filesystem>
#include <string>
#include <vector>

namespace fs = std::filesystem;

TEST(CppPrinterTest, GoldenFiles) {
  auto matches = tests::glob("tests/testdata/comp/in/*.in");
  ASSERT_FALSE(matches.empty());

  for (const auto& inFile : matches) {
    if (tests::path_base(inFile) == "order.in") {
      continue;
    }

    ctx.run(inFile, [&](tests::TestContext& sub_ctx) {
      std::error_code ec;
      fs::path temp_dir = fs::temp_directory_path() / ("bapel_test_codegen_" + std::to_string(std::rand()) + "_" + std::to_string(std::chrono::system_clock::now().time_since_epoch().count()));
      fs::create_directories(temp_dir, ec);
      if (ec) {
        sub_ctx.add_error("Failed to create temporary directory: " + ec.message(), __FILE__, __LINE__);
        return;
      }

      std::string baseName = tests::trim_extension(tests::path_base(inFile));
      std::string gotFilenameBase = (temp_dir / baseName).string();

      int status = codegen::compile_unit(inFile, gotFilenameBase);
      if (status != 0) {
        sub_ctx.add_error("codegen::compile_unit failed for " + inFile, __FILE__, __LINE__);
        fs::remove_all(temp_dir, ec);
        return;
      }

      std::string wantFileH = tests::replace_string(tests::replace_extension(inFile, ".h"), "/in/", "/cpp/");
      std::string wantFilePrivH = tests::replace_string(tests::replace_extension(inFile, "_private.h"), "/in/", "/cpp/");
      std::string wantFileCc = tests::replace_string(tests::replace_extension(inFile, ".cc"), "/in/", "/cpp/");

      std::string diff_str;
      if (!tests::diff_out_regen_file(gotFilenameBase + ".h", wantFileH, diff_str)) {
        sub_ctx.add_error(".h diff in " + inFile + ":\n" + diff_str, __FILE__, __LINE__);
      }

      if (!tests::diff_out_regen_file(gotFilenameBase + "_private.h", wantFilePrivH, diff_str)) {
        sub_ctx.add_error("_private.h diff in " + inFile + ":\n" + diff_str, __FILE__, __LINE__);
      }

      if (!tests::diff_out_regen_file(gotFilenameBase + ".cc", wantFileCc, diff_str)) {
        sub_ctx.add_error(".cc diff in " + inFile + ":\n" + diff_str, __FILE__, __LINE__);
      }

      fs::remove_all(temp_dir, ec);
    });
  }
}

TEST(CppPrinterTest, IsValidCpp) {
  auto matches = tests::glob("tests/testdata/comp/cpp/*.cc");
  ASSERT_FALSE(matches.empty());

  for (const auto& inFile : matches) {
    std::string base = tests::path_base(inFile);
    if (base == "array.cc" || base == "context1.cc" || base == "loops.cc" || base == "polymorphism.cc") {
      // These tests import 'bapel.core' which requires external runtime headers.
      continue;
    }

    ctx.run(inFile, [&](tests::TestContext& sub_ctx) {
      std::error_code ec;
      fs::path temp_obj = fs::temp_directory_path() / ("bapel_obj_" + std::to_string(std::rand()) + "_" + base + ".o");

      std::string inDir = tests::path_dir(inFile);
      std::string cmd = "clang++ -std=c++17 -c " + inFile + " -o " + temp_obj.string() + " -I" + inDir + " -I. 2>&1";

      FILE* fp = popen(cmd.c_str(), "r");
      if (!fp) {
        sub_ctx.add_error("Failed to execute clang++ for " + inFile, __FILE__, __LINE__);
        return;
      }

      char buf[1024];
      std::string output;
      while (fgets(buf, sizeof(buf), fp) != nullptr) {
        output += buf;
      }
      int ret = pclose(fp);
      fs::remove(temp_obj, ec);

      if (ret != 0) {
        sub_ctx.add_error("clang++ compilation failed for " + inFile + ":\n" + output, __FILE__, __LINE__);
      }
    });
  }
}
