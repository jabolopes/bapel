#include "cpp_parser/parser.h"
#include <iostream>
#include <string>
#include <vector>

int main(int argc, char* argv[]) {
  std::vector<std::string> args;
  for (int i = 1; i < argc; ++i) {
    args.push_back(argv[i]);
  }
  auto [code, out] = bapel_parser::run(args);
  if (code == 0) {
    std::cout << out;
  } else {
    std::cerr << out;
  }
  return static_cast<int>(code);
}

