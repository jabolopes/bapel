#pragma once

#include <charconv>
#include <cstdint>
#include <string>
#include <string_view>
#include <variant>
#include "bapel/stl_optional.h"


using String = std::string;
using StringView = std::string_view;

// @bpl: pub to_string: forall ['a] 'a -> String
template <typename T>
inline String to_string(T value) { return std::to_string(value); }

// @bpl: pub getline: forall ['s] (Ptr 's, Ptr String) -> bool
template <typename Stream>
inline bool getline(Stream* is, String* s) {
  return static_cast<bool>(std::getline(*is, *s));
}

// @bpl: pub parse_i64: StringView -> Optional i64
inline Optional<int64_t> parse_i64(StringView s) {
  if (s.empty()) {
    return std::nullopt;
  }
  int base = 10;
  if (s.size() > 2 && s[0] == '0' && (s[1] == 'x' || s[1] == 'X')) {
    base = 16;
    s.remove_prefix(2);
  }
  int64_t val = 0;
  auto res = std::from_chars(s.data(), s.data() + s.size(), val, base);
  if (res.ec != std::errc() || res.ptr != s.data() + s.size()) {
    return std::nullopt;
  }
  return std::make_optional(val);
}

// @bpl: pub StringImpl::empty: &String -> bool
// @bpl: pub StringImpl::front: &String -> i8
// @bpl: pub StringImpl::at: (&String, i64) -> i8
// @bpl: pub StringImpl::size: &String -> i64
// @bpl: pub StringImpl::view: &String -> StringView
// @bpl: pub StringImpl::from_view: StringView -> String
// @bpl: pub StringImpl::from_char: i8 -> String
// @bpl: pub StringImpl::concat: (&String, &String) -> String
// @bpl: pub StringImpl::append: (&String, StringView) -> ()
// @bpl: pub StringImpl::push_back: (&String, i8) -> ()
// @bpl: pub StringImpl::reserve: (&String, i64) -> ()
// @bpl: pub StringImpl::substr: (&String, i64, i64) -> String
// @bpl: pub StringImpl::find: (&String, &String, i64) -> i64
// @bpl: pub StringImpl::replace: (&String, i64, i64, &String) -> String
// @bpl: pub StringImpl::starts_with: (&String, &String) -> bool
// @bpl: pub StringImpl::ends_with: (&String, &String) -> bool
// @bpl: pub StringImpl::rfind: (&String, &String) -> i64
struct StringImpl {
  StringImpl() = delete;

  static inline String from_view(StringView v) { return String(v); }
  static inline String from_char(char c) { return String(1, c); }

  static inline bool empty(const String* s) { return s->empty(); }
  static inline char front(const String* s) { return s->front(); }
  static inline char at(const String* s, int64_t i) { return s->at(i); }
  static inline int64_t size(const String* s) { return s->size(); }
  static inline StringView view(const String* s) { return *s; }
  static inline String concat(const String* a, const String* b) { return *a + *b; }

  static inline std::monostate append(String* s, StringView other) {
    s->append(other);
    return std::monostate();
  }

  static inline std::monostate push_back(String* s, char c) {
    s->push_back(c);
    return std::monostate();
  }

  static inline std::monostate reserve(String* s, int64_t cap) {
    if (cap > 0) {
      s->reserve(static_cast<size_t>(cap));
    }
    return std::monostate();
  }

  static inline String substr(const String* s, int64_t pos, int64_t len) {
    if (pos < 0 || pos > static_cast<int64_t>(s->size()) || len < 0) {
      return String();
    }
    return s->substr(static_cast<size_t>(pos), static_cast<size_t>(len));
  }

  static inline int64_t find(const String* s, const String* target, int64_t pos) {
    if (pos < 0 || pos > static_cast<int64_t>(s->size())) {
        return -1;
    }
    size_t res = s->find(*target, pos);
    if (res == std::string::npos) {
        return -1;
    }
    return static_cast<int64_t>(res);
  }

  // TODO: rfind should return an Optional i64 instead.
  static inline int64_t rfind(const String* s, const String* target) {
    size_t res = s->rfind(*target);
    if (res == std::string::npos) {
        return -1;
    }
    return static_cast<int64_t>(res);
  }

  static inline String replace(const String* s, int64_t pos, int64_t count, const String* to) {
    String res = *s;
    // TODO: This is only needed because C++ uses size_t (which is
    // unsigned) whereas Bapel is using int64_t.
    if (pos < 0 || pos > static_cast<int64_t>(res.size()) || count < 0) {
        return res;
    }
    res.replace(pos, count, *to);
    return res;
  }

  static inline bool starts_with(const String* s, const String* prefix) {
    if (s->size() < prefix->size()) return false;
    return StringView(*s).substr(0, prefix->size()) == StringView(*prefix);
  }

  static inline bool ends_with(const String* s, const String* suffix) {
    if (s->size() < suffix->size()) return false;
    return StringView(*s).substr(s->size() - suffix->size(), suffix->size()) == StringView(*suffix);
  }
};

// @bpl: pub StringViewImpl::at: (StringView, i64) -> i8
// @bpl: pub StringViewImpl::empty: StringView -> bool
// @bpl: pub StringViewImpl::front: StringView -> i8
// @bpl: pub StringViewImpl::size: StringView -> i64
// @bpl: pub StringViewImpl::substr: (StringView, i64, i64) -> StringView
// @bpl: pub StringViewImpl::starts_with: (StringView, StringView) -> bool
// @bpl: pub StringViewImpl::ends_with: (StringView, StringView) -> bool
// @bpl: pub StringViewImpl::remove_prefix: (&StringView, i64) -> ()
// @bpl: pub StringViewImpl::remove_suffix: (&StringView, i64) -> ()
struct StringViewImpl {
  StringViewImpl() = delete;

  static inline char at(StringView s, int64_t i) { return s.at(i); }
  static inline bool empty(StringView s) { return s.empty(); }
  static inline char front(StringView s) { return s.front(); }
  static inline int64_t size(StringView s) { return s.size(); }

  static inline StringView substr(StringView s, int64_t pos, int64_t size) {
    if (pos > s.size()) {
      return StringView();
    }
    return s.substr(pos, size);
  }

  static inline bool starts_with(StringView s, StringView prefix) {
    if (s.size() < prefix.size()) return false;
    return s.substr(0, prefix.size()) == prefix;
  }

  static inline bool ends_with(StringView s, StringView suffix) {
    if (s.size() < suffix.size()) return false;
    return s.substr(s.size() - suffix.size(), suffix.size()) == suffix;
  }

  // Check n > 0 to guard against negative int64_t inputs implicitly converting
  // to massive unsigned size_t values, which would cause undefined behavior.
  static inline std::monostate remove_prefix(StringView* s, int64_t n) {
    if (n > 0 && n <= static_cast<int64_t>(s->size())) {
      s->remove_prefix(n);
    } else if (n > static_cast<int64_t>(s->size())) {
      *s = StringView();
    }
    return std::monostate();
  }

  // Check n > 0 to guard against negative int64_t inputs implicitly converting
  // to massive unsigned size_t values, which would cause undefined behavior.
  static inline std::monostate remove_suffix(StringView* s, int64_t n) {
    if (n > 0 && n <= static_cast<int64_t>(s->size())) {
      s->remove_suffix(n);
    } else if (n > static_cast<int64_t>(s->size())) {
      *s = StringView();
    }
    return std::monostate();
  }
};
