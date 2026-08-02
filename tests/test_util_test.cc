#include "tests/test_util.h"

TEST(TestUtilTest, FilePathUtilities) {
  EXPECT_EQ(tests::path_base("foo/bar/baz.bpl"), "baz.bpl");
  EXPECT_EQ(tests::path_base("baz.bpl"), "baz.bpl");
  EXPECT_EQ(tests::path_dir("foo/bar/baz.bpl"), "foo/bar");
  EXPECT_EQ(tests::path_dir("baz.bpl"), ".");

  EXPECT_EQ(tests::trim_extension("foo/bar/baz.bpl"), "foo/bar/baz");
  EXPECT_EQ(tests::replace_extension("foo/bar/baz.in", ".out"), "foo/bar/baz.out");
  EXPECT_EQ(tests::replace_extension("baz", ".out"), "baz.out");
}

TEST(TestUtilTest, DiffIdenticalStrings) {
  std::string s = "line 1\nline 2\nline 3\n";
  std::string d = tests::diff(s, s);
  EXPECT_TRUE(d.empty());
}

TEST(TestUtilTest, DiffDifferentStrings) {
  std::string want = "line 1\nline 2\nline 3\n";
  std::string got = "line 1\nline 2 modified\nline 3\n";
  std::string d = tests::diff(got, want);
  EXPECT_FALSE(d.empty());
  EXPECT_TRUE(d.find("- line 2") != std::string::npos);
  EXPECT_TRUE(d.find("+ line 2 modified") != std::string::npos);
}

TEST(TestUtilTest, GlobMatching) {
  auto matches = tests::glob("tests/testdata/in/*.in");
  EXPECT_TRUE(!matches.empty());
  for (const auto& m : matches) {
    EXPECT_TRUE(m.rfind("tests/testdata/in/", 0) == 0);
    EXPECT_TRUE(m.substr(m.size() - 3) == ".in");
  }
}

TEST(TestUtilTest, ReadAndWriteFile) {
  std::string tmp_path = "/tmp/bapel_test_util_test.txt";
  std::string content = "Hello Bapel Test Framework!\nLine 2\n";
  std::string err;

  EXPECT_TRUE(tests::write_file(tmp_path, content, &err));
  EXPECT_TRUE(err.empty());

  std::string read_back = tests::read_file(tmp_path, &err);
  EXPECT_TRUE(err.empty());
  EXPECT_EQ(read_back, content);
}

TEST(TestUtilTest, DirectivesParsing) {
  std::string tmp_path = "/tmp/bapel_test_directives.in";
  std::string content =
      "// @expect-error: typecheck\n"
      "// @skip-stage: cpp_codegen, cpp_compile\n"
      "// @typecheck-option: skip_undefined_term_checks=true\n"
      "module test_mod\n";
  std::string err;
  EXPECT_TRUE(tests::write_file(tmp_path, content, &err));

  auto dir = tests::TestDirectives::parse_from_file(tmp_path);
  EXPECT_TRUE(dir.expects_error("typecheck"));
  EXPECT_FALSE(dir.expects_error("parse"));
  EXPECT_TRUE(dir.typecheck_options.skip_undefined_term_checks);
  EXPECT_FALSE(dir.typecheck_options.skip_default_context);

  EXPECT_TRUE(dir.should_run_stage("parse"));
  EXPECT_TRUE(dir.should_run_stage("normalize"));
  EXPECT_TRUE(dir.should_run_stage("infer"));
  EXPECT_TRUE(dir.should_run_stage("typecheck"));
  EXPECT_FALSE(dir.should_run_stage("cpp_codegen"));
  EXPECT_FALSE(dir.should_run_stage("cpp_compile"));
}

