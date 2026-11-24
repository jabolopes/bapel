#pragma once

#include <cstddef>
#include <utility>
#include <variant>

#define BAPEL_UNWRAP(...) __VA_ARGS__
#define BAPEL_CALL(M, ...) M(__VA_ARGS__)
#define BAPEL_COMMA() ,
#define BAPEL_EMPTY()

#define BAPEL_V_APPLY(M, Name, Idx, Pair) BAPEL_CALL(M, Name, Idx, BAPEL_UNWRAP Pair)

#define BAPEL_V_TYPE(Name, Idx, Tag, TypeTuple) BAPEL_UNWRAP TypeTuple
#define BAPEL_V_ENUM(Name, Idx, Tag, TypeTuple) Tag = Idx
#define BAPEL_V_METHODS(Name, Idx, Tag, TypeTuple)                     \
  static Name mk_##Tag(BAPEL_UNWRAP TypeTuple v = {}) {                \
    return Name(std::in_place_index<Idx>, std::move(v));               \
  }                                                                    \
  BAPEL_UNWRAP TypeTuple& as_##Tag() { return std::get<Idx>(*this); }  \
  const BAPEL_UNWRAP TypeTuple& as_##Tag() const {                     \
    return std::get<Idx>(*this);                                       \
  }

#define BAPEL_FE_1(M, S, N, I, x) BAPEL_V_APPLY(M, N, I, x)
#define BAPEL_FE_2(M, S, N, I, x, ...) BAPEL_V_APPLY(M, N, I, x) S() BAPEL_FE_1(M, S, N, I + 1, __VA_ARGS__)
#define BAPEL_FE_3(M, S, N, I, x, ...) BAPEL_V_APPLY(M, N, I, x) S() BAPEL_FE_2(M, S, N, I + 1, __VA_ARGS__)
#define BAPEL_FE_4(M, S, N, I, x, ...) BAPEL_V_APPLY(M, N, I, x) S() BAPEL_FE_3(M, S, N, I + 1, __VA_ARGS__)
#define BAPEL_FE_5(M, S, N, I, x, ...) BAPEL_V_APPLY(M, N, I, x) S() BAPEL_FE_4(M, S, N, I + 1, __VA_ARGS__)
#define BAPEL_FE_6(M, S, N, I, x, ...) BAPEL_V_APPLY(M, N, I, x) S() BAPEL_FE_5(M, S, N, I + 1, __VA_ARGS__)
#define BAPEL_FE_7(M, S, N, I, x, ...) BAPEL_V_APPLY(M, N, I, x) S() BAPEL_FE_6(M, S, N, I + 1, __VA_ARGS__)
#define BAPEL_FE_8(M, S, N, I, x, ...) BAPEL_V_APPLY(M, N, I, x) S() BAPEL_FE_7(M, S, N, I + 1, __VA_ARGS__)
#define BAPEL_FE_9(M, S, N, I, x, ...) BAPEL_V_APPLY(M, N, I, x) S() BAPEL_FE_8(M, S, N, I + 1, __VA_ARGS__)
#define BAPEL_FE_10(M, S, N, I, x, ...) BAPEL_V_APPLY(M, N, I, x) S() BAPEL_FE_9(M, S, N, I + 1, __VA_ARGS__)
#define BAPEL_FE_11(M, S, N, I, x, ...) BAPEL_V_APPLY(M, N, I, x) S() BAPEL_FE_10(M, S, N, I + 1, __VA_ARGS__)
#define BAPEL_FE_12(M, S, N, I, x, ...) BAPEL_V_APPLY(M, N, I, x) S() BAPEL_FE_11(M, S, N, I + 1, __VA_ARGS__)
#define BAPEL_FE_13(M, S, N, I, x, ...) BAPEL_V_APPLY(M, N, I, x) S() BAPEL_FE_12(M, S, N, I + 1, __VA_ARGS__)
#define BAPEL_FE_14(M, S, N, I, x, ...) BAPEL_V_APPLY(M, N, I, x) S() BAPEL_FE_13(M, S, N, I + 1, __VA_ARGS__)
#define BAPEL_FE_15(M, S, N, I, x, ...) BAPEL_V_APPLY(M, N, I, x) S() BAPEL_FE_14(M, S, N, I + 1, __VA_ARGS__)
#define BAPEL_FE_16(M, S, N, I, x, ...) BAPEL_V_APPLY(M, N, I, x) S() BAPEL_FE_15(M, S, N, I + 1, __VA_ARGS__)
#define BAPEL_FE_17(M, S, N, I, x, ...) BAPEL_V_APPLY(M, N, I, x) S() BAPEL_FE_16(M, S, N, I + 1, __VA_ARGS__)
#define BAPEL_FE_18(M, S, N, I, x, ...) BAPEL_V_APPLY(M, N, I, x) S() BAPEL_FE_17(M, S, N, I + 1, __VA_ARGS__)
#define BAPEL_FE_19(M, S, N, I, x, ...) BAPEL_V_APPLY(M, N, I, x) S() BAPEL_FE_18(M, S, N, I + 1, __VA_ARGS__)
#define BAPEL_FE_20(M, S, N, I, x, ...) BAPEL_V_APPLY(M, N, I, x) S() BAPEL_FE_19(M, S, N, I + 1, __VA_ARGS__)
#define BAPEL_FE_21(M, S, N, I, x, ...) BAPEL_V_APPLY(M, N, I, x) S() BAPEL_FE_20(M, S, N, I + 1, __VA_ARGS__)
#define BAPEL_FE_22(M, S, N, I, x, ...) BAPEL_V_APPLY(M, N, I, x) S() BAPEL_FE_21(M, S, N, I + 1, __VA_ARGS__)
#define BAPEL_FE_23(M, S, N, I, x, ...) BAPEL_V_APPLY(M, N, I, x) S() BAPEL_FE_22(M, S, N, I + 1, __VA_ARGS__)
#define BAPEL_FE_24(M, S, N, I, x, ...) BAPEL_V_APPLY(M, N, I, x) S() BAPEL_FE_23(M, S, N, I + 1, __VA_ARGS__)

#define BAPEL_GET_FE(_1, _2, _3, _4, _5, _6, _7, _8, _9, _10, _11, _12, \
                     _13, _14, _15, _16, _17, _18, _19, _20, _21, _22,  \
                     _23, _24, NAME, ...) NAME

#define BAPEL_FOR_EACH(M, S, N, ...)                                      \
  BAPEL_GET_FE(__VA_ARGS__, BAPEL_FE_24, BAPEL_FE_23, BAPEL_FE_22,        \
               BAPEL_FE_21, BAPEL_FE_20, BAPEL_FE_19, BAPEL_FE_18,        \
               BAPEL_FE_17, BAPEL_FE_16, BAPEL_FE_15, BAPEL_FE_14,        \
               BAPEL_FE_13, BAPEL_FE_12, BAPEL_FE_11, BAPEL_FE_10,        \
               BAPEL_FE_9, BAPEL_FE_8, BAPEL_FE_7, BAPEL_FE_6,            \
               BAPEL_FE_5, BAPEL_FE_4, BAPEL_FE_3, BAPEL_FE_2,            \
               BAPEL_FE_1)(M, S, N, 0, __VA_ARGS__)

#define BAPEL_VARIANT(Name, ...)                                          \
  struct Name                                                             \
      : std::variant<BAPEL_FOR_EACH(BAPEL_V_TYPE, BAPEL_COMMA, Name,      \
                                    __VA_ARGS__)> {                       \
    using std::variant<BAPEL_FOR_EACH(BAPEL_V_TYPE, BAPEL_COMMA, Name,    \
                                      __VA_ARGS__)>::variant;             \
    enum class Tag : size_t {                                             \
      BAPEL_FOR_EACH(BAPEL_V_ENUM, BAPEL_COMMA, Name, __VA_ARGS__)        \
    };                                                                    \
    Tag tag() const { return static_cast<Tag>(this->index()); }           \
    bool is(Tag _bapel_tag) const { return tag() == _bapel_tag; }         \
    BAPEL_FOR_EACH(BAPEL_V_METHODS, BAPEL_EMPTY, Name, __VA_ARGS__)       \
  }
