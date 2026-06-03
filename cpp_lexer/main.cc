#include "lexer.h"

#include <iostream>
#include <iomanip>
#include <fstream>
#include <sstream>

using namespace bapel::lex;

std::string TokenTypeToString(TokenType type) {
  switch (type) {
    case WordToken: return "Word";
    case NumberToken: return "Number";
    case RuneToken: return "Rune";
    case StringToken: return "String";
    case EOFToken: return "EOF";
    default:
      if (type > 0 && type < 128) {
        return "Char('" + std::string(1, static_cast<char>(type)) + "')";
      }
      return "Unknown(" + std::to_string(static_cast<int>(type)) + ")";
  }
}

int main(int argc, char* argv[]) {
  std::string filename = "stdin";
  std::string code;
  bool raw_mode = false;
  std::string file_to_read = "";

  // Simple argument parsing
  for (int i = 1; i < argc; ++i) {
    std::string arg = argv[i];
    if (arg == "--raw") {
      raw_mode = true;
    } else if (arg == "--filename" && i + 1 < argc) {
      filename = argv[++i];
    } else {
      file_to_read = arg;
    }
  }

  if (!file_to_read.empty()) {
    std::ifstream file(file_to_read);
    if (!file.is_open()) {
      std::cerr << "Failed to open file: " << file_to_read << "\n";
      return 1;
    }
    std::stringstream buffer;
    buffer << file.rdbuf();
    code = buffer.str();
    if (filename == "stdin") {
      filename = file_to_read; // Use it as filename if not explicitly set
    }
  } else {
    std::stringstream buffer;
    buffer << std::cin.rdbuf();
    code = buffer.str();
  }

  Lexer lexer(filename, code);
  
  while (auto tok = lexer.NextToken()) {
    if (raw_mode) {
      // Print header: line type size
      std::cout << tok->line_num << " " 
                << static_cast<int>(tok->type) << " " 
                << tok->value.size() << "\n";
      // Write raw value bytes directly
      std::cout.write(tok->value.data(), tok->value.size());
    } else {
      // Keep the human-readable debug output
      static int line = 0;
      if (line != tok->line_num) {
        line = tok->line_num;
        std::cout << "LINE " << line << ":\n";
      }
      if (tok->value.size() > 0) {
        std::cout << "TOKEN: " << tok->value << "\n";
      } else {
        std::cout << "TOKEN: " << TokenTypeToString(tok->type) << "\n";
      }
    }
  }

  std::string err = lexer.scan_err();
  if (!err.empty()) {
    std::cerr << err << "\n";
    return 1;
  }

  return 0;
}
