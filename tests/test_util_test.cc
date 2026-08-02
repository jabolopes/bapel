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
  auto matches = tests::glob("tests/testdata/parse/in/*.in");
  EXPECT_TRUE(!matches.empty());
  for (const auto& m : matches) {
    EXPECT_TRUE(m.rfind("tests/testdata/parse/in/", 0) == 0);
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
