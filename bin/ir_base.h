#pragma once

#include "bapel/core_variant.h"

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

enum class ImplCase {
  TraitImpl = 0,
  InherentImpl = 1
};

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

  std::string to_filename() const {
    std::string s = name;
    std::replace(s.begin(), s.end(), '.', '/');
    return s;
  }
};

inline ModuleID new_module_id(std::string name, Pos pos) {
  return ModuleID{std::move(name), std::move(pos)};
}

struct IrKind;

struct ArrowKindData {
  std::shared_ptr<IrKind> left;
  std::shared_ptr<IrKind> right;
};

// TODO: Move member helpers (is_type_kind, is_arrow_kind, to_string) to free
// functions or inherents::IrKind and define IrKind directly with BAPEL_VARIANT.
BAPEL_VARIANT(IrKindBase,
  (type_kind, (std::monostate)),
  (arrow_kind, (ArrowKindData))
);

struct IrKind : IrKindBase {
  using IrKindBase::IrKindBase;
  IrKind() : IrKindBase(IrKindBase::mk_type_kind()) {}
  IrKind(IrKindBase b) : IrKindBase(std::move(b)) {}

  bool is_type_kind() const { return is(Tag::type_kind); }
  bool is_arrow_kind() const { return is(Tag::arrow_kind); }

  std::string to_string() const {
    if (is_type_kind()) return "∗";
    const auto& ak = as_arrow_kind();
    std::string l = ak.left ? ak.left->to_string() : "∗";
    std::string r = ak.right ? ak.right->to_string() : "∗";
    if (ak.left && ak.left->is_arrow_kind()) {
      l = "(" + l + ")";
    }
    return l + " -> " + r;
  }
};

inline IrKind new_type_kind() {
  return IrKind(IrKindBase::mk_type_kind());
}

inline IrKind new_arrow_kind(IrKind left, IrKind right) {
  return IrKind(IrKindBase::mk_arrow_kind(ArrowKindData{
      std::make_shared<IrKind>(std::move(left)),
      std::make_shared<IrKind>(std::move(right)),
  }));
}

inline bool equals_kind(const IrKind& k1, const IrKind& k2) {
  if (k1.tag() != k2.tag()) return false;
  if (k1.is_type_kind()) return true;
  if (k1.is_arrow_kind()) {
    const auto& a1 = k1.as_arrow_kind();
    const auto& a2 = k2.as_arrow_kind();
    if (!a1.left || !a2.left || !a1.right || !a2.right) {
      return a1.left == a2.left && a1.right == a2.right;
    }
    return equals_kind(*a1.left, *a2.left) && equals_kind(*a1.right, *a2.right);
  }
  return false;
}

inline bool operator==(const IrKind& k1, const IrKind& k2) {
  return equals_kind(k1, k2);
}

inline bool operator!=(const IrKind& k1, const IrKind& k2) {
  return !(k1 == k2);
}

// Forward declarations
struct IrType;
struct IrTerm;
struct IrDecl;
struct IrFunction;
struct IrTraitImpl;
struct IrFile;
struct IrModule;
struct IrTypedModule;

enum class IrFileCase {
  Base = 0,
  Impl = 1,
  Ffi = 2,
};

struct TypeParam {
  std::string var;
  IrKind kind;
  std::vector<IrType> bounds;

  TypeParam();
  TypeParam(std::string v, IrKind k);
  TypeParam(std::string v, IrKind k, std::vector<IrType> b);
  TypeParam(const TypeParam&);
  TypeParam(TypeParam&&) noexcept;
  TypeParam& operator=(const TypeParam&);
  TypeParam& operator=(TypeParam&&) noexcept;
  ~TypeParam();

  std::string to_string() const;
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

BAPEL_VARIANT(IrLiteralData,
  (int_lit, (int64_t)),
  (float_lit, (FloatLit)),
  (rune_lit, (std::string)),
  (str_lit, (std::string))
);

struct IrLiteral {
  IrLiteralData data{IrLiteralData::mk_int_lit(0)};
  Pos pos;

  bool is_int() const { return data.is(IrLiteralData::Tag::int_lit); }
  bool is_float() const { return data.is(IrLiteralData::Tag::float_lit); }
  bool is_rune() const { return data.is(IrLiteralData::Tag::rune_lit); }
  bool is_str() const { return data.is(IrLiteralData::Tag::str_lit); }

  std::string to_string(bool with_pos = false) const {
    std::string s;
    switch (data.tag()) {
      case IrLiteralData::Tag::int_lit:
        s = std::to_string(data.as_int_lit());
        break;
      case IrLiteralData::Tag::float_lit:
        s = std::to_string(data.as_float_lit().integer) + "." + std::to_string(data.as_float_lit().decimal);
        break;
      case IrLiteralData::Tag::rune_lit:
        s = "'" + data.as_rune_lit() + "'";
        break;
      case IrLiteralData::Tag::str_lit:
        s = "\"" + data.as_str_lit() + "\"";
        break;
    }
    if (with_pos && !pos.filename.empty()) {
      return pos.to_string(true) + s;
    }
    return s;
  }
};

inline IrLiteral new_int_literal(Pos pos, int64_t val) {
  return IrLiteral{IrLiteralData::mk_int_lit(val), std::move(pos)};
}

inline IrLiteral new_float_literal(Pos pos, int64_t integer, int64_t decimal) {
  return IrLiteral{IrLiteralData::mk_float_lit(FloatLit{integer, decimal}), std::move(pos)};
}

inline IrLiteral new_rune_literal(Pos pos, std::string val) {
  return IrLiteral{IrLiteralData::mk_rune_lit(std::move(val)), std::move(pos)};
}

inline IrLiteral new_str_literal(Pos pos, std::string val) {
  return IrLiteral{IrLiteralData::mk_str_lit(std::move(val)), std::move(pos)};
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
