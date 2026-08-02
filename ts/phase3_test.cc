#include "bin/ir_base.h"
#include "bin/ir_decl.h"
#include "bin/ir_function.h"
#include "bin/ir_term.h"
#include "bin/ir_type.h"
#include "ts/bind.h"
#include "ts/context.h"
#include "ts/infer_kind.h"
#include "ts/inferencer.h"
#include "ts/predicative.h"
#include "ts/reduce_type.h"
#include "ts/returns.h"
#include "ts/solver.h"
#include "ts/subtype.h"
#include "ts/typecheck.h"
#include "ts/unify.h"
#include "ts/wellformed.h"

#include <cassert>
#include <iostream>
#include <string>
#include <vector>

void test_subtyping() {
  std::cout << "Running test_subtyping..." << std::endl;

  ts::Context ctx;
  ctx = ctx.add_bind(ts::new_const_bind("i32", ir::new_type_kind()));
  ctx = ctx.add_bind(ts::new_const_bind("bool", ir::new_type_kind()));
  ctx = ctx.add_bind(ts::new_const_bind("String", ir::new_type_kind()));
  ctx = ctx.add_bind(ts::new_alias_bind("Integer", ir::new_name_type("i32")));
  ctx = ctx.add_bind(ts::new_trait_bind("Size", {}, {}));
  ctx = ctx.add_bind(ts::new_trait_impl_bind({}, ir::new_name_type("Size"), ir::new_name_type("String")));

  std::string err;

  // 1. Primitive and alias subtyping
  assert(ts::subtype(ctx, ir::new_name_type("i32"), ir::new_name_type("i32"), err));
  assert(ts::subtype(ctx, ir::new_name_type("Integer"), ir::new_name_type("i32"), err));
  assert(ts::subtype(ctx, ir::new_name_type("i32"), ir::new_name_type("Integer"), err));

  // 2. Function subtyping: contravariant arg, covariant ret
  // (Integer -> String) <: (i32 -> String)
  auto fn1 = ir::new_function_type(ir::new_name_type("Integer"), ir::new_name_type("String"));
  auto fn2 = ir::new_function_type(ir::new_name_type("i32"), ir::new_name_type("String"));
  assert(ts::subtype(ctx, fn1, fn2, err));

  // 3. Struct subtyping
  auto s1 = ir::new_struct_type({
      ir::StructField{"x", std::make_shared<ir::IrType>(ir::new_name_type("Integer"))},
      ir::StructField{"y", std::make_shared<ir::IrType>(ir::new_name_type("String"))},
  });
  auto s2 = ir::new_struct_type({
      ir::StructField{"x", std::make_shared<ir::IrType>(ir::new_name_type("i32"))},
      ir::StructField{"y", std::make_shared<ir::IrType>(ir::new_name_type("String"))},
  });
  assert(ts::subtype(ctx, s1, s2, err));

  // 4. Tuple subtyping
  auto t1 = ir::new_tuple_type({ir::new_name_type("Integer"), ir::new_name_type("String")});
  auto t2 = ir::new_tuple_type({ir::new_name_type("i32"), ir::new_name_type("String")});
  assert(ts::subtype(ctx, t1, t2, err));

  // 5. Variant subtyping
  auto v1 = ir::new_variant_type({
      ir::VariantTag{"Val", std::make_shared<ir::IrType>(ir::new_name_type("Integer"))},
  });
  auto v2 = ir::new_variant_type({
      ir::VariantTag{"Val", std::make_shared<ir::IrType>(ir::new_name_type("i32"))},
  });
  assert(ts::subtype(ctx, v1, v2, err));

  // 6. Trait bounds satisfaction
  assert(ts::satisfies_bound(ctx, ir::new_name_type("String"), ir::new_name_type("Size"), err));
  assert(!ts::satisfies_bound(ctx, ir::new_name_type("i32"), ir::new_name_type("Size"), err));

  std::cout << "test_subtyping PASSED!" << std::endl;
}

void test_unification() {
  std::cout << "Running test_unification..." << std::endl;

  ts::Context ctx;
  ctx = ctx.add_bind(ts::new_const_bind("i32", ir::new_type_kind()));
  ctx = ctx.add_bind(ts::new_const_bind("bool", ir::new_type_kind()));

  ts::UnificationState state;
  ir::IrType evar0 = state.new_evar(ctx); // ^0
  ir::IrType evar1 = state.new_evar(ctx); // ^1

  std::string err;

  // Unify ^0 with i32
  assert(ts::unify(ctx, state, evar0, ir::new_name_type("i32"), err));
  assert(state.is_exist_var_assigned(evar0));
  assert(ir::equals_type(state.exist_var_solution(evar0), ir::new_name_type("i32")));

  // Unify (^1 -> bool) with (i32 -> bool)
  auto fn_evar = ir::new_function_type(evar1, ir::new_name_type("bool"));
  auto fn_concrete = ir::new_function_type(ir::new_name_type("i32"), ir::new_name_type("bool"));
  assert(ts::unify(ctx, state, fn_evar, fn_concrete, err));
  assert(state.is_exist_var_assigned(evar1));
  assert(ir::equals_type(state.exist_var_solution(evar1), ir::new_name_type("i32")));

  // Incompatible unification: bool with i32
  assert(!ts::unify(ctx, state, ir::new_name_type("bool"), ir::new_name_type("i32"), err));

  std::cout << "test_unification PASSED!" << std::endl;
}

void test_inference() {
  std::cout << "Running test_inference..." << std::endl;

  ts::Context ctx;
  ctx = ctx.add_bind(ts::new_const_bind("i32", ir::new_type_kind()));
  ctx = ctx.add_bind(ts::new_const_bind("bool", ir::new_type_kind()));
  ctx = ctx.add_bind(ts::new_const_bind("String", ir::new_type_kind()));

  // (+): (i32, i32) -> i32
  auto plus_type = ir::new_function_type(
      ir::new_tuple_type({ir::new_name_type("i32"), ir::new_name_type("i32")}),
      ir::new_name_type("i32"));
  ctx = ctx.add_bind(ts::new_term_decl_bind("+", plus_type));

  // Test inferring function:
  // fn add_one(x: i32) -> i32 { x + 1 }
  ir::IrFunction fn;
  fn.id = "add_one";
  fn.args = {ir::FunctionArg{"x", ir::new_name_type("i32")}};
  fn.ret_type = ir::new_name_type("i32");

  // Body: Block { (+)(x, 1) }
  auto add_call = ir::new_app_term(
      ir::new_var_term("+"),
      ir::new_tuple_term({ir::new_var_term("x"), ir::new_const_int_term(1)}));
  fn.body = ir::new_block_term({add_call});

  ts::Inferencer inferencer(ctx);
  ts::Context updated_ctx;
  assert(inferencer.infer_function(&fn, updated_ctx));

  // Check solved types
  assert(fn.body.type.has_value());
  assert(ir::equals_type(*fn.body.type, ir::new_name_type("i32")));

  std::cout << "test_inference PASSED!" << std::endl;
}

void test_typechecker() {
  std::cout << "Running test_typechecker..." << std::endl;

  ts::Context ctx;
  ctx = ctx.add_bind(ts::new_const_bind("i32", ir::new_type_kind()));
  ctx = ctx.add_bind(ts::new_const_bind("bool", ir::new_type_kind()));

  // (+): (i32, i32) -> i32
  auto plus_type = ir::new_function_type(
      ir::new_tuple_type({ir::new_name_type("i32"), ir::new_name_type("i32")}),
      ir::new_name_type("i32"));
  ctx = ctx.add_bind(ts::new_term_decl_bind("+", plus_type));

  // Inherent method String::len
  ctx = ctx.add_bind(ts::new_const_bind("String", ir::new_type_kind()));
  ctx = ctx.add_bind(ts::new_term_decl_bind("String::len", ir::new_function_type(ir::new_name_type("String"), ir::new_name_type("i32"))));

  // Function with let statement and typecheck
  ir::IrFunction fn;
  fn.id = "test_fn";
  fn.args = {ir::FunctionArg{"s", ir::new_name_type("String")}};
  fn.ret_type = ir::new_name_type("i32");

  // Body: Block { let l: i32 = s.len; l }
  auto proj = ir::new_projection_term(ir::new_var_term("s"), "len");
  auto let_stmt = ir::new_let_term("l", ir::new_name_type("i32"), proj);
  fn.body = ir::new_block_term({let_stmt, ir::new_var_term("l")});

  ts::Typechecker tc(ctx);
  ts::Context updated_ctx;
  assert(tc.infer_function(&fn, updated_ctx));

  // Verify typechecking passes
  ts::Typechecker tc2(ctx);
  ts::Context out_ctx;
  assert(tc2.typecheck_function(&fn, out_ctx));

  std::cout << "test_typechecker PASSED!" << std::endl;
}

int main() {
  std::cout << "========================================" << std::endl;
  std::cout << "Running Phase 3 Native C++ Test Suite..." << std::endl;
  std::cout << "========================================" << std::endl;

  test_subtyping();
  test_unification();
  test_inference();
  test_typechecker();

  std::cout << "========================================" << std::endl;
  std::cout << "ALL PHASE 3 TESTS PASSED SUCCESSFULLY!" << std::endl;
  std::cout << "========================================" << std::endl;
  return 0;
}
