#pragma once

#include "ir_unit.h"
#include <cctype>
#include <map>
#include <memory>
#include <sstream>
#include <string>
#include <variant>
#include <vector>

namespace ir {

class JsonValue {
public:
  enum class Type { Null, Bool, Number, String, Array, Object };

  using ArrayType = std::vector<JsonValue>;
  using ObjectType = std::map<std::string, JsonValue>;

  Type type = Type::Null;
  bool bool_val = false;
  double num_val = 0.0;
  std::string str_val;
  ArrayType arr_val;
  ObjectType obj_val;

  bool is_null() const { return type == Type::Null; }
  bool is_bool() const { return type == Type::Bool; }
  bool is_number() const { return type == Type::Number; }
  bool is_string() const { return type == Type::String; }
  bool is_array() const { return type == Type::Array; }
  bool is_object() const { return type == Type::Object; }

  int64_t as_int() const { return static_cast<int64_t>(num_val); }
  double as_double() const { return num_val; }
  bool as_bool() const { return bool_val; }
  const std::string& as_string() const { return str_val; }
  const ArrayType& as_array() const { return arr_val; }
  const ObjectType& as_object() const { return obj_val; }

  bool has_field(const std::string& key) const {
    if (!is_object()) return false;
    auto it = obj_val.find(key);
    return it != obj_val.end() && !it->second.is_null();
  }

  const JsonValue& get(const std::string& key) const {
    static const JsonValue null_val;
    if (!is_object()) return null_val;
    auto it = obj_val.find(key);
    return it != obj_val.end() ? it->second : null_val;
  }
};

class JsonParser {
public:
  static JsonValue parse(std::string_view json) {
    JsonParser p(json);
    return p.parse_value();
  }

private:
  std::string_view src_;
  size_t pos_ = 0;

  explicit JsonParser(std::string_view json) : src_(json) {}

  void skip_whitespace() {
    while (pos_ < src_.size() && (std::isspace(static_cast<unsigned char>(src_[pos_])) != 0)) {
      pos_++;
    }
  }

  char peek() {
    skip_whitespace();
    return pos_ < src_.size() ? src_[pos_] : '\0';
  }

  char get() {
    skip_whitespace();
    return pos_ < src_.size() ? src_[pos_++] : '\0';
  }

  JsonValue parse_value() {
    char c = peek();
    if (c == 'n') return parse_null();
    if (c == 't' || c == 'f') return parse_bool();
    if (c == '"') return parse_string();
    if (c == '[') return parse_array();
    if (c == '{') return parse_object();
    if (c == '-' || (c >= '0' && c <= '9')) return parse_number();
    return JsonValue{};
  }

  JsonValue parse_null() {
    if (src_.substr(pos_, 4) == "null") {
      pos_ += 4;
    }
    return JsonValue{};
  }

  JsonValue parse_bool() {
    JsonValue v;
    v.type = JsonValue::Type::Bool;
    if (src_.substr(pos_, 4) == "true") {
      pos_ += 4;
      v.bool_val = true;
    } else if (src_.substr(pos_, 5) == "false") {
      pos_ += 5;
      v.bool_val = false;
    }
    return v;
  }

  JsonValue parse_number() {
    skip_whitespace();
    size_t start = pos_;
    if (pos_ < src_.size() && src_[pos_] == '-') pos_++;
    while (pos_ < src_.size() && (std::isdigit(static_cast<unsigned char>(src_[pos_])) != 0)) pos_++;
    if (pos_ < src_.size() && src_[pos_] == '.') {
      pos_++;
      while (pos_ < src_.size() && (std::isdigit(static_cast<unsigned char>(src_[pos_])) != 0)) pos_++;
    }
    if (pos_ < src_.size() && (src_[pos_] == 'e' || src_[pos_] == 'E')) {
      pos_++;
      if (pos_ < src_.size() && (src_[pos_] == '+' || src_[pos_] == '-')) pos_++;
      while (pos_ < src_.size() && (std::isdigit(static_cast<unsigned char>(src_[pos_])) != 0)) pos_++;
    }

    std::string num_str(src_.substr(start, pos_ - start));
    JsonValue v;
    v.type = JsonValue::Type::Number;
    try {
      v.num_val = std::stod(num_str);
    } catch (...) {
      v.num_val = 0.0;
    }
    return v;
  }

  JsonValue parse_string() {
    JsonValue v;
    v.type = JsonValue::Type::String;
    if (get() != '"') return v;

    std::string s;
    while (pos_ < src_.size()) {
      char c = src_[pos_++];
      if (c == '"') {
        v.str_val = s;
        return v;
      }
      if (c == '\\' && pos_ < src_.size()) {
        char esc = src_[pos_++];
        switch (esc) {
          case '"': s += '"'; break;
          case '\\': s += '\\'; break;
          case '/': s += '/'; break;
          case 'b': s += '\b'; break;
          case 'f': s += '\f'; break;
          case 'n': s += '\n'; break;
          case 'r': s += '\r'; break;
          case 't': s += '\t'; break;
          case 'u': {
            if (pos_ + 4 <= src_.size()) {
              std::string hex_str(src_.substr(pos_, 4));
              pos_ += 4;
              try {
                int code = std::stoi(hex_str, nullptr, 16);
                if (code <= 127) {
                  s += static_cast<char>(code);
                } else if (code <= 0x7FF) {
                  s += static_cast<char>(0xC0 | (code >> 6));
                  s += static_cast<char>(0x80 | (code & 0x3F));
                } else {
                  s += static_cast<char>(0xE0 | (code >> 12));
                  s += static_cast<char>(0x80 | ((code >> 6) & 0x3F));
                  s += static_cast<char>(0x80 | (code & 0x3F));
                }
              } catch (...) {
                s += '?';
              }
            }
            break;
          }
          default: s += esc; break;
        }
      } else {
        s += c;
      }
    }
    v.str_val = s;
    return v;
  }

  JsonValue parse_array() {
    JsonValue v;
    v.type = JsonValue::Type::Array;
    get(); // consume '['

    if (peek() == ']') {
      get();
      return v;
    }

    while (true) {
      v.arr_val.push_back(parse_value());
      char c = peek();
      if (c == ',') {
        get();
      } else if (c == ']') {
        get();
        break;
      } else {
        break;
      }
    }
    return v;
  }

  JsonValue parse_object() {
    JsonValue v;
    v.type = JsonValue::Type::Object;
    get(); // consume '{'

    if (peek() == '}') {
      get();
      return v;
    }

    while (true) {
      if (peek() != '"') break;
      JsonValue key_val = parse_string();
      if (peek() != ':') break;
      get(); // consume ':'
      JsonValue val = parse_value();
      v.obj_val[key_val.str_val] = std::move(val);

      char c = peek();
      if (c == ',') {
        get();
      } else if (c == '}') {
        get();
        break;
      } else {
        break;
      }
    }
    return v;
  }
};

// Deserialization functions
inline Pos deserialize_pos(const JsonValue& j) {
  Pos p;
  if (!j.is_object()) return p;
  p.filename = j.get("Filename").as_string();
  p.begin_line_num = j.get("BeginLineNum").as_int();
  p.end_line_num = j.get("EndLineNum").as_int();
  return p;
}

inline Filename deserialize_filename(const JsonValue& j) {
  Filename f;
  if (!j.is_object()) return f;
  f.value = j.get("Value").as_string();
  f.pos = deserialize_pos(j.get("Pos"));
  return f;
}

inline ModuleID deserialize_module_id(const JsonValue& j) {
  ModuleID m;
  if (!j.is_object()) return m;
  m.name = j.get("Name").as_string();
  m.pos = deserialize_pos(j.get("Pos"));
  return m;
}

inline IrKind deserialize_kind(const JsonValue& j) {
  IrKind k;
  if (!j.is_object()) return k;
  k.case_val = static_cast<IrKindCase>(j.get("Case").as_int());
  if (j.has_field("Arrow")) {
    const auto& a = j.get("Arrow");
    k.left = std::make_shared<IrKind>(deserialize_kind(a.has_field("Arg") ? a.get("Arg") : a.get("Left")));
    k.right = std::make_shared<IrKind>(deserialize_kind(a.has_field("Ret") ? a.get("Ret") : a.get("Right")));
  }
  return k;
}

inline IrType deserialize_type(const JsonValue& j);

inline TypeParam deserialize_type_param(const JsonValue& j) {
  TypeParam tp;
  if (!j.is_object()) return tp;
  tp.var = j.get("Var").as_string();
  tp.kind = deserialize_kind(j.get("Kind"));
  if (j.has_field("Bounds") && j.get("Bounds").is_array()) {
    for (const auto& b : j.get("Bounds").as_array()) {
      tp.bounds.push_back(deserialize_type(b));
    }
  }
  return tp;
}

inline IrType deserialize_type(const JsonValue& j) {
  IrType t;
  if (!j.is_object()) return t;
  t.case_val = static_cast<IrTypeCase>(j.get("Case").as_int());
  t.pos = deserialize_pos(j.get("Pos"));

  switch (t.case_val) {
    case IrTypeCase::AppType: {
      if (j.has_field("App")) {
        const auto& app_j = j.get("App");
        t.app = std::make_shared<AppTypeData>();
        t.app->fun = std::make_shared<IrType>(deserialize_type(app_j.get("Fun")));
        t.app->arg = std::make_shared<IrType>(deserialize_type(app_j.get("Arg")));
      }
      break;
    }
    case IrTypeCase::ArrayType: {
      if (j.has_field("Array")) {
        const auto& arr_j = j.get("Array");
        t.array = std::make_shared<ArrayTypeData>();
        t.array->elem_type = std::make_shared<IrType>(deserialize_type(arr_j.get("ElemType")));
        t.array->size = arr_j.get("Size").as_int();
      }
      break;
    }
    case IrTypeCase::ExistVarType:
      t.exist_var = j.get("ExistVar").as_int();
      break;
    case IrTypeCase::ForallType: {
      if (j.has_field("Forall")) {
        const auto& f_j = j.get("Forall");
        t.forall = std::make_shared<ForallTypeData>();
        t.forall->type_param = deserialize_type_param(f_j.has_field("TypeParam") ? f_j.get("TypeParam") : f_j);
        t.forall->type = std::make_shared<IrType>(deserialize_type(f_j.get("Type")));
      }
      break;
    }
    case IrTypeCase::FunType: {
      if (j.has_field("Fun")) {
        const auto& fun_j = j.get("Fun");
        t.fun = std::make_shared<FunctionTypeData>();
        t.fun->arg = std::make_shared<IrType>(deserialize_type(fun_j.get("Arg")));
        t.fun->ret = std::make_shared<IrType>(deserialize_type(fun_j.get("Ret")));
      }
      break;
    }
    case IrTypeCase::LambdaType: {
      if (j.has_field("Lambda")) {
        const auto& l_j = j.get("Lambda");
        t.lambda = std::make_shared<LambdaTypeData>();
        t.lambda->var = l_j.get("Var").as_string();
        t.lambda->kind = deserialize_kind(l_j.get("Kind"));
        t.lambda->type = std::make_shared<IrType>(deserialize_type(l_j.get("Type")));
      }
      break;
    }
    case IrTypeCase::NameType:
      t.name = j.get("Name").as_string();
      break;
    case IrTypeCase::StructType: {
      if (j.has_field("Struct") && j.get("Struct").has_field("Fields")) {
        t.struct_data = std::make_shared<StructTypeData>();
        for (const auto& f_j : j.get("Struct").get("Fields").as_array()) {
          StructField f;
          f.id = f_j.get("ID").as_string();
          f.type = std::make_shared<IrType>(deserialize_type(f_j.get("Type")));
          t.struct_data->fields.push_back(std::move(f));
        }
      }
      break;
    }
    case IrTypeCase::TupleType: {
      if (j.has_field("Tuple") && j.get("Tuple").has_field("Elems")) {
        t.tuple_data = std::make_shared<TupleTypeData>();
        for (const auto& e_j : j.get("Tuple").get("Elems").as_array()) {
          t.tuple_data->elems.push_back(deserialize_type(e_j));
        }
      }
      break;
    }
    case IrTypeCase::VariantType: {
      if (j.has_field("Variant") && j.get("Variant").has_field("Tags")) {
        t.variant_data = std::make_shared<VariantTypeData>();
        for (const auto& tag_j : j.get("Variant").get("Tags").as_array()) {
          VariantTag tag;
          tag.id = tag_j.get("ID").as_string();
          tag.type = std::make_shared<IrType>(deserialize_type(tag_j.get("Type")));
          t.variant_data->tags.push_back(std::move(tag));
        }
      }
      break;
    }
    case IrTypeCase::VarType:
      t.var = j.get("Var").as_string();
      break;
  }
  return t;
}

inline FunctionArg deserialize_function_arg(const JsonValue& j) {
  FunctionArg arg;
  if (!j.is_object()) return arg;
  arg.id = j.get("ID").as_string();
  arg.type = deserialize_type(j.get("Type"));
  return arg;
}

inline IrTerm deserialize_term(const JsonValue& j);

inline MatchArm deserialize_match_arm(const JsonValue& j) {
  MatchArm arm;
  if (!j.is_object()) return arm;
  arm.tag = j.get("Tag").as_string();
  arm.arg = j.get("Arg").as_string();
  arm.body = std::make_shared<IrTerm>(deserialize_term(j.get("Body")));
  if (j.has_field("Index")) {
    arm.index = static_cast<int>(j.get("Index").as_int());
  }
  return arm;
}

inline LabelValue deserialize_label_value(const JsonValue& j) {
  LabelValue lv;
  if (!j.is_object()) return lv;
  lv.label = j.get("Label").as_string();
  lv.value = std::make_shared<IrTerm>(deserialize_term(j.get("Value")));
  return lv;
}

inline IrTerm deserialize_term(const JsonValue& j) {
  IrTerm t;
  if (!j.is_object()) return t;
  t.case_val = static_cast<IrTermCase>(j.get("Case").as_int());
  t.pos = deserialize_pos(j.get("Pos"));
  if (j.has_field("Type")) {
    t.type = deserialize_type(j.get("Type"));
  }

  switch (t.case_val) {
    case IrTermCase::AppTermTerm: {
      if (j.has_field("AppTerm")) {
        const auto& at = j.get("AppTerm");
        t.app_term = std::make_shared<AppTermData>();
        t.app_term->fun = std::make_shared<IrTerm>(deserialize_term(at.get("Fun")));
        t.app_term->arg = std::make_shared<IrTerm>(deserialize_term(at.get("Arg")));
      }
      break;
    }
    case IrTermCase::AppTypeTerm: {
      if (j.has_field("AppType")) {
        const auto& at = j.get("AppType");
        t.app_type = std::make_shared<AppTypeTermData>();
        t.app_type->fun = std::make_shared<IrTerm>(deserialize_term(at.get("Fun")));
        t.app_type->arg = deserialize_type(at.get("Arg"));
      }
      break;
    }
    case IrTermCase::AssignTerm: {
      if (j.has_field("Assign")) {
        const auto& a = j.get("Assign");
        t.assign = std::make_shared<AssignTermData>();
        t.assign->ret = std::make_shared<IrTerm>(deserialize_term(a.get("Ret")));
        t.assign->arg = std::make_shared<IrTerm>(deserialize_term(a.get("Arg")));
      }
      break;
    }
    case IrTermCase::BlockTerm: {
      if (j.has_field("Block") && j.get("Block").has_field("Terms")) {
        t.block = std::make_shared<BlockTermData>();
        for (const auto& term_j : j.get("Block").get("Terms").as_array()) {
          t.block->terms.push_back(deserialize_term(term_j));
        }
      }
      break;
    }
    case IrTermCase::ConstTerm: {
      if (j.has_field("Const")) {
        const auto& c = j.get("Const");
        t.const_data = std::make_shared<ConstTermData>();
        if (c.has_field("Int")) {
          t.const_data->literal.int_val = c.get("Int").as_int();
        } else if (c.has_field("Float")) {
          const auto& f = c.get("Float");
          t.const_data->literal.float_val = FloatLit{f.get("Integer").as_int(), f.get("Decimal").as_int()};
        } else if (c.has_field("Rune")) {
          t.const_data->literal.rune_val = c.get("Rune").as_string();
        } else if (c.has_field("Str")) {
          t.const_data->literal.str_val = c.get("Str").as_string();
        }
      }
      break;
    }
    case IrTermCase::InjectionTerm: {
      if (j.has_field("Injection")) {
        const auto& inj = j.get("Injection");
        t.injection = std::make_shared<InjectionTermData>();
        t.injection->variant_type = deserialize_type(inj.get("VariantType"));
        t.injection->tag = inj.get("Tag").as_string();
        t.injection->value = std::make_shared<IrTerm>(deserialize_term(inj.get("Value")));
        if (inj.has_field("TagIndex")) {
          t.injection->tag_index = static_cast<int>(inj.get("TagIndex").as_int());
        }
      }
      break;
    }
    case IrTermCase::LambdaTerm: {
      if (j.has_field("Lambda")) {
        const auto& l = j.get("Lambda");
        t.lambda = std::make_shared<LambdaTermData>();
        t.lambda->arg = deserialize_function_arg(l.get("Arg"));
        t.lambda->body = std::make_shared<IrTerm>(deserialize_term(l.get("Body")));
      }
      break;
    }
    case IrTermCase::LetTerm: {
      if (j.has_field("Let")) {
        const auto& lt = j.get("Let");
        t.let_data = std::make_shared<LetTermData>();
        t.let_data->var = lt.get("Var").as_string();
        if (lt.has_field("VarType")) {
          t.let_data->var_type = deserialize_type(lt.get("VarType"));
        }
        t.let_data->value = std::make_shared<IrTerm>(deserialize_term(lt.get("Value")));
      }
      break;
    }
    case IrTermCase::MatchTerm: {
      if (j.has_field("Match")) {
        const auto& m = j.get("Match");
        t.match_data = std::make_shared<MatchTermData>();
        t.match_data->term = std::make_shared<IrTerm>(deserialize_term(m.get("Term")));
        if (m.has_field("Arms")) {
          for (const auto& arm_j : m.get("Arms").as_array()) {
            t.match_data->arms.push_back(deserialize_match_arm(arm_j));
          }
        }
      }
      break;
    }
    case IrTermCase::ProjectionTerm: {
      if (j.has_field("Projection")) {
        const auto& p = j.get("Projection");
        t.projection = std::make_shared<ProjectionTermData>();
        t.projection->term = std::make_shared<IrTerm>(deserialize_term(p.get("Term")));
        t.projection->label = p.get("Label").as_string();
        if (p.has_field("ReducedType")) {
          t.projection->reduced_type = deserialize_type(p.get("ReducedType"));
        }
      }
      break;
    }
    case IrTermCase::ReturnTerm: {
      if (j.has_field("Return")) {
        t.return_data = std::make_shared<ReturnTermData>();
        t.return_data->expr = std::make_shared<IrTerm>(deserialize_term(j.get("Return").get("Expr")));
      }
      break;
    }
    case IrTermCase::SetTerm: {
      if (j.has_field("Set")) {
        const auto& st = j.get("Set");
        t.set_data = std::make_shared<SetTermData>();
        t.set_data->term = std::make_shared<IrTerm>(deserialize_term(st.get("Term")));
        if (st.has_field("Values")) {
          for (const auto& lv_j : st.get("Values").as_array()) {
            t.set_data->values.push_back(deserialize_label_value(lv_j));
          }
        }
        if (st.has_field("ReducedType")) {
          t.set_data->reduced_type = deserialize_type(st.get("ReducedType"));
        }
      }
      break;
    }
    case IrTermCase::StructTerm: {
      if (j.has_field("Struct") && j.get("Struct").has_field("Values")) {
        t.struct_data = std::make_shared<StructTermData>();
        for (const auto& lv_j : j.get("Struct").get("Values").as_array()) {
          t.struct_data->values.push_back(deserialize_label_value(lv_j));
        }
      }
      break;
    }
    case IrTermCase::TupleTerm: {
      if (j.has_field("Tuple") && j.get("Tuple").has_field("Elems")) {
        t.tuple_data = std::make_shared<TupleTermData>();
        for (const auto& elem_j : j.get("Tuple").get("Elems").as_array()) {
          t.tuple_data->elems.push_back(deserialize_term(elem_j));
        }
      }
      break;
    }
    case IrTermCase::TypeAbsTerm: {
      if (j.has_field("TypeAbs")) {
        const auto& ta = j.get("TypeAbs");
        t.type_abs = std::make_shared<TypeAbsTermData>();
        t.type_abs->type_param = deserialize_type_param(ta.has_field("Arg") ? ta.get("Arg") : ta.get("TypeParam"));
        t.type_abs->body = std::make_shared<IrTerm>(deserialize_term(ta.get("Body")));
      }
      break;
    }
    case IrTermCase::VarTerm: {
      t.var_data = std::make_shared<VarTermData>();
      if (j.has_field("Var")) {
        t.var_data->id = j.get("Var").get("ID").as_string();
      }
      break;
    }
  }
  return t;
}

inline IrSignature deserialize_signature(const JsonValue& j) {
  IrSignature s;
  if (!j.is_object()) return s;
  s.id = j.get("ID").as_string();
  s.ret_type = deserialize_type(j.get("RetType"));
  if (j.has_field("Args")) {
    for (const auto& arg_j : j.get("Args").as_array()) {
      s.args.push_back(deserialize_function_arg(arg_j));
    }
  }
  return s;
}

inline IrDecl deserialize_decl(const JsonValue& j) {
  IrDecl d;
  if (!j.is_object()) return d;
  d.case_val = static_cast<IrDeclCase>(j.get("Case").as_int());
  d.export_flag = j.get("Export").as_bool();
  d.pos = deserialize_pos(j.get("Pos"));

  switch (d.case_val) {
    case IrDeclCase::TermDecl: {
      if (j.has_field("Term")) {
        d.term = std::make_shared<TermDeclData>();
        d.term->id = j.get("Term").get("ID").as_string();
        d.term->type = deserialize_type(j.get("Term").get("Type"));
      }
      break;
    }
    case IrDeclCase::AliasDecl: {
      if (j.has_field("Alias")) {
        const auto& al = j.get("Alias");
        d.alias = std::make_shared<AliasDeclData>();
        d.alias->id = al.get("ID").as_string();
        d.alias->kind = deserialize_kind(al.get("Kind"));
        d.alias->type = deserialize_type(al.get("Type"));
      }
      break;
    }
    case IrDeclCase::NameDecl: {
      if (j.has_field("Name")) {
        d.name = std::make_shared<NameDeclData>();
        d.name->id = j.get("Name").get("ID").as_string();
        d.name->kind = deserialize_kind(j.get("Name").get("Kind"));
      }
      break;
    }
    case IrDeclCase::TraitDecl: {
      if (j.has_field("Trait")) {
        const auto& tr = j.get("Trait");
        d.trait = std::make_shared<TraitDeclData>();
        d.trait->id = tr.get("ID").as_string();
        if (tr.has_field("TypeParams")) {
          for (const auto& tp_j : tr.get("TypeParams").as_array()) {
            d.trait->type_params.push_back(deserialize_type_param(tp_j));
          }
        }
        if (tr.has_field("Methods")) {
          for (const auto& m_j : tr.get("Methods").as_array()) {
            d.trait->methods.push_back(deserialize_signature(m_j));
          }
        }
      }
      break;
    }
  }
  return d;
}

inline IrFunction deserialize_function(const JsonValue& j) {
  IrFunction f;
  if (!j.is_object()) return f;
  f.export_flag = j.get("Export").as_bool();
  f.id = j.get("ID").as_string();
  f.ret_type = deserialize_type(j.get("RetType"));
  f.body = deserialize_term(j.get("Body"));
  f.pos = deserialize_pos(j.get("Pos"));

  if (j.has_field("TypeParams")) {
    for (const auto& tp_j : j.get("TypeParams").as_array()) {
      f.type_params.push_back(deserialize_type_param(tp_j));
    }
  }
  if (j.has_field("Args")) {
    for (const auto& arg_j : j.get("Args").as_array()) {
      f.args.push_back(deserialize_function_arg(arg_j));
    }
  }
  return f;
}

inline IrTraitImpl deserialize_trait_impl(const JsonValue& j) {
  IrTraitImpl impl;
  if (!j.is_object()) return impl;
  impl.case_val = static_cast<ImplCase>(j.get("Case").as_int());
  impl.trait_type = deserialize_type(j.get("TraitType"));
  impl.type_name = deserialize_type(j.get("TypeName"));
  impl.pos = deserialize_pos(j.get("Pos"));

  if (j.has_field("TypeParams")) {
    for (const auto& tp_j : j.get("TypeParams").as_array()) {
      impl.type_params.push_back(deserialize_type_param(tp_j));
    }
  }
  if (j.has_field("Methods")) {
    for (const auto& m_j : j.get("Methods").as_array()) {
      impl.methods.push_back(deserialize_function(m_j));
    }
  }
  return impl;
}

inline IrUnit parse_ir_unit_from_json(std::string_view json_str) {
  JsonValue j = JsonParser::parse(json_str);
  IrUnit unit;
  if (!j.is_object()) return unit;

  unit.case_val = static_cast<IrUnitCase>(j.get("Case").as_int());
  unit.module_id = deserialize_module_id(j.get("ModuleID"));
  unit.filename = deserialize_filename(j.get("Filename"));

  if (j.has_field("Imports")) {
    for (const auto& imp_j : j.get("Imports").as_array()) {
      IrImport imp;
      imp.module_id = deserialize_module_id(imp_j.get("ModuleID"));
      unit.imports.push_back(imp);
    }
  }

  if (j.has_field("Impls")) {
    for (const auto& impl_j : j.get("Impls").as_array()) {
      IrImpl imp;
      imp.relative_filename = deserialize_filename(impl_j.get("RelativeFilename"));
      unit.impls.push_back(imp);
    }
  }

  if (j.has_field("ImportDecls")) {
    for (const auto& d_j : j.get("ImportDecls").as_array()) {
      unit.import_decls.push_back(deserialize_decl(d_j));
    }
  }

  if (j.has_field("ImplDecls")) {
    for (const auto& d_j : j.get("ImplDecls").as_array()) {
      unit.impl_decls.push_back(deserialize_decl(d_j));
    }
  }

  if (j.has_field("Decls")) {
    for (const auto& d_j : j.get("Decls").as_array()) {
      unit.decls.push_back(deserialize_decl(d_j));
    }
  }

  if (j.has_field("Functions")) {
    for (const auto& f_j : j.get("Functions").as_array()) {
      unit.functions.push_back(deserialize_function(f_j));
    }
  }

  if (j.has_field("TraitImpls")) {
    for (const auto& ti_j : j.get("TraitImpls").as_array()) {
      unit.trait_impls.push_back(deserialize_trait_impl(ti_j));
    }
  }

  if (j.has_field("ImportedTraitImpls")) {
    for (const auto& ti_j : j.get("ImportedTraitImpls").as_array()) {
      unit.imported_trait_impls.push_back(deserialize_trait_impl(ti_j));
    }
  }

  return unit;
}

} // namespace ir
