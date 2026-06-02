#ifndef CPP_LEXER_LEXER_H_
#define CPP_LEXER_LEXER_H_

#include <string>
#include <string_view>
#include <vector>
#include <queue>
#include <optional>
#include <cstdint>

namespace bapel::lex {

enum TokenType : int32_t {
  WordToken = 0,
  NumberToken = 1,
  RuneToken = 2,
  StringToken = 3,
  SymbolToken = WordToken,
  OperatorToken = WordToken,
  EOFToken = -1,
};

struct Token {
  int line_num;
  TokenType type;
  std::string value;

  bool operator==(const Token& o) const {
    return line_num == o.line_num && type == o.type && value == o.value;
  }
};

class Scanner {
public:
  Scanner(std::string filename, std::string file);

  std::string_view rest_file() const;
  int line_num() const { return left_line_num_; }
  
  std::optional<char> peek_rune() const;
  bool peek_string(std::string_view str) const;
  
  std::optional<char> read_rune();
  bool read_string(std::string_view str);
  
  std::string_view current() const;
  void ignore();

  std::string pos_string() const;

private:
  std::string filename_;
  std::string file_;
  size_t left_;
  size_t right_;
  int left_line_num_;
  int right_line_num_;
};

enum class StateId {
  Initial,
  LineComment,
  BlockComment,
  Whitespace,
  Word,
  Number,
  Symbol,
  Rune,
  String,
  Stop, // Equivalent to returning nil/nullptr
};

class Lexer {
public:
  Lexer(std::string filename, std::string file);

  std::optional<Token> NextToken();
  std::string scan_err() const;
  std::string PosString() const;

private:
  // Helpers for states
  std::optional<char> PeekRune() const { return scanner_.peek_rune(); }
  bool PeekString(std::string_view str) const { return scanner_.peek_string(str); }
  std::optional<char> ReadRune() { return scanner_.read_rune(); }
  bool ReadString(std::string_view str) { return scanner_.read_string(str); }
  std::string_view Current() const { return scanner_.current(); }
  void Ignore() { scanner_.ignore(); }
  void Emit(TokenType type);
  void Error(const std::string& err);

  // State execution step
  StateId Step(StateId state);

  // State implementations (now returning next StateId)
  StateId InitialState();
  StateId LineCommentState();
  StateId BlockCommentState();
  StateId WhitespaceState();
  StateId WordState();
  StateId NumberState();
  StateId SymbolState();
  StateId RuneState();
  StateId StringState();

  Scanner scanner_;
  std::queue<Token> tokens_;
  std::vector<std::string> errors_;
  StateId state_;

  // Parameters for StringState
  char string_delimiter_;
  std::string string_name_;
};

} // namespace bapel::lex

#endif // CPP_LEXER_LEXER_H_
