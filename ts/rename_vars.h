#pragma once

#include "bin/ir_base.h"
#include "bin/ir_decl.h"
#include "bin/ir_type.h"
#include "ts/bind.h"
#include "ts/context.h"
#include "ts/list.h"

#include <string>
#include <tuple>
#include <vector>

namespace ts {

struct TypeSubstitution {
  ir::IrType source;
  ir::IrType target;
};

class TypeVarRenamer {
 public:
  explicit TypeVarRenamer(Context context) : context_(std::move(context)) {}

  bool lookup_substitution(const std::string& source_tvar, ir::IrType& out_target) const {
    const auto* sub = substitutions_.find_if([&](const TypeSubstitution& s) {
      return s.source.data.is(ir::IrTypeData::Tag::var_t) && s.source.data.as_var_t() == source_tvar;
    });
    if (sub) {
      out_target = sub->target;
      return true;
    }
    return false;
  }

  template <typename Callback>
  ir::IrType with_type_substitution(
      const std::string& source_tvar,
      const ir::IrKind& source_kind,
      const std::vector<ir::IrType>& source_bounds,
      Callback&& callback) {
    Context orig_context = context_;
    List<TypeSubstitution> orig_substitutions = substitutions_;

    ir::IrType target_tvar = context_.gen_fresh_var_type();
    const std::string target_name = target_tvar.data.as_var_t();
    substitutions_ = substitutions_.add(TypeSubstitution{ir::new_var_type(source_tvar), target_tvar});

    std::vector<ir::IrType> renamed_bounds;
    renamed_bounds.reserve(source_bounds.size());
    for (const auto& b : source_bounds) {
      renamed_bounds.push_back(rename(b));
    }

    try {
      context_ = context_.add_bind(new_type_param_bind(target_name, source_kind, renamed_bounds));
    } catch (const std::exception& e) {
      if (!err_.empty()) err_ += "\n";
      err_ += e.what();
      return ir::IrType{};
    }

    ir::IrType result = callback(target_name);

    context_ = orig_context;
    substitutions_ = orig_substitutions;
    return result;
  }

  ir::IrType rename(const ir::IrType& typ) {
    if (!err_.empty()) {
      return ir::IrType{};
    }
    return rename_impl(typ);
  }

  const std::string& error() const { return err_; }
  const Context& context() const { return context_; }

 private:
  ir::IrType rename_impl(const ir::IrType& typ) {
    switch (typ.data.tag()) {
      case ir::IrTypeData::Tag::forall_t: {
        const auto& fa = typ.data.as_forall_t();
        const auto& tp = fa.param;
        return with_type_substitution(tp.var, tp.kind, tp.bounds, [&](const std::string& target_tvar) {
          std::vector<ir::IrType> bounds;
          bounds.reserve(tp.bounds.size());
          for (const auto& b : tp.bounds) {
            bounds.push_back(rename(b));
          }
          ir::TypeParam new_tp{target_tvar, tp.kind, std::move(bounds)};
          ir::IrType body = fa.body ? rename(*fa.body) : ir::IrType{};
          return ir::new_forall_type(std::move(new_tp), std::move(body));
        });
      }
      case ir::IrTypeData::Tag::lambda_t: {
        const auto& l = typ.data.as_lambda_t();
        return with_type_substitution(l.var, l.kind, {}, [&](const std::string& target_tvar) {
          ir::IrType body = l.body ? rename(*l.body) : ir::IrType{};
          return ir::new_lambda_type(target_tvar, l.kind, std::move(body));
        });
      }
      case ir::IrTypeData::Tag::mu_t: {
        const auto& m = typ.data.as_mu_t();
        return with_type_substitution(m.var, m.kind, {}, [&](const std::string& target_tvar) {
          ir::IrType body = m.body ? rename(*m.body) : ir::IrType{};
          return ir::new_mu_type(target_tvar, m.kind, std::move(body));
        });
      }
      case ir::IrTypeData::Tag::var_t: {
        ir::IrType target;
        if (lookup_substitution(typ.data.as_var_t(), target)) {
          return target;
        }
        return typ;
      }
      default:
        return ir::map_type(typ, [&](const ir::IrType& child) {
          return rename(child);
        });
    }
  }

  Context context_;
  List<TypeSubstitution> substitutions_;
  std::string err_;
};

inline std::tuple<Context, ir::IrType, std::string> rename_type_vars(
    const Context& context,
    const ir::IrType& typ) {
  TypeVarRenamer renamer(context);
  ir::IrType result = renamer.rename(typ);
  return {renamer.context(), std::move(result), renamer.error()};
}

} // namespace ts
