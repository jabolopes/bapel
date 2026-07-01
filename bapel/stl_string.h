#pragma once

#include <string>
#include <string_view>


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

// @bpl: pub StringImpl::empty: String -> bool
// @bpl: pub StringImpl::front: String -> i8
// @bpl: pub StringImpl::size: String -> i64
// @bpl: pub StringImpl::view: String -> StringView
// @bpl: pub StringImpl::from_view: StringView -> String
// @bpl: pub StringImpl::from_char: i8 -> String
// @bpl: pub StringImpl::concat: (String, String) -> String
// @bpl: pub StringImpl::find: (String, String, i64) -> i64
// @bpl: pub StringImpl::replace: (String, i64, i64, String) -> String
struct StringImpl {
  StringImpl() = delete;

  static inline String from_view(StringView v) { return String(v); }
  static inline String from_char(char c) { return String(1, c); }

  static inline bool empty(const String& s) { return s.empty(); }
  static inline char front(const String& s) { return s.front(); }
  static inline int64_t size(const String& s) { return s.size(); }
  static inline StringView view(const String& s) { return s; }
  static inline String concat(const String& a, const String& b) { return a + b; }

  static inline int64_t find(const String& s, const String& target, int64_t pos) {
    if (pos < 0 || pos > static_cast<int64_t>(s.size())) {
        return -1;
    }
    size_t res = s.find(target, pos);
    if (res == std::string::npos) {
        return -1;
    }
    return static_cast<int64_t>(res);
  }

  static inline String replace(const String& s, int64_t pos, int64_t count, const String& to) {
    String res = s;
    // TODO: This is only needed because C++ uses size_t (which is
    // unsigned) whereas Bapel is using int64_t.
    if (pos < 0 || pos > static_cast<int64_t>(res.size()) || count < 0) {
        return res;
    }
    res.replace(pos, count, to);
    return res;
  }
};

// @bpl: pub StringViewImpl::at: (StringView, i64) -> i8
// @bpl: pub StringViewImpl::empty: StringView -> bool
// @bpl: pub StringViewImpl::front: StringView -> i8
// @bpl: pub StringViewImpl::size: StringView -> i64
// @bpl: pub StringViewImpl::substr: (StringView, i64, i64) -> StringView
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
};
