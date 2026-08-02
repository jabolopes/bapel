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
    EXPECT_EQ(t.to_string(), "forall ['a] 'a");
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
    EXPECT_EQ(t.to_string(), "forall ['a, 'b] 'a -> 'b");
  }
}

TEST(IrTypeTest, Formatting) {
  // Primitive types
  EXPECT_EQ(ir::new_name_type("i32").to_string(), "i32");
  EXPECT_EQ(ir::new_var_type("a").to_string(), "'a");
  EXPECT_EQ(ir::new_var_type("'a").to_string(), "'a");

  // Tuple types
  EXPECT_EQ(ir::new_tuple_type({}).to_string(), "()");
  EXPECT_EQ(ir::new_tuple_type({ir::new_name_type("i8"), ir::new_name_type("i16")}).to_string(), "(i8, i16)");

  // Struct types
  {
    std::vector<ir::StructField> fields = {
        {"x", std::make_shared<ir::IrType>(ir::new_name_type("i32"))},
        {"y", std::make_shared<ir::IrType>(ir::new_name_type("i32"))}};
    EXPECT_EQ(ir::new_struct_type(std::move(fields)).to_string(), "struct{x: i32, y: i32}");
  }

  // Variant types
  {
    std::vector<ir::VariantTag> tags = {
        {"none", std::make_shared<ir::IrType>(ir::new_tuple_type({}))},
        {"some", std::make_shared<ir::IrType>(ir::new_name_type("i64"))}};
    EXPECT_EQ(ir::new_variant_type(std::move(tags)).to_string(), "variant{none (), some i64}");
  }
}
