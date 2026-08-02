#include "bin/ir_base.h"
#include "bin/ir_decl.h"
#include "bin/ir_function.h"
#include "bin/ir_term.h"
#include "bin/ir_type.h"
#include "bin/ir_unit.h"
#include "comp/desugar.h"
#include "comp/module_finder.h"
#include "comp/querier.h"
#include "comp/resolver.h"

#include <cassert>
#include <iostream>
#include <string>
#include <vector>

void test_desugar() {
  std::cout << "Running test_desugar..." << std::endl;

  // 1. Desugar binop a + b -> (+) (a, b)
  auto a = ir::new_var_term("a");
  auto b = ir::new_var_term("b");
  auto add_term = comp::desugar_binop_term("+", a, b);
  assert(add_term.is(ir::IrTermCase::AppTermTerm));
  assert(add_term.app_term->fun->is(ir::IrTermCase::VarTerm));
  assert(add_term.app_term->fun->var_data->id == "+");
  assert(add_term.app_term->arg->is(ir::IrTermCase::TupleTerm));
  assert(add_term.app_term->arg->tuple_data->elems.size() == 2);

  // 2. Desugar if condition -> match cond { True => t, False => e }
  auto cond = ir::new_var_term("cond");
  auto then_b = ir::new_const_int_term(1);
  auto else_b = ir::new_const_int_term(0);
  auto if_term = comp::desugar_if_term(cond, then_b, else_b);
  assert(if_term.is(ir::IrTermCase::MatchTerm));
  assert(if_term.match_data->arms.size() == 2);
  assert(if_term.match_data->arms[0].tag == "True");
  assert(if_term.match_data->arms[1].tag == "False");

  // 3. Desugar short-circuit boolean &&
  auto c1 = ir::new_var_term("x");
  auto c2 = ir::new_var_term("y");
  auto and_term = comp::desugar_binop_term("&&", c1, c2);
  assert(and_term.is(ir::IrTermCase::MatchTerm));

  std::cout << "test_desugar PASSED!" << std::endl;
}

void test_module_finder() {
  std::cout << "Running test_module_finder..." << std::endl;

  comp::ModuleFinder finder;
  finder.add_module("core.io", "lib");
  finder.add_prefix("std", "std_lib");

  // Lookup by exact name
  auto core_io_fn = finder.base_source_filename(ir::new_module_id("core.io", {}));
  assert(core_io_fn.value == "lib/core/io.bpl");

  // Lookup by prefix
  auto std_net_fn = finder.base_source_filename(ir::new_module_id("std.net.http", {}));
  assert(std_net_fn.value == "std_lib/std/net/http.bpl");

  // Impl source filename
  auto base_fn = ir::new_filename("lib/core/io.bpl", {});
  auto rel_impl = ir::new_filename("io_impl.bpl", {});
  auto impl_fn = finder.impl_source_filename(base_fn, rel_impl);
  assert(impl_fn.value == "lib/core/io_impl.bpl");

  std::cout << "test_module_finder PASSED!" << std::endl;
}

void test_querier_and_resolver() {
  std::cout << "Running test_querier_and_resolver..." << std::endl;

  comp::ModuleFinder finder;
  comp::Querier querier(finder);

  // Register an exported module "math"
  comp::ModuleQuery math_mod;
  math_mod.decls.push_back(ir::new_name_decl("Num", ir::new_type_kind(), true));
  math_mod.decls.push_back(ir::new_term_decl("math::pi", ir::new_name_type("Num"), true));
  math_mod.decls.push_back(ir::new_term_decl("math::private_helper", ir::new_name_type("Num"), false));
  querier.register_module("math", math_mod);

  // Unit under test importing "math"
  ir::IrUnit unit;
  unit.module_id = ir::new_module_id("main", {});
  unit.imports.push_back(ir::new_import(ir::new_module_id("math", {})));

  // Define local type MyInt
  unit.decls.push_back(ir::new_alias_decl("MyInt", ir::new_type_kind(), ir::new_name_type("i32"), true));

  // Trait coherence test: Valid local impl
  unit.trait_impls.push_back(ir::new_trait_impl({}, ir::new_name_type("Show"), ir::new_name_type("MyInt")));

  std::string err;
  comp::Resolver resolver(querier);
  assert(resolver.resolve_source_unit(unit, err));

  // Verify only exported declarations from "math" are imported
  assert(unit.import_decls.size() == 2);
  assert(unit.import_decls[0].id() == "Num");
  assert(unit.import_decls[1].id() == "math::pi");

  // Verify Trait Coherence rejection for foreign types:
  // Try implementing a trait for foreign type "Num" which is not defined in this unit
  ir::IrUnit bad_unit;
  bad_unit.module_id = ir::new_module_id("main2", {});
  bad_unit.trait_impls.push_back(ir::new_trait_impl({}, ir::new_name_type("Show"), ir::new_name_type("Num")));

  comp::Resolver resolver2(querier);
  assert(!resolver2.resolve_source_unit(bad_unit, err));
  assert(err.find("trait coherence violation") != std::string::npos);

  std::cout << "test_querier_and_resolver PASSED!" << std::endl;
}

int main() {
  std::cout << "========================================" << std::endl;
  std::cout << "Running Phase 4 Native C++ Test Suite..." << std::endl;
  std::cout << "========================================" << std::endl;

  test_desugar();
  test_module_finder();
  test_querier_and_resolver();

  std::cout << "========================================" << std::endl;
  std::cout << "ALL PHASE 4 TESTS PASSED SUCCESSFULLY!" << std::endl;
  std::cout << "========================================" << std::endl;
  return 0;
}
