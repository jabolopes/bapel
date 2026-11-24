#include "tests/test_util.h"
#include "bin/ir_base.h"
#include "bin/ir_type.h"

TEST(IrTypeTest, ForallVars) {
  // Empty params
  {
    auto t = ir::forall_vars({}, ir::new_name_type("i8"));
    EXPECT_TRUE(ir::equals_type(t, ir::new_name_type("i8")));
  }

  // Single param
  {
    std::vector<ir::TypeParam> tps = {ir::TypeParam{"a", ir::new_type_kind(), {}}};
    auto t = ir::forall_vars(tps, ir::new_var_type("a"));
    auto want = ir::new_forall_type(ir::TypeParam{"a", ir::new_type_kind(), {}}, ir::new_var_type("a"));
    EXPECT_TRUE(ir::equals_type(t, want));
    EXPECT_EQ(ir::type_to_string(t), "forall ['a] 'a");
  }

  // Multiple params
  {
    std::vector<ir::TypeParam> tps = {
        ir::TypeParam{"a", ir::new_type_kind(), {}},
        ir::TypeParam{"b", ir::new_type_kind(), {}}};
    auto t = ir::forall_vars(tps, ir::new_function_type(ir::new_var_type("a"), ir::new_var_type("b")));
    auto want = ir::new_forall_type(
        ir::TypeParam{"a", ir::new_type_kind(), {}},
        ir::new_forall_type(
            ir::TypeParam{"b", ir::new_type_kind(), {}},
            ir::new_function_type(ir::new_var_type("a"), ir::new_var_type("b"))));
    EXPECT_TRUE(ir::equals_type(t, want));
    EXPECT_EQ(ir::type_to_string(t), "forall ['a, 'b] 'a -> 'b");
  }
}

TEST(IrTypeTest, Formatting) {
  // Primitive types
  EXPECT_EQ(ir::type_to_string(ir::new_name_type("i32")), "i32");
  EXPECT_EQ(ir::type_to_string(ir::new_var_type("a")), "'a");
  EXPECT_EQ(ir::type_to_string(ir::new_var_type("'a")), "'a");

  // Tuple types
  EXPECT_EQ(ir::type_to_string(ir::new_tuple_type({})), "()");
  EXPECT_EQ(ir::type_to_string(ir::new_tuple_type({ir::new_name_type("i8"), ir::new_name_type("i16")})), "(i8, i16)");

  // Struct types
  {
    std::vector<ir::StructField> fields = {
        {"x", std::make_shared<ir::IrType>(ir::new_name_type("i32"))},
        {"y", std::make_shared<ir::IrType>(ir::new_name_type("i32"))}};
    EXPECT_EQ(ir::type_to_string(ir::new_struct_type(std::move(fields))), "struct{x: i32, y: i32}");
  }

  // Variant types
  {
    std::vector<ir::VariantTag> tags = {
        {"none", std::make_shared<ir::IrType>(ir::new_tuple_type({}))},
        {"some", std::make_shared<ir::IrType>(ir::new_name_type("i64"))}};
    EXPECT_EQ(ir::type_to_string(ir::new_variant_type(std::move(tags))), "variant{none (), some i64}");
  }

  // Mu types
  {
    auto t = ir::new_mu_type(
        "self",
        ir::new_type_kind(),
        ir::new_tuple_type({ir::new_name_type("i64"), ir::new_var_type("self")}));
    EXPECT_EQ(ir::type_to_string(t), "mu 'self. (i64, 'self)");
  }
}

TEST(IrTypeTest, SubstituteType) {
  // Simple substitution
  {
    auto t = ir::new_var_type("a");
    auto res = ir::substitute_type(t, ir::new_var_type("a"), ir::new_name_type("i32"));
    EXPECT_TRUE(ir::equals_type(res, ir::new_name_type("i32")));
  }

  // Shadowing in ForallType
  {
    auto t = ir::new_forall_type(ir::TypeParam{"a", ir::new_type_kind(), {}}, ir::new_var_type("a"));
    auto res = ir::substitute_type(t, ir::new_var_type("a"), ir::new_name_type("i32"));
    EXPECT_TRUE(ir::equals_type(res, t));
  }

  // Capture avoidance in ForallType: forall ['b] 'a -> ('a, 'b) substituted with ['a -> 'b]
  {
    auto inner = ir::new_forall_type(
        ir::TypeParam{"b", ir::new_type_kind(), {}},
        ir::new_function_type(ir::new_var_type("a"), ir::new_tuple_type({ir::new_var_type("a"), ir::new_var_type("b")})));
    auto res = ir::substitute_type(inner, ir::new_var_type("a"), ir::new_var_type("b"));
    // 'b' in the inner quantifier should be alpha-renamed to 'c' to prevent capturing 'b'
    EXPECT_EQ(ir::type_to_string(res), "forall ['c] 'b -> ('b, 'c)");
  }

  // Shadowing in LambdaType
  {
    auto t = ir::new_lambda_type("a", ir::new_type_kind(), ir::new_var_type("a"));
    auto res = ir::substitute_type(t, ir::new_var_type("a"), ir::new_name_type("i32"));
    EXPECT_TRUE(ir::equals_type(res, t));
  }

  // Capture avoidance in LambdaType
  {
    auto t = ir::new_lambda_type("b", ir::new_type_kind(), ir::new_var_type("a"));
    auto res = ir::substitute_type(t, ir::new_var_type("a"), ir::new_var_type("b"));
    EXPECT_EQ(ir::type_to_string(res), "fun (c) ('b)");
  }

  // Free variable substitution in MuType
  {
    auto t = ir::new_mu_type(
        "self",
        ir::new_type_kind(),
        ir::new_tuple_type({ir::new_var_type("a"), ir::new_var_type("self")}));
    auto res = ir::substitute_type(t, ir::new_var_type("a"), ir::new_name_type("i64"));
    auto want = ir::new_mu_type(
        "self",
        ir::new_type_kind(),
        ir::new_tuple_type({ir::new_name_type("i64"), ir::new_var_type("self")}));
    EXPECT_TRUE(ir::equals_type(res, want));
  }

  // Shadowing in MuType
  {
    auto t = ir::new_mu_type("a", ir::new_type_kind(), ir::new_var_type("a"));
    auto res = ir::substitute_type(t, ir::new_var_type("a"), ir::new_name_type("i32"));
    EXPECT_TRUE(ir::equals_type(res, t));
  }

  // Capture avoidance in MuType: mu 'b. ('a, 'b) substituted with ['a -> 'b]
  {
    auto t = ir::new_mu_type(
        "b",
        ir::new_type_kind(),
        ir::new_tuple_type({ir::new_var_type("a"), ir::new_var_type("b")}));
    auto res = ir::substitute_type(t, ir::new_var_type("a"), ir::new_var_type("b"));
    EXPECT_EQ(ir::type_to_string(res), "mu 'c. ('b, 'c)");
  }
}

TEST(IrTypeTest, MuTypeEquality) {
  auto t1 = ir::new_mu_type(
      "self",
      ir::new_type_kind(),
      ir::new_tuple_type({ir::new_name_type("i64"), ir::new_var_type("self")}));
  auto t2 = ir::new_mu_type(
      "self",
      ir::new_type_kind(),
      ir::new_tuple_type({ir::new_name_type("i64"), ir::new_var_type("self")}));
  auto t3 = ir::new_mu_type(
      "other",
      ir::new_type_kind(),
      ir::new_tuple_type({ir::new_name_type("i64"), ir::new_var_type("other")}));
  auto t4 = ir::new_mu_type(
      "self",
      ir::new_type_kind(),
      ir::new_tuple_type({ir::new_name_type("i32"), ir::new_var_type("self")}));

  EXPECT_TRUE(ir::equals_type(t1, t2));
  EXPECT_FALSE(ir::equals_type(t1, t3));
  EXPECT_FALSE(ir::equals_type(t1, t4));
}

