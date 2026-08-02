#include "tests/test_util.h"
#include "bin/ir_base.h"
#include "bin/ir_type.h"
#include "ts/bind.h"
#include "ts/context.h"
#include "ts/rename_vars.h"

TEST(RenameVarsTest, CanonicalNaming) {
  // 1. VarType - Single variable
  {
    ts::Context type_ctx;
    auto input = ir::new_var_type("a");
    auto [_, res, err] = ts::rename_type_vars(type_ctx, input);
    EXPECT_TRUE(err.empty());
    EXPECT_TRUE(ir::equals_type(res, ir::new_var_type("a")));
  }

  // 2. ForallType - Single quantified variable - No rename
  {
    ts::Context type_ctx;
    auto input = ir::new_forall_type(ir::TypeParam{"a", ir::new_type_kind(), {}}, ir::new_var_type("a"));
    auto [_, res, err] = ts::rename_type_vars(type_ctx, input);
    EXPECT_TRUE(err.empty());
    EXPECT_TRUE(ir::equals_type(res, input));
  }

  // 3. ForallType - Single quantified variable - Renamed
  {
    ts::Context type_ctx;
    type_ctx = type_ctx.add_bind(ts::new_type_param_bind("a", ir::new_type_kind()));
    auto input = ir::new_forall_type(ir::TypeParam{"a", ir::new_type_kind(), {}}, ir::new_var_type("a"));
    auto [_, res, err] = ts::rename_type_vars(type_ctx, input);
    EXPECT_TRUE(err.empty());
    auto want = ir::new_forall_type(ir::TypeParam{"b", ir::new_type_kind(), {}}, ir::new_var_type("b"));
    EXPECT_TRUE(ir::equals_type(res, want));
  }

  // 4. ForallType - Multiple quantified variables - No rename
  {
    ts::Context type_ctx;
    auto input = ir::new_forall_type(
        ir::TypeParam{"a", ir::new_type_kind(), {}},
        ir::new_forall_type(
            ir::TypeParam{"b", ir::new_type_kind(), {}},
            ir::new_function_type(ir::new_var_type("a"), ir::new_var_type("b"))));
    auto [_, res, err] = ts::rename_type_vars(type_ctx, input);
    EXPECT_TRUE(err.empty());
    EXPECT_TRUE(ir::equals_type(res, input));
  }

  // 5. ForallType - Partially bound variables - Renamed
  {
    ts::Context type_ctx;
    type_ctx = type_ctx.add_bind(ts::new_type_param_bind("a", ir::new_type_kind()));
    auto input = ir::new_forall_type(
        ir::TypeParam{"a", ir::new_type_kind(), {}},
        ir::new_forall_type(
            ir::TypeParam{"b", ir::new_type_kind(), {}},
            ir::new_function_type(ir::new_var_type("a"), ir::new_var_type("b"))));
    auto [_, res, err] = ts::rename_type_vars(type_ctx, input);
    EXPECT_TRUE(err.empty());
    auto want = ir::new_forall_type(
        ir::TypeParam{"b", ir::new_type_kind(), {}},
        ir::new_forall_type(
            ir::TypeParam{"c", ir::new_type_kind(), {}},
            ir::new_function_type(ir::new_var_type("b"), ir::new_var_type("c"))));
    EXPECT_TRUE(ir::equals_type(res, want));
  }

  // 6. ForallType - Multiple bound variables - Renamed
  {
    ts::Context type_ctx;
    type_ctx = type_ctx.add_bind(ts::new_type_param_bind("a", ir::new_type_kind()));
    type_ctx = type_ctx.add_bind(ts::new_type_param_bind("b", ir::new_type_kind()));
    auto input = ir::new_forall_type(
        ir::TypeParam{"a", ir::new_type_kind(), {}},
        ir::new_forall_type(
            ir::TypeParam{"b", ir::new_type_kind(), {}},
            ir::new_function_type(ir::new_var_type("a"), ir::new_var_type("b"))));
    auto [_, res, err] = ts::rename_type_vars(type_ctx, input);
    EXPECT_TRUE(err.empty());
    auto want = ir::new_forall_type(
        ir::TypeParam{"c", ir::new_type_kind(), {}},
        ir::new_forall_type(
            ir::TypeParam{"d", ir::new_type_kind(), {}},
            ir::new_function_type(ir::new_var_type("c"), ir::new_var_type("d"))));
    EXPECT_TRUE(ir::equals_type(res, want));
  }

  // 7. LambdaType - Single abstracted variable - No rename
  {
    ts::Context type_ctx;
    auto input = ir::new_lambda_type("a", ir::new_type_kind(), ir::new_var_type("a"));
    auto [_, res, err] = ts::rename_type_vars(type_ctx, input);
    EXPECT_TRUE(err.empty());
    EXPECT_TRUE(ir::equals_type(res, input));
  }

  // 8. LambdaType - Single abstracted variable - Renamed
  {
    ts::Context type_ctx;
    type_ctx = type_ctx.add_bind(ts::new_type_param_bind("a", ir::new_type_kind()));
    auto input = ir::new_lambda_type("a", ir::new_type_kind(), ir::new_var_type("a"));
    auto [_, res, err] = ts::rename_type_vars(type_ctx, input);
    EXPECT_TRUE(err.empty());
    auto want = ir::new_lambda_type("b", ir::new_type_kind(), ir::new_var_type("b"));
    EXPECT_TRUE(ir::equals_type(res, want));
  }
}
