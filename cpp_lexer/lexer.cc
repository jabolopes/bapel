#include "lexer.h"

#include <cctype>
#include <stdexcept>
#include <iostream>

namespace bapel::lex {

static const std::vector<std::string_view> operators = {
	"<-",
	"->",
	"||",
	"&&",
	"!=",
	"::",
};

// Helper to match Go's unicode.IsSymbol for ASCII
static bool IsSymbol(char c) {
  return c == '+' || c == '<' || c == '=' || c == '>' || c == '|' || c == '~' || c == '^' || c == '$';
}

// Helper to match Go's unicode.IsLetter or '_'
static bool IsLetter(char c) {
  return std::isalpha(static_cast<unsigned char>(c)) || c == '_';
}

// Helper to match Go's unicode.IsDigit
static bool IsDigit(char c) {
  return std::isdigit(static_cast<unsigned char>(c));
}

// Helper to match Go's unicode.IsSpace
static bool IsSpace(char c) {
  return std::isspace(static_cast<unsigned char>(c));
}

// Helper to match Go's unicode.IsPrint
static bool IsPrint(char c) {
  return std::isprint(static_cast<unsigned char>(c));
}

// =============================================================================
// Scanner Implementation
// =============================================================================

Scanner::Scanner(std::string filename, std::string file)
    : filename_(std::move(filename)), file_(std::move(file)),
      left_(0), right_(0), left_line_num_(1), right_line_num_(1) {}

std::string_view Scanner::rest_file() const {
  if (right_ <= file_.size()) {
    return std::string_view(file_).substr(right_);
  }
  return "";
}

std::optional<char> Scanner::peek_rune() const {
  std::string_view rest = rest_file();
  if (rest.empty()) {
    return std::nullopt;
  }
  return rest.front();
}

bool Scanner::peek_string(std::string_view str) const {
  std::string_view rest = rest_file();
  return str.size() <= rest.size() && rest.substr(0, str.size()) == str;
}

std::optional<char> Scanner::read_rune() {
  auto r = peek_rune();
  if (!r) {
    return std::nullopt;
  }
  right_++;
  if (*r == '\n') {
    right_line_num_++;
  }
  return r;
}

bool Scanner::read_string(std::string_view str) {
  if (peek_string(str)) {
    for (size_t i = 0; i < str.size(); ++i) {
      read_rune();
    }
    return true;
  }
  return false;
}

std::string_view Scanner::current() const {
  if (left_ <= file_.size()) {
    return std::string_view(file_).substr(left_, right_ - left_);
  }
  return "";
}

void Scanner::ignore() {
  left_ = right_;
  left_line_num_ = right_line_num_;
}

std::string Scanner::pos_string() const {
  if (left_line_num_ == right_line_num_) {
    return "in \"" + filename_ + "\" in line " + std::to_string(left_line_num_);
  }
  return "in \"" + filename_ + "\" in lines " + std::to_string(left_line_num_) + "-" + std::to_string(right_line_num_);
}

// Helper to match Go's unexpected token message
static std::string UnexpectedTokenMessage(char c, int line) {
  std::string char_str;
  if (c == '\n') char_str = "'\\n'";
  else if (c == '\r') char_str = "'\\r'";
  else if (c == '\t') char_str = "'\\t'";
  else if (c == '\'') char_str = "'\\''";
  else if (c == '\\') char_str = "'\\\\'";
  else if (std::isprint(static_cast<unsigned char>(c))) {
    char_str = "'" + std::string(1, c) + "'";
  } else {
    char buf[10];
    snprintf(buf, sizeof(buf), "'\\x%02x'", static_cast<unsigned char>(c));
    char_str = buf;
  }
  
  return "unexpected token " + char_str + " (" + std::to_string(static_cast<int>(c)) + ") at line " + std::to_string(line);
}

// =============================================================================
// Lexer Implementation
// =============================================================================

Lexer::Lexer(std::string filename, std::string file)
    : scanner_(std::move(filename), std::move(file)),
      state_(StateId::Initial), string_delimiter_(0) {}

std::optional<Token> Lexer::NextToken() {
  while (tokens_.empty() && state_ != StateId::Stop) {
    state_ = Step(state_);
  }
  if (tokens_.empty()) {
    return std::nullopt;
  }
  Token tok = tokens_.front();
  tokens_.pop();
  return tok;
}

void Lexer::Emit(TokenType type) {
  tokens_.push(Token{scanner_.line_num(), type, std::string(scanner_.current())});
  scanner_.ignore();
}

void Lexer::Error(const std::string& err) {
  errors_.push_back(PosString() + ": " + err);
}

std::string Lexer::PosString() const {
  return scanner_.pos_string();
}

std::string Lexer::scan_err() const {
  std::string res;
  for (const auto& err : errors_) {
    if (!res.empty()) res += "\n";
    res += err;
  }
  return res;
}

StateId Lexer::Step(StateId state) {
  switch (state) {
    case StateId::Initial:      return InitialState();
    case StateId::LineComment:  return LineCommentState();
    case StateId::BlockComment: return BlockCommentState();
    case StateId::Whitespace:   return WhitespaceState();
    case StateId::Word:         return WordState();
    case StateId::Number:       return NumberState();
    case StateId::Symbol:       return SymbolState();
    case StateId::Rune:         return RuneState();
    case StateId::String:       return StringState();
    case StateId::Stop:         return StateId::Stop;
  }
  return StateId::Stop;
}

// =============================================================================
// Lexer States
// =============================================================================

StateId Lexer::InitialState() {
  auto r = PeekRune();
  if (!r) {
    return StateId::Stop; // EOF
  }

  char c = *r;

  if (c == '"') {
    string_delimiter_ = '"';
    string_name_ = "string";
    return StateId::String;
  }
  if (c == '`') {
    string_delimiter_ = '`';
    string_name_ = "raw string";
    return StateId::String;
  }

  if (PeekString("//")) {
    return StateId::LineComment;
  }
  if (PeekString("/*")) {
    return StateId::BlockComment;
  }

  for (const auto& op : operators) {
    if (ReadString(op)) {
      Emit(OperatorToken);
      return StateId::Initial;
    }
  }

  if (IsSpace(c)) {
    return StateId::Whitespace;
  }

  if (IsLetter(c)) {
    return StateId::Word;
  }

  if (IsDigit(c)) {
    return StateId::Number;
  }

  if (PeekString("'\\")) {
    return StateId::Rune;
  }

  if (IsSymbol(c)) {
    return StateId::Symbol;
  }

  if (IsPrint(c)) {
    ReadRune();
    Emit(static_cast<TokenType>(c));
    return StateId::Initial;
  }

  Error(UnexpectedTokenMessage(c, scanner_.line_num()));
  return StateId::Stop;
}

StateId Lexer::LineCommentState() {
  ReadRune(); // consume '/'
  ReadRune(); // consume '/'
  while (true) {
    auto r = ReadRune();
    Ignore();
    if (!r || *r == '\n') {
      return StateId::Initial;
    }
  }
}

StateId Lexer::BlockCommentState() {
  ReadString("/*");
  while (true) {
    if (ReadString("*/")) {
      break;
    }
    auto r = ReadRune();
    if (!r) {
      Error("unterminated block comment (/* ... */) starting at line " + std::to_string(scanner_.line_num()));
      return StateId::Stop;
    }
  }
  Ignore();
  return StateId::Initial;
}

StateId Lexer::WhitespaceState() {
  while (auto r = PeekRune()) {
    if (IsSpace(*r)) {
      ReadRune();
      Ignore();
    } else {
      break;
    }
  }
  return StateId::Initial;
}

StateId Lexer::WordState() {
  while (auto r = PeekRune()) {
    if (IsLetter(*r) || IsDigit(*r)) {
      ReadRune();
    } else {
      break;
    }
  }
  Emit(WordToken);
  return StateId::Initial;
}

StateId Lexer::NumberState() {
  while (auto r = PeekRune()) {
    if (IsDigit(*r)) {
      ReadRune();
    } else {
      break;
    }
  }
  Emit(NumberToken);
  return StateId::Initial;
}

StateId Lexer::SymbolState() {
  while (auto r = PeekRune()) {
    if (IsSymbol(*r)) {
      ReadRune();
    } else {
      break;
    }
  }
  if (Current().size() == 1) {
    Emit(static_cast<TokenType>(Current()[0]));
  } else {
    Emit(SymbolToken);
  }
  return StateId::Initial;
}

StateId Lexer::RuneState() {
  ReadRune(); // Consume '\''
  while (true) {
    auto r = ReadRune();
    if (!r) {
      Error("unterminated rune (' ... ') starting at line " + std::to_string(scanner_.line_num()));
      return StateId::Stop;
    }
    if (*r == '\'') {
      Emit(RuneToken);
      return StateId::Initial;
    }
  }
}

StateId Lexer::StringState() {
  ReadRune(); // Consume delimiter
  while (true) {
    auto r = ReadRune();
    if (!r) {
      Error("unterminated " + string_name_ + " (" + string_delimiter_ + " ... " + string_delimiter_ + ") starting at line " + std::to_string(scanner_.line_num()));
      return StateId::Stop;
    }
    if (*r == string_delimiter_) {
      Emit(StringToken);
      return StateId::Initial;
    }
  }
}

} // namespace bapel::lex
