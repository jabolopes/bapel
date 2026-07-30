
// Generated from cpp_parser/bapel.g4 by ANTLR 4.9.2

#pragma once


#include "antlr4-runtime.h"




class  bapelLexer : public antlr4::Lexer {
public:
  enum {
    WORKSPACE = 1, PACKAGES = 2, PREFIX = 3, MODULE = 4, IMPLEMENTS = 5, 
    IMPORTS = 6, IMPLS = 7, FLAGS = 8, IN = 9, PUB = 10, DECL = 11, FN = 12, 
    TYPE = 13, FORALL = 14, STRUCT = 15, VARIANT = 16, MATCH = 17, SET = 18, 
    LET = 19, RETURN = 20, IF = 21, ELSE = 22, FOR = 23, TRAIT = 24, IMPL = 25, 
    DOUBLE_COLON = 26, ARROW = 27, FAT_ARROW = 28, LARROW = 29, OR = 30, 
    AND = 31, NE = 32, EQ = 33, GE = 34, LE = 35, GT = 36, LT = 37, PLUS = 38, 
    MINUS = 39, MUL = 40, DIV = 41, NOT = 42, AMP = 43, DOT = 44, ASSIGN = 45, 
    COMMA = 46, SEMI = 47, COLON = 48, LBRACE = 49, RBRACE = 50, LPAREN = 51, 
    RPAREN = 52, LBRACKET = 53, RBRACKET = 54, SINGLE_QUOTE = 55, IDENTIFIER = 56, 
    INT_LITERAL = 57, FLOAT_LITERAL = 58, RUNE_LITERAL = 59, STRING_LITERAL = 60, 
    RAW_STRING_LITERAL = 61, UNTERMINATED_RUNE_LITERAL = 62, UNTERMINATED_STRING_LITERAL = 63, 
    UNTERMINATED_RAW_STRING_LITERAL = 64, UNTERMINATED_BLOCK_COMMENT = 65, 
    WS = 66, LINE_COMMENT = 67, BLOCK_COMMENT = 68
  };

  explicit bapelLexer(antlr4::CharStream *input);
  ~bapelLexer();

  virtual std::string getGrammarFileName() const override;
  virtual const std::vector<std::string>& getRuleNames() const override;

  virtual const std::vector<std::string>& getChannelNames() const override;
  virtual const std::vector<std::string>& getModeNames() const override;
  virtual const std::vector<std::string>& getTokenNames() const override; // deprecated, use vocabulary instead
  virtual antlr4::dfa::Vocabulary& getVocabulary() const override;

  virtual const std::vector<uint16_t> getSerializedATN() const override;
  virtual const antlr4::atn::ATN& getATN() const override;

private:
  static std::vector<antlr4::dfa::DFA> _decisionToDFA;
  static antlr4::atn::PredictionContextCache _sharedContextCache;
  static std::vector<std::string> _ruleNames;
  static std::vector<std::string> _tokenNames;
  static std::vector<std::string> _channelNames;
  static std::vector<std::string> _modeNames;

  static std::vector<std::string> _literalNames;
  static std::vector<std::string> _symbolicNames;
  static antlr4::dfa::Vocabulary _vocabulary;
  static antlr4::atn::ATN _atn;
  static std::vector<uint16_t> _serializedATN;


  // Individual action functions triggered by action() above.

  // Individual semantic predicate functions triggered by sempred() above.

  struct Initializer {
    Initializer();
  };
  static Initializer _init;
};

