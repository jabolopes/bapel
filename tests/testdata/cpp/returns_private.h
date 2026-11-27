#pragma once

#include "returns.h"

BAPEL_VARIANT(Choice, (a, (int64_t)), (b, (int64_t)));
int8_t return1();
int8_t return2();
int8_t return3();
std::monostate noop();
std::monostate unit_if_no_else(bool cond);
std::monostate unit_early_return_if(bool cond);
std::monostate unit_for_loop(int64_t n);
int64_t value_if_else_implicit(bool cond);
int64_t value_for_early_return(int64_t n);
int64_t match_exhaustive_value(::Choice t);
int64_t match_wildcard_value(::Choice t);
std::monostate unit_match_exhaustive(::Choice t);
