#include "bin/ir_base.h"
#include "bin/ir_type.h"
#include "ts/bind.h"
#include "ts/context.h"
#include "ts/infer_kind.h"
#include "ts/kind.h"
#include "ts/list.h"
#include "ts/reduce.h"
#include "ts/reduce_type.h"
#include "ts/rename_vars.h"
#include "ts/wellformed.h"

#include <cassert>
#include <iostream>
#include <string>
#include <vector>

void test_rename_type_vars() {
  std::cout << "Running test_rename_type_vars..." << std::endl;

  // 1. VarType - Single variable
  {
    ts::Context ctx;
    auto input = ir::new_var_type("a");
    auto [_, res, err] = ts::rename_type_vars(ctx, input);
    assert(err.empty());
    assert(ir::equals_type(res, ir::new_var_type("a")));
  }

  // 2. ForallType - Single quantified variable - No rename
  {
    ts::Context ctx;
    auto input = ir::new_forall_type(ir::TypeParam{"a", ir::new_type_kind(), {}}, ir::new_var_type("a"));
    auto [_, res, err] = ts::rename_type_vars(ctx, input);
    assert(err.empty());
    assert(ir::equals_type(res, input));
  }

  // 3. ForallType - Single quantified variable - Renamed
  {
    ts::Context ctx;
    ctx = ctx.add_bind(ts::new_type_param_bind("a", ir::new_type_kind()));
    auto input = ir::new_forall_type(ir::TypeParam{"a", ir::new_type_kind(), {}}, ir::new_var_type("a"));
    auto [_, res, err] = ts::rename_type_vars(ctx, input);
    assert(err.empty());
    auto want = ir::new_forall_type(ir::TypeParam{"b", ir::new_type_kind(), {}}, ir::new_var_type("b"));
    assert(ir::equals_type(res, want));
  }

  // 4. ForallType - Multiple quantified variables - No rename
  {
    ts::Context ctx;
    auto input = ir::new_forall_type(
        ir::TypeParam{"a", ir::new_type_kind(), {}},
        ir::new_forall_type(
            ir::TypeParam{"b", ir::new_type_kind(), {}},
            ir::new_function_type(ir::new_var_type("a"), ir::new_var_type("b"))));
    auto [_, res, err] = ts::rename_type_vars(ctx, input);
    assert(err.empty());
    assert(ir::equals_type(res, input));
  }

  // 5. ForallType - Partially bound variables - Renamed
  {
    ts::Context ctx;
    ctx = ctx.add_bind(ts::new_type_param_bind("a", ir::new_type_kind()));
    auto input = ir::new_forall_type(
        ir::TypeParam{"a", ir::new_type_kind(), {}},
        ir::new_forall_type(
            ir::TypeParam{"b", ir::new_type_kind(), {}},
            ir::new_function_type(ir::new_var_type("a"), ir::new_var_type("b"))));
    auto [_, res, err] = ts::rename_type_vars(ctx, input);
    assert(err.empty());
    auto want = ir::new_forall_type(
        ir::TypeParam{"b", ir::new_type_kind(), {}},
        ir::new_forall_type(
            ir::TypeParam{"c", ir::new_type_kind(), {}},
            ir::new_function_type(ir::new_var_type("b"), ir::new_var_type("c"))));
    assert(ir::equals_type(res, want));
  }

  // 6. ForallType - Multiple bound variables - Renamed
  {
    ts::Context ctx;
    ctx = ctx.add_bind(ts::new_type_param_bind("a", ir::new_type_kind()));
    ctx = ctx.add_bind(ts::new_type_param_bind("b", ir::new_type_kind()));
    auto input = ir::new_forall_type(
        ir::TypeParam{"a", ir::new_type_kind(), {}},
        ir::new_forall_type(
            ir::TypeParam{"b", ir::new_type_kind(), {}},
            ir::new_function_type(ir::new_var_type("a"), ir::new_var_type("b"))));
    auto [_, res, err] = ts::rename_type_vars(ctx, input);
    assert(err.empty());
    auto want = ir::new_forall_type(
        ir::TypeParam{"c", ir::new_type_kind(), {}},
        ir::new_forall_type(
            ir::TypeParam{"d", ir::new_type_kind(), {}},
            ir::new_function_type(ir::new_var_type("c"), ir::new_var_type("d"))));
    assert(ir::equals_type(res, want));
  }

  // 7. LambdaType - Single abstracted variable - No rename
  {
    ts::Context ctx;
    auto input = ir::new_lambda_type("a", ir::new_type_kind(), ir::new_var_type("a"));
    auto [_, res, err] = ts::rename_type_vars(ctx, input);
    assert(err.empty());
    assert(ir::equals_type(res, input));
  }

  // 8. LambdaType - Single abstracted variable - Renamed
  {
    ts::Context ctx;
    ctx = ctx.add_bind(ts::new_type_param_bind("a", ir::new_type_kind()));
    auto input = ir::new_lambda_type("a", ir::new_type_kind(), ir::new_var_type("a"));
    auto [_, res, err] = ts::rename_type_vars(ctx, input);
    assert(err.empty());
    auto want = ir::new_lambda_type("b", ir::new_type_kind(), ir::new_var_type("b"));
    assert(ir::equals_type(res, want));
  }

  // 9. StructType - Fields freshened with substitutions
  {
    ts::Context ctx;
    auto input = ir::new_struct_type({
        ir::StructField{"f1", std::make_shared<ir::IrType>(ir::new_var_type("a"))},
        ir::StructField{"f2", std::make_shared<ir::IrType>(ir::new_function_type(ir::new_var_type("b"), ir::new_var_type("c")))},
    });
    std::vector<ts::TypeSubstitution> subs = {
        {ir::new_var_type("a"), ir::new_var_type("b")},
        {ir::new_var_type("b"), ir::new_var_type("c")},
    };
    auto [_, res, err] = ts::rename_type_vars_with_substitutions(ctx, input, subs);
    assert(err.empty());
    auto want = ir::new_struct_type({
        ir::StructField{"f1", std::make_shared<ir::IrType>(ir::new_var_type("b"))},
        ir::StructField{"f2", std::make_shared<ir::IrType>(ir::new_function_type(ir::new_var_type("c"), ir::new_var_type("c")))},
    });
    assert(ir::equals_type(res, want));
  }

  // 10. TupleType - Elements freshened
  {
    ts::Context ctx;
    auto input = ir::new_tuple_type({ir::new_var_type("a"), ir::new_var_type("b")});
    std::vector<ts::TypeSubstitution> subs = {
        {ir::new_var_type("a"), ir::new_var_type("b")},
        {ir::new_var_type("b"), ir::new_var_type("c")},
    };
    auto [_, res, err] = ts::rename_type_vars_with_substitutions(ctx, input, subs);
    assert(err.empty());
    auto want = ir::new_tuple_type({ir::new_var_type("b"), ir::new_var_type("c")});
    assert(ir::equals_type(res, want));
  }

  // 11. VariantType - Tags freshened
  {
    ts::Context ctx;
    auto input = ir::new_variant_type({
        ir::VariantTag{"TagA", std::make_shared<ir::IrType>(ir::new_var_type("a"))},
        ir::VariantTag{"TagB", std::make_shared<ir::IrType>(ir::new_function_type(ir::new_var_type("b"), ir::new_var_type("c")))},
    });
    std::vector<ts::TypeSubstitution> subs = {
        {ir::new_var_type("a"), ir::new_var_type("b")},
        {ir::new_var_type("b"), ir::new_var_type("c")},
    };
    auto [_, res, err] = ts::rename_type_vars_with_substitutions(ctx, input, subs);
    assert(err.empty());
    auto want = ir::new_variant_type({
        ir::VariantTag{"TagA", std::make_shared<ir::IrType>(ir::new_var_type("b"))},
        ir::VariantTag{"TagB", std::make_shared<ir::IrType>(ir::new_function_type(ir::new_var_type("c"), ir::new_var_type("c")))},
    });
    assert(ir::equals_type(res, want));
  }

  // 12. Complex nested type with context collision
  {
    ts::Context ctx;
    ctx = ctx.add_bind(ts::new_type_param_bind("a", ir::new_type_kind()));
    auto input = ir::new_forall_type(
        ir::TypeParam{"a", ir::new_type_kind(), {}},
        ir::new_function_type(
            ir::new_var_type("a"),
            ir::new_array_type(ir::new_app_type(ir::new_name_type("List"), ir::new_var_type("a")), 5)));
    auto [_, res, err] = ts::rename_type_vars(ctx, input);
    assert(err.empty());
    auto want = ir::new_forall_type(
        ir::TypeParam{"b", ir::new_type_kind(), {}},
        ir::new_function_type(
            ir::new_var_type("b"),
            ir::new_array_type(ir::new_app_type(ir::new_name_type("List"), ir::new_var_type("b")), 5)));
    assert(ir::equals_type(res, want));
  }

  std::cout << "test_rename_type_vars PASSED!" << std::endl;
}

void test_wellformed_type() {
  std::cout << "Running test_wellformed_type..." << std::endl;

  ts::Context ctx;
  ctx = ctx.add_bind(ts::new_const_bind("i32", ir::new_type_kind()));
  ctx = ctx.add_bind(ts::new_const_bind("bool", ir::new_type_kind()));
  ctx = ctx.add_bind(ts::new_alias_bind("Integer", ir::new_name_type("i32")));
  ctx = ctx.add_bind(ts::new_trait_bind("Size", {}, {}));
  ctx = ctx.add_bind(ts::new_type_param_bind("T", ir::new_type_kind()));

  std::string err;

  // Primitives and defined names
  assert(ts::is_wellformed_type(ctx, ir::new_name_type("i32"), err));
  assert(ts::is_wellformed_type(ctx, ir::new_name_type("Integer"), err));
  assert(ts::is_wellformed_type(ctx, ir::new_name_type("Size"), err));
  assert(ts::is_wellformed_type(ctx, ir::new_var_type("T"), err));

  // Undefined name
  assert(!ts::is_wellformed_type(ctx, ir::new_name_type("Float"), err));
  assert(err.find("undefined") != std::string::npos);

  // Undefined var
  assert(!ts::is_wellformed_type(ctx, ir::new_var_type("U"), err));

  // Function, tuple, array
  assert(ts::is_wellformed_type(ctx, ir::new_function_type(ir::new_name_type("i32"), ir::new_name_type("bool")), err));
  assert(ts::is_wellformed_type(ctx, ir::new_tuple_type({ir::new_name_type("i32"), ir::new_var_type("T")}), err));
  assert(ts::is_wellformed_type(ctx, ir::new_array_type(ir::new_name_type("i32"), 10), err));

  // Struct with unique fields
  auto valid_struct = ir::new_struct_type({
      ir::StructField{"x", std::make_shared<ir::IrType>(ir::new_name_type("i32"))},
      ir::StructField{"y", std::make_shared<ir::IrType>(ir::new_name_type("bool"))},
  });
  assert(ts::is_wellformed_type(ctx, valid_struct, err));

  // Struct with duplicate fields
  auto dup_struct = ir::new_struct_type({
      ir::StructField{"x", std::make_shared<ir::IrType>(ir::new_name_type("i32"))},
      ir::StructField{"x", std::make_shared<ir::IrType>(ir::new_name_type("bool"))},
  });
  assert(!ts::is_wellformed_type(ctx, dup_struct, err));
  assert(err.find("duplicate") != std::string::npos);

  // Variant with unique tags
  auto valid_variant = ir::new_variant_type({
      ir::VariantTag{"None", std::make_shared<ir::IrType>(ir::new_tuple_type({}))},
      ir::VariantTag{"Some", std::make_shared<ir::IrType>(ir::new_name_type("i32"))},
  });
  assert(ts::is_wellformed_type(ctx, valid_variant, err));

  // Variant with duplicate tags
  auto dup_variant = ir::new_variant_type({
      ir::VariantTag{"Some", std::make_shared<ir::IrType>(ir::new_name_type("i32"))},
      ir::VariantTag{"Some", std::make_shared<ir::IrType>(ir::new_name_type("bool"))},
  });
  assert(!ts::is_wellformed_type(ctx, dup_variant, err));
  assert(err.find("duplicate") != std::string::npos);

  // Forall and Lambda well-formedness
  auto valid_forall = ir::new_forall_type(ir::TypeParam{"a", ir::new_type_kind(), {}}, ir::new_function_type(ir::new_var_type("a"), ir::new_name_type("i32")));
  assert(ts::is_wellformed_type(ctx, valid_forall, err));

  auto valid_lambda = ir::new_lambda_type("a", ir::new_type_kind(), ir::new_function_type(ir::new_var_type("a"), ir::new_name_type("i32")));
  assert(ts::is_wellformed_type(ctx, valid_lambda, err));

  std::cout << "test_wellformed_type PASSED!" << std::endl;
}

void test_infer_kind() {
  std::cout << "Running test_infer_kind..." << std::endl;

  ts::Context ctx;
  ctx = ctx.add_bind(ts::new_const_bind("i32", ir::new_type_kind()));
  ctx = ctx.add_bind(ts::new_const_bind("Option", ir::new_arrow_kind(ir::new_type_kind(), ir::new_type_kind())));
  ctx = ctx.add_bind(ts::new_type_param_bind("T", ir::new_type_kind()));
  ctx = ctx.add_bind(ts::new_type_param_bind("F", ir::new_arrow_kind(ir::new_type_kind(), ir::new_type_kind())));

  // Primitive kind
  assert(ir::equals_kind(ts::infer_kind(ctx, ir::new_name_type("i32")), ir::new_type_kind()));

  // Type constructor kind
  assert(ir::equals_kind(ts::infer_kind(ctx, ir::new_name_type("Option")), ir::new_arrow_kind(ir::new_type_kind(), ir::new_type_kind())));

  // Type variable kinds
  assert(ir::equals_kind(ts::infer_kind(ctx, ir::new_var_type("T")), ir::new_type_kind()));
  assert(ir::equals_kind(ts::infer_kind(ctx, ir::new_var_type("F")), ir::new_arrow_kind(ir::new_type_kind(), ir::new_type_kind())));

  // Type Application: Option i32 -> *
  auto app_type = ir::new_app_type(ir::new_name_type("Option"), ir::new_name_type("i32"));
  assert(ir::equals_kind(ts::infer_kind(ctx, app_type), ir::new_type_kind()));

  // Higher-Kinded Application: F i32 -> *
  auto hk_app = ir::new_app_type(ir::new_var_type("F"), ir::new_name_type("i32"));
  assert(ir::equals_kind(ts::infer_kind(ctx, hk_app), ir::new_type_kind()));

  // Lambda Kind: fun ('a :: *) -> 'a -> * -> *
  auto lambda_t = ir::new_lambda_type("a", ir::new_type_kind(), ir::new_var_type("a"));
  assert(ir::equals_kind(ts::infer_kind(ctx, lambda_t), ir::new_arrow_kind(ir::new_type_kind(), ir::new_type_kind())));

  // Higher-order Lambda: fun ('f :: * -> *) -> 'f i32 -> (* -> *) -> *
  auto ho_lambda = ir::new_lambda_type(
      "f",
      ir::new_arrow_kind(ir::new_type_kind(), ir::new_type_kind()),
      ir::new_app_type(ir::new_var_type("f"), ir::new_name_type("i32")));
  auto want_ho_kind = ir::new_arrow_kind(ir::new_arrow_kind(ir::new_type_kind(), ir::new_type_kind()), ir::new_type_kind());
  assert(ir::equals_kind(ts::infer_kind(ctx, ho_lambda), want_ho_kind));

  // Kind mismatch in application: i32 i32 -> error
  ir::IrKind err_kind;
  std::string err_msg;
  assert(!ts::infer_kind(ctx, ir::new_app_type(ir::new_name_type("i32"), ir::new_name_type("i32")), err_kind, err_msg));
  assert(err_msg.find("expected arrow kind") != std::string::npos);

  std::cout << "test_infer_kind PASSED!" << std::endl;
}

void test_reduce_type() {
  std::cout << "Running test_reduce_type..." << std::endl;

  ts::Context ctx;
  ctx = ctx.add_bind(ts::new_const_bind("i32", ir::new_type_kind()));
  ctx = ctx.add_bind(ts::new_const_bind("String", ir::new_type_kind()));
  ctx = ctx.add_bind(ts::new_alias_bind("MyInt", ir::new_name_type("i32")));
  ctx = ctx.add_bind(ts::new_alias_bind("MyPair", ir::new_lambda_type("a", ir::new_type_kind(), ir::new_tuple_type({ir::new_var_type("a"), ir::new_var_type("a")}))));

  // 1. Alias reduction: MyInt -> i32
  assert(ir::equals_type(ts::reduce_type(ctx, ir::new_name_type("MyInt")), ir::new_name_type("i32")));

  // 2. Beta reduction: (fun ('a) -> ('a, 'a)) String -> (String, String)
  auto app_lambda = ir::new_app_type(
      ir::new_lambda_type("a", ir::new_type_kind(), ir::new_tuple_type({ir::new_var_type("a"), ir::new_var_type("a")})),
      ir::new_name_type("String"));
  auto reduced_pair = ts::reduce_type(ctx, app_lambda);
  auto want_pair = ir::new_tuple_type({ir::new_name_type("String"), ir::new_name_type("String")});
  assert(ir::equals_type(reduced_pair, want_pair));

  // 3. Nested Beta reduction:
  // (fun ('f) -> fun ('x) -> 'f 'x) (fun ('a) -> ['a, 5]) i32 -> [i32, 5]
  auto id_array = ir::new_lambda_type("a", ir::new_type_kind(), ir::new_array_type(ir::new_var_type("a"), 5));
  auto apply_fn = ir::new_lambda_type(
      "f",
      ir::new_arrow_kind(ir::new_type_kind(), ir::new_type_kind()),
      ir::new_lambda_type("x", ir::new_type_kind(), ir::new_app_type(ir::new_var_type("f"), ir::new_var_type("x"))));
  auto step1 = ir::new_app_type(apply_fn, id_array);
  auto step2 = ir::new_app_type(step1, ir::new_name_type("i32"));
  auto fully_reduced = ts::reduce_type(ctx, step2);
  auto want_array = ir::new_array_type(ir::new_name_type("i32"), 5);
  assert(ir::equals_type(fully_reduced, want_array));

  // 4. Reduction inside struct
  auto struct_with_alias = ir::new_struct_type({
      ir::StructField{"count", std::make_shared<ir::IrType>(ir::new_name_type("MyInt"))},
  });
  auto reduced_struct = ts::reduce_type(ctx, struct_with_alias);
  auto want_struct = ir::new_struct_type({
      ir::StructField{"count", std::make_shared<ir::IrType>(ir::new_name_type("i32"))},
  });
  assert(ir::equals_type(reduced_struct, want_struct));

  std::cout << "test_reduce_type PASSED!" << std::endl;
}

int main() {
  std::cout << "========================================" << std::endl;
  std::cout << "Running Phase 2 Native C++ Test Suite..." << std::endl;
  std::cout << "========================================" << std::endl;

  test_rename_type_vars();
  test_wellformed_type();
  test_infer_kind();
  test_reduce_type();

  std::cout << "========================================" << std::endl;
  std::cout << "ALL PHASE 2 TESTS PASSED SUCCESSFULLY!" << std::endl;
  std::cout << "========================================" << std::endl;
  return 0;
}
