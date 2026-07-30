#pragma once

#include <algorithm>
#include <cstdint>
#include <functional>
#include <memory>
#include <optional>
#include <string>
#include <string_view>
#include <tuple>
#include <utility>
#include <vector>

namespace ir {

inline const std::string NamespaceSeparator = "::";
inline const std::string ModuleIDSeparator = ".";

enum class Position {
  TypePosition,
  BindPosition
};

enum class PrinterMode {
  ModePublicHeader,
  ModePrivateHeader,
  ModeSource
};

enum class IrTypeCase {
  AppType = 0,
  ArrayType = 1,
  ExistVarType = 2,
  ForallType = 3,
  FunType = 4,
  LambdaType = 5,
  NameType = 6,
  StructType = 7,
  TupleType = 8,
  VariantType = 9,
  VarType = 10
};

enum class IrTermCase {
  AppTermTerm = 0,
  AppTypeTerm = 1,
  AssignTerm = 2,
  BlockTerm = 3,
  ConstTerm = 4,
  InjectionTerm = 5,
  LambdaTerm = 6,
  LetTerm = 7,
  MatchTerm = 8,
  ProjectionTerm = 9,
  ReturnTerm = 10,
  SetTerm = 11,
  StructTerm = 12,
  TupleTerm = 13,
  TypeAbsTerm = 14,
  VarTerm = 15
};

enum class IrDeclCase {
  TermDecl = 0,
  AliasDecl = 1,
  NameDecl = 2,
  TraitDecl = 3
};

enum class ImplCase {
  TraitImpl = 0,
  InherentImpl = 1
};

enum class IrUnitCase {
  BaseUnit = 0,
  ImplUnit = 1
};

enum class IrKindCase {
  TypeKind = 0,
  ArrowKind = 1
};

inline std::string json_escape(std::string_view s) {
  std::string out;
  for (char c : s) {
    switch (c) {
      case '"': out += "\\\""; break;
      case '\\': out += "\\\\"; break;
      case '\b': out += "\\b"; break;
      case '\f': out += "\\f"; break;
      case '\n': out += "\\n"; break;
      case '\r': out += "\\r"; break;
      case '\t': out += "\\t"; break;
      default:
        if (static_cast<unsigned char>(c) < 32) {
          char buf[8];
          snprintf(buf, sizeof(buf), "\\u%04x", static_cast<unsigned char>(c));
          out += buf;
        } else {
          out += c;
        }
        break;
    }
  }
  return out;
}

struct Pos {
  std::string filename;
  int64_t begin_line_num = 0;
  int64_t end_line_num = 0;

  std::string to_string(bool with_comment = false) const {
    if (filename.empty() && begin_line_num == 0) return "";
    if (with_comment) {
      if (begin_line_num == end_line_num) {
        return "/* " + filename + ":" + std::to_string(begin_line_num) + " */";
      }
      return "/* " + filename + ":" + std::to_string(begin_line_num) + "-" + std::to_string(end_line_num) + " */";
    }
    if (begin_line_num == end_line_num) {
      return "in \"" + filename + "\" in line " + std::to_string(begin_line_num);
    }
    return "in \"" + filename + "\" in lines " + std::to_string(begin_line_num) + "-" + std::to_string(end_line_num);
  }

  std::string to_json() const {
    return "{\"Filename\":\"" + json_escape(filename) + "\",\"BeginLineNum\":" +
           std::to_string(begin_line_num) + ",\"EndLineNum\":" + std::to_string(end_line_num) + "}";
  }
};

inline Pos new_line_pos(std::string filename, int64_t line) {
  return Pos{std::move(filename), line, line};
}

inline Pos new_range_pos(std::string filename, int64_t begin_line, int64_t end_line) {
  return Pos{std::move(filename), begin_line, end_line};
}

struct Filename {
  std::string value;
  Pos pos;

  std::string to_string(bool with_pos = false) const {
    if (with_pos && !pos.filename.empty()) {
      return pos.to_string(true) + " \"" + value + "\"";
    }
    return "\"" + value + "\"";
  }

  std::string to_json() const {
    return "{\"Value\":\"" + json_escape(value) + "\",\"Pos\":" + pos.to_json() + "}";
  }
};

inline Filename new_filename(std::string value, Pos pos) {
  return Filename{std::move(value), std::move(pos)};
}

struct ModuleID {
  std::string name;
  Pos pos;

  std::string to_string(bool with_pos = false) const {
    if (with_pos && !pos.filename.empty()) {
      return pos.to_string(true) + name;
    }
    return name;
  }

  std::string to_json() const {
    return "{\"Name\":\"" + json_escape(name) + "\",\"Pos\":" + pos.to_json() + "}";
  }
};

inline ModuleID new_module_id(std::string name, Pos pos) {
  return ModuleID{std::move(name), std::move(pos)};
}

struct IrKind {
  IrKindCase case_val = IrKindCase::TypeKind;
  std::shared_ptr<IrKind> left;
  std::shared_ptr<IrKind> right;

  bool is_type_kind() const { return case_val == IrKindCase::TypeKind; }
  bool is_arrow_kind() const { return case_val == IrKindCase::ArrowKind; }

  std::string to_string() const {
    if (is_type_kind()) return "*";
    std::string l = left ? left->to_string() : "*";
    std::string r = right ? right->to_string() : "*";
    if (left && left->is_arrow_kind()) {
      l = "(" + l + ")";
    }
    return l + " -> " + r;
  }

  std::string to_json() const {
    if (is_type_kind()) {
      return "{\"Case\":0}";
    }
    return "{\"Case\":1,\"Left\":" + (left ? left->to_json() : "null") +
           ",\"Right\":" + (right ? right->to_json() : "null") + "}";
  }
};

inline IrKind new_type_kind() {
  return IrKind{IrKindCase::TypeKind, nullptr, nullptr};
}

inline IrKind new_arrow_kind(IrKind left, IrKind right) {
  return IrKind{
      IrKindCase::ArrowKind,
      std::make_shared<IrKind>(std::move(left)),
      std::make_shared<IrKind>(std::move(right)),
  };
}

// Forward declarations
struct IrType;
struct IrTerm;
struct IrDecl;
struct IrFunction;
struct IrTraitImpl;
struct IrUnit;

struct TypeParam {
  std::string var;
  IrKind kind;
  std::vector<IrType> bounds;

  std::string to_string() const;
  std::string to_json() const;
};

struct FloatLit {
  int64_t integer = 0;
  int64_t decimal = 0;
};

enum class IrLiteralCase {
  IntLiteral = 0,
  FloatLiteral = 1,
  RuneLiteral = 2,
  StrLiteral = 3,
};

struct IrLiteral {
  IrLiteralCase case_val = IrLiteralCase::IntLiteral;
  std::optional<int64_t> int_val;
  std::optional<FloatLit> float_val;
  std::optional<std::string> rune_val;
  std::optional<std::string> str_val;
  Pos pos;

  bool is(IrLiteralCase c) const { return case_val == c; }
  bool is_int() const { return int_val.has_value(); }
  bool is_float() const { return float_val.has_value(); }
  bool is_rune() const { return rune_val.has_value(); }
  bool is_str() const { return str_val.has_value(); }

  std::string to_string(bool with_pos = false) const {
    std::string s;
    if (is_int()) {
      s = std::to_string(int_val.value_or(0));
    } else if (is_float()) {
      s = std::to_string(float_val->integer) + "." + std::to_string(float_val->decimal);
    } else if (is_rune()) {
      s = "'" + rune_val.value_or("") + "'";
    } else if (is_str()) {
      s = "\"" + str_val.value_or("") + "\"";
    }
    if (with_pos && !pos.filename.empty()) {
      return pos.to_string(true) + s;
    }
    return s;
  }

  std::string to_json() const {
    std::string s = "{\"Case\":" + std::to_string(static_cast<int>(case_val));
    if (is_int()) {
      s += ",\"Int\":" + std::to_string(int_val.value_or(0));
    } else if (is_float()) {
      s += ",\"Float\":{\"Integer\":" + std::to_string(float_val->integer) +
           ",\"Decimal\":" + std::to_string(float_val->decimal) + "}";
    } else if (is_rune()) {
      s += ",\"Rune\":\"" + json_escape(rune_val.value_or("")) + "\"";
    } else if (is_str()) {
      s += ",\"Str\":\"" + json_escape(str_val.value_or("")) + "\"";
    }
    s += ",\"Pos\":" + pos.to_json() + "}";
    return s;
  }
};

inline IrLiteral new_int_literal(Pos pos, int64_t val) {
  IrLiteral lit;
  lit.case_val = IrLiteralCase::IntLiteral;
  lit.int_val = val;
  lit.pos = std::move(pos);
  return lit;
}

inline IrLiteral new_float_literal(Pos pos, int64_t integer, int64_t decimal) {
  IrLiteral lit;
  lit.case_val = IrLiteralCase::FloatLiteral;
  lit.float_val = FloatLit{integer, decimal};
  lit.pos = std::move(pos);
  return lit;
}

inline IrLiteral new_rune_literal(Pos pos, std::string val) {
  IrLiteral lit;
  lit.case_val = IrLiteralCase::RuneLiteral;
  lit.rune_val = std::move(val);
  lit.pos = std::move(pos);
  return lit;
}

inline IrLiteral new_str_literal(Pos pos, std::string val) {
  IrLiteral lit;
  lit.case_val = IrLiteralCase::StrLiteral;
  lit.str_val = std::move(val);
  lit.pos = std::move(pos);
  return lit;
}

template <typename T, typename SepFunc, typename ElemFunc>
inline void Interleave(const std::vector<T>& elems, SepFunc&& sep_func, ElemFunc&& elem_func) {
  for (size_t i = 0; i < elems.size(); ++i) {
    if (i > 0) {
      sep_func();
    }
    elem_func(static_cast<int>(i), elems[i]);
  }
}

inline bool is_operator(const std::string& id) {
  static const std::vector<std::string> ops = {
      "+", "-", "*", "/", "%", "==", "!=", "<", "<=", ">", ">=",
      "&&", "||", "!", "&", "|", "^", "~", "<<", ">>"
  };
  return std::find(ops.begin(), ops.end(), id) != ops.end();
}

} // namespace ir
