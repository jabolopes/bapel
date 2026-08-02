#include "ts/bind.h"
#include "ts/context.h"
#include "ts/list.h"

#include <cassert>
#include <iostream>
#include <string>
#include <vector>

void test_list() {
  std::cout << "Running test_list..." << std::endl;

  ts::List<int> l;
  assert(l.empty());
  assert(l.size() == 0);

  auto l1 = l.add(1);
  auto l2 = l1.add(2);
  auto l3 = l2.add(3);
  auto l4 = l3.remove();
  auto l5 = l4.remove();
  auto l6 = l5.remove();
  auto l7 = l6.remove();

  // Test sizes and empty states
  assert(l.empty() && l.size() == 0);
  assert(!l1.empty() && l1.size() == 1);
  assert(!l2.empty() && l2.size() == 2);
  assert(!l3.empty() && l3.size() == 3);
  assert(!l4.empty() && l4.size() == 2);
  assert(!l5.empty() && l5.size() == 1);
  assert(l6.empty() && l6.size() == 0);
  assert(l7.empty() && l7.size() == 0);

  // Test values (head of stack)
  int val = 0;
  assert(!l.value(val));
  assert(l1.value(val) && val == 1);
  assert(l2.value(val) && val == 2);
  assert(l3.value(val) && val == 3);
  assert(l4.value(val) && val == 2);
  assert(l5.value(val) && val == 1);
  assert(!l6.value(val));

  // Test collect (oldest to newest)
  std::vector<int> want3 = {1, 2, 3};
  assert(l3.collect() == want3);
  std::vector<int> want2 = {1, 2};
  assert(l2.collect() == want2);
  std::vector<int> want1 = {1};
  assert(l1.collect() == want1);
  assert(l.collect().empty());

  // Test collect_reverse (newest to oldest)
  std::vector<int> want3_rev = {3, 2, 1};
  assert(l3.collect_reverse() == want3_rev);

  // Test from_vector
  auto from_v = ts::List<int>::from_vector({1, 2, 3});
  assert(from_v.collect() == want3);
  assert(from_v == l3);

  // Test map
  auto mapped = l3.map([](int x) { return x * 10; });
  std::vector<int> want_mapped = {10, 20, 30};
  assert(mapped.collect() == want_mapped);

  // Test filter
  auto filtered = l3.filter([](int x) { return x % 2 != 0; });
  std::vector<int> want_filtered = {1, 3};
  assert(filtered.collect() == want_filtered);

  // Test reverse
  auto reversed = l3.reverse();
  std::vector<int> want_reversed = {3, 2, 1};
  assert(reversed.collect() == want_reversed);

  // Test iterator
  auto it = l3.iterate();
  size_t idx = 0;
  int item = 0;
  assert(it.next(idx, item) && idx == 2 && item == 3);
  assert(it.next(idx, item) && idx == 1 && item == 2);
  assert(it.next(idx, item) && idx == 0 && item == 1);
  assert(!it.next(idx, item));

  // Test structural sharing
  auto branch_a = l1.add(10);
  auto branch_b = l1.add(20);
  std::vector<int> want_a = {1, 10};
  std::vector<int> want_b = {1, 20};
  assert(branch_a.collect() == want_a);
  assert(branch_b.collect() == want_b);
  assert(l1.collect() == want1);

  std::cout << "test_list PASSED!" << std::endl;
}

void test_bind() {
  std::cout << "Running test_bind..." << std::endl;

  // AliasBind
  auto b_alias = ts::new_alias_bind("Int", ir::new_name_type("i64"));
  assert(b_alias.is_alias());
  assert(b_alias.to_string() == "type Int = i64");

  // ConstBind
  auto b_const = ts::new_const_bind("i64", ir::new_type_kind());
  assert(b_const.is_const());
  assert(b_const.to_string() == "type i64 :: ∗");

  // ScopeBind
  auto b_scope = ts::new_scope_bind(1);
  assert(b_scope.is_scope());
  assert(b_scope.to_string() == "scope 1");

  // TermDeclBind
  auto b_term_decl = ts::new_term_decl_bind("x", ir::new_name_type("i32"));
  assert(b_term_decl.is_term_decl());
  assert(b_term_decl.to_string() == "x: i32");

  // TermDefBind
  auto b_term_def = ts::new_term_def_bind("x", ir::new_name_type("i32"));
  assert(b_term_def.is_term_def());
  assert(b_term_def.to_string() == "let x: i32");

  // TypeParamBind
  auto b_tp = ts::new_type_param_bind("a", ir::new_type_kind(), {ir::new_name_type("Size")});
  assert(b_tp.is_type_param());
  assert(b_tp.to_string() == "type 'a: Size");

  // TraitBind
  auto b_trait = ts::new_trait_bind("Size", {}, {});
  assert(b_trait.is_trait());
  assert(b_trait.to_string() == "trait Size");

  // TraitImplBind
  auto b_trait_impl = ts::new_trait_impl_bind({}, ir::new_name_type("Size"), ir::new_name_type("String"));
  assert(b_trait_impl.is_trait_impl());
  assert(b_trait_impl.to_string() == "impl Size for String");

  // ExistVarBind
  auto b_evar = ts::new_exist_var_bind(10);
  assert(b_evar.is_exist_var());
  assert(b_evar.to_string() == "^10");

  // SolvedExistVarBind
  auto b_solved = ts::new_solved_exist_var_bind(10, ir::new_name_type("String"));
  assert(b_solved.is_solved_exist_var());
  assert(b_solved.to_string() == "^10 = String");

  // MarkerBind
  auto b_marker = ts::new_marker_bind(10);
  assert(b_marker.is_marker());
  assert(b_marker.to_string() == "|> ^10");

  // TermVarBind
  auto b_term_var = ts::new_term_var_bind("y", ir::new_name_type("bool"), true);
  assert(b_term_var.is_term_var());
  assert(b_term_var.to_string() == "let y: bool");

  // DeclBind
  auto b_decl = ts::new_decl_bind(ir::new_term_decl("my_func", ir::new_name_type("i32")));
  assert(b_decl.is_decl());

  // Test equality
  auto b_alias_dup = ts::new_alias_bind("Int", ir::new_name_type("i64"));
  assert(b_alias == b_alias_dup);
  assert(b_alias != b_const);

  std::cout << "test_bind PASSED!" << std::endl;
}

void test_context() {
  std::cout << "Running test_context..." << std::endl;

  ts::Context ctx = ts::new_context();
  assert(ctx.empty());
  assert(ctx.size() == 0);

  // Add Const & Alias
  ctx = ctx.add_bind(ts::new_const_bind("i32", ir::new_type_kind()));
  ctx = ctx.add_bind(ts::new_const_bind("String", ir::new_type_kind()));
  ctx = ctx.add_bind(ts::new_alias_bind("MyInt", ir::new_name_type("i32")));

  assert(ctx.size() == 3);
  assert(ctx.contains_const_bind("i32"));
  assert(ctx.contains_const_bind("String"));
  assert(ctx.contains_alias_bind("MyInt"));
  assert(!ctx.contains_const_bind("Unknown"));

  // Check string representation
  assert(ctx.to_string() == "type i32 :: ∗, type String :: ∗, type MyInt = i32");

  // Enter Scope and Term Decl / Def
  ctx = ctx.enter_scope();
  assert(ctx.size() == 4);

  ctx = ctx.add_bind(ts::new_term_decl_bind("x", ir::new_name_type("i32")));
  assert(ctx.contains_term_decl_bind_in_scope("x"));
  assert(!ctx.contains_term_def_bind_in_scope("x"));

  ctx = ctx.add_bind(ts::new_term_def_bind("y", ir::new_name_type("String")));
  assert(ctx.contains_term_def_bind_in_scope("y"));

  // Enter function
  std::vector<ir::TypeParam> tps = {ir::TypeParam{"T", ir::new_type_kind(), {}}};
  std::vector<ir::FunctionArg> args = {ir::FunctionArg{"arg0", ir::new_var_type("T")}};
  ts::Context fn_ctx = ctx.enter_function(tps, args);

  assert(fn_ctx.contains_type_param_bind_in_scope("T"));
  assert(fn_ctx.contains_term_def_bind_in_scope("arg0"));
  assert(fn_ctx.contains_term_decl_bind_in_scope("x") == false); // x was in outer scope

  // Pop
  auto [popped, popped_ctx] = fn_ctx.pop();
  assert(popped.is_term_def());

  // Test Symbol Addition: Trait with methods
  ir::IrSignature size_sig;
  size_sig.id = "size";
  size_sig.args = {ir::FunctionArg{"self", ir::new_name_type("Self")}};
  size_sig.ret_type = ir::new_name_type("i32");

  ir::IrDecl trait_decl;
  trait_decl.case_val = ir::IrDeclCase::TraitDecl;
  trait_decl.trait = std::make_shared<ir::TraitDeclData>();
  trait_decl.trait->id = "Size";
  trait_decl.trait->methods = {size_sig};

  ts::Context trait_ctx = ctx.add_symbol(trait_decl);
  assert(trait_ctx.contains_trait_bind("Size"));
  ts::Binding method_bind;
  assert(trait_ctx.lookup_term_decl_or_def_bind("Size::size", method_bind));
  ir::IrType term_var_t;
  assert(trait_ctx.lookup_term_var("Size::size", term_var_t));

  // Test Decl Lookup
  ts::Context decl_ctx = ctx.add_bind(ts::new_decl_bind(ir::new_term_decl("my_global", ir::new_name_type("i32"))));
  ir::IrDecl found_decl;
  assert(decl_ctx.lookup_decl("my_global", found_decl));
  assert(found_decl.id() == "my_global");

  // Test Inherent Method Lookup
  ts::Context method_ctx = ctx;
  method_ctx = method_ctx.add_bind(ts::new_term_decl_bind("String::len", ir::new_function_type(ir::new_name_type("String"), ir::new_name_type("i32"))));
  std::string found_method;
  ir::IrType found_type;
  assert(method_ctx.lookup_method(ir::new_name_type("String"), "len", found_method, found_type));
  assert(found_method == "String::len");

  // Test Trait Impl Method Lookup
  method_ctx = method_ctx.add_bind(ts::new_trait_bind("Printable", {}, {}));
  method_ctx = method_ctx.add_bind(ts::new_term_decl_bind("Printable::print", ir::new_function_type(ir::new_name_type("Self"), ir::new_name_type("()"))));
  method_ctx = method_ctx.add_bind(ts::new_trait_impl_bind({}, ir::new_name_type("Printable"), ir::new_name_type("String")));
  assert(method_ctx.lookup_method(ir::new_name_type("String"), "print", found_method, found_type));
  assert(found_method == "Printable::print");

  // Test Marker and Existential Variable Scoping
  ts::Context evar_ctx = ctx;
  evar_ctx = evar_ctx.add_exist_var(1);
  evar_ctx = evar_ctx.add_marker(1);
  evar_ctx = evar_ctx.add_exist_var(2);
  evar_ctx = evar_ctx.add_solved_exist_var(2, ir::new_name_type("i32"));
  assert(evar_ctx.contains_exist_var(1));
  assert(evar_ctx.contains_exist_var(2));

  // Split at marker
  auto [before_marker, after_marker] = evar_ctx.split_at_marker(1);
  assert(after_marker.size() == 2); // ^2 and ^2 = i32
  assert(!before_marker.contains_exist_var(2));
  assert(before_marker.contains_exist_var(1));

  // Drop marker
  ts::Context dropped = evar_ctx.drop_marker(1);
  assert(!dropped.contains_exist_var(2));
  assert(dropped.contains_exist_var(1));

  // Test Existential Variable Substitution
  ts::Context solve_ctx = ctx;
  solve_ctx = solve_ctx.add_solved_exist_var(1, ir::new_exist_var_type(2));
  solve_ctx = solve_ctx.add_solved_exist_var(2, ir::new_name_type("String"));

  ir::IrType input_type = ir::new_function_type(ir::new_exist_var_type(1), ir::new_name_type("i32"));
  ir::IrType substituted = solve_ctx.substitute_exist_vars(input_type);
  ir::IrType expected_type = ir::new_function_type(ir::new_name_type("String"), ir::new_name_type("i32"));
  assert(ir::equals_type(substituted, expected_type));

  // Test Fresh Variable Generation
  ir::IrType fresh_v = ctx.gen_fresh_var_type();
  assert(fresh_v.is(ir::IrTypeCase::VarType));
  assert(fresh_v.var == "a");

  ir::IrType fresh_e = ctx.gen_fresh_exist_var();
  assert(fresh_e.is(ir::IrTypeCase::ExistVarType));

  // Test Well-formedness Rejection of Duplicates
  bool threw = false;
  try {
    ctx.add_bind(ts::new_alias_bind("i32", ir::new_name_type("i64"))); // "i32" is already defined as Const
  } catch (const std::runtime_error&) {
    threw = true;
  }
  assert(threw);

  std::cout << "test_context PASSED!" << std::endl;
}

int main() {
  std::cout << "========================================" << std::endl;
  std::cout << "Running Phase 1 Native C++ Test Suite..." << std::endl;
  std::cout << "========================================" << std::endl;

  test_list();
  test_bind();
  test_context();

  std::cout << "========================================" << std::endl;
  std::cout << "ALL PHASE 1 TESTS PASSED SUCCESSFULLY!" << std::endl;
  std::cout << "========================================" << std::endl;
  return 0;
}
