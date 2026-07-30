#pragma once

#include "antlr4-runtime.h"
#include "generated/bapelParser.h"
#include <algorithm>
#include <cstdio>
#include <sstream>
#include <string>
#include <vector>

namespace ast {

class BapelErrorListener : public antlr4::BaseErrorListener {
public:
  explicit BapelErrorListener(std::string filename) : filename_(std::move(filename)) {}

  const std::vector<std::string>& errors() const { return errors_; }
  bool has_errors() const { return !errors_.empty(); }

  void syntaxError(antlr4::Recognizer* recognizer, antlr4::Token* offendingSymbol,
                   size_t line, size_t charPositionInLine,
                   const std::string& msg_orig, std::exception_ptr e) override {
    (void)recognizer;
    (void)charPositionInLine;
    (void)e;
    std::string msg = msg_orig;

    if (offendingSymbol != nullptr) {
      size_t tokenType = offendingSymbol->getType();
      size_t startLine = offendingSymbol->getLine();
      std::string text = offendingSymbol->getText();
      size_t newlines = std::count(text.begin(), text.end(), '\n');
      size_t endLine = startLine + newlines;

      if (tokenType == bapelParser::UNTERMINATED_STRING_LITERAL) {
        if (startLine == endLine) {
          msg = "unterminated string (\" ... \") starting at line " + std::to_string(startLine);
        } else {
          errors_.push_back("in \"" + filename_ + "\" in lines " + std::to_string(startLine) + "-" +
                            std::to_string(endLine) + ": unterminated string (\" ... \") starting at line " +
                            std::to_string(startLine));
          return;
        }
      } else if (tokenType == bapelParser::UNTERMINATED_RAW_STRING_LITERAL) {
        if (startLine == endLine) {
          msg = "unterminated raw string (` ... `) starting at line " + std::to_string(startLine);
        } else {
          errors_.push_back("in \"" + filename_ + "\" in lines " + std::to_string(startLine) + "-" +
                            std::to_string(endLine) + ": unterminated raw string (` ... `) starting at line " +
                            std::to_string(startLine));
          return;
        }
      } else if (tokenType == bapelParser::UNTERMINATED_BLOCK_COMMENT) {
        if (startLine == endLine) {
          msg = "unterminated block comment (/* ... */) starting at line " + std::to_string(startLine);
        } else {
          errors_.push_back("in \"" + filename_ + "\" in lines " + std::to_string(startLine) + "-" +
                            std::to_string(endLine) + ": unterminated block comment (/* ... */) starting at line " +
                            std::to_string(startLine));
          return;
        }
      } else if (tokenType == bapelParser::UNTERMINATED_RUNE_LITERAL) {
        if (startLine == endLine) {
          msg = "unterminated rune (' ... ') starting at line " + std::to_string(startLine);
        } else {
          errors_.push_back("in \"" + filename_ + "\" in lines " + std::to_string(startLine) + "-" +
                            std::to_string(endLine) + ": unterminated rune (' ... ') starting at line " +
                            std::to_string(startLine));
          return;
        }
      }
    }

    const std::string prefix = "token recognition error at: ";
    if (msg.rfind(prefix, 0) == 0) {
      std::string raw = msg.substr(prefix.size());
      if (raw.size() >= 2 && raw.front() == '\'' && raw.back() == '\'') {
        raw = raw.substr(1, raw.size() - 2);
      }
      if (raw.rfind("\\x", 0) == 0) {
        int val = 0;
        if (sscanf(raw.c_str(), "\\x%x", &val) == 1) {
          char buf[128];
          snprintf(buf, sizeof(buf), "unexpected token '\\x%02x' (%d) at line %zu", val, val, line);
          msg = buf;
        }
      } else if (!raw.empty()) {
        unsigned char char_val = static_cast<unsigned char>(raw[0]);
        char buf[128];
        if (char_val < 32 || char_val > 126) {
          snprintf(buf, sizeof(buf), "unexpected token '\\x%02x' (%d) at line %zu", char_val, static_cast<int>(char_val), line);
        } else {
          snprintf(buf, sizeof(buf), "unexpected token '%s' (%d) at line %zu", raw.c_str(), static_cast<int>(char_val), line);
        }
        msg = buf;
      }
    }

    errors_.push_back("in \"" + filename_ + "\" in line " + std::to_string(line) + ": " + msg);
  }

private:
  std::string filename_;
  std::vector<std::string> errors_;
};

} // namespace ast
