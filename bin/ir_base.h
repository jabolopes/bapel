#pragma once

#include "bapel/compiler/ir.h"
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

using ::Pos;
using ::Filename;
using ::ModuleID;

inline std::string pos_to_string(const Pos& p, bool with_comment = false) {
  if (p.filename.empty() && p.begin_line == 0) return "";
  if (with_comment) {
    if (p.begin_line == p.end_line) {
      return "/* " + p.filename + ":" + std::to_string(p.begin_line) + " */";
    }
    return "/* " + p.filename + ":" + std::to_string(p.begin_line) + "-" + std::to_string(p.end_line) + " */";
  }
  if (p.begin_line == p.end_line) {
    return "in \"" + p.filename + "\" in line " + std::to_string(p.begin_line);
  }
  return "in \"" + p.filename + "\" in lines " + std::to_string(p.begin_line) + "-" + std::to_string(p.end_line);
}

inline Pos new_line_pos(std::string filename, int64_t line) {
  return Pos{std::move(filename), line, line};
}

inline Pos new_range_pos(std::string filename, int64_t begin_line, int64_t end_line) {
  return Pos{std::move(filename), begin_line, end_line};
}

inline std::string filename_to_string(const Filename& f, bool with_pos = false) {
  if (with_pos && !f.pos.filename.empty()) {
    return pos_to_string(f.pos, true) + " \"" + f.value + "\"";
  }
  return "\"" + f.value + "\"";
}

inline Filename new_filename(std::string value, Pos pos = Pos{}) {
  return Filename{std::move(value), std::move(pos)};
}

inline std::string module_id_to_string(const ModuleID& m, bool with_pos = false) {
  if (with_pos && !m.pos.filename.empty()) {
    return pos_to_string(m.pos, true) + m.name;
  }
  return m.name;
}

inline std::string module_id_to_filename(const ModuleID& m) {
  std::string s = m.name;
  std::replace(s.begin(), s.end(), '.', '/');
  return s;
}

inline ModuleID new_module_id(std::string name, Pos pos = Pos{}) {
  return ModuleID{std::move(name), std::move(pos)};
}

using ::IrKind;

inline bool is_type_kind(const IrKind& k) {
  return k.is(IrKind::Tag::type_kind);
}

inline bool is_arrow_kind(const IrKind& k) {
  return k.is(IrKind::Tag::arrow_kind);
}

inline std::string kind_to_string(const IrKind& k) {
  if (is_type_kind(k)) return "∗";
  const auto& ak = k.as_arrow_kind();
  std::string l = ak.left ? kind_to_string(*ak.left) : "∗";
  std::string r = ak.right ? kind_to_string(*ak.right) : "∗";
  if (ak.left && is_arrow_kind(*ak.left)) {
    l = "(" + l + ")";
  }
  return l + " -> " + r;
}

inline IrKind new_type_kind() {
  return IrKind::mk_type_kind();
}

inline IrKind new_arrow_kind(IrKind left, IrKind right) {
  return IrKind::mk_arrow_kind({
      std::make_shared<IrKind>(std::move(left)),
      std::make_shared<IrKind>(std::move(right)),
  });
}

inline bool equals_kind(const IrKind& k1, const IrKind& k2) {
  if (k1.tag() != k2.tag()) return false;
  if (is_type_kind(k1)) return true;
  if (is_arrow_kind(k1)) {
    const auto& a1 = k1.as_arrow_kind();
    const auto& a2 = k2.as_arrow_kind();
    if (!a1.left || !a2.left || !a1.right || !a2.right) {
      return a1.left == a2.left && a1.right == a2.right;
    }
    return equals_kind(*a1.left, *a2.left) && equals_kind(*a1.right, *a2.right);
  }
  return false;
}

} // namespace ir

inline bool operator==(const ::IrKind& k1, const ::IrKind& k2) {
  return ir::equals_kind(k1, k2);
}

namespace ir {

inline bool operator!=(const IrKind& k1, const IrKind& k2) {
  return !(k1 == k2);
}

// Forward declarations
using ::IrType;
using TypeParam = ::TypeParam<::IrType>;
using ::IrTerm;
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

using ::FloatLit;
using ::IrLiteral;

inline bool is_int_literal(const IrLiteral& lit) { return lit.is(IrLiteral::Tag::int_lit); }
inline bool is_float_literal(const IrLiteral& lit) { return lit.is(IrLiteral::Tag::float_lit); }
inline bool is_rune_literal(const IrLiteral& lit) { return lit.is(IrLiteral::Tag::rune_lit); }
inline bool is_str_literal(const IrLiteral& lit) { return lit.is(IrLiteral::Tag::str_lit); }

inline std::string literal_to_string(const IrLiteral& lit) {
  switch (lit.tag()) {
    case IrLiteral::Tag::int_lit:
      return std::to_string(lit.as_int_lit());
    case IrLiteral::Tag::float_lit:
      return std::to_string(lit.as_float_lit().integer) + "." + std::to_string(lit.as_float_lit().decimal);
    case IrLiteral::Tag::rune_lit:
      return "'" + lit.as_rune_lit() + "'";
    case IrLiteral::Tag::str_lit:
      return "\"" + lit.as_str_lit() + "\"";
  }
  return "";
}

inline IrLiteral new_int_literal(Pos /*pos*/, int64_t val) {
  return IrLiteral::mk_int_lit(val);
}

inline IrLiteral new_float_literal(Pos /*pos*/, int64_t integer, int64_t decimal) {
  return IrLiteral::mk_float_lit(FloatLit{integer, decimal});
}

inline IrLiteral new_rune_literal(Pos /*pos*/, std::string val) {
  return IrLiteral::mk_rune_lit(std::move(val));
}

inline IrLiteral new_str_literal(Pos /*pos*/, std::string val) {
  return IrLiteral::mk_str_lit(std::move(val));
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
