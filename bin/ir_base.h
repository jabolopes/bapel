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

struct Pos {
  std::string filename;
  int64_t begin_line_num = 0;
  int64_t end_line_num = 0;
};

struct Filename {
  std::string value;
  Pos pos;
};

struct ModuleID {
  std::string name;
  Pos pos;
};

struct IrKind {
  IrKindCase case_val = IrKindCase::TypeKind;
  std::shared_ptr<IrKind> left;
  std::shared_ptr<IrKind> right;

  bool is_type_kind() const { return case_val == IrKindCase::TypeKind; }
  bool is_arrow_kind() const { return case_val == IrKindCase::ArrowKind; }
};

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
};

struct FloatLit {
  int64_t integer = 0;
  int64_t decimal = 0;
};

struct IrLiteral {
  std::optional<int64_t> int_val;
  std::optional<FloatLit> float_val;
  std::optional<std::string> rune_val;
  std::optional<std::string> str_val;

  bool is_int() const { return int_val.has_value(); }
  bool is_float() const { return float_val.has_value(); }
  bool is_rune() const { return rune_val.has_value(); }
  bool is_str() const { return str_val.has_value(); }
};

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
